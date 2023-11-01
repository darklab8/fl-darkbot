package base

import (
	"darkbot/app/settings/logus"
	"testing"

	"github.com/stretchr/testify/assert"
)

func FixtureBaseStorageMockified() *BaseStorage {
	return NewBaseStorage(FixtureBaseApiMock())
}

func TestGetBases(t *testing.T) {
	storage := FixtureBaseStorageMockified()
	storage.Update()

	bases, err := storage.GetLatestRecord()
	logus.CheckFatal(err, "not found latest base record")

	assert.True(t, len(bases.List) > 0)
	logus.Debug("", logus.Items(bases.List, "bases.List"))
}

func TestAddManyRecords(t *testing.T) {

	storage := FixtureBaseStorageMockified()
	storage.Update()
	storage.Update()

	assert.Equal(t, 2, storage.Length())
}
