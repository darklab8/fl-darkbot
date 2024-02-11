/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package management

import (
	"github.com/darklab8/fl-darkbot/app/configurator"
	"github.com/darklab8/fl-darkbot/app/settings"
	"github.com/darklab8/fl-darkbot/app/settings/logus"

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
