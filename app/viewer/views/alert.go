package views

import (
	"darkbot/app/settings/logus"
	"darkbot/app/settings/types"
	"darkbot/app/settings/utils"
	"darkbot/app/viewer/apis"
	_ "embed"
	"fmt"
	"text/template"
	"time"
)

//go:embed alert_template.md
var alertMarkup string
var alertTemplate *template.Template

func init() {
	alertTemplate = utils.TmpInit(alertMarkup)
}

type TemplateAlertInput struct {
	Header      string
	LastUpdated string
	PingMessage types.PingMessage
	Msg         string
}

func RenderAlertTemplate(Header string, ChannelID types.DiscordChannelID, Msg string, api *apis.API) string {

	pingMessage, err := api.Alerts.PingMessage.Status(ChannelID)
	logus.Debug("RenderAlertTemplate.PingMessage.Status", logus.OptError(err), logus.PingMessage(pingMessage))
	if err != nil {
		ownerID, err := api.Discorder.GetOwnerID(ChannelID)
		if logus.CheckWarn(err, "unable to acquire Discorder Channel Owner", logus.ChannelID(ChannelID)) {
			ownerID = "TestOwnerID"
		}
		pingMessage = types.PingMessage(fmt.Sprintf("<@%s>", ownerID))
	}

	input := TemplateAlertInput{
		Header:      Header,
		LastUpdated: time.Now().String(),
		PingMessage: pingMessage,
		Msg:         Msg,
	}
	return utils.TmpRender(alertTemplate, input)
}
