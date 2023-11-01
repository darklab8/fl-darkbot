package baseattack

import (
	"darkbot/scrappy/tests"
	"darkbot/settings/logus"
	"darkbot/settings/utils"
	"os"
	"path"
	"testing"
)

func TestRegenerateBaseData(t *testing.T) {
	utils.RegenerativeTest(
		func() error {
			data, _ := NewBaseAttackAPI().GetBaseAttackData()
			path_testdata := tests.FixtureCreateTestDataFolder()
			path_testfile := path.Join(path_testdata, "data.json")
			err := os.WriteFile(path_testfile, data, os.ModePerm)
			logus.CheckFatal(err, "unable to write file")
			return nil
		})
}
