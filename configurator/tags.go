package configurator

import (
	"darkbot/configurator/models"
	"darkbot/settings/utils"
	"darkbot/settings/utils/logger"
	"fmt"
)

type IConfiguratorTags interface {
	TagsAdd(channelID string, tags ...string) *ConfiguratorError
	TagsRemove(channelID string, tags ...string) *ConfiguratorError
	TagsList(channelID string) ([]string, *ConfiguratorError)
	TagsClear(channelID string) *ConfiguratorError
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
func (c ConfiguratorTags[T]) TagsAdd(channelID string, tags ...string) *ConfiguratorError {
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
	logger.CheckWarn(res.Error, "unsuccesful result of c.db.Create")
	errors.AppendSQLError(res)
	return errors
}

func (c ConfiguratorTags[T]) TagsRemove(channelID string, tags ...string) *ConfiguratorError {
	errors := &ConfiguratorError{}
	for _, tag := range tags {
		result := c.db.Where("channel_id = ? AND tag = ?", channelID, tag).Delete(&T{})
		logger.CheckWarn(result.Error, "unsuccesful result of c.db.Delete")
		errors.AppendSQLError(result)
	}
	return errors
}

func (c ConfiguratorTags[T]) TagsList(channelID string) ([]string, *ConfiguratorError) {
	objs := []T{}
	result := c.db.Where("channel_id = ?", channelID).Find(&objs)
	logger.CheckWarn(result.Error, "unsuccesful result of c.db.Find")

	return utils.CompL(objs,
		func(x T) string { return x.GetTag() }), (&ConfiguratorError{}).AppendSQLError(result)
}

func (c ConfiguratorTags[T]) TagsClear(channelID string) *ConfiguratorError {
	tags := []T{}
	result := c.db.Unscoped().Where("channel_id = ?", channelID).Find(&tags)
	fmt.Println("Clear.Find.rowsAffected=", result.RowsAffected)
	fmt.Println("Clear.Find.err=", result.Error)
	result = c.db.Unscoped().Delete(&tags)
	fmt.Println("Clear.Delete.rowsAffected=", result.RowsAffected)
	fmt.Println("Clear.Delete.err=", result.Error)
	return (&ConfiguratorError{}).AppendSQLError(result)
}
