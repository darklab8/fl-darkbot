package templ

import (
	"darkbot/configurator"
	"darkbot/dtypes"
	"darkbot/scrappy"
	"darkbot/scrappy/base"
	"darkbot/scrappy/shared/records"
	"fmt"
	"testing"
)

func TestBaseViewerMocked(t *testing.T) {
	configurator.FixtureMigrator(func(dbpath dtypes.Dbpath) {
		channelID, _ := configurator.FixtureChannel(dbpath)

		cg := configurator.ConfiguratorBase{Configurator: configurator.NewConfigurator(dbpath)}
		cg.TagsAdd(channelID, []string{"Station"}...)

		bases := base.BaseStorage{}
		scrappy.Storage = &scrappy.ScrappyStorage{BaseStorage: &bases}
		record := records.StampedObjects[base.Base]{}.New()
		record.Add(base.Base{Name: "Station1", Affiliation: "Abc", Health: 100})
		record.Add(base.Base{Name: "Station2", Affiliation: "Qwe", Health: 100})
		bases.Add(record)

		base := NewTemplateBase(channelID, dbpath)
		base.Render()
		fmt.Println(base.main.Content)
	})
}

// TODO fix those tests... for some reason memory ref error :smile:
// func TestBaseViewerRealData(t *testing.T) {
// 	configurator.FixtureMigrator(func(dbpath dtypes.Dbpath) {
// 		channelID, _ := configurator.FixtureChannel(dbpath)

// 		cg := configurator.ConfiguratorBase{Configurator: configurator.NewConfigurator(dbpath)}
// 		cg.TagsAdd(channelID, []string{"Station"}...)

// 		scrappy.Storage.BaseStorage.Api = base.APIBasespy{}
// 		scrappy.Storage.PlayerStorage.Api = player.APIPlayerSpy{}
// 		scrappy.Storage.Update()

// 		base := NewTemplateBase(channelID, dbpath)
// 		base.Render()
// 		fmt.Println(base.main.Content)
// 	})
// }
