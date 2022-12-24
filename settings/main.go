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
	workdir := filepath.Dir(utils.GetCurrrentFolder())
	AutogitSettingsPath := filepath.Join(workdir, ".settings.yml")

	file, err := ioutil.ReadFile(AutogitSettingsPath)
	utils.CheckFatal(err, "Could not read the file due to error, settings_path=%s\n", AutogitSettingsPath)

	err = yaml.Unmarshal(file, &Config)
	utils.CheckFatal(err, "unable to unmarshal settings")
}

func init() {
	load()
}
