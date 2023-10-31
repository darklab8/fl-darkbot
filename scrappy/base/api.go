package base

import (
	"darkbot/scrappy/shared/api"
	"darkbot/settings"
)

type basesAPI struct {
	api.APIrequest
	url api.APIurl
}

func (b basesAPI) GetBaseData() ([]byte, error) {
	return b.GetData(b.url)
}

type IbaseAPI interface {
	GetBaseData() ([]byte, error)
}

func NewBaseApi() IbaseAPI {
	b := basesAPI{}
	b.url = api.APIurl(settings.Config.ScrappyBaseUrl)
	return b
}
