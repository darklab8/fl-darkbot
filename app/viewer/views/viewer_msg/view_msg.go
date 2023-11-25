package viewer_msg

import (
	"darkbot/app/settings/logus"
	"darkbot/app/settings/types"
	"darkbot/app/settings/utils"
	"darkbot/app/viewer/apis"
	"fmt"
	"strings"
	"time"
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
		for _, record := range v.records {
			content.WriteString(string(*record))
		}
		content.WriteString("\n" + string(v.viewEnumeratedID) + string(v.GetTimestamp()))
	}
	return content.String()
}

func (v *Msg) Send(api *apis.API) {
	utils.TimeMeasure(func() {

		utils.TimeMeasure(func() {
			if len(v.records) == 0 && v.messageID != "" {
				api.Discorder.DeleteMessage(v.channelID, v.messageID)
			}
		}, fmt.Sprintf("msg.Send DeleteMessage. Msg=%v", *v), logus.ChannelID(v.channelID))

		if len(v.records) == 0 {
			return
		}

		time_render := utils.NewTimeMeasure(fmt.Sprintf("msg.Render, msg=%v", v))
		content := v.Render()
		time_render.Close()

		utils.TimeMeasure(func() {
			var err error
			if v.messageID == "" {
				utils.TimeMeasure(func() {
					err = api.Discorder.SengMessage(v.channelID, content)
					logus.CheckWarn(err, "unable to send msg", logus.ChannelID(v.channelID))
					CheckTooLongMsgErr(err, api, v.channelID, v.viewEnumeratedID, ActSend, "")
				}, fmt.Sprintf("Msg.Send().SendMessage + ChecTooLongMsg msg=%v", v))
			} else {
				utils.TimeMeasure(func() {
					utils.TimeMeasure(func() {
						err = api.Discorder.EditMessage(v.channelID, v.messageID, content)
					}, fmt.Sprintf("Msg.Send().EditMessage.only msg=%v", v))
					utils.TimeMeasure(func() {
						logus.CheckWarn(err, "unable to edit msg", logus.ChannelID(v.channelID))
						CheckTooLongMsgErr(err, api, v.channelID, v.viewEnumeratedID, ActEdit, v.messageID)
					}, fmt.Sprintf("Msg.Send().EditMessage.CheckTooLongMsgErr msg=%v", v))
				}, fmt.Sprintf("Msg.Send().EditMessage + ChecTooLongMsg msg=%v", v))
			}
		}, fmt.Sprintf("Msg.Send().SecondSection msg=%v", v))
	}, fmt.Sprintf("Msg.Send() msg=%v", v))
}

func CheckTooLongMsgErr(err error, api *apis.API, channeID types.DiscordChannelID, header types.ViewEnumeratedID, action MsgAction, MessageID types.DiscordMessageID) {
	utils.TimeMeasure(func() {

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
	}, fmt.Sprintf("CheckTooLongMsgErr. header=%s, messageID=%s, action=%d", header, MessageID, action), logus.ChannelID(channeID))
}

const (
	ActSend MsgAction = iota
	ActEdit
)

type MsgAction int
