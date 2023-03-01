package templ

import (
	"darkbot/scrappy/base"
	"darkbot/scrappy/shared/records"
	"darkbot/utils/logger"
	"darkbot/viewer/apis"
	"math"
)

type ErrorCalculatingDerivative struct {
	msg string
}

func (t ErrorCalculatingDerivative) Error() string {
	err_msg := "Some error calculating derivative :), msg=" + t.msg
	logger.Warn("ErrorCalculatingDerivative.err_mesg=", err_msg)
	return err_msg
}

type NoNonZeroDerivativesWarning struct {
}

func (t NoNonZeroDerivativesWarning) Error() string {
	logger.Warn("No unzero derivatives")
	return "No unzero derivatives"
}

func CalculateDerivates(tags []string, api apis.API) (map[string]float64, error) {
	baseHealths := make(map[string][]float64)
	var res_records []records.StampedObjects[base.Base]
	api.Scrappy.BaseStorage.Records.List(func(records2 []records.StampedObjects[base.Base]) {
		res_records = records2
	})

	if len(res_records) < 2 {
		return map[string]float64{}, ErrorCalculatingDerivative{msg: "amoung of records less than 2"}
	}

	TimeDiff := res_records[len(res_records)-1].Timestamp.Sub(res_records[0].Timestamp)

	if TimeDiff.Seconds() == 0 {
		return map[string]float64{}, ErrorCalculatingDerivative{msg: "insffucuient time diff. it is zero"}
	}

	for _, record := range res_records {
		bases := record.List

		for _, base := range bases {
			if _, ok := baseHealths[base.Name]; !ok {
				baseHealths[base.Name] = make([]float64, 0)
			}
			baseHealths[base.Name] = append(baseHealths[base.Name], base.Health)
		}
	}

	baseDerivatives := make(map[string]float64)
	wasThereNonZeroDeravatives := false
	for baseName, baseHealthHistory := range baseHealths {
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
				wasThereNonZeroDeravatives = true
			}
		}
	}

	if !wasThereNonZeroDeravatives {
		return baseDerivatives, NoNonZeroDerivativesWarning{}
	}

	return baseDerivatives, nil
}
