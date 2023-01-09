package configurator

import "darkbot/configurator/models"

func (c Configurator) ActionBaseTagsAdd(channelID string, tags ...string) {

	c.db.FirstOrCreate(&models.Channel{ChannelID: channelID}, models.Channel{ChannelID: channelID})

	objs := []models.TagBase{}

	for _, tag := range tags {
		objs = append(objs, models.TagBase{
			TagTemplate: models.TagTemplate{
				Tag:       tag,
				FKChannel: models.FKChannel{ChannelID: channelID},
			},
		})
	}

	c.db.Create(objs)
}
