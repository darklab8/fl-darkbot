package templ

import (
	"darkbot/settings/utils"
	"darkbot/settings/utils/logger"
	"darkbot/viewer/apis"
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
	PingMessage string
	Msg         string
}

func RenderAlertTemplate(Header string, ChannelID string, Msg string, api apis.API) string {

	pingMessage, err := api.Alerts.PingMessage.Status(ChannelID)
	logger.Debug("RenderAlertTemplate.PingMessage.Status.err=", err, " pingMessage=", pingMessage)
	if err.GetError() != nil {
		ownerID, err := api.Discorder.GetOwnerID(ChannelID)
		if err != nil {
			ownerID = "TestOwnerID"
		}
		pingMessage = fmt.Sprintf("<@%s>", ownerID)
	}

	input := TemplateAlertInput{
		Header:      Header,
		LastUpdated: time.Now().String(),
		PingMessage: pingMessage,
		Msg:         Msg,
	}
	return utils.TmpRender(alertTemplate, input)
}
