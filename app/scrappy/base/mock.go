package base

import (
	"os"
	"path"

	"github.com/darklab/fl-darkbot/app/scrappy/tests"
	"github.com/darklab/fl-darkbot/app/settings/logus"
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

func (a apiBaseSpy) GetBaseData() ([]byte, error) {
	path_testdata := tests.FixtureCreateTestDataFolder()
	path_testfile := path.Join(path_testdata, a.Filename)
	data, err := os.ReadFile(path_testfile)
	logus.Log.CheckFatal(err, "unable to read file")
	return data, nil
}
