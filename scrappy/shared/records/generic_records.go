package records

import (
	"darkbot/utils"
	"darkbot/utils/logger"
	"sync"
)

type Deletable interface {
}

type Records[T Deletable] struct {
	records []*T
	mu      sync.Mutex
}

func (b *Records[T]) Add(record T) []*T {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.records = append(b.records, &record)

	cutterStart := len(b.records) - 10
	if cutterStart < 0 {
		cutterStart = 0
	}

	returnable := b.records[:cutterStart]
	b.records = b.records[cutterStart:]
	return returnable
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
