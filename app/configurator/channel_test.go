package configurator

import (
	"darkbot/app/settings/darkbot_logus"
	"darkbot/app/settings/types"
	"testing"

	"github.com/darklab8/darklab_goutils/goutils/logus"
	"github.com/stretchr/testify/assert"
)

func TestChannels(t *testing.T) {
	FixtureMigrator(func(dbpath types.Dbpath) {
		channelID, cg := FixtureChannel(dbpath)
		cg.Remove(channelID)

		cg.Add("1")
		cg.Add("2")
		cg.Add("3")

		channels, _ := cg.List()
		darkbot_logus.Log.Debug("invoked List", logus.Items(channels, "channels"))
		assert.Len(t, channels, 3)

		cg.Remove("3")

		channels, _ = cg.List()
		darkbot_logus.Log.Debug("", logus.Items(channels, "channels"))
		assert.Len(t, channels, 2)

		cg.Add("3")

		channels, _ = cg.List()
		darkbot_logus.Log.Debug("", logus.Items(channels, "channels"))
		assert.Len(t, channels, 3)
	})
}
