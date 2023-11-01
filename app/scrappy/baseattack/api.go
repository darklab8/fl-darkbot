package baseattack

import (
	"darkbot/app/scrappy/shared/api"
	"darkbot/app/settings"
	"darkbot/app/settings/logus"
	"darkbot/app/settings/types"
)

type basesattackAPI struct {
	api.APIrequest
	url types.APIurl
}

func (b basesattackAPI) GetBaseAttackData() ([]byte, error) {
	return b.GetData(b.url)
}

type IbaseAttackAPI interface {
	GetBaseAttackData() ([]byte, error)
}

func NewBaseAttackAPI() IbaseAttackAPI {
	a := basesattackAPI{}
	a.url = types.APIurl(settings.Config.ScrappyBaseAttackUrl)
	return a
}

type BaseAttackStorage struct {
	data BaseAttackData
	api  IbaseAttackAPI
}

func NewBaseAttackStorage(api IbaseAttackAPI) *BaseAttackStorage {
	b := &BaseAttackStorage{}
	b.api = api
	return b
}

type BaseAttackData string

func (b *BaseAttackStorage) GetData() BaseAttackData { return BaseAttackData(b.data) }

func (b *BaseAttackStorage) Update() {
	data, err := b.api.GetBaseAttackData()
	if logus.CheckWarn(err, "quering API with error in BaseAttackStorage") {
		return
	}
	b.data = BaseAttackData(string((data)))
}
