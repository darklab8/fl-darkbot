package scrappy

import (
	"time"

	"github.com/darklab/fl-darkbot/app/scrappy/base"
	"github.com/darklab/fl-darkbot/app/scrappy/baseattack"
	"github.com/darklab/fl-darkbot/app/scrappy/player"
	"github.com/darklab/fl-darkbot/app/settings"
	"github.com/darklab/fl-darkbot/app/settings/logus"
)

type ScrappyStorage struct {
	baseStorage       *base.BaseStorage
	playerStorage     *player.PlayerStorage
	baseAttackStorage *baseattack.BaseAttackStorage
}

func NewScrapyStorage(base_api base.IbaseAPI, player_api player.IPlayerAPI, base_attack baseattack.IbaseAttackAPI, opts ...storageParam) *ScrappyStorage {
	s := &ScrappyStorage{}
	s.baseStorage = base.NewBaseStorage(base_api)
	s.playerStorage = player.NewPlayerStorage(player_api)
	s.baseAttackStorage = baseattack.NewBaseAttackStorage(base_attack)

	for _, opt := range opts {
		opt(s)
	}
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

func NewScrappyWithApis() *ScrappyStorage {
	return NewScrapyStorage(base.NewBaseApi(), player.NewPlayerAPI(), baseattack.NewBaseAttackAPI())
}

func (s *ScrappyStorage) Run() {
	logus.Log.Info("initialized scrappy")
	logus.Log.Info("starting scrappy infinity update loop")
	for {
		s.Update()
		time.Sleep(time.Duration(settings.ScrappyLoopDelay) * time.Second)
	}
	logus.Log.Info("gracefully shutdown scrappy infinity loop")
}
