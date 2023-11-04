package apis

import (
	"darkbot/app/configurator"
	"darkbot/app/discorder"
	"darkbot/app/scrappy"
	"darkbot/app/settings/types"
)

type Players struct {
	Systems configurator.ConfiguratorSystem
	Regions configurator.ConfiguratorRegion
	Enemies configurator.ConfiguratorPlayerEnemy
	Friends configurator.ConfiguratorPlayerFriend
	Events  configurator.ConfiguratorPlayerEvent
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
	Configur  *configurator.Configurator
}

type apiParam func(api *API)

func WithStorage(storage *scrappy.ScrappyStorage) apiParam {
	return func(api *API) {
		api.Scrappy = storage
	}
}

// Reusing API is important. db connections are memory leaking or some other stuff
func (api *API) SetChannelID(ChannelID types.DiscordChannelID) *API {
	api.ChannelID = ChannelID
	return api
}

func NewAPI(ChannelID types.DiscordChannelID, dbpath types.Dbpath, opts ...apiParam) *API {
	configur := configurator.NewConfigurator(dbpath)
	api := &API{
		Configur:  configur,
		ChannelID: ChannelID,
		Discorder: discorder.NewClient(),
		Scrappy:   scrappy.Storage,
		Bases:     configurator.NewConfiguratorBase(configur),
		Players: Players{
			Systems: configurator.NewConfiguratorSystem(configur),
			Regions: configurator.NewConfiguratorRegion(configur),
			Enemies: configurator.NewConfiguratorPlayerEnemy(configur),
			Friends: configurator.NewConfiguratorPlayerFriend(configur),
			Events:  configurator.NewConfiguratorPlayerEvent(configur),
		},
		Alerts: Alerts{
			NeutralsGreaterThan:    configurator.NewCfgAlertNeutralPlayersGreaterThan(configur),
			EnemiesGreaterThan:     configurator.NewCfgAlertEnemyPlayersGreaterThan(configur),
			FriendsGreaterThan:     configurator.NewCfgAlertFriendPlayersGreaterThan(configur),
			BaseHealthLowerThan:    configurator.NewCfgAlertBaseHealthLowerThan(configur),
			BaseHealthIsDecreasing: configurator.NewCfgAlertBaseHealthIsDecreasing(configur),
			BaseIsUnderAttack:      configurator.NewCfgAlertBaseIsUnderAttack(configur),
			PingMessage:            configurator.NewCfgAlertPingMessage(configur),
		},
	}

	for _, opt := range opts {
		opt(api)
	}

	return api
}
