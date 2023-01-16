package configurator

import (
	"darkbot/configurator/models"
	"darkbot/utils"
)

type ConfiguratorChannel struct {
	Configurator
}

func (c ConfiguratorChannel) Add(names ...string) {
	objs := utils.CompL(names,
		func(channelID string) models.Channel { return models.Channel{ChannelID: channelID} })

	c.db.Create(objs)
}

func (c ConfiguratorChannel) List() []string {
	objs := []models.Channel{}
	c.db.Find(&objs)

	return utils.CompL(objs,
		func(x models.Channel) string { return x.ChannelID })
}
