package base

import (
	"darkbot/utils"
	"encoding/json"
)

type baseSerializer struct {
	Affiliation string  `json:"affiliation"`
	Health      float64 `json:"health"`
	Tid         int     `json:"tid"`
}

type baseParser struct {
}

func (b baseParser) Parse(body []byte) basesStampedRecord {
	record := basesStampedRecord{}.new()

	var bases map[string]baseSerializer
	if err := json.Unmarshal(body, &bases); err != nil {
		utils.CheckPanic(err, "unable to unmarshal base request")
	}

	for name, serializer := range bases {
		record.addBase(Base{
			name:        name,
			Affiliation: serializer.Affiliation,
			Health:      serializer.Health,
			Tid:         serializer.Tid,
		})
	}
	return record
}
