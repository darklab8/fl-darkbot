package records

import (
	"darkbot/utils"
	"sync"
)

type Records[T interface{}] struct {
	records []*T
	mu      sync.Mutex
}

func (b *Records[T]) Add(record T) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.records = append(b.records, &record)
}

func (b *Records[T]) GetLatestRecord() (*T, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if len(b.records) == 0 {
		return nil, utils.ErrorNotFound{}
	}

	return b.records[len(b.records)-1], nil
}

func (b *Records[T]) Length() int {
	b.mu.Lock()
	defer b.mu.Unlock()

	return len(b.records)
}
