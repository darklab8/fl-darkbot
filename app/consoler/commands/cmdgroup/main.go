package cmdgroup

import (
	"github.com/darklab/fl-darkbot/app/configurator"
	"github.com/darklab/fl-darkbot/app/consoler/consoler_types"

	"github.com/spf13/cobra"
)

type CmdGroup struct {
	ParentCmd    *cobra.Command
	CurrentCmd   *cobra.Command
	Configurator *configurator.Configurator
	*consoler_types.ChannelParams
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

func NewCmdGroup(
	rootCmdPrefix *cobra.Command,
	channelParams *consoler_types.ChannelParams,
	configur *configurator.Configurator,
	command Command, short_desc ShortDesc,
) CmdGroup {
	result := CmdGroup{
		ParentCmd:     rootCmdPrefix,
		Configurator:  configur,
		ChannelParams: channelParams,
	}
	result.setProps(command, short_desc)

	return result
}
