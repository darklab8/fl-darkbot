package viewer_msg

import (
	"fmt"
	"strings"
	"time"

	"github.com/darklab8/fl-darkbot/app/settings/logus"
	"github.com/darklab8/fl-darkbot/app/settings/types"
	"github.com/darklab8/fl-darkbot/app/viewer/apis"

	"github.com/darklab8/go-utils/utils/timeit"
)

type MsgShared struct {
	viewID        types.ViewID // ID against msg resend
	viewHeader    types.ViewHeader
	viewBeginning types.ViewBeginning
	viewEnd       types.ViewEnd
	viewType      viewType
}

type viewType int

const (
	MsgAlert viewType = iota
	MsgTable
)

func NewTableMsg(
	viewID types.ViewID,
	viewHeader types.ViewHeader,
	viewBeginning types.ViewBeginning,
	viewEnd types.ViewEnd,
) *MsgShared {
	return &MsgShared{
		viewID:        viewID,
		viewHeader:    viewHeader,
		viewBeginning: viewBeginning,
		viewEnd:       viewEnd,
		viewType:      MsgTable,
	}
}

func NewAlertMsg(
	viewID types.ViewID,
) *MsgShared {
	return &MsgShared{
		viewID:   viewID,
		viewType: MsgAlert,
	}
}

func (v *MsgShared) SetHeader(header types.ViewHeader) { v.viewHeader = header }

func (m *MsgShared) GetTimestamp() types.ViewTimeStamp {
	return types.ViewTimeStamp(fmt.Sprintf(" (last updated: %s)", time.Now().String()))
}
func (m *MsgShared) GetViewID() types.ViewID { return m.viewID }

func (m *MsgShared) LenShared() int {
	size := len(m.viewID)
	size += len(m.GetTimestamp())
	size += len(m.viewHeader)
	size += len(m.viewBeginning)
	size += len(m.viewEnd)
	return size
}

type Msg struct {
	*MsgShared
	messageID        types.DiscordMessageID
	records          []*types.ViewRecord
	viewEnumeratedID types.ViewEnumeratedID
	channelID        types.DiscordChannelID
}

func NewMsg(msgShared *MsgShared, number int, channelID types.DiscordChannelID) *Msg {
	m := &Msg{
		MsgShared:        msgShared,
		viewEnumeratedID: types.ViewEnumeratedID(fmt.Sprintf("%s-%d", msgShared.viewID, number)),
		channelID:        channelID,
	}
	return m
}

func (m *Msg) GetMessageID() types.DiscordMessageID { return m.messageID }

func (m *Msg) GetViewEnumeratedID() types.ViewEnumeratedID { return m.viewEnumeratedID }

func (m *Msg) HasRecords() bool { return len(m.records) > 0 }

func (m *Msg) SetMessageID(messagID types.DiscordMessageID) { m.messageID = messagID }

func (m *Msg) AppendRecordToMsg(record *types.ViewRecord) { m.records = append(m.records, record) }

func (m *Msg) Len() int {
	size := m.LenShared()
	for _, record := range m.records {
		size += len(*record)
	}
	return size
}

func (v *Msg) Render() string {
	var content strings.Builder

	if v.viewType == MsgTable {
		content.WriteString(string(v.viewEnumeratedID) + string(v.GetTimestamp()) + "\n")
		content.WriteString(string(v.viewHeader))
		content.WriteString(string(v.viewBeginning))
		for _, record := range v.records {
			content.WriteString(string(*record))
		}
		content.WriteString(string(v.viewEnd))
	} else {
		// Mobile friendly way to render alert
		content.WriteString(string(v.viewHeader))
		content.WriteString("\n" + string(v.viewEnumeratedID) + string(v.GetTimestamp()))
		// content.WriteString(string(v.viewBeginning))
		// for _, record := range v.records {
		// 	content.WriteString(string(*record))
		// }
		// content.WriteString(string(v.viewEnd))
	}
	return content.String()
}

func (v *Msg) Send(api *apis.API) {
	timeit.NewTimerMF(fmt.Sprintf("Msg.Send() msg=%v", v), func() {

		timeit.NewTimerMFL(fmt.Sprintf("msg.Send DeleteMessage. Msg=%v", *v), func() {
			if len(v.records) == 0 && v.messageID != "" {
				api.Discorder.DeleteMessage(v.channelID, v.messageID)
			}
		}, logus.ChannelID(v.channelID))

		if len(v.records) == 0 {
			return
		}

		time_render := timeit.NewTimer(fmt.Sprintf("msg.Render, msg=%v", v))
		content := v.Render()
		time_render.Close()

		timeit.NewTimerMF(fmt.Sprintf("Msg.Send().SecondSection msg=%v", v), func() {
			var err error
			if v.messageID == "" {
				timeit.NewTimerMF(fmt.Sprintf("Msg.Send().SendMessage + ChecTooLongMsg msg=%v", v), func() {
					err = api.Discorder.SengMessage(v.channelID, content)
					logus.Log.CheckWarn(err, "unable to send msg", logus.ChannelID(v.channelID))
					CheckTooLongMsgErr(err, api, v.channelID, v.viewEnumeratedID, ActSend, "")
				})
			} else {
				timeit.NewTimerMF(fmt.Sprintf("Msg.Send().EditMessage + ChecTooLongMsg msg=%v", v), func() {
					timeit.NewTimerMF(fmt.Sprintf("Msg.Send().EditMessage.only msg=%v", v), func() {
						err = api.Discorder.EditMessage(v.channelID, v.messageID, content)
					})
					timeit.NewTimerMF(fmt.Sprintf("Msg.Send().EditMessage.CheckTooLongMsgErr msg=%v", v), func() {
						logus.Log.CheckWarn(err, "unable to edit msg", logus.ChannelID(v.channelID))
						CheckTooLongMsgErr(err, api, v.channelID, v.viewEnumeratedID, ActEdit, v.messageID)
					})
				})
			}
		})
	})
}

func CheckTooLongMsgErr(err error, api *apis.API, channeID types.DiscordChannelID, header types.ViewEnumeratedID, action MsgAction, MessageID types.DiscordMessageID) {
	timeit.NewTimerMFL(fmt.Sprintf("CheckTooLongMsgErr. header=%s, messageID=%s, action=%d", header, MessageID, action), func() {

		if err == nil {
			return
		}

		if !strings.Contains(err.Error(), "BASE_TYPE_MAX_LENGTH") &&
			!strings.Contains(err.Error(), "or fewer in length") {
			return
		}

		msg := fmt.Sprintf("%s, %s, %s", header, time.Now(), err)

		switch action {
		case ActSend:
			api.Discorder.SengMessage(channeID, msg)
		case ActEdit:
			api.Discorder.EditMessage(channeID, MessageID, msg)
		}
	}, logus.ChannelID(channeID))
}

const (
	ActSend MsgAction = iota
	ActEdit
)

type MsgAction int
