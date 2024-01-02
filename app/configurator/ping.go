package configurator

import (
	"darkbot/app/discorder"
	"darkbot/app/settings/darkbot_logus"
	"darkbot/app/settings/types"
	"fmt"

	"github.com/darklab8/darklab_goutils/goutils/logus"
)

func GetPingingMessage(ChannelID types.DiscordChannelID, configurator *Configurators, Discorder *discorder.Discorder) types.PingMessage {
	pingMessageStr, err := configurator.Alerts.PingMessage.Status(ChannelID)
	pingMessage := types.PingMessage(pingMessageStr)
	darkbot_logus.Log.Debug("RenderAlertTemplate.PingMessage.Status", logus.OptError(err), darkbot_logus.PingMessage(pingMessage))
	if err != nil {
		ownerID, err := Discorder.GetOwnerID(ChannelID)
		if darkbot_logus.Log.CheckWarn(err, "unable to acquire Discorder Channel Owner", darkbot_logus.ChannelID(ChannelID)) {
			ownerID = "TestOwnerID"
		}
		pingMessage = types.PingMessage(fmt.Sprintf("<@%s>", ownerID))
	}
	return pingMessage
}
