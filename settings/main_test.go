package settings

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain(t *testing.T) {
	load()
	isContains1 := strings.Contains(Config.ScrappyForumUsername, "darkwind")
	isContains2 := strings.Contains(Config.ScrappyForumUsername, "example")
	assert.True(t, isContains1 || isContains2)
}
