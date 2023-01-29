package templ

import (
	"darkbot/configurator"
	"darkbot/dtypes"
	"testing"
)

func TestPlayerViewer(t *testing.T) {
	configurator.FixtureMigrator(func(dbpath dtypes.Dbpath) {
		channelID, _ := configurator.FixtureChannel(dbpath)
		_ = channelID
	})

	// cg := configurator.ConfiguratorBase{Configurator: configurator.NewConfigurator()}
	// cg.TagsAdd(channelID, []string{"Station"}...)

	// bases := base.BaseStorage{}
	// scrappy.Storage = &scrappy.ScrappyStorage{BaseStorage: &bases}
	// record := records.StampedObjects[base.Base]{}.New()
	// record.Add("Station1", base.Base{Name: "Station1", Affiliation: "Abc", Health: 100})
	// record.Add("Station2", base.Base{Name: "Station2", Affiliation: "Qwe", Health: 100})
	// bases.Add(record)

	// base := NewTemplateBase(channelID)
	// base.Render()
	// fmt.Println(base.main.Content)
}
