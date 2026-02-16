package models

import (
	"github.com/darklab8/fl-darkbot/app/settings/types"

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

type TagForumSubforumTrack struct {
	TagTemplate
}

type TagForumSubforumIgnore struct {
	TagTemplate
}

type TagForumContentWatch struct {
	TagTemplate
}

type TagForumContentIgnore struct {
	TagTemplate
}

// =========== Alerts ===============

type AlertTresholdShared struct {
	Threshold int `gorm:"check:threshold > 0; check:threshold <= 100"`
}

func (t AlertTresholdShared) GetThreshold() int {
	return t.Threshold
}

// ====== Shared alerts for all bases =========
type AlertBaseHealthLowerThan struct {
	OneValueTemplate
	AlertTresholdShared
}

type AlertBaseIfHealthDecreasing struct {
	OneValueTemplate
}
type AlertBaseIfUnderAttack struct {
	OneValueTemplate
}

type AlertPingMessage struct {
	OneValueTemplate
	Value string
}

func (a AlertPingMessage) GetValue() string {
	return a.Value
}

// ====== Configs =========

type ConfigBaseOrderingKey struct {
	OneValueTemplate
	Value string
}

// i know it can be Constraint. But if i add it as `struct tag` it breaks typing
// and kind of hard to figure out how to fix nicely
const (
	BaseKeyName        types.OrderKey = "name"
	BaseKeyAffiliation types.OrderKey = "affiliation"
)

var ConfigBaseOrderingKeyAllowedTags = []types.OrderKey{BaseKeyName, BaseKeyAffiliation}

func (a ConfigBaseOrderingKey) GetValue() string {
	return a.Value
}
