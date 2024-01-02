package commands

import (
	"darkbot/app/configurator"
	"darkbot/app/consoler/commands/cmdgroup"
	"darkbot/app/consoler/printer"
	"darkbot/app/settings/darkbot_logus"
	"darkbot/app/settings/types"
	"darkbot/app/settings/utils"
	"strconv"
	"strings"

	"github.com/darklab8/darklab_goutils/goutils/utils_logus"
	"github.com/spf13/cobra"
)

type alertThresholdCommands[T configurator.AlertThresholdType] struct {
	*cmdgroup.CmdGroup
	cfgTags  configurator.IConfiguratorAlertThreshold[T]
	channels configurator.ConfiguratorChannel
}

func NewAlertThresholdCommands[T configurator.AlertThresholdType](
	cmdGroup *cmdgroup.CmdGroup,
	cfgTags configurator.IConfiguratorAlertThreshold[T],
	channels configurator.ConfiguratorChannel,
) *alertThresholdCommands[T] {
	t := &alertThresholdCommands[T]{CmdGroup: cmdGroup, cfgTags: cfgTags, channels: channels}
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
			darkbot_logus.Log.Debug("CreateSetAlertCmd.consoler running with args=", darkbot_logus.Args(args))
			if !CheckCommandAllowedToRun(cmd, t.channels, t.GetChannelID()) {
				return
			}

			printer.Println(cmd, "Attempting to parse input into integer number")
			rawInteger := args[0]
			integer, err := strconv.Atoi(rawInteger)
			if darkbot_logus.Log.CheckWarn(err, "Atoi result with warning", utils_logus.OptError(err)) {
				printer.Println(cmd, "failed to parse value to integer. Value="+rawInteger)
			}

			printer.Println(cmd, "Parsed integer = "+strconv.Itoa(integer))
			err = t.cfgTags.Set(t.GetChannelID(), integer)
			if err != nil {
				printer.Println(cmd, "ERR msg="+err.Error())
				return
			}
			darkbot_logus.Log.Debug("checking args again?", darkbot_logus.Args(args))

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
			darkbot_logus.Log.Debug("CreateUnsetCmd.consoler running with args=", darkbot_logus.Args(args))
			if !CheckCommandAllowedToRun(cmd, t.channels, t.GetChannelID()) {
				return
			}

			err := t.cfgTags.Unset(t.GetChannelID())
			if err != nil {
				if _, ok := err.(configurator.ErrorZeroAffectedRows); ok {
					printer.Println(cmd, "ERR it is already unset")
				} else {
					printer.Println(cmd, "ERR ="+err.Error())
				}
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
			darkbot_logus.Log.Debug("CreateStatusCmd.consoler running with args=", darkbot_logus.Args(args))
			if !CheckCommandAllowedToRun(cmd, t.channels, t.GetChannelID()) {
				return
			}

			integer, err := t.cfgTags.Status(t.GetChannelID())
			if err != nil {
				errMsg := err.Error()

				if _, ok := err.(configurator.ErrorZeroAffectedRows); ok {
					printer.Println(cmd, "OK state is disabled")
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
	cfgTags  configurator.IConfiguratorAlertBool[T]
	channels configurator.ConfiguratorChannel
}

func NewAlertBoolCommands[T configurator.AlertBoolType](
	cmdGroup *cmdgroup.CmdGroup,
	cfgTags configurator.IConfiguratorAlertBool[T],
	channels configurator.ConfiguratorChannel,
) *AlertBoolCommands[T] {
	t := &AlertBoolCommands[T]{CmdGroup: cmdGroup, cfgTags: cfgTags, channels: channels}
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
			darkbot_logus.Log.Debug("CreateEnableCmd.consoler running with args=", darkbot_logus.Args(args))
			if !CheckCommandAllowedToRun(cmd, t.channels, t.GetChannelID()) {
				return
			}

			err := t.cfgTags.Enable(t.GetChannelID())
			if err != nil {
				if _, ok := err.(configurator.ErrorZeroAffectedRows); ok {
					printer.Println(cmd, "ERR state was already enabled")
				} else {
					printer.Println(cmd, "ERR ="+err.Error())
				}
				return
			}
			darkbot_logus.Log.Debug("Create Enable is finished", darkbot_logus.Args(args))

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
			darkbot_logus.Log.Debug("CreateDisableCmd.consoler running with args=", darkbot_logus.Args(args))
			if !CheckCommandAllowedToRun(cmd, t.channels, t.GetChannelID()) {
				return
			}

			err := t.cfgTags.Disable(t.GetChannelID())
			if err != nil {
				if _, ok := err.(configurator.ErrorZeroAffectedRows); ok {
					printer.Println(cmd, "ERR state was already disabled")
				} else {
					printer.Println(cmd, "ERR ="+err.Error())
				}
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
			darkbot_logus.Log.Debug("CreateStatusCmd.consoler running with args=", darkbot_logus.Args(args))
			if !CheckCommandAllowedToRun(cmd, t.channels, t.GetChannelID()) {
				return
			}

			_, err := t.cfgTags.Status(t.GetChannelID())
			if err != nil {
				errMsg := err.Error()

				if _, ok := err.(configurator.ErrorZeroAffectedRows); ok {
					printer.Println(cmd, "OK alert status is disabled")
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
	cfgTags  configurator.IConfiguratorAlertString[T]
	channels configurator.ConfiguratorChannel
}

func NewAlertSetStringCommand[T configurator.AlertStringType](
	cmdGroup *cmdgroup.CmdGroup,
	cfgTags configurator.IConfiguratorAlertString[T],
	channels configurator.ConfiguratorChannel,
	allowed_order_keys []types.OrderKey,
) *AlertSetStringCommand[T] {
	t := &AlertSetStringCommand[T]{CmdGroup: cmdGroup, cfgTags: cfgTags, channels: channels}
	t.CreateSetCmd(allowed_order_keys)
	t.CreateUnsetCmd()
	t.CreateStatusCmd()
	return t
}

func (t *AlertSetStringCommand[T]) CreateSetCmd(allowed_order_keys []types.OrderKey) {
	command := &cobra.Command{
		Use:   "set",
		Short: "Set Value (provide 'set StringValue')",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			darkbot_logus.Log.Debug("CreateSetAlertCmd.consoler running with args=", darkbot_logus.Args(args))
			if !CheckCommandAllowedToRun(cmd, t.channels, t.GetChannelID()) {
				return
			}

			str := args[0]

			if len(allowed_order_keys) > 0 {
				is_allowed_tag := false
				for _, tag := range allowed_order_keys {
					if string(tag) == str {
						is_allowed_tag = true
					}
				}

				if !is_allowed_tag {
					printer.Println(cmd, "ERR only next values are allowed: "+strings.Join(utils.CompL(allowed_order_keys,
						func(x types.OrderKey) string { return string(x) }), ", "))
					return
				}
			}

			err := t.cfgTags.Set(t.GetChannelID(), str)
			if err != nil {
				printer.Println(cmd, "ERR msg="+err.Error())
				return
			}
			darkbot_logus.Log.Debug("finished CreateSetCmd", darkbot_logus.Args(args))

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
			darkbot_logus.Log.Debug("CreateUnsetCmd.consoler running with args=", darkbot_logus.Args(args))
			if !CheckCommandAllowedToRun(cmd, t.channels, t.GetChannelID()) {
				return
			}

			err := t.cfgTags.Unset(t.GetChannelID())
			if err != nil {
				if _, ok := err.(configurator.ErrorZeroAffectedRows); ok {
					printer.Println(cmd, "ERR state was already unset")
				} else {
					printer.Println(cmd, "ERR ="+err.Error())
				}
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
			darkbot_logus.Log.Debug("CreateStatusCmd.consoler running with args=", darkbot_logus.Args(args))
			if !CheckCommandAllowedToRun(cmd, t.channels, t.GetChannelID()) {
				return
			}

			str, err := t.cfgTags.Status(t.GetChannelID())
			if err != nil {
				errMsg := err.Error()
				if _, ok := err.(configurator.ErrorZeroAffectedRows); ok {
					printer.Println(cmd, "OK state of alert is disabled")
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
