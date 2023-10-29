package base

import (
	"darkbot/settings/utils/logger"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func FixtureBaseStorageMockified() *BaseStorage {
	storage := (&BaseStorage{}).New()
	storage.Api = APIBasespy{}.New()
	return storage
}

func TestGetBases(t *testing.T) {
	storage := FixtureBaseStorageMockified()
	storage.Update()

	bases, err := storage.GetLatestRecord()
	logger.CheckPanic(err, "not found latest base record")

	assert.True(t, len(bases.List) > 0)
	fmt.Println(bases.List)
}

func TestAddManyRecords(t *testing.T) {

	storage := FixtureBaseStorageMockified()
	storage.Update()
	storage.Update()

	assert.Equal(t, 2, storage.Length())
}
