/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package consoler

import (
	"darkbot/consoler/commands"
	"strings"
)

type Consoler struct {
	cmd        string
	buffStdout Writer
	buffStderr Writer
}

func (c Consoler) New(cmd string) *Consoler {
	c.cmd = cmd
	c.buffStdout = Writer{}.New()
	c.buffStderr = Writer{}.New()
	return &c
}

func (c *Consoler) Execute() *Consoler {
	rootCmd := commands.Create()
	rootCmd.SetArgs(strings.Split(c.cmd, " "))

	rootCmd.SetOut(c.buffStdout)
	rootCmd.SetErr(c.buffStderr)

	rootCmd.Execute()
	return c
}

func (c *Consoler) String() string {
	return c.buffStdout.String() + c.buffStderr.String()
}
