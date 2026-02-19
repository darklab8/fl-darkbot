package base

import (
	"encoding/json"
	"os"
	"path"
	"testing"

	"github.com/darklab8/fl-darkbot/app/scrappy/tests"
	"github.com/darklab8/fl-darkbot/app/settings/logus"
)

func TestRegenerateBaseData(t *testing.T) {
	pobs, err := NewBaseApi().GetPobs()
	logus.Log.CheckPanic(err, "failed to query api")
	path_testdata := tests.FixtureCreateTestDataFolder()
	path_testfile := path.Join(path_testdata, "basedata2.json")
	data, err := json.Marshal(pobs)
	logus.Log.CheckPanic(err, "failed to marshal data")
	err = os.WriteFile(path_testfile, data, os.ModePerm)
	logus.Log.CheckFatal(err, "unable to write file")
}
