package consoler

import (
	"testing"

	"github.com/darklab8/fl-darkbot/app/configurator"
	"github.com/darklab8/fl-darkbot/app/settings"
	"github.com/darklab8/fl-darkbot/app/settings/types"

	"github.com/stretchr/testify/assert"
)

func TestGettingOutput(t *testing.T) {
	configurator.FixtureMigrator(func(dbpath types.Dbpath) {
		channelID, _ := configurator.FixtureChannel(dbpath)
		assert.Contains(t, NewConsoler(dbpath).Execute(settings.Env.ConsolerPrefix+" ping", channelID), "Pong!")
	})
}

func TestGrabStdout(t *testing.T) {
	configurator.FixtureMigrator(func(dbpath types.Dbpath) {
		channelID, _ := configurator.FixtureChannel(dbpath)
		c := NewConsoler(dbpath)
		result := c.Execute(settings.Env.ConsolerPrefix+" ping --help", channelID)

		assert.Contains(t, result, "\nFlags:\n  -h, --help   ")
	})
}

func TestAddBaseTag(t *testing.T) {
	configurator.FixtureMigrator(func(dbpath types.Dbpath) {
		channelID, _ := configurator.FixtureChannel(dbpath)
		assert.Contains(t, NewConsoler(dbpath).Execute(settings.Env.ConsolerPrefix+` base tags add "bla bla" sdf`, channelID), "OK tags are added")
	})
}

func TestSystemCommands(t *testing.T) {
	configurator.FixtureMigrator(func(dbpath types.Dbpath) {
		channelID, _ := configurator.FixtureChannel(dbpath)
		cons := NewConsoler(dbpath)
		result := cons.Execute(settings.Env.ConsolerPrefix+` player --help`, channelID)
		_ = result
		assert.Contains(t, result, "System commands")
	})
}

func TestAddForumIgnore(t *testing.T) {
	configurator.FixtureMigrator(func(dbpath types.Dbpath) {
		channelID, _ := configurator.FixtureChannel(dbpath)
		response := NewConsoler(dbpath).Execute(settings.Env.ConsolerPrefix+` forum thread ignore add Soundtrack`, channelID)
		response2 := NewConsoler(dbpath).Execute(settings.Env.ConsolerPrefix+` forum thread watch list`, channelID)

		_ = response
		_ = response2
		// assert.Contains(t, )
	})
}
