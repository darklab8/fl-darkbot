package viewer

import (
	"darkbot/scrappy/base"
	"darkbot/utils"
	_ "embed"
	"fmt"
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

type TemplateView struct {
	MessageID string
	Content   string
	Header    string
	ViewConfig
}

type BaseView struct {
	TemplateView
}

func CheckTooLongMsgErr(err error, api ViewConfig, header string, action MsgAction, MessageID string) {
	if err == nil {
		return
	}

	if !strings.Contains(err.Error(), "BASE_TYPE_MAX_LENGTH") &&
		!strings.Contains(err.Error(), "or fewer in length") {
		return
	}

	msg := fmt.Sprintf("%s, %s, %s", header, time.Now(), err)

	switch action {
	case ActSend:
		api.discorder.SengMessage(api.channelID, msg)
	case ActEdit:
		api.discorder.EditMessage(api.channelID, MessageID, msg)
	}

}

func ChannelCheckWarn(err error, channelID string, msg string) {
	utils.CheckWarn(err, "channelID=", channelID, msg)
}

func (v *TemplateView) Send() {
	if v.Content == "" && v.MessageID != "" {
		v.discorder.DeleteMessage(v.channelID, v.MessageID)
	}

	if v.Content == "" {
		return
	}

	var err error
	if v.MessageID == "" {
		err = v.discorder.SengMessage(v.channelID, v.Content)
		ChannelCheckWarn(err, v.channelID, "unable to send msg")
		CheckTooLongMsgErr(err, v.ViewConfig, v.Header, ActSend, "")

	} else {
		err = v.discorder.EditMessage(v.channelID, v.MessageID, v.Content)
		ChannelCheckWarn(err, v.channelID, "unable to edit msg")
		CheckTooLongMsgErr(err, v.ViewConfig, v.Header, ActEdit, v.MessageID)
	}
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
		Header:      b.Header,
		LastUpdated: time.Now().String(),
	}

	record, err := b.scrappy.BaseStorage.GetLatestRecord()
	if err != nil {
		return
	}
	sort.Slice(record.List, func(i, j int) bool {
		return record.List[i].Name < record.List[j].Name
	})

	tags := b.bases.TagsList(b.channelID)

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
