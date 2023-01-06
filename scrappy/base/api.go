package base

import (
	"darkbot/scrappy/api"
	"darkbot/settings"
)

type BasesAPI struct {
	api.API
}

func (a BasesAPI) New() api.APIinterface {
	a.Init(settings.Config.Scrappy.Base.URL)
	return a
}
