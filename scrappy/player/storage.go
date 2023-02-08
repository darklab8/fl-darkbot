package player

import (
	"darkbot/scrappy/shared/api"
	"darkbot/scrappy/shared/parser"
	"darkbot/scrappy/shared/records"
	"darkbot/utils/logger"
)

type Player struct {
	Time   string
	Name   string
	System string
	Region string
}

type PlayerStorage struct {
	records.Records[records.StampedObjects[Player]]
	Api    api.APIinterface
	parser parser.Parser[records.StampedObjects[Player]]
}

// Conveniently born some factory
func (b *PlayerStorage) Update() {
	data := b.Api.New().GetData()
	record := b.parser.Parse(data)
	b.Add(record)
	logger.Info("updated player storage")
}

func (b *PlayerStorage) New() *PlayerStorage {
	b.parser = playerParser{}
	b.Api = PlayerAPI{}
	return b
}
