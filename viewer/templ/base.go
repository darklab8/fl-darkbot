package templ

import (
	"darkbot/scrappy/base"
	"darkbot/utils"
	"darkbot/viewer/apis"
	_ "embed"
	"sort"
	"strings"
	"text/template"
	"time"
)

const (
	BaseViewHeader = "#darkbot-base-view"
)

//go:embed base_template.md
var baseMarkup string
var baseTemplate *template.Template

func init() {
	baseTemplate = utils.TmpInit(baseMarkup)
}

// Base

type TemplateBase struct {
	TemplateShared
}

func NewTemplateBase(channelID string) *TemplateBase {
	base := TemplateBase{}
	base.API = apis.NewAPI(channelID)
	base.Header = BaseViewHeader
	return &base
}

type BaseInput struct {
	Header      string
	LastUpdated string
	Bases       []base.Base
}

func BaseContainsTag(bas *base.Base, tags []string) bool {
	for _, tag := range tags {
		if strings.Contains(bas.Name, tag) {
			return true
		}
	}

	return false
}

func (b *TemplateBase) Render() {
	input := BaseInput{
		Header:      b.Header,
		LastUpdated: time.Now().String(),
	}

	record, err := b.Scrappy.BaseStorage.GetLatestRecord()
	if err != nil {
		return
	}
	sort.Slice(record.List, func(i, j int) bool {
		return record.List[i].Name < record.List[j].Name
	})

	tags := b.Bases.TagsList(b.ChannelID)

	for _, base := range record.List {

		if !BaseContainsTag(base, tags) {
			continue
		}

		input.Bases = append(input.Bases, *base)
	}

	if len(input.Bases) == 0 {
		return
	}

	b.Content = utils.TmpRender(baseTemplate, input)
}
