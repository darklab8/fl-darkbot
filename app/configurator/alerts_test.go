package configurator

import (
	"fmt"
	"testing"

	"github.com/darklab8/fl-darkbot/app/configurator/models"
	"github.com/darklab8/fl-darkbot/app/settings/logus"
	"github.com/darklab8/fl-darkbot/app/settings/types"

	"github.com/stretchr/testify/assert"
)

func TestAlertTreshold(t *testing.T) {
	FixtureMigrator(func(dbpath types.Dbpath) {
		channelID, _ := FixtureChannel(dbpath)
		genericCfg := NewConfigurator(dbpath).AutoMigrateSchema()
		_ = channelID

		cfg := NewCfgAlertBaseHealthLowerThan(genericCfg)
		status, err := cfg.Status(channelID)
		logus.Log.Debug(fmt.Sprintf("status=%v", status))
		assert.Error(t, err, ErrorZeroAffectedRowsMsg)
		assert.ErrorContains(t, err, ErrorZeroAffectedRowsMsg)

		cfg.Set(channelID, models.ThresholdIntegerPercentage, 5)
		status, err = cfg.Status(channelID)
		assert.Nil(t, err, "result of status operation is without errors")
		assert.Equal(t, 5, status)

		cfg.Unset(channelID)
		_, err = cfg.Status(channelID)
		assert.Error(t, err)
		assert.ErrorContains(t, err, ErrorZeroAffectedRowsMsg)
	})
}

func TestAlertBool(t *testing.T) {
	FixtureMigrator(func(dbpath types.Dbpath) {
		channelID, _ := FixtureChannel(dbpath)
		genericCfg := NewConfigurator(dbpath).AutoMigrateSchema()
		_ = channelID

		cfg := NewCfgAlertBaseIsUnderAttack(genericCfg)
		status, _ := cfg.Status(channelID)
		logus.Log.Debug(fmt.Sprintf("status=%t", status))
		assert.False(t, status, "status is not true. failed aert")

		cfg.Enable(channelID)
		status, _ = cfg.Status(channelID)
		assert.True(t, status)

		cfg.Disable(channelID)
		status, _ = cfg.Status(channelID)
		assert.False(t, status)
	})
}
