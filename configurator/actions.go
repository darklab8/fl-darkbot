package configurator

import (
	"darkbot/configurator/models"
	"fmt"
)

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

func (c Configurator) ActionBaseTagsList(channelID string) []models.TagBase {

	objs := []models.TagBase{}
	c.db.Where("channel_id = ?", channelID).Find(&objs)
	return objs
}

func (c Configurator) ActionBaseTagsClear(channelID string) {
	var tags []models.TagBase
	result := c.db.Unscoped().Where("channel_id = ?", channelID).Find(&tags)
	fmt.Println("Clear.Find.rowsAffected=", result.RowsAffected)
	fmt.Println("Clear.Find.err=", result.Error)
	result = c.db.Unscoped().Delete(&tags)
	fmt.Println("Clear.Delete.rowsAffected=", result.RowsAffected)
	fmt.Println("Clear.Delete.err=", result.Error)
}
