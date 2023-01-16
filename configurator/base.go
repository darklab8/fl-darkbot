package configurator

import (
	"darkbot/configurator/models"
	"darkbot/utils"
	"fmt"
)

type IConfiguratorTags interface {
	TagsAdd(channelID string, tags ...string)
	TagsRemove(channelID string, tags ...string)
	TagsList(channelID string) []string
	TagsClear(channelID string)
}

type ConfiguratorBase struct {
	Configurator
}

func (c ConfiguratorBase) TagsAdd(channelID string, tags ...string) {
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

func (c ConfiguratorBase) TagsRemove(channelID string, tags ...string) {
	for _, tag := range tags {
		c.db.Where("channel_id = ? AND tag = ?", channelID, tag).Delete(&models.TagBase{})
	}
}

func (c ConfiguratorBase) TagsList(channelID string) []string {
	objs := []models.TagBase{}
	c.db.Where("channel_id = ?", channelID).Find(&objs)

	return utils.Comprehension(objs, func(x models.TagBase) string { return x.Tag })
}

func (c ConfiguratorBase) TagsClear(channelID string) {
	var tags []models.TagBase
	result := c.db.Unscoped().Where("channel_id = ?", channelID).Find(&tags)
	fmt.Println("Clear.Find.rowsAffected=", result.RowsAffected)
	fmt.Println("Clear.Find.err=", result.Error)
	result = c.db.Unscoped().Delete(&tags)
	fmt.Println("Clear.Delete.rowsAffected=", result.RowsAffected)
	fmt.Println("Clear.Delete.err=", result.Error)
}
