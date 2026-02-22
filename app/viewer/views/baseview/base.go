package baseview

import (
	_ "embed"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"text/template"

	"github.com/darklab8/fl-darkbot/app/configurator/models"
	"github.com/darklab8/fl-darkbot/app/settings/logus"
	"github.com/darklab8/fl-darkbot/app/settings/types"
	"github.com/darklab8/fl-darkbot/app/viewer/apis"
	"github.com/darklab8/fl-darkbot/app/viewer/views"
	"github.com/darklab8/fl-darkbot/app/viewer/views/viewer_msg"
	"github.com/darklab8/fl-darkstat/darkstat/configs_export"

	"github.com/darklab8/go-utils/typelog"
	"github.com/darklab8/go-utils/utils"
	"github.com/darklab8/go-utils/utils/utils_types"
)

//go:embed base_template.md
var baseMarkup utils_types.TemplateExpression
var baseTemplate *template.Template

func init() {
	baseTemplate = utils.TmpInit(baseMarkup)
}

// Base

type TemplateBase struct {
	main                     *views.ViewTable
	alertHealthLowerThan     *views.ViewTable
	alertHealthIsDecreasing  *views.ViewTable
	alertBaseUnderAttack     *views.ViewTable
	alertMoneyBelow          *views.ViewTable
	AlertCargospaceLeftBelow *views.ViewTable
	WarningMissconfiguration *views.ViewTable

	api *apis.API
	*views.SharedViewTableSplitter
	channelID types.DiscordChannelID
}

const BaseViewID types.ViewID = "#darkbot-base-view"

func NewTemplateBase(api *apis.API, channelID types.DiscordChannelID) *TemplateBase {
	base := TemplateBase{}
	base.api = api
	base.channelID = channelID
	base.main = views.NewViewTable(viewer_msg.NewTableMsg(
		types.ViewID(BaseViewID),
		types.ViewHeader("**Bases:**\n"),
		types.ViewBeginning(""),
		types.ViewEnd(""),
	))

	//
	base.alertHealthLowerThan = views.NewViewTable(viewer_msg.NewAlertMsg(
		types.ViewID("#darkbot-base-alert-health-lower-than"),
	))
	base.alertHealthIsDecreasing = views.NewViewTable(viewer_msg.NewAlertMsg(
		types.ViewID("#darkbot-base-health-is-decreasing"),
	))
	base.alertBaseUnderAttack = views.NewViewTable(viewer_msg.NewAlertMsg(
		types.ViewID("#darkbot-base-base-under-attack"),
	))

	base.alertMoneyBelow = views.NewViewTable(viewer_msg.NewAlertMsg(
		types.ViewID("#darkbot-base-money-below"),
	))
	base.AlertCargospaceLeftBelow = views.NewViewTable(viewer_msg.NewAlertMsg(
		types.ViewID("#darkbot-base-cargospace-below"),
	))

	base.WarningMissconfiguration = views.NewViewTable(viewer_msg.NewAlertMsg(
		types.ViewID("#darkbot-base-warning-missconfiguration"),
		viewer_msg.WithHeader(":warning: :warning: :warning: "),
	))

	base.SharedViewTableSplitter = views.NewSharedViewSplitter(
		api,
		channelID,
		&base,
		base.main,
		base.alertHealthLowerThan,
		base.alertHealthIsDecreasing,
		base.alertBaseUnderAttack,
		base.alertMoneyBelow,
		base.AlertCargospaceLeftBelow,
		base.WarningMissconfiguration,
	)
	return &base
}

type TemplateAugmentedBase struct {
	*configs_export.PoB
	HealthChange         string
	IsHealthDecreasing   bool
	IsUnderAttack        bool
	HealthDecreasePhrase string
	UnderAttackPhrase    string
}

