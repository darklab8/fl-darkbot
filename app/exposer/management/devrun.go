package main

import (
	_ "embed"

	"github.com/darklab/fl-darkbot/app/exposer"
)

func main() {
	exposer.NewExposer()
}
