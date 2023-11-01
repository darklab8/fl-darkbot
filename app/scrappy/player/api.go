package player

import (
	"darkbot/app/scrappy/shared/api"
	"darkbot/app/settings"
	"darkbot/app/settings/types"
)

type PlayerAPI struct {
	api.APIrequest
	url types.APIurl
}

func (b PlayerAPI) GetPlayerData() ([]byte, error) {
	return b.GetData(b.url)
}

type IPlayerAPI interface {
	GetPlayerData() ([]byte, error)
}

func NewPlayerAPI() PlayerAPI {
	a := PlayerAPI{}
	a.url = types.APIurl(settings.Config.ScrappyPlayerUrl)
	return a
}
