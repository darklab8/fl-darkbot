package pobgoodsview

import (
	_ "embed"
	"fmt"
	"sort"
	"text/template"

	"github.com/darklab8/fl-darkbot/app/settings/logus"
	"github.com/darklab8/fl-darkbot/app/settings/types"
	"github.com/darklab8/fl-darkbot/app/viewer/apis"
	"github.com/darklab8/fl-darkbot/app/viewer/views"
	"github.com/darklab8/fl-darkbot/app/viewer/views/baseview"
	"github.com/darklab8/fl-darkbot/app/viewer/views/viewer_msg"
	"github.com/darklab8/fl-darkstat/darkstat/configs_export"

	"github.com/darklab8/go-utils/utils"
	"github.com/darklab8/go-utils/utils/utils_types"
)

//go:embed pobgood_template.md
var pobGoodsTemplateText utils_types.TemplateExpression
var pobGoodsTemplate *template.Template

//go:embed pobgood_template_base.md
var pobGoodsTemplateBaseText utils_types.TemplateExpression
var pobGoodsTemplateBase *template.Template

func init() {
	pobGoodsTemplate = utils.TmpInit(pobGoodsTemplateText)
	pobGoodsTemplateBase = utils.TmpInit(pobGoodsTemplateBaseText)

}

// Base

type TemplatePoBGood struct {
	main                     *views.ViewTable
	alertPoBGoodLowerThan    *views.ViewTable
	alertPoBGoodAboveThan    *views.ViewTable
	WarningMissconfiguration *views.ViewTable

	api *apis.API
	*views.SharedViewTableSplitter
	channelID types.DiscordChannelID
}

const PobGoodViewID types.ViewID = "#darkbot-pobgood-view"

func NewTemplatePoBGood(api *apis.API, channelID types.DiscordChannelID) *TemplatePoBGood {
	base := TemplatePoBGood{}
	base.api = api
	base.channelID = channelID
	base.main = views.NewViewTable(viewer_msg.NewTableMsg(
		types.ViewID(PobGoodViewID),
		types.ViewHeader("**PoB Goods:** (Base, Quantity, Category, Item name)\n"),
		types.ViewBeginning("```scss\n"),
		types.ViewEnd("```"),
	))

	base.alertPoBGoodLowerThan = views.NewViewTable(viewer_msg.NewAlertMsg(
		types.ViewID("#darkbot-pobgood-below-than"),
	))
	base.alertPoBGoodAboveThan = views.NewViewTable(viewer_msg.NewAlertMsg(
		types.ViewID("#darkbot-pobgood-above-than"),
	))
	base.WarningMissconfiguration = views.NewViewTable(viewer_msg.NewAlertMsg(
		types.ViewID("#darkbot-pobgood-warning-missconfiguration"),
		viewer_msg.WithHeader(":warning: :warning: :warning: "),
	))

	base.SharedViewTableSplitter = views.NewSharedViewSplitter(
		api,
		channelID,
		&base,
		base.main,
		base.alertPoBGoodLowerThan,
		base.alertPoBGoodAboveThan,
		base.WarningMissconfiguration,
	)
	return &base
}

type PoBGood struct {
	AmountValue int
	GoodName    string
	Category    string
	IsEnd       bool
}
type BaseWithGoods struct {
	BaseName string
	Goods    []PoBGood
}

func SortGoods(goods []PoBGood) ([]PoBGood, error) {

	sort.Slice(goods, func(i, j int) bool {
		// if goods[i].BaseName != goods[j].BaseName {
		// 	return goods[i].BaseName < goods[j].BaseName
		// }
		if goods[i].Category != goods[j].Category {
			return goods[i].Category < goods[j].Category
		}
		return goods[i].GoodName < goods[j].GoodName
	})

	return goods, nil
}

type BaseGoods struct {
	Base  *configs_export.PoB
	Goods []*configs_export.ShopItem
}

