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

func findOrFake(field string, args map[string]string) string {
	if p, ok := args[field]; !ok {
		s := "Unspecified"
		if field == "Type" {
			// Default Type should be dynamic
			s = "dynamic"
		} else if field == "RequiredFeatures" {
			// Default RequiredFeatures should be empty string
			s = ""
		}
		return s
	} else {
		return p
	}

}

func replaceContent(path string) error {
	layer := &models.Content{}
	if err := into(path, layer); err != nil {
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
			return replaceContent(args[0])
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

	// Bundlize - takes a list of objects and makes them a bundle - deleting them optionaly.- interactive.
	var delete = false
	var reload = false
	bundlize := &cobra.Command{
		Use:   "bundlize [file] [meta fields] [objects]",
		Short: "Bundle the specified object into [file]. [meta fields] allows for the specification of the meta data. [objects] define which objects to record.",
		Long:  "Bundlize assumes that the objects are read-write.",
		Args: func(c *cobra.Command, args []string) error {
			if len(args) == 0 {
				return fmt.Errorf("Must provide a file")
			}
			for i := 1; i < len(args); i++ {
				if !strings.ContainsAny(args[i], "=") && !strings.ContainsAny(args[i], ":") {
					return fmt.Errorf("Meta fields must have '=' in them.  Objects must have a :.")
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

			params := map[string]string{}
			objects := map[string][]string{}
			for i := 1; i < len(args); i++ {
				arg := args[i]
				ci := strings.IndexRune(arg, ':')
				ei := strings.IndexRune(arg, '=')

				if ci == -1 {
					ci = 10000
				}
				if ei == -1 {
					ei = 10000
				}

				if ci < ei {
					// if colon first, then it is an object key
					parts := strings.SplitN(args[i], ":", 2)
					objects[parts[0]] = append(objects[parts[0]], parts[1])
				} else {
					// if equal first, then it is an meta key
					parts := strings.SplitN(args[i], "=", 2)
					params[parts[0]] = parts[1]
				}
			}

			storeURI := fmt.Sprintf("file:%s.tmp?codec=%s", target, codec)
			s, err := store.Open(storeURI)
			if err != nil {
				return fmt.Errorf("Failed to open store %s: %v", target, err)
			}
			defer os.Remove(target + ".tmp")
			defer s.Close()

			if dm, ok := s.(store.MetaSaver); ok {
				meta := map[string]string{
					"Name":             findOrFake("Name", params),
					"Description":      findOrFake("Description", params),
					"Documentation":    findOrFake("Documentation", params),
					"RequiredFeatures": findOrFake("RequiredFeatures", params),
					"Version":          findOrFake("Version", params),
					"Source":           findOrFake("Source", params),
					"Type":             findOrFake("Type", params),
				}
				dm.SetMetaData(meta)
			}

			if len(objects) == 0 {
				// interactive mode??
				return fmt.Errorf("No object specified")
			}

			// Get objects
			deleteObjects := map[string][]string{}
			for prefix, list := range objects {
				sub, err := s.MakeSub(prefix)
				if err != nil {
					return fmt.Errorf("Cannot make substore %s: %v", prefix, err)
				}

				for _, lookupKey := range list {
					if strings.ContainsAny(lookupKey, "=") {
						items, err := session.ListModel(prefix, strings.SplitN(lookupKey, "=", 2)...)
						if err != nil {
							return fmt.Errorf("Failed to list: %s: %v", lookupKey, err)
						}
						for _, item := range items {
							if err := sub.Save(item.Key(), item); err != nil {
								return fmt.Errorf("Failed to save from list %s:%s: %v", item.Prefix(), item.Key(), err)
							}
							deleteObjects[prefix] = append(deleteObjects[prefix], item.Key())
						}
					} else {
						item, _ := models.New(prefix)
						if err := session.FillModel(item, lookupKey); err != nil {
							return fmt.Errorf("Failed to get: %s: %v", lookupKey, err)
						}
						if err := sub.Save(item.Key(), item); err != nil {
							return fmt.Errorf("Failed to save %s:%s: %v", item.Prefix(), item.Key(), err)
						}
						deleteObjects[prefix] = append(deleteObjects[prefix], item.Key())
					}
				}
			}

			// delete the objects
			if delete {
				deleteOrder := []string{
					"machines",
					"leases",
					"reservations",
					"subnets",
					"roles",
					"users",
					"workflows",
					"stages",
					"bootenvs",
					"tasks",
					"templates",
					"profiles",
					"params",
					"plugins",
					"jobs",
				}
				for _, prefix := range deleteOrder {
					list, ok := deleteObjects[prefix]
					if !ok {
						continue
					}
					for _, lookupKey := range list {
						if _, err := session.DeleteModel(prefix, lookupKey); err != nil {
							return fmt.Errorf("Failed to delete %s:%s: %v", prefix, lookupKey, err)
						}
					}
				}

				if reload {
					if err := replaceContent(target + ".tmp"); err != nil {
						return err
					}
				}
			}

			os.Rename(target+".tmp", target)
			return nil
		},
	}
	bundlize.Flags().BoolVar(&delete, "delete", false, "Delete bundlized content")
	bundlize.Flags().BoolVar(&reload, "reload", false, "Load the bundle as a content package (requires delete)")
	content.AddCommand(bundlize)

	// Convert - load a yaml content as read=write objects.
	content.AddCommand(&cobra.Command{
		Use:   "convert [file]",
		Short: "Expand the content bundle [file] into DRP as read-write objects",
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

			for prefix, vals := range content.Sections {
				for _, v := range vals {
					item, _ := models.New(prefix)
					if err := models.Remarshal(v, item); err != nil {
						return fmt.Errorf("Failed to remarshal %s:%v: %v", prefix, v, err)
					}
					if err := session.CreateModel(item); err != nil {
						return fmt.Errorf("Failed to create %s:%v: %v", prefix, v, err)
					}
				}
			}
			return nil
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
