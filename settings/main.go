package settings

import (
	"darkbot/utils"
	"darkbot/utils/logger"
	"path/filepath"

	"darkbot/dtypes"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

const (
	EnvFalse = "false"
	EnvTrue  = "true"
)

type ConfigScheme struct {
	// Example how to add `env:"S3_BUCKET_NAME" envDefault:"default_value"`
	ScrappyBaseUrl   string `env:"SCRAPPY_BASE_URL"`
	ScrappyPlayerUrl string `env:"SCRAPPY_PLAYER_URL"`

	DiscorderBotToken string `env:"DISCORDER_BOT_TOKEN"`

	ConfiguratorDbname string `env:"CONFIGURATOR_DBNAME"`

	ConsolerPrefix   string `env:"CONSOLER_PREFIX" envDefault:","`
	ProfilingEnabled string `env:"PROFILING" envDefault:"false"`
}

var Config ConfigScheme

type dbpath dtypes.Dbpath

var Dbpath dtypes.Dbpath
var Workdir string

func NewDBPath(dbname string) dtypes.Dbpath {
	return dtypes.Dbpath(filepath.Join(Workdir, "data", dbname+".sqlite3"))
}

func load() {
	logger.Info("identifying folder of settings")
	Workdir = filepath.Dir(utils.GetCurrrentFolder())

	err := godotenv.Load(filepath.Join(Workdir, ".env"))
	if err == nil {
		logger.Info("loadded settings from .env")
	}

	opts := env.Options{RequiredIfNoDef: true}
	err = env.Parse(&Config, opts)

	logger.CheckPanic(err, "settings have unset variable")

	logger.Info("settings were downloaded. Scrappy base url=", Config.ScrappyBaseUrl)

	Dbpath = NewDBPath(Config.ConfiguratorDbname)
}

func init() {
	logger.Info("attempt to load settings")
	load()
}
