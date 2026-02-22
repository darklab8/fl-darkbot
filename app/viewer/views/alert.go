package views

import (
	_ "embed"
	"text/template"

	"github.com/darklab8/fl-darkbot/app/configurator"
	"github.com/darklab8/fl-darkbot/app/settings/types"
	"github.com/darklab8/fl-darkbot/app/viewer/apis"

	"github.com/darklab8/go-utils/utils"
	"github.com/darklab8/go-utils/utils/utils_types"
)

//go:embed alert_template.md
var alertMarkup utils_types.TemplateExpression
var alertTemplate *template.Template

func init() {
	alertTemplate = utils.TmpInit(alertMarkup)
}

type TemplateAlertInput struct {
	PingMessage types.PingMessage
	Msg         string
}

type Alert struct {
	PingMessage types.PingMessage
}
type AlertOpt func(a *Alert)

func WithAlertOverride(PingMessage string) AlertOpt {
	return func(a *Alert) {
		a.PingMessage = types.PingMessage(PingMessage)
	}
}

func RenderAlertTemplate(ChannelID types.DiscordChannelID, Msg string, api *apis.API, opts ...AlertOpt) types.ViewHeader {
	// pingMessage, err := api.Alerts.PingMessage.Status(ChannelID)
	// ownerID, err := api.Discorder.GetOwnerID(ChannelID)

	a := &Alert{
		PingMessage: configurator.GetPingingMessage(ChannelID, api.Configurators, api.Discorder),
	}
	for _, opt := range opts {
		opt(a)
	}

	input := TemplateAlertInput{
		PingMessage: a.PingMessage,
		Msg:         Msg,
	}
	return types.ViewHeader(utils.TmpRender(alertTemplate, input))
}
