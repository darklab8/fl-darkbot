package scrappy

import (
	"darkbot/app/scrappy/base"
	"darkbot/app/scrappy/baseattack"
	"darkbot/app/scrappy/player"
)

func FixtureNewStorageWithPlayers(players *player.PlayerStorage) *ScrappyStorage {
	return &ScrappyStorage{playerStorage: players}
}

func FixtureMockedStorage() *ScrappyStorage {
	return NewScrapyStorage(
		base.FixtureBaseApiMock(),
		player.FixturePlayerAPIMock(),
		baseattack.FixtureBaseAttackAPIMock(),
	)
}
