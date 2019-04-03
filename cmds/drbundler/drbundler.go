package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"

	"github.com/digitalrebar/provision/api"
	"github.com/digitalrebar/store"
)

func outputBuffer(filename string, buf []byte) error {
	ext := path.Ext(filename)
	if ext == ".go" {
		s := fmt.Sprintf("package main\n\nvar contentYamlString = %s\n", strconv.Quote(string(buf)))
		buf = []byte(s)
	}

	if err := ioutil.WriteFile(filename, buf, 0644); err != nil {
		return err
	}
	return nil

}

func cleanUp(filename, output string) {
	os.Remove(filename + ".tmp")
	fmt.Printf(output)
	os.Exit(1)
}

func main() {
	args := os.Args

	if len(args) != 3 {
		cleanUp("", "Must provide a directory and a file\n")
	}

	directory := args[1]
	filename := args[2]

	if _, err := os.Stat(directory); os.IsNotExist(err) {
		if err := outputBuffer(filename, []byte("")); err != nil {
			cleanUp(filename, fmt.Sprintf("Failed to write file: %v\n", err))
		}
		os.Exit(0)
	}

	ext := path.Ext(filename)
	codec := "yaml"
	if ext == ".go" {
		codec = "yaml"
	} else if ext == ".yaml" || ext == ".yml" {
		codec = "yaml"
	} else if ext == ".json" {
		codec = "json"
	} else {
		cleanUp(filename, fmt.Sprintf("Unknown extension: %s\n", ext))
	}

	if dst, err := store.Open(fmt.Sprintf("file:%s.tmp?codec=%s", filename, codec)); err != nil {
		cleanUp(filename, fmt.Sprintf("Failed to open store: %v\n", err))
	} else {
		client := &api.Client{}
		if err := client.BundleContent(directory, dst, map[string]string{}); err != nil {
			cleanUp(filename, fmt.Sprintf("Failed to load: %v\n", err))
		}
		dst.Close()
	}

	if contents, err := ioutil.ReadFile(filename + ".tmp"); err != nil {
		cleanUp(filename, fmt.Sprintf("Failed to readfile: %v\n", err))
	} else {
		if err := outputBuffer(filename, contents); err != nil {
			cleanUp(filename, fmt.Sprintf("Failed to write file: %v\n", err))
		}
	}
	os.Remove(filename + ".tmp")
}
