package configurator

import "github.com/darklab8/fl-darkbot/app/settings/types"

type Alerts struct {
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
		Alerts: Alerts{
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
