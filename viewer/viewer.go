package viewer

import (
	"darkbot/configurator"
	"darkbot/dtypes"
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
}

func NewViewer(dbpath dtypes.Dbpath) Viewer {
	return Viewer{
		dbpath:   dbpath,
		channels: configurator.ConfiguratorChannel{Configurator: configurator.NewConfigurator(dbpath)},
		delays: ViewerDelays{
			betweenChannels: 1,
			betweenLoops:    10,
		},
	}
}

func (v Viewer) Update() {
	logger.Info("Viewer.Update")

	// Query all channels
	channelIDs, _ := v.channels.List()
	logger.Info("Viewer.Update.channelIDs=", fmt.Sprintf("%v", channelIDs))

	// For each channel
	for _, channelID := range channelIDs {
		view := NewChannelView(channelID, v.dbpath)
		view.Discover()
		view.Render()
		view.Send()
		view.DeleteOld()
		time.Sleep(time.Duration(v.delays.betweenChannels) * time.Second)
	}
	time.Sleep(time.Duration(v.delays.betweenLoops) * time.Second)
}
