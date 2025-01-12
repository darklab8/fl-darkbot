package configurator

import (
	"testing"

	"github.com/darklab8/fl-darkbot/app/settings/logus"
	"github.com/darklab8/fl-darkbot/app/settings/types"
	"github.com/darklab8/go-typelog/typelog"

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
		logus.Log.Debug("err=", typelog.OptError(err))

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

func TestAddForumIgnore(t *testing.T) {
	FixtureMigrator(func(dbname types.Dbpath) {
		configur := FixtureConfigurator(dbname)
		cg := NewConfiguratorForumIgnore(configur)

		ConfiguratorChannel{Configurator: configur}.Add("c1")
		cg.TagsAdd("c1", []types.Tag{"t1", "t2"}...)

		tags, _ := cg.TagsList("c1")
		assert.Len(t, tags, 2)

		watch_cfg := NewConfiguratorForumWatch(configur)
		watch_tags, _ := watch_cfg.TagsList("c1")
		assert.Len(t, watch_tags, 0)
	})

}
