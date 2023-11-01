package records

import "time"

type StampedObjects[T interface{}] struct {
	List      []T
	Timestamp time.Time
}

func NewStampedObjects[T interface{}]() StampedObjects[T] {
	b := StampedObjects[T]{}
	b.Timestamp = time.Now()
	return b
}

func (b *StampedObjects[T]) Add(obj T) {
	b.List = append(b.List, obj)
}
