package base

import (
	"darkbot/utils"
	"time"
)

type Base struct {
	Affiliation string
	Health      float64
	Tid         int
	name        string
}

type basesStampedRecord struct {
	dict      map[string]*Base
	list      []*Base
	timestamp time.Time
}

func (b basesStampedRecord) new() basesStampedRecord {
	b.timestamp = time.Now()
	b.dict = make(map[string]*Base)
	b.list = make([]*Base, 0)
	return b
}

func (b *basesStampedRecord) addBase(base Base) {
	b.dict[base.name] = &base
	b.list = append(b.list, &base)
}

type baseRecords struct {
	records []*basesStampedRecord
}

func (b *baseRecords) add(record basesStampedRecord) {
	b.records = append(b.records, &record)
}

func (b *baseRecords) GetLatestRecord() (*basesStampedRecord, error) {
	if len(b.records) == 0 {
		return nil, utils.ErrorNotFound{}
	}

	return b.records[len(b.records)-1], nil
}

func (b *baseRecords) Length() int {
	return len(b.records)
}
