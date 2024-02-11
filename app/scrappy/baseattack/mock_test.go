package baseattack

import (
	"os"
	"path"
	"testing"

	"github.com/darklab/fl-darkbot/app/scrappy/tests"
	"github.com/darklab/fl-darkbot/app/settings/logus"

	"github.com/darklab8/darklab_goutils/goutils/utils"
)

func TestRegenerateBaseData(t *testing.T) {
	utils.RegenerativeTest(
		func() error {
			data, _ := NewBaseAttackAPI().GetBaseAttackData()
			path_testdata := tests.FixtureCreateTestDataFolder()
			path_testfile := path.Join(path_testdata, "data.json")
			err := os.WriteFile(path_testfile, data, os.ModePerm)
			logus.Log.CheckFatal(err, "unable to write file")
			return nil
		})
}
