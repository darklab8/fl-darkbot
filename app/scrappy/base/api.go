package base

import (
	"darkbot/app/scrappy/shared/api"
	"darkbot/app/settings"
	"darkbot/app/settings/types"
)

type basesAPI struct {
	api.APIrequest
	url types.APIurl
}

func (b basesAPI) GetBaseData() ([]byte, error) {
	return b.GetData(b.url)
}

type IbaseAPI interface {
	GetBaseData() ([]byte, error)
}

func NewBaseApi() IbaseAPI {
	b := basesAPI{}
	b.url = types.APIurl(settings.Config.ScrappyBaseUrl)
	return b
}
