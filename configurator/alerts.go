package configurator

import (
	"darkbot/configurator/models"
	"darkbot/settings/types"
)

type iConfiguratorBoolAlert interface {
	Enable(channelID string) *ConfiguratorError
	Disable(channelID string) *ConfiguratorError
	Status(channelID string) *ConfiguratorError
}

type iConfiguratorThresholdAlert interface {
	Set(channelID string, value int) *ConfiguratorError
	Unset(channelID string) *ConfiguratorError
	Status(channelID string) (*int, *ConfiguratorError)
}

type iConfiguratorStringValue interface {
	Set(channelID string, value string) *ConfiguratorError
	Unset(channelID string) *ConfiguratorError
	Status(channelID string) (string, *ConfiguratorError)
}

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
	models.AlertPingMessage

	GetValue() string
}

type IConfiguratorAlertThreshold[T AlertThresholdType] struct {
	Configurator
}

type IConfiguratorAlertBool[T AlertBoolType] struct {
	Configurator
}

type IConfiguratorAlertString[T AlertStringType] struct {
	Configurator
}

type CfgAlertNeutralPlayersGreaterThan = IConfiguratorAlertThreshold[models.AlertNeutralPlayersEqualOrGreater]
type CfgAlertEnemyPlayersGreaterThan = IConfiguratorAlertThreshold[models.AlertEnemiesEqualOrGreater]
type CfgAlertFriendPlayersGreaterThan = IConfiguratorAlertThreshold[models.AlertFriendsEqualOrGreater]
type CfgAlertBaseHealthLowerThan = IConfiguratorAlertThreshold[models.AlertBaseHealthLowerThan]
type CfgAlertBaseHealthIsDecreasing = IConfiguratorAlertBool[models.AlertBaseIfHealthDecreasing]
type CfgAlertBaseIsUnderAttack = IConfiguratorAlertBool[models.AlertBaseIfUnderAttack]
type CfgAlertPingMessage = IConfiguratorAlertString[models.AlertPingMessage]

func (c IConfiguratorAlertThreshold[T]) Set(channelID types.DiscordChannelID, value int) *ConfiguratorError {
	c.Unset(channelID)
	obj := T{
		AlertTemplate:       models.AlertTemplate{ChannelID: channelID},
		AlertTresholdShared: models.AlertTresholdShared{Threshold: value},
	}
	result2 := c.db.Create(&obj)
	return (&ConfiguratorError{}).AppendSQLError(result2)
}

func (c IConfiguratorAlertThreshold[T]) Unset(channelID types.DiscordChannelID) *ConfiguratorError {
	objs := []T{}
	result := c.db.Unscoped().Where("channel_id = ?", channelID).Find(&objs)
	result = c.db.Unscoped().Delete(&objs)
	return (&ConfiguratorError{}).AppendSQLError(result)
}

func (c IConfiguratorAlertThreshold[T]) Status(channelID types.DiscordChannelID) (*int, *ConfiguratorError) {
	var obj T
	result := c.db.Where("channel_id = ?", channelID).First(&obj)
	if result.Error != nil {
		return nil, (&ConfiguratorError{}).AppendSQLError(result)
	}

	integer := obj.GetThreshold()
	return &integer, (&ConfiguratorError{}).AppendSQLError(result)
}

///////////////////////////

func (c IConfiguratorAlertBool[T]) Enable(channelID types.DiscordChannelID) *ConfiguratorError {
	obj := T{
		AlertTemplate: models.AlertTemplate{ChannelID: channelID},
	}
	result := c.db.Create(&obj)
	return (&ConfiguratorError{}).AppendSQLError(result)
}

func (c IConfiguratorAlertBool[T]) Disable(channelID types.DiscordChannelID) *ConfiguratorError {
	objs := []T{}
	result := c.db.Unscoped().Where("channel_id = ?", channelID).Find(&objs)
	result = c.db.Unscoped().Delete(&objs)
	return (&ConfiguratorError{}).AppendSQLError(result)
}

func (c IConfiguratorAlertBool[T]) Status(channelID types.DiscordChannelID) (bool, *ConfiguratorError) {
	obj := T{}
	result := c.db.Where("channel_id = ?", channelID).First(&obj)

	return result.Error == nil, (&ConfiguratorError{}).AppendSQLError(result)
}

////////////////////////////

func (c IConfiguratorAlertString[T]) Set(channelID types.DiscordChannelID, value string) *ConfiguratorError {
	c.Unset(channelID)
	obj := T{
		AlertTemplate: models.AlertTemplate{ChannelID: channelID},
		Value:         value,
	}
	result2 := c.db.Create(&obj)
	return (&ConfiguratorError{}).AppendSQLError(result2)
}

func (c IConfiguratorAlertString[T]) Unset(channelID types.DiscordChannelID) *ConfiguratorError {
	objs := []T{}
	result := c.db.Unscoped().Where("channel_id = ?", channelID).Find(&objs)
	result = c.db.Unscoped().Delete(&objs)
	return (&ConfiguratorError{}).AppendSQLError(result)
}

func (c IConfiguratorAlertString[T]) Status(channelID types.DiscordChannelID) (types.PingMessage, *ConfiguratorError) {
	var obj T
	result := c.db.Where("channel_id = ?", channelID).First(&obj)
	if result.Error != nil {
		return "", (&ConfiguratorError{}).AppendSQLError(result)
	}

	str := obj.GetValue()
	return types.PingMessage(str), (&ConfiguratorError{}).AppendSQLError(result)
}
