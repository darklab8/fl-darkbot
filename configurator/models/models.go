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

type AlertNeutralPlayersEqualOrGreater struct {
	AlertTemplate
	Threshold int `gorm:"check:threshold > 0"`
}

type AlertEnemiesEqualOrGreater struct {
	AlertTemplate
	Threshold int `gorm:"check:threshold > 0"`
}

type AlertFriendsEqualOrGreater struct {
	AlertTemplate
	Threshold int `gorm:"check:threshold > 0"`
}

// Shared alerts for all bases
type AlertBaseHealthLowerThan struct {
	AlertTemplate
	Threshold int `gorm:"check:threshold > 0; check:threshold <= 100"`
}

func (t AlertNeutralPlayersEqualOrGreater) GetThreshold() int {
	return t.Threshold
}
func (t AlertEnemiesEqualOrGreater) GetThreshold() int {
	return t.Threshold
}
func (t AlertFriendsEqualOrGreater) GetThreshold() int {
	return t.Threshold
}
func (t AlertBaseHealthLowerThan) GetThreshold() int {
	return t.Threshold
}

type alertThreshold interface {
	AlertNeutralPlayersEqualOrGreater |
		AlertEnemiesEqualOrGreater |
		AlertFriendsEqualOrGreater |
		AlertBaseHealthLowerThan
}

func (t AlertNeutralPlayersEqualOrGreater) SetThreshold(channelID string, value int) AlertNeutralPlayersEqualOrGreater {
	t.ChannelID = channelID
	t.Threshold = value
	return t
}
func (t AlertEnemiesEqualOrGreater) SetThreshold(channelID string, value int) AlertEnemiesEqualOrGreater {
	t.ChannelID = channelID
	t.Threshold = value
	return t
}
func (t AlertFriendsEqualOrGreater) SetThreshold(channelID string, value int) AlertFriendsEqualOrGreater {
	t.ChannelID = channelID
	t.Threshold = value
	return t
}
func (t AlertBaseHealthLowerThan) SetThreshold(channelID string, value int) AlertBaseHealthLowerThan {
	t.ChannelID = channelID
	t.Threshold = value
	return t
}

type AlertBaseIfHealthDecreasing struct {
	AlertTemplate
}
type AlertBaseIfUnderAttack struct {
	AlertTemplate
}
