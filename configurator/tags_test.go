package configurator

import (
	"darkbot/dtypes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTags(t *testing.T) {
	FixtureMigrator(func(dbname dtypes.Dbpath) {
		channelID, _ := FixtureChannel(dbname)

		cg := ConfiguratorBase{Configurator: NewConfigurator(dbname).Migrate()}
		cg.TagsAdd(channelID, []string{"4"}...)
		cg.TagsAdd(channelID, []string{"5", "6"}...)

		baseTags, _ := cg.TagsList(channelID)
		assert.Len(t, baseTags, 3)

		cg.TagsClear(channelID)
		baseTags, _ = cg.TagsList(channelID)
		assert.Len(t, baseTags, 0)
	})
}

func TestCanWriteRepeatedTagsPerChannels(t *testing.T) {
	FixtureMigrator(func(dbname dtypes.Dbpath) {
		configur := FixtureConfigurator(dbname)
		cg := ConfiguratorBase{Configurator: configur}

		ConfiguratorChannel{Configurator: NewConfigurator(dbname)}.Add("c1")
		ConfiguratorChannel{Configurator: NewConfigurator(dbname)}.Add("c2")
		cg.TagsAdd("c1", []string{"t1"}...)
		cg.TagsAdd("c2", []string{"t1"}...)

		tags, _ := cg.TagsList("c1")
		assert.Len(t, tags, 1)
		tags, _ = cg.TagsList("c2")
		assert.Len(t, tags, 1)

		err := cg.TagsAdd("c2", []string{"t1"}...)
		assert.Contains(t, err.GetError().Error(), "database already has those items")
		fmt.Println("err=", err.GetError().Error())

		// make a test to check? :thinking:
		tags, _ = cg.TagsList("c2")
		assert.Len(t, tags, 1)
	})

}

func TestDoNotInputRepeatedTags(t *testing.T) {
	FixtureMigrator(func(dbname dtypes.Dbpath) {
		configur := FixtureConfigurator(dbname)
		cg := ConfiguratorBase{Configurator: configur}

		ConfiguratorChannel{Configurator: configur}.Add("c1")
		cg.TagsAdd("c1", []string{"t1", "t2"}...)
		cg.TagsAdd("c1", []string{"t1"}...)

		tags, _ := cg.TagsList("c1")
		assert.Len(t, tags, 2)
	})

}