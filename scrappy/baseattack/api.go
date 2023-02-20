package baseattack

import (
	"darkbot/scrappy/shared/api"
	"darkbot/settings"
)

type basesattackAPI struct {
	api.APIrequest
}

func (a basesattackAPI) New() api.APIinterface {
	a.Init(settings.Config.ScrappyBaseAttackUrl)
	return a
}

type BaseAttackStorage struct {
	Data string
	Api  api.APIinterface
}

func (b *BaseAttackStorage) New() *BaseAttackStorage {
	b.Api = basesattackAPI{}.New()
	return b
}

func (b *BaseAttackStorage) Update() {
	b.Data = string(b.Api.GetData())
}
