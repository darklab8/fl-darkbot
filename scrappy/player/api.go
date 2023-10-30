package player

import (
	"darkbot/scrappy/shared/api"
	"darkbot/settings"
)

type PlayerAPI struct {
	api.APIrequest
}

func NewPlayerAPI() PlayerAPI {
	a := PlayerAPI{}
	a.Init(settings.Config.ScrappyPlayerUrl)
	return a
}
