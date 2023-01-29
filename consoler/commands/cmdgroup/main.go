package cmdgroup

import (
	"darkbot/configurator"
	"darkbot/consoler/helper"

	"github.com/spf13/cobra"
)

type CmdGroup struct {
	ParentCmd    *cobra.Command
	CurrentCmd   *cobra.Command
	Configurator configurator.Configurator
	ChannelInfo  helper.ChannelInfo
}

func (c CmdGroup) GetChild(
	parentCmd *cobra.Command,
	props CmdGroupProps,
) *CmdGroup {
	c.ParentCmd = parentCmd
	c.setProps(props)
	return &c
}

func (c *CmdGroup) setProps(
	props CmdGroupProps,
) {

	c.CurrentCmd = &cobra.Command{
		Use:   props.Command,
		Short: props.ShortDesc,
	}
	c.ParentCmd.AddCommand(c.CurrentCmd)
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
		ParentCmd:    rootCmdPrefix,
		Configurator: configurator.NewConfigurator(channelInfo.Dbpath),
		ChannelInfo:  channelInfo,
	}
	result.setProps(props)

	return result
}
