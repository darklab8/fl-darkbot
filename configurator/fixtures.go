package configurator

import (
	"darkbot/settings"
	"darkbot/settings/types"
	"darkbot/settings/utils"
	"fmt"
	"os"
)

func FixtureConfigurator(dbpath types.Dbpath) Configurator {
	cfg := NewConfigurator(dbpath)
	return cfg
}

func FixtureMigrator(callback func(dbpath types.Dbpath)) Configurator {
	dbname := utils.TokenHex(8)
	dbpath := types.Dbpath(settings.NewDBPath(dbname))
	// setup
	fmt.Println(dbpath)
	os.Remove(string(dbpath))
	os.Remove(string(dbpath) + "-shm")
	os.Remove(string(dbpath) + "-wal")
	cfg := FixtureConfigurator(dbpath)
	cfg.Migrate()

	// teardown
	defer os.Remove(string(dbpath))
	defer os.Remove(string(dbpath) + "-shm")
	defer os.Remove(string(dbpath) + "-wal")

	callback(dbpath)

	return cfg
}

func FixtureChannel(dbpath types.Dbpath) (types.DiscordChannelID, ConfiguratorChannel) {
	channelID := types.DiscordChannelID("123")
	configurator_ := FixtureConfigurator(dbpath)
	cfg_channel := ConfiguratorChannel{Configurator: configurator_}
	cfg_channel.Add(channelID)

	return channelID, cfg_channel
}
