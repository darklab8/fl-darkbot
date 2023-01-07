package base

import (
	"darkbot/scrappy/apiRawData"
	"darkbot/settings"
)

type basesAPI struct {
	apiRawData.APIrequest
}

func (a basesAPI) New() apiRawData.APIinterface {
	a.Init(settings.Config.Scrappy.Base.URL)
	return a
}
