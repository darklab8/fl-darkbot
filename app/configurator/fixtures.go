package configurator

import (
	"darkbot/app/settings"
	"darkbot/app/settings/logus"
	"darkbot/app/settings/types"
	"darkbot/app/settings/utils"
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
	logus.Debug("", logus.Dbpath(dbpath))
	os.Remove(string(dbpath))
	os.Remove(string(dbpath) + "-shm")
	os.Remove(string(dbpath) + "-wal")
	cfg := FixtureConfigurator(dbpath)
	cfg.AutoMigrateSchema()

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
	cfg_channel := NewConfiguratorChannel(configurator_)
	cfg_channel.Add(channelID)

	return channelID, cfg_channel
}