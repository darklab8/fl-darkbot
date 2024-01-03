package configurator

import (
	"darkbot/app/configurator/models"
	"darkbot/app/settings/logus"
	"darkbot/app/settings/types"

	"github.com/darklab8/darklab_goutils/goutils/utils"
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
		models.TagPlayerEnemy |
		models.TagPlayerEvent |
		models.TagForumPostTrack |
		models.TagForumPostIgnore |
		models.TagForumSubforumTrack |
		models.TagForumSubforumIgnore
	GetTag() types.Tag
}

type ConfiguratorTags[T taggable] struct {
	*Configurator
}

func NewConfiguratorTags[T taggable](configurator *Configurator) ConfiguratorTags[T] {
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

type ConfiguratorPlayerEvent = ConfiguratorTags[models.TagPlayerEvent]

var NewConfiguratorPlayerEvent = NewConfiguratorTags[models.TagPlayerEvent]

type ConfiguratorForumWatch = ConfiguratorTags[models.TagForumPostTrack]

var NewConfiguratorForumWatch = NewConfiguratorTags[models.TagForumPostTrack]

type ConfiguratorForumIgnore = ConfiguratorTags[models.TagForumPostIgnore]

var NewConfiguratorForumIgnore = NewConfiguratorTags[models.TagForumPostIgnore]

type ConfiguratorSubForumWatch = ConfiguratorTags[models.TagForumSubforumTrack]

var NewConfiguratorSubForumWatch = NewConfiguratorTags[models.TagForumSubforumTrack]

type ConfiguratorSubForumIgnore = ConfiguratorTags[models.TagForumSubforumIgnore]

var NewConfiguratorSubForumIgnore = NewConfiguratorTags[models.TagForumSubforumIgnore]

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
				logus.Log.Info("TagsAdd. Tag %s is already present in channelID=%s\n", logus.Tag(tag), logus.ChannelID(channelID))
				return StorageErrorExists{items: []string{string(tag)}}
			}
		}
	}

	res := c.db.Create(objs)
	logus.Log.CheckWarn(res.Error, "unsuccesful result of c.db.Create")
	return res.Error
}

func (c ConfiguratorTags[T]) TagsRemove(channelID types.DiscordChannelID, tags ...types.Tag) error {
	errors := NewErrorAggregator()
	TotalRowsAffected := 0
	for _, tag := range tags {
		result := c.db.Where("channel_id = ? AND tag = ?", channelID, tag).Delete(&T{})
		logus.Log.CheckWarn(result.Error, "unsuccesful result of c.db.Delete")
		errors.Append(result.Error)
		TotalRowsAffected += int(result.RowsAffected)
	}

	if TotalRowsAffected == 0 {
		return ErrorZeroAffectedRows{}
	}

	return errors.TryToGetError()
}

func (c ConfiguratorTags[T]) TagsList(channelID types.DiscordChannelID) ([]types.Tag, error) {
	objs := []T{}
	result := c.db.Where("channel_id = ?", channelID).Find(&objs)
	logus.Log.CheckWarn(result.Error, "unsuccesful result of c.db.Find")

	tags := utils.CompL(objs,
		func(x T) types.Tag { return x.GetTag() })

	if result.RowsAffected == 0 {
		return tags, ErrorZeroAffectedRows{}
	}

	return tags, result.Error
}

func (c ConfiguratorTags[T]) TagsClear(channelID types.DiscordChannelID) error {
	tags := []T{}
	result := c.db.Unscoped().Where("channel_id = ?", channelID).Find(&tags)
	if len(tags) == 0 {
		return ErrorZeroAffectedRows{ExtraMsg: "no tags found"}
	}
	logus.Log.Debug("Clear.Find", logus.GormResult(result))
	result = c.db.Unscoped().Delete(&tags)
	logus.Log.Debug("Clear.Detete", logus.GormResult(result))
	return result.Error
}
