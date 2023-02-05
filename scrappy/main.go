package scrappy

import (
	"darkbot/scrappy/base"
	"darkbot/scrappy/player"
	"darkbot/utils/logger"
	"time"
)

type ScrappyStorage struct {
	BaseStorage   *base.BaseStorage
	PlayerStorage *player.PlayerStorage
}

func (s *ScrappyStorage) New() *ScrappyStorage {
	s.BaseStorage = (&base.BaseStorage{}).New()
	s.PlayerStorage = (&player.PlayerStorage{}).New()
	return s
}

func (s *ScrappyStorage) Update() {
	s.BaseStorage.Update()
	s.PlayerStorage.Update()
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
		time.Sleep(10 * time.Second)
	}
	logger.Info("gracefully shutdown scrappy infinity loop")
}
