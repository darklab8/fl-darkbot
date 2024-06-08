package base

import (
	"github.com/darklab8/fl-darkbot/app/scrappy/shared/api"
	"github.com/darklab8/fl-darkbot/app/settings"
	"github.com/darklab8/fl-darkbot/app/settings/types"
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
	b.url = types.APIurl(settings.Env.ScrappyBaseUrl)
	return b
}
