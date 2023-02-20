package templ

import (
	"darkbot/dtypes"
	"darkbot/scrappy/base"
	"darkbot/scrappy/shared/records"
	"darkbot/utils"
	"darkbot/viewer/apis"
	_ "embed"
	"fmt"
	"sort"
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

func NewTemplateBase(channelID string, dbpath dtypes.Dbpath) TemplateBase {
	base := TemplateBase{}
	base.API = apis.NewAPI(channelID, dbpath)
	base.main.Header = "#darkbot-base-view"
	base.AlertHealthLowerThan.Header = "#darkbot-base-alert-health-lower-than"
	base.AlertHealthIsDecreasing.Header = "#darkbot-base-health-is-decreasing"
	base.AlertBaseUnderAttack.Header = "#darkbot-base-base-under-attack"
	return base
}

type AugmentedBase struct {
	base.Base
	HealthChange       float64
	IsHealthDecreasing bool
	IsUnderAttack      bool
}

type TemplateRendererBaseInput struct {
	Header      string
	LastUpdated string
	Bases       []AugmentedBase
}

func BaseContainsTag(bas base.Base, tags []string) bool {
	for _, tag := range tags {
		if strings.Contains(bas.Name, tag) {
			return true
		}
	}

	return false
}

func MatchBases(record records.StampedObjects[base.Base], tags []string) []base.Base {
	result := []base.Base{}
	for _, base := range record.List {

		if !BaseContainsTag(base, tags) {
			continue
		}

		result = append(result, base)
	}
	return result
}

func (b *TemplateBase) Setup(channelID string) {
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

const HealthRateDecreasingThreshold = -0.002200 * 2

func (b *TemplateBase) Render() {
	input := TemplateRendererBaseInput{
		Header:      b.main.Header,
		LastUpdated: time.Now().String(),
	}

	record, err := b.API.Scrappy.BaseStorage.GetLatestRecord()
	if err != nil {
		return
	}
	sort.Slice(record.List, func(i, j int) bool {
		return record.List[i].Name < record.List[j].Name
	})

	tags, _ := b.API.Bases.TagsList(b.API.ChannelID)

	matchedBases := MatchBases(record, tags)
	healthDeritives := CalculateDerivates(tags, b.API)

	for _, base := range matchedBases {
		healthDeritive, _ := healthDeritives[base.Name]

		HealthDecreasing := healthDeritive < 0
		UnderAttack := healthDeritive < HealthRateDecreasingThreshold
		input.Bases = append(input.Bases, AugmentedBase{
			Base:               base,
			HealthChange:       healthDeritive,
			IsHealthDecreasing: HealthDecreasing,
			IsUnderAttack:      UnderAttack,
		})
	}

	if len(input.Bases) != 0 {
		b.main.Content = utils.TmpRender(baseTemplate, input)
	}

	// Alerts
	if healthThreshold, _ := b.API.Alerts.BaseHealthLowerThan.Status(b.API.ChannelID); healthThreshold != nil {
		for _, base := range record.List {
			if int(base.Health) < *healthThreshold {
				b.AlertHealthLowerThan.Content = RenderAlertTemplate(b.AlertHealthLowerThan.Header, b.API.ChannelID, fmt.Sprintf("Base %s has health %d lower than threshold %d", base.Name, int(base.Health), *healthThreshold), b.API)
			}
		}
	}

	for _, base := range input.Bases {
		if base.IsHealthDecreasing {
			b.AlertHealthIsDecreasing.Content = RenderAlertTemplate(b.AlertHealthIsDecreasing.Header, b.API.ChannelID, fmt.Sprintf("Base %s health %d is decreasing with value %f", base.Name, int(base.Health), base.HealthChange), b.API)
			break
		}
	}

	for _, base := range input.Bases {
		if base.IsUnderAttack {
			b.AlertBaseUnderAttack.Content = RenderAlertTemplate(b.AlertBaseUnderAttack.Header, b.API.ChannelID, fmt.Sprintf("Base %s health %d is probably under attack because health change %f is dropping faster than %f", base.Name, int(base.Health), base.HealthChange, HealthRateDecreasingThreshold), b.API)
			break
		}
	}
}

func (t *TemplateBase) Send() {
	t.main.Send(t.API)
	t.AlertHealthLowerThan.Send(t.API)
	t.AlertHealthIsDecreasing.Send(t.API)
	t.AlertBaseUnderAttack.Send(t.API)
}

func (t *TemplateBase) MatchMessageID(messageID string) bool {

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

func (t *TemplateBase) DiscoverMessageID(content string, msgID string) {
	if strings.Contains(content, t.main.Header) {
		t.main.MessageID = msgID
	}
	if strings.Contains(content, t.AlertHealthLowerThan.Header) {
		t.AlertHealthLowerThan.MessageID = msgID
	}
	if strings.Contains(content, t.AlertHealthIsDecreasing.Header) {
		t.AlertHealthIsDecreasing.MessageID = msgID
	}
	if strings.Contains(content, t.AlertBaseUnderAttack.Header) {
		t.AlertBaseUnderAttack.MessageID = msgID
	}
}
