/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package consoler

import (
	"strings"

	"github.com/darklab8/fl-darkbot/app/configurator"
	"github.com/darklab8/fl-darkbot/app/consoler/commands"
	"github.com/darklab8/fl-darkbot/app/consoler/consoler_types"
	"github.com/darklab8/fl-darkbot/app/prometheuser"
	"github.com/darklab8/fl-darkbot/app/settings"
	"github.com/darklab8/fl-darkbot/app/settings/types"
	"github.com/prometheus/client_golang/prometheus"
)

type Consoler struct {
	dbpath   types.Dbpath
	params   *consoler_types.ChannelParams
	configur *configurator.Configurator
}

func NewConsoler(
	dbpath types.Dbpath,
) *Consoler {
	c := &Consoler{}

	c.dbpath = dbpath
	c.params = consoler_types.NewChannelParams("", dbpath)
	c.configur = configurator.NewConfigurator(dbpath)

	return c
}

type ConsolerOut struct {
	stdout string
}

func (c *Consoler) Execute(
	cmd string,
	channelID types.DiscordChannelID,
) string {
	// only commands starting from prefix are allowed
	if !strings.HasPrefix(cmd, settings.Env.ConsolerPrefix) {
		return ""
	}

	rootCmd := commands.CreateConsoler(c.params, c.configur)
	buffStdout := NewWriter()
	buffStderr := NewWriter()
	rootCmd.SetOut(buffStdout)
	rootCmd.SetErr(buffStderr)

	c.params.SetChannelID(channelID)
	args := strings.Split(cmd, " ")
	rootCmd.SetArgs(args)
	rootCmd.Execute()

	_, remainingArgs, _ := rootCmd.Find(args)
	remained := make(map[string]bool)
	for _, arg := range remainingArgs {
		remained[arg] = true
	}
	var command_args []string
	for _, arg := range args[1:] {
		_, ok := remained[arg]
		if !ok {
			command_args = append(command_args, arg)
		}
	}
	prometheuser.ListenerKindOperation.With(prometheus.Labels{
		"command":    strings.Join(command_args, " "),
		"channel_id": string(channelID),
	}).Inc()

	return buffStdout.String() + buffStderr.String()
}
