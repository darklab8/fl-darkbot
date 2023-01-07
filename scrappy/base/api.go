package base

import (
	"darkbot/scrappy/shared/api"
	"darkbot/settings"
)

type basesAPI struct {
	api.APIrequest
}

func (a basesAPI) New() api.APIinterface {
	a.Init(settings.Config.Scrappy.Base.URL)
	return a
}
