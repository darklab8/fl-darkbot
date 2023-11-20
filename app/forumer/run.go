package forumer

import (
	"darkbot/app/configurator"
	"darkbot/app/discorder"
	"darkbot/app/settings/types"
)

type Forumer struct {
	Discorder discorder.Discorder
	*configurator.Configurators
}

func NewForumer(dbpath types.Dbpath) *Forumer {
	forum := &Forumer{
		Discorder:     discorder.NewClient(),
		Configurators: configurator.NewConfigugurators(dbpath),
	}
	return forum
}

func (f *Forumer) update() {
}

func (f *Forumer) Run() {

	for {

	}
}
