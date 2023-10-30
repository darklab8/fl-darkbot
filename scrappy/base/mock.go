package base

import (
	"darkbot/scrappy/shared/api"
	"darkbot/scrappy/tests"
	"darkbot/settings/utils/logger"
	"io/ioutil"
	"path"
)

// SPY

type APIspy struct {
	Filename string
}

type apiBaseSpy struct {
	APIspy
}

func NewMock(filename string) api.APIinterface {
	return apiBaseSpy{APIspy{Filename: filename}}
}

func NewBaseApiMock() api.APIinterface {
	return NewMock("basedata.json")
}

func (a apiBaseSpy) GetData() ([]byte, error) {
	path_testdata := tests.FixtureCreateTestDataFolder()
	path_testfile := path.Join(path_testdata, a.Filename)
	data, err := ioutil.ReadFile(path_testfile)
	logger.CheckPanic(err, "unable to read file")
	return data, nil
}
