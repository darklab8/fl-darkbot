package viewer

import (
	"darkbot/discorder"
	"darkbot/utils"
	"fmt"
	"strings"
)

type ChannelView struct {
	discorder discorder.Discorder
	BaseView  BaseView
	channelID string
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
}

// Render new messages (ensure preserved Message ID)
func (v *ChannelView) Render() {
	v.BaseView.Content = fmt.Sprintf("%s\n```Content :)```", BaseViewHeader)
}

// Edit if message ID is present.
// Send if not present.
func (v ChannelView) Send() {
	if v.BaseView.MessageID == "" {
		v.discorder.SengMessage(v.channelID, v.BaseView.Content)
	} else {
		v.discorder.EditMessage(v.channelID, v.BaseView.MessageID, v.BaseView.Content)
	}
}
