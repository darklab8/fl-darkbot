package viewer

import (
	"darkbot/app/scrappy"
	"darkbot/app/settings"
	"darkbot/app/settings/logus"
	"darkbot/app/settings/types"
	"darkbot/app/viewer/apis"
	"time"
)

type ViewerDelays struct {
	betweenChannels int
	betweenLoops    types.ViewerLoopDelay
}

type Viewer struct {
	delays ViewerDelays
	api    *apis.API
}

func NewViewer(dbpath types.Dbpath, scrappy_storage *scrappy.ScrappyStorage) *Viewer {
	api := apis.NewAPI("", dbpath, scrappy_storage)
	return &Viewer{
		api: api,
		delays: ViewerDelays{
			betweenChannels: 1,
			betweenLoops:    settings.ViewerLoopDelay,
		},
	}
}

func (v *Viewer) Run() {
	logus.Info("Viewer is now running.")

	for {
		v.Update()
	}
}

func (v Viewer) Update() {
	logus.Info("Viewer.Update")

	// Query all channels
	channelIDs, _ := v.api.Channels.List()
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
