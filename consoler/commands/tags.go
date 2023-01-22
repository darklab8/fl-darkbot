package commands

import (
	"darkbot/consoler/helper"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

type TagCommands struct {
	CommandGroup
}

func (t TagCommands) Bootstrap() {
	t.CreateTagAdd()
	t.CreateTagRemove()
	t.CreateTagClear()
	t.CreateTagList()
}

func (t TagCommands) CreateTagAdd() {
	command := &cobra.Command{
		Use:   "add",
		Short: "Add tags",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("CreateTagAdd.consoler running with args=", args)
			t.CfgTags.TagsAdd(t.ChannelInfo.ChannelID, args...)
			fmt.Println(len(args))

			helper.Printer{Cmd: cmd}.Println("OK tags are added")
		},
	}
	t.CurrentCmd.AddCommand(command)
}

func (t TagCommands) CreateTagRemove() {
	command := &cobra.Command{
		Use:   "remove",
		Short: "Remove tags",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("CreateTagRemove.consoler running with args=", args)
			t.CfgTags.TagsRemove(t.ChannelInfo.ChannelID, args...)

			helper.Printer{Cmd: cmd}.Println("OK tags are removed")
		},
	}
	t.CurrentCmd.AddCommand(command)
}

func (t TagCommands) CreateTagClear() {
	command := &cobra.Command{
		Use:   "clear",
		Short: "Clear tags",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("CreateTagClear.consoler running with args=", args)
			t.CfgTags.TagsClear(t.ChannelInfo.ChannelID)

			helper.Printer{Cmd: cmd}.Println("OK tags are cleared")
		},
	}
	t.CurrentCmd.AddCommand(command)
}

func (t TagCommands) CreateTagList() {
	command := &cobra.Command{
		Use:   "list",
		Short: "List tags",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("CreateTagList.consoler running with args=", args)
			tags := t.CfgTags.TagsList(t.ChannelInfo.ChannelID)
			fmt.Println("tags=", tags)
			var sb strings.Builder
			for number, tag := range tags {
				sb.WriteString(tag)

				if number != len(tags)-1 {
					sb.WriteString(", ")
				}
			}
			printer := helper.Printer{Cmd: cmd}
			printer.Println("OK tags are listed")
			printer.Println(sb.String())
		},
	}
	t.CurrentCmd.AddCommand(command)
}
