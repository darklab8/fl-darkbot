package player

import (
	"os"
	"path"
	"testing"

	"github.com/darklab/fl-darkbot/app/scrappy/tests"
	"github.com/darklab/fl-darkbot/app/settings/logus"

	"github.com/darklab8/darklab_goutils/goutils/utils"
)

func TestRegeneratePlayerData(t *testing.T) {
	utils.RegenerativeTest(func() error {
		data, err := NewPlayerAPI().GetPlayerData()
		logus.Log.CheckFatal(err, "new player api get data errored")
		path_testdata := tests.FixtureCreateTestDataFolder()
		path_testfile := path.Join(path_testdata, "playerdata.json")
		err = os.WriteFile(path_testfile, data, os.ModePerm)
		logus.Log.CheckFatal(err, "unable to write file")
		return nil
	})
}

// SPY

type apiPlayerSpy struct {
}

func FixturePlayerAPIMock() IPlayerAPI {
	return apiPlayerSpy{}
}

func (a apiPlayerSpy) GetPlayerData() ([]byte, error) {
	path_testdata := tests.FixtureCreateTestDataFolder()
	path_testfile := path.Join(path_testdata, "playerdata.json")
	data, err := os.ReadFile(path_testfile)
	logus.Log.CheckFatal(err, "unable to read file")
	return data, nil
}
