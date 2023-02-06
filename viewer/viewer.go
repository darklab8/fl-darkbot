package viewer

import (
	"darkbot/configurator"
	"darkbot/dtypes"
	"darkbot/settings"
	"darkbot/utils/logger"
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
	dbpath   dtypes.Dbpath
	channel  ChannelView
}

func NewViewer(dbpath dtypes.Dbpath) Viewer {
	return Viewer{
		dbpath:   dbpath,
		channels: configurator.ConfiguratorChannel{Configurator: configurator.NewConfigurator(dbpath)},
		delays: ViewerDelays{
			betweenChannels: 1,
			betweenLoops:    settings.LoopDelay,
		},
		channel: NewChannelView(dbpath),
	}
}

func (v Viewer) Update() {
	logger.Info("Viewer.Update")

	// Query all channels
	channelIDs, _ := v.channels.List()
	logger.Info("Viewer.Update.channelIDs=", fmt.Sprintf("%v", channelIDs))

	// For each channel
	for _, channelID := range channelIDs {
		v.channel.Setup(channelID)
		v.channel.Discover()
		v.channel.Render()
		v.channel.Send()
		v.channel.DeleteOld()
		time.Sleep(time.Duration(v.delays.betweenChannels) * time.Second)
	}
	time.Sleep(time.Duration(v.delays.betweenLoops) * time.Second)
}
