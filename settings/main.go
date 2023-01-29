package settings

import (
	"darkbot/utils"
	"path/filepath"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

type ConfigScheme struct {
	// Example how to add `env:"S3_BUCKET_NAME" envDefault:"default_value"`
	ScrappyBaseUrl   string `env:"SCRAPPY_BASE_URL"`
	ScrappyPlayerUrl string `env:"SCRAPPY_PLAYER_URL"`

	DiscorderBotToken string `env:"DISCORDER_BOT_TOKEN"`

	ConfiguratorDbname string `env:"CONFIGURATOR_DBNAME"`

	ConsolerPrefix string `env:"CONSOLER_PREFIX" envDefault:","`
}

var Config ConfigScheme

var Dbpath string

func load() {
	utils.LogInfo("identifying folder of settings")
	workdir := filepath.Dir(utils.GetCurrrentFolder())

	err := godotenv.Load(filepath.Join(workdir, ".env"))
	if err == nil {
		utils.LogInfo("loadded settings from .env")
	}

	opts := env.Options{RequiredIfNoDef: true}
	err = env.Parse(&Config, opts)

	utils.CheckPanic(err, "settings have unset variable")

	utils.LogInfo("settings were downloaded. Scrappy base url=", Config.ScrappyBaseUrl)

	Dbpath = filepath.Join(workdir, "data", Config.ConfiguratorDbname+".sqlite3")
}

func init() {
	utils.LogInfo("attempt to load settings")
	load()
}
