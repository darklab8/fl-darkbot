package commands

import (
	"darkbot/app/configurator"
	"darkbot/app/consoler/printer"
	"darkbot/app/settings"
	"darkbot/app/settings/types"
	"fmt"

	"github.com/spf13/cobra"
)

func CheckCommandAllowedToRun(cmd *cobra.Command, channels configurator.ConfiguratorChannel, channelID types.DiscordChannelID) bool {
	isChannelEnabled, err := channels.IsEnabled(channelID)

	if err != nil {
		printer.Println(cmd, "ERR ="+err.Error())
		return false
	}

	if !isChannelEnabled {
		printer.Println(cmd, fmt.Sprintf("Darkbot is not connected to this channel. Run `%s connect`", settings.Config.ConsolerPrefix))
		return false
	}

	return true
}
