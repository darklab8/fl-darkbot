package configurator

import (
	"darkbot/app/discorder"
	"darkbot/app/settings/logus"
	"darkbot/app/settings/types"
	"fmt"
)

func GetPingingMessage(ChannelID types.DiscordChannelID, configurator *Configurators, Discorder *discorder.Discorder) types.PingMessage {
	pingMessage, err := configurator.Alerts.PingMessage.Status(ChannelID)
	logus.Debug("RenderAlertTemplate.PingMessage.Status", logus.OptError(err), logus.PingMessage(pingMessage))
	if err != nil {
		ownerID, err := Discorder.GetOwnerID(ChannelID)
		if logus.CheckWarn(err, "unable to acquire Discorder Channel Owner", logus.ChannelID(ChannelID)) {
			ownerID = "TestOwnerID"
		}
		pingMessage = types.PingMessage(fmt.Sprintf("<@%s>", ownerID))
	}
	return pingMessage
}
