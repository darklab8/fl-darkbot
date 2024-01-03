package eventview

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

func TestPlayerEvent(t *testing.T) {
	configurator.FixtureMigrator(func(dbpath types.Dbpath) {
		channelID, _ := configurator.FixtureChannel(dbpath)

		configurator.NewConfiguratorPlayerEvent(configurator.NewConfigurator(dbpath)).TagsAdd(channelID, []types.Tag{"player1", "player2"}...)

		players := player.NewPlayerStorage(player.FixturePlayerAPIMock())
		storage := scrappy.FixtureMockedStorage(scrappy.WithPlayerStorage(players))
		api := apis.NewAPI(dbpath, storage)
		record := records.NewStampedObjects[player.Player]()
		record.Add(player.Player{Name: "player1", System: "system1", Region: "region1"})
		record.Add(player.Player{Name: "player2", System: "system2", Region: "region2"})
		record.Add(player.Player{Name: "player3", System: "system3", Region: "region3"})
		record.Add(player.Player{Name: "player4", System: "system4", Region: "region4"})
		players.Add(record)

		playerView := NewEventRenderer(api, channelID)
		playerView.RenderView()
		//logus.Debug(len(playerView.main.ViewRecords))
		logus.Log.Debug("test TestPlayerEvent is finished")

		assert.Equal(t, 2, playerView.main.RecordCount())
		assert.Equal(t, 1, len(playerView.main.GetMsgs()))
	})
}
