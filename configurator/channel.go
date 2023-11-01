package configurator

import (
	"darkbot/configurator/models"
	"darkbot/settings/logus"
	"darkbot/settings/types"
	"darkbot/settings/utils"
)

type ConfiguratorChannel struct {
	Configurator
}

func (c ConfiguratorChannel) Add(channelID types.DiscordChannelID) *ConfiguratorError {

	result := c.db.Table("channels").Where("channel_id = ? AND deleted_at IS NOT NULL", channelID).Update("deleted_at", nil)
	if result.RowsAffected > 0 {
		return (&ConfiguratorError{}).AppendSQLError(result)
	}

	if result.Error != nil {
		logus.Info("channels.Add", logus.OptError(result.Error))
	}

	channel := models.Channel{ChannelID: channelID}
	result = c.db.Create(&channel)
	if result.Error != nil {
		logus.Info("channels.Add.Error2=", logus.OptError(result.Error))
	}
	return (&ConfiguratorError{}).AppendSQLError(result)
}

func (c ConfiguratorChannel) Remove(channelID types.DiscordChannelID) *ConfiguratorError {
	result := c.db.Where("channel_id = ?", channelID).Delete(&models.Channel{})
	return (&ConfiguratorError{}).AppendSQLError(result)
}

func (c ConfiguratorChannel) List() ([]types.DiscordChannelID, *ConfiguratorError) {
	objs := []models.Channel{}
	result := c.db.Find(&objs)

	return utils.CompL(objs,
		func(x models.Channel) types.DiscordChannelID { return x.ChannelID }), (&ConfiguratorError{}).AppendSQLError(result)
}
