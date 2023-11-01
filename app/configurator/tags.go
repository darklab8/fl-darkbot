package configurator

import (
	"darkbot/app/configurator/models"
	"darkbot/app/settings/logus"
	"darkbot/app/settings/types"
	"darkbot/app/settings/utils"
)

type IConfiguratorTags interface {
	TagsAdd(channelID types.DiscordChannelID, tags ...types.Tag) error
	TagsRemove(channelID types.DiscordChannelID, tags ...types.Tag) error
	TagsList(channelID types.DiscordChannelID) ([]types.Tag, error)
	TagsClear(channelID types.DiscordChannelID) error
}

type taggable interface {
	models.TagBase |
		models.TagSystem |
		models.TagRegion |
		models.TagPlayerFriend |
		models.TagPlayerEnemy
	GetTag() types.Tag
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

func (c ConfiguratorTags[T]) TagsAdd(channelID types.DiscordChannelID, tags ...types.Tag) error {
	objs := []T{}
	for _, tag := range tags {
		objs = append(objs, T{
			TagTemplate: models.TagTemplate{
				ChannelID: channelID,
				Tag:       string(tag),
			},
		})
	}

	presentTags, _ := c.TagsList(channelID)
	for _, tag := range presentTags {
		for _, newtag := range tags {
			if tag == newtag {
				logus.Info("TagsAdd. Tag %s is already present in channelID=%s\n", logus.Tag(tag), logus.ChannelID(channelID))
				return StorageErrorExists{items: []string{string(tag)}}
			}
		}
	}

	res := c.db.Create(objs)
	logus.CheckWarn(res.Error, "unsuccesful result of c.db.Create")
	return res.Error
}

func (c ConfiguratorTags[T]) TagsRemove(channelID types.DiscordChannelID, tags ...types.Tag) error {
	errors := NewErrorAggregator()
	for _, tag := range tags {
		result := c.db.Where("channel_id = ? AND tag = ?", channelID, tag).Delete(&T{})
		logus.CheckWarn(result.Error, "unsuccesful result of c.db.Delete")
		errors.Append(result.Error)
	}
	return errors.TryToGetError()
}

func (c ConfiguratorTags[T]) TagsList(channelID types.DiscordChannelID) ([]types.Tag, error) {
	objs := []T{}
	result := c.db.Where("channel_id = ?", channelID).Find(&objs)
	logus.CheckWarn(result.Error, "unsuccesful result of c.db.Find")

	return utils.CompL(objs,
		func(x T) types.Tag { return x.GetTag() }), result.Error
}

func (c ConfiguratorTags[T]) TagsClear(channelID types.DiscordChannelID) error {
	tags := []T{}
	result := c.db.Unscoped().Where("channel_id = ?", channelID).Find(&tags)
	logus.Debug("Clear.Find", logus.GormResult(result))
	result = c.db.Unscoped().Delete(&tags)
	logus.Debug("Clear.Detete", logus.GormResult(result))
	return result.Error
}
