package helper

import "github.com/spf13/cobra"

type Printer struct {
	Cmd *cobra.Command
}

func (p Printer) Print(msg string) {
	p.Cmd.OutOrStdout().Write([]byte("OK tags are added"))
}
