package commands

import (
	"fmt"
	"strings"

	"github.com/darklab8/fl-darkbot/app/configurator"
	"github.com/darklab8/fl-darkbot/app/consoler/commands/cmdgroup"
	"github.com/darklab8/fl-darkbot/app/consoler/printer"
	"github.com/darklab8/fl-darkbot/app/settings/logus"
	"github.com/darklab8/fl-darkbot/app/settings/types"

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
