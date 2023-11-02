package commands

import (
	"darkbot/app/configurator"
	"darkbot/app/consoler/commands/cmdgroup"
	"darkbot/app/consoler/printer"
	"darkbot/app/settings/logus"
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
			logus.Debug("CreateSetAlertCmd.consoler running with args=", logus.Args(args))
			printer.Println(cmd, "Attempting to parse input into integer number")
			rawInteger := args[0]
			integer, err := strconv.Atoi(rawInteger)
			if logus.CheckWarn(err, "Atoi result with warning", logus.OptError(err)) {
				printer.Println(cmd, "failed to parse value to integer. Value="+rawInteger)
			}

			printer.Println(cmd, "Parsed integer = "+strconv.Itoa(integer))
			err = t.cfgTags.Set(t.GetChannelID(), integer)
			if err != nil {
				printer.Println(cmd, "ERR msg="+err.Error())
				return
			}
			logus.Debug("checking args again?", logus.Args(args))

			printer.Println(cmd, "OK alert threshold is set")
		},
	}
	t.CurrentCmd.AddCommand(command)
}

func (t *alertThresholdCommands[T]) CreateUnsetCmd() {
	command := &cobra.Command{
		Use:   "unset",
		Short: "Unsert alert / Clear alert",
		Run: func(cmd *cobra.Command, args []string) {
			logus.Debug("CreateUnsetCmd.consoler running with args=", logus.Args(args))
			err := t.cfgTags.Unset(t.GetChannelID())
			if err != nil {
				printer.Println(cmd, "ERR msg="+err.Error())
				return
			}
			printer.Println(cmd, "OK Alert is unset")
		},
	}
	t.CurrentCmd.AddCommand(command)
}

func (t *alertThresholdCommands[T]) CreateStatusCmd() {
	command := &cobra.Command{
		Use:   "status",
		Short: "Status of alert",
		Run: func(cmd *cobra.Command, args []string) {
			logus.Debug("CreateStatusCmd.consoler running with args=", logus.Args(args))
			integer, err := t.cfgTags.Status(t.GetChannelID())
			if err != nil {
				errMsg := err.Error()

				if strings.Contains(errMsg, "record not found") {
					printer.Println(cmd, "OK status of alert is disabled")
					return
				} else {
					printer.Println(cmd, "ERR ="+errMsg)
					return
				}
			}

			printer.Println(cmd, "OK status of alert threshold = "+strconv.Itoa(integer))
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
			logus.Debug("CreateEnableCmd.consoler running with args=", logus.Args(args))
			err := t.cfgTags.Enable(t.GetChannelID())
			if err != nil {
				printer.Println(cmd, "ERR msg="+err.Error())
				return
			}
			logus.Debug("Create Enable is finished", logus.Args(args))

			printer.Println(cmd, "OK alert is enabled")
		},
	}
	t.CurrentCmd.AddCommand(command)
}

func (t *AlertBoolCommands[T]) CreateDisableCmd() {
	command := &cobra.Command{
		Use:   "disable",
		Short: "Disable alert / Clear alert",
		Run: func(cmd *cobra.Command, args []string) {
			logus.Debug("CreateDisableCmd.consoler running with args=", logus.Args(args))
			err := t.cfgTags.Disable(t.GetChannelID())
			if err != nil {
				printer.Println(cmd, "ERR msg="+err.Error())
				return
			}
			printer.Println(cmd, "OK Alert is disabled")
		},
	}
	t.CurrentCmd.AddCommand(command)
}

func (t *AlertBoolCommands[T]) CreateStatusCmd() {
	command := &cobra.Command{
		Use:   "status",
		Short: "Status of alert",
		Run: func(cmd *cobra.Command, args []string) {
			logus.Debug("CreateStatusCmd.consoler running with args=", logus.Args(args))
			_, err := t.cfgTags.Status(t.GetChannelID())
			if err != nil {
				errMsg := err.Error()

				if strings.Contains(errMsg, "record not found") {
					printer.Println(cmd, "OK status of alert is disabled")
					return
				} else {
					printer.Println(cmd, "ERR ="+errMsg)
					return
				}
			}

			printer.Println(cmd, "OK status of alert is enabled")
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
			logus.Debug("CreateSetAlertCmd.consoler running with args=", logus.Args(args))
			str := args[0]
			err := t.cfgTags.Set(t.GetChannelID(), str)
			if err != nil {
				printer.Println(cmd, "ERR msg="+err.Error())
				return
			}
			logus.Debug("finished CreateSetCmd", logus.Args(args))

			printer.Println(cmd, "OK value is set")
		},
	}
	t.CurrentCmd.AddCommand(command)
}

func (t *AlertSetStringCommand[T]) CreateUnsetCmd() {
	command := &cobra.Command{
		Use:   "unset",
		Short: "Unsert / Clear ",
		Run: func(cmd *cobra.Command, args []string) {
			logus.Debug("CreateUnsetCmd.consoler running with args=", logus.Args(args))
			err := t.cfgTags.Unset(t.GetChannelID())
			if err != nil {
				printer.Println(cmd, "ERR msg="+err.Error())
				return
			}
			printer.Println(cmd, "OK value is unset")
		},
	}
	t.CurrentCmd.AddCommand(command)
}

func (t *AlertSetStringCommand[T]) CreateStatusCmd() {
	command := &cobra.Command{
		Use:   "status",
		Short: "Status",
		Run: func(cmd *cobra.Command, args []string) {
			logus.Debug("CreateStatusCmd.consoler running with args=", logus.Args(args))
			str, err := t.cfgTags.Status(t.GetChannelID())
			if err != nil {
				errMsg := err.Error()

				if strings.Contains(errMsg, "record not found") {
					printer.Println(cmd, "OK status of alert is disabled")
					return
				} else {
					printer.Println(cmd, "ERR ="+errMsg)
					return
				}
			}

			printer.Println(cmd, "OK value is = "+string(str))
		},
	}
	t.CurrentCmd.AddCommand(command)
}
