package commands

import (
	"darkbot/configurator"
	"darkbot/consoler/commands/cmdgroup"
	"darkbot/consoler/helper"
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

type alertThresholdCommands[T configurator.AlertThresholdType] struct {
	*cmdgroup.CmdGroup
	cfgTags configurator.IConfiguratorAlertThreshold[T]
}

func NewAlertThresholdCommands[T configurator.AlertThresholdType](
	cmdGroup *cmdgroup.CmdGroup,
	cfgTags configurator.IConfiguratorAlertThreshold[T],
) *alertThresholdCommands[T] {
	t := &alertThresholdCommands[T]{CmdGroup: cmdGroup, cfgTags: cfgTags}
	t.CreateSetAlertCmd()
	t.CreateUnsetCmd()
	t.CreateStatusCmd()
	return t
}

func (t *alertThresholdCommands[T]) CreateSetAlertCmd() {
	command := &cobra.Command{
		Use:   "set",
		Short: "Set alert (Works as set {number})",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("CreateSetAlertCmd.consoler running with args=", args)
			helper.Printer{Cmd: cmd}.Println("Attempting to parse input into integer number")
			rawInteger := args[0]
			integer, err := strconv.Atoi(rawInteger)
			helper.Printer{Cmd: cmd}.Println("Parsed integer = " + strconv.Itoa(integer))
			err = t.cfgTags.Set(t.ChannelInfo.ChannelID, integer).GetError()
			if err != nil {
				cmd.OutOrStdout().Write([]byte("ERR msg=" + err.Error()))
				return
			}
			fmt.Println(len(args))

			helper.Printer{Cmd: cmd}.Println("OK alert threshold is set")
		},
	}
	t.CurrentCmd.AddCommand(command)
}

func (t *alertThresholdCommands[T]) CreateUnsetCmd() {
	command := &cobra.Command{
		Use:   "unset",
		Short: "Unsert alert / Clear alert",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("CreateUnsetCmd.consoler running with args=", args)
			err := t.cfgTags.Unset(t.ChannelInfo.ChannelID).GetError()
			if err != nil {
				cmd.OutOrStdout().Write([]byte("ERR msg=" + err.Error()))
				return
			}
			helper.Printer{Cmd: cmd}.Println("OK Alert is unset")
		},
	}
	t.CurrentCmd.AddCommand(command)
}

func (t *alertThresholdCommands[T]) CreateStatusCmd() {
	command := &cobra.Command{
		Use:   "status",
		Short: "Status of alert",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("CreateStatusCmd.consoler running with args=", args)
			integer, err := t.cfgTags.Status(t.ChannelInfo.ChannelID)
			if err.GetError() != nil {
				errMsg := err.GetError().Error()

				if strings.Contains(errMsg, "record not found") {
					helper.Printer{Cmd: cmd}.Println("OK status of alert is disabled")
					return
				} else {
					cmd.OutOrStdout().Write([]byte("ERR =" + errMsg))
					return
				}
			}

			helper.Printer{Cmd: cmd}.Println("OK status of alert threshold = " + strconv.Itoa(*integer))
		},
	}
	t.CurrentCmd.AddCommand(command)
}

type AlertBoolCommands[T configurator.AlertBoolType] struct {
	*cmdgroup.CmdGroup
	cfgTags configurator.IConfiguratorAlertBool[T]
}

func NewAlertBoolCommands[T configurator.AlertBoolType](
	cmdGroup *cmdgroup.CmdGroup,
	cfgTags configurator.IConfiguratorAlertBool[T],
) *AlertBoolCommands[T] {
	t := &AlertBoolCommands[T]{CmdGroup: cmdGroup, cfgTags: cfgTags}
	t.CreateEnableCmd()
	t.CreateDisableCmd()
	t.CreateStatusCmd()
	return t
}

func (t *AlertBoolCommands[T]) CreateEnableCmd() {
	command := &cobra.Command{
		Use:   "enable",
		Short: "Enable alert",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("CreateEnableCmd.consoler running with args=", args)
			err := t.cfgTags.Enable(t.ChannelInfo.ChannelID).GetError()
			if err != nil {
				cmd.OutOrStdout().Write([]byte("ERR msg=" + err.Error()))
				return
			}
			fmt.Println(len(args))

			helper.Printer{Cmd: cmd}.Println("OK alert is enabled")
		},
	}
	t.CurrentCmd.AddCommand(command)
}

