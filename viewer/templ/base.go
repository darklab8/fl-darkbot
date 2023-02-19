package templ

import (
	"darkbot/dtypes"
	"darkbot/scrappy/base"
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

type TemplateRendererBaseInput struct {
	Header      string
	LastUpdated string
	Bases       []base.Base
}

func BaseContainsTag(bas base.Base, tags []string) bool {
	for _, tag := range tags {
		if strings.Contains(bas.Name, tag) {
			return true
		}
	}

	return false
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

	for _, base := range record.List {

		if !BaseContainsTag(base, tags) {
			continue
		}

		input.Bases = append(input.Bases, base)
	}

	if len(input.Bases) == 0 {
		return
	}

	b.main.Content = utils.TmpRender(baseTemplate, input)

	// Alerts
	if healthThreshold, _ := b.API.Alerts.BaseHealthLowerThan.Status(b.API.ChannelID); healthThreshold != nil {
		for _, base := range record.List {
			if int(base.Health) < *healthThreshold {
				b.AlertHealthLowerThan.Content = RenderAlertTemplate(b.AlertHealthLowerThan.Header, b.API.ChannelID, fmt.Sprintf("Base %s has health %d lower than threshold %d", base.Name, int(base.Health), *healthThreshold), b.API)
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
