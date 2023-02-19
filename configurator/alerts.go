package configurator

import "darkbot/configurator/models"

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

type alertMethods[T any] interface {
	GetThreshold() int
	SetThreshold(channelID string, value int) T
}

type alertThreshold[T any] interface {
	models.AlertNeutralPlayersEqualOrGreater |
		models.AlertEnemiesEqualOrGreater |
		models.AlertFriendsEqualOrGreater |
		models.AlertBaseHealthLowerThan

	alertMethods[T]
}

type alertBool interface {
	models.AlertBaseIfHealthDecreasing |
		models.AlertBaseIfUnderAttack
}

type ConfiguratorAlertThreshold[T alertThreshold[T]] struct {
	Configurator
}

type ConfiguratorAlertBool[T alertBool] struct {
	Configurator
}

type CfgAlertNeutralPlayersGreaterThan = ConfiguratorAlertThreshold[models.AlertNeutralPlayersEqualOrGreater]
type CfgAlertEnemyPlayersGreaterThan = ConfiguratorAlertThreshold[models.AlertEnemiesEqualOrGreater]
type CfgAlertFriendPlayersGreaterThan = ConfiguratorAlertThreshold[models.AlertFriendsEqualOrGreater]
type CfgAlertBaseHealthLowerThan = ConfiguratorAlertThreshold[models.AlertBaseHealthLowerThan]
type CfgAlertBaseHealthIsDecreasing = ConfiguratorAlertBool[models.AlertBaseIfHealthDecreasing]
type CfgAlertBaseIsUnderAttack = ConfiguratorAlertBool[models.AlertBaseIfUnderAttack]

func (c ConfiguratorAlertThreshold[T]) Set(channelID string, value int) *ConfiguratorError {
	var obj T
	result := obj.SetThreshold(channelID, value)
	result2 := c.db.Create(&result)
	return (&ConfiguratorError{}).AppendSQLError(result2)
}

func (c ConfiguratorAlertThreshold[T]) Unset(channelID string) *ConfiguratorError {
	objs := []T{}
	result := c.db.Unscoped().Where("channel_id = ?", channelID).Find(&objs)
	result = c.db.Unscoped().Delete(&objs)
	return (&ConfiguratorError{}).AppendSQLError(result)
}

func (c ConfiguratorAlertBool[T]) Disable(channelID string) *ConfiguratorError {
	objs := []T{}
	result := c.db.Unscoped().Where("channel_id = ?", channelID).Find(&objs)
	result = c.db.Unscoped().Delete(&objs)
	return (&ConfiguratorError{}).AppendSQLError(result)
}

func (c ConfiguratorAlertBool[T]) Enable(channelID string) *ConfiguratorError {
	obj := T{
		AlertTemplate: models.AlertTemplate{ChannelID: channelID},
	}
	result := c.db.Create(&obj)
	return (&ConfiguratorError{}).AppendSQLError(result)
}

func (c ConfiguratorAlertBool[T]) Status(channelID string) (bool, *ConfiguratorError) {
	obj := T{}
	result := c.db.Where("channel_id = ?", channelID).First(&obj)

	return result.Error == nil, (&ConfiguratorError{}).AppendSQLError(result)
}

func (c ConfiguratorAlertThreshold[T]) Status(channelID string) (*int, *ConfiguratorError) {
	var obj T
	result := c.db.Where("channel_id = ?", channelID).First(&obj)
	if result.Error != nil {
		return nil, (&ConfiguratorError{}).AppendSQLError(result)
	}

	integer := obj.GetThreshold()
	return &integer, (&ConfiguratorError{}).AppendSQLError(result)
}
