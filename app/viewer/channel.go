package viewer

import (
	"fmt"
	"strings"
	"time"

	"github.com/darklab8/fl-darkbot/app/discorder"
	"github.com/darklab8/fl-darkbot/app/settings/logus"
	"github.com/darklab8/fl-darkbot/app/settings/types"
	"github.com/darklab8/fl-darkbot/app/viewer/apis"
	"github.com/darklab8/fl-darkbot/app/viewer/views"
	"github.com/darklab8/fl-darkbot/app/viewer/views/baseview"
	"github.com/darklab8/fl-darkbot/app/viewer/views/eventview"
	"github.com/darklab8/fl-darkbot/app/viewer/views/playerview"
	"github.com/darklab8/go-utils/utils/timeit"
)

type ChannelView struct {
	Msgs      []*discorder.DiscordMessage
	api       *apis.API
	ChannelID types.DiscordChannelID
	views     []views.View
}

func NewChannelView(api *apis.API, channelID types.DiscordChannelID) ChannelView {
	view := ChannelView{api: api}
	view.views = append(view.views, baseview.NewTemplateBase(api, channelID))
	view.views = append(view.views, playerview.NewTemplatePlayers(api, channelID))
	view.views = append(view.views, eventview.NewEventRenderer(api, channelID))
	view.ChannelID = channelID

	return view
}

// Query all Discord messages
// Try to grab already sent message by ID, if yes, assign to found objects with message ID.
func (v *ChannelView) Discover() error {
	logus.Log.Debug("viewer.Init.channelID=", logus.ChannelID(v.ChannelID))
	msgs, err := v.api.Discorder.GetLatestMessages(v.ChannelID)
	if logus.Log.CheckWarn(err, "unable to grab latst msg", logus.ChannelID(v.ChannelID)) {
		TryChannelAutoRemoval(v.api, err, v.ChannelID)
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

// RenderViews new messages (ensure preserved Message ID)
func (v *ChannelView) RenderViews() {
	for _, view := range v.views {
		view.RenderView()
	}
}

// Edit if message ID is present.
// Send if not present.
func (v ChannelView) Send() {
	for view_num, view := range v.views {
		timeit.NewTimerMFL(fmt.Sprintf("view.Send view_num=%d, view=%v", view_num, view), func() {
			view.Send()
		}, logus.ChannelID(v.ChannelID))
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

		// embeding msgs contain forumer emails
		if len(msg.Embeds) > 0 {
			continue
		}

		if deleteLimit <= 0 {
			break
		}

		timeDiff := time.Now().Sub(msg.Timestamp)
		if timeDiff.Seconds() < 40 {
			continue
		}

		err := v.api.Discorder.DeleteMessage(v.ChannelID, msg.ID)
		if err != nil {
			// No point to continue deleting in this channel if first one failed
			// it will make more optimized amount of network requests
			return
		}
		logus.Log.Info("deleted message with id", logus.MessageID(msg.ID))
		deleteLimit--
	}
}
