package player

import (
	"darkbot/scrappy/api"
	"darkbot/settings"
)

type PlayerAPI struct {
	api.API
}

func (a PlayerAPI) New() api.APIinterface {
	a.Init(settings.Config.Scrappy.Player.URL)
	return a
}
