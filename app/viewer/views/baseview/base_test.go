package baseview

import (
	"darkbot/app/configurator"
	"darkbot/app/scrappy"
	"darkbot/app/scrappy/base"
	"darkbot/app/scrappy/baseattack"
	"darkbot/app/scrappy/player"
	"darkbot/app/scrappy/shared/records"
	"darkbot/app/settings/logus"
	"darkbot/app/settings/types"
	"darkbot/app/viewer/apis"
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
		cg.TagsAdd(channelID, []types.Tag{"Station"}...)

		scrapper := scrappy.FixtureMockedStorage()
		scrapper.Update()
		record := records.NewStampedObjects[base.Base]()
		record.Add(base.Base{Name: "Station1", Affiliation: "Abc", Health: 100})
		record.Add(base.Base{Name: "Station2", Affiliation: "Qwe", Health: 100})
		scrapper.GetBaseStorage().Add(record)

		api := apis.NewAPI(dbpath, scrapper)

		render := NewTemplateBase(api, channelID)
		render.RenderView()
		logus.Log.Debug(fmt.Sprintf("render.main.Content=%v", render.main.GetMsgs()))

		assert.True(t, render.main.HasRecords())
		assert.False(t, render.alertHealthLowerThan.HasRecords())
		assert.False(t, render.alertHealthIsDecreasing.HasRecords())
		assert.False(t, render.alertBaseUnderAttack.HasRecords())

		// alerts
		baseAlertDecreasing := configurator.NewCfgAlertBaseHealthIsDecreasing(configurator.NewConfigurator(dbpath))
		isEnabled, _ := baseAlertDecreasing.Status(channelID)
		assert.False(t, isEnabled)

		record = records.NewStampedObjects[base.Base]()
		record.Add(base.Base{Name: "Station1", Affiliation: "Abc", Health: 100})
		record.Add(base.Base{Name: "Station2", Affiliation: "Qwe", Health: 50})
		scrapper.GetBaseStorage().Add(record)

		baseAlertDecreasing.Enable(channelID)
		isEnabled, _ = baseAlertDecreasing.Status(channelID)
		assert.True(t, isEnabled)

		render = NewTemplateBase(api, channelID)
		render.RenderView()

		assert.True(t, render.main.HasRecords())
		assert.True(t, render.alertHealthIsDecreasing.HasRecords())
		assert.False(t, render.alertHealthLowerThan.HasRecords())
		assert.False(t, render.alertBaseUnderAttack.HasRecords())

		baseAlertBelowThreshold := configurator.NewCfgAlertBaseHealthLowerThan(configurator.NewConfigurator(dbpath))
		_, err := baseAlertBelowThreshold.Status(channelID)
		assert.Error(t, err)

		baseAlertBelowThreshold.Set(channelID, 40)
		render = NewTemplateBase(api, channelID)
		render.RenderView()

		assert.True(t, render.main.HasRecords())
		assert.True(t, render.alertHealthIsDecreasing.HasRecords())
		assert.False(t, render.alertHealthLowerThan.HasRecords())
		assert.False(t, render.alertBaseUnderAttack.HasRecords())

		baseAlertBelowThreshold.Set(channelID, 60)
		render = NewTemplateBase(api, channelID)
		render.RenderView()

		assert.True(t, render.main.HasRecords())
		assert.True(t, render.alertHealthIsDecreasing.HasRecords())
		assert.True(t, render.alertHealthLowerThan.HasRecords())
		assert.False(t, render.alertBaseUnderAttack.HasRecords())

		record = records.NewStampedObjects[base.Base]()
		record.Add(base.Base{Name: "Bank of Bretonia", Affiliation: "Abc", Health: 100})
		scrapper.GetBaseStorage().Add(record)
		cg.TagsAdd(channelID, []types.Tag{"Bank"}...)
		render = NewTemplateBase(api, channelID)
		render.RenderView()

		assert.False(t, render.alertBaseUnderAttack.HasRecords())

		baseUnderAttackalert := configurator.NewCfgAlertBaseIsUnderAttack(configurator.NewConfigurator(dbpath))
		baseUnderAttackalert.Enable(channelID)
		render = NewTemplateBase(api, channelID)
		render.RenderView()

		assert.True(t, render.alertBaseUnderAttack.HasRecords())
	})
}

