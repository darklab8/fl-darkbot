/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package consoler

import (
	"darkbot/app/configurator"
	"darkbot/app/consoler/commands"
	"darkbot/app/consoler/consoler_types"
	"darkbot/app/settings"
	"darkbot/app/settings/types"
	"strings"
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

func (c *Consoler) Execute(
	cmd string,
	channelID types.DiscordChannelID,
) string {
	// only commands starting from prefix are allowed
	if !strings.HasPrefix(cmd, settings.Config.ConsolerPrefix) {
		return ""
	}

	rootCmd := commands.CreateConsoler(c.params, c.configur)
	buffStdout := NewWriter()
	buffStderr := NewWriter()
	rootCmd.SetOut(buffStdout)
	rootCmd.SetErr(buffStderr)

	c.params.SetChannelID(channelID)
	rootCmd.SetArgs(strings.Split(cmd, " "))
	rootCmd.Execute()

	return buffStdout.String() + buffStderr.String()
}
