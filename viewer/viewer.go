package viewer

import (
	"darkbot/configurator"
	"darkbot/utils"
	"fmt"
	"time"
)

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
	channelIDs, _ := v.channels.List()
	utils.LogInfo("Viewer.Update.channelIDs=", fmt.Sprintf("%v", channelIDs))

	// For each channel
	for _, channelID := range channelIDs {
		view := NewChannelView(channelID)
		view.Discover()
		view.Render()
		view.Send()
		view.DeleteOld()
		time.Sleep(time.Duration(v.delays.betweenChannels) * time.Second)
	}
	time.Sleep(time.Duration(v.delays.betweenLoops) * time.Second)
}
