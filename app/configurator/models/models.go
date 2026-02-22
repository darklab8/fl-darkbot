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
type TagPoBGood struct {
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

type TagForumAuthorWatch struct {
	TagTemplate
}

type TagForumAuthorIgnore struct {
	TagTemplate
}

// =========== Alerts ===============

type AlertTresholdInteger struct {
	Threshold int `gorm:"check:threshold > 0; check:threshold >= 0"`
}

func (t AlertTresholdInteger) GetThreshold() int {
	return t.Threshold
}

type AlertPoBGood struct {
	GoodNickname string
}

func (t AlertPoBGood) GetGoodNickname() string {
	return t.GoodNickname
}

type AlertPobGoodBelowThan struct {
	MultiValueTemplate
	AlertTresholdInteger
	ThresholdIntegerKind `gorm:"-"`
	AlertPoBGood
}
type AlertPobGoodAboveThan struct {
	MultiValueTemplate
	AlertTresholdInteger
	ThresholdIntegerKind `gorm:"-"`
	AlertPoBGood
}

// ====== Shared alerts for all bases =========
type AlertBaseHealthLowerThan struct {
	OneValueTemplate
	AlertTresholdInteger
	ThresholdIntegerKind `gorm:"-"`
}

type AlertBaseIfHealthDecreasing struct {
	OneValueTemplate
}
type AlertBaseIfUnderAttack struct {
	OneValueTemplate
}

type AlertBaseMoneyBelow struct {
	OneValueTemplate
	AlertTresholdInteger
	ThresholdIntegerKind `gorm:"-"`
}

type AlertBaseCargoBelow struct {
	OneValueTemplate
	AlertTresholdInteger
	ThresholdIntegerKind `gorm:"-"`
}

type ThresholdIntegerKind int64

const (
	ThresholdIntegerPercentage ThresholdIntegerKind = iota
	ThresholdIntegerNotConstrained
)

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
