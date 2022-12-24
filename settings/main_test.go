package settings

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain(t *testing.T) {
	load()
	isContains := strings.Contains(Config.Scrappy.Forum.Username, "darkwind")
	assert.True(t, isContains)
}
