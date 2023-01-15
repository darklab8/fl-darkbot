package configurator

import (
	"darkbot/settings"
	"fmt"
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

func TestChannels(t *testing.T) {
	os.Remove(settings.Dbpath)
	cg := ConfiguratorChannel{Configurator: NewConfigurator()}

	cg.Add("1", "2", "3")

	channels := cg.List()
	fmt.Println(channels)

	assert.Len(t, channels, 3)
}
