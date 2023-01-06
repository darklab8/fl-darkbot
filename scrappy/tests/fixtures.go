package tests

import (
	"os"
	"path"
	"path/filepath"
	"runtime"
)

func FixtureCreateTestDataFolder() string {
	_, filename, _, _ := runtime.Caller(2)
	path_curr_folder := filepath.Dir(filename)
	path_testdata := path.Join(path_curr_folder, "testdata")
	os.Mkdir(path_testdata, os.ModePerm)
	return path_testdata
}
