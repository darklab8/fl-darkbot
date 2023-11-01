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
		Bases:     configurator.NewConfiguratorBase(dbconnection),
		Players: Players{
			Systems: configurator.NewConfiguratorSystem(dbconnection),
			Regions: configurator.NewConfiguratorRegion(dbconnection),
			Enemies: configurator.NewConfiguratorPlayerEnemy(dbconnection),
			Friends: configurator.NewConfiguratorPlayerFriend(dbconnection),
		},
		Alerts: Alerts{
			NeutralsGreaterThan:    configurator.NewCfgAlertNeutralPlayersGreaterThan(dbconnection),
			EnemiesGreaterThan:     configurator.NewCfgAlertEnemyPlayersGreaterThan(dbconnection),
			FriendsGreaterThan:     configurator.NewCfgAlertFriendPlayersGreaterThan(dbconnection),
			BaseHealthLowerThan:    configurator.NewCfgAlertBaseHealthLowerThan(dbconnection),
			BaseHealthIsDecreasing: configurator.NewCfgAlertBaseHealthIsDecreasing(dbconnection),
			BaseIsUnderAttack:      configurator.NewCfgAlertBaseIsUnderAttack(dbconnection),
			PingMessage:            configurator.NewCfgAlertPingMessage(dbconnection),
		},
	}
}
