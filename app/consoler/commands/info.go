package commands

import (
	"fmt"
	"strings"

	"github.com/darklab8/fl-darkbot/app/configurator"
	"github.com/darklab8/fl-darkbot/app/consoler/commands/cmdgroup"
	"github.com/darklab8/fl-darkbot/app/consoler/printer"
	"github.com/darklab8/fl-darkbot/app/settings"
	"github.com/darklab8/fl-darkbot/app/settings/logus"
	"github.com/darklab8/fl-darkstat/darkapis/darkhttp"
	"github.com/spf13/cobra"
)

type InfoCommands struct {
	*cmdgroup.CmdGroup
	channels configurator.ConfiguratorChannel
}

func NewInfoCommands(
	cmdGroup *cmdgroup.CmdGroup,
	channels configurator.ConfiguratorChannel,
) *InfoCommands {
	t := &InfoCommands{CmdGroup: cmdGroup, channels: channels}
	t.CreateGetInfoCmd()
	return t
}

func (t *InfoCommands) CreateGetInfoCmd() {
	command := &cobra.Command{
		Use:   "info",
		Short: "Get info about some Freelancer entity",
		// Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			logus.Log.Debug("Get info about some Freelancer entity", logus.Args(args))

			// printer.Println(cmd, "Attempting to parse input into integer number")
			query_args := args
			query := strings.Join(query_args, " ")

			answer := "OK giving info about(" + fmt.Sprintf("%d", len(query)) + "): " + query
			printer.Println(cmd, answer)

			client := darkhttp.NewClient(settings.Env.DarkstatApiUrl)
			reply, err := client.GetInfo(darkhttp.GetInfoArgs{Query: query})

			if err != nil {
				printer.Println(cmd, err.Error())
			} else {
				printer.Println(cmd, "Got info!")
				for _, line := range reply.Content {
					printer.Println(cmd, line)
				}
			}
		},
	}
	t.CurrentCmd.AddCommand(command)
}
