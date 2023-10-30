package scrappy

import (
	"darkbot/scrappy/base"
	"darkbot/scrappy/baseattack"
	"darkbot/scrappy/player"
	"darkbot/scrappy/shared/api"
	"darkbot/settings"
	"darkbot/settings/logus"
	"time"
)

type ScrappyStorage struct {
	baseStorage       *base.BaseStorage
	playerStorage     *player.PlayerStorage
	baseAttackStorage *baseattack.BaseAttackStorage
}

func NewScrapyStorage(base_api api.APIinterface, player_api api.APIinterface, base_attack api.APIinterface) *ScrappyStorage {
	s := &ScrappyStorage{}
	s.baseStorage = base.NewBaseStorage(base_api)
	s.playerStorage = player.NewPlayerStorage(player_api)
	s.baseAttackStorage = baseattack.NewBaseAttackStorage(base_attack)
	return s
}

func (s *ScrappyStorage) Update() {
	s.baseStorage.Update()
	s.playerStorage.Update()
	s.baseAttackStorage.Update()
}

func (s *ScrappyStorage) GetBaseStorage() *base.BaseStorage {
	return s.baseStorage
}
func (s *ScrappyStorage) GetPlayerStorage() *player.PlayerStorage {
	return s.playerStorage
}

func (s *ScrappyStorage) GetBaseAttackStorage() *baseattack.BaseAttackStorage {
	return s.baseAttackStorage
}

var Storage *ScrappyStorage

func init() {
	logus.Info("initialized scrappy")
	Storage = NewScrapyStorage(base.NewBaseApi(), player.NewPlayerAPI(), base.NewBaseApi())
}

func Run() {
	logus.Info("starting scrappy infinity update loop")
	for {
		Storage.Update()
		time.Sleep(time.Duration(settings.LoopDelay) * time.Second)
	}
	logus.Info("gracefully shutdown scrappy infinity loop")
}
