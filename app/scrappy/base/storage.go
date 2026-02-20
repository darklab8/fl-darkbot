package base

import (
	"github.com/darklab8/fl-darkbot/app/scrappy/shared/parser"
	"github.com/darklab8/fl-darkbot/app/scrappy/shared/records"
	"github.com/darklab8/fl-darkstat/darkstat/configs_export"

	"github.com/darklab8/fl-darkbot/app/settings/logus"
)

type Base struct {
	Affiliation string
	Health      float64
	Tid         int
	Name        string
}

type BaseStorage struct {
	records.Records[records.StampedObjects[*configs_export.PoB]]
	api    IbaseAPI
	parser parser.Parser[records.StampedObjects[*configs_export.PoB]]
}

// Conveniently born some factory
func (b *BaseStorage) Update() {
	pobs, err := b.api.GetPobs()
	if logus.Log.CheckWarn(err, "quering API with error in BaseStorage") {
		return
	}

	if len(pobs) == 0 {
		logus.Log.Warn("skip pobs update if no pobs was received")
		return
	}

	record := records.NewStampedObjects[*configs_export.PoB]()
	for _, pob := range pobs {
		record.Add(pob)
	}
	if logus.Log.CheckWarn(err, "received bad parser parsing result in BaseStorage. Ignoring.") {
		return
	}
	b.Add(record)
	logus.Log.Info("updated base storage")
}

func NewBaseStorage(api IbaseAPI) *BaseStorage {
	b := &BaseStorage{}
	b.api = api
	return b
}
