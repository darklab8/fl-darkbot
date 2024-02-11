package settings

import (
	"path/filepath"
	"strconv"

	"github.com/darklab8/fl-darkbot/app/settings/logus"
	"github.com/darklab8/fl-darkbot/app/settings/types"

	"github.com/darklab8/go-utils/goutils/utils"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

const (
	EnvFalse = "false"
	EnvTrue  = "true"
)

type ConfigScheme struct {
	ScrappyBaseUrl       types.APIurl `env:"SCRAPPY_BASE_URL" envDefault:"undefined"`
	ScrappyPlayerUrl     types.APIurl `env:"SCRAPPY_PLAYER_URL" envDefault:"undefined"`
	ScrappyBaseAttackUrl types.APIurl `env:"SCRAPPY_BASE_ATTACK_URL" envDefault:"https://discoverygc.com/forums/showthread.php?tid=110046&action=lastpost"`

	DiscorderBotToken string `env:"DISCORDER_BOT_TOKEN" envDefault:"undefined"`

	ConfiguratorDbname string `env:"CONFIGURATOR_DBNAME" envDefault:"dev"`

	ConsolerPrefix   string `env:"CONSOLER_PREFIX" envDefault:";"`
	ProfilingEnabled string `env:"PROFILING" envDefault:"false"`

	ScrappyLoopDelay string `env:"SCRAPPY_LOOP_DELAY" envDefault:"10"`
	ViewerLoopDelay  string `env:"VIEWER_LOOP_DELAY" envDefault:"10"`
	DevEnvMockApi    string `env:"DEVENV_MOCK_API" envDefault:"true"`
}

var ScrappyLoopDelay types.ScrappyLoopDelay
var ViewerLoopDelay types.ViewerLoopDelay
var Config ConfigScheme

var Dbpath types.Dbpath
var Workdir string

func NewDBPath(dbname string) types.Dbpath {
	return types.Dbpath(filepath.Join(Workdir, "data", dbname+".sqlite3"))
}

func load() {
	logus.Log.Info("identifying folder of settings")
	Workdir = filepath.Dir(filepath.Dir(utils.GetCurrentFolder().ToString()))

	err := godotenv.Load(filepath.Join(Workdir, ".env"))
	if err == nil {
		logus.Log.Info("loadded settings from .env")
	}

	opts := env.Options{RequiredIfNoDef: true}
	err = env.Parse(&Config, opts)

	logus.Log.CheckFatal(err, "settings have unset variable")

	logus.Log.Debug("settings were downloaded. Scrappy base url=", logus.APIUrl(Config.ScrappyBaseUrl))

	Dbpath = NewDBPath(Config.ConfiguratorDbname)

	scrappy_loop_delay, err := strconv.Atoi(Config.ScrappyLoopDelay)
	logus.Log.CheckFatal(err, "failed to parse ScrappyLoopDelay")
	ScrappyLoopDelay = types.ScrappyLoopDelay(scrappy_loop_delay)

	viewer_loop_delay, err := strconv.Atoi(Config.ViewerLoopDelay)
	logus.Log.CheckFatal(err, "failed to parse ViewerLoopDelay")
	ViewerLoopDelay = types.ViewerLoopDelay(viewer_loop_delay)

	logus.Log.Info("settings.ScrappyLoopDelay=", logus.ScrappyLoopDelay(ScrappyLoopDelay))
}

func init() {
	logus.Log.Info("attempt to load settings")
	load()
}
