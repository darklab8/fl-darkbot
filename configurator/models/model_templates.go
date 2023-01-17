package models

import "gorm.io/gorm"

type FKChannel struct {
	ChannelID string  `gorm:"uniqueIndex:idx_tag_per_channel"`
	Channel   Channel `gorm:"references:ChannelID"`
}

type TagTemplate struct {
	gorm.Model
	FKChannel
	Tag string `gorm:"uniqueIndex:idx_tag_per_channel"`
}

type AlertTemplate struct {
	gorm.Model
	FKChannel
	Threshold int
}
