package models

import "gorm.io/gorm"

type Channel struct {
	gorm.Model
	Ref string `gorm:"index"` // Discord channel reference
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
}

type AlertPlayerEnemy struct {
	AlertTemplate
}

type AlertPlayerFriend struct {
	AlertTemplate
}

// Shared alerts for all bases
type AlertBase struct {
	FKChannel
	HealthIsLowerThan  float64
	HealthIsDecreasing float64
}
