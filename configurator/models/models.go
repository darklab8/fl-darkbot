package models

import "gorm.io/gorm"

type Channel struct {
	gorm.DeletedAt
	ChannelID string `gorm:"primarykey"` // Discord channel reference
}

// ================== Tag Tracking ====================

type TagBase struct {
	TagTemplate
}

type TagPlayerFriend struct {
	TagTemplate
}

type TagPlayerEnemy struct {
	TagTemplate
}

type TagSystem struct {
	TagTemplate
}

type TagRegion struct {
	TagTemplate
}

type TagForumPostTrack struct {
	TagTemplate
}

type TagForumPostIgnore struct {
	TagTemplate
}

// =========== Alerts ===============

type AlertPlayerUnrecognized struct {
	AlertTemplate
	Threshold int
}

type AlertPlayerEnemy struct {
	AlertTemplate
	Threshold int
}

type AlertPlayerFriend struct {
	AlertTemplate
	Threshold int
}

// Shared alerts for all bases
type AlertBase struct {
	AlertTemplate
	HealthIsLowerThan  float64
	HealthIsDecreasing float64
}
