package scrappy

import (
	"darkbot/utils"
	"encoding/json"
	"time"
)

type baseSerializer struct {
	Affiliation string  `json:"affiliation"`
	Health      float64 `json:"health"`
	Tid         int     `json:"tid"`
}

type Base struct {
	Affiliation string
	Health      float64
	Tid         int
	name        string
}

type BasesStampedRecord struct {
	dict      map[string]*Base
	list      []*Base
	timestamp time.Time
}

func (b BasesStampedRecord) New() BasesStampedRecord {
	b.timestamp = time.Now()
	b.dict = make(map[string]*Base)
	b.list = make([]*Base, 0)
	return b
}

func (b BasesStampedRecord) NewFromAPI(body []byte) BasesStampedRecord {
	b = b.New()

	var bases map[string]Base
	if err := json.Unmarshal(body, &bases); err != nil {
		utils.CheckPanic(err, "unable to unmarshal base request")
	}

	for name, base := range bases {
		base.name = name

		b.dict[name] = &base
		b.list = append(b.list, &base)
	}

	return b
}

type BaseRecords struct {
	records []*BasesStampedRecord
	api     APIinterface
}

var BaseStorage BaseRecords

func (b *BaseRecords) addFromAPI() {
	data := b.api.New().GetData()
	record := BasesStampedRecord{}.NewFromAPI(data)
	b.add(record)
}

func (b *BaseRecords) add(record BasesStampedRecord) {
	b.records = append(b.records, &record)
}

func (b *BaseRecords) getLatestRecord() (*BasesStampedRecord, error) {
	if len(b.records) == 0 {
		return nil, utils.ErrorNotFound{}
	}

	return b.records[len(b.records)-1], nil
}
