package helper

import (
	"fmt"

	"github.com/spf13/cobra"
)

type Printer struct {
	Cmd *cobra.Command
}

func (p Printer) Print(msg string) {
	p.Cmd.OutOrStdout().Write([]byte(msg))
}

func (p Printer) Println(msg string) {
	p.Cmd.OutOrStdout().Write([]byte(fmt.Sprintf("%s\n", msg)))
}
