package commands

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/darklab8/fl-darkbot/app/configurator"
	"github.com/darklab8/fl-darkbot/app/consoler/commands/cmdgroup"
	"github.com/darklab8/fl-darkbot/app/consoler/printer"
	"github.com/darklab8/fl-darkbot/app/settings/logus"
	"github.com/darklab8/fl-darkbot/app/settings/types"
	"github.com/darklab8/go-utils/typelog"

	"github.com/spf13/cobra"
)

type tagCommands struct {
	*cmdgroup.CmdGroup
	cfgTags          configurator.IConfiguratorTags
	channels         configurator.ConfiguratorChannel
	disabled_add_cmd bool
}

type tagCommandOpt func(p *tagCommands)

func WithDisabledAdd() tagCommandOpt {
	return func(p *tagCommands) {
		p.disabled_add_cmd = true
	}
}

func NewTagCommands(
	cmd *cmdgroup.CmdGroup,
	cfgTags configurator.IConfiguratorTags,
	channels configurator.ConfiguratorChannel,
	opts ...tagCommandOpt,
) *tagCommands {
	t := &tagCommands{CmdGroup: cmd, cfgTags: cfgTags, channels: channels}

	for _, opt := range opts {
		opt(t)
	}
	if !t.disabled_add_cmd {
		t.CreateTagAdd()
	}

	t.CreateTagRemove()
	t.CreateTagClear()
	t.CreateTagList()
	return t
}

func (t *tagCommands) CreateTagAdd() {
	command := &cobra.Command{
		Use:   "add",
		Short: "Add tags",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			logus.Log.Debug("CreateTagAdd.consoler running with args=", logus.Args(args))
			if !CheckCommandAllowedToRun(cmd, t.channels, t.GetChannelID()) {
				return
			}

			tags := types.Tag(strings.Join(args, " "))
			err := t.cfgTags.TagsAdd(t.GetChannelID(), tags)
			if err != nil {
				printer.Println(cmd, "ERR msg="+err.Error())
				return
			}
			logus.Log.Debug("CreateTagAdd", logus.Args(args))

			printer.Println(cmd, fmt.Sprintf("OK tags are added\n```\n%#v\n```\n", tags))
		},
	}
	t.CurrentCmd.AddCommand(command)
}

func (t *tagCommands) CreateTagRemove() {
	command := &cobra.Command{
		Use:   "remove",
		Short: "Remove tags",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			logus.Log.Debug("CreateTagRemove.consoler running with args=", logus.Args(args))
			if !CheckCommandAllowedToRun(cmd, t.channels, t.GetChannelID()) {
				return
			}

			if len(args) == 0 {
				printer.Println(cmd, "No tags found to remove. Expected at least one")
				return
			}

			err := t.cfgTags.TagsRemove(t.GetChannelID(), types.Tag(strings.Join(args, " ")))
			if err != nil {
				if _, ok := err.(configurator.ErrorZeroAffectedRows); ok {
					printer.Println(cmd, "ERR removed nothing, because inserted value did not match anything present in the list")
				} else {
					printer.Println(cmd, "ERR ="+err.Error())
				}
				return
			}

			printer.Println(cmd, "OK tags are removed: "+strings.Join(args, " "))
			logus.Log.Debug("executed Create Tag Remove with args", logus.Args(args))
		},
	}
	t.CurrentCmd.AddCommand(command)

}

func (t *tagCommands) CreateTagClear() {
	command := &cobra.Command{
		Use:   "clear",
		Short: "Clear tags",
		Run: func(cmd *cobra.Command, args []string) {
			logus.Log.Debug("CreateTagClear.consoler running with args=", logus.Args(args))
			if !CheckCommandAllowedToRun(cmd, t.channels, t.GetChannelID()) {
				return
			}

			err := t.cfgTags.TagsClear(t.GetChannelID())
			if err != nil {
				if _, ok := err.(configurator.ErrorZeroAffectedRows); ok {
					printer.Println(cmd, "ERR list is already empty. nothing to clear.")
				} else {
					printer.Println(cmd, "ERR ="+err.Error())
				}
				return
			}

			printer.Println(cmd, "OK tags are cleared")
		},
	}
	t.CurrentCmd.AddCommand(command)
}

func (t *tagCommands) CreateTagList() {
	command := &cobra.Command{
		Use:   "list",
		Short: "List tags",
		Run: func(cmd *cobra.Command, args []string) {
			logus.Log.Debug("CreateTagList.consoler running with args=", logus.Args(args))
			if !CheckCommandAllowedToRun(cmd, t.channels, t.GetChannelID()) {
				return
			}

			tags, cfgErr := t.cfgTags.TagsList(t.GetChannelID())
			err := cfgErr
			if err != nil {
				if _, ok := err.(configurator.ErrorZeroAffectedRows); ok {
					printer.Println(cmd, "OK tag list is empty")
				} else {
					printer.Println(cmd, "ERR ="+err.Error())
				}
				return
			}

			logus.Log.Debug("CreateTagList continuied", logus.Tags(tags))
			var sb strings.Builder
			for number, tag := range tags {
				sb.WriteString(fmt.Sprintf("\"%s\"", tag))

				if number != len(tags)-1 {
					sb.WriteString(", ")
				}
			}
			printer.Println(cmd, "OK tags are listed")
			printer.Println(cmd, sb.String())
		},
	}
	t.CurrentCmd.AddCommand(command)
}

type PoBGoodCommands[T configurator.AlertPoBGoodType] struct {
	*cmdgroup.CmdGroup
	cfgGoodThresholds configurator.IConfiguratorAlertPoBGood[T]
	channels          configurator.ConfiguratorChannel
}

