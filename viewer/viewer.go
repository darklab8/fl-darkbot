package viewer

import (
	"darkbot/configurator"
	"darkbot/settings"
	"darkbot/settings/logus"
	"darkbot/settings/types"
	"time"
)

type ViewerDelays struct {
	betweenChannels int
	betweenLoops    types.ScrappyLoopDelay
}

type Viewer struct {
	channels configurator.ConfiguratorChannel
	delays   ViewerDelays
	dbpath   types.Dbpath
	channel  ChannelView
}

func NewViewer(dbpath types.Dbpath) Viewer {
	return Viewer{
		dbpath:   dbpath,
		channels: configurator.NewConfiguratorChannel(configurator.NewConfigurator(dbpath)),
		delays: ViewerDelays{
			betweenChannels: 1,
			betweenLoops:    settings.LoopDelay,
		},
		channel: NewChannelView(dbpath),
	}
}

func (v Viewer) Update() {
	logus.Info("Viewer.Update")

	// Query all channels
	channelIDs, _ := v.channels.List()
	logus.Info("Viewer.Update.channelIDs=", logus.ChannelIDs(channelIDs))

	// For each channel
	for _, channelID := range channelIDs {
		v.channel.Setup(channelID)
		err := v.channel.Discover()
		if logus.CheckWarn(err, "unable to grab Discord msgs", logus.ChannelID(channelID)) {
			continue
		}
		v.channel.Render()
		v.channel.Send()
		v.channel.DeleteOld()
		time.Sleep(time.Duration(v.delays.betweenChannels) * time.Second)
	}
	time.Sleep(time.Duration(v.delays.betweenLoops) * time.Second)
}
