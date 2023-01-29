package commands

import (
	"darkbot/configurator"
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

	rootGroup := cmdgroup.New(
		consolerCmd,
		channelInfo,
		cmdgroup.CmdGroupProps{
			Command:   settings.Config.ConsolerPrefix,
			ShortDesc: "Welcome to darkbot!",
		},
	)
	root := (&rootCommands{CmdGroup: &rootGroup}).Bootstrap()

	(&TagCommands{
		cfgTags:  configurator.ConfiguratorBase{Configurator: configurator.NewConfigurator()},
		CmdGroup: root.GetChild(root.CurrentCmd, cmdgroup.CmdGroupProps{Command: "base", ShortDesc: "Base commands"}),
	}).Bootstrap()

	playerGroup := root.GetChild(root.CurrentCmd, cmdgroup.CmdGroupProps{Command: "player", ShortDesc: "Player commands"})
	(&TagCommands{
		cfgTags:  configurator.ConfiguratorSystem{Configurator: configurator.NewConfigurator()},
		CmdGroup: playerGroup.GetChild(playerGroup.CurrentCmd, cmdgroup.CmdGroupProps{Command: "system", ShortDesc: "System commands"}),
	}).Bootstrap()

	(&TagCommands{
		cfgTags:  configurator.ConfiguratorRegion{Configurator: configurator.NewConfigurator()},
		CmdGroup: playerGroup.GetChild(playerGroup.CurrentCmd, cmdgroup.CmdGroupProps{Command: "region", ShortDesc: "Region commands"}),
	}).Bootstrap()

	(&TagCommands{
		cfgTags:  configurator.ConfiguratorPlayerFriend{Configurator: configurator.NewConfigurator()},
		CmdGroup: playerGroup.GetChild(playerGroup.CurrentCmd, cmdgroup.CmdGroupProps{Command: "friend", ShortDesc: "Player friend commands"}),
	}).Bootstrap()

	(&TagCommands{
		cfgTags:  configurator.ConfiguratorPlayerEnemy{Configurator: configurator.NewConfigurator()},
		CmdGroup: playerGroup.GetChild(playerGroup.CurrentCmd, cmdgroup.CmdGroupProps{Command: "enemy", ShortDesc: "Player enemy commands"}),
	}).Bootstrap()

	return consolerCmd
}

type rootCommands struct {
	*cmdgroup.CmdGroup
	channels configurator.ConfiguratorChannel
}

func (r *rootCommands) Bootstrap() *rootCommands {
	r.channels = configurator.ConfiguratorChannel{Configurator: r.Configurator}
	r.CreatePing()
	r.CreateConnect()
	r.CreateDisconnect()
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

func (r *rootCommands) CreateConnect() {

	command := &cobra.Command{
		Use:   "connect",
		Short: "Connect bot to channel",
		Run: func(cmd *cobra.Command, args []string) {
			err := r.channels.Add(r.ChannelInfo.ChannelID).GetError()
			if err != nil {
				cmd.OutOrStdout().Write([]byte("ERR channel may be already connected, msg=" + err.Error()))
				return
			}
			cmd.OutOrStdout().Write([]byte("OK Channel is connected"))
		},
	}
	r.CurrentCmd.AddCommand(command)
}

func (r *rootCommands) CreateDisconnect() {
	command := &cobra.Command{
		Use:   "disconnect",
		Short: "Disconnect bot from channel",
		Run: func(cmd *cobra.Command, args []string) {
			err := r.channels.Remove(r.ChannelInfo.ChannelID).GetError()
			if err != nil {
				cmd.OutOrStdout().Write([]byte("ERR channel may be already disconnected, msg=" + err.Error()))
				return
			}
			cmd.OutOrStdout().Write([]byte("OK Channel is disconnected"))
		},
	}
	r.CurrentCmd.AddCommand(command)
}
