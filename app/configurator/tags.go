package configurator

import (
	"darkbot/app/configurator/models"
	"darkbot/app/settings/logus"
	"darkbot/app/settings/types"
	"darkbot/app/settings/utils"
	"fmt"
)

type IConfiguratorTags interface {
	TagsAdd(channelID types.DiscordChannelID, tags ...string) error
	TagsRemove(channelID types.DiscordChannelID, tags ...string) error
	TagsList(channelID types.DiscordChannelID) ([]string, error)
	TagsClear(channelID types.DiscordChannelID) error
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

func NewConfiguratorTags[T taggable](configurator Configurator) ConfiguratorTags[T] {
	t := ConfiguratorTags[T]{Configurator: configurator}
	return t
}

type ConfiguratorBase = ConfiguratorTags[models.TagBase]

var NewConfiguratorBase = NewConfiguratorTags[models.TagBase]

type ConfiguratorSystem = ConfiguratorTags[models.TagSystem]

var NewConfiguratorSystem = NewConfiguratorTags[models.TagSystem]

type ConfiguratorRegion = ConfiguratorTags[models.TagRegion]

var NewConfiguratorRegion = NewConfiguratorTags[models.TagRegion]

type ConfiguratorPlayerFriend = ConfiguratorTags[models.TagPlayerFriend]

var NewConfiguratorPlayerFriend = NewConfiguratorTags[models.TagPlayerFriend]

type ConfiguratorPlayerEnemy = ConfiguratorTags[models.TagPlayerEnemy]

var NewConfiguratorPlayerEnemy = NewConfiguratorTags[models.TagPlayerEnemy]

func (c ConfiguratorTags[T]) TagsAdd(channelID types.DiscordChannelID, tags ...string) error {
	objs := []T{}
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
				return StorageErrorExists{items: []string{tag}}
			}
		}
	}

	res := c.db.Create(objs)
	logus.CheckWarn(res.Error, "unsuccesful result of c.db.Create")
	return res.Error
}

func (c ConfiguratorTags[T]) TagsRemove(channelID types.DiscordChannelID, tags ...string) error {
	errors := NewErrorAggregator()
	for _, tag := range tags {
		result := c.db.Where("channel_id = ? AND tag = ?", channelID, tag).Delete(&T{})
		logus.CheckWarn(result.Error, "unsuccesful result of c.db.Delete")
		errors.Append(result.Error)
	}
	return errors.TryToGetError()
}

func (c ConfiguratorTags[T]) TagsList(channelID types.DiscordChannelID) ([]string, error) {
	objs := []T{}
	result := c.db.Where("channel_id = ?", channelID).Find(&objs)
	logus.CheckWarn(result.Error, "unsuccesful result of c.db.Find")

	return utils.CompL(objs,
		func(x T) string { return x.GetTag() }), result.Error
}

func (c ConfiguratorTags[T]) TagsClear(channelID types.DiscordChannelID) error {
	tags := []T{}
	result := c.db.Unscoped().Where("channel_id = ?", channelID).Find(&tags)
	fmt.Println("Clear.Find.rowsAffected=", result.RowsAffected)
	fmt.Println("Clear.Find.err=", result.Error)
	result = c.db.Unscoped().Delete(&tags)
	fmt.Println("Clear.Delete.rowsAffected=", result.RowsAffected)
	fmt.Println("Clear.Delete.err=", result.Error)
	return result.Error
}
