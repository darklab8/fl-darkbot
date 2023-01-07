package settings

import (
	"darkbot/utils"
	"io/ioutil"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type ConfigScheme struct {
	Scrappy struct {
		Player struct {
			URL string `yaml:"url"`
		} `yaml:"player"`
		Base struct {
			URL string `yaml:"url"`
		} `yaml:"base"`
		Forum struct {
			Username string `yaml:"username"`
			Password string `yaml:"password"`
		} `yaml:"forum"`
	} `yaml:"scrappy"`
	Listener struct {
		Discord struct {
			Bot struct {
				Token string `yaml:"token"`
			} `yaml:"bot"`
		} `yaml:"discord"`
	} `yaml:"listener"`
	Discorder struct {
		Discord struct {
			Bot struct {
				Token string `yaml:"token"`
			} `yaml:"bot"`
		} `yaml:"discord"`
	} `yaml:"discorder"`
}

var Config ConfigScheme

func load() {
	utils.LogInfo("identifying folder of settings")
	workdir := filepath.Dir(utils.GetCurrrentFolder())
	AutogitSettingsPath := filepath.Join(workdir, ".settings.yml")
	utils.LogInfo("settings folder is ", AutogitSettingsPath)

	file, err := ioutil.ReadFile(AutogitSettingsPath)
	utils.LogInfo("Reading file settings")
	utils.CheckFatal(err, "Could not read the file due to error, settings_path=%s\n", AutogitSettingsPath)

	err = yaml.Unmarshal(file, &Config)
	utils.CheckFatal(err, "unable to unmarshal settings")

	utils.LogInfo("settings were downloaded. Scrappy base url=", Config.Scrappy.Base.URL)
}

func init() {
	utils.LogInfo("attempt to load settings")
	load()
}
