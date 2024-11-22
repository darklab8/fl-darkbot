package playerview

import (
	"fmt"
	"testing"

	"github.com/darklab8/fl-darkbot/app/configurator"
	"github.com/darklab8/fl-darkbot/app/scrappy"
	"github.com/darklab8/fl-darkbot/app/settings/logus"
	"github.com/darklab8/fl-darkbot/app/settings/types"
	"github.com/darklab8/fl-darkbot/app/viewer/apis"
	"github.com/stretchr/testify/assert"
)

func TestPlayerViewerRealDataDebug(t *testing.T) {
	return

	dbpath := types.Dbpath("/home/naa/repos/pet_projects/fl-darkbot/data/dev.sqlite3")
	// configurator.FixtureMigrator(func(dbpath types.Dbpath) {

	channelID := types.DiscordChannelID("1079189823098724433")
	// cfg_channel := NewConfiguratorChannel(configurator_)
	// cfg_channel.Add(channelID)
	// channelID, _ := configurator.FixtureChannel(dbpath)

	// configurator.NewConfigugurators(dbpath)

	// storage := scrappy.FixtureMockedStorage()
	storage := scrappy.NewScrappyWithApis()
	storage.Update()

	configur := configurator.NewConfigurator(dbpath)
	// R\Nickname\-Tester
	// player.Player
	// Visible as
	// Time = "4m"
	// Name = "R\\Nickname\\-Tester"
	// System = "Pennsylvania"
	// Region = "Liberty Space"

	var tags []types.Tag = []types.Tag{`\-Tester`, `R\`, `R\Nickname`, `\-`}
	for _, tag := range tags {
		fmt.Println("testing for tag=", tag)
		players := configurator.NewConfiguratorPlayerFriend(configur)
		// players.TagsAdd(channelID, tag)

		tags, err := players.TagsList(channelID)
		// seen as Nickname\\-Tester, Nickname\-Tester
		// [0] = "Nickname\\\\\\\\-Tester"
		// [1] = "Nickname\\\\-Tester"

		// ; player friend add R\N, gets as R\\N
		// fmt.Println(tags, err)
		fmt.Printf("%#v, %#v", tags, err)

		playerView := NewTemplatePlayers(apis.NewAPI(dbpath, storage), channelID)
		playerView.RenderView()

		assert.True(t, playerView.friends.mainTable.HasRecords())
		assert.False(t, playerView.enemies.mainTable.HasRecords())
		assert.False(t, playerView.neutral.mainTable.HasRecords())

		fmt.Println(playerView.friends.mainTable.GetMsgs()[0].Render())

		logus.Log.Debug("test TestPlayerViewer is finished")
	}

	// })

	/* Confusion
	Main problem that player get extra confused with extra \\ handling of discord.
	Converting things to ``` ``` escape will be a good solution.
	*/

}
