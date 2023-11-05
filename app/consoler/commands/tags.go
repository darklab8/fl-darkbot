package commands

import (
	"darkbot/app/configurator"
	"darkbot/app/consoler/commands/cmdgroup"
	"darkbot/app/consoler/printer"
	"darkbot/app/settings/logus"
	"darkbot/app/settings/types"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

type tagCommands struct {
	*cmdgroup.CmdGroup
	cfgTags  configurator.IConfiguratorTags
	channels configurator.ConfiguratorChannel
}

func NewTagCommands(cmd *cmdgroup.CmdGroup, cfgTags configurator.IConfiguratorTags, channels configurator.ConfiguratorChannel) *tagCommands {
	t := &tagCommands{CmdGroup: cmd, cfgTags: cfgTags, channels: channels}
	t.CreateTagAdd()
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
			logus.Debug("CreateTagAdd.consoler running with args=", logus.Args(args))
			if !CheckCommandAllowedToRun(cmd, t.channels, t.GetChannelID()) {
				return
			}

			err := t.cfgTags.TagsAdd(t.GetChannelID(), types.Tag(strings.Join(args, " ")))
			if err != nil {
				printer.Println(cmd, "ERR msg="+err.Error())
				return
			}
			logus.Debug("CreateTagAdd", logus.Args(args))

			printer.Println(cmd, "OK tags are added")
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
			logus.Debug("CreateTagRemove.consoler running with args=", logus.Args(args))
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
			logus.Debug("executed Create Tag Remove with args", logus.Args(args))
		},
	}
	t.CurrentCmd.AddCommand(command)

}

func (t *tagCommands) CreateTagClear() {
	command := &cobra.Command{
		Use:   "clear",
		Short: "Clear tags",
		Run: func(cmd *cobra.Command, args []string) {
			logus.Debug("CreateTagClear.consoler running with args=", logus.Args(args))
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
			logus.Debug("CreateTagList.consoler running with args=", logus.Args(args))
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

			logus.Debug("CreateTagList continuied", logus.Tags(tags))
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
