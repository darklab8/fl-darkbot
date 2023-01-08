/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"darkbot/listener"
	"darkbot/scrappy"
	"darkbot/utils"
	"fmt"

	"github.com/spf13/cobra"
)

// runCmd represents the check command
var runCmd = &cobra.Command{
	Use: "run",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("run called")

		go scrappy.Run()
		go listener.Run()
		utils.SleepAwaitCtrlC()
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
