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
	command Command, short_desc ShortDesc,
) *CmdGroup {
	c.ParentCmd = parentCmd
	c.setProps(command, short_desc)
	return &c
}

func (c *CmdGroup) setProps(
	command Command, short_desc ShortDesc,
) {

	c.CurrentCmd = &cobra.Command{
		Use:   string(command),
		Short: string(short_desc),
	}
	c.ParentCmd.AddCommand(c.CurrentCmd)
}

type Command string
type ShortDesc string

func New(
	rootCmdPrefix *cobra.Command,
	channelInfo helper.ChannelInfo,
	command Command, short_desc ShortDesc,
) CmdGroup {
	result := CmdGroup{
		ParentCmd:    rootCmdPrefix,
		Configurator: configurator.NewConfigurator(channelInfo.Dbpath),
		ChannelInfo:  channelInfo,
	}
	result.setProps(command, short_desc)

	return result
}
