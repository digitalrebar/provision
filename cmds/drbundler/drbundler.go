package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/digitalrebar/provision/api"
	"github.com/digitalrebar/store"
)

func main() {
	args := os.Args

	if len(args) != 3 {
		fmt.Printf("Must provide a directory and a file\n")
		os.Exit(1)
	}

	directory := args[1]
	filename := args[2]

	if _, err := os.Stat(directory); os.IsNotExist(err) {
		s := fmt.Sprintf("package main\n\nvar contentYamlString = \"\"\n")
		buf := []byte(s)
		if err := ioutil.WriteFile(filename, buf, 0644); err != nil {
			fmt.Printf("Failed to write file: %v\n", err)
			os.Exit(1)
		}
		os.Exit(0)
	}

	if dst, err := store.Open(fmt.Sprintf("file:%s.tmp?codec=yaml", filename)); err != nil {
		fmt.Printf("Failed to open store: %v\n", err)
		os.Exit(1)

	} else {
		client := &api.Client{}
		if err := client.BundleContent(directory, dst, map[string]string{}); err != nil {
			fmt.Printf("Failed to load: %v\n", err)
			os.Exit(1)
		}
		dst.Close()
	}

	if contents, err := ioutil.ReadFile(filename + ".tmp"); err != nil {
		fmt.Printf("Failed to readfile: %v\n", err)
		os.Exit(1)
	} else {
		s := fmt.Sprintf("package main\n\nvar contentYamlString = `\n%s\n`\n", string(contents))
		buf := []byte(s)
		if err := ioutil.WriteFile(filename, buf, 0644); err != nil {
			fmt.Printf("Failed to write file: %v\n", err)
			os.Exit(1)
		}
	}
	os.Remove(filename + ".tmp")
}
