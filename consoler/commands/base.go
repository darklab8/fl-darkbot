package commands

import (
	"darkbot/configurator"
	"darkbot/consoler/helper"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

type TagCommands struct {
	rootBase    *cobra.Command
	cfgTags     configurator.ConfiguratorTags
	channelInfo helper.ChannelInfo
}

func (t TagCommands) Init(rootCmd *cobra.Command, cfgTags configurator.ConfiguratorTags, channelInfo helper.ChannelInfo) {
	rootBase := &cobra.Command{
		Use:   "base",
		Short: "Base Commands",
	}
	rootCmd.AddCommand(rootBase)

	t.rootBase = rootBase
	t.cfgTags = cfgTags
	t.channelInfo = channelInfo
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
			t.cfgTags.TagsAdd(t.channelInfo.ChannelID, args...)
			fmt.Println(len(args))

			helper.Printer{Cmd: cmd}.Println("OK tags are added")
		},
	}
	t.rootBase.AddCommand(command)
}

func (t TagCommands) CreateTagRemove() {
	command := &cobra.Command{
		Use:   "remove",
		Short: "Remove tags",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("CreateTagRemove.consoler running with args=", args)
			t.cfgTags.TagsRemove(t.channelInfo.ChannelID, args...)

			helper.Printer{Cmd: cmd}.Println("OK tags are removed")
		},
	}
	t.rootBase.AddCommand(command)
}

func (t TagCommands) CreateTagClear() {
	command := &cobra.Command{
		Use:   "clear",
		Short: "Clear tags",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("CreateTagClear.consoler running with args=", args)
			t.cfgTags.TagsClear(t.channelInfo.ChannelID)

			helper.Printer{Cmd: cmd}.Println("OK tags are cleared")
		},
	}
	t.rootBase.AddCommand(command)
}

func (t TagCommands) CreateTagList() {
	command := &cobra.Command{
		Use:   "list",
		Short: "List tags",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("CreateTagList.consoler running with args=", args)
			tags := t.cfgTags.TagsList(t.channelInfo.ChannelID)
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
	t.rootBase.AddCommand(command)
}
