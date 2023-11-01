package templ

import (
	"darkbot/configurator"
	"darkbot/scrappy"
	"darkbot/scrappy/base"
	"darkbot/scrappy/baseattack"
	"darkbot/scrappy/player"
	"darkbot/scrappy/shared/records"
	"darkbot/settings/logus"
	"darkbot/settings/types"
	"darkbot/viewer/apis"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBaseViewerMocked(t *testing.T) {
	configurator.FixtureMigrator(func(dbpath types.Dbpath) {
		channelID, _ := configurator.FixtureChannel(dbpath)

		cg := configurator.NewConfiguratorBase(configurator.NewConfigurator(dbpath))
		cg.TagsAdd(channelID, []string{"Station"}...)

		scrappy.Storage = scrappy.FixtureMockedStorage()
		scrappy.Storage.Update()
		record := records.NewStampedObjects[base.Base]()
		record.Add(base.Base{Name: "Station1", Affiliation: "Abc", Health: 100})
		record.Add(base.Base{Name: "Station2", Affiliation: "Qwe", Health: 100})
		scrappy.Storage.GetBaseStorage().Add(record)

		render := NewTemplateBase(channelID, dbpath)
		render.Render()
		fmt.Println(render.main.Content)

		assert.NotEmpty(t, render.main.Content)
		assert.Empty(t, render.AlertHealthLowerThan.Content)
		assert.Empty(t, render.AlertHealthIsDecreasing.Content)
		assert.Empty(t, render.AlertBaseUnderAttack.Content)

		// alerts
		baseAlertDecreasing := configurator.NewCfgAlertBaseHealthIsDecreasing(configurator.NewConfigurator(dbpath))
		isEnabled, _ := baseAlertDecreasing.Status(channelID)
		assert.False(t, isEnabled)

		record = records.NewStampedObjects[base.Base]()
		record.Add(base.Base{Name: "Station1", Affiliation: "Abc", Health: 100})
		record.Add(base.Base{Name: "Station2", Affiliation: "Qwe", Health: 50})
		scrappy.Storage.GetBaseStorage().Add(record)

		baseAlertDecreasing.Enable(channelID)
		isEnabled, _ = baseAlertDecreasing.Status(channelID)
		assert.True(t, isEnabled)

		render = NewTemplateBase(channelID, dbpath)
		render.Render()

		assert.NotEmpty(t, render.main.Content)
		assert.NotEmpty(t, render.AlertHealthIsDecreasing.Content)
		assert.Empty(t, render.AlertHealthLowerThan.Content)
		assert.Empty(t, render.AlertBaseUnderAttack.Content)

		baseAlertBelowThreshold := configurator.NewCfgAlertBaseHealthLowerThan(configurator.NewConfigurator(dbpath))
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

		record = records.NewStampedObjects[base.Base]()
		record.Add(base.Base{Name: "Bank of Bretonia", Affiliation: "Abc", Health: 100})
		scrappy.Storage.GetBaseStorage().Add(record)
		cg.TagsAdd(channelID, []string{"Bank"}...)
		render = NewTemplateBase(channelID, dbpath)
		render.Render()

		assert.Empty(t, render.AlertBaseUnderAttack.Content)

		baseUnderAttackalert := configurator.NewCfgAlertBaseIsUnderAttack(configurator.NewConfigurator(dbpath))
		baseUnderAttackalert.Enable(channelID)
		render = NewTemplateBase(channelID, dbpath)
		render.Render()

		assert.NotEmpty(t, render.AlertBaseUnderAttack.Content)
	})
}

// TODO fix those tests... for some reason memory ref error :smile:
// func TestBaseViewerRealData(t *testing.T) {
// 	configurator.FixtureMigrator(func(dbpath types.Dbpath) {
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
	configurator.FixtureMigrator(func(dbpath types.Dbpath) {
		logus.Debug("1")
		channelID, _ := configurator.FixtureChannel(dbpath)
		api := apis.NewAPI(channelID, dbpath)

		tags := []string{""}
		logus.Debug("2")
		scrappy.Storage = scrappy.NewScrapyStorage(base.NewMock("basedata.json"), player.FixturePlayerAPIMock(), baseattack.FixtureBaseAttackAPIMock())
		logus.Debug("2.1")
		logus.Debug("2.2")
		scrappy.Storage.Update()
		logus.Debug("2.3")

		scrappy.FixtureSetBaseStorageAPI(base.NewMock("basedata2.json"))
		scrappy.Storage.Update()
		scrappy.Storage.GetBaseStorage().Records.List(func(values []records.StampedObjects[base.Base]) {
			values[1].Timestamp = values[0].Timestamp.Add(time.Minute * 15)
		})

		logus.Debug("3")

		result1 := make(map[string]base.Base)
		result2 := make(map[string]base.Base)
		scrappy.Storage.GetBaseStorage().Records.List(func(values []records.StampedObjects[base.Base]) {
			for _, base := range values[0].List {
				result1[base.Name] = base
			}

			for _, base := range values[1].List {
				result2[base.Name] = base
			}
		})
		logus.Debug("4")
		res1 := result1["Stockholm Base"]
		res2 := result2["Stockholm Base"]
		_ = res1
		_ = res2

		logus.Debug("5")
		baseDerivatives, _ := CalculateDerivates(tags, api)
		for baseName, baseDeravative := range baseDerivatives {
			logus.Info(fmt.Sprintf("baseName=%s, baseDeravative=%f", baseName, baseDeravative))
		}
		logus.Debug("6")
	})
}

func TestDetectAttackOnLPBase(t *testing.T) {

	configurator.FixtureMigrator(func(dbpath types.Dbpath) {
		channelID, _ := configurator.FixtureChannel(dbpath)

		cg := configurator.NewConfiguratorBase(configurator.NewConfigurator(dbpath))
		cg.TagsAdd(channelID, []string{"LP-7743"}...)

		scrappy.Storage = scrappy.NewScrapyStorage(base.FixtureBaseApiMock(), player.FixturePlayerAPIMock(), baseattack.NewMock("data_lp.json"))
		scrappy.Storage.Update()

		assert.True(t, strings.Contains(string(scrappy.Storage.GetBaseAttackStorage().GetData()), "LP-7743"))

		bases := scrappy.Storage.GetBaseStorage()
		record := records.NewStampedObjects[base.Base]()
		record.Add(base.Base{Name: "LP-7743", Affiliation: "Abc", Health: 5})
		bases.Add(record)
		record2 := records.NewStampedObjects[base.Base]()
		record2.Add(base.Base{Name: "LP-7743", Affiliation: "Abc", Health: 6})
		record2.Timestamp = record2.Timestamp.Add(time.Hour * 1)
		bases.Add(record2)

		baseUnderAttackalert := configurator.NewCfgAlertBaseIsUnderAttack(configurator.NewConfigurator(dbpath))
		baseUnderAttackalert.Enable(channelID)

		render := NewTemplateBase(channelID, dbpath)
		render.Render()

		assert.NotEmpty(t, render.AlertBaseUnderAttack.Content)
	})
}
