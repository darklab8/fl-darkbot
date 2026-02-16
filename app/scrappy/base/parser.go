package base

import (
	"github.com/darklab8/fl-darkbot/app/scrappy/shared/records"
	"github.com/darklab8/fl-darkstat/darkstat/configs_export"
)

type baseParser struct {
}

func (b baseParser) Parse(pobs []*configs_export.PoB) (records.StampedObjects[*configs_export.PoB], error) {
	record := records.NewStampedObjects[*configs_export.PoB]()

	for _, base := range pobs {
		record.Add(base)
	}
	return record, nil
}
