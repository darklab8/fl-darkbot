package baseview

import (
	"darkbot/app/scrappy/base"
	"darkbot/app/settings/logus"
	"darkbot/app/settings/types"
	"darkbot/app/settings/utils"
	"darkbot/app/viewer/apis"
	"darkbot/app/viewer/views"
	_ "embed"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"text/template"
)

//go:embed base_template.md
var baseMarkup string
var baseTemplate *template.Template

func init() {
	baseTemplate = utils.TmpInit(baseMarkup)
}

// Base

type TemplateBase struct {
	main                    views.ViewTable
	alertHealthLowerThan    views.ViewTable
	alertHealthIsDecreasing views.ViewTable
	alertBaseUnderAttack    views.ViewTable
	api                     *apis.API
	*views.SharedViewTableSplitter
	channelID types.DiscordChannelID
}

func NewTemplateBase(api *apis.API, channelID types.DiscordChannelID) *TemplateBase {
	base := TemplateBase{}
	base.api = api
	base.channelID = channelID
	base.main.ViewID = "#darkbot-base-view"
	base.alertHealthLowerThan.ViewID = "#darkbot-base-alert-health-lower-than"
	base.alertHealthIsDecreasing.ViewID = "#darkbot-base-health-is-decreasing"
	base.alertBaseUnderAttack.ViewID = "#darkbot-base-base-under-attack"

	base.SharedViewTableSplitter = views.NewSharedViewSplitter(
		api,
		channelID,
		&base,
		&base.main,
		&base.alertHealthLowerThan,
		&base.alertHealthIsDecreasing,
		&base.alertBaseUnderAttack,
	)
	return &base
}

type TemplateAugmentedBase struct {
	base.Base
	HealthChange         string
	IsHealthDecreasing   bool
	IsUnderAttack        bool
	HealthDecreasePhrase string
	UnderAttackPhrase    string
}

func BaseContainsTag(bas base.Base, tags []types.Tag) bool {
	for _, tag := range tags {
		if strings.Contains(bas.Name, string(tag)) {
			return true
		}
	}

	return false
}

func MatchBases(bases []base.Base, tags []types.Tag) []base.Base {
	result := []base.Base{}
	for _, base := range bases {

		if !BaseContainsTag(base, tags) {
			continue
		}

		result = append(result, base)
	}
	return result
}

const HealthRateDecreasingThreshold = -0.01

func (b *TemplateBase) GenerateRecords() error {
	record, err := b.api.Scrappy.GetBaseStorage().GetLatestRecord()
	if logus.CheckWarn(err, "unable to query bases from storage in Template base Generate records") {
		return err
	}
	sort.Slice(record.List, func(i, j int) bool {
		return record.List[i].Name < record.List[j].Name
	})

	var beginning strings.Builder
	beginning.WriteString("**Bases:**\n")
	b.main.ViewBeginning = types.ViewBeginning(beginning.String())

	HealthDecreasePhrase := "\n@healthDecreasing;"
	UnderAttackPhrase := "\n@underAttack;"
	bases := []TemplateAugmentedBase{}

	tags, _ := b.api.Bases.TagsList(b.channelID)

	matchedBases := MatchBases(record.List, tags)
	healthDeritives, healthDerivativeErr := CalculateDerivates(tags, b.api)
	DerivativesInitializing := healthDerivativeErr != nil

	for _, base := range matchedBases {
		healthDeritiveNumber := healthDeritives[base.Name]

		healthDeritive := strconv.FormatFloat(healthDeritiveNumber, 'f', 4, 64)
		var HealthDecreasing, UnderAttack bool

		if DerivativesInitializing {
			healthDeritive = "initializing"
		} else {
			HealthDecreasing = healthDeritiveNumber < 0
			UnderAttack = healthDeritiveNumber < HealthRateDecreasingThreshold || strings.Contains(string(b.api.Scrappy.GetBaseAttackStorage().GetData()), base.Name)
		}

		baseVars := TemplateAugmentedBase{
			Base:                 base,
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
			if int(base.Health) < healthThreshold {
				b.alertHealthLowerThan.AppendRecord(views.RenderAlertTemplate(
					b.channelID,
					fmt.Sprintf("Base %s has health %d lower than threshold %d", base.Name, int(base.Health), healthThreshold),
					b.api,
				))
				break
			}
		}
	}

	if isAlertEnabled, err := b.api.Alerts.BaseHealthIsDecreasing.Status(b.channelID); err == nil && isAlertEnabled {
		for _, base := range bases {
			if base.IsHealthDecreasing {
				b.alertHealthIsDecreasing.AppendRecord(views.RenderAlertTemplate(
					b.channelID,
					fmt.Sprintf("Base %s health %d is decreasing with value %s", base.Name, int(base.Health), base.HealthChange),
					b.api,
				))
				break
			}
		}
	}

	if isAlertEnabled, _ := b.api.Alerts.BaseIsUnderAttack.Status(b.channelID); isAlertEnabled {
		for _, base := range bases {
			if base.IsUnderAttack {
				b.alertBaseUnderAttack.AppendRecord(views.RenderAlertTemplate(
					b.channelID,
					fmt.Sprintf("Base %s health %d is probably under attack because health change %s is dropping faster than %f. Or it was detected at forum attack declaration thread.",
						base.Name,
						int(base.Health),
						base.HealthChange,
						HealthRateDecreasingThreshold,
					),
					b.api,
				))
				break
			}
		}
	}

	return nil
}
