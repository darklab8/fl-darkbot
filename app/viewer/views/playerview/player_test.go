package playerview

import (
	"darkbot/app/configurator"
	"darkbot/app/scrappy"
	"darkbot/app/scrappy/player"
	"darkbot/app/scrappy/shared/records"
	"darkbot/app/settings/logus"
	"darkbot/app/settings/types"
	"darkbot/app/viewer/apis"
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
		storage := scrappy.FixtureMockedStorage(scrappy.WithPlayerStorage(players))
		api := apis.NewAPI(dbpath, storage)
		record := records.NewStampedObjects[player.Player]()
		record.Add(player.Player{Name: "player1", System: "system1", Region: "region1"})
		record.Add(player.Player{Name: "player2", System: "system2", Region: "region2"})
		record.Add(player.Player{Name: "player3", System: "system3", Region: "region3"})
		record.Add(player.Player{Name: "player4", System: "system4", Region: "region4"})
		players.Add(record)

		playerView := NewTemplatePlayers(api, channelID)
		playerView.RenderView()
		logus.Debug(playerView.friends.mainTable.GetMsgs()[0].Render())
		logus.Debug(playerView.enemies.mainTable.GetMsgs()[0].Render())
		logus.Debug(playerView.neutral.mainTable.GetMsgs()[0].Render())
		logus.Debug("test TestPlayerViewer is finished")

		assert.True(t, playerView.friends.mainTable.HasRecords())
		assert.True(t, playerView.enemies.mainTable.HasRecords())
		assert.True(t, playerView.neutral.mainTable.HasRecords())

		assert.False(t, playerView.friends.alertTmpl.HasRecords())
		assert.False(t, playerView.enemies.alertTmpl.HasRecords())
		assert.False(t, playerView.neutral.alertTmpl.HasRecords())

		enemyAlerts := configurator.NewCfgAlertEnemyPlayersGreaterThan(configurator.NewConfigurator(dbpath))
		_, err := enemyAlerts.Status(channelID)
		assert.Error(t, err, configurator.ErrorZeroAffectedRowsMsg)

		enemyAlerts.Set(channelID, 1)
		integer, err := enemyAlerts.Status(channelID)
		assert.Equal(t, 1, integer)
		assert.Nil(t, err)

		playerView = NewTemplatePlayers(api, channelID)
		playerView.RenderView()

		assert.True(t, playerView.enemies.alertTmpl.HasRecords())
		assert.False(t, playerView.friends.alertTmpl.HasRecords())
		assert.False(t, playerView.neutral.alertTmpl.HasRecords())

	})
}

func TestPlayerViewerRealData(t *testing.T) {
	configurator.FixtureMigrator(func(dbpath types.Dbpath) {
		channelID, _ := configurator.FixtureChannel(dbpath)

		storage := scrappy.FixtureMockedStorage()
		storage.Update()

		configur := configurator.NewConfigurator(dbpath)
		configurator.NewConfiguratorPlayerFriend(configur).TagsAdd(channelID, []types.Tag{"RM"}...)

		playerView := NewTemplatePlayers(apis.NewAPI(dbpath, storage), channelID)
		playerView.RenderView()

		assert.True(t, playerView.friends.mainTable.HasRecords())
		assert.False(t, playerView.enemies.mainTable.HasRecords())
		assert.False(t, playerView.neutral.mainTable.HasRecords())

		logus.Debug(playerView.friends.mainTable.GetMsgs()[0].Render())

		logus.Debug("test TestPlayerViewer is finished")
	})
}
