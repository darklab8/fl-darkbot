package scrappy

import (
	"darkbot/scrappy/base"
	"darkbot/scrappy/player"
	"darkbot/utils"
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
	utils.LogInfo("initialized scrappy")
	Storage = (&ScrappyStorage{}).New()
}

func Run() {
	utils.LogInfo("starting scrappy infinity update loop")
	for {
		Storage.Update()
		time.Sleep(10 * time.Second)
	}
	utils.LogInfo("gracefully shutdown scrappy infinity loop")
}
