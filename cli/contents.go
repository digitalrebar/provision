package cli

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"text/template"

	"github.com/digitalrebar/provision/api"
	"github.com/digitalrebar/provision/models"
	"github.com/digitalrebar/store"
	"github.com/spf13/cobra"
)

func init() {
	addRegistrar(registerContent)
}

func registerContent(app *cobra.Command) {
	content := &cobra.Command{
		Use:   "contents",
		Short: "Access CLI commands relating to content",
	}
	content.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "List the installed content bundles",
		Long:  "Provides a summarized version of the content bundles installed on the server",
		Args: func(c *cobra.Command, args []string) error {
			if len(args) == 0 {
				return nil
			}
			return fmt.Errorf("%v does not support filtering", c.UseLine())
		},
		RunE: func(c *cobra.Command, args []string) error {
			summary, err := session.GetContentSummary()
			if err != nil {
				return generateError(err, "listing contents")
			}
			return prettyPrint(summary)
		},
	})
	content.AddCommand(&cobra.Command{
		Use:   "show [id]",
		Short: "Show a single content layer referenced by [id]",
		Args: func(c *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("%v requires 1 argument", c.UseLine())
			}
			return nil
		},
		RunE: func(c *cobra.Command, args []string) error {
			layer, err := session.GetContentItem(args[0])
			if err != nil {
				return generateError(err, "Failed to fetch content: %s", args[0])
			}
			return prettyPrint(layer)
		},
	})
	content.AddCommand(&cobra.Command{
		Use:   "exists [id]",
		Short: "See if content layer referenced by [id] exists",
		Args: func(c *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("%v requires 1 argument", c.UseLine())
			}
			return nil
		},
		RunE: func(c *cobra.Command, args []string) error {
			_, err := session.GetContentItem(args[0])
			if err != nil {
				return fmt.Errorf("content:%s does not exist", args[0])
			}
			return nil
		},
	})
	content.AddCommand(&cobra.Command{
		Use:   "create [json]",
		Short: "Add a new content layer to the system",
		Args: func(c *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("%v requires 1 argument", c.UseLine())
			}
			return nil
		},
		RunE: func(c *cobra.Command, args []string) error {
			layer := &models.Content{}
			if err := into(args[0], layer); err != nil {
				return generateError(err, "Error parsing layer")
			}
			if res, err := session.CreateContent(layer); err != nil {
				return generateError(err, "Error adding content layer")
			} else {
				return prettyPrint(res)
			}
		},
	})
	content.AddCommand(&cobra.Command{
		Use:   "update [id] [json]",
		Short: "Replace a content layer in the system.",
		Args: func(c *cobra.Command, args []string) error {
			if len(args) != 2 {
				return fmt.Errorf("%v requires 2 arguments", c.UseLine())
			}
			return nil
		},
		RunE: func(c *cobra.Command, args []string) error {
			layer := &models.Content{}
			id := args[0]
			if err := into(args[1], layer); err != nil {
				return generateError(err, "Error parsing layer")
			}
			if id != layer.Meta.Name {
				return fmt.Errorf("Passed ID %s does not match layer ID %s", id, layer.Meta.Name)
			}
			if res, err := session.ReplaceContent(layer); err != nil {
				return generateError(err, "Error replacing content layer")
			} else {
				return prettyPrint(res)
			}
		},
	})
	content.AddCommand(&cobra.Command{
		Use:   "upload [json]",
		Short: "Upload a content layer into the system, replacing the earlier one if needed.",
		Args: func(c *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("%v requires 1 arguments", c.UseLine())
			}
			return nil
		},
		RunE: func(c *cobra.Command, args []string) error {
			layer := &models.Content{}
			if err := into(args[0], layer); err != nil {
				return generateError(err, "Error parsing layer")
			}
			if res, err := session.ReplaceContent(layer); err == nil {
				return prettyPrint(res)
			}
			if res, err := session.CreateContent(layer); err == nil {
				return prettyPrint(res)
			} else {
				return generateError(err, "Error uploading layer")
			}
		},
	})
	content.AddCommand(&cobra.Command{
		Use:   "destroy [id]",
		Short: "Remove the content layer [id] from the system.",
		Args: func(c *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("%v requires 1 arguments", c.UseLine())
			}
			return nil
		},
		RunE: func(c *cobra.Command, args []string) error {
			if err := session.DeleteContent(args[0]); err != nil {
				return generateError(err, "Error deleting content layer")
			}
			fmt.Printf("Deleted content %s", args[0])
			return nil
		},
	})
	content.AddCommand(&cobra.Command{
		Use:   "bundle [file] [meta fields]",
		Short: "Bundle the current directory into [file].  [meta fields] allows for the specification of the meta data.",
		Long:  "Bundle assumes that the directories are the object types of the system.",
		Args: func(c *cobra.Command, args []string) error {
			if len(args) == 0 {
				return fmt.Errorf("Must provide a file")
			}
			for i := 1; i < len(args); i++ {
				if !strings.ContainsAny(args[i], "=") {
					return fmt.Errorf("Meta fields must have '=' in them")
				}
			}
			return nil
		},
		RunE: func(c *cobra.Command, args []string) error {
			target := args[0]
			ext := path.Ext(target)
			codec := ""
			switch ext {
			case ".go", ".yaml", ".yml":
				codec = "yaml"
			case ".json":
				codec = "json"
			default:
				return fmt.Errorf("Unknown store extension %s", ext)
			}
			storeURI := fmt.Sprintf("file:%s.tmp?codec=%s", target, codec)
			params := map[string]string{}
			for i := 1; i < len(args); i++ {
				parts := strings.SplitN(args[i], "=", 2)
				params[parts[0]] = parts[1]
			}
			s, err := store.Open(storeURI)
			if err != nil {
				return fmt.Errorf("Failed to open store %s: %v", target, err)
			}
			defer os.Remove(target + ".tmp")
			defer s.Close()
			cc := &api.Client{}
			if err := cc.BundleContent(".", s, params); err != nil {
				return fmt.Errorf("Failed to load: %v", err)
			}
			os.Rename(target+".tmp", target)
			return nil
		},
	})
	content.AddCommand(&cobra.Command{
		Use:   "unbundle [file]",
		Short: "Expand the content bundle [file] into the current directory",
		Args: func(c *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("Must provide a file")
			}
			return nil
		},
		RunE: func(c *cobra.Command, args []string) error {
			src := args[0]
			ext := path.Ext(src)
			switch ext {
			case ".yaml", ".yml", ".json":
			default:
				return fmt.Errorf("Unknown store extension %s", ext)
			}
			buf, err := ioutil.ReadFile(src)
			if err != nil {
				return fmt.Errorf("Failed to open store %s: %v", src, err)
			}
			content := &models.Content{}
			if err := api.DecodeYaml(buf, content); err != nil {
				return fmt.Errorf("Failed to unmarshal store content: %v", err)
			}
			s, _ := store.Open("memory:///")
			if err := content.ToStore(s); err != nil {
				return fmt.Errorf("Failed to open store %s: %v", src, err)
			}
			defer s.Close()
			cc := &api.Client{}
			return cc.UnbundleContent(s, ".")
		},
	})
	content.AddCommand(&cobra.Command{
		Use:   "document [file]",
		Short: "Expand the content bundle [file] into documentation",
		Args: func(c *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("Must provide a file")
			}
			return nil
		},
		RunE: func(c *cobra.Command, args []string) error {
			src := args[0]
			ext := path.Ext(src)
			switch ext {
			case ".yaml", ".yml", ".json":
			default:
				return fmt.Errorf("Unknown store extension %s", ext)
			}
			buf, err := ioutil.ReadFile(src)
			if err != nil {
				return fmt.Errorf("Failed to open store %s: %v", src, err)
			}
			content := &models.Content{}
			if err := api.DecodeYaml(buf, content); err != nil {
				return fmt.Errorf("Failed to unmarshal store content: %v", err)
			}

			tempData := &DocData{
				Name:          content.Meta.Name,
				Version:       content.Meta.Version,
				Documentation: content.Meta.Documentation,
				Objects:       [][]models.Docer{},
			}

			for pref, section := range content.Sections {
				dlist := []models.Docer{}
				for key, obj := range section {
					m, e := models.New(pref)
					if e != nil {
						return fmt.Errorf("Failed to create new %s: %v", pref, e)
					}

					e = models.Remarshal(obj, m)
					if e != nil {
						return fmt.Errorf("Failed to remarshal new %s: %v", key, e)
					}
					if d, ok := m.(models.Docer); ok {
						if d.GetDocumentation() != "" {
							dlist = append(dlist, d)
						}
					}
				}
				if len(dlist) > 0 {
					tempData.Objects = append(tempData.Objects, dlist)
				}
			}

			tmpl := template.New("installLines").Funcs(models.DrpSafeFuncMap()).Option("missingkey=error")
			tmpl, err = tmpl.Parse(docTemplate)
			if err != nil {
				return err
			}

			buf2 := &bytes.Buffer{}
			err = tmpl.Execute(buf2, tempData)
			if err == nil {
				fmt.Println(string(buf2.Bytes()))
			}
			return err
		},
	})
	app.AddCommand(content)
}

type DocData struct {
	Name          string
	Version       string
	Documentation string
	Objects       [][]models.Docer
}

func (dd *DocData) LengthString(str, pat string) string {
	l := len(str)
	out := ""
	for i := 0; i < l; i++ {
		out += pat
	}
	return out
}

var docTemplate = `
.. Copyright (c) 2017 RackN Inc.
.. Licensed under the Apache License, Version 2.0 (the "License");
.. Digital Rebar Provision documentation under Digital Rebar master license
.. index::
  pair: {{.Name}}; Content Packages

.. _rs_cp_{{.Name}}:

{{ $topdot := . }}
{{.Name}}
{{ .LengthString .Name  "~" }}

The following documentation is for {{.Name}} content package at version {{.Version}}.

{{.Documentation}}

{{ range .Objects }}
{{ $b := index . 0 }}
{{ $b := $b.Prefix }}

{{$b}}
{{$topdot.LengthString $b "-"}}

The content package provides the following {{$b}}.

{{ range . }}

{{.Name}}
{{$topdot.LengthString .Name "="}}

{{.GetDocumentation}}

{{ end }}

{{ end }}

`
