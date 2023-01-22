/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package management

import (
	"darkbot/configurator"
	"darkbot/discorder"
	"darkbot/viewer/templ"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

// migrateCmd represents the migrate command
var amounceCmd = &cobra.Command{
	Use:   "anounce",
	Short: "Anounce something",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Anounce is called with args=", args)
		dg := discorder.NewClient()

		// go listener.Run()
		// time.Sleep(5 * time.Second)

		channels := configurator.ConfiguratorChannel{Configurator: configurator.NewConfigurator()}
		channelIDs, _ := channels.List()

		for _, channeID := range channelIDs {
			dg.SengMessage(channeID, templ.MsgViewHeader+": "+strings.Join(args, " "))
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
