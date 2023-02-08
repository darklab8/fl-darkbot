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

func TestPlayerViewer(t *testing.T) {
	configurator.FixtureMigrator(func(dbpath dtypes.Dbpath) {
		channelID, _ := configurator.FixtureChannel(dbpath)

		playerCfg := configurator.ConfiguratorBase{Configurator: configurator.NewConfigurator(dbpath)}
		playerCfg.TagsAdd(channelID, []string{"Station"}...)

		scrappy.Storage.BaseStorage.Api = base.APIBasespy{}
		scrappy.Storage.PlayerStorage.Api = player.APIPlayerSpy{}
		scrappy.Storage.Update()
		// See TestBaseViewer for recipe to add manually players

		playerView := NewTemplatePlayers(channelID, dbpath)
		playerView.Render()
		fmt.Println(playerView.friends.Content)
		fmt.Println(playerView.enemies.Content)
		fmt.Println(playerView.all.Content)
	})
}
