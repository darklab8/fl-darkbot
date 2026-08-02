package configurator

import (
	"github.com/darklab8/fl-darkbot/app/configurator/models"
	"github.com/darklab8/fl-darkbot/app/settings/logus"
	"github.com/darklab8/fl-darkbot/app/settings/types"
)

type IConfiguratorBasePriority[T BasePriorityType] struct {
	*Configurator
}

func NewConfiguratorBasePriority[T BasePriorityType](configurator *Configurator) IConfiguratorBasePriority[T] {
	t := IConfiguratorBasePriority[T]{Configurator: configurator}
	return t
}

type CfgConfigBasePriorty = IConfiguratorBasePriority[models.ConfigBasePriorty]

var NewCfgConfigBasePriorty = NewConfiguratorBasePriority[models.ConfigBasePriorty]

func (c IConfiguratorBasePriority[T]) Add(channelID types.DiscordChannelID, base_nickname string, value int) error {
	obj := T{
		MultiValueTemplate: models.MultiValueTemplate{ChannelID: channelID},
		PriorityValue:      models.PriorityValue{Priority: value},
		BaseNick:           models.BaseNick{BaseNickname: base_nickname},
	}
	result2 := c.db.Create(&obj)

	return result2.Error
}

func (c IConfiguratorBasePriority[T]) Remove(channelID types.DiscordChannelID, base_nickname string) error {
	objs := []T{}
	result := c.db.Unscoped().Where("channel_id = ?", channelID).Find(&objs)
	if result.RowsAffected == 0 {
		return ErrorZeroAffectedRows{}
	}
	result = c.db.Unscoped().Delete(&objs)

	return result.Error
}

func (c IConfiguratorBasePriority[T]) Get(channelID types.DiscordChannelID) (map[string]int, error) {
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
		output[obj.GetBaseNickname()] = obj.GetPriority()
	}

	return output, nil
}

func (c IConfiguratorBasePriority[T]) Clear(channelID types.DiscordChannelID) error {
	tags := []T{}
	result := c.db.Unscoped().Where("channel_id = ?", channelID).Find(&tags)
	if len(tags) == 0 {
		return ErrorZeroAffectedRows{ExtraMsg: "no pob priority configs found"}
	}
	logus.Log.Debug("Clear.Find", logus.GormResult(result))
	result = c.db.Unscoped().Delete(&tags)
	logus.Log.Debug("Clear.Detete", logus.GormResult(result))
	return result.Error
}
