package configurator

import (
	"darkbot/configurator/models"
	"darkbot/utils"
)

type ConfiguratorChannel struct {
	Configurator
}

func (c ConfiguratorChannel) Add(channelID string) error {

	result := c.db.Table("channels").Where("channel_id = ?", channelID).Update("deleted_at", nil)
	if result.RowsAffected > 0 {
		return result.Error
	}

	if result.Error != nil {
		utils.LogInfo("channels.Add.Error1=", result.Error.Error())
	}

	channel := models.Channel{ChannelID: channelID}
	result = c.db.FirstOrCreate(&channel)
	if result.Error != nil {
		utils.LogInfo("channels.Add.Error2=", result.Error.Error())
	}
	return result.Error
}

func (c ConfiguratorChannel) Remove(channelID string) error {
	return c.db.Where("channel_id = ?", channelID).Delete(&models.Channel{}).Error
}

func (c ConfiguratorChannel) List() ([]string, error) {
	var err error
	objs := []models.Channel{}
	err = c.db.Find(&objs).Error

	return utils.CompL(objs,
		func(x models.Channel) string { return x.ChannelID }), err
}
