package main

import (
	"os"
	"strings"

	"github.com/digitalrebar/provision/v4/cli"
	gojq "github.com/itchyny/gojq/cli"
)

func main() {
	pgname := os.Args[0]
	pgname = strings.TrimSuffix(pgname, ".exe")
	if strings.HasSuffix(pgname, "jq") {
		os.Exit(gojq.Run())
	}
	err := cli.NewApp().Execute()
	if err != nil {
		os.Exit(1)
	}
}
