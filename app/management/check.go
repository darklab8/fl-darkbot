/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package management

import (
	"github.com/darklab/fl-darkbot/app/discorder"
	"github.com/darklab/fl-darkbot/app/settings/logus"

	"github.com/spf13/cobra"
)

// checkCmd represents the check command
var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Experimental command",
	Run: func(cmd *cobra.Command, args []string) {
		logus.Log.Info("check called")
		dg := discorder.NewClient()
		// dg.SengMessage("838802002582175756", "123message")
		dg.GetLatestMessages("838802002582175756")
	},
}

func init() {
	rootCmd.AddCommand(checkCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// checkCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// checkCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
