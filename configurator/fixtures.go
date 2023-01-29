package configurator

import (
	"darkbot/dtypes"
	"darkbot/settings"
	"darkbot/utils"
	"fmt"
	"os"
)

func FixtureConfigurator(dbpath dtypes.Dbpath) Configurator {
	cfg := NewConfigurator(dbpath)
	return cfg
}

func FixtureMigrator(callback func(dbpath dtypes.Dbpath)) Configurator {
	dbname := utils.TokenHex(8)
	dbpath := dtypes.Dbpath(settings.NewDBPath(dbname))
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

func FixtureChannel(dbpath dtypes.Dbpath) (string, ConfiguratorChannel) {
	channelID := "123"
	configurator_ := FixtureConfigurator(dbpath)
	cfg_channel := ConfiguratorChannel{Configurator: configurator_}
	cfg_channel.Add(channelID)

	return channelID, cfg_channel
}
