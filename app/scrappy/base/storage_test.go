package base

import (
	"darkbot/app/settings/darkbot_logus"
	"testing"

	"github.com/darklab8/darklab_goutils/goutils/utils_logus"
	"github.com/stretchr/testify/assert"
)

func FixtureBaseStorageMockified() *BaseStorage {
	return NewBaseStorage(FixtureBaseApiMock())
}

func TestGetBases(t *testing.T) {
	storage := FixtureBaseStorageMockified()
	storage.Update()

	bases, err := storage.GetLatestRecord()
	darkbot_logus.Log.CheckFatal(err, "not found latest base record")

	assert.True(t, len(bases.List) > 0)
	darkbot_logus.Log.Debug("", utils_logus.Items(bases.List, "bases.List"))
}

func TestAddManyRecords(t *testing.T) {

	storage := FixtureBaseStorageMockified()
	storage.Update()
	storage.Update()

	assert.Equal(t, 2, storage.Length())
}
