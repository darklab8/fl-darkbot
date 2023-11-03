package consoler

import (
	"darkbot/app/configurator"
	"darkbot/app/settings"
	"darkbot/app/settings/types"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGettingOutput(t *testing.T) {
	configurator.FixtureMigrator(func(dbpath types.Dbpath) {
		channelID, _ := configurator.FixtureChannel(dbpath)
		assert.Contains(t, NewConsoler(dbpath).Execute(settings.Config.ConsolerPrefix+" ping", channelID).String(), "Pong!")
	})
}

func TestGrabStdout(t *testing.T) {
	configurator.FixtureMigrator(func(dbpath types.Dbpath) {
		channelID, _ := configurator.FixtureChannel(dbpath)
		c := NewConsoler(dbpath)
		result := c.Execute(settings.Config.ConsolerPrefix+" ping --help", channelID).String()

		assert.Contains(t, result, "\nFlags:\n  -h, --help   ")
	})
}

func TestAddBaseTag(t *testing.T) {
	configurator.FixtureMigrator(func(dbpath types.Dbpath) {
		channelID, _ := configurator.FixtureChannel(dbpath)
		assert.Contains(t, NewConsoler(dbpath).Execute(settings.Config.ConsolerPrefix+` base add "bla bla" sdf`, channelID).String(), "OK tags are added")
	})
}

func TestSystemCommands(t *testing.T) {
	configurator.FixtureMigrator(func(dbpath types.Dbpath) {
		channelID, _ := configurator.FixtureChannel(dbpath)
		cons := NewConsoler(dbpath)
		result := cons.Execute(settings.Config.ConsolerPrefix+` player --help`, channelID).String()
		_ = result
		assert.Contains(t, result, "System commands")
	})
}
