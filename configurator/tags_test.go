package configurator

import (
	"darkbot/settings"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTags(t *testing.T) {
	os.Remove(settings.Dbpath)
	channelID := "123"
	ConfiguratorChannel{Configurator: NewConfigurator().Migrate()}.Add(channelID)
	ConfiguratorChannel{Configurator: NewConfigurator()}.Add(channelID)

	cg := ConfiguratorBase{Configurator: NewConfigurator().Migrate()}
	cg.TagsAdd(channelID, []string{"4"}...)
	cg.TagsAdd(channelID, []string{"5", "6"}...)

	baseTags, _ := cg.TagsList(channelID)
	assert.Len(t, baseTags, 3)

	cg.TagsClear(channelID)
	baseTags, _ = cg.TagsList(channelID)
	assert.Len(t, baseTags, 0)
}

func TestCanWriteRepeatedTagsPerChannels(t *testing.T) {
	os.Remove(settings.Dbpath)
	cg := ConfiguratorBase{Configurator: NewConfigurator().Migrate()}

	ConfiguratorChannel{Configurator: NewConfigurator()}.Add("c1")
	ConfiguratorChannel{Configurator: NewConfigurator()}.Add("c2")
	cg.TagsAdd("c1", []string{"t1"}...)
	cg.TagsAdd("c2", []string{"t1"}...)

	tags, _ := cg.TagsList("c1")
	assert.Len(t, tags, 1)
	tags, _ = cg.TagsList("c2")
	assert.Len(t, tags, 1)

	cg.TagsAdd("c2", []string{"t1"}...)

	// make a test to check? :thinking:
	tags, _ = cg.TagsList("c2")
	assert.Len(t, tags, 2)
}
