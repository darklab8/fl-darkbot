package models

import (
	"darkbot/settings/types"

	"gorm.io/gorm"
)

type TagTemplate struct {
	gorm.Model
	ChannelID types.DiscordChannelID
	Channel   Channel `gorm:"references:ChannelID"`
	Tag       string
}

func (t TagTemplate) GetTag() string {
	return t.Tag
}

type AlertPlayerMoreThan struct {
	PlayersMoreThan int
}

type AlertTemplate struct {
	gorm.Model
	ChannelID types.DiscordChannelID
	Channel   Channel `gorm:"references:ChannelID,unique"`
}