func TestBaseViewerRealData(t *testing.T) {
	configurator.FixtureMigrator(func(dbpath types.Dbpath) {
		channelID, _ := configurator.FixtureChannel(dbpath)

		cg := configurator.NewConfiguratorBase(configurator.NewConfigurator(dbpath))
		cg.TagsAdd(channelID, []types.Tag{"Station"}...)

		scrapper := scrappy.NewScrapyStorage(
			base.NewMock("basedata.json"),
			player.FixturePlayerAPIMock(),
			baseattack.FixtureBaseAttackAPIMock(),
		)
		api := apis.NewAPI(dbpath, scrapper)
		scrapper.Update()
		scrapper.GetBaseStorage().FixtureSetAPI(base.NewMock("basedata2.json"))
		scrapper.Update()
		scrapper.GetBaseStorage().Records.List(func(values []records.StampedObjects[base.Base]) {
			values[1].Timestamp = values[0].Timestamp.Add(time.Minute * 15)
		})

		base := NewTemplateBase(api, channelID)
		base.RenderView()
		logus.Log.Debug(fmt.Sprintf("base.main.Msgs=%v", base.main.GetMsgs()))
	})
}

func TestGetDerivativeBaseHealth(t *testing.T) {
	configurator.FixtureMigrator(func(dbpath types.Dbpath) {
		logus.Log.Debug("1")
		tags := []types.Tag{""}
		logus.Log.Debug("2")
		scrapper := scrappy.NewScrapyStorage(base.NewMock("basedata.json"), player.FixturePlayerAPIMock(), baseattack.FixtureBaseAttackAPIMock())
		logus.Log.Debug("2.1")
		logus.Log.Debug("2.2")
		scrapper.Update()
		logus.Log.Debug("2.3")

		api := apis.NewAPI(dbpath, scrapper)

		scrapper.GetBaseStorage().FixtureSetAPI(base.NewMock("basedata2.json"))
		scrapper.Update()
		scrapper.GetBaseStorage().Records.List(func(values []records.StampedObjects[base.Base]) {
			values[1].Timestamp = values[0].Timestamp.Add(time.Minute * 15)
		})

		logus.Log.Debug("3")

		result1 := make(map[string]base.Base)
		result2 := make(map[string]base.Base)
		scrapper.GetBaseStorage().Records.List(func(values []records.StampedObjects[base.Base]) {
			for _, base := range values[0].List {
				result1[base.Name] = base
			}

			for _, base := range values[1].List {
				result2[base.Name] = base
			}
		})
		logus.Log.Debug("4")
		res1 := result1["Stockholm Base"]
		res2 := result2["Stockholm Base"]
		_ = res1
		_ = res2

		logus.Log.Debug("5")
		baseDerivatives, _ := CalculateDerivates(tags, api)
		for baseName, baseDeravative := range baseDerivatives {
			logus.Log.Info(fmt.Sprintf("baseName=%s, baseDeravative=%f", baseName, baseDeravative))
		}
		logus.Log.Debug("6")
	})
}

func TestDetectAttackOnLPBase(t *testing.T) {
	configurator.FixtureMigrator(func(dbpath types.Dbpath) {
		channelID, _ := configurator.FixtureChannel(dbpath)

		cg := configurator.NewConfiguratorBase(configurator.NewConfigurator(dbpath))
		cg.TagsAdd(channelID, []types.Tag{"LP-7743"}...)

		scrapper := scrappy.NewScrapyStorage(base.FixtureBaseApiMock(), player.FixturePlayerAPIMock(), baseattack.NewMock("data_lp.json"))
		scrapper.Update()
		api := apis.NewAPI(dbpath, scrapper)

		assert.True(t, strings.Contains(string(scrapper.GetBaseAttackStorage().GetData()), "LP-7743"))

		bases := scrapper.GetBaseStorage()
		record := records.NewStampedObjects[base.Base]()
		record.Add(base.Base{Name: "LP-7743", Affiliation: "Abc", Health: 5})
		bases.Add(record)
		record2 := records.NewStampedObjects[base.Base]()
		record2.Add(base.Base{Name: "LP-7743", Affiliation: "Abc", Health: 6})
		record2.Timestamp = record2.Timestamp.Add(time.Hour * 1)
		bases.Add(record2)

		baseUnderAttackalert := configurator.NewCfgAlertBaseIsUnderAttack(configurator.NewConfigurator(dbpath))
		baseUnderAttackalert.Enable(channelID)

		render := NewTemplateBase(api, channelID)
		render.RenderView()

		assert.True(t, render.alertBaseUnderAttack.HasRecords())
	})
}
