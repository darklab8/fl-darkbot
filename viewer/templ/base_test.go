package templ

import (
	"darkbot/configurator"
	"darkbot/dtypes"
	"darkbot/scrappy"
	"darkbot/scrappy/base"
	"darkbot/scrappy/player"
	"fmt"
	"testing"
)

func TestBaseViewer(t *testing.T) {
	configurator.FixtureMigrator(func(dbpath dtypes.Dbpath) {
		channelID, _ := configurator.FixtureChannel(dbpath)

		cg := configurator.ConfiguratorBase{Configurator: configurator.NewConfigurator(dbpath)}
		cg.TagsAdd(channelID, []string{"Station"}...)

		scrappy.Storage.BaseStorage.Api = base.APIBasespy{}
		scrappy.Storage.PlayerStorage.Api = player.APIPlayerSpy{}
		scrappy.Storage.Update()

		// If you want to add information explicitely
		// bases := base.BaseStorage{}
		// scrappy.Storage = &scrappy.ScrappyStorage{BaseStorage: &bases}
		// record := records.StampedObjects[base.Base]{}.New()
		// record.Add(base.Base{Name: "Station1", Affiliation: "Abc", Health: 100})
		// record.Add(base.Base{Name: "Station2", Affiliation: "Qwe", Health: 100})
		// bases.Add(record)

		base := NewTemplateBase(channelID, dbpath)
		base.Render()
		fmt.Println(base.main.Content)
	})
}

// func TestIntegrationTesting(t *testing.T) {
// 	os.Remove(settings.Dbpath)
// 	channelID := "838802002582175756"

// 	cg := configurator.ConfiguratorBase{Configurator: configurator.NewConfigurator()}
// 	cg.TagsAdd(channelID, []string{"Station"}...)

// 	scrappy.Storage.Update()

// 	base := BaseView{}
// 	base.ViewConfig = apis.NewAPI(channelID)
// 	base.Render()
// 	fmt.Println(base.Content)
// }
