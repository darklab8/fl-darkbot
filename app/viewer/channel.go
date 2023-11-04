package viewer

import (
	"darkbot/app/discorder"
	"darkbot/app/settings/logus"
	"darkbot/app/settings/types"
	"darkbot/app/viewer/apis"
	"darkbot/app/viewer/views"
	"darkbot/app/viewer/views/baseview"
	"darkbot/app/viewer/views/eventview"
	"darkbot/app/viewer/views/playerview"
	"strings"
	"time"
)

type ChannelView struct {
	Msgs      []discorder.DiscordMessage
	api       *apis.API
	ChannelID types.DiscordChannelID
	views     []views.View
}

// apis.NewAPI(view.ChannelID, dbpath)
func NewChannelView(api *apis.API, channelID types.DiscordChannelID) ChannelView {
	view := ChannelView{api: api}
	view.views = append(view.views, baseview.NewTemplateBase(api))
	view.views = append(view.views, playerview.NewTemplatePlayers(api))
	view.views = append(view.views, eventview.NewEventRenderer(api))
	view.ChannelID = channelID

	return view
}

// Query all Discord messages
// Try to grab already sent message by ID, if yes, assign to found objects with message ID.
func (v *ChannelView) Discover() error {
	logus.Info("viewer.Init.channelID=", logus.ChannelID(v.ChannelID))
	msgs, err := v.api.Discorder.GetLatestMessages(v.ChannelID)
	if logus.CheckWarn(err, "unable to grab latst msg", logus.ChannelID(v.ChannelID)) {
		return err
	}

	for _, msg := range msgs {
		for _, view := range v.views {
			view.DiscoverMessageID(msg.Content, msg.ID)
		}

	}

	v.Msgs = msgs

	return nil
}

// Render new messages (ensure preserved Message ID)
func (v *ChannelView) Render() {
	for _, view := range v.views {
		view.Render()
	}
}

// Edit if message ID is present.
// Send if not present.
func (v ChannelView) Send() {
	for _, view := range v.views {
		view.Send()
	}
}

func (v ChannelView) DeleteOld() {
	deleteLimit := 10
	for _, msg := range v.Msgs {

		matched_smth := false
		for _, view := range v.views {
			if view.MatchMessageID(msg.ID) {
				matched_smth = true
			}
		}
		if matched_smth {
			continue
		}

		// forbidding to delete messages that aren't having their own template renderer
		if strings.Contains(msg.Content, views.MsgViewHeader) {
			continue
		}

		if deleteLimit <= 0 {
			break
		}

		timeDiff := time.Now().Sub(msg.Timestamp)
		if timeDiff.Seconds() < 40 {
			continue
		}

		v.api.Discorder.DeleteMessage(v.ChannelID, msg.ID)
		logus.Info("deleted message with id", logus.MessageID(msg.ID))
		deleteLimit--
	}
}
