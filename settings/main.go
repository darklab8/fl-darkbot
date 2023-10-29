package settings

import (
	"darkbot/utils"
	"darkbot/utils/logger"
	"path/filepath"
	"strconv"

	"darkbot/dtypes"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

const (
	EnvFalse = "false"
	EnvTrue  = "true"
)

type ConfigScheme struct {
	ScrappyBaseUrl       string `env:"SCRAPPY_BASE_URL" envDefault:"undefined"`
	ScrappyPlayerUrl     string `env:"SCRAPPY_PLAYER_URL" envDefault:"undefined"`
	ScrappyBaseAttackUrl string `env:"SCRAPPY_BASE_ATTACK_URL" envDefault:"https://discoverygc.com/forums/showthread.php?tid=110046&action=lastpost"`

	DiscorderBotToken string `env:"DISCORDER_BOT_TOKEN" envDefault:"undefined"`

	ConfiguratorDbname string `env:"CONFIGURATOR_DBNAME" envDefault:"dev"`

	ConsolerPrefix   string `env:"CONSOLER_PREFIX" envDefault:","`
	ProfilingEnabled string `env:"PROFILING" envDefault:"false"`

	LoopDelay     string `env:"LOOP_DELAY" envDefault:"10"`
	DevEnvMockApi string `env:"DEVENV_MOCK_API" envDefault:"true"`
}

var LoopDelay int
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

	LoopDelay, err = strconv.Atoi(Config.LoopDelay)
	logger.CheckPanic(err, "failed to parse LoopDelay")
	logger.Info("settings.LoopDelay=", LoopDelay)
}

func init() {
	logger.Info("attempt to load settings")
	load()
}
