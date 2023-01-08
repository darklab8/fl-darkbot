/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package consoler

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
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

func (c *Consoler) CreateRoot() *cobra.Command {
	createdCmd := &cobra.Command{
		Use:   "consoler",
		Short: "A brief description of your application",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("consoler running with args=", args)
		},
	}
	createdCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	return createdCmd
}

func (c *Consoler) CreatePing(rootCmd *cobra.Command) {
	pingCMD := &cobra.Command{
		Use:   "ping",
		Short: "Check stuff is working",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("ping called with args=", args)
			cmd.OutOrStdout().Write([]byte("Pong! from consoler"))
		},
	}
	rootCmd.AddCommand(pingCMD)
}

func (c *Consoler) Create() *cobra.Command {
	rootCmd := c.CreateRoot()
	c.CreatePing(rootCmd)

	return rootCmd
}

func (c *Consoler) Execute() *Consoler {
	rootCmd := c.Create()
	rootCmd.SetArgs(strings.Split(c.cmd, " "))

	rootCmd.SetOut(c.buffStdout)
	rootCmd.SetErr(c.buffStderr)

	rootCmd.Execute()
	return c
}

func (c *Consoler) GetResult() string {
	return c.buffStdout.String() + c.buffStderr.String()
}
