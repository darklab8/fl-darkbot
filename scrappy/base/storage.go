package base

import "darkbot/scrappy/shared/api"

type BaseStorage struct {
	baseRecords
	api    api.APIinterface
	parser baseParser
}

// Conveniently born some factory
func (b *BaseStorage) Update() {
	data := b.api.New().GetData()
	record := b.parser.Parse(data)
	b.add(record)
}

var Storage baseRecords
