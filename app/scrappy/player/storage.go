package player

import (
	"github.com/darklab8/fl-darkbot/app/scrappy/shared/parser"
	"github.com/darklab8/fl-darkbot/app/scrappy/shared/records"
	"github.com/darklab8/fl-darkbot/app/settings/logus"
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

	observers []Observer
}
type Observer interface {
	ReceivePlayers(p *PlayerStorage)
}

func (b *PlayerStorage) RegisterObserve(obs Observer) {
	b.observers = append(b.observers, obs)
}

func (b *PlayerStorage) UpdateObservers() {
	for _, obs := range b.observers {
		obs.ReceivePlayers(b)
	}
}

// Conveniently born some factory
func (b *PlayerStorage) Update() {
	data, err := b.api.GetPlayerData()
	if logus.Log.CheckWarn(err, "quering API with error in PlayerStorage") {
		return
	}
	record, err := b.parser.Parse(data)
	if logus.Log.CheckWarn(err, "received bad parser parsing result in PlayerStorage. Ignoring.") {
		return
	}
	b.Add(record)
	logus.Log.Info("updated player storage")

	b.UpdateObservers()
}

func NewPlayerStorage(api IPlayerAPI) *PlayerStorage {
	b := &PlayerStorage{}
	b.parser = playerParser{}
	b.api = api
	return b
}
