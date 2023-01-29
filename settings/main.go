package settings

import (
	"darkbot/utils"
	"path/filepath"

	"darkbot/dtypes"

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

type dbpath dtypes.Dbpath

var Dbpath dtypes.Dbpath
var Workdir string

func NewDBPath(dbname string) dtypes.Dbpath {
	return dtypes.Dbpath(filepath.Join(Workdir, "data", dbname+".sqlite3"))
}

func load() {
	utils.LogInfo("identifying folder of settings")
	Workdir = filepath.Dir(utils.GetCurrrentFolder())

	err := godotenv.Load(filepath.Join(Workdir, ".env"))
	if err == nil {
		utils.LogInfo("loadded settings from .env")
	}

	opts := env.Options{RequiredIfNoDef: true}
	err = env.Parse(&Config, opts)

	utils.CheckPanic(err, "settings have unset variable")

	utils.LogInfo("settings were downloaded. Scrappy base url=", Config.ScrappyBaseUrl)

	Dbpath = NewDBPath(Config.ConfiguratorDbname)
}

func init() {
	utils.LogInfo("attempt to load settings")
	load()
}
