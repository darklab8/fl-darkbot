package player

import (
	"darkbot/scrappy/shared/api"
	"darkbot/settings"
)

type PlayerAPI struct {
	api.APIrequest
}

func (a PlayerAPI) New() api.APIinterface {
	a.Init(settings.Config.ScrappyPlayerUrl)
	return a
}
