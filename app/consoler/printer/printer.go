package printer

/*
Those functions are capable to print back to user to Discord via Cobra
*/

import (
	"fmt"

	"github.com/spf13/cobra"
)

func Print(Cmd *cobra.Command, msg string) {
	Cmd.OutOrStdout().Write([]byte(msg))
}

func Println(Cmd *cobra.Command, msg string) {
	Cmd.OutOrStdout().Write([]byte(fmt.Sprintf("%s\n", msg)))
}
