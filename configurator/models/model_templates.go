package models

import "gorm.io/gorm"

type ChannelShared struct {
	ChannelID string
	Channel   Channel `gorm:"references:ChannelID"`
}

type TagTemplate struct {
	gorm.Model
	ChannelShared
	Tag string
}

type AlertTemplate struct {
	gorm.Model
	ChannelShared
}
