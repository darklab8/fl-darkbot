package base

import (
	"darkbot/scrappy/shared/api"
	"darkbot/scrappy/shared/parser"
	"darkbot/scrappy/shared/records"

	"darkbot/settings/utils/logger"
)

type Base struct {
	Affiliation string
	Health      float64
	Tid         int
	Name        string
}

type BaseStorage struct {
	records.Records[records.StampedObjects[Base]]
	Api    api.APIinterface
	parser parser.Parser[records.StampedObjects[Base]]
}

// Conveniently born some factory
func (b *BaseStorage) Update() {
	data, err := b.Api.GetData()
	if err != nil {
		logger.CheckWarn(err, "quering API with error in BaseStorage")
		return
	}
	record, err := b.parser.Parse(data)
	if err != nil {
		logger.CheckWarn(err, "received bad parser parsing result in BaseStorage. Ignoring.")
		return
	}
	b.Add(record)
	logger.Info("updated base storage")
}

func (b *BaseStorage) New() *BaseStorage {
	b.parser = baseParser{}
	b.Api = NewBaseApi()
	return b
}
