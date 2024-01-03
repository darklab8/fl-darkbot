package configurator

import (
	"darkbot/app/configurator/models"
	"darkbot/app/settings/logus"
	"darkbot/app/settings/types"

	"github.com/darklab8/darklab_goutils/goutils/logus_core"
	"github.com/darklab8/darklab_goutils/goutils/utils"
)

func NewConfiguratorChannel(con *Configurator) ConfiguratorChannel {
	return ConfiguratorChannel{Configurator: con}
}

type ConfiguratorChannel struct {
	*Configurator
}

func (c ConfiguratorChannel) Add(channelID types.DiscordChannelID) error {

	result := c.db.Table("channels").Where("channel_id = ? AND deleted_at IS NOT NULL", channelID).Update("deleted_at", nil)
	if result.RowsAffected > 0 {
		return result.Error
	}

	if result.Error != nil {
		logus.Log.Info("channels.Add", logus_core.OptError(result.Error))
	}

	channel := models.Channel{ChannelID: channelID}
	result = c.db.Create(&channel)
	if result.Error != nil {
		logus.Log.Info("channels.Add.Error2=", logus_core.OptError(result.Error))
	}
	return result.Error
}

func (c ConfiguratorChannel) Remove(channelID types.DiscordChannelID) error {
	result := c.db.Where("channel_id = ?", channelID).Delete(&models.Channel{})
	if result.Error == nil && result.RowsAffected == 0 {
		return ErrorZeroAffectedRows{}
	}
	return result.Error
}

func (c ConfiguratorChannel) List() ([]types.DiscordChannelID, error) {
	objs := []models.Channel{}
	result := c.db.Find(&objs)

	return utils.CompL(objs,
		func(x models.Channel) types.DiscordChannelID { return x.ChannelID }), result.Error
}

func (c ConfiguratorChannel) IsEnabled(channelID types.DiscordChannelID) (bool, error) {
	objs := []models.Channel{}
	result := c.db.Where("channel_id = ?", channelID).Find(&objs)

	return result.RowsAffected != 0, result.Error
}
