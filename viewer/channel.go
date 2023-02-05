package viewer

import (
	"darkbot/dtypes"
	"darkbot/utils/logger"
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

func NewChannelView(channelID string, dbpath dtypes.Dbpath) ChannelView {
	view := ChannelView{}
	view.API = apis.NewAPI(channelID, dbpath)
	view.BaseView = templ.NewTemplateBase(channelID, dbpath)
	view.PlayersView = templ.NewTemplatePlayers(channelID, dbpath)
	return view
}

// Query all Discord messages
// Try to grab already sent message by ID, if yes, assign to found objects with message ID.
func (v *ChannelView) Discover() {
	logger.Info("viewer.Init.channelID=", v.ChannelID)
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
		logger.Info("deleted message with id", msg.ID)
		deleteLimit--
	}
}

func (v *ChannelView) Delete() {
	for index, _ := range v.Msgs {
		v.Msgs[index] = nil
	}

	v.BaseView = nil
	v.PlayersView = nil
	v.API = nil
}