func BaseContainsTag(bas *configs_export.PoB, tags []types.Tag) bool {
	for _, tag := range tags {
		if strings.Contains(strings.ToLower(bas.Name), strings.ToLower(string(tag))) {
			return true
		}
	}

	return false
}

func MatchBases(bases []*configs_export.PoB, tags []types.Tag) []*configs_export.PoB {
	result := []*configs_export.PoB{}
	for _, base := range bases {

		if !BaseContainsTag(base, tags) {
			continue
		}

		result = append(result, base)
	}
	return result
}

type ForbiddenOrderKey struct{ order_key types.OrderKey }

func ErrorForbiddenOrderKey(order_key types.OrderKey) ForbiddenOrderKey {
	return ForbiddenOrderKey{order_key: order_key}
}

func (f ForbiddenOrderKey) Error() string { return fmt.Sprintf("Forbidden order key=%s", f.order_key) }

func SortBases(bases []*configs_export.PoB, order_key types.OrderKey) ([]*configs_export.PoB, error) {

	switch order_key {
	case models.BaseKeyName:
		sort.Slice(bases, func(i, j int) bool {
			return bases[i].Name < bases[j].Name
		})
	case models.BaseKeyAffiliation:
		sort.Slice(bases, func(i, j int) bool {
			return types.GetS(bases[i].FactionName, "") < types.GetS(bases[j].FactionName, "")
		})
	default:
		logus.Log.Error(fmt.Sprintf("forbidden order order_key=%s, only keys=%v are allowed", order_key, models.ConfigBaseOrderingKeyAllowedTags))
		return bases, ErrorForbiddenOrderKey(order_key)
	}

	return bases, nil
}

