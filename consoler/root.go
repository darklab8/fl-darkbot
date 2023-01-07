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
	cmd    string
	result string
}

func (c Consoler) New(cmd string) *Consoler {
	c.cmd = cmd
	return &c
}

func (c *Consoler) CreateRoot() *cobra.Command {
	createdCmd := &cobra.Command{
		Use:   "consoler",
		Short: "A brief description of your application",
		// Uncomment the following line if your bare application
		// has an action associated with it:
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
			c.result = "Pong! (from consoler)"
		},
	}
	rootCmd.AddCommand(pingCMD)
}

func (c *Consoler) Create() *cobra.Command {
	rootCmd := c.CreateRoot()
	c.CreatePing(rootCmd)

	return rootCmd
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func (c *Consoler) Execute() *Consoler {
	rootCmd := c.Create()
	rootCmd.SetArgs(strings.Split(c.cmd, " "))
	rootCmd.Execute()
	return c
}

func (c *Consoler) GetResult() string {
	return c.result
}
