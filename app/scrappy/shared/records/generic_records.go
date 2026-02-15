package records

import (
	"sync"

	"github.com/darklab8/fl-darkbot/app/settings/logus"

	"github.com/darklab8/go-utils/utils"
)

type Deletable interface {
}

type Records[T Deletable] struct {
	records []T
	mu      sync.Mutex
}

const recordLimit = 10

func (b *Records[T]) Add(record T) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.records = append(b.records, record)

	cutterStart := len(b.records) - recordLimit
	if cutterStart < 0 {
		cutterStart = 0
	}
	b.records = b.records[cutterStart:]
}

func (b *Records[T]) GetLatestRecord() (T, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if len(b.records) == 0 {
		var obj T
		return obj, utils.ErrorNotFound{}
	}
	logus.Log.Debug("records.GetLatestRecord", logus.Records(b.records))

	return b.records[len(b.records)-1], nil
}

func (b *Records[T]) List(callback func(values []T)) {
	b.mu.Lock()
	defer b.mu.Unlock()

	callback(b.records)
}

func (b *Records[T]) Length() int {
	b.mu.Lock()
	defer b.mu.Unlock()

	return len(b.records)
}
