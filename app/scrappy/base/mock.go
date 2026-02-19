package base

import (
	"encoding/json"
	"os"
	"path"

	"github.com/darklab8/fl-darkbot/app/scrappy/tests"
	"github.com/darklab8/fl-darkbot/app/settings/logus"
	"github.com/darklab8/fl-darkstat/darkstat/configs_export"
)

// SPY

type APIspy struct {
	Filename string
}

type apiBaseSpy struct {
	APIspy
}

func NewMock(filename string) IbaseAPI {
	return apiBaseSpy{APIspy{Filename: filename}}
}

func FixtureBaseApiMock() IbaseAPI {
	return NewMock("basedata.json")
}

func (a apiBaseSpy) GetPobs() ([]*configs_export.PoB, error) {
	path_testdata := tests.FixtureCreateTestDataFolder()
	path_testfile := path.Join(path_testdata, a.Filename)
	data, err := os.ReadFile(path_testfile)
	logus.Log.CheckFatal(err, "unable to read file")

	var pobs []*configs_export.PoB
	err = json.Unmarshal(data, &pobs)

	logus.Log.CheckFatal(err, "unable to unmrashal data")
	return pobs, nil
}
