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

func TestPlayerViewer(t *testing.T) {
	configurator.FixtureMigrator(func(dbpath dtypes.Dbpath) {
		channelID, _ := configurator.FixtureChannel(dbpath)

		playerCfg := configurator.ConfiguratorBase{Configurator: configurator.NewConfigurator(dbpath)}
		playerCfg.TagsAdd(channelID, []string{"Station"}...)

		players := player.PlayerStorage{}
		scrappy.Storage = &scrappy.ScrappyStorage{PlayerStorage: &players}
		record := records.StampedObjects[player.Player]{}.New()
		record.Add(player.Player{Name: "abc_123", System: "New York"})
		record.Add(player.Player{Name: "abc_456", System: "Hamburg"})
		record.Add(player.Player{Name: "456456", System: "Hamburg"})
		players.Add(record)

		playerView := NewTemplatePlayers(channelID, dbpath)
		playerView.Render()
		fmt.Println(playerView.friends.Content)
		fmt.Println(playerView.enemies.Content)
		fmt.Println(playerView.all.Content)
	})
}
