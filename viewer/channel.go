package viewer

import (
	"darkbot/utils"
	"strings"
)

type ChannelView struct {
	ViewConfig
	BaseView BaseView
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
}

// Render new messages (ensure preserved Message ID)
func (v *ChannelView) Render() {
	v.BaseView.Render()
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