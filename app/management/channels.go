/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package management

import (
	"darkbot/app/configurator"
	"darkbot/app/settings"
	"darkbot/app/settings/logus"
	"darkbot/app/settings/types"
	"fmt"

	"github.com/spf13/cobra"
)

var channelDeleteCMD = &cobra.Command{
	Use:   "channel_delete",
	Short: "Delete channelID",
	Run: func(cmd *cobra.Command, args []string) {
		logus.Info("Cmd is called with args=", logus.Args(args))
		channels := configurator.NewConfiguratorChannel(configurator.NewConfigurator(settings.Dbpath))

		for _, channeID := range args {
			channels.Remove(types.DiscordChannelID(channeID))
		}
	},
}
var channelListCMD = &cobra.Command{
	Use:   "channel_list",
	Short: "List channelIDs",
	Run: func(cmd *cobra.Command, args []string) {
		logus.Info("Cmd is called with args=", logus.Args(args))
		channels := configurator.NewConfiguratorChannel(configurator.NewConfigurator(settings.Dbpath))

		channelIDs, _ := channels.List()
		fmt.Println("channelIDs=", channelIDs)
	},
}

func init() {
	rootCmd.AddCommand(channelDeleteCMD)
	rootCmd.AddCommand(channelListCMD)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// migrateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// migrateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
