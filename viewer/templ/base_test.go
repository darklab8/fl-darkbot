package templ

import (
	"darkbot/configurator"
	"darkbot/dtypes"
	"darkbot/scrappy"
	"darkbot/scrappy/base"
	"darkbot/scrappy/baseattack"
	"darkbot/scrappy/player"
	"darkbot/scrappy/shared/records"
	"darkbot/utils/logger"
	"darkbot/viewer/apis"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBaseViewerMocked(t *testing.T) {
	configurator.FixtureMigrator(func(dbpath dtypes.Dbpath) {
		channelID, _ := configurator.FixtureChannel(dbpath)

		cg := configurator.ConfiguratorBase{Configurator: configurator.NewConfigurator(dbpath)}
		cg.TagsAdd(channelID, []string{"Station"}...)

		bases := (&base.BaseStorage{}).New()
		baseAttack := (&baseattack.BaseAttackStorage{}).New()
		scrappy.Storage = &scrappy.ScrappyStorage{BaseStorage: bases, BaseAttackStorage: baseAttack}
		record := records.StampedObjects[base.Base]{}.New()
		record.Add(base.Base{Name: "Station1", Affiliation: "Abc", Health: 100})
		record.Add(base.Base{Name: "Station2", Affiliation: "Qwe", Health: 100})
		bases.Add(record)

		render := NewTemplateBase(channelID, dbpath)
		render.Render()
		fmt.Println(render.main.Content)

		assert.NotEmpty(t, render.main.Content)
		assert.Empty(t, render.AlertHealthLowerThan.Content)
		assert.Empty(t, render.AlertHealthIsDecreasing.Content)
		assert.Empty(t, render.AlertBaseUnderAttack.Content)

		// alerts
		baseAlertDecreasing := configurator.CfgAlertBaseHealthIsDecreasing{Configurator: configurator.NewConfigurator(dbpath)}
		isEnabled, _ := baseAlertDecreasing.Status(channelID)
		assert.False(t, isEnabled)

		record = records.StampedObjects[base.Base]{}.New()
		record.Add(base.Base{Name: "Station1", Affiliation: "Abc", Health: 100})
		record.Add(base.Base{Name: "Station2", Affiliation: "Qwe", Health: 50})
		bases.Add(record)

		baseAlertDecreasing.Enable(channelID)
		isEnabled, _ = baseAlertDecreasing.Status(channelID)
		assert.True(t, isEnabled)

		render = NewTemplateBase(channelID, dbpath)
		render.Render()

		assert.NotEmpty(t, render.main.Content)
		assert.NotEmpty(t, render.AlertHealthIsDecreasing.Content)
		assert.Empty(t, render.AlertHealthLowerThan.Content)
		assert.Empty(t, render.AlertBaseUnderAttack.Content)

		baseAlertBelowThreshold := configurator.CfgAlertBaseHealthLowerThan{Configurator: configurator.NewConfigurator(dbpath)}
		threshold, _ := baseAlertBelowThreshold.Status(channelID)
		assert.Nil(t, threshold)

		baseAlertBelowThreshold.Set(channelID, 40)
		render = NewTemplateBase(channelID, dbpath)
		render.Render()

		assert.NotEmpty(t, render.main.Content)
		assert.NotEmpty(t, render.AlertHealthIsDecreasing.Content)
		assert.Empty(t, render.AlertHealthLowerThan.Content)
		assert.Empty(t, render.AlertBaseUnderAttack.Content)

		baseAlertBelowThreshold.Set(channelID, 60)
		render = NewTemplateBase(channelID, dbpath)
		render.Render()

		assert.NotEmpty(t, render.main.Content)
		assert.NotEmpty(t, render.AlertHealthIsDecreasing.Content)
		assert.NotEmpty(t, render.AlertHealthLowerThan.Content)
		assert.Empty(t, render.AlertBaseUnderAttack.Content)
	})
}

// TODO fix those tests... for some reason memory ref error :smile:
// func TestBaseViewerRealData(t *testing.T) {
// 	configurator.FixtureMigrator(func(dbpath dtypes.Dbpath) {
// 		channelID, _ := configurator.FixtureChannel(dbpath)

// 		cg := configurator.ConfiguratorBase{Configurator: configurator.NewConfigurator(dbpath)}
// 		cg.TagsAdd(channelID, []string{"Station"}...)

// 		scrappy.Storage.BaseStorage.Api = base.NewMock("basedata.json")
// 		scrappy.Storage.PlayerStorage.Api = player.APIPlayerSpy{}
// 		scrappy.Storage.Update()
// 		scrappy.Storage.BaseStorage.Api = base.NewMock("basedata2.json")
// 		scrappy.Storage.Update()
// 		scrappy.Storage.BaseStorage.Records.List(func(values []records.StampedObjects[base.Base]) {
// 			values[1].Timestamp = values[0].Timestamp.Add(time.Minute * 15)
// 		})

// 		base := NewTemplateBase(channelID, dbpath)
// 		base.Render()
// 		fmt.Println(base.main.Content)
// 	})
// }

// TEST TO FIND OUT derivative of base health
func TestGetDerivative(t *testing.T) {
	configurator.FixtureMigrator(func(dbpath dtypes.Dbpath) {
		logger.Debug("1")
		channelID, _ := configurator.FixtureChannel(dbpath)
		api := apis.NewAPI(channelID, dbpath)

		tags := []string{""}
		logger.Debug("2")
		scrappy.Storage.BaseStorage = (&base.BaseStorage{}).New()
		scrappy.Storage.PlayerStorage = (&player.PlayerStorage{}).New()
		scrappy.Storage.BaseStorage.Api = base.NewMock("basedata.json")
		logger.Debug("2.1")
		scrappy.Storage.PlayerStorage.Api = player.APIPlayerSpy{}.New()
		logger.Debug("2.2")
		scrappy.Storage.Update()
		logger.Debug("2.3")
		scrappy.Storage.BaseStorage.Api = base.NewMock("basedata2.json")
		scrappy.Storage.Update()
		scrappy.Storage.BaseStorage.Records.List(func(values []records.StampedObjects[base.Base]) {
			values[1].Timestamp = values[0].Timestamp.Add(time.Minute * 15)
		})

		logger.Debug("3")

		result1 := make(map[string]base.Base)
		result2 := make(map[string]base.Base)
		scrappy.Storage.BaseStorage.Records.List(func(values []records.StampedObjects[base.Base]) {
			for _, base := range values[0].List {
				result1[base.Name] = base
			}

			for _, base := range values[1].List {
				result2[base.Name] = base
			}
		})
		logger.Debug("4")
		res1 := result1["Stockholm Base"]
		res2 := result2["Stockholm Base"]
		_ = res1
		_ = res2

		logger.Debug("5")
		baseDerivatives := CalculateDerivates(tags, api)
		for baseName, baseDeravative := range baseDerivatives {
			logger.Info("baseName=", baseName, " baseDeravative=", baseDeravative)
		}
		logger.Debug("6")
	})
}
