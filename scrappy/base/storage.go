package base

import (
	"darkbot/scrappy/shared/api"
	"darkbot/scrappy/shared/parser"
	"darkbot/scrappy/shared/records"
	"darkbot/utils"
)

type Base struct {
	Affiliation string
	Health      float64
	Tid         int
	name        string
}

type BaseStorage struct {
	records.Records[records.StampedObjects[Base]]
	api    api.APIinterface
	parser parser.Parser[records.StampedObjects[Base]]
}

// Conveniently born some factory
func (b *BaseStorage) Update() {
	data := b.api.New().GetData()
	record := b.parser.Parse(data)
	b.Add(record)
	utils.LogInfo("updated base storage")
}

func (b BaseStorage) New() BaseStorage {
	b.parser = baseParser{}
	b.api = basesAPI{}
	return b
}
