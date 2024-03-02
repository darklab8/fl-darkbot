package commands

import (
	"fmt"

	"github.com/darklab8/fl-darkbot/app/configurator"
	"github.com/darklab8/fl-darkbot/app/consoler/printer"
	"github.com/darklab8/fl-darkbot/app/settings"
	"github.com/darklab8/fl-darkbot/app/settings/types"

	"github.com/spf13/cobra"
)

func CheckCommandAllowedToRun(cmd *cobra.Command, channels configurator.ConfiguratorChannel, channelID types.DiscordChannelID) bool {
	isChannelEnabled, err := channels.IsEnabled(channelID)

	if err != nil {
		printer.Println(cmd, "ERR ="+err.Error())
		return false
	}

	if !isChannelEnabled {
		printer.Println(cmd, fmt.Sprintf("darkbot is not connected to this channel. Run `%s connect`", settings.Config.ConsolerPrefix))
		return false
	}

	return true
}
