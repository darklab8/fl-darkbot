package consoler

import (
	"darkbot/configurator"
	"darkbot/consoler/helper"
	"darkbot/settings"
	"darkbot/settings/types"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGettingOutput(t *testing.T) {
	configurator.FixtureMigrator(func(dbpath types.Dbpath) {
		channelID, _ := configurator.FixtureChannel(dbpath)
		assert.Contains(t, Consoler{}.New(settings.Config.ConsolerPrefix+" ping").Execute(helper.ChannelInfo{ChannelID: channelID, Dbpath: dbpath}).String(), "Pong!")
	})
}

func TestGrabStdout(t *testing.T) {
	configurator.FixtureMigrator(func(dbpath types.Dbpath) {
		channelID, _ := configurator.FixtureChannel(dbpath)
		c := Consoler{}.New(settings.Config.ConsolerPrefix + " ping --help")
		result := c.Execute(helper.ChannelInfo{ChannelID: channelID, Dbpath: dbpath}).String()

		assert.Contains(t, result, "\nFlags:\n  -h, --help   ")
	})
}

func TestAddBaseTag(t *testing.T) {
	configurator.FixtureMigrator(func(dbpath types.Dbpath) {
		channelID, _ := configurator.FixtureChannel(dbpath)
		assert.Contains(t, Consoler{}.New(settings.Config.ConsolerPrefix+` base add "bla bla" sdf`).Execute(helper.ChannelInfo{ChannelID: channelID, Dbpath: dbpath}).String(), "OK tags are added")
	})
}

func TestSystemCommands(t *testing.T) {
	configurator.FixtureMigrator(func(dbpath types.Dbpath) {
		channelID, _ := configurator.FixtureChannel(dbpath)
		cons := Consoler{}.New(settings.Config.ConsolerPrefix + ` player --help`)
		result := cons.Execute(helper.ChannelInfo{ChannelID: channelID, Dbpath: dbpath}).String()
		_ = result
		assert.Contains(t, result, "System commands")
	})
}
