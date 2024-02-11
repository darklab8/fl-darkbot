package printer

/*
Those functions are capable to print back to user to Discord via Cobra
*/

import (
	"fmt"

	"github.com/darklab8/fl-darkbot/app/settings/logus"

	"github.com/spf13/cobra"
)

func Print(Cmd *cobra.Command, msg string) {
	Cmd.OutOrStdout().Write([]byte(msg))
}

func Println(Cmd *cobra.Command, msg string) {
	logus.Log.Debug(fmt.Sprintf("printer.Println msg=%s", msg))
	Cmd.OutOrStdout().Write([]byte(fmt.Sprintf("%s\n", msg)))
}
