package models

import (
	"github.com/darklab8/fl-darkbot/app/settings/types"

	"gorm.io/gorm"
)

type TagTemplate struct {
	gorm.Model
	ChannelID types.DiscordChannelID
	Channel   Channel `gorm:"references:ChannelID"`
	Tag       string
}

func (t TagTemplate) GetTag() types.Tag {
	return types.Tag(t.Tag)
}

type OneValueTemplate struct {
	gorm.Model
	ChannelID types.DiscordChannelID
	Channel   Channel `gorm:"references:ChannelID,unique"`
}

type MultiValueTemplate struct {
	gorm.Model
	ChannelID types.DiscordChannelID
	Channel   Channel `gorm:"references:ChannelID"`
}
