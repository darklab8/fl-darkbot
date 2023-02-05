package records

import (
	"darkbot/utils"
	"darkbot/utils/logger"
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

	cutterStart := len(b.records) - 10
	if cutterStart < 0 {
		cutterStart = 0
	}

	// Golang has bug with array of pointers, they will not be free until nilled
	// This code ensures proper deletion
	for index, _ := range b.records[:cutterStart] {
		b.records[index] = nil
	}

	b.records = b.records[cutterStart:]
}

func (b *Records[T]) GetLatestRecord() (*T, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if len(b.records) == 0 {
		return nil, utils.ErrorNotFound{}
	}
	logger.Info("records.GetLatestRecord.len(b.records)=", len(b.records))

	return b.records[len(b.records)-1], nil
}

func (b *Records[T]) Length() int {
	b.mu.Lock()
	defer b.mu.Unlock()

	return len(b.records)
}
