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

	cg := ConfiguratorBase{Configurator: NewConfigurator()}
	cg.TagsAdd(channelID, []string{"4"}...)
	cg.TagsAdd(channelID, []string{"5", "6"}...)

	baseTags := cg.TagsList(channelID)
	assert.Len(t, baseTags, 3)

	cg.TagsClear(channelID)
	baseTags = cg.TagsList(channelID)
	assert.Len(t, baseTags, 0)
}

func TestCanWriteRepeatedTagsPerChannels(t *testing.T) {
	os.Remove(settings.Dbpath)
	cg := ConfiguratorBase{Configurator: NewConfigurator()}
	cg.TagsAdd("c1", []string{"t1"}...)
	cg.TagsAdd("c2", []string{"t1"}...)

	assert.Len(t, cg.TagsList("c1"), 1)
	assert.Len(t, cg.TagsList("c2"), 1)

	cg.TagsAdd("c2", []string{"t1"}...)

	// make a test to check? :thinking:
	assert.Len(t, cg.TagsList("c2"), 2)
}
