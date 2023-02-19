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

type alertThreshold interface {
	models.AlertNeutralPlayersEqualOrGreater |
		models.AlertEnemiesEqualOrGreater |
		models.AlertFriendsEqualOrGreater |
		models.AlertBaseHealthLowerThan

	GetThreshold() int
}

type alertBool interface {
	models.AlertBaseIfHealthDecreasing |
		models.AlertBaseIfUnderAttack
}

type IConfiguratorAlertThreshold[T alertThreshold] struct {
	Configurator
}

type IConfiguratorAlertBool[T alertBool] struct {
	Configurator
}

type CfgAlertNeutralPlayersGreaterThan = IConfiguratorAlertThreshold[models.AlertNeutralPlayersEqualOrGreater]
type CfgAlertEnemyPlayersGreaterThan = IConfiguratorAlertThreshold[models.AlertEnemiesEqualOrGreater]
type CfgAlertFriendPlayersGreaterThan = IConfiguratorAlertThreshold[models.AlertFriendsEqualOrGreater]
type CfgAlertBaseHealthLowerThan = IConfiguratorAlertThreshold[models.AlertBaseHealthLowerThan]
type CfgAlertBaseHealthIsDecreasing = IConfiguratorAlertBool[models.AlertBaseIfHealthDecreasing]
type CfgAlertBaseIsUnderAttack = IConfiguratorAlertBool[models.AlertBaseIfUnderAttack]

func (c IConfiguratorAlertThreshold[T]) Set(channelID string, value int) *ConfiguratorError {
	obj := T{
		AlertTemplate:       models.AlertTemplate{ChannelID: channelID},
		AlertTresholdShared: models.AlertTresholdShared{Threshold: value},
	}
	result2 := c.db.Create(&obj)
	return (&ConfiguratorError{}).AppendSQLError(result2)
}

func (c IConfiguratorAlertThreshold[T]) Unset(channelID string) *ConfiguratorError {
	objs := []T{}
	result := c.db.Unscoped().Where("channel_id = ?", channelID).Find(&objs)
	result = c.db.Unscoped().Delete(&objs)
	return (&ConfiguratorError{}).AppendSQLError(result)
}

func (c IConfiguratorAlertBool[T]) Disable(channelID string) *ConfiguratorError {
	objs := []T{}
	result := c.db.Unscoped().Where("channel_id = ?", channelID).Find(&objs)
	result = c.db.Unscoped().Delete(&objs)
	return (&ConfiguratorError{}).AppendSQLError(result)
}

func (c IConfiguratorAlertBool[T]) Enable(channelID string) *ConfiguratorError {
	obj := T{
		AlertTemplate: models.AlertTemplate{ChannelID: channelID},
	}
	result := c.db.Create(&obj)
	return (&ConfiguratorError{}).AppendSQLError(result)
}

func (c IConfiguratorAlertBool[T]) Status(channelID string) (bool, *ConfiguratorError) {
	obj := T{}
	result := c.db.Where("channel_id = ?", channelID).First(&obj)

	return result.Error == nil, (&ConfiguratorError{}).AppendSQLError(result)
}

func (c IConfiguratorAlertThreshold[T]) Status(channelID string) (*int, *ConfiguratorError) {
	var obj T
	result := c.db.Where("channel_id = ?", channelID).First(&obj)
	if result.Error != nil {
		return nil, (&ConfiguratorError{}).AppendSQLError(result)
	}

	integer := obj.GetThreshold()
	return &integer, (&ConfiguratorError{}).AppendSQLError(result)
}
