package settings

import (
	"path/filepath"

	"github.com/darklab8/fl-darkbot/app/settings/logus"
	"github.com/darklab8/fl-darkbot/app/settings/types"

	"github.com/darklab8/go-utils/utils/enverant"
	"github.com/darklab8/go-utils/utils/utils_os"
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

	envs := enverant.NewEnverant(
		enverant.WithEnvFile(utils_os.GetCurrentFolder().Dir().Dir().Join(".vscode", "enverant.json").ToString()),
	)
	Env = DarkbotEnv{
		UtilsEnvs:            utils_settings.GetEnvs(envs),
		DevEnvMockApi:        envs.GetBoolOr("DEV_ENV_MOCK_API", true),
		ScrappyBaseUrl:       envs.GetStrOr("SCRAPPY_BASE_URL", ""),
		ScrappyPlayerUrl:     envs.GetStrOr("SCRAPPY_PLAYER_URL", ""),
		ScrappyBaseAttackUrl: envs.GetStrOr("SCRAPPY_BASE_ATTACK_URL", "https://discoverygc.com/forums/showthread.php?tid=110046&action=lastpost"),

		DiscorderBotToken: envs.GetStrOr("DISCORDER_BOT_TOKEN", ""),

		ConfiguratorDbname: envs.GetStrOr("CONFIGURATOR_DBNAME", "dev"),

		ConsolerPrefix:   envs.GetStrOr("CONSOLER_PREFIX", ";"),
		ProfilingEnabled: envs.GetBoolOr("PROFILING", false),

		ScrappyLoopDelay: envs.GetIntOr("SCRAPPY_LOOP_DELAY", 10),
		ViewerLoopDelay:  envs.GetIntOr("VIEWER_LOOP_DELAY", 10),
	}
	Workdir = filepath.Dir(filepath.Dir(utils_os.GetCurrentFolder().ToString()))
	Dbpath = NewDBPath(Env.ConfiguratorDbname)

	if !Env.DevEnvMockApi {
		if Env.ScrappyBaseUrl == "" {
			logus.Log.Panic("DevEnvMockApi=false, Expected SCRAPPY_BASE_URL env var to be defined")
		}
		if Env.ScrappyPlayerUrl == "" {
			logus.Log.Panic("DevEnvMockApi=false, Expected SCRAPPY_PLAYER_URL env var to be defined")
		}
	}
}

var Dbpath types.Dbpath
var Workdir string

func NewDBPath(dbname string) types.Dbpath {
	return types.Dbpath(filepath.Join(Workdir, "data", dbname+".sqlite3"))
}
