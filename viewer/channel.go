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

// Edit if message ID is present.
// Send if not present.
func (v ChannelView) Send() {
	if v.BaseView.Content == "" {
		return
	}

	var err error
	if v.BaseView.MessageID == "" {
		err = v.discorder.SengMessage(v.channelID, v.BaseView.Content)

		if err != nil {
			err = v.discorder.SengMessage(v.channelID, fmt.Sprintf("%s, %s, %s", BaseViewHeader, time.Now(), err))
		}
	} else {
		err = v.discorder.EditMessage(v.channelID, v.BaseView.MessageID, v.BaseView.Content)

		if err != nil {
			err = v.discorder.EditMessage(v.channelID, v.BaseView.MessageID, fmt.Sprintf("%s, %s, %s", BaseViewHeader, time.Now(), err))
		}
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
