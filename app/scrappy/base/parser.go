package base

import (
	"encoding/json"

	"github.com/darklab8/fl-darkbot/app/scrappy/shared/records"
	"github.com/darklab8/fl-darkbot/app/settings/logus"
)

type baseSerializer struct {
	Affiliation string  `json:"affiliation"`
	Health      float64 `json:"health"`
	Tid         int     `json:"tid"`
}

type baseParser struct {
}

func (b baseParser) Parse(body []byte) (records.StampedObjects[Base], error) {
	record := records.NewStampedObjects[Base]()

	bases := map[string]baseSerializer{}
	if err := json.Unmarshal(body, &bases); err != nil {
		logus.Log.CheckWarn(err, "unable to unmarshal base request", logus.Body(body))
		logus.Log.Warn("unable to marshal base body=")
		return record, err
	}

	for name, serializer := range bases {
		record.Add(
			Base{
				Name:        name,
				Affiliation: serializer.Affiliation,
				Health:      serializer.Health,
				Tid:         serializer.Tid,
			},
		)
	}
	return record, nil
}
