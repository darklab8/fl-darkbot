package base

import "darkbot/scrappy/apiRawData"

type BaseStorage struct {
	baseRecords
	api    apiRawData.APIinterface
	parser baseParser
}

// Conveniently born some factory
func (b *BaseStorage) Update() {
	data := b.api.New().GetData()
	record := b.parser.Parse(data)
	b.add(record)
}

var Storage baseRecords
