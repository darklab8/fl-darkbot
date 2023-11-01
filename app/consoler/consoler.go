/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package consoler

import (
	"darkbot/app/consoler/commands"
	"darkbot/app/settings"
	"darkbot/app/settings/types"
	"strings"
)

type Consoler struct {
	cmd        string
	buffStdout Writer
	buffStderr Writer
	channelID  types.DiscordChannelID
	dbpath     types.Dbpath
}

func NewConsoler(
	cmd string,
	channelID types.DiscordChannelID,
	dbpath types.Dbpath,
) *Consoler {
	c := &Consoler{}
	c.cmd = cmd
	c.buffStdout = NewWriter()
	c.buffStderr = NewWriter()
	c.channelID = channelID
	c.dbpath = dbpath
	return c
}

func (c *Consoler) Execute() *Consoler {
	// only commands starting from prefix are allowed
	if !strings.HasPrefix(c.cmd, settings.Config.ConsolerPrefix) {
		return c
	}

	rootCmd := commands.CreateConsoler(c.channelID, c.dbpath)
	rootCmd.SetArgs(strings.Split(c.cmd, " "))

	rootCmd.SetOut(c.buffStdout)
	rootCmd.SetErr(c.buffStderr)

	rootCmd.Execute()
	return c
}

func (c *Consoler) String() string {
	return c.buffStdout.String() + c.buffStderr.String()
}
