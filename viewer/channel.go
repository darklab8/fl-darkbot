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
	BaseView    *templ.TemplateBase
	Msgs        []*discordgo.Message
	PlayersView *templ.PlayersTemplates
}

func NewChannelView(channelID string) ChannelView {
	view := ChannelView{}
	view.API = apis.NewAPI(channelID)
	view.BaseView = templ.NewTemplateBase(channelID)
	view.PlayersView = templ.NewTemplatePlayers(channelID)
	return view
}

// Query all Discord messages
// Try to grab already sent message by ID, if yes, assign to found objects with message ID.
func (v *ChannelView) Discover() {
	utils.LogInfo("viewer.Init.channelID=", v.ChannelID)
	msgs := v.Discorder.GetLatestMessages(v.ChannelID)
	for _, msg := range msgs {
		v.BaseView.DiscoverMessageID(msg.Content, msg.ID)
		v.PlayersView.DiscoverMessageID(msg.Content, msg.ID)
	}

	v.Msgs = msgs
}

// Render new messages (ensure preserved Message ID)
func (v *ChannelView) Render() {
	v.BaseView.Render()
	v.PlayersView.Render()
}

// Edit if message ID is present.
// Send if not present.
func (v ChannelView) Send() {
	v.BaseView.Send()
	v.PlayersView.Send()
}

func (v ChannelView) DeleteOld() {
	deleteLimit := 10
	for _, msg := range v.Msgs {

		if v.BaseView.MatchMessageID(msg.ID) {
			continue
		}

		if v.PlayersView.MatchMessageID(msg.ID) {
			continue
		}

		// forbidding to delete messages that aren't having their own template renderer
		if strings.Contains(msg.Content, templ.MsgViewHeader) {
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
