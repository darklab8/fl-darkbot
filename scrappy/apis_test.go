package scrappy

import (
	"darkbot/utils"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

func FixtureCreateTestDataFolder() string {
	path_curr_folder := utils.GetCurrrentFolder()
	path_testdata := path.Join(path_curr_folder, "testdata")
	os.Mkdir(path_testdata, os.ModePerm)
	return path_testdata
}

func TestRegenerateBaseData(t *testing.T) {
	utils.RegenerativeTest(func() error {
		data := BasesAPI{}.New().GetData()
		path_testdata := FixtureCreateTestDataFolder()
		path_testfile := path.Join(path_testdata, "basedata.json")
		err := ioutil.WriteFile(path_testfile, data, os.ModePerm)
		utils.CheckPanic(err, "unable to write file")
		return nil
	})
}

func TestRegeneratePlayerData(t *testing.T) {
	utils.RegenerativeTest(func() error {
		data := PlayerAPI{}.New().GetData()
		path_testdata := FixtureCreateTestDataFolder()
		path_testfile := path.Join(path_testdata, "playerdata.json")
		err := ioutil.WriteFile(path_testfile, data, os.ModePerm)
		utils.CheckPanic(err, "unable to write file")
		return nil
	})
}

// Test Data

type APIspy struct {
}

func (a APIBasespy) New() APIinterface {
	return a
}

func (a APIPlayerSpy) New() APIinterface {
	return a
}

type APIBasespy struct {
	APIspy
}

func (a APIBasespy) GetData() []byte {
	path_testdata := FixtureCreateTestDataFolder()
	path_testfile := path.Join(path_testdata, "basedata.json")
	data, err := ioutil.ReadFile(path_testfile)
	utils.CheckPanic(err, "unable to read file")
	return data
}

type APIPlayerSpy struct {
	APIspy
}

func (a APIPlayerSpy) GetData() []byte {
	path_testdata := FixtureCreateTestDataFolder()
	path_testfile := path.Join(path_testdata, "playerdata.json")
	data, err := ioutil.ReadFile(path_testfile)
	utils.CheckPanic(err, "unable to read file")
	return data
}
