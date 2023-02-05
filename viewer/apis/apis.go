package apis

import (
	"darkbot/configurator"
	"darkbot/discorder"
	"darkbot/dtypes"
	"darkbot/scrappy"
)

type API struct {
	Discorder discorder.Discorder
	ChannelID string
	Scrappy   *scrappy.ScrappyStorage
	Bases     configurator.ConfiguratorBase
}

func NewAPI(channelID string, dbpath dtypes.Dbpath) API {
	return API{
		Discorder: discorder.NewClient(),
		ChannelID: channelID,
		Scrappy:   scrappy.Storage,
		Bases:     configurator.ConfiguratorBase{Configurator: configurator.NewConfigurator(dbpath)},
	}
}
