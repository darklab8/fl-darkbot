package baseattack

import (
	"darkbot/scrappy/shared/api"
	"darkbot/scrappy/tests"
	"darkbot/utils/logger"
	"io/ioutil"
	"path"
)

// SPY

type APIbasis struct {
	Filename string
}

type BaseAttackAPISpy struct {
	APIbasis
}

func NewMock(filename string) api.APIinterface {
	return BaseAttackAPISpy{APIbasis{Filename: filename}}
}

func (a BaseAttackAPISpy) New() api.APIinterface {
	return NewMock("basedata.json")
}

func (a BaseAttackAPISpy) GetData() []byte {
	path_testdata := tests.FixtureCreateTestDataFolder()
	path_testfile := path.Join(path_testdata, a.Filename)
	data, err := ioutil.ReadFile(path_testfile)
	logger.CheckPanic(err, "unable to read file")
	return data
}
