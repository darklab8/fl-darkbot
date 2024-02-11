package base

import (
	"os"
	"path"
	"testing"

	"github.com/darklab8/fl-darkbot/app/scrappy/tests"
	"github.com/darklab8/fl-darkbot/app/settings/logus"

	"github.com/darklab8/go-utils/goutils/utils"
)

func TestRegenerateBaseData(t *testing.T) {
	utils.RegenerativeTest(func() error {
		data, _ := NewBaseApi().GetBaseData()
		path_testdata := tests.FixtureCreateTestDataFolder()
		path_testfile := path.Join(path_testdata, "basedata.json")
		err := os.WriteFile(path_testfile, data, os.ModePerm)
		logus.Log.CheckFatal(err, "unable to write file")
		return nil
	})
}
