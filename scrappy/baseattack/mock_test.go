package baseattack

import (
	"darkbot/scrappy/tests"
	"darkbot/settings/utils"
	"darkbot/settings/utils/logger"
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
			logger.CheckPanic(err, "unable to write file")
			return nil
		})
}
