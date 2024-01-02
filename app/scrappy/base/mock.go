package base

import (
	"darkbot/app/scrappy/tests"
	"darkbot/app/settings/darkbot_logus"
	"os"
	"path"
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
	darkbot_logus.Log.CheckFatal(err, "unable to read file")
	return data, nil
}
