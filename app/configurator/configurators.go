package configurator

import "github.com/darklab8/fl-darkbot/app/settings/types"

type Players struct {
	Systems ConfiguratorSystem
	Regions ConfiguratorRegion
	Enemies ConfiguratorPlayerEnemy
	Friends ConfiguratorPlayerFriend
	Events  ConfiguratorPlayerEvent
}
type Alerts struct {
	NeutralsGreaterThan    CfgAlertNeutralPlayersGreaterThan
	EnemiesGreaterThan     CfgAlertEnemyPlayersGreaterThan
	FriendsGreaterThan     CfgAlertFriendPlayersGreaterThan
	BaseHealthLowerThan    CfgAlertBaseHealthLowerThan
	BaseHealthIsDecreasing CfgAlertBaseHealthIsDecreasing
	BaseIsUnderAttack      CfgAlertBaseIsUnderAttack
	PingMessage            CfgAlertPingMessage
}

type ForumThread struct {
	Watch  ConfiguratorForumWatch
	Ignore ConfiguratorForumIgnore
}

type ForumSubforum struct {
	Watch  ConfiguratorSubForumWatch
	Ignore ConfiguratorSubForumIgnore
}

type Forum struct {
	Thread   ForumThread
	Subforum ForumSubforum
}
type Base struct {
	Tags    ConfiguratorBase
	OrderBy CfgBaseOrderingKey
}
type Configurators struct {
	Bases    Base
	Players  Players
	Alerts   Alerts
	Forum    Forum
	Configur *Configurator
	Channels ConfiguratorChannel
}

func NewConfigugurators(dbpath types.Dbpath) *Configurators {

	configur := NewConfigurator(dbpath)
	return NewConfiguratorsFromConfigur(configur)
}

func NewConfiguratorsFromConfigur(configur *Configurator) *Configurators {
	return &Configurators{
		Configur: configur,
		Channels: NewConfiguratorChannel(configur),
		Bases: Base{
			Tags:    NewConfiguratorBase(configur),
			OrderBy: NewCfgBaseOrderingKey(configur),
		},
		Players: Players{
			Systems: NewConfiguratorSystem(configur),
			Regions: NewConfiguratorRegion(configur),
			Enemies: NewConfiguratorPlayerEnemy(configur),
			Friends: NewConfiguratorPlayerFriend(configur),
			Events:  NewConfiguratorPlayerEvent(configur),
		},
		Alerts: Alerts{
			NeutralsGreaterThan:    NewCfgAlertNeutralPlayersGreaterThan(configur),
			EnemiesGreaterThan:     NewCfgAlertEnemyPlayersGreaterThan(configur),
			FriendsGreaterThan:     NewCfgAlertFriendPlayersGreaterThan(configur),
			BaseHealthLowerThan:    NewCfgAlertBaseHealthLowerThan(configur),
			BaseHealthIsDecreasing: NewCfgAlertBaseHealthIsDecreasing(configur),
			BaseIsUnderAttack:      NewCfgAlertBaseIsUnderAttack(configur),
			PingMessage:            NewCfgAlertPingMessage(configur),
		},
		Forum: Forum{
			Thread: ForumThread{
				Watch:  NewConfiguratorForumWatch(configur),
				Ignore: NewConfiguratorForumIgnore(configur),
			},
			Subforum: ForumSubforum{
				Watch:  NewConfiguratorSubForumWatch(configur),
				Ignore: NewConfiguratorSubForumIgnore(configur),
			},
		},
	}
}
