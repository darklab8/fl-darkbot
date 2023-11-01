package templ

import (
	"darkbot/app/configurator"
	"darkbot/app/scrappy"
	"darkbot/app/scrappy/player"
	"darkbot/app/scrappy/shared/records"
	"darkbot/app/settings/logus"
	"darkbot/app/settings/types"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPlayerViewerMadeUpData(t *testing.T) {
	configurator.FixtureMigrator(func(dbpath types.Dbpath) {
		channelID, _ := configurator.FixtureChannel(dbpath)

		configurator.NewConfiguratorRegion(configurator.NewConfigurator(dbpath)).TagsAdd(channelID, []types.Tag{"region1", "region0"}...)
		configurator.NewConfiguratorSystem(configurator.NewConfigurator(dbpath)).TagsAdd(channelID, []types.Tag{"system1", "system2"}...)
		configurator.NewConfiguratorPlayerEnemy(configurator.NewConfigurator(dbpath)).TagsAdd(channelID, []types.Tag{"player2"}...)
		configurator.NewConfiguratorPlayerFriend(configurator.NewConfigurator(dbpath)).TagsAdd(channelID, []types.Tag{"player4"}...)

		players := player.NewPlayerStorage(player.FixturePlayerAPIMock())

		scrappy.Storage = scrappy.FixtureNewStorageWithPlayers(players)
		record := records.NewStampedObjects[player.Player]()
		record.Add(player.Player{Name: "player1", System: "system1", Region: "region1"})
		record.Add(player.Player{Name: "player2", System: "system2", Region: "region2"})
		record.Add(player.Player{Name: "player3", System: "system3", Region: "region3"})
		record.Add(player.Player{Name: "player4", System: "system4", Region: "region4"})
		players.Add(record)

		playerView := NewTemplatePlayers(channelID, dbpath)
		playerView.Render()
		logus.Debug(playerView.friends.MainTable.Content)
		logus.Debug(playerView.enemies.MainTable.Content)
		logus.Debug(playerView.neutral.MainTable.Content)
		logus.Debug("test TestPlayerViewer is finished")

		assert.NotEmpty(t, playerView.friends.MainTable.Content)
		assert.NotEmpty(t, playerView.enemies.MainTable.Content)
		assert.NotEmpty(t, playerView.neutral.MainTable.Content)

		assert.Empty(t, playerView.friends.AlertTmpl.Content)
		assert.Empty(t, playerView.enemies.AlertTmpl.Content)
		assert.Empty(t, playerView.neutral.AlertTmpl.Content)

		enemyAlerts := configurator.NewCfgAlertEnemyPlayersGreaterThan(configurator.NewConfigurator(dbpath))
		_, err := enemyAlerts.Status(channelID)
		assert.ErrorContains(t, err, "not found")

		enemyAlerts.Set(channelID, 1)
		integer, err := enemyAlerts.Status(channelID)
		assert.Equal(t, 1, integer)
		assert.Nil(t, err)

		playerView = NewTemplatePlayers(channelID, dbpath)
		playerView.Render()

		assert.NotEmpty(t, playerView.enemies.AlertTmpl.Content)
		assert.Empty(t, playerView.friends.AlertTmpl.Content)
		assert.Empty(t, playerView.neutral.AlertTmpl.Content)

	})
}

func TestPlayerViewerRealData(t *testing.T) {
	configurator.FixtureMigrator(func(dbpath types.Dbpath) {
		channelID, _ := configurator.FixtureChannel(dbpath)

		scrappy.Storage = scrappy.FixtureMockedStorage()
		scrappy.Storage.Update()

		configurator.NewConfiguratorPlayerFriend(configurator.NewConfigurator(dbpath)).TagsAdd(channelID, []types.Tag{"RM"}...)

		playerView := NewTemplatePlayers(channelID, dbpath)
		playerView.Render()

		logus.Debug(playerView.friends.MainTable.Content)
		logus.Debug(playerView.enemies.MainTable.Content)
		logus.Debug(playerView.neutral.MainTable.Content)

		logus.Debug("test TestPlayerViewer is finished")
	})
}
