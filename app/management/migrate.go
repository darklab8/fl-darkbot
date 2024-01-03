/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package management

import (
	"darkbot/app/configurator"
	"darkbot/app/settings"
	"darkbot/app/settings/logus"

	"github.com/spf13/cobra"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Migrate db",
	Run: func(cmd *cobra.Command, args []string) {
		logus.Log.Info("migrate called")
		configurator.NewConfigurator(settings.Dbpath).AutoMigrateSchema()
	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)
}
