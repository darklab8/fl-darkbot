/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package management

import (
	"darkbot/configurator"
	"darkbot/listener"
	"darkbot/scrappy"
	"darkbot/utils"
	"darkbot/viewer"
	"fmt"

	"github.com/spf13/cobra"
)

// runCmd represents the check command
var runCmd = &cobra.Command{
	Use: "run",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("run called")

		// migrate db
		_ = configurator.ConfiguratorChannel{Configurator: configurator.NewConfigurator()}

		go scrappy.Run()
		go listener.Run()
		go viewer.Run()
		utils.SleepAwaitCtrlC()
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
