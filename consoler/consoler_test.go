package consoler

import (
	"darkbot/configurator"
	"darkbot/settings"
	"darkbot/settings/types"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGettingOutput(t *testing.T) {
	configurator.FixtureMigrator(func(dbpath types.Dbpath) {
		channelID, _ := configurator.FixtureChannel(dbpath)
		assert.Contains(t, NewConsoler(settings.Config.ConsolerPrefix+" ping", channelID, dbpath).Execute().String(), "Pong!")
	})
}

func TestGrabStdout(t *testing.T) {
	configurator.FixtureMigrator(func(dbpath types.Dbpath) {
		channelID, _ := configurator.FixtureChannel(dbpath)
		c := NewConsoler(settings.Config.ConsolerPrefix+" ping --help", channelID, dbpath)
		result := c.Execute().String()

		assert.Contains(t, result, "\nFlags:\n  -h, --help   ")
	})
}

func TestAddBaseTag(t *testing.T) {
	configurator.FixtureMigrator(func(dbpath types.Dbpath) {
		channelID, _ := configurator.FixtureChannel(dbpath)
		assert.Contains(t, NewConsoler(settings.Config.ConsolerPrefix+` base add "bla bla" sdf`, channelID, dbpath).Execute().String(), "OK tags are added")
	})
}

func TestSystemCommands(t *testing.T) {
	configurator.FixtureMigrator(func(dbpath types.Dbpath) {
		channelID, _ := configurator.FixtureChannel(dbpath)
		cons := NewConsoler(settings.Config.ConsolerPrefix+` player --help`, channelID, dbpath)
		result := cons.Execute().String()
		_ = result
		assert.Contains(t, result, "System commands")
	})
}
