/*
Loop configurator settings, and send view to coresponding channels
*/

package viewer

import (
	"darkbot/configurator"
	"darkbot/discorder"
	"darkbot/utils"
	"fmt"
	"strings"
	"time"
)

const (
	BaseViewHeader = "#darkbot-base-view"
)

type BaseView struct {
	MessageID string
	Content   string
}

type ChannelView struct {
	discorder discorder.Discorder
	BaseView  BaseView
	channelID string
}

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

type ViewerDelays struct {
	betweenChannels int
	betweenLoops    int
}

type Viewer struct {
	channels configurator.ConfiguratorChannel
	delays   ViewerDelays
}

func NewViewer() Viewer {
	return Viewer{
		channels: configurator.ConfiguratorChannel{Configurator: configurator.NewConfigurator()},
		delays: ViewerDelays{
			betweenChannels: 1,
			betweenLoops:    10,
		},
	}
}

func (v Viewer) Update() {
	utils.LogInfo("Viewer.Update")

	// Query all channels
	channelIDs := v.channels.List()

	// For each channel
	// Query all Discord messages
	for _, channelID := range channelIDs {
		view := ChannelView{discorder: discorder.NewClient(), channelID: channelID}
		view.Discover()
		view.Render()
		view.Send()
		time.Sleep(time.Duration(v.delays.betweenChannels) * time.Second)
	}
	time.Sleep(time.Duration(v.delays.betweenLoops) * time.Second)
}

func Run() {
	utils.LogInfo("Viewer is now running.")

	for {
		NewViewer().Update()
	}
}