func MatchGoods(bases []*configs_export.PoB, good_tags map[string]bool) []*BaseGoods {

	results := []*BaseGoods{}

	for _, base := range bases {
		goods_unique := make(map[string]*configs_export.ShopItem)
		for _, good := range base.ShopItems {

			_, ok := good_tags[good.Nickname]
			if !ok {
				continue
			}

			goods_unique[good.Nickname+base.Nickname] = good
		}

		base_goods := &BaseGoods{
			Base: base,
		}
		for _, good := range goods_unique {
			base_goods.Goods = append(base_goods.Goods, good)
		}
		if len(base_goods.Goods) == 0 {
			continue
		}
		results = append(results, base_goods)
	}
	return results
}
func (b *TemplatePoBGood) GenerateRecords() error {
	record, err := b.api.Scrappy.GetBaseStorage().GetLatestRecord()
	if logus.Log.CheckWarn(err, "unable to query bases from storage in Template base Generate records") {
		return err
	}
	sort.Slice(record.List, func(i, j int) bool {
		return record.List[i].Name < record.List[j].Name
	})

	tags, _ := b.api.Bases.Tags.TagsList(b.channelID)
	matchedBases := baseview.MatchBases(record.List, tags)

	order_key, err := b.api.Bases.OrderBy.Status(b.channelID)
	base_priorities, err2 := b.api.Bases.Priorities.Get(b.channelID)
	logus.Log.CheckDebug(err, "failed to query Order by key")
	if !logus.Log.CheckDebug(err2, "failed to base priorities") {
		base_priorities = make(map[string]int)
	}
	matchedBases, err = baseview.SortBases(matchedBases, types.OrderKey(order_key), base_priorities)

	base_table_will_be_rendered := len(matchedBases) > 0
	if err != nil && base_table_will_be_rendered {
		b.main.AppendRecord(types.ViewRecord(fmt.Sprintf("ERR %s", err.Error())))
		return err
	}

	// match goods
	good_tags, _ := b.api.PoBGood.Tags.TagsList(b.channelID)
	good_tags_map := make(map[string]bool)
	for _, tag := range good_tags {
		good_tags_map[string(tag)] = true
	}
	base_goods := MatchGoods(matchedBases, good_tags_map)

	bases := []BaseWithGoods{}
	for _, base := range base_goods {
		base_with_goods := BaseWithGoods{
			BaseName: base.Base.Name,
		}
		for _, good := range base.Goods {
			base_with_goods.Goods = append(base_with_goods.Goods, PoBGood{
				AmountValue: good.Quantity,
				GoodName:    good.Name,
				Category:    good.Category,
			})
		}
		bases = append(bases, base_with_goods)
	}
	for i, base := range bases {
		base.Goods, err = SortGoods(base.Goods)

		if len(base.Goods) > 0 {
			// makes sure base name remains white in scss markdown view, and pob good records green
			bases[i].Goods[len(base.Goods)-1].IsEnd = true
		}
	}

	var base_names map[string]bool = make(map[string]bool)

	is_first := true
	for _, base := range bases {
		if ok := base_names[base.BaseName]; !ok {
			if !is_first {
				b.main.AppendRecord(types.ViewRecord("\n"))
			}
			b.main.AppendRecord(types.ViewRecord(utils.TmpRender(pobGoodsTemplateBase, base)))
			base_names[base.BaseName] = true
			is_first = false
		}

		for _, good := range base.Goods {
			b.main.AppendRecord(types.ViewRecord(utils.TmpRender(pobGoodsTemplate, good)))
		}

	}

	if pobgood_thresholds, err := b.api.Alerts.PoBGoodsBelowThan.Get(b.channelID); err == nil {
		matched_pobgoods := make(map[string]bool)
		for pobgood_nickname, _ := range pobgood_thresholds {
			matched_pobgoods[pobgood_nickname] = false
		}

		for _, base := range base_goods {
			for _, good := range base.Goods {

				alert_threshold, ok := pobgood_thresholds[good.Nickname]
				if ok {
					matched_pobgoods[good.Nickname] = true
					if good.Quantity < alert_threshold {
						b.alertPoBGoodLowerThan.AppendRecord(types.ViewRecord(views.RenderAlertTemplate(
							b.channelID,
							fmt.Sprintf("Good %s has quantity %d below threshold %d at base %s\n", good.Name, good.Quantity, alert_threshold, base.Base.Name),
							b.api,
						)))
					}
				}
			}
		}

		for alert_config_pobgood_nickname, alert_config_is_matched := range matched_pobgoods {
			if !alert_config_is_matched {
				b.WarningMissconfiguration.AppendRecord(types.ViewRecord(views.RenderAlertTemplate(
					b.channelID,
					fmt.Sprintf("alert for pob good %s to be below threshold is configured, but all monitored bases have no such pob good data exposed! Add pob good exposure, see here for details: https://darkstat.dd84ai.com/#how_to_turn_pob_feature_on\\n",
						alert_config_pobgood_nickname,
					),
					b.api,
					views.WithAlertOverride(""),
				)))
			}
		}
	}

	if pobgood_thresholds, err := b.api.Alerts.PoBGoodsAboveThan.Get(b.channelID); err == nil {
		matched_pobgoods := make(map[string]bool)
		for pobgood_nickname, _ := range pobgood_thresholds {
			matched_pobgoods[pobgood_nickname] = false
		}

		for _, base := range base_goods {
			for _, good := range base.Goods {

				alert_threshold, ok := pobgood_thresholds[good.Nickname]
				if ok {
					matched_pobgoods[good.Nickname] = true
					if good.Quantity > alert_threshold {
						b.alertPoBGoodAboveThan.AppendRecord(types.ViewRecord(views.RenderAlertTemplate(
							b.channelID,
							fmt.Sprintf("Good %s has quantity %d above threshold %d at base %s\n", good.Name, good.Quantity, alert_threshold, base.Base.Name),
							b.api,
						)))
					}
				}
			}
		}

		for alert_config_pobgood_nickname, alert_config_is_matched := range matched_pobgoods {
			if !alert_config_is_matched {
				b.WarningMissconfiguration.AppendRecord(types.ViewRecord(views.RenderAlertTemplate(
					b.channelID,
					fmt.Sprintf("alert for pob good %s to be above threshold is configured, but all monitored bases have no such pob good data exposed! Add pob good exposure, see here for details: https://darkstat.dd84ai.com/#how_to_turn_pob_feature_on\n",
						alert_config_pobgood_nickname,
					),
					b.api,
					views.WithAlertOverride(""),
				)))
			}
		}
	}

	return nil
}
