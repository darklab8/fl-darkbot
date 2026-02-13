package settings

import (
	"log"
	"path/filepath"

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

	DarkstatHost string
	DarkstatPort int

	PlayerViewDeprecated bool

	PrometheuserOn bool
}

var Env DarkbotEnv

var Environ *enverant.Enverant

func init() {
	log.Println("attempt to load settings")

	Environ = enverant.NewEnverant()
	LoadEnv(Environ)
}

func LoadEnv(envs *enverant.Enverant) {
	Env = DarkbotEnv{
		UtilsEnvs:            utils_settings.GetEnvs(),
		DevEnvMockApi:        envs.GetBoolOr("DEV_ENV_MOCK_API", true),
		ScrappyBaseUrl:       envs.GetStrOr("SCRAPPY_BASE_URL", ""),
		ScrappyPlayerUrl:     envs.GetStrOr("SCRAPPY_PLAYER_URL", ""),
		ScrappyBaseAttackUrl: envs.GetStrOr("SCRAPPY_BASE_ATTACK_URL", "https://discoverygc.com/forums/showthread.php?tid=110046&action=lastpost"),

		DiscorderBotToken: envs.GetStr("DISCORDER_BOT_TOKEN"),

		ConfiguratorDbname: envs.GetStrOr("CONFIGURATOR_DBNAME", "dev"),

		ConsolerPrefix:   envs.GetStrOr("CONSOLER_PREFIX", ";"),
		ProfilingEnabled: envs.GetBoolOr("PROFILING", false),

		ScrappyLoopDelay: envs.GetIntOr("SCRAPPY_LOOP_DELAY", 10),
		ViewerLoopDelay:  envs.GetIntOr("VIEWER_LOOP_DELAY", 10),

		DarkstatHost: envs.GetStrOr("DARKBOT_DARKSTAT_HOST", "127.0.0.1"),
		DarkstatPort: envs.GetIntOr("DARKBOT_DARKSTAT_PORT", 8100),

		PrometheuserOn:       envs.GetBoolOr("PROMETHEUSER_ON", true),
		PlayerViewDeprecated: true,
	}
	Workdir = filepath.Dir(filepath.Dir(utils_os.GetCurrentFolder().ToString()))
	Dbpath = NewDBPath(Env.ConfiguratorDbname)

	if !Env.DevEnvMockApi {
		if Env.ScrappyBaseUrl == "" {
			log.Panic("DevEnvMockApi=false, Expected SCRAPPY_BASE_URL env var to be defined")
		}
		if Env.ScrappyPlayerUrl == "" {
			log.Panic("DevEnvMockApi=false, Expected SCRAPPY_PLAYER_URL env var to be defined")
		}
	}
}

var Dbpath types.Dbpath
var Workdir string

func NewDBPath(dbname string) types.Dbpath {
	return types.Dbpath(filepath.Join(Workdir, "data", dbname+".sqlite3"))
}
