/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package management

import (
	"darkbot/app/configurator"
	"darkbot/app/discorder"
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

		fmt.Println("trying to delete ", len(args), args)
		var errors []error
		for _, channeID := range args {
			errors = append(errors, channels.Remove(types.DiscordChannelID(channeID)))
		}

		fmt.Println("error count", len(errors), " | ", errors)
	},
}
var channelListCMD = &cobra.Command{
	Use:   "channel_list",
	Short: "List channelIDs",
	Run: func(cmd *cobra.Command, args []string) {
		dg := discorder.NewClient()
		logus.Info("Cmd is called with args=", logus.Args(args))
		channels := configurator.NewConfiguratorChannel(configurator.NewConfigurator(settings.Dbpath))

		channelIDs, _ := channels.List()
		fmt.Println("all channelIDs", len(channelIDs), " | ", channelIDs)

		var accessable_channels []types.DiscordChannelID
		var error_channels []types.DiscordChannelID
		var error_reason []error

		for _, channel := range channelIDs {
			_, err := dg.GetLatestMessages("838802002582175756")

			if err == nil {
				accessable_channels = append(accessable_channels, channel)
			} else {
				error_channels = append(error_channels, channel)
				error_reason = append(error_reason, err)
			}
		}

		fmt.Println("accessable_channels", len(accessable_channels), " | ", accessable_channels)
		fmt.Println("error_channels", len(error_channels), " | ", error_channels)
		fmt.Println("error_reason", len(error_reason), " | ", error_reason)
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
