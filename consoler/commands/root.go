package commands

import (
	"darkbot/configurator"
	"darkbot/consoler/helper"
	"darkbot/settings"
	"fmt"

	"github.com/spf13/cobra"
)

func Create(channelInfo helper.ChannelInfo) *cobra.Command {
	rootCmd := CreateRoot()
	rootCmdPrefix := &cobra.Command{
		Use:   settings.Config.ConsolerPrefix,
		Short: "Welcome to darkbot!",
	}
	rootCmd.AddCommand(rootCmdPrefix)

	CreatePing(rootCmdPrefix)

	TagCommands{}.Init(
		rootCmdPrefix,
		configurator.ConfiguratorBase{Configurator: configurator.NewConfigurator()},
		channelInfo,
	)

	return rootCmd
}

func CreateRoot() *cobra.Command {
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

func CreatePing(parentCommand *cobra.Command) {
	command := &cobra.Command{
		Use:   "ping",
		Short: "Check stuff is working",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("ping called with args=", args)
			cmd.OutOrStdout().Write([]byte("Pong! from consoler"))
		},
	}
	parentCommand.AddCommand(command)
}
