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
	"time"
)

//go:embed base_template.md
var baseMarkup string
var baseTemplate *template.Template

func init() {
	baseTemplate = utils.TmpInit(baseMarkup)
}

// Base

type TemplateBase struct {
	main                    views.TemplateShared
	alertHealthLowerThan    views.TemplateShared
	alertHealthIsDecreasing views.TemplateShared
	alertBaseUnderAttack    views.TemplateShared
	api                     *apis.API
}

func NewTemplateBase(api *apis.API) *TemplateBase {
	base := TemplateBase{}
	base.api = api
	base.main.Header = "#darkbot-base-view"
	base.alertHealthLowerThan.Header = "#darkbot-base-alert-health-lower-than"
	base.alertHealthIsDecreasing.Header = "#darkbot-base-health-is-decreasing"
	base.alertBaseUnderAttack.Header = "#darkbot-base-base-under-attack"
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

type TemplateRendererBaseInput struct {
	Header               string
	LastUpdated          string
	Bases                []TemplateAugmentedBase
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

func (b *TemplateBase) Render() error {
	input := TemplateRendererBaseInput{
		Header:               b.main.Header,
		LastUpdated:          time.Now().String(),
		HealthDecreasePhrase: "\n@healthDecreasing;",
		UnderAttackPhrase:    "\n@underAttack;",
	}

	record, err := b.api.Scrappy.GetBaseStorage().GetLatestRecord()
	if logus.CheckWarn(err, "unable to render TemplateBase") {
		return err
	}
	sort.Slice(record.List, func(i, j int) bool {
		return record.List[i].Name < record.List[j].Name
	})

	tags, _ := b.api.Bases.TagsList(b.api.ChannelID)

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

		input.Bases = append(input.Bases, TemplateAugmentedBase{
			Base:                 base,
			HealthChange:         healthDeritive,
			IsHealthDecreasing:   HealthDecreasing,
			IsUnderAttack:        UnderAttack,
			HealthDecreasePhrase: input.HealthDecreasePhrase,
			UnderAttackPhrase:    input.UnderAttackPhrase,
		})
	}

	if len(input.Bases) != 0 {
		b.main.Content = utils.TmpRender(baseTemplate, input)
	}

	// Alerts
	if DerivativesInitializing {
		// Don't update alerts until bases are properly initalized. To avoid extra pings to players
		return nil
	}

	b.alertHealthLowerThan.Content = ""
	if healthThreshold, err := b.api.Alerts.BaseHealthLowerThan.Status(b.api.ChannelID); err == nil {
		for _, base := range input.Bases {
			if int(base.Health) < healthThreshold {
				b.alertHealthLowerThan.Content = views.RenderAlertTemplate(b.alertHealthLowerThan.Header, b.api.ChannelID, fmt.Sprintf("Base %s has health %d lower than threshold %d", base.Name, int(base.Health), healthThreshold), b.api)
				break
			}
		}
	}

	b.alertHealthIsDecreasing.Content = ""
	if isAlertEnabled, err := b.api.Alerts.BaseHealthIsDecreasing.Status(b.api.ChannelID); err == nil && isAlertEnabled {
		for _, base := range input.Bases {
			if base.IsHealthDecreasing {
				b.alertHealthIsDecreasing.Content = views.RenderAlertTemplate(b.alertHealthIsDecreasing.Header, b.api.ChannelID, fmt.Sprintf("Base %s health %d is decreasing with value %s", base.Name, int(base.Health), base.HealthChange), b.api)
				break
			}
		}
	}

	b.alertBaseUnderAttack.Content = ""
	if isAlertEnabled, _ := b.api.Alerts.BaseIsUnderAttack.Status(b.api.ChannelID); isAlertEnabled {
		for _, base := range input.Bases {
			if base.IsUnderAttack {
				b.alertBaseUnderAttack.Content = views.RenderAlertTemplate(b.alertBaseUnderAttack.Header, b.api.ChannelID, fmt.Sprintf("Base %s health %d is probably under attack because health change %s is dropping faster than %f. Or it was detected at forum attack declaration thread.", base.Name, int(base.Health), base.HealthChange, HealthRateDecreasingThreshold), b.api)
				break
			}
		}
	}

	return nil
}

func (t *TemplateBase) Send() {
	t.main.Send(t.api)
	t.alertHealthLowerThan.Send(t.api)
	t.alertHealthIsDecreasing.Send(t.api)
	t.alertBaseUnderAttack.Send(t.api)
}

func (t *TemplateBase) MatchMessageID(messageID types.DiscordMessageID) bool {

	if messageID == t.main.MessageID {
		return true
	}
	if messageID == t.alertHealthLowerThan.MessageID {
		return true
	}
	if messageID == t.alertHealthIsDecreasing.MessageID {
		return true
	}
	if messageID == t.alertBaseUnderAttack.MessageID {
		return true
	}
	return false
}

func (t *TemplateBase) DiscoverMessageID(content string, msgID types.DiscordMessageID) {
	if strings.Contains(content, t.main.Header) {
		t.main.MessageID = msgID
		t.main.Content = content
	}
	if strings.Contains(content, t.alertHealthLowerThan.Header) {
		t.alertHealthLowerThan.MessageID = msgID
		t.alertHealthLowerThan.Content = content
	}
	if strings.Contains(content, t.alertHealthIsDecreasing.Header) {
		t.alertHealthIsDecreasing.MessageID = msgID
		t.alertHealthIsDecreasing.Content = content
	}
	if strings.Contains(content, t.alertBaseUnderAttack.Header) {
		t.alertBaseUnderAttack.MessageID = msgID
		t.alertBaseUnderAttack.Content = content
	}
}