func (b *TemplateBase) GenerateRecords() error {
	record, err := b.api.Scrappy.GetBaseStorage().GetLatestRecord()
	if logus.Log.CheckWarn(err, "unable to query bases from storage in Template base Generate records") {
		return err
	}
	sort.Slice(record.List, func(i, j int) bool {
		return record.List[i].Name < record.List[j].Name
	})

	HealthDecreasePhrase := "\n@healthDecreasing;"
	UnderAttackPhrase := "\n@underAttack;"
	bases := []TemplateAugmentedBase{}

	tags, _ := b.api.Bases.Tags.TagsList(b.channelID)
	matchedBases := MatchBases(record.List, tags)

	order_key, err := b.api.Bases.OrderBy.Status(b.channelID)
	if !logus.Log.CheckDebug(err, "failed to query Order by key") {
		matchedBases, err = SortBases(matchedBases, types.OrderKey(order_key))

		base_table_will_be_rendered := len(matchedBases) > 0
		if err != nil && base_table_will_be_rendered {
			b.main.AppendRecord(types.ViewRecord(fmt.Sprintf("ERR %s", err.Error())))
			return err
		}
	}

	healthDeritives, healthDerivativeErr := CalculateDerivates(tags, b.api)
	DerivativesInitializing := healthDerivativeErr != nil

	for _, base := range matchedBases {
		healthDeritiveNumber := healthDeritives[base.Name]

		healthDeritive := strconv.FormatFloat(healthDeritiveNumber, 'f', 4, 64)
		var HealthDecreasing, UnderAttack bool

		if DerivativesInitializing {
			healthDeritive = "initializing"
			logus.Log.Warn("healthDeritive is initializing, because", typelog.OptError(healthDerivativeErr))
		} else {
			HealthDecreasing = healthDeritiveNumber < 0
			UnderAttack = strings.Contains(string(b.api.Scrappy.GetBaseAttackStorage().GetData()), base.Name)
		}

		baseVars := TemplateAugmentedBase{
			PoB:                  base,
			HealthChange:         healthDeritive,
			IsHealthDecreasing:   HealthDecreasing,
			IsUnderAttack:        UnderAttack,
			HealthDecreasePhrase: HealthDecreasePhrase,
			UnderAttackPhrase:    UnderAttackPhrase,
		}
		bases = append(bases, baseVars)

	}

	for _, base := range bases {
		b.main.AppendRecord(types.ViewRecord(utils.TmpRender(baseTemplate, base)))
	}

	if healthThreshold, err := b.api.Alerts.BaseHealthLowerThan.Status(b.channelID); err == nil {
		for _, base := range bases {
			if base.Health == nil {
				continue
			}
			if int(*base.Health) < healthThreshold {
				b.alertHealthLowerThan.AppendRecord(types.ViewRecord(views.RenderAlertTemplate(
					b.channelID,
					fmt.Sprintf("Base %s has health %d lower than threshold %d", base.Name, int(*base.Health), healthThreshold),
					b.api,
				)))
				break
			}
		}
	}

	if isAlertEnabled, err := b.api.Alerts.BaseHealthIsDecreasing.Status(b.channelID); err == nil && isAlertEnabled {
		for _, base := range bases {
			if base.IsHealthDecreasing && base.Health != nil {
				b.alertHealthIsDecreasing.AppendRecord(types.ViewRecord(views.RenderAlertTemplate(
					b.channelID,
					fmt.Sprintf("Base %s health %d is decreasing with value %s", base.Name, int(*base.Health), base.HealthChange),
					b.api,
				)))
				break
			}
		}
	}

	if isAlertEnabled, _ := b.api.Alerts.BaseIsUnderAttack.Status(b.channelID); isAlertEnabled {
		for _, base := range bases {
			if base.IsUnderAttack {
				b.alertBaseUnderAttack.AppendRecord(types.ViewRecord(views.RenderAlertTemplate(
					b.channelID,
					fmt.Sprintf("Base %s health %d is under attack, because we detected base name at forum attack declaration thread.",
						base.Name,
						int(types.GetF(base.Health, -1)),
					),
					b.api,
				)))
				break
			}
		}
	}

	if money_threshold, err := b.api.Alerts.BaseMoneyBelowThan.Status(b.channelID); err == nil {
		for _, base := range bases {

			if base.Money == nil {
				b.WarningMissconfiguration.AppendRecord(types.ViewRecord(views.RenderAlertTemplate(
					b.channelID,
					fmt.Sprintf("Base %s has no data for money, but has configured money alert. No permission present to see it. Fix in Player Owned Base account manager at forum\n",
						base.Name,
					),
					b.api,
					views.WithAlertOverride(""),
				)))
				continue
			}

			if *base.Money < money_threshold {
				b.alertMoneyBelow.AppendRecord(types.ViewRecord(views.RenderAlertTemplate(
					b.channelID,
					fmt.Sprintf("Base %s money %d is below threshold %d. Fix money at your base :)\n",
						base.Name,
						int(types.GetI(base.Money, -1)),
						money_threshold,
					),
					b.api,
				)))
				break
			}
		}
	}

	if cargo_threshold, err := b.api.Alerts.BaseCargoBelowThan.Status(b.channelID); err == nil {
		for _, base := range bases {

			if base.CargoSpaceLeft == nil {
				b.WarningMissconfiguration.AppendRecord(types.ViewRecord(views.RenderAlertTemplate(
					b.channelID,
					fmt.Sprintf("Base %s has no data for cargo space left, but has alert for cargo configured. No permission present to see it. Fix in Player Owned Base account manager at forum\n",
						base.Name,
					),
					b.api,
					views.WithAlertOverride(""),
				)))
				continue
			}

			if *base.CargoSpaceLeft < cargo_threshold {
				b.AlertCargospaceLeftBelow.AppendRecord(types.ViewRecord(views.RenderAlertTemplate(
					b.channelID,
					fmt.Sprintf("Base %s cargo space left %d is below threshold %d. Fix cargo left space at your base :)\n",
						base.Name,
						int(types.GetI(base.Money, -1)),
						cargo_threshold,
					),
					b.api,
				)))
				break
			}
		}
	}

	return nil
}
