package main

import (
	_ "embed"

	"github.com/darklab8/fl-darkbot/app/exposer"
)

func main() {
	exposer.NewExposer()
}
