package player

import (
	"darkbot/scrappy/shared/api"
	"darkbot/scrappy/shared/parser"
	"darkbot/scrappy/shared/records"
	"darkbot/settings/utils/logger"
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
	data, err := b.api.GetData()
	if err != nil {
		logger.CheckWarn(err, "quering API with error in PlayerStorage")
		return
	}
	record, err := b.parser.Parse(data)
	if err != nil {
		logger.CheckWarn(err, "received bad parser parsing result in PlayerStorage. Ignoring.")
		return
	}
	b.Add(record)
	logger.Info("updated player storage")
}

func NewPlayerStorage(api api.APIinterface) *PlayerStorage {
	b := &PlayerStorage{}
	b.parser = playerParser{}
	b.api = api
	return b
}
