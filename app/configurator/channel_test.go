package configurator

import (
	"testing"

	"github.com/darklab/fl-darkbot/app/settings/logus"
	"github.com/darklab/fl-darkbot/app/settings/types"

	"github.com/darklab8/darklab_goutils/goutils/logus_core"
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
		logus.Log.Debug("invoked List", logus_core.Items(channels, "channels"))
		assert.Len(t, channels, 3)

		cg.Remove("3")

		channels, _ = cg.List()
		logus.Log.Debug("", logus_core.Items(channels, "channels"))
		assert.Len(t, channels, 2)

		cg.Add("3")

		channels, _ = cg.List()
		logus.Log.Debug("", logus_core.Items(channels, "channels"))
		assert.Len(t, channels, 3)
	})
}
