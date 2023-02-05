package records

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCleanup(t *testing.T) {
	storage := Records[int]{}

	for i := 0; i < 12; i++ {
		storage.Add(i)
	}

	latest, _ := storage.GetLatestRecord()
	assert.Equal(t, 11, latest)
	assert.Equal(t, 10, storage.Length())

	storage.Add(100)
	latest, _ = storage.GetLatestRecord()
	assert.Equal(t, 100, latest)
	assert.Equal(t, 10, storage.Length())
}
