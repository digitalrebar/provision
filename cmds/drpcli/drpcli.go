package main

import (
	"os"

	"github.com/digitalrebar/provision/cli"
)

func main() {
	err := cli.NewApp().Execute()
	if err != nil {
		os.Exit(1)
	}
}
