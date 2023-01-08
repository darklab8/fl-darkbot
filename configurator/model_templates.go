package configurator

import "gorm.io/gorm"

type FKChannel struct {
	ChannelID uint
	Channel   Channel `gorm:"foreignKey:ChannelID"`
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
