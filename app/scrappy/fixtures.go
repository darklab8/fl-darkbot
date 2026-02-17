package scrappy

import (
	"github.com/darklab8/fl-darkbot/app/scrappy/base"
	"github.com/darklab8/fl-darkbot/app/scrappy/baseattack"
	"github.com/darklab8/fl-darkbot/app/scrappy/player"
)

func FixtureMockedStorage(opts ...storageParam) *ScrappyStorage {
	return NewScrapyStorage(
		base.FixtureBaseApiMock(),
		player.FixturePlayerAPIMock(),
		baseattack.FixtureBaseAttackAPIMock(),
		opts...,
	)
}

type storageParam func(storage *ScrappyStorage)

func WithPlayerStorage(playerStorage *player.PlayerStorage) storageParam {
	return func(storage *ScrappyStorage) {
		storage.playerStorage = playerStorage
	}
}
