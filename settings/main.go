package settings

import (
	"darkbot/utils"
	"path/filepath"

	"github.com/spf13/viper"
)

type ConfigScheme struct {
	ScrappyBaseUrl   string `mapstructure:"SCRAPPY_BASE_URL"`
	ScrappyPlayerUrl string `mapstructure:"SCRAPPY_PLAYER_URL"`

	ScrappyForumUsername string `mapstructure:"SCRAPPY_FORUM_USERNAME"`
	ScrappyForumPassword string `mapstructure:"SCRAPPY_FORUM_PASSWORD"`

	DiscorderBotToken string `mapstructure:"DISCORDER_BOT_TOKEN"`

	ConfiguratorDbname string `mapstructure:"CONFIGURATOR_DBNAME"`

	ConsolerPrefix string `mapstructure:"CONSOLER_PREFIX"`
}

var Config ConfigScheme

var Dbpath string

func load() {
	utils.LogInfo("identifying folder of settings")
	workdir := filepath.Dir(utils.GetCurrrentFolder())
	viper.AddConfigPath(workdir)
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AutomaticEnv()
	viper.ReadInConfig()

	viper.Unmarshal(&Config)

	utils.LogInfo("settings were downloaded. Scrappy base url=", Config.ScrappyBaseUrl)

	Dbpath = filepath.Join(workdir, "data", Config.ConfiguratorDbname+".sqlite3")

	if Config.ConsolerPrefix == "" {
		Config.ConsolerPrefix = "."
	}
}

func init() {
	utils.LogInfo("attempt to load settings")
	load()
}
