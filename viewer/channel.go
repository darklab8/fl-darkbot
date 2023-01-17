package viewer

import (
	"darkbot/utils"
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

type ChannelView struct {
	ViewConfig
	BaseView BaseView
	Msgs     []*discordgo.Message
}

func NewChannelView(channelID string) ChannelView {
	view := ChannelView{}
	view.ViewConfig = NewViewerConfig(channelID)
	view.BaseView.ViewConfig = view.ViewConfig
	return view
}

// Query all Discord messages
// Try to grab already sent message by ID, if yes, assign to found objects with message ID.
func (v *ChannelView) Discover() {
	utils.LogInfo("viewer.Init.channelID=", v.channelID)
	msgs := v.discorder.GetLatestMessages(v.channelID)
	for _, msg := range msgs {
		if strings.Contains(msg.Content, BaseViewHeader) {
			v.BaseView.MessageID = msg.ID
		}
	}

	v.Msgs = msgs
}

// Render new messages (ensure preserved Message ID)
func (v *ChannelView) Render() {
	v.BaseView.Render()
}

type MsgAction int

const (
	ActSend MsgAction = iota
	ActEdit
)

func CheckTooLongMsgErr(err error, api ViewConfig, header string, action MsgAction, MessageID string) {
	if err == nil {
		return
	}

	if !strings.Contains(err.Error(), "BASE_TYPE_MAX_LENGTH") &&
		!strings.Contains(err.Error(), "or fewer in length") {
		return
	}

	msg := fmt.Sprintf("%s, %s, %s", BaseViewHeader, time.Now(), err)

	switch action {
	case ActSend:
		api.discorder.SengMessage(api.channelID, msg)
	case ActEdit:
		api.discorder.EditMessage(api.channelID, MessageID, msg)
	}

}

func ChannelCheckWarn(err error, channelID string, msg string) {
	utils.CheckWarn(err, "channelID=", channelID, msg)
}

// Edit if message ID is present.
// Send if not present.
func (v ChannelView) Send() {
	if v.BaseView.Content == "" {
		return
	}

	var err error
	if v.BaseView.MessageID == "" {
		err = v.discorder.SengMessage(v.channelID, v.BaseView.Content)
		ChannelCheckWarn(err, v.channelID, "unable to send msg")
		CheckTooLongMsgErr(err, v.ViewConfig, BaseViewHeader, ActSend, "")

	} else {
		err = v.discorder.EditMessage(v.channelID, v.BaseView.MessageID, v.BaseView.Content)
		ChannelCheckWarn(err, v.channelID, "unable to edit msg")
		CheckTooLongMsgErr(err, v.ViewConfig, BaseViewHeader, ActEdit, v.BaseView.MessageID)
	}
}

func (v ChannelView) DeleteOld() {
	deleteLimit := 10
	for _, msg := range v.Msgs {
		if msg.ID == v.BaseView.MessageID {
			continue
		}

		if deleteLimit <= 0 {
			break
		}

		timeDiff := time.Now().Sub(msg.Timestamp)
		if timeDiff.Seconds() < 40 {
			continue
		}

		v.discorder.DeleteMessage(v.channelID, msg.ID)
		utils.LogInfo("deleted message with id", msg.ID)
		deleteLimit--
	}
}
