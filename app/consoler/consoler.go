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

	"github.com/spf13/cobra"
)

type Consoler struct {
	buffStdout Writer
	buffStderr Writer
	dbpath     types.Dbpath
	rootCmd    *cobra.Command
	params     *consoler_types.ChannelParams
}

func NewConsoler(
	dbpath types.Dbpath,
) *Consoler {
	c := &Consoler{}

	c.dbpath = dbpath
	c.params = consoler_types.NewChannelParams("", dbpath)
	configur := configurator.NewConfigurator(dbpath)

	c.rootCmd = commands.CreateConsoler(c.params, configur)
	return c
}

func (c *Consoler) Execute(
	cmd string,
	channelID types.DiscordChannelID,
) *Consoler {
	// only commands starting from prefix are allowed
	if !strings.HasPrefix(cmd, settings.Config.ConsolerPrefix) {
		return c
	}

	c.buffStdout = NewWriter()
	c.buffStderr = NewWriter()
	c.rootCmd.SetOut(c.buffStdout)
	c.rootCmd.SetErr(c.buffStderr)

	c.params.SetChannelID(channelID)
	c.rootCmd.SetArgs(strings.Split(cmd, " "))
	c.rootCmd.Execute()

	return c
}

func (c *Consoler) String() string {
	return c.buffStdout.String() + c.buffStderr.String()
}
