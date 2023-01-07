package player

import (
	"darkbot/scrappy/apiRawData"
	"darkbot/settings"
)

type PlayerAPI struct {
	apiRawData.APIrequest
}

func (a PlayerAPI) New() apiRawData.APIinterface {
	a.Init(settings.Config.Scrappy.Player.URL)
	return a
}
