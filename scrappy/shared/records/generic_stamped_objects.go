package records

import "time"

type StampedObjects[T interface{}] struct {
	Dict      map[string]*T
	List      []*T
	Timestamp time.Time
}

func (b StampedObjects[T]) New() StampedObjects[T] {
	b.Timestamp = time.Now()
	b.Dict = make(map[string]*T)
	b.List = make([]*T, 0)
	return b
}

func (b *StampedObjects[T]) Add(key string, obj T) {
	b.Dict[key] = &obj
	b.List = append(b.List, &obj)
}