package configurator

import (
	"fmt"

	"github.com/darklab8/fl-darkbot/app/discorder"
	"github.com/darklab8/fl-darkbot/app/settings/logus"
	"github.com/darklab8/fl-darkbot/app/settings/types"
	"github.com/darklab8/go-utils/typelog"
)

func GetPingingMessage(ChannelID types.DiscordChannelID, configurator *Configurators, Discorder *discorder.Discorder) types.PingMessage {
	pingMessageStr, err := configurator.Alerts.PingMessage.Status(ChannelID)
	pingMessage := types.PingMessage(pingMessageStr)
	logus.Log.Debug("RenderAlertTemplate.PingMessage.Status", typelog.OptError(err), logus.PingMessage(pingMessage))
	if err != nil {
		ownerID, err := Discorder.GetOwnerID(ChannelID)
		if logus.Log.CheckWarn(err, "unable to acquire Discorder Channel Owner", logus.ChannelID(ChannelID)) {
			ownerID = "TestOwnerID"
		}
		pingMessage = types.PingMessage(fmt.Sprintf("<@%s>", ownerID))
	}
	return pingMessage
}
