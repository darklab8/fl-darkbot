package models

import "gorm.io/gorm"

type TagTemplate struct {
	gorm.Model
	Channel Channel `gorm:"references:ChannelID"`
}

type AlertTemplate struct {
	gorm.Model
	ChannelID string
	Channel   Channel `gorm:"references:ChannelID"`
}
