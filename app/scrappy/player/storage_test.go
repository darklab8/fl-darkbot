package player

import (
	"testing"

	"github.com/darklab8/fl-darkbot/app/settings/logus"
	"github.com/darklab8/go-typelog/typelog"

	"github.com/stretchr/testify/assert"
)

func FixturePlayerStorageMockified() *PlayerStorage {
	return NewPlayerStorage(FixturePlayerAPIMock())
}

func TestGetPlayers(t *testing.T) {
	storage := FixturePlayerStorageMockified()
	storage.Update()

	bases, err := storage.GetLatestRecord()
	logus.Log.CheckFatal(err, "not found latest base record")

	assert.True(t, len(bases.List) > 0)
	logus.Log.Debug("", typelog.Items("bases.List", bases.List))
}
