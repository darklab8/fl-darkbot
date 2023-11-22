package apis

import (
	"darkbot/app/configurator"
	"darkbot/app/discorder"
	"darkbot/app/scrappy"
	"darkbot/app/settings/types"
)

type API struct {
	Discorder *discorder.Discorder
	ChannelID types.DiscordChannelID
	Scrappy   *scrappy.ScrappyStorage
	*configurator.Configurators
}

type apiParam func(api *API)

func WithStorage(storage *scrappy.ScrappyStorage) apiParam {
	return func(api *API) {
		api.Scrappy = storage
	}
}

// Reusing API is important. db connections are memory leaking or some other stuff
func (api *API) SetChannelID(ChannelID types.DiscordChannelID) *API {
	api.ChannelID = ChannelID
	return api
}

func NewAPI(ChannelID types.DiscordChannelID, dbpath types.Dbpath, scrappy_storage *scrappy.ScrappyStorage, opts ...apiParam) *API {
	api := &API{
		ChannelID:     ChannelID,
		Discorder:     discorder.NewClient(),
		Scrappy:       scrappy_storage,
		Configurators: configurator.NewConfigugurators(dbpath),
	}

	for _, opt := range opts {
		opt(api)
	}

	return api
}
