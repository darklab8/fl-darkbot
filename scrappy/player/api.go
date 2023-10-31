package player

import (
	"darkbot/scrappy/shared/api"
	"darkbot/settings"
)

type PlayerAPI struct {
	api.APIrequest
	url api.APIurl
}

func (b PlayerAPI) GetPlayerData() ([]byte, error) {
	return b.GetData(b.url)
}

type IPlayerAPI interface {
	GetPlayerData() ([]byte, error)
}

func NewPlayerAPI() PlayerAPI {
	a := PlayerAPI{}
	a.url = api.APIurl(settings.Config.ScrappyPlayerUrl)
	return a
}
