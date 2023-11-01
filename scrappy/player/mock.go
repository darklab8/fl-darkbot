package player

import (
	"darkbot/scrappy/tests"
	"darkbot/settings/logus"
	"darkbot/settings/utils"
	"os"
	"path"
	"testing"
)

func TestRegeneratePlayerData(t *testing.T) {
	utils.RegenerativeTest(func() error {
		data, err := NewPlayerAPI().GetPlayerData()
		logus.CheckError(err, "new player api get data errored")
		path_testdata := tests.FixtureCreateTestDataFolder()
		path_testfile := path.Join(path_testdata, "playerdata.json")
		err = os.WriteFile(path_testfile, data, os.ModePerm)
		logus.CheckFatal(err, "unable to write file")
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
	logus.CheckFatal(err, "unable to read file")
	return data, nil
}
