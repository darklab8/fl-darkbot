package configurator

import (
	"darkbot/configurator/models"
	"darkbot/settings/logus"
	"darkbot/settings/types"
	"darkbot/settings/utils"
	"fmt"
)

type IConfiguratorTags interface {
	TagsAdd(channelID types.DiscordChannelID, tags ...string) *ConfiguratorError
	TagsRemove(channelID types.DiscordChannelID, tags ...string) *ConfiguratorError
	TagsList(channelID types.DiscordChannelID) ([]string, *ConfiguratorError)
	TagsClear(channelID types.DiscordChannelID) *ConfiguratorError
}

type taggable interface {
	models.TagBase |
		models.TagSystem |
		models.TagRegion |
		models.TagPlayerFriend |
		models.TagPlayerEnemy
	GetTag() string
}

type ConfiguratorTags[T taggable] struct {
	Configurator
}

type ConfiguratorBase = ConfiguratorTags[models.TagBase]
type ConfiguratorSystem = ConfiguratorTags[models.TagSystem]
type ConfiguratorRegion = ConfiguratorTags[models.TagRegion]
type ConfiguratorPlayerFriend = ConfiguratorTags[models.TagPlayerFriend]
type ConfiguratorPlayerEnemy = ConfiguratorTags[models.TagPlayerEnemy]

// T =
func (c ConfiguratorTags[T]) TagsAdd(channelID types.DiscordChannelID, tags ...string) *ConfiguratorError {
	objs := []T{}
	errors := &ConfiguratorError{}
	for _, tag := range tags {
		objs = append(objs, T{
			TagTemplate: models.TagTemplate{
				ChannelID: channelID,
				Tag:       tag,
			},
		})
	}

	presentTags, _ := c.TagsList(channelID)
	for _, tag := range presentTags {
		for _, newtag := range tags {
			if tag == newtag {
				fmt.Printf("TagsAdd. Tag %s is already present in channelID=%s\n", tag, channelID)
				errors.AppenError(StorageErrorExists{items: []string{tag}})
				return errors
			}
		}
	}

	res := c.db.Create(objs)
	logus.CheckWarn(res.Error, "unsuccesful result of c.db.Create")
	errors.AppendSQLError(res)
	return errors
}

func (c ConfiguratorTags[T]) TagsRemove(channelID types.DiscordChannelID, tags ...string) *ConfiguratorError {
	errors := &ConfiguratorError{}
	for _, tag := range tags {
		result := c.db.Where("channel_id = ? AND tag = ?", channelID, tag).Delete(&T{})
		logus.CheckWarn(result.Error, "unsuccesful result of c.db.Delete")
		errors.AppendSQLError(result)
	}
	return errors
}

func (c ConfiguratorTags[T]) TagsList(channelID types.DiscordChannelID) ([]string, *ConfiguratorError) {
	objs := []T{}
	result := c.db.Where("channel_id = ?", channelID).Find(&objs)
	logus.CheckWarn(result.Error, "unsuccesful result of c.db.Find")

	return utils.CompL(objs,
		func(x T) string { return x.GetTag() }), (&ConfiguratorError{}).AppendSQLError(result)
}

func (c ConfiguratorTags[T]) TagsClear(channelID types.DiscordChannelID) *ConfiguratorError {
	tags := []T{}
	result := c.db.Unscoped().Where("channel_id = ?", channelID).Find(&tags)
	fmt.Println("Clear.Find.rowsAffected=", result.RowsAffected)
	fmt.Println("Clear.Find.err=", result.Error)
	result = c.db.Unscoped().Delete(&tags)
	fmt.Println("Clear.Delete.rowsAffected=", result.RowsAffected)
	fmt.Println("Clear.Delete.err=", result.Error)
	return (&ConfiguratorError{}).AppendSQLError(result)
}
