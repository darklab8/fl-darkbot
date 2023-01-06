package base

import (
	"darkbot/scrappy/api"
	"darkbot/scrappy/tests"
	"darkbot/utils"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

func TestRegenerateBaseData(t *testing.T) {
	utils.RegenerativeTest(func() error {
		data := BasesAPI{}.New().GetData()
		path_testdata := tests.FixtureCreateTestDataFolder()
		path_testfile := path.Join(path_testdata, "basedata.json")
		err := ioutil.WriteFile(path_testfile, data, os.ModePerm)
		utils.CheckPanic(err, "unable to write file")
		return nil
	})
}

// SPY

type APIspy struct {
}

type APIBasespy struct {
	APIspy
}

func (a APIBasespy) New() api.APIinterface {
	return a
}

func (a APIBasespy) GetData() []byte {
	path_testdata := tests.FixtureCreateTestDataFolder()
	path_testfile := path.Join(path_testdata, "basedata.json")
	data, err := ioutil.ReadFile(path_testfile)
	utils.CheckPanic(err, "unable to read file")
	return data
}
