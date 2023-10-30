package player

import (
	"darkbot/scrappy/shared/api"
	"darkbot/scrappy/tests"
	"darkbot/settings/logus"
	"darkbot/settings/utils"
	"darkbot/settings/utils/logger"
	"os"
	"path"
	"testing"
)

func TestRegeneratePlayerData(t *testing.T) {
	utils.RegenerativeTest(func() error {
		data, err := NewPlayerAPI().GetData()
		logus.CheckError(err, "new player api get data errored")
		path_testdata := tests.FixtureCreateTestDataFolder()
		path_testfile := path.Join(path_testdata, "playerdata.json")
		err = os.WriteFile(path_testfile, data, os.ModePerm)
		logger.CheckPanic(err, "unable to write file")
		return nil
	})
}

// SPY

type APIPlayerSpy struct {
}

func (a APIPlayerSpy) New() api.APIinterface {
	return a
}

func (a APIPlayerSpy) GetData() ([]byte, error) {
	path_testdata := tests.FixtureCreateTestDataFolder()
	path_testfile := path.Join(path_testdata, "playerdata.json")
	data, err := os.ReadFile(path_testfile)
	logger.CheckPanic(err, "unable to read file")
	return data, nil
}
