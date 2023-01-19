package templ

import (
	"darkbot/utils"
	"darkbot/viewer/apis"
	"fmt"
	"strings"
	"time"
)

type MsgAction int

const (
	ActSend MsgAction = iota
	ActEdit
)

type TemplateShared struct {
	MessageID string
	Content   string
	Header    string
	*apis.API
}

func CheckTooLongMsgErr(err error, api *apis.API, header string, action MsgAction, MessageID string) {
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

func ChannelCheckWarn(err error, channelID string, msg string) {
	utils.CheckWarn(err, "channelID=", channelID, msg)
}

func (v *TemplateShared) Send() {
	if v.Content == "" && v.MessageID != "" {
		v.Discorder.DeleteMessage(v.ChannelID, v.MessageID)
	}

	if v.Content == "" {
		return
	}

	var err error
	if v.MessageID == "" {
		err = v.Discorder.SengMessage(v.ChannelID, v.Content)
		ChannelCheckWarn(err, v.ChannelID, "unable to send msg")
		CheckTooLongMsgErr(err, v.API, v.Header, ActSend, "")

	} else {
		err = v.Discorder.EditMessage(v.ChannelID, v.MessageID, v.Content)
		ChannelCheckWarn(err, v.ChannelID, "unable to edit msg")
		CheckTooLongMsgErr(err, v.API, v.Header, ActEdit, v.MessageID)
	}
}
