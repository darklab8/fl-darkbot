package scrappy

import (
	"darkbot/scrappy/base"
	"darkbot/scrappy/baseattack"
	"darkbot/scrappy/player"
	"darkbot/scrappy/shared/api"
)

func FixtureNewStorage(players *player.PlayerStorage) *ScrappyStorage {
	return &ScrappyStorage{playerStorage: players}
}

func FixtureMockedStorage() *ScrappyStorage {
	return NewScrapyStorage(base.NewBaseApi(), player.NewPlayerMockAPI(), baseattack.NewBaseAttackAPIMock())
}

func FixtureSetBaseStorageAPI(base_api api.APIinterface) {
	Storage.baseStorage.FixtureSetAPI(base_api)
}
