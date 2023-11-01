package scrappy

import (
	"darkbot/scrappy/base"
	"darkbot/scrappy/baseattack"
	"darkbot/scrappy/player"
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

func FixtureSetBaseStorageAPI(base_api base.IbaseAPI) {
	Storage.baseStorage.FixtureSetAPI(base_api)
}
