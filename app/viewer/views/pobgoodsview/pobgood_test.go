package pobgoodsview

import (
	"fmt"
	"testing"

	"github.com/darklab8/fl-darkbot/app/configurator"
	"github.com/darklab8/fl-darkbot/app/scrappy"
	"github.com/darklab8/fl-darkbot/app/scrappy/base"
	"github.com/darklab8/fl-darkbot/app/scrappy/baseattack"
	"github.com/darklab8/fl-darkbot/app/scrappy/player"
	"github.com/darklab8/fl-darkbot/app/settings/logus"
	"github.com/darklab8/fl-darkbot/app/settings/types"
	"github.com/darklab8/fl-darkbot/app/viewer/apis"
	"github.com/darklab8/fl-darkbot/app/viewer/views/baseview"
	"github.com/darklab8/fl-darkstat/darkstat/configs_export"

	"github.com/stretchr/testify/assert"
)

func TestPoBGoodsRender(t *testing.T) {
	configurator.FixtureMigrator(func(dbpath types.Dbpath) {
		channelID, _ := configurator.FixtureChannel(dbpath)

		cg := configurator.NewConfiguratorBase(configurator.NewConfigurator(dbpath))
		cg.TagsAdd(channelID, []types.Tag{"Station"}...)

		scrapper := scrappy.FixtureMockedStorage()
		scrapper.Update()

		api := apis.NewAPI(dbpath, scrapper)

		tags, _ := api.Bases.Tags.TagsList(channelID)
		record, err := api.Scrappy.GetBaseStorage().GetLatestRecord()
		logus.Log.CheckPanic(err, "unable to query bases from storage in Template base Generate records")
		matchedBases := baseview.MatchBases(record.List, tags)

		goods := make(map[string]*configs_export.ShopItem)
		for _, base := range matchedBases {
			for _, good := range base.ShopItems {
				goods[good.Nickname] = good
			}
		}

		i := 0
		cg_pobgood := configurator.NewConfiguratorPoBGood(configurator.NewConfigurator(dbpath))
		for _, good := range goods {
			cg_pobgood.TagsAdd(channelID, []types.Tag{types.Tag(good.Nickname)}...)
			i++
			if i > 10 {
				continue
			}
		}

		render := NewTemplatePoBGood(api, channelID)
		render.RenderView()
		logus.Log.Debug(fmt.Sprintf("render.main.Content=%v", render.main.GetMsgs()))

		assert.True(t, render.main.HasRecords())
		if msgs := render.main.GetMsgs(); len(msgs) > 0 {
			for _, msg := range msgs {
				rendered := msg.Render()
				fmt.Println(rendered)
			}
		}
		// // alerts
		// baseAlertDecreasing := configurator.NewCfgAlertBaseHealthIsDecreasing(configurator.NewConfigurator(dbpath))
		// isEnabled, _ := baseAlertDecreasing.Status(channelID)
		// assert.False(t, isEnabled)

	})
}

func TestPoBGoodiewerRealData(t *testing.T) {
	if false { // test pur efor debugging on real data

	}
	configurator.FixtureMigrator(func(dbpath types.Dbpath) {
		channelID, _ := configurator.FixtureChannel(dbpath)

		cg := configurator.NewConfiguratorBase(configurator.NewConfigurator(dbpath))
		cg.TagsAdd(channelID, []types.Tag{"Monte Cinto"}...)

		scrapper := scrappy.NewScrapyStorage(
			base.NewBaseApi(),
			player.FixturePlayerAPIMock(),
			baseattack.FixtureBaseAttackAPIMock(),
		)
		scrapper.Update()

		api := apis.NewAPI(dbpath, scrapper)

		cg_pobgood := configurator.NewConfiguratorPoBGood(configurator.NewConfigurator(dbpath))
		cg_pobgood.TagsAdd(channelID, []types.Tag{"commodity_nox_fake01", "commodity_military_salvage", "commodity_food"}...)

		render := NewTemplatePoBGood(api, channelID)
		render.RenderView()
		logus.Log.Debug(fmt.Sprintf("base.main.Msgs=%v", render.main.GetMsgs()))

		assert.True(t, render.main.HasRecords())
		if msgs := render.main.GetMsgs(); len(msgs) > 0 {
			for _, msg := range msgs {
				rendered := msg.Render()
				fmt.Println(rendered)
			}
		}
	})
}
