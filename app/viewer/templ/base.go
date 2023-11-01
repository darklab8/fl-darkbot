package templ

import (
	"darkbot/app/scrappy/base"
	"darkbot/app/settings/logus"
	"darkbot/app/settings/types"
	"darkbot/app/settings/utils"
	"darkbot/app/viewer/apis"
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
	main                    TemplateShared
	AlertHealthLowerThan    TemplateShared
	AlertHealthIsDecreasing TemplateShared
	AlertBaseUnderAttack    TemplateShared
	API                     apis.API
}

func NewTemplateBase(channelID types.DiscordChannelID, dbpath types.Dbpath) TemplateBase {
	base := TemplateBase{}
	base.API = apis.NewAPI(channelID, dbpath)
	base.main.Header = "#darkbot-base-view"
	base.AlertHealthLowerThan.Header = "#darkbot-base-alert-health-lower-than"
	base.AlertHealthIsDecreasing.Header = "#darkbot-base-health-is-decreasing"
	base.AlertBaseUnderAttack.Header = "#darkbot-base-base-under-attack"
	return base
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

func BaseContainsTag(bas base.Base, tags []string) bool {
	for _, tag := range tags {
		if strings.Contains(bas.Name, tag) {
			return true
		}
	}

	return false
}

func MatchBases(bases []base.Base, tags []string) []base.Base {
	result := []base.Base{}
	for _, base := range bases {

		if !BaseContainsTag(base, tags) {
			continue
		}

		result = append(result, base)
	}
	return result
}

func (b *TemplateBase) Setup(channelID types.DiscordChannelID) {
	b.API.ChannelID = channelID
	b.main.MessageID = ""
	b.main.Content = ""
	b.AlertHealthLowerThan.MessageID = ""
	b.AlertHealthLowerThan.Content = ""
	b.AlertHealthIsDecreasing.MessageID = ""
	b.AlertHealthIsDecreasing.Content = ""
	b.AlertBaseUnderAttack.MessageID = ""
	b.AlertBaseUnderAttack.Content = ""
}

const HealthRateDecreasingThreshold = -0.01

func (b *TemplateBase) Render() {
	input := TemplateRendererBaseInput{
		Header:               b.main.Header,
		LastUpdated:          time.Now().String(),
		HealthDecreasePhrase: "\n@healthDecreasing;",
		UnderAttackPhrase:    "\n@underAttack;",
	}

	record, err := b.API.Scrappy.GetBaseStorage().GetLatestRecord()
	if logus.CheckWarn(err, "unable to render TemplateBase") {
		return
	}
	sort.Slice(record.List, func(i, j int) bool {
		return record.List[i].Name < record.List[j].Name
	})

	tags, _ := b.API.Bases.TagsList(b.API.ChannelID)

	matchedBases := MatchBases(record.List, tags)
	healthDeritives, healthDerivativeErr := CalculateDerivates(tags, b.API)
	DerivativesInitializing := healthDerivativeErr != nil

	for _, base := range matchedBases {
		healthDeritiveNumber, _ := healthDeritives[base.Name]
		healthDeritive := strconv.FormatFloat(healthDeritiveNumber, 'f', 4, 64)
		var HealthDecreasing, UnderAttack bool

		if DerivativesInitializing {
			healthDeritive = "initializing"
		} else {
			HealthDecreasing = healthDeritiveNumber < 0
			UnderAttack = healthDeritiveNumber < HealthRateDecreasingThreshold || strings.Contains(string(b.API.Scrappy.GetBaseAttackStorage().GetData()), base.Name)
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
		return
	}

	b.AlertHealthLowerThan.Content = ""
	if healthThreshold, _ := b.API.Alerts.BaseHealthLowerThan.Status(b.API.ChannelID); healthThreshold != nil {
		for _, base := range input.Bases {
			if int(base.Health) < *healthThreshold {
				b.AlertHealthLowerThan.Content = RenderAlertTemplate(b.AlertHealthLowerThan.Header, b.API.ChannelID, fmt.Sprintf("Base %s has health %d lower than threshold %d", base.Name, int(base.Health), *healthThreshold), b.API)
				break
			}
		}
	}

	b.AlertHealthIsDecreasing.Content = ""
	if isAlertEnabled, _ := b.API.Alerts.BaseHealthIsDecreasing.Status(b.API.ChannelID); isAlertEnabled {
		for _, base := range input.Bases {
			if base.IsHealthDecreasing {
				b.AlertHealthIsDecreasing.Content = RenderAlertTemplate(b.AlertHealthIsDecreasing.Header, b.API.ChannelID, fmt.Sprintf("Base %s health %d is decreasing with value %s", base.Name, int(base.Health), base.HealthChange), b.API)
				break
			}
		}
	}

	b.AlertBaseUnderAttack.Content = ""
	if isAlertEnabled, _ := b.API.Alerts.BaseIsUnderAttack.Status(b.API.ChannelID); isAlertEnabled {
		for _, base := range input.Bases {
			if base.IsUnderAttack {
				b.AlertBaseUnderAttack.Content = RenderAlertTemplate(b.AlertBaseUnderAttack.Header, b.API.ChannelID, fmt.Sprintf("Base %s health %d is probably under attack because health change %s is dropping faster than %f. Or it was detected at forum attack declaration thread.", base.Name, int(base.Health), base.HealthChange, HealthRateDecreasingThreshold), b.API)
				break
			}
		}
	}
}

func (t *TemplateBase) Send() {
	t.main.Send(t.API)
	t.AlertHealthLowerThan.Send(t.API)
	t.AlertHealthIsDecreasing.Send(t.API)
	t.AlertBaseUnderAttack.Send(t.API)
}

func (t *TemplateBase) MatchMessageID(messageID types.DiscordMessageID) bool {

	if messageID == t.main.MessageID {
		return true
	}
	if messageID == t.AlertHealthLowerThan.MessageID {
		return true
	}
	if messageID == t.AlertHealthIsDecreasing.MessageID {
		return true
	}
	if messageID == t.AlertBaseUnderAttack.MessageID {
		return true
	}
	return false
}

func (t *TemplateBase) DiscoverMessageID(content string, msgID types.DiscordMessageID) {
	if strings.Contains(content, t.main.Header) {
		t.main.MessageID = msgID
		t.main.Content = content
	}
	if strings.Contains(content, t.AlertHealthLowerThan.Header) {
		t.AlertHealthLowerThan.MessageID = msgID
		t.AlertHealthLowerThan.Content = content
	}
	if strings.Contains(content, t.AlertHealthIsDecreasing.Header) {
		t.AlertHealthIsDecreasing.MessageID = msgID
		t.AlertHealthIsDecreasing.Content = content
	}
	if strings.Contains(content, t.AlertBaseUnderAttack.Header) {
		t.AlertBaseUnderAttack.MessageID = msgID
		t.AlertBaseUnderAttack.Content = content
	}
}
