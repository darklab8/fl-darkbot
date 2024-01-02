package player

import (
	"darkbot/app/scrappy/shared/parser"
	"darkbot/app/scrappy/shared/records"
	"darkbot/app/settings/darkbot_logus"
)

type Player struct {
	Time   string
	Name   string
	System string
	Region string
}

type PlayerStorage struct {
	records.Records[records.StampedObjects[Player]]
	api    IPlayerAPI
	parser parser.Parser[records.StampedObjects[Player]]
}

// Conveniently born some factory
func (b *PlayerStorage) Update() {
	data, err := b.api.GetPlayerData()
	if darkbot_logus.Log.CheckWarn(err, "quering API with error in PlayerStorage") {
		return
	}
	record, err := b.parser.Parse(data)
	if darkbot_logus.Log.CheckWarn(err, "received bad parser parsing result in PlayerStorage. Ignoring.") {
		return
	}
	b.Add(record)
	darkbot_logus.Log.Info("updated player storage")
}

func NewPlayerStorage(api IPlayerAPI) *PlayerStorage {
	b := &PlayerStorage{}
	b.parser = playerParser{}
	b.api = api
	return b
}
