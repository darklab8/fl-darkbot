package commands

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/darklab8/fl-darkbot/app/configurator"
	"github.com/darklab8/fl-darkbot/app/consoler/commands/cmdgroup"
	"github.com/darklab8/fl-darkbot/app/consoler/printer"
	"github.com/darklab8/fl-darkbot/app/settings/logus"
	"github.com/darklab8/go-utils/typelog"
	"github.com/spf13/cobra"
)

type BasePriorityCommands[T configurator.BasePriorityType] struct {
	*cmdgroup.CmdGroup
	cfg      configurator.IConfiguratorBasePriority[T]
	channels configurator.ConfiguratorChannel
}

func NewBasePriorityCommands[T configurator.BasePriorityType](
	cmd *cmdgroup.CmdGroup,
	cfgGoodThresholds configurator.IConfiguratorBasePriority[T],
	channels configurator.ConfiguratorChannel,
) *BasePriorityCommands[T] {
	t := &BasePriorityCommands[T]{CmdGroup: cmd, cfg: cfgGoodThresholds, channels: channels}
	t.Create()
	t.Remove()
	t.Clear()
	t.List()
	return t
}

func (t *BasePriorityCommands[T]) Create() {
	command := &cobra.Command{
		Use:   "add",
		Short: "Add base priority config",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			logus.Log.Debug("BasePriorityCommands.Create.consoler running with args=", logus.Args(args))
			if !CheckCommandAllowedToRun(cmd, t.channels, t.GetChannelID()) {
				return
			}

			if len(args) != 2 {
				printer.Println(cmd, fmt.Sprintln("expected to get 2 args base_nickname(string) and priority(number). got=", len(args)))
				return
			}

			base_nickname := args[0]
			priority_str := args[1]

			priority, err := strconv.Atoi(priority_str)
			if err != nil {
				printer.Println(cmd, "ERR failed to convert second argument to number err="+err.Error())
				return
			}

			err = t.cfg.Add(t.GetChannelID(), base_nickname, priority)
			if err != nil {
				printer.Println(cmd, "ERR failed to add to table msg="+err.Error())
				return
			}
			logus.Log.Debug("CreateTagAdd", logus.Args(args))

			printer.Println(cmd, fmt.Sprintf("OK pob priority is added.\n```\nbase_nickname=%s, priority=%d\n```\n", base_nickname, priority))
		},
	}
	t.CurrentCmd.AddCommand(command)
}

func (t *BasePriorityCommands[T]) Remove() {
	command := &cobra.Command{
		Use:   "remove",
		Short: "Remove pob priority configs",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			logus.Log.Debug("BasePriorityCommands.Remove.consoler running with args=", logus.Args(args))
			if !CheckCommandAllowedToRun(cmd, t.channels, t.GetChannelID()) {
				return
			}

			if len(args) == 0 {
				printer.Println(cmd, "No pob priority configs found to remove. Expected at least one base_nickname")
				return
			}

			for _, base_nickname := range args {
				err := t.cfg.Remove(t.GetChannelID(), base_nickname)
				if err != nil {
					if _, ok := err.(configurator.ErrorZeroAffectedRows); ok {
						printer.Println(cmd, "ERR removed nothing, because inserted value did not match anything present in the pob good alert configs")
					} else {
						printer.Println(cmd, "ERR ="+err.Error())
					}
					return
				}
			}

			printer.Println(cmd, "OK pob priority configs are removed: "+strings.Join(args, " "))
			logus.Log.Debug("executed BasePriorityCommands Remove with args", logus.Args(args))
		},
	}
	t.CurrentCmd.AddCommand(command)

}

func (t *BasePriorityCommands[T]) Clear() {
	command := &cobra.Command{
		Use:   "clear",
		Short: "Clear pob priority configs",
		Run: func(cmd *cobra.Command, args []string) {
			logus.Log.Debug("BasePriorityCommands.consoler running with args=", logus.Args(args))
			if !CheckCommandAllowedToRun(cmd, t.channels, t.GetChannelID()) {
				return
			}

			err := t.cfg.Clear(t.GetChannelID())
			if err != nil {
				if _, ok := err.(configurator.ErrorZeroAffectedRows); ok {
					printer.Println(cmd, "ERR pob priority configs list is already empty. nothing to clear.")
				} else {
					printer.Println(cmd, "ERR ="+err.Error())
				}
				return
			}

			printer.Println(cmd, "OK pob priority configs are cleared")
		},
	}
	t.CurrentCmd.AddCommand(command)
}

func (t *BasePriorityCommands[T]) List() {
	command := &cobra.Command{
		Use:   "list",
		Short: "List pob priority configs",
		Run: func(cmd *cobra.Command, args []string) {
			logus.Log.Debug("BasePriorityCommands.List.consoler running with args=", logus.Args(args))
			if !CheckCommandAllowedToRun(cmd, t.channels, t.GetChannelID()) {
				return
			}

			tags, cfgErr := t.cfg.Get(t.GetChannelID())
			err := cfgErr
			if err != nil {
				if _, ok := err.(configurator.ErrorZeroAffectedRows); ok {
					printer.Println(cmd, "OK pob priority config list is empty")
				} else {
					printer.Println(cmd, "ERR ="+err.Error())
				}
				return
			}

			logus.Log.Debug("BasePriorityCommands.List continuied", typelog.Any("tags", tags))
			var sb strings.Builder
			for tag, priority := range tags {
				sb.WriteString(fmt.Sprintf("\"{%s:%d}\"", tag, priority))

				if priority != len(tags)-1 {
					sb.WriteString(", ")
				}
			}
			printer.Println(cmd, "OK pob priority configs are listed")
			printer.Println(cmd, sb.String())
		},
	}
	t.CurrentCmd.AddCommand(command)
}
