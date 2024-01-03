package configurator

import (
	"darkbot/app/settings/logus"
	"darkbot/app/settings/types"
	"testing"

	"github.com/darklab8/darklab_goutils/goutils/logus_core"
	"github.com/stretchr/testify/assert"
)

func TestTags(t *testing.T) {
	FixtureMigrator(func(dbname types.Dbpath) {
		channelID, _ := FixtureChannel(dbname)

		cg := NewConfiguratorBase(NewConfigurator(dbname).AutoMigrateSchema())
		cg.TagsAdd(channelID, []types.Tag{"4"}...)
		cg.TagsAdd(channelID, []types.Tag{"5", "6"}...)

		baseTags, _ := cg.TagsList(channelID)
		assert.Len(t, baseTags, 3)

		cg.TagsClear(channelID)
		baseTags, _ = cg.TagsList(channelID)
		assert.Len(t, baseTags, 0)
	})
}

func TestCanWriteRepeatedTagsPerChannels(t *testing.T) {
	FixtureMigrator(func(dbname types.Dbpath) {
		configur := FixtureConfigurator(dbname)
		cg := NewConfiguratorBase(configur)

		ConfiguratorChannel{Configurator: NewConfigurator(dbname)}.Add("c1")
		ConfiguratorChannel{Configurator: NewConfigurator(dbname)}.Add("c2")
		cg.TagsAdd("c1", []types.Tag{"t1"}...)
		cg.TagsAdd("c2", []types.Tag{"t1"}...)

		tags, _ := cg.TagsList("c1")
		assert.Len(t, tags, 1)
		tags, _ = cg.TagsList("c2")
		assert.Len(t, tags, 1)

		err := cg.TagsAdd("c2", []types.Tag{"t1"}...)

		assert.Error(t, err, "expected error to get in test")
		assert.Contains(t, err.Error(), "database already has those items")
		logus.Log.Debug("err=", logus_core.OptError(err))

		// make a test to check? :thinking:
		tags, _ = cg.TagsList("c2")
		assert.Len(t, tags, 1)
	})

}

func TestDoNotInputRepeatedTags(t *testing.T) {
	FixtureMigrator(func(dbname types.Dbpath) {
		configur := FixtureConfigurator(dbname)
		cg := NewConfiguratorBase(configur)

		ConfiguratorChannel{Configurator: configur}.Add("c1")
		cg.TagsAdd("c1", []types.Tag{"t1", "t2"}...)
		cg.TagsAdd("c1", []types.Tag{"t1"}...)

		tags, _ := cg.TagsList("c1")
		assert.Len(t, tags, 2)
	})

}
