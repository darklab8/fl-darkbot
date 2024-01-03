package views

import (
	"darkbot/app/configurator"
	"darkbot/app/settings/types"
	"darkbot/app/viewer/apis"
	_ "embed"
	"text/template"

	"github.com/darklab8/darklab_goutils/goutils/utils"
	"github.com/darklab8/darklab_goutils/goutils/utils/utils_types"
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