func (t *AlertBoolCommands[T]) CreateDisableCmd() {
	command := &cobra.Command{
		Use:   "disable",
		Short: "Disable alert / Clear alert",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("CreateDisableCmd.consoler running with args=", args)
			err := t.cfgTags.Disable(t.ChannelInfo.ChannelID).GetError()
			if err != nil {
				cmd.OutOrStdout().Write([]byte("ERR msg=" + err.Error()))
				return
			}
			helper.Printer{Cmd: cmd}.Println("OK Alert is disabled")
		},
	}
	t.CurrentCmd.AddCommand(command)
}

func (t *AlertBoolCommands[T]) CreateStatusCmd() {
	command := &cobra.Command{
		Use:   "status",
		Short: "Status of alert",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("CreateStatusCmd.consoler running with args=", args)
			_, err := t.cfgTags.Status(t.ChannelInfo.ChannelID)
			if err.GetError() != nil {
				errMsg := err.GetError().Error()

				if strings.Contains(errMsg, "record not found") {
					helper.Printer{Cmd: cmd}.Println("OK status of alert is disabled")
					return
				} else {
					cmd.OutOrStdout().Write([]byte("ERR =" + errMsg))
					return
				}
			}

			helper.Printer{Cmd: cmd}.Println("OK status of alert is enabled")
		},
	}
	t.CurrentCmd.AddCommand(command)
}

///////// string //////////////

type AlertSetStringCommand[T configurator.AlertStringType] struct {
	*cmdgroup.CmdGroup
	cfgTags configurator.IConfiguratorAlertString[T]
}

func NewAlertSetStringCommand[T configurator.AlertStringType](
	cmdGroup *cmdgroup.CmdGroup,
	cfgTags configurator.IConfiguratorAlertString[T],
) *AlertSetStringCommand[T] {
	t := &AlertSetStringCommand[T]{CmdGroup: cmdGroup, cfgTags: cfgTags}
	t.CreateSetCmd()
	t.CreateUnsetCmd()
	t.CreateStatusCmd()
	return t
}

func (t *AlertSetStringCommand[T]) CreateSetCmd() {
	command := &cobra.Command{
		Use:   "set",
		Short: "Set Value (provide 'set StringValue')",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("CreateSetAlertCmd.consoler running with args=", args)
			str := args[0]
			err := t.cfgTags.Set(t.ChannelInfo.ChannelID, str).GetError()
			if err != nil {
				cmd.OutOrStdout().Write([]byte("ERR msg=" + err.Error()))
				return
			}
			fmt.Println(len(args))

			helper.Printer{Cmd: cmd}.Println("OK value is set")
		},
	}
	t.CurrentCmd.AddCommand(command)
}

func (t *AlertSetStringCommand[T]) CreateUnsetCmd() {
	command := &cobra.Command{
		Use:   "unset",
		Short: "Unsert / Clear ",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("CreateUnsetCmd.consoler running with args=", args)
			err := t.cfgTags.Unset(t.ChannelInfo.ChannelID).GetError()
			if err != nil {
				cmd.OutOrStdout().Write([]byte("ERR msg=" + err.Error()))
				return
			}
			helper.Printer{Cmd: cmd}.Println("OK value is unset")
		},
	}
	t.CurrentCmd.AddCommand(command)
}

func (t *AlertSetStringCommand[T]) CreateStatusCmd() {
	command := &cobra.Command{
		Use:   "status",
		Short: "Status",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("CreateStatusCmd.consoler running with args=", args)
			str, err := t.cfgTags.Status(t.ChannelInfo.ChannelID)
			if err.GetError() != nil {
				errMsg := err.GetError().Error()

				if strings.Contains(errMsg, "record not found") {
					helper.Printer{Cmd: cmd}.Println("OK status of alert is disabled")
					return
				} else {
					cmd.OutOrStdout().Write([]byte("ERR =" + errMsg))
					return
				}
			}

			helper.Printer{Cmd: cmd}.Println("OK value is = " + string(str))
		},
	}
	t.CurrentCmd.AddCommand(command)
}
