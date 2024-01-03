package base

import (
	"darkbot/app/scrappy/shared/parser"
	"darkbot/app/scrappy/shared/records"

	"darkbot/app/settings/logus"
)

type Base struct {
	Affiliation string
	Health      float64
	Tid         int
	Name        string
}

type BaseStorage struct {
	records.Records[records.StampedObjects[Base]]
	api    IbaseAPI
	parser parser.Parser[records.StampedObjects[Base]]
}

// Conveniently born some factory
func (b *BaseStorage) Update() {
	data, err := b.api.GetBaseData()
	if logus.Log.CheckWarn(err, "quering API with error in BaseStorage") {
		return
	}
	record, err := b.parser.Parse(data)
	if logus.Log.CheckWarn(err, "received bad parser parsing result in BaseStorage. Ignoring.") {
		return
	}
	b.Add(record)
	logus.Log.Info("updated base storage")
}

func NewBaseStorage(api IbaseAPI) *BaseStorage {
	b := &BaseStorage{}
	b.parser = baseParser{}
	b.api = api
	return b
}
