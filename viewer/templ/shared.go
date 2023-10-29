package templ

import (
	"darkbot/settings/utils/logger"
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
}

func CheckTooLongMsgErr(err error, api apis.API, header string, action MsgAction, MessageID string) {
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
	logger.CheckWarn(err, "channelID=", channelID, msg)
}

func (v *TemplateShared) Send(api apis.API) {
	if v.Content == "" && v.MessageID != "" {
		api.Discorder.DeleteMessage(api.ChannelID, v.MessageID)
	}

	if v.Content == "" {
		return
	}

	var err error
	if v.MessageID == "" {
		err = api.Discorder.SengMessage(api.ChannelID, v.Content)
		ChannelCheckWarn(err, api.ChannelID, "unable to send msg")
		CheckTooLongMsgErr(err, api, v.Header, ActSend, "")

	} else {
		err = api.Discorder.EditMessage(api.ChannelID, v.MessageID, v.Content)
		ChannelCheckWarn(err, api.ChannelID, "unable to edit msg")
		CheckTooLongMsgErr(err, api, v.Header, ActEdit, v.MessageID)
	}
}
