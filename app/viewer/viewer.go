package viewer

import (
	"darkbot/app/configurator"
	"darkbot/app/settings"
	"darkbot/app/settings/logus"
	"darkbot/app/settings/types"
	"darkbot/app/viewer/apis"
	"time"
)

type ViewerDelays struct {
	betweenChannels int
	betweenLoops    types.ScrappyLoopDelay
}

type Viewer struct {
	channels configurator.ConfiguratorChannel
	delays   ViewerDelays
	api      *apis.API
}

func NewViewer(dbpath types.Dbpath) Viewer {
	api := apis.NewAPI("", dbpath)
	return Viewer{
		api:      api,
		channels: configurator.NewConfiguratorChannel(api.Configur),
		delays: ViewerDelays{
			betweenChannels: 1,
			betweenLoops:    settings.LoopDelay,
		},
	}
}

func (v Viewer) Update() {
	logus.Info("Viewer.Update")

	// Query all channels
	channelIDs, _ := v.channels.List()
	logus.Info("Viewer.Update.channelIDs=", logus.ChannelIDs(channelIDs))

	// For each channel
	for _, channelID := range channelIDs {
		v.api.SetChannelID(channelID)
		channel := NewChannelView(v.api, channelID)
		channel.Render()
		err := channel.Discover()
		if logus.CheckWarn(err, "unable to grab Discord msgs", logus.ChannelID(channelID)) {
			continue
		}
		channel.Send()
		channel.DeleteOld()
		time.Sleep(time.Duration(v.delays.betweenChannels) * time.Second)
	}
	time.Sleep(time.Duration(v.delays.betweenLoops) * time.Second)
}
