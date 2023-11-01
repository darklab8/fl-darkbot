package configurator

import (
	"darkbot/app/settings/types"
	"fmt"
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
		fmt.Println(channels)
		assert.Len(t, channels, 3)

		cg.Remove("3")

		channels, _ = cg.List()
		fmt.Println(channels)
		assert.Len(t, channels, 2)

		cg.Add("3")

		channels, _ = cg.List()
		fmt.Println(channels)
		assert.Len(t, channels, 3)
	})
}
