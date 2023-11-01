package consoler

import (
	"darkbot/configurator"
	"darkbot/consoler/printer"
	"darkbot/settings"
	"darkbot/settings/types"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGettingOutput(t *testing.T) {
	configurator.FixtureMigrator(func(dbpath types.Dbpath) {
		channelID, _ := configurator.FixtureChannel(dbpath)
		assert.Contains(t, NewConsoler(settings.Config.ConsolerPrefix+" ping").Execute(printer.ChannelInfo{ChannelID: channelID, Dbpath: dbpath}).String(), "Pong!")
	})
}

func TestGrabStdout(t *testing.T) {
	configurator.FixtureMigrator(func(dbpath types.Dbpath) {
		channelID, _ := configurator.FixtureChannel(dbpath)
		c := NewConsoler(settings.Config.ConsolerPrefix + " ping --help")
		result := c.Execute(printer.ChannelInfo{ChannelID: channelID, Dbpath: dbpath}).String()

		assert.Contains(t, result, "\nFlags:\n  -h, --help   ")
	})
}

func TestAddBaseTag(t *testing.T) {
	configurator.FixtureMigrator(func(dbpath types.Dbpath) {
		channelID, _ := configurator.FixtureChannel(dbpath)
		assert.Contains(t, NewConsoler(settings.Config.ConsolerPrefix+` base add "bla bla" sdf`).Execute(printer.ChannelInfo{ChannelID: channelID, Dbpath: dbpath}).String(), "OK tags are added")
	})
}

func TestSystemCommands(t *testing.T) {
	configurator.FixtureMigrator(func(dbpath types.Dbpath) {
		channelID, _ := configurator.FixtureChannel(dbpath)
		cons := NewConsoler(settings.Config.ConsolerPrefix + ` player --help`)
		result := cons.Execute(printer.ChannelInfo{ChannelID: channelID, Dbpath: dbpath}).String()
		_ = result
		assert.Contains(t, result, "System commands")
	})
}
