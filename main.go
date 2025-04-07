/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"github.com/darklab8/fl-darkbot/app/management"
	"github.com/darklab8/fl-darkbot/app/prometheuser"
)

func main() {
	go prometheuser.Prometheuser()
	management.Execute()
}
