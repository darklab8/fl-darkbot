package viewer

import (
	"darkbot/utils"
	"darkbot/viewer/apis"
	"darkbot/viewer/templ"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

type ChannelView struct {
	*apis.API
	BaseView *templ.TemplateBase
	Msgs     []*discordgo.Message
}

func NewChannelView(channelID string) ChannelView {
	view := ChannelView{}
	view.API = apis.NewAPI(channelID)
	view.BaseView = templ.NewTemplateBase(channelID)
	return view
}

// Query all Discord messages
// Try to grab already sent message by ID, if yes, assign to found objects with message ID.
func (v *ChannelView) Discover() {
	utils.LogInfo("viewer.Init.channelID=", v.ChannelID)
	msgs := v.Discorder.GetLatestMessages(v.ChannelID)
	for _, msg := range msgs {
		if strings.Contains(msg.Content, v.BaseView.Header) {
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
	v.BaseView.Send()
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

		v.Discorder.DeleteMessage(v.ChannelID, msg.ID)
		utils.LogInfo("deleted message with id", msg.ID)
		deleteLimit--
	}
}
