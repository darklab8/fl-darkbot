package player

import (
	"darkbot/scrappy/shared/api"
	"darkbot/scrappy/tests"
	"darkbot/utils"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

func TestRegeneratePlayerData(t *testing.T) {
	utils.RegenerativeTest(func() error {
		data := PlayerAPI{}.New().GetData()
		path_testdata := tests.FixtureCreateTestDataFolder()
		path_testfile := path.Join(path_testdata, "playerdata.json")
		err := ioutil.WriteFile(path_testfile, data, os.ModePerm)
		utils.CheckPanic(err, "unable to write file")
		return nil
	})
}

// SPY

type APIPlayerSpy struct {
}

func (a APIPlayerSpy) New() api.APIinterface {
	return a
}

func (a APIPlayerSpy) GetData() []byte {
	path_testdata := tests.FixtureCreateTestDataFolder()
	path_testfile := path.Join(path_testdata, "playerdata.json")
	data, err := ioutil.ReadFile(path_testfile)
	utils.CheckPanic(err, "unable to read file")
	return data
}
