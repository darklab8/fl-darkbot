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

func init() {
	pobGoodsTemplate = utils.TmpInit(pobGoodsTemplateText)
}

// Base

type TemplatePoBGood struct {
	main *views.ViewTable
	// alertPoBGoodLowerThan    *views.ViewTable
	// alertPoBGoodAboveThan    *views.ViewTable
	// WarningMissconfiguration *views.ViewTable

	api *apis.API
	*views.SharedViewTableSplitter
	channelID types.DiscordChannelID
}

func NewTemplatePoBGood(api *apis.API, channelID types.DiscordChannelID) *TemplatePoBGood {
	base := TemplatePoBGood{}
	base.api = api
	base.channelID = channelID
	base.main = views.NewViewTable(viewer_msg.NewTableMsg(
		types.ViewID("#darkbot-pobgood-view"),
		types.ViewHeader("**PoB Goods:**\n"),
		types.ViewBeginning("```scss\n"),
		types.ViewEnd("```"),
	))

	//
	// base.alertPoBGoodLowerThan = views.NewViewTable(viewer_msg.NewAlertMsg(
	// 	types.ViewID("#darkbot-pobgood-below-than"),
	// ))
	// base.alertPoBGoodAboveThan = views.NewViewTable(viewer_msg.NewAlertMsg(
	// 	types.ViewID("#darkbot-pobgood-above-than"),
	// ))
	// base.WarningMissconfiguration = views.NewViewTable(viewer_msg.NewAlertMsg(
	// 	types.ViewID("#darkbot-pobgood-warning-missconfiguration"),
	// ))

	base.SharedViewTableSplitter = views.NewSharedViewSplitter(
		api,
		channelID,
		&base,
		base.main,
		// base.alertPoBGoodLowerThan,
		// base.alertPoBGoodAboveThan,
	)
	return &base
}

type PoBGood struct {
	AmountValue int
	GoodName    string
	BaseName    string
	Category    string
}

type ForbiddenOrderKey struct{ order_key types.OrderKey }

func ErrorForbiddenOrderKey(order_key types.OrderKey) ForbiddenOrderKey {
	return ForbiddenOrderKey{order_key: order_key}
}

func (f ForbiddenOrderKey) Error() string { return fmt.Sprintf("Forbidden order key=%s", f.order_key) }

func SortGoods(goods []PoBGood) ([]PoBGood, error) {

	sort.Slice(goods, func(i, j int) bool {
		if goods[i].Category != goods[j].Category {
			return goods[i].Category < goods[j].Category
		}
		return goods[i].GoodName < goods[j].GoodName
	})

	return goods, nil
}

type MatchedGood struct {
	Base *configs_export.PoB
	Good *configs_export.ShopItem
}

func MatchGoods(bases []*configs_export.PoB, good_tags map[string]bool) map[string]MatchedGood {
	result := make(map[string]MatchedGood)
	for _, base := range bases {

		for _, good := range base.ShopItems {

			_, ok := good_tags[good.Nickname]
			if !ok {
				continue
			}

			matched_good := MatchedGood{
				Base: base,
				Good: good,
			}
			result[matched_good.Good.Nickname+matched_good.Base.Nickname] = matched_good
		}

	}
	return result
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
	if !logus.Log.CheckDebug(err, "failed to query Order by key") {
		matchedBases, err = baseview.SortBases(matchedBases, types.OrderKey(order_key))

		base_table_will_be_rendered := len(matchedBases) > 0
		if err != nil && base_table_will_be_rendered {
			b.main.AppendRecord(types.ViewRecord(fmt.Sprintf("ERR %s", err.Error())))
			return err
		}
	}

	goods := []PoBGood{}

	// match goods
	good_tags, _ := b.api.PoBGood.Tags.TagsList(b.channelID)
	good_tags_map := make(map[string]bool)
	for _, tag := range good_tags {
		good_tags_map[string(tag)] = true
	}
	matchedGoods := MatchGoods(matchedBases, good_tags_map)

	for _, matched_good := range matchedGoods {
		goods = append(goods, PoBGood{
			AmountValue: matched_good.Good.Quantity,
			GoodName:    matched_good.Good.Name,
			BaseName:    matched_good.Base.Name,
			Category:    matched_good.Good.Category,
		})
	}
	goods, err = SortGoods(goods)

	for _, good := range goods {
		b.main.AppendRecord(types.ViewRecord(utils.TmpRender(pobGoodsTemplate, good)))
	}

	return nil
}
