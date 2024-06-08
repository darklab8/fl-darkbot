package settings

import (
	"path/filepath"

	"github.com/darklab8/fl-darkbot/app/settings/logus"
	"github.com/darklab8/fl-darkbot/app/settings/types"

	"github.com/darklab8/go-utils/utils"
	"github.com/darklab8/go-utils/utils/utils_env"
	"github.com/darklab8/go-utils/utils/utils_settings"
)

type DarkbotEnv struct {
	utils_settings.UtilsEnvs

	ScrappyBaseUrl       string
	ScrappyPlayerUrl     string
	ScrappyBaseAttackUrl string

	DiscorderBotToken string

	ConfiguratorDbname string

	ConsolerPrefix   string
	ProfilingEnabled bool

	ScrappyLoopDelay int
	ViewerLoopDelay  int
	DevEnvMockApi    bool
}

var Env DarkbotEnv

func init() {
	logus.Log.Info("attempt to load settings")

	envs := utils_env.NewEnvConfig()
	Env = DarkbotEnv{
		UtilsEnvs:            utils_settings.Envs,
		ScrappyBaseUrl:       envs.GetEnvWithDefault("SCRAPPY_BASE_URL", "undefined"),
		ScrappyPlayerUrl:     envs.GetEnvWithDefault("SCRAPPY_PLAYER_URL", "undefined"),
		ScrappyBaseAttackUrl: envs.GetEnvWithDefault("SCRAPPY_BASE_ATTACK_URL", "https://discoverygc.com/forums/showthread.php?tid=110046&action=lastpost"),

		DiscorderBotToken: envs.GetEnvWithDefault("DISCORDER_BOT_TOKEN", "undefined"),

		ConfiguratorDbname: envs.GetEnvWithDefault("CONFIGURATOR_DBNAME", "dev"),

		ConsolerPrefix:   envs.GetEnvWithDefault("CONSOLER_PREFIX", ";"),
		ProfilingEnabled: envs.GetEnvBool("PROFILING"),

		ScrappyLoopDelay: envs.GetIntWithDefault("SCRAPPY_LOOP_DELAY", 10),
		ViewerLoopDelay:  envs.GetIntWithDefault("VIEWER_LOOP_DELAY", 10),
		DevEnvMockApi:    envs.GetEnvBoolWithDefault("DEVENV_MOCK_API", true),
	}
	Workdir = filepath.Dir(filepath.Dir(utils.GetCurrentFolder().ToString()))
	Dbpath = NewDBPath(Env.ConfiguratorDbname)
}

var Dbpath types.Dbpath
var Workdir string

func NewDBPath(dbname string) types.Dbpath {
	return types.Dbpath(filepath.Join(Workdir, "data", dbname+".sqlite3"))
}
