package commands

import (
	"darkbot/app/configurator"
	"darkbot/app/consoler/commands/cmdgroup"
	"darkbot/app/consoler/printer"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

type tagCommands struct {
	*cmdgroup.CmdGroup
	cfgTags configurator.IConfiguratorTags
}

func NewTagCommands(cmd *cmdgroup.CmdGroup, cfgTags configurator.IConfiguratorTags) *tagCommands {
	t := &tagCommands{CmdGroup: cmd, cfgTags: cfgTags}
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
			fmt.Println("CreateTagAdd.consoler running with args=", args)
			err := t.cfgTags.TagsAdd(t.GetChannelID(), strings.Join(args, " "))
			if err != nil {
				cmd.OutOrStdout().Write([]byte("ERR msg=" + err.Error()))
				return
			}
			fmt.Println(len(args))

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
			fmt.Println("CreateTagRemove.consoler running with args=", args)
			err := t.cfgTags.TagsRemove(t.GetChannelID(), strings.Join(args, " "))
			if err != nil {
				cmd.OutOrStdout().Write([]byte("ERR msg=" + err.Error()))
				return
			}

			printer.Println(cmd, "OK tags are removed")
		},
	}
	t.CurrentCmd.AddCommand(command)
}

func (t *tagCommands) CreateTagClear() {
	command := &cobra.Command{
		Use:   "clear",
		Short: "Clear tags",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("CreateTagClear.consoler running with args=", args)
			err := t.cfgTags.TagsClear(t.GetChannelID())
			if err != nil {
				cmd.OutOrStdout().Write([]byte("ERR msg=" + err.Error()))
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
			fmt.Println("CreateTagList.consoler running with args=", args)
			tags, cfgErr := t.cfgTags.TagsList(t.GetChannelID())
			err := cfgErr
			if err != nil {
				cmd.OutOrStdout().Write([]byte("ERR msg=" + err.Error()))
				return
			}

			fmt.Println("tags=", tags)
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
