package views

import (
	"darkbot/app/configurator"
	"darkbot/app/settings/types"
	"darkbot/app/settings/utils"
	"darkbot/app/viewer/apis"
	_ "embed"
	"text/template"
)

//go:embed alert_template.md
var alertMarkup string
var alertTemplate *template.Template

func init() {
	alertTemplate = utils.TmpInit(alertMarkup)
}

type TemplateAlertInput struct {
	PingMessage types.PingMessage
	Msg         string
}

func RenderAlertTemplate(ChannelID types.DiscordChannelID, Msg string, api *apis.API) types.ViewHeader {
	// pingMessage, err := api.Alerts.PingMessage.Status(ChannelID)
	// ownerID, err := api.Discorder.GetOwnerID(ChannelID)

	pingMessage := configurator.GetPingingMessage(ChannelID, api.Configurators, api.Discorder)
	input := TemplateAlertInput{
		PingMessage: pingMessage,
		Msg:         Msg,
	}
	return types.ViewHeader(utils.TmpRender(alertTemplate, input))
}
