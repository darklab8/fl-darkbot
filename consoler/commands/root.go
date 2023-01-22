package commands

import (
	"darkbot/consoler/commands/cmdgroup"
	"darkbot/consoler/helper"
	"darkbot/settings"
	"fmt"

	"github.com/spf13/cobra"
)

// Entrance into CLI
func createEntrance() *cobra.Command {
	command := &cobra.Command{
		Use:   "consoler",
		Short: "A brief description of your application",
		// Args:  cobra.MinimumNArgs(1),
		// // When commented out, HELP info is rendered
		// Run: func(cmd *cobra.Command, args []string) {
		// 	// Ignoring message and rendering nothing
		// },
		SilenceUsage:  true,
		SilenceErrors: true,
		Hidden:        true,
	}
	command.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	return command
}

func CreateConsoler(channelInfo helper.ChannelInfo) *cobra.Command {
	consolerCmd := createEntrance()

	root := (&rootCommands{
		CmdGroup: cmdgroup.New(
			consolerCmd,
			channelInfo,
			cmdgroup.CmdGroupProps{
				Command:   settings.Config.ConsolerPrefix,
				ShortDesc: "Welcome to darkbot!",
			},
		)}).Bootstrap()

	(&TagCommands{
		CmdGroup: cmdgroup.New(
			root.CurrentCmd,
			channelInfo,
			cmdgroup.CmdGroupProps{Command: "base", ShortDesc: "Base commands"},
		)}).Bootstrap()

	return consolerCmd
}

type rootCommands struct {
	cmdgroup.CmdGroup
}

func (r *rootCommands) Bootstrap() *rootCommands {
	r.CreatePing()
	return r
}

func (r *rootCommands) CreatePing() {
	command := &cobra.Command{
		Use:   "ping",
		Short: "Check stuff is working",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("ping called with args=", args)
			cmd.OutOrStdout().Write([]byte("Pong! from consoler"))
		},
	}
	r.CurrentCmd.AddCommand(command)
}
