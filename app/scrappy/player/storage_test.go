package player

import (
	"darkbot/app/settings/darkbot_logus"
	"testing"

	"github.com/darklab8/darklab_goutils/goutils/utils_logus"
	"github.com/stretchr/testify/assert"
)

func FixturePlayerStorageMockified() *PlayerStorage {
	return NewPlayerStorage(FixturePlayerAPIMock())
}

func TestGetPlayers(t *testing.T) {
	storage := FixturePlayerStorageMockified()
	storage.Update()

	bases, err := storage.GetLatestRecord()
	darkbot_logus.Log.CheckFatal(err, "not found latest base record")

	assert.True(t, len(bases.List) > 0)
	darkbot_logus.Log.Debug("", utils_logus.Items(bases.List, "bases.List"))
}
