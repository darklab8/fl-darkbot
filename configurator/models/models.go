package models

type Channel struct {
	ChannelID string `gorm:"primarykey"` // Discord channel reference
}

// ================== Tag Tracking ====================

type TagBase struct {
	TagTemplate

	ChannelID string `gorm:"uniqueIndex:idx_1_tag_per_channel"`
	Tag       string `gorm:"uniqueIndex:idx_1_tag_per_channel"`
}

type TagPlayerFriend struct {
	TagTemplate

	ChannelID string `gorm:"uniqueIndex:idx_2_tag_per_channel"`
	Tag       string `gorm:"uniqueIndex:idx_2_tag_per_channel"`
}

type TagPlayerEnemy struct {
	TagTemplate

	ChannelID string `gorm:"uniqueIndex:idx_3_tag_per_channel"`
	Tag       string `gorm:"uniqueIndex:idx_3_tag_per_channel"`
}

type TagSystem struct {
	TagTemplate

	ChannelID string `gorm:"uniqueIndex:idx_4_tag_per_channel"`
	Tag       string `gorm:"uniqueIndex:idx_4_tag_per_channel"`
}

type TagRegion struct {
	TagTemplate

	ChannelID string `gorm:"uniqueIndex:idx_5_tag_per_channel"`
	Tag       string `gorm:"uniqueIndex:idx_5_tag_per_channel"`
}

type TagForumPostTrack struct {
	TagTemplate

	ChannelID string `gorm:"uniqueIndex:idx_6_tag_per_channel"`
	Tag       string `gorm:"uniqueIndex:idx_6_tag_per_channel"`
}

type TagForumPostIgnore struct {
	TagTemplate

	ChannelID string `gorm:"uniqueIndex:idx_7_tag_per_channel"`
	Tag       string `gorm:"uniqueIndex:idx_7_tag_per_channel"`
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
