package discorder

import (
	"testing"

	"github.com/darklab8/fl-darkbot/app/scrappy/player"
	"github.com/darklab8/fl-darkbot/app/settings/logus"
)

func TestLogging(t *testing.T) {

	long_test := "long_test_111111111111111111"
	short_test := "short"

	logus.Log.Warn("send long msg", logus.MsgContent(long_test))
	logus.Log.Warn("send shot msg", logus.MsgContent(short_test))
}

func TestReceivePlayersNumbersUpdate(t *testing.T) {
	// unit test by static typing check it is present :)

	if false {
		return
	}

	player_api := player.FixturePlayerAPIMock()
	player_storage := player.NewPlayerStorage(player_api)
	dg := NewClient()
	player_storage.RegisterObserve(dg) // updates number of players in Bot description

	player_storage.Update() // has at the end call of observers
	player_storage.UpdateObservers()

}
