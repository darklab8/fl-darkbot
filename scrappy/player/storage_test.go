package player

import (
	"darkbot/utils/logger"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func FixtureBaseStorageMockified() *PlayerStorage {
	storage := (&PlayerStorage{}).New()
	storage.api = APIPlayerSpy{}
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
