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

	cg := Base{Configurator: NewConfigurator()}
	cg.TagsAdd(channelID, []string{"4"}...)
	cg.TagsAdd(channelID, []string{"5", "6"}...)

	baseTags := cg.TagsList(channelID)
	assert.Len(t, baseTags, 3)

	cg.TagsClear(channelID)
	baseTags = cg.TagsList(channelID)
	assert.Len(t, baseTags, 0)
}
