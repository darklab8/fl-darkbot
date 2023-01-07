package records

import "darkbot/utils"

type Records[T interface{}] struct {
	records []*T
}

func (b *Records[T]) Add(record T) {
	b.records = append(b.records, &record)
}

func (b *Records[T]) GetLatestRecord() (*T, error) {
	if len(b.records) == 0 {
		return nil, utils.ErrorNotFound{}
	}

	return b.records[len(b.records)-1], nil
}

func (b *Records[T]) Length() int {
	return len(b.records)
}
