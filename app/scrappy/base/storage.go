package base

import (
	"github.com/darklab8/fl-darkstat/darkstat/configs_export"

	"github.com/darklab8/fl-darkbot/app/scrappy/shared/records"
	"github.com/darklab8/fl-darkbot/app/settings/logus"
)

type Base struct {
	Affiliation string
	Health      float64
	Tid         int
	Name        string
}

// converted for old tests
func NewPoB1(base Base) *configs_export.PoB {
	return &configs_export.PoB{
		PoBCore: configs_export.PoBCore{
			Name:        base.Name,
			FactionName: &base.Affiliation,
			Health:      &base.Health,
		},
	}
}

type BaseStorage struct {
	records.Records[records.StampedObjects[*configs_export.PoB]]
	parser baseParser

	api BaseApi

	// pobs []*configs_export.PoB
	// err  error
}

type BaseData struct {
	List []*configs_export.PoB
}

type BaseApi interface {
	GetPobs() ([]*configs_export.PoB, error)
}

// func (b *BaseStorage) GetLatestRecord() (BaseData, error) {
// 	return BaseData{
// 		List: b.pobs,
// 	}, b.err
// }
// func (b *BaseStorage) Length() int {
// 	return len(b.pobs)
// }

// Conveniently born some factory
func (b *BaseStorage) Update() {
	pobs, err := b.api.GetPobs()
	if logus.Log.CheckWarn(err, "quering API with error in BaseStorage") {
		err = err
		return
	}

	b.parser.Parse(pobs)

	logus.Log.Info("updated base storage")
}

func NewBaseStorage(api BaseApi) *BaseStorage {
	b := &BaseStorage{}
	b.api = api
	return b
}
