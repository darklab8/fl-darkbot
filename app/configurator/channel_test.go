package configurator

import (
	"darkbot/app/settings/logus"
	"darkbot/app/settings/types"
	"testing"

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
		logus.Debug("invoked List", logus.Items(channels, "channels"))
		assert.Len(t, channels, 3)

		cg.Remove("3")

		channels, _ = cg.List()
		logus.Debug("", logus.Items(channels, "channels"))
		assert.Len(t, channels, 2)

		cg.Add("3")

		channels, _ = cg.List()
		logus.Debug("", logus.Items(channels, "channels"))
		assert.Len(t, channels, 3)
	})
}
