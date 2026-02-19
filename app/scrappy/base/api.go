package base

import (
	"github.com/darklab8/fl-darkbot/app/scrappy/shared/api"
	"github.com/darklab8/fl-darkbot/app/settings"
	"github.com/darklab8/fl-darkbot/app/settings/types"
	"github.com/darklab8/fl-darkstat/darkapis/darkhttp"
	"github.com/darklab8/fl-darkstat/darkstat/configs_export"
)

type basesAPI struct {
	api.APIrequest
	url types.APIurl
}

func (b basesAPI) GetBaseData() ([]byte, error) {
	return b.GetData(b.url)
}

type IbaseAPI interface {
	GetPobs() ([]*configs_export.PoB, error)
}

func NewBaseApi() IbaseAPI {
	return darkhttp.NewClient(settings.Env.DarkstatApiUrl)
}
