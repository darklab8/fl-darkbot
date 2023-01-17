package viewer

import (
	"darkbot/scrappy/base"
	"darkbot/utils"
	_ "embed"
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

type BaseView struct {
	MessageID string
	Content   string
	ViewConfig
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

func (b *BaseView) Render() {
	input := BaseInput{
		Header:      BaseViewHeader,
		LastUpdated: time.Now().String(),
	}

	record, err := b.scrappy.BaseStorage.GetLatestRecord()
	if err != nil {
		return
	}
	tags := b.bases.TagsList(b.channelID)

	for _, base := range record.List {

		if !BaseContainsTag(base, tags) {
			continue
		}

		input.Bases = append(input.Bases, *base)
	}

	b.Content = utils.TmpRender(baseTemplate, input)
}
