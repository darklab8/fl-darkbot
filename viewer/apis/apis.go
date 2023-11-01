package apis

import (
	"darkbot/configurator"
	"darkbot/discorder"
	"darkbot/scrappy"
	"darkbot/settings/types"
)

type Players struct {
	Systems configurator.ConfiguratorSystem
	Regions configurator.ConfiguratorRegion
	Enemies configurator.ConfiguratorPlayerEnemy
	Friends configurator.ConfiguratorPlayerFriend
}
type Alerts struct {
	NeutralsGreaterThan    configurator.CfgAlertNeutralPlayersGreaterThan
	EnemiesGreaterThan     configurator.CfgAlertEnemyPlayersGreaterThan
	FriendsGreaterThan     configurator.CfgAlertFriendPlayersGreaterThan
	BaseHealthLowerThan    configurator.CfgAlertBaseHealthLowerThan
	BaseHealthIsDecreasing configurator.CfgAlertBaseHealthIsDecreasing
	BaseIsUnderAttack      configurator.CfgAlertBaseIsUnderAttack
	PingMessage            configurator.CfgAlertPingMessage
}
type API struct {
	Discorder discorder.Discorder
	ChannelID types.DiscordChannelID
	Scrappy   *scrappy.ScrappyStorage
	Bases     configurator.ConfiguratorBase
	Players   Players
	Alerts    Alerts
}

func NewAPI(channelID types.DiscordChannelID, dbpath types.Dbpath) API {
	dbconnection := configurator.NewConfigurator(dbpath)
	return API{
		Discorder: discorder.NewClient(),
		ChannelID: channelID,
		Scrappy:   scrappy.Storage,
		Bases:     configurator.ConfiguratorBase{Configurator: dbconnection},
		Players: Players{
			Systems: configurator.ConfiguratorSystem{Configurator: dbconnection},
			Regions: configurator.ConfiguratorRegion{Configurator: dbconnection},
			Enemies: configurator.ConfiguratorPlayerEnemy{Configurator: dbconnection},
			Friends: configurator.ConfiguratorPlayerFriend{Configurator: dbconnection},
		},
		Alerts: Alerts{
			NeutralsGreaterThan:    configurator.CfgAlertNeutralPlayersGreaterThan{Configurator: dbconnection},
			EnemiesGreaterThan:     configurator.CfgAlertEnemyPlayersGreaterThan{Configurator: dbconnection},
			FriendsGreaterThan:     configurator.CfgAlertFriendPlayersGreaterThan{Configurator: dbconnection},
			BaseHealthLowerThan:    configurator.CfgAlertBaseHealthLowerThan{Configurator: dbconnection},
			BaseHealthIsDecreasing: configurator.CfgAlertBaseHealthIsDecreasing{Configurator: dbconnection},
			BaseIsUnderAttack:      configurator.CfgAlertBaseIsUnderAttack{Configurator: dbconnection},
			PingMessage:            configurator.CfgAlertPingMessage{Configurator: dbconnection},
		},
	}
}
