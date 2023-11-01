package models

import (
	"darkbot/settings/types"

	"gorm.io/gorm"
)

type Channel struct {
	gorm.DeletedAt
	ChannelID types.DiscordChannelID `gorm:"primarykey"`
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

type AlertTresholdShared struct {
	Threshold int `gorm:"check:threshold > 0; check:threshold <= 100"`
}

func (t AlertTresholdShared) GetThreshold() int {
	return t.Threshold
}

type AlertNeutralPlayersEqualOrGreater struct {
	AlertTemplate
	AlertTresholdShared
}

type AlertEnemiesEqualOrGreater struct {
	AlertTemplate
	AlertTresholdShared
}

type AlertFriendsEqualOrGreater struct {
	AlertTemplate
	AlertTresholdShared
}

// ====== Shared alerts for all bases =========
type AlertBaseHealthLowerThan struct {
	AlertTemplate
	AlertTresholdShared
}

type AlertBaseIfHealthDecreasing struct {
	AlertTemplate
}
type AlertBaseIfUnderAttack struct {
	AlertTemplate
}

type AlertPingMessage struct {
	AlertTemplate
	Value string
}

func (a AlertPingMessage) GetValue() string {
	return a.Value
}
