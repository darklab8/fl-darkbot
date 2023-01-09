package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

func AddStdCommands(rootCmd *cobra.Command) {
	CreateTagAdd(rootCmd)
}

func CreateTagAdd(rootCmd *cobra.Command) {
	command := &cobra.Command{
		Use:   "add",
		Short: "Add tags",
		// When commented out, HELP info is rendered
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("consoler running with args=", args)
		},
	}
	rootCmd.AddCommand(command)
}

func CreateTagRemove(rootCmd *cobra.Command) {
	command := &cobra.Command{
		Use:   "remove",
		Short: "Remove tags",
		// When commented out, HELP info is rendered
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("consoler running with args=", args)
		},
	}
	rootCmd.AddCommand(command)
}

func CreateTagClear(rootCmd *cobra.Command) {
	command := &cobra.Command{
		Use:   "clear",
		Short: "Clear tags",
		// When commented out, HELP info is rendered
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("consoler running with args=", args)
		},
	}
	rootCmd.AddCommand(command)
}
