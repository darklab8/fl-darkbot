package records

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type DeletableInt struct {
	Value int
}

func (d DeletableInt) Delete() {

}

func TestCleanup(t *testing.T) {
	storage := Records[DeletableInt]{}

	for i := 0; i < 12; i++ {
		storage.Add(DeletableInt{i})
	}

	latest, _ := storage.GetLatestRecord()
	assert.Equal(t, 11, latest.Value)
	assert.Equal(t, 10, storage.Length())

	storage.Add(DeletableInt{100})
	latest, _ = storage.GetLatestRecord()
	assert.Equal(t, 100, latest.Value)
	assert.Equal(t, 10, storage.Length())
}
