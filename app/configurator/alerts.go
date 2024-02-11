package configurator

import (
	"github.com/darklab8/fl-darkbot/app/configurator/models"
	"github.com/darklab8/fl-darkbot/app/settings/logus"
	"github.com/darklab8/fl-darkbot/app/settings/types"
)

type AlertThresholdType interface {
	models.AlertNeutralPlayersEqualOrGreater |
		models.AlertEnemiesEqualOrGreater |
		models.AlertFriendsEqualOrGreater |
		models.AlertBaseHealthLowerThan

	GetThreshold() int
}

type AlertBoolType interface {
	models.AlertBaseIfHealthDecreasing |
		models.AlertBaseIfUnderAttack
}

type AlertStringType interface {
	models.AlertPingMessage |
		models.ConfigBaseOrderingKey
	GetValue() string
}

type IConfiguratorAlertThreshold[T AlertThresholdType] struct {
	*Configurator
}

type IConfiguratorAlertBool[T AlertBoolType] struct {
	*Configurator
}

type IConfiguratorAlertString[T AlertStringType] struct {
	*Configurator
}

func NewConfiguratorAlertThreshold[T AlertThresholdType](configurator *Configurator) IConfiguratorAlertThreshold[T] {
	t := IConfiguratorAlertThreshold[T]{Configurator: configurator}
	return t
}
func NewConfiguratorAlertBool[T AlertBoolType](configurator *Configurator) IConfiguratorAlertBool[T] {
	t := IConfiguratorAlertBool[T]{Configurator: configurator}
	return t
}
func NewConfiguratorAlertString[T AlertStringType](configurator *Configurator) IConfiguratorAlertString[T] {
	t := IConfiguratorAlertString[T]{Configurator: configurator}
	return t
}

type CfgAlertNeutralPlayersGreaterThan = IConfiguratorAlertThreshold[models.AlertNeutralPlayersEqualOrGreater]

var NewCfgAlertNeutralPlayersGreaterThan = NewConfiguratorAlertThreshold[models.AlertNeutralPlayersEqualOrGreater]

type CfgAlertEnemyPlayersGreaterThan = IConfiguratorAlertThreshold[models.AlertEnemiesEqualOrGreater]

var NewCfgAlertEnemyPlayersGreaterThan = NewConfiguratorAlertThreshold[models.AlertEnemiesEqualOrGreater]

type CfgAlertFriendPlayersGreaterThan = IConfiguratorAlertThreshold[models.AlertFriendsEqualOrGreater]

var NewCfgAlertFriendPlayersGreaterThan = NewConfiguratorAlertThreshold[models.AlertFriendsEqualOrGreater]

type CfgAlertBaseHealthLowerThan = IConfiguratorAlertThreshold[models.AlertBaseHealthLowerThan]

var NewCfgAlertBaseHealthLowerThan = NewConfiguratorAlertThreshold[models.AlertBaseHealthLowerThan]

type CfgAlertBaseHealthIsDecreasing = IConfiguratorAlertBool[models.AlertBaseIfHealthDecreasing]

var NewCfgAlertBaseHealthIsDecreasing = NewConfiguratorAlertBool[models.AlertBaseIfHealthDecreasing]

type CfgAlertBaseIsUnderAttack = IConfiguratorAlertBool[models.AlertBaseIfUnderAttack]

var NewCfgAlertBaseIsUnderAttack = NewConfiguratorAlertBool[models.AlertBaseIfUnderAttack]

type CfgAlertPingMessage = IConfiguratorAlertString[models.AlertPingMessage]

var NewCfgAlertPingMessage = NewConfiguratorAlertString[models.AlertPingMessage]

type CfgBaseOrderingKey = IConfiguratorAlertString[models.ConfigBaseOrderingKey]

var NewCfgBaseOrderingKey = NewConfiguratorAlertString[models.ConfigBaseOrderingKey]

func (c IConfiguratorAlertThreshold[T]) Set(channelID types.DiscordChannelID, value int) error {
	c.Unset(channelID)
	obj := T{
		OneValueTemplate:    models.OneValueTemplate{ChannelID: channelID},
		AlertTresholdShared: models.AlertTresholdShared{Threshold: value},
	}
	result2 := c.db.Create(&obj)

	return result2.Error
}

func (c IConfiguratorAlertThreshold[T]) Unset(channelID types.DiscordChannelID) error {
	objs := []T{}
	result := c.db.Unscoped().Where("channel_id = ?", channelID).Find(&objs)
	if result.RowsAffected == 0 {
		return ErrorZeroAffectedRows{}
	}
	result = c.db.Unscoped().Delete(&objs)

	return result.Error
}

func (c IConfiguratorAlertThreshold[T]) Status(channelID types.DiscordChannelID) (int, error) {
	var obj T
	result := c.db.Where("channel_id = ?", channelID).Limit(1).Find(&obj)

	if result.Error != nil {
		return 0, result.Error
	}

	if result.RowsAffected == 0 {
		return 0, ErrorZeroAffectedRows{}
	}

	integer := obj.GetThreshold()
	return integer, result.Error
}

///////////////////////////

func (c IConfiguratorAlertBool[T]) Enable(channelID types.DiscordChannelID) error {
	obj := T{OneValueTemplate: models.OneValueTemplate{ChannelID: channelID}}
	result := c.db.Create(&obj)
	return result.Error
}

func (c IConfiguratorAlertBool[T]) Disable(channelID types.DiscordChannelID) error {
	objs := []T{}
	c.db.Unscoped().Where("channel_id = ?", channelID).Find(&objs)
	if len(objs) == 0 {
		return ErrorZeroAffectedRows{}
	}
	result := c.db.Unscoped().Delete(&objs)
	return result.Error
}

func (c IConfiguratorAlertBool[T]) Status(channelID types.DiscordChannelID) (bool, error) {
	obj := T{}
	result := c.db.Where("channel_id = ?", channelID).Limit(1).Find(&obj)

	if result.Error != nil {
		return false, result.Error
	}

	if result.RowsAffected == 0 {
		return false, ErrorZeroAffectedRows{}
	}

	return true, nil
}

////////////////////////////

func (c IConfiguratorAlertString[T]) Set(channelID types.DiscordChannelID, value string) error {
	c.Unset(channelID)
	obj := T{
		OneValueTemplate: models.OneValueTemplate{ChannelID: channelID},
		Value:            value,
	}
	result2 := c.db.Create(&obj)
	return result2.Error
}

func (c IConfiguratorAlertString[T]) Unset(channelID types.DiscordChannelID) error {
	objs := []T{}
	result := c.db.Unscoped().Where("channel_id = ?", channelID).Find(&objs)
	logus.Log.CheckWarn(result.Error, "attempted to unset with errors", logus.GormResult(result))
	result = c.db.Unscoped().Delete(&objs)
	return result.Error
}

func (c IConfiguratorAlertString[T]) Status(channelID types.DiscordChannelID) (string, error) {
	var obj T
	result := c.db.Where("channel_id = ?", channelID).Limit(1).Find(&obj)
	if result.Error != nil {
		return "", result.Error
	}

	if result.RowsAffected == 0 {
		return "", ErrorZeroAffectedRows{}
	}

	str := obj.GetValue()
	return str, result.Error
}
