package baseattack

import (
	"darkbot/scrappy/shared/api"
	"darkbot/settings"
	"darkbot/settings/utils/logger"
)

type basesattackAPI struct {
	api.APIrequest
}

func NewBaseAttackAPI() api.APIinterface {
	a := basesattackAPI{}
	a.Init(settings.Config.ScrappyBaseAttackUrl)
	return a
}

type BaseAttackStorage struct {
	data BaseAttackData
	api  api.APIinterface
}

func NewBaseAttackStorage(api api.APIinterface) *BaseAttackStorage {
	b := &BaseAttackStorage{}
	b.api = api
	return b
}

type BaseAttackData string

func (b *BaseAttackStorage) GetData() BaseAttackData { return BaseAttackData(b.data) }

func (b *BaseAttackStorage) Update() {
	data, err := b.api.GetData()
	if err != nil {
		logger.CheckWarn(err, "quering API with error in BaseAttackStorage")
		return
	}
	b.data = BaseAttackData(string((data)))
}
