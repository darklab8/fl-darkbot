package scrappy

import (
	"time"

	"github.com/darklab8/fl-darkbot/app/scrappy/base"
	"github.com/darklab8/fl-darkbot/app/scrappy/baseattack"
	"github.com/darklab8/fl-darkbot/app/settings"
	"github.com/darklab8/fl-darkbot/app/settings/logus"
)

type ScrappyStorage struct {
	baseStorage       *base.BaseStorage
	baseAttackStorage *baseattack.BaseAttackStorage
}

func NewScrapyStorage(base_api base.IbaseAPI, base_attack baseattack.IbaseAttackAPI, opts ...StorageParam) *ScrappyStorage {
	s := &ScrappyStorage{}
	s.baseStorage = base.NewBaseStorage(base_api)
	s.baseAttackStorage = baseattack.NewBaseAttackStorage(base_attack)

	for _, opt := range opts {
		opt(s)
	}
	return s
}

func (s *ScrappyStorage) Update() {
	s.baseStorage.Update()
	s.baseAttackStorage.Update()
}

func (s *ScrappyStorage) GetBaseStorage() *base.BaseStorage {
	return s.baseStorage
}

func (s *ScrappyStorage) GetBaseAttackStorage() *baseattack.BaseAttackStorage {
	return s.baseAttackStorage
}

func NewScrappyWithApis() *ScrappyStorage {
	return NewScrapyStorage(base.NewBaseApi(), baseattack.NewBaseAttackAPI())
}

func (s *ScrappyStorage) Run() {
	logus.Log.Info("initialized scrappy")
	logus.Log.Info("starting scrappy infinity update loop")
	for {
		s.Update()
		time.Sleep(time.Duration(settings.Env.ScrappyLoopDelay) * time.Second)
	}
	logus.Log.Info("gracefully shutdown scrappy infinity loop")
}
