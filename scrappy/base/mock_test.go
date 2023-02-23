package base

import (
	"darkbot/scrappy/tests"
	"darkbot/utils"
	"darkbot/utils/logger"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

func TestRegenerateBaseData(t *testing.T) {
	utils.RegenerativeTest(func() error {
		data, _ := basesAPI{}.New().GetData()
		path_testdata := tests.FixtureCreateTestDataFolder()
		path_testfile := path.Join(path_testdata, "basedata.json")
		err := ioutil.WriteFile(path_testfile, data, os.ModePerm)
		logger.CheckPanic(err, "unable to write file")
		return nil
	})
}
