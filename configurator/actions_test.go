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

	cg := NewConfigurator()
	cg.ActionBaseTagsAdd(channelID, []string{"4"}...)
	cg.ActionBaseTagsAdd(channelID, []string{"5", "6"}...)

	baseTags := cg.ActionBaseTagsList(channelID)
	assert.Len(t, baseTags, 3)

	cg.ActionBaseTagsClear(channelID)
	baseTags = cg.ActionBaseTagsList(channelID)
	assert.Len(t, baseTags, 0)
}
