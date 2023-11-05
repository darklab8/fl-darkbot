package views

import (
	"darkbot/app/settings/logus"
	"darkbot/app/settings/types"
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
}

func NewMsg(msgShared MsgShared, number int) *Msg {
	m := &Msg{MsgShared: msgShared}
	m.ViewEnumeratedID = types.ViewEnumeratedID(fmt.Sprintf("%s-%d", m.ViewID, number))
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
	content.WriteString(string(v.ViewEnumeratedID) + "\n")
	content.WriteString(string(v.ViewBeginning))
	for _, record := range v.Records {
		content.WriteString(string(*record))
	}
	content.WriteString(string(v.ViewEnd))
	return content.String()
}

func (v *Msg) Send(api *apis.API) {
	if len(v.Records) == 0 && v.MessageID != "" {
		api.Discorder.DeleteMessage(api.ChannelID, v.MessageID)
	}

	if len(v.Records) == 0 {
		return
	}

	content := v.Render()

	var err error
	if v.MessageID == "" {
		err = api.Discorder.SengMessage(api.ChannelID, content)
		logus.CheckWarn(err, "unable to send msg", logus.ChannelID(api.ChannelID))
		CheckTooLongMsgErr(err, api, v.ViewEnumeratedID, ActSend, "")

	} else {
		err = api.Discorder.EditMessage(api.ChannelID, v.MessageID, content)
		logus.CheckWarn(err, "unable to edit msg", logus.ChannelID(api.ChannelID))
		CheckTooLongMsgErr(err, api, v.ViewEnumeratedID, ActEdit, v.MessageID)
	}
}

func CheckTooLongMsgErr(err error, api *apis.API, header types.ViewEnumeratedID, action MsgAction, MessageID types.DiscordMessageID) {
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
		api.Discorder.SengMessage(api.ChannelID, msg)
	case ActEdit:
		api.Discorder.EditMessage(api.ChannelID, MessageID, msg)
	}

}
