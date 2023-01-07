package scrappy

import (
	"darkbot/scrappy/base"
	"darkbot/scrappy/player"
	"darkbot/utils"
	"time"
)

type ScrappyStorage struct {
	baseStorage   base.BaseStorage
	playerStorage player.PlayerStorage
}

func (s *ScrappyStorage) New() *ScrappyStorage {
	s.baseStorage = base.BaseStorage{}.New()
	s.playerStorage = player.PlayerStorage{}.New()
	return s
}

func (s *ScrappyStorage) Update() {
	s.baseStorage.Update()
	s.playerStorage.Update()
}

var storage *ScrappyStorage

func init() {
	utils.LogInfo("initialized scrappy")
	storage = (&ScrappyStorage{}).New()
}

func Run() {
	utils.LogInfo("starting scrappy infinity update loop")
	for {
		storage.Update()
		time.Sleep(10 * time.Second)
	}
	utils.LogInfo("gracefully shutdown scrappy infinity loop")
}
