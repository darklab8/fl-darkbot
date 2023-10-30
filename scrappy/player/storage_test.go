package player

import (
	"darkbot/settings/utils/logger"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func FixturePlayerStorageMockified() *PlayerStorage {
	storage := (&PlayerStorage{}).New()
	storage.Api = NewPlayerMockAPI()
	return storage
}

func TestGetPlayers(t *testing.T) {
	storage := FixturePlayerStorageMockified()
	storage.Update()

	bases, err := storage.GetLatestRecord()
	logger.CheckPanic(err, "not found latest base record")

	assert.True(t, len(bases.List) > 0)
	fmt.Println(bases.List)
}
