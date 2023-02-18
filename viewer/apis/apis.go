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
	Systems   configurator.ConfiguratorSystem
	Regions   configurator.ConfiguratorRegion
	Enemies   configurator.ConfiguratorPlayerEnemy
	Friends   configurator.ConfiguratorPlayerFriend
}

func NewAPI(channelID string, dbpath dtypes.Dbpath) API {
	return API{
		Discorder: discorder.NewClient(),
		ChannelID: channelID,
		Scrappy:   scrappy.Storage,
		Bases:     configurator.ConfiguratorBase{Configurator: configurator.NewConfigurator(dbpath)},
		Systems:   configurator.ConfiguratorSystem{Configurator: configurator.NewConfigurator(dbpath)},
		Regions:   configurator.ConfiguratorRegion{Configurator: configurator.NewConfigurator(dbpath)},
		Enemies:   configurator.ConfiguratorPlayerEnemy{Configurator: configurator.NewConfigurator(dbpath)},
		Friends:   configurator.ConfiguratorPlayerFriend{Configurator: configurator.NewConfigurator(dbpath)},
	}
}
