package scrappy

import (
	"github.com/darklab8/fl-darkbot/app/scrappy/base"
	"github.com/darklab8/fl-darkbot/app/scrappy/baseattack"
)

func FixtureMockedStorage(opts ...StorageParam) *ScrappyStorage {
	return NewScrapyStorage(
		base.FixtureBaseApiMock(),
		baseattack.FixtureBaseAttackAPIMock(),
		opts...,
	)
}

type StorageParam func(storage *ScrappyStorage)
