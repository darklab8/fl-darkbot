package configurator

import (
	"darkbot/configurator/models"
	"darkbot/utils"
	"darkbot/utils/logger"
)

type ConfiguratorChannel struct {
	Configurator
}

func (c ConfiguratorChannel) Add(channelID string) *ConfiguratorError {

	result := c.db.Table("channels").Where("channel_id = ? AND deleted_at IS NOT NULL", channelID).Update("deleted_at", nil)
	if result.RowsAffected > 0 {
		return (&ConfiguratorError{}).AppendSQLError(result)
	}

	if result.Error != nil {
		logger.Info("channels.Add.Error1=", result.Error.Error())
	}

	channel := models.Channel{ChannelID: channelID}
	result = c.db.Create(&channel)
	if result.Error != nil {
		logger.Info("channels.Add.Error2=", result.Error.Error())
	}
	return (&ConfiguratorError{}).AppendSQLError(result)
}

func (c ConfiguratorChannel) Remove(channelID string) *ConfiguratorError {
	result := c.db.Where("channel_id = ?", channelID).Delete(&models.Channel{})
	return (&ConfiguratorError{}).AppendSQLError(result)
}

func (c ConfiguratorChannel) List() ([]string, *ConfiguratorError) {
	objs := []models.Channel{}
	result := c.db.Find(&objs)

	return utils.CompL(objs,
		func(x models.Channel) string { return x.ChannelID }), (&ConfiguratorError{}).AppendSQLError(result)
}
