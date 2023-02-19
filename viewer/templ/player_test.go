package templ

import (
	"darkbot/configurator"
	"darkbot/dtypes"
	"darkbot/scrappy"
	"darkbot/scrappy/player"
	"darkbot/scrappy/shared/records"
	"fmt"
	"testing"
)

func TestPlayerViewerMadeUpData(t *testing.T) {
	configurator.FixtureMigrator(func(dbpath dtypes.Dbpath) {
		channelID, _ := configurator.FixtureChannel(dbpath)

		configurator.ConfiguratorRegion{Configurator: configurator.NewConfigurator(dbpath)}.TagsAdd(channelID, []string{"region1", "region0"}...)
		configurator.ConfiguratorSystem{Configurator: configurator.NewConfigurator(dbpath)}.TagsAdd(channelID, []string{"system1", "system2"}...)
		configurator.ConfiguratorPlayerEnemy{Configurator: configurator.NewConfigurator(dbpath)}.TagsAdd(channelID, []string{"player2"}...)
		configurator.ConfiguratorPlayerFriend{Configurator: configurator.NewConfigurator(dbpath)}.TagsAdd(channelID, []string{"player4"}...)

		players := player.PlayerStorage{}
		scrappy.Storage = &scrappy.ScrappyStorage{PlayerStorage: &players}
		record := records.StampedObjects[player.Player]{}.New()
		record.Add(player.Player{Name: "player1", System: "system1", Region: "region1"})
		record.Add(player.Player{Name: "player2", System: "system2", Region: "region2"})
		record.Add(player.Player{Name: "player3", System: "system3", Region: "region3"})
		record.Add(player.Player{Name: "player4", System: "system4", Region: "region4"})
		players.Add(record)

		playerView := NewTemplatePlayers(channelID, dbpath)
		playerView.Render()
		fmt.Println(playerView.friends.MainTable.Content)
		fmt.Println(playerView.enemies.MainTable.Content)
		fmt.Println(playerView.neutral.MainTable.Content)
		fmt.Println("test TestPlayerViewer is finished")
	})
}

// TODO fix those tests... for some reason memory ref error :smile:
// func TestPlayerViewerRealData(t *testing.T) {
// 	configurator.FixtureMigrator(func(dbpath dtypes.Dbpath) {
// 		channelID, _ := configurator.FixtureChannel(dbpath)

// 		scrappy.Storage.BaseStorage.Api = base.APIBasespy{}
// 		scrappy.Storage.PlayerStorage.Api = player.APIPlayerSpy{}
// 		scrappy.Storage.Update()

// 		configurator.ConfiguratorPlayerFriend{Configurator: configurator.NewConfigurator(dbpath)}.TagsAdd(channelID, []string{"RM"}...)

// 		playerView := NewTemplatePlayers(channelID, dbpath)
// 		playerView.Render()
// 		fmt.Println(playerView.friends.Content)
// 		fmt.Println(playerView.enemies.Content)
// 		fmt.Println(playerView.neutral.Content)
// 		fmt.Println("test TestPlayerViewer is finished")
// 	})
// }
