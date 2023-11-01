/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package consoler

import (
	"darkbot/consoler/commands"
	"darkbot/consoler/helper"
	"darkbot/settings"
	"strings"
)

type Consoler struct {
	cmd        string
	buffStdout Writer
	buffStderr Writer
}

func NewConsoler(cmd string) *Consoler {
	c := &Consoler{}
	c.cmd = cmd
	c.buffStdout = NewWriter()
	c.buffStderr = NewWriter()
	return c
}

func (c *Consoler) Execute(channelInfo helper.ChannelInfo) *Consoler {
	// only commands starting from prefix are allowed
	if !strings.HasPrefix(c.cmd, settings.Config.ConsolerPrefix) {
		return c
	}

	rootCmd := commands.CreateConsoler(channelInfo)
	rootCmd.SetArgs(strings.Split(c.cmd, " "))

	rootCmd.SetOut(c.buffStdout)
	rootCmd.SetErr(c.buffStderr)

	rootCmd.Execute()
	return c
}

func (c *Consoler) String() string {
	return c.buffStdout.String() + c.buffStderr.String()
}
