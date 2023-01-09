package configurator

import (
	"darkbot/configurator/models"
	"darkbot/settings"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateTags(t *testing.T) {
	os.Remove(settings.Dbpath)

	cg := NewConfigurator()
	cg.ActionBaseTagsAdd("123", []string{"4"}...)
	cg.ActionBaseTagsAdd("123", []string{"5", "6"}...)

	baseTags := []models.TagBase{}
	cg.GetClient().Find(&baseTags)
	assert.Len(t, baseTags, 3)
}
