package models

import (
	"darkbot/app/settings/types"

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

type AlertPlayerMoreThan struct {
	PlayersMoreThan int
}

type AlertTemplate struct {
	gorm.Model
	ChannelID types.DiscordChannelID
	Channel   Channel `gorm:"references:ChannelID,unique"`
}
