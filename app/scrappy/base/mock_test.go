package base

import (
	"darkbot/app/scrappy/tests"
	"darkbot/app/settings/logus"
	"darkbot/app/settings/utils"
	"os"
	"path"
	"testing"
)

func TestRegenerateBaseData(t *testing.T) {
	utils.RegenerativeTest(func() error {
		data, _ := NewBaseApi().GetBaseData()
		path_testdata := tests.FixtureCreateTestDataFolder()
		path_testfile := path.Join(path_testdata, "basedata.json")
		err := os.WriteFile(path_testfile, data, os.ModePerm)
		logus.CheckFatal(err, "unable to write file")
		return nil
	})
}
