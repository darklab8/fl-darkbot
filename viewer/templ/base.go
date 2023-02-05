package templ

import (
	"darkbot/dtypes"
	"darkbot/scrappy/base"
	"darkbot/utils"
	"darkbot/viewer/apis"
	_ "embed"
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
	main TemplateShared
}

func NewTemplateBase(channelID string, dbpath dtypes.Dbpath) *TemplateBase {
	base := TemplateBase{}
	base.main.API = apis.NewAPI(channelID, dbpath)
	base.main.Header = "#darkbot-base-view"
	return &base
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

func (b *TemplateBase) Render() {
	input := TemplateRendererBaseInput{
		Header:      b.main.Header,
		LastUpdated: time.Now().String(),
	}

	record, err := b.main.Scrappy.BaseStorage.GetLatestRecord()
	if err != nil {
		return
	}
	sort.Slice(record.List, func(i, j int) bool {
		return record.List[i].Name < record.List[j].Name
	})

	tags, _ := b.main.Bases.TagsList(b.main.ChannelID)

	for _, base := range record.List {

		if !BaseContainsTag(*base, tags) {
			continue
		}

		input.Bases = append(input.Bases, *base)
	}

	if len(input.Bases) == 0 {
		return
	}

	b.main.Content = utils.TmpRender(baseTemplate, input)
}

func (t *TemplateBase) Send() {
	t.main.Send()
}

func (t *TemplateBase) MatchMessageID(messageID string) bool {

	if messageID == t.main.MessageID {
		return true
	}
	return false
}

func (t *TemplateBase) DiscoverMessageID(content string, msgID string) {
	if strings.Contains(content, t.main.Header) {
		t.main.MessageID = msgID
	}
}
