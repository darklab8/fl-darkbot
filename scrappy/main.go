package scrappy

import (
	"darkbot/scrappy/base"
	"darkbot/scrappy/baseattack"
	"darkbot/scrappy/player"
	"darkbot/settings"
	"darkbot/settings/utils/logger"
	"time"
)

type ScrappyStorage struct {
	BaseStorage       *base.BaseStorage
	PlayerStorage     *player.PlayerStorage
	BaseAttackStorage *baseattack.BaseAttackStorage
}

func (s *ScrappyStorage) New() *ScrappyStorage {
	s.BaseStorage = (&base.BaseStorage{}).New()
	s.PlayerStorage = (&player.PlayerStorage{}).New()
	s.BaseAttackStorage = (&baseattack.BaseAttackStorage{}).New()
	return s
}

func (s *ScrappyStorage) Update() {
	s.BaseStorage.Update()
	s.PlayerStorage.Update()
	s.BaseAttackStorage.Update()
}

var Storage *ScrappyStorage

func init() {
	logger.Info("initialized scrappy")
	Storage = (&ScrappyStorage{}).New()
}

func Run() {
	logger.Info("starting scrappy infinity update loop")
	for {
		Storage.Update()
		time.Sleep(time.Duration(settings.LoopDelay) * time.Second)
	}
	logger.Info("gracefully shutdown scrappy infinity loop")
}
