package base

import (
	"darkbot/scrappy/shared/api"
	"darkbot/settings"
	"darkbot/settings/types"
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
