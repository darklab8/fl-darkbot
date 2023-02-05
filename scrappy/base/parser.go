package base

import (
	"darkbot/scrappy/shared/records"
	"darkbot/utils/logger"
	"encoding/json"
)

type baseSerializer struct {
	Affiliation string  `json:"affiliation"`
	Health      float64 `json:"health"`
	Tid         int     `json:"tid"`
}

type baseParser struct {
}

func (b baseParser) Parse(body []byte) records.StampedObjects[Base] {
	record := records.StampedObjects[Base]{}.New()

	bases := map[string]baseSerializer{}
	if err := json.Unmarshal(body, &bases); err != nil {
		logger.CheckPanic(err, "unable to unmarshal base request")
	}

	for name, serializer := range bases {
		record.Add(
			name,
			Base{
				Name:        name,
				Affiliation: serializer.Affiliation,
				Health:      serializer.Health,
				Tid:         serializer.Tid,
			},
		)
	}
	return record
}
