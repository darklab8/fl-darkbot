package base

import (
	"darkbot/scrappy/shared/api"
	"darkbot/settings"
)

type basesAPI struct {
	api.APIrequest
}

func NewBaseApi() api.APIinterface {
	a := basesAPI{}
	a.Init(settings.Config.ScrappyBaseUrl)
	return a
}
