package baseattack

import (
	"darkbot/app/scrappy/tests"
	"darkbot/app/settings/darkbot_logus"
	"os"
	"path"
)

// SPY

type APIbasis struct {
	Filename string
}

type BaseAttackAPISpy struct {
	APIbasis
}

func NewMock(filename string) IbaseAttackAPI {
	return BaseAttackAPISpy{APIbasis{Filename: filename}}
}

func FixtureBaseAttackAPIMock() IbaseAttackAPI {
	return NewMock("data.json")
}

func (a BaseAttackAPISpy) GetBaseAttackData() ([]byte, error) {
	path_testdata := tests.FixtureCreateTestDataFolder()
	path_testfile := path.Join(path_testdata, a.Filename)
	data, err := os.ReadFile(path_testfile)
	darkbot_logus.Log.CheckFatal(err, "unable to read file")
	return data, nil
}