func NewPoBGoodCommand[T configurator.AlertPoBGoodType](
	cmd *cmdgroup.CmdGroup,
	cfgGoodThresholds configurator.IConfiguratorAlertPoBGood[T],
	channels configurator.ConfiguratorChannel,
) *PoBGoodCommands[T] {
	t := &PoBGoodCommands[T]{CmdGroup: cmd, cfgGoodThresholds: cfgGoodThresholds, channels: channels}
	t.CreatePobGoodAdd()
	t.CreatePobGoodRemove()
	t.CreatePobGoodClear()
	t.CreatePobGoodThresholdList()
	return t
}

func (t *PoBGoodCommands[T]) CreatePobGoodAdd() {
	command := &cobra.Command{
		Use:   "add",
		Short: "Add pob good alert config",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			logus.Log.Debug("CreatePobGoodAdd.consoler running with args=", logus.Args(args))
			if !CheckCommandAllowedToRun(cmd, t.channels, t.GetChannelID()) {
				return
			}

			if len(args) != 2 {
				printer.Println(cmd, fmt.Sprintln("expected to get 2 args good_nickname(string) and threshold(number). got=", len(args)))
				return
			}

			good_nickname := args[0]
			threshold_str := args[1]

			threshold, err := strconv.Atoi(threshold_str)
			if err != nil {
				printer.Println(cmd, "ERR failed to convert second argument to number err="+err.Error())
				return
			}

			err = t.cfgGoodThresholds.Add(t.GetChannelID(), good_nickname, threshold)
			if err != nil {
				printer.Println(cmd, "ERR failed to add to table msg="+err.Error())
				return
			}
			logus.Log.Debug("CreateTagAdd", logus.Args(args))

			printer.Println(cmd, fmt.Sprintf("OK pob good alert is added.\n```\ngood_nickname=%s, threshold=%d\n```\n", good_nickname, threshold))
		},
	}
	t.CurrentCmd.AddCommand(command)
}

func (t *PoBGoodCommands[T]) CreatePobGoodRemove() {
	command := &cobra.Command{
		Use:   "remove",
		Short: "Remove pob good alert configs",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			logus.Log.Debug("CreatePobGoodRemove.consoler running with args=", logus.Args(args))
			if !CheckCommandAllowedToRun(cmd, t.channels, t.GetChannelID()) {
				return
			}

			if len(args) == 0 {
				printer.Println(cmd, "No pob good alert configs found to remove. Expected at least one good_nickname")
				return
			}

			for _, good_nickname := range args {
				err := t.cfgGoodThresholds.Remove(t.GetChannelID(), good_nickname)
				if err != nil {
					if _, ok := err.(configurator.ErrorZeroAffectedRows); ok {
						printer.Println(cmd, "ERR removed nothing, because inserted value did not match anything present in the pob good alert configs")
					} else {
						printer.Println(cmd, "ERR ="+err.Error())
					}
					return
				}
			}

			printer.Println(cmd, "OK pob good alert configs are removed: "+strings.Join(args, " "))
			logus.Log.Debug("executed CreatePobGoodRemove Remove with args", logus.Args(args))
		},
	}
	t.CurrentCmd.AddCommand(command)

}

func (t *PoBGoodCommands[T]) CreatePobGoodClear() {
	command := &cobra.Command{
		Use:   "clear",
		Short: "Clear pob good alert configs",
		Run: func(cmd *cobra.Command, args []string) {
			logus.Log.Debug("CreatePobGoodClear.consoler running with args=", logus.Args(args))
			if !CheckCommandAllowedToRun(cmd, t.channels, t.GetChannelID()) {
				return
			}

			err := t.cfgGoodThresholds.Clear(t.GetChannelID())
			if err != nil {
				if _, ok := err.(configurator.ErrorZeroAffectedRows); ok {
					printer.Println(cmd, "ERR pob good alert configs list is already empty. nothing to clear.")
				} else {
					printer.Println(cmd, "ERR ="+err.Error())
				}
				return
			}

			printer.Println(cmd, "OK pob good alert configs are cleared")
		},
	}
	t.CurrentCmd.AddCommand(command)
}

func (t *PoBGoodCommands[T]) CreatePobGoodThresholdList() {
	command := &cobra.Command{
		Use:   "list",
		Short: "List pob good alert configs",
		Run: func(cmd *cobra.Command, args []string) {
			logus.Log.Debug("CreatePobGoodThresholdList.consoler running with args=", logus.Args(args))
			if !CheckCommandAllowedToRun(cmd, t.channels, t.GetChannelID()) {
				return
			}

			tags, cfgErr := t.cfgGoodThresholds.Get(t.GetChannelID())
			err := cfgErr
			if err != nil {
				if _, ok := err.(configurator.ErrorZeroAffectedRows); ok {
					printer.Println(cmd, "OK pob good alert config list is empty")
				} else {
					printer.Println(cmd, "ERR ="+err.Error())
				}
				return
			}

			logus.Log.Debug("CreatePobGoodThresholdList continuied", typelog.Any("tags", tags))
			var sb strings.Builder
			for tag, threshold := range tags {
				sb.WriteString(fmt.Sprintf("\"{%s:%d}\"", tag, threshold))

				if threshold != len(tags)-1 {
					sb.WriteString(", ")
				}
			}
			printer.Println(cmd, "OK pob good alert configs are listed")
			printer.Println(cmd, sb.String())
		},
	}
	t.CurrentCmd.AddCommand(command)
}
