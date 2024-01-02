package base

import (
	"darkbot/app/scrappy/tests"
	"darkbot/app/settings/darkbot_logus"
	"os"
	"path"
	"testing"

	"github.com/darklab8/darklab_goutils/goutils/utils"
)

func TestRegenerateBaseData(t *testing.T) {
	utils.RegenerativeTest(func() error {
		data, _ := NewBaseApi().GetBaseData()
		path_testdata := tests.FixtureCreateTestDataFolder()
		path_testfile := path.Join(path_testdata, "basedata.json")
		err := os.WriteFile(path_testfile, data, os.ModePerm)
		darkbot_logus.Log.CheckFatal(err, "unable to write file")
		return nil
	})
}
