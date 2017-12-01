package main

import (
	"log"
	"strings"

	"github.com/digitalrebar/provision/cli"
	"github.com/spf13/cobra/doc"
)

func main() {
	linkHandler := func(name string) string {
		return strings.TrimSuffix(name, ".md") + ".html"
	}
	filePrepender := func(name string) string {
		return ""
	}
	err := doc.GenMarkdownTreeCustom(cli.NewApp(), "./doc/cli", filePrepender, linkHandler)
	if err != nil {
		log.Fatal(err)
	}
}
