package consoler

import (
	"darkbot/configurator"
	"darkbot/consoler/helper"
	"darkbot/dtypes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGettingOutput(t *testing.T) {
	configurator.FixtureMigrator(func(dbpath dtypes.Dbpath) {
		channelID, _ := configurator.FixtureChannel(dbpath)
		assert.Contains(t, Consoler{}.New(", ping").Execute(helper.ChannelInfo{ChannelID: channelID, Dbpath: dbpath}).String(), "Pong!")
	})
}

func TestGrabStdout(t *testing.T) {
	configurator.FixtureMigrator(func(dbpath dtypes.Dbpath) {
		channelID, _ := configurator.FixtureChannel(dbpath)
		c := Consoler{}.New(", ping --help")
		result := c.Execute(helper.ChannelInfo{ChannelID: channelID, Dbpath: dbpath}).String()

		assert.Contains(t, result, "\nFlags:\n  -h, --help   ")
	})
}

func TestAddBaseTag(t *testing.T) {
	configurator.FixtureMigrator(func(dbpath dtypes.Dbpath) {
		channelID, _ := configurator.FixtureChannel(dbpath)
		assert.Contains(t, Consoler{}.New(`, base add "bla bla" sdf`).Execute(helper.ChannelInfo{ChannelID: channelID, Dbpath: dbpath}).String(), "OK tags are added")
	})
}

func TestSystemCommands(t *testing.T) {
	configurator.FixtureMigrator(func(dbpath dtypes.Dbpath) {
		channelID, _ := configurator.FixtureChannel(dbpath)
		cons := Consoler{}.New(`, --help`)
		result := cons.Execute(helper.ChannelInfo{ChannelID: channelID, Dbpath: dbpath}).String()
		_ = result
		// TODO activate for new test
		// assert.Contains(t, result, "System commands")
	})
}
