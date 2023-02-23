package base

import (
	"darkbot/scrappy/shared/api"
	"darkbot/scrappy/tests"
	"darkbot/utils/logger"
	"io/ioutil"
	"path"
)

// SPY

type APIspy struct {
	Filename string
}

type APIBasespy struct {
	APIspy
}

func NewMock(filename string) api.APIinterface {
	return APIBasespy{APIspy{Filename: filename}}
}

func (a APIBasespy) New() api.APIinterface {
	return NewMock("basedata.json")
}

func (a APIBasespy) GetData() ([]byte, error) {
	path_testdata := tests.FixtureCreateTestDataFolder()
	path_testfile := path.Join(path_testdata, a.Filename)
	data, err := ioutil.ReadFile(path_testfile)
	logger.CheckPanic(err, "unable to read file")
	return data, nil
}
