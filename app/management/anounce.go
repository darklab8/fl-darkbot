/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package management

import (
	"darkbot/app/configurator"
	"darkbot/app/discorder"
	"darkbot/app/settings"
	"darkbot/app/settings/darkbot_logus"
	"darkbot/app/viewer/views"
	"strings"

	"github.com/spf13/cobra"
)

// migrateCmd represents the migrate command
var amounceCmd = &cobra.Command{
	Use:   "anounce",
	Short: "Anounce something",
	Run: func(cmd *cobra.Command, args []string) {
		darkbot_logus.Log.Info("Anounce is called with args=", darkbot_logus.Args(args))
		dg := discorder.NewClient()

		channels := configurator.NewConfiguratorChannel(configurator.NewConfigurator(settings.Dbpath))
		channelIDs, _ := channels.List()

		for _, channeID := range channelIDs {
			dg.SengMessage(channeID, views.MsgViewHeader+": "+strings.Join(args, " "))
		}

	},
}

func init() {
	rootCmd.AddCommand(amounceCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// migrateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// migrateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
