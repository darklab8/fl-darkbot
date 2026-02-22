package configurator

import (
	"errors"

	"github.com/darklab8/fl-darkbot/app/configurator/models"
	"github.com/darklab8/fl-darkbot/app/settings/logus"
	"github.com/darklab8/fl-darkbot/app/settings/types"
)

type AlertThresholdType interface {
	// models.SomeAlert |
	models.AlertBaseHealthLowerThan |
		models.AlertBaseMoneyBelow |
		models.AlertBaseCargoBelow

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

type AlertPoBGoodType interface {
	models.AlertPobGoodBelowThan | models.AlertPobGoodAboveThan
	GetGoodNickname() string
	GetThreshold() int
}

type IConfiguratorAlertThreshold[T AlertThresholdType] struct {
	*Configurator
}
type IConfiguratorAlertPoBGood[T AlertPoBGoodType] struct {
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
func NewConfiguratorAlertPoBGood[T AlertPoBGoodType](configurator *Configurator) IConfiguratorAlertPoBGood[T] {
	t := IConfiguratorAlertPoBGood[T]{Configurator: configurator}
	return t
}

type CfgAlertPoBGoodBelowThan = IConfiguratorAlertPoBGood[models.AlertPobGoodBelowThan]

var NewCfgAlertPoBGoodBelowThan = NewConfiguratorAlertPoBGood[models.AlertPobGoodBelowThan]

type CfgAlertPoBGoodAboveThan = IConfiguratorAlertPoBGood[models.AlertPobGoodAboveThan]

var NewCfgAlertPoBGoodAboveThan = NewConfiguratorAlertPoBGood[models.AlertPobGoodAboveThan]

type CfgAlertBaseHealthLowerThan = IConfiguratorAlertThreshold[models.AlertBaseHealthLowerThan]

var NewCfgAlertBaseHealthLowerThan = NewConfiguratorAlertThreshold[models.AlertBaseHealthLowerThan]

type CfgAlertBaseHealthIsDecreasing = IConfiguratorAlertBool[models.AlertBaseIfHealthDecreasing]

var NewCfgAlertBaseHealthIsDecreasing = NewConfiguratorAlertBool[models.AlertBaseIfHealthDecreasing]

type CfgAlertBaseIsUnderAttack = IConfiguratorAlertBool[models.AlertBaseIfUnderAttack]

var NewCfgAlertBaseIsUnderAttack = NewConfiguratorAlertBool[models.AlertBaseIfUnderAttack]

type CfgAlertBaseMoneyBelowThan = IConfiguratorAlertThreshold[models.AlertBaseMoneyBelow]

var NewCfgAlertBaseMoneyBelowThan = NewConfiguratorAlertThreshold[models.AlertBaseMoneyBelow]

type CfgAlertBaseCargoBelowThan = IConfiguratorAlertThreshold[models.AlertBaseCargoBelow]

var NewCfgAlertBaseCargoBelowThan = NewConfiguratorAlertThreshold[models.AlertBaseCargoBelow]

type CfgAlertPingMessage = IConfiguratorAlertString[models.AlertPingMessage]

var NewCfgAlertPingMessage = NewConfiguratorAlertString[models.AlertPingMessage]

type CfgBaseOrderingKey = IConfiguratorAlertString[models.ConfigBaseOrderingKey]

var NewCfgBaseOrderingKey = NewConfiguratorAlertString[models.ConfigBaseOrderingKey]

func (c IConfiguratorAlertThreshold[T]) Set(channelID types.DiscordChannelID, kind models.ThresholdIntegerKind, value int) error {
	c.Unset(channelID)

	if kind == models.ThresholdIntegerPercentage {
		if value < 0 || value > 100 {
			return errors.New("value should be within 0 to 100 range")
		}
	}

	obj := T{
		OneValueTemplate:     models.OneValueTemplate{ChannelID: channelID},
		AlertTresholdInteger: models.AlertTresholdInteger{Threshold: value},
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

////////////////////////////////

func (c IConfiguratorAlertPoBGood[T]) Add(channelID types.DiscordChannelID, good_nickname string, value int) error {
	obj := T{
		MultiValueTemplate:   models.MultiValueTemplate{ChannelID: channelID},
		AlertTresholdInteger: models.AlertTresholdInteger{Threshold: value},
		AlertPoBGood:         models.AlertPoBGood{GoodNickname: good_nickname},
	}
	result2 := c.db.Create(&obj)

	return result2.Error
}

func (c IConfiguratorAlertPoBGood[T]) Remove(channelID types.DiscordChannelID, good_nickname string) error {
	objs := []T{}
	result := c.db.Unscoped().Where("channel_id = ?", channelID).Find(&objs)
	if result.RowsAffected == 0 {
		return ErrorZeroAffectedRows{}
	}
	result = c.db.Unscoped().Delete(&objs)

	return result.Error
}

type PoBGoodStatus struct {
	Threshold    int
	GoodNickname string
}

func (c IConfiguratorAlertPoBGood[T]) Get(channelID types.DiscordChannelID) (map[string]int, error) {
	var objs []T = []T{}

	result := c.db.Where("channel_id = ?", channelID).Find(&objs)
	var output map[string]int

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, ErrorZeroAffectedRows{}
	}

	output = make(map[string]int)

	for _, obj := range objs {
		output[obj.GetGoodNickname()] = obj.GetThreshold()
	}

	return output, nil
}

func (c IConfiguratorAlertPoBGood[T]) Clear(channelID types.DiscordChannelID) error {
	tags := []T{}
	result := c.db.Unscoped().Where("channel_id = ?", channelID).Find(&tags)
	if len(tags) == 0 {
		return ErrorZeroAffectedRows{ExtraMsg: "no pob good alert configs found"}
	}
	logus.Log.Debug("Clear.Find", logus.GormResult(result))
	result = c.db.Unscoped().Delete(&tags)
	logus.Log.Debug("Clear.Detete", logus.GormResult(result))
	return result.Error
}
