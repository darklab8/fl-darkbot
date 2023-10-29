package templ

import (
	"darkbot/configurator"
	"darkbot/scrappy"
	"darkbot/scrappy/player"
	"darkbot/scrappy/shared/records"
	"darkbot/settings/types"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPlayerViewerMadeUpData(t *testing.T) {
	configurator.FixtureMigrator(func(dbpath types.Dbpath) {
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

		assert.NotEmpty(t, playerView.friends.MainTable.Content)
		assert.NotEmpty(t, playerView.enemies.MainTable.Content)
		assert.NotEmpty(t, playerView.neutral.MainTable.Content)

		assert.Empty(t, playerView.friends.AlertTmpl.Content)
		assert.Empty(t, playerView.enemies.AlertTmpl.Content)
		assert.Empty(t, playerView.neutral.AlertTmpl.Content)

		enemyAlerts := configurator.CfgAlertEnemyPlayersGreaterThan{Configurator: configurator.NewConfigurator(dbpath)}
		integer, _ := enemyAlerts.Status(channelID)
		assert.Nil(t, integer)

		enemyAlerts.Set(channelID, 1)
		integer, _ = enemyAlerts.Status(channelID)
		assert.Equal(t, 1, *integer)

		playerView = NewTemplatePlayers(channelID, dbpath)
		playerView.Render()

		assert.NotEmpty(t, playerView.enemies.AlertTmpl.Content)
		assert.Empty(t, playerView.friends.AlertTmpl.Content)
		assert.Empty(t, playerView.neutral.AlertTmpl.Content)

	})
}

// TODO fix those tests... for some reason memory ref error :smile:
// func TestPlayerViewerRealData(t *testing.T) {
// 	configurator.FixtureMigrator(func(dbpath types.Dbpath) {
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
