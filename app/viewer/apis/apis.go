package apis

import (
	"darkbot/app/configurator"
	"darkbot/app/discorder"
	"darkbot/app/scrappy"
	"darkbot/app/settings/types"
)

type API struct {
	Discorder *discorder.Discorder
	Scrappy   *scrappy.ScrappyStorage
	*configurator.Configurators
}

type apiParam func(api *API)

func WithStorage(storage *scrappy.ScrappyStorage) apiParam {
	return func(api *API) {
		api.Scrappy = storage
	}
}

func NewAPI(dbpath types.Dbpath, scrappy_storage *scrappy.ScrappyStorage, opts ...apiParam) *API {
	api := &API{
		Discorder:     discorder.NewClient(),
		Scrappy:       scrappy_storage,
		Configurators: configurator.NewConfigugurators(dbpath),
	}

	for _, opt := range opts {
		opt(api)
	}

	return api
}
