package templ

import (
	"darkbot/scrappy/base"
	"darkbot/scrappy/shared/records"
	"darkbot/viewer/apis"
	"math"
)

func CalculateDerivates(tags []string, api apis.API) map[string]float64 {
	baseHealths := make(map[string][]float64)
	var res_records []records.StampedObjects[base.Base]
	api.Scrappy.BaseStorage.Records.List(func(records2 []records.StampedObjects[base.Base]) {
		res_records = records2
	})

	if len(res_records) < 2 {
		return map[string]float64{}
	}

	TimeDiff := res_records[len(res_records)-1].Timestamp.Sub(res_records[0].Timestamp)

	if TimeDiff.Seconds() == 0 {
		return map[string]float64{}
	}

	for _, record := range res_records {
		bases := MatchBases(record, tags)

		for _, base := range bases {
			if _, ok := baseHealths[base.Name]; !ok {
				baseHealths[base.Name] = make([]float64, 0)
			}
			baseHealths[base.Name] = append(baseHealths[base.Name], base.Health)
		}
	}

	baseDerivatives := make(map[string]float64)
	for baseName, baseHealthHistory := range baseHealths {
		var normalizer int = 0
		if len(baseHealthHistory) <= 1 {
			continue
		}

		if _, ok := baseDerivatives[baseName]; !ok {
			baseDerivatives[baseName] = 0
		}

		for i := 0; i < len(baseHealthHistory)-1; i++ {
			derivative := baseHealthHistory[i+1] - baseHealthHistory[i]
			if math.Abs(derivative) > 1e-10 {
				baseDerivatives[baseName] = derivative
				normalizer++
			}
		}

		if normalizer == 0 {
			normalizer = 1
		}
		baseDerivatives[baseName] = baseDerivatives[baseName] / float64(normalizer)

	}

	return baseDerivatives
}
