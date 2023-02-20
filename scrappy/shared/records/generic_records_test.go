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

	for i := 0; i < recordLimit+2; i++ {
		storage.Add(DeletableInt{i})
	}

	latest, _ := storage.GetLatestRecord()
	assert.Equal(t, recordLimit+1, (latest).Value)
	assert.Equal(t, recordLimit, storage.Length())

	storage.Add(DeletableInt{100})
	latest, _ = storage.GetLatestRecord()
	assert.Equal(t, 100, (latest).Value)
	assert.Equal(t, recordLimit, storage.Length())
}
