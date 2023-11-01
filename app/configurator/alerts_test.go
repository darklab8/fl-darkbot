package configurator

import (
	"darkbot/app/settings/types"
	"darkbot/app/settings/utils"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAlertTreshold(t *testing.T) {
	FixtureMigrator(func(dbpath types.Dbpath) {
		channelID, _ := FixtureChannel(dbpath)
		genericCfg := NewConfigurator(dbpath).AutoMigrateSchema()
		_ = channelID

		cfg := NewCfgAlertNeutralPlayersGreaterThan(genericCfg)
		status, _ := cfg.Status(channelID)
		fmt.Println("status=", status)
		assert.Nil(t, status, "status is not Nil. failed aert")

		cfg.Set(channelID, 5)
		status, _ = cfg.Status(channelID)
		assert.Equal(t, utils.Ptr(5), status)

		cfg.Unset(channelID)
		status, _ = cfg.Status(channelID)
		assert.Nil(t, status)
	})
}

func TestAlertBool(t *testing.T) {
	FixtureMigrator(func(dbpath types.Dbpath) {
		channelID, _ := FixtureChannel(dbpath)
		genericCfg := NewConfigurator(dbpath).AutoMigrateSchema()
		_ = channelID

		cfg := NewCfgAlertBaseIsUnderAttack(genericCfg)
		status, _ := cfg.Status(channelID)
		fmt.Println("status=", status)
		assert.False(t, status, "status is not true. failed aert")

		cfg.Enable(channelID)
		status, _ = cfg.Status(channelID)
		assert.True(t, status)

		cfg.Disable(channelID)
		status, _ = cfg.Status(channelID)
		assert.False(t, status)
	})
}
