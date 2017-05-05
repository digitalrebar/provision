package main

import (
	"os"

	"github.com/digitalrebar/provision/cli"
)

func main() {
	err := cli.App.Execute()
	if err != nil {
		os.Exit(1)
	}
}
