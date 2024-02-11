package baseattack

import (
	"os"
	"path"

	"github.com/darklab/fl-darkbot/app/scrappy/tests"
	"github.com/darklab/fl-darkbot/app/settings/logus"
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
	logus.Log.CheckFatal(err, "unable to read file")
	return data, nil
}
