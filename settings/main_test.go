package settings

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain(t *testing.T) {
	load()
	isContains1 := strings.Contains(Config.Scrappy.Forum.Username, "darkwind")
	isContains2 := strings.Contains(Config.Scrappy.Forum.Username, "example")
	assert.True(t, isContains1 || isContains2)
}