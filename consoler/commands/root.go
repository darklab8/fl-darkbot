package commands

import (
	"darkbot/configurator"
	"fmt"

	"github.com/spf13/cobra"
)

func Create(channelID string) *cobra.Command {
	rootCmd := CreateRoot()
	CreatePing(rootCmd)
	TagCommands{}.Init(rootCmd, configurator.Base{Configurator: configurator.NewConfigurator()}, channelID)

	return rootCmd
}

func CreateRoot() *cobra.Command {
	command := &cobra.Command{
		Use:   "consoler",
		Short: "A brief description of your application",
		// When commented out, HELP info is rendered
		// Run: func(cmd *cobra.Command, args []string) {
		// 	fmt.Println("consoler running with args=", args)
		// },
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
