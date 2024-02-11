package views

import (
	_ "embed"
	"text/template"

	"github.com/darklab/fl-darkbot/app/configurator"
	"github.com/darklab/fl-darkbot/app/settings/types"
	"github.com/darklab/fl-darkbot/app/viewer/apis"

	"github.com/darklab8/go-utils/goutils/utils"
	"github.com/darklab8/go-utils/goutils/utils/utils_types"
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
