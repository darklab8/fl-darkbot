package views

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
	ViewID        types.ViewID
	ViewBeginning types.ViewBeginning
	ViewEnd       types.ViewEnd
}
type Msg struct {
	MsgShared
	MessageID        types.DiscordMessageID
	Records          []*types.ViewRecord
	ViewEnumeratedID types.ViewEnumeratedID
	channelID        types.DiscordChannelID
}

func NewMsg(msgShared MsgShared, number int, channelID types.DiscordChannelID) *Msg {
	m := &Msg{MsgShared: msgShared}
	m.ViewEnumeratedID = types.ViewEnumeratedID(fmt.Sprintf("%s-%d", m.ViewID, number))
	m.channelID = channelID
	return m
}

func (m Msg) Len() int {
	size := len(m.ViewID)
	size += len(m.ViewBeginning)
	size += len(m.ViewEnd)

	for _, record := range m.Records {
		size += len(*record)
	}
	return size
}

func (v *Msg) Render() string {
	var content strings.Builder
	content.WriteString(string(v.ViewEnumeratedID) + fmt.Sprintf(" (last updated: %s)", time.Now().String()) + "\n")
	content.WriteString(string(v.ViewBeginning))
	for _, record := range v.Records {
		content.WriteString(string(*record))
	}
	content.WriteString(string(v.ViewEnd))
	return content.String()
}

func (v *Msg) Send(api *apis.API) {
	utils.TimeMeasure(func() {

		utils.TimeMeasure(func() {
			if len(v.Records) == 0 && v.MessageID != "" {
				api.Discorder.DeleteMessage(v.channelID, v.MessageID)
			}
		}, fmt.Sprintf("msg.Send DeleteMessage. Msg=%v", *v), logus.ChannelID(v.channelID))

		if len(v.Records) == 0 {
			return
		}

		time_render := utils.NewTimeMeasure(fmt.Sprintf("msg.Render, msg=%v", v))
		content := v.Render()
		time_render.Close()

		utils.TimeMeasure(func() {
			var err error
			if v.MessageID == "" {
				utils.TimeMeasure(func() {
					err = api.Discorder.SengMessage(v.channelID, content)
					logus.CheckWarn(err, "unable to send msg", logus.ChannelID(v.channelID))
					CheckTooLongMsgErr(err, api, v.channelID, v.ViewEnumeratedID, ActSend, "")
				}, fmt.Sprintf("Msg.Send().SendMessage + ChecTooLongMsg msg=%v", v))
			} else {
				utils.TimeMeasure(func() {
					utils.TimeMeasure(func() {
						err = api.Discorder.EditMessage(v.channelID, v.MessageID, content)
					}, fmt.Sprintf("Msg.Send().EditMessage.only msg=%v", v))
					utils.TimeMeasure(func() {
						logus.CheckWarn(err, "unable to edit msg", logus.ChannelID(v.channelID))
						CheckTooLongMsgErr(err, api, v.channelID, v.ViewEnumeratedID, ActEdit, v.MessageID)
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
