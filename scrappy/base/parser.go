package base

import (
	"darkbot/scrappy/shared/records"
	"darkbot/settings/logus"
	"encoding/json"
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
		logus.CheckWarn(err, "unable to unmarshal base request", logus.Body(body))
		logus.Warn("unable to marshal base body=")
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
