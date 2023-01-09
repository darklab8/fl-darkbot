package models

import "gorm.io/gorm"

type FKChannel struct {
	ChannelID string
	Channel   Channel `gorm:"references:ChannelID"`
}

type TagTemplate struct {
	gorm.Model
	FKChannel
	Tag string
}

type AlertTemplate struct {
	gorm.Model
	FKChannel
	Threshold int
}
