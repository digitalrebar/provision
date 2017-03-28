package main

import (
	"os"

	"github.com/rackn/rocket-skates/cli"
)

func main() {
	err := cli.App.Execute()
	if err != nil {
		os.Exit(1)
	}
}
