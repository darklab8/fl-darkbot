package consoler

import (
	"fmt"

	"github.com/spf13/cobra"
)

var pingCMD = &cobra.Command{
	Use:   "ping",
	Short: "Check stuff is working",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("check called")
	},
}

func init() {
	rootCmd.AddCommand(pingCMD)
}
