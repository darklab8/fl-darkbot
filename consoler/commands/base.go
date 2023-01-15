package commands

import (
	"darkbot/configurator"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

type TagCommands struct {
	rootBase  *cobra.Command
	cfgTags   configurator.ConfiguratorTags
	channelID string
}

func (t TagCommands) Init(rootCmd *cobra.Command, cfgTags configurator.ConfiguratorTags, channelID string) {
	rootBase := &cobra.Command{
		Use:   "base",
		Short: "Base Commands",
	}
	rootCmd.AddCommand(rootBase)

	t.rootBase = rootBase
	t.cfgTags = cfgTags
	t.channelID = channelID
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
			fmt.Println("consoler running with args=", args)
			t.cfgTags.TagsAdd(t.channelID, args...)
			fmt.Println(len(args))

			cmd.OutOrStdout().Write([]byte("OK tags are added"))
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
			fmt.Println("consoler running with args=", args)
		},
	}
	t.rootBase.AddCommand(command)
}

func (t TagCommands) CreateTagClear() {
	command := &cobra.Command{
		Use:   "clear",
		Short: "Clear tags",
		Run: func(cmd *cobra.Command, args []string) {
			t.cfgTags.TagsClear(t.channelID)
		},
	}
	t.rootBase.AddCommand(command)
}

func (t TagCommands) CreateTagList() {
	command := &cobra.Command{
		Use:   "list",
		Short: "List tags",
		Run: func(cmd *cobra.Command, args []string) {
			tags := t.cfgTags.TagsList(t.channelID)

			var sb strings.Builder
			for number, tag := range tags {
				sb.WriteString(tag)

				if number != len(tags)-1 {
					sb.WriteString(", ")
				}
			}
			cmd.OutOrStdout().Write([]byte(sb.String()))
		},
	}
	t.rootBase.AddCommand(command)
}
