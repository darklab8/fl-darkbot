package player

import (
	"darkbot/scrappy/shared/api"
	"darkbot/scrappy/shared/parser"
	"darkbot/scrappy/shared/records"
	"darkbot/utils"
)

type Player struct {
	Time   string
	Name   string
	System string
	Region string
}

type PlayerStorage struct {
	records.Records[records.StampedObjects[Player]]
	api    api.APIinterface
	parser parser.Parser[records.StampedObjects[Player]]
}

// Conveniently born some factory
func (b *PlayerStorage) Update() {
	data := b.api.New().GetData()
	record := b.parser.Parse(data)
	b.Add(record)
	utils.LogInfo("updated player storage")
}

func (b *PlayerStorage) New() *PlayerStorage {
	b.parser = playerParser{}
	b.api = PlayerAPI{}
	return b
}
