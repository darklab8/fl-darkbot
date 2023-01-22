package cmdgroup

import (
	"darkbot/configurator"
	"darkbot/consoler/helper"

	"github.com/spf13/cobra"
)

type CmdGroup struct {
	ParentCmd   *cobra.Command
	CurrentCmd  *cobra.Command
	CfgTags     configurator.IConfiguratorTags
	ChannelInfo helper.ChannelInfo
}

type CmdGroupProps struct {
	Command   string
	ShortDesc string
}

func New(
	rootCmdPrefix *cobra.Command,
	channelInfo helper.ChannelInfo,
	props CmdGroupProps,
) CmdGroup {
	result := CmdGroup{
		ParentCmd:   rootCmdPrefix,
		CfgTags:     configurator.ConfiguratorBase{Configurator: configurator.NewConfigurator()},
		ChannelInfo: channelInfo,
	}
	result.CurrentCmd = &cobra.Command{
		Use:   props.Command,
		Short: props.ShortDesc,
	}
	result.ParentCmd.AddCommand(result.CurrentCmd)

	return result
}
