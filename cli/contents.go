package cli

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"
	"text/template"

	"github.com/digitalrebar/provision/v4/api"
	"github.com/digitalrebar/provision/v4/models"
	"github.com/digitalrebar/provision/v4/store"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/nacl/box"
)

func init() {
	addRegistrar(registerContent)
}

func outputGoBuffer(filename string, buf []byte) error {
	s := fmt.Sprintf("package main\n\nvar contentYamlString = %s\n", strconv.Quote(string(buf)))
	return ioutil.WriteFile(filename, []byte(s), 0644)
}

func decryptForUpload(c *models.Content, key string) error {
	if s, e := Session.Info(); e != nil || !s.HasFeature("secure-params-in-content-packs") {
		return nil
	}
	if key == "" {
		return nil
	}
	pk := []byte{}
	if err := into(key, &pk); err != nil {
		return err
	}
	return c.Mangle(func(prefix string, obj interface{}) (interface{}, error) {
		v, _ := models.New(prefix)
		if err := models.Remarshal(obj, v); err != nil {
			return nil, err
		}
		paramer, ok := v.(models.Paramer)
		if !ok {
			return nil, nil
		}
		params := paramer.GetParams()
		for k := range params {
			sd := &models.SecureData{}
			if err := models.Remarshal(params[k], sd); err != nil || sd.Validate() != nil {
				continue
			}
			if len(pk) != 32 {
				return nil, models.BadKey
			}
			var v interface{}
			if err := sd.Unmarshal(pk, &v); err != nil {
				return nil, err
			}
			params[k] = v
		}
		paramer.SetParams(params)
		return paramer, nil
	})
}

func encryptAfterDownload(c *models.Content) (key []byte, err error) {
	if s, e := Session.Info(); e != nil || !s.HasFeature("secure-params-in-content-packs") {
		return
	}
	sp := []models.Param{}
	if err = Session.Req().Filter("params", "Secure", "Eq", "true").Do(&sp); err != nil {
		return
	}
	if len(sp) == 0 {
		return
	}
	secureParams := map[string]struct{}{}
	for i := range sp {
		secureParams[sp[i].Name] = struct{}{}
	}
	var pubKey, privKey *[32]byte
	pubKey, privKey, err = box.GenerateKey(rand.Reader)
	if err != nil {
		return
	}
	seenSP := false
	if key, err = json.Marshal(privKey[:]); err != nil {
		return
	}
	err = c.Mangle(func(prefix string, obj interface{}) (interface{}, error) {
		v, _ := models.New(prefix)
		if err := models.Remarshal(obj, &v); err != nil {
			return nil, nil
		}
		paramer, ok := v.(models.Paramer)
		if !ok {
			return nil, nil
		}
		params := paramer.GetParams()
		for k := range params {
			if _, ok := secureParams[k]; !ok {
				continue
			}
			sd := &models.SecureData{}
			if err := models.Remarshal(params[k], sd); err == nil && sd.Validate() == nil {
				continue
			}
			sd = &models.SecureData{}
			if err = sd.Marshal(pubKey[:], params[k]); err != nil {
				return nil, err
			}
			seenSP = true
			params[k] = sd
		}
		paramer.SetParams(params)
		return paramer, nil
	})
	if !seenSP {
		key = []byte{}
	}
	return
}

func doReplaceContent(layer *models.Content, key string, replaceWritable bool) error {
	if err := decryptForUpload(layer, key); err != nil {
		return generateError(err, "Error preparing layer")
	}
	summary, err := Session.GetContentSummary()
	if err != nil {
		return generateError(err, "Error uploading layer")
	}
	exists := false
	for _, cs := range summary {
		if cs.Meta.Name == layer.Meta.Name {
			exists = true
			break
		}
	}
	var res interface{}
	if exists {
		res, err = Session.ReplaceContent(layer, replaceWritable)
	} else {
		res, err = Session.CreateContent(layer, replaceWritable)
	}
	if err == nil {
		return prettyPrint(res)
	} else {
		return generateError(err, "Error uploading layer")
	}
}

func replaceContent(path, key string, replaceWritable bool) error {
	layer := &models.Content{}
	if err := into(path, layer); err != nil {
		return generateError(err, "Error parsing layer")
	}
	return doReplaceContent(layer, key, replaceWritable)
}

func registerContent(app *cobra.Command) {
	content := &cobra.Command{
		Use:   "contents",
		Short: "Access CLI commands relating to content",
	}
	var key string
	var replaceWritable bool
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
			summary, err := Session.GetContentSummary()
			if err != nil {
				return generateError(err, "listing contents")
			}
			return prettyPrint(summary)
		},
	})
	cshow := &cobra.Command{
		Use:   "show [id]",
		Short: "Show a single content layer referenced by [id]",
		Args: func(c *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("%v requires 1 argument", c.UseLine())
			}
			return nil
		},
		RunE: func(c *cobra.Command, args []string) error {
			layer, err := Session.GetContentItem(args[0])
			if err != nil {
				return generateError(err, "Failed to fetch content: %s", args[0])
			}
			genKey, err := encryptAfterDownload(layer)
			if err != nil {
				return generateError(err, "Failed to postprocess comments")
			}
			if len(genKey) > 0 {
				if key == "" {
					return fmt.Errorf("Content has secure parameters, but cannot save key!  Use --key=path/to/keyfile.json")
				}
				if err = ioutil.WriteFile(key, genKey, 0600); err != nil {
					return generateError(err, "Error saving key for secure params")
				}
			}
			return prettyPrint(layer)
		},
	}
	cshow.Flags().StringVar(&key, "key", "", "Location to save key for embedded secure parameters")
	content.AddCommand(cshow)
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
			_, err := Session.GetContentItem(args[0])
			if err != nil {
				return fmt.Errorf("content:%s does not exist", args[0])
			}
			return nil
		},
	})
	create := &cobra.Command{
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
			if err := decryptForUpload(layer, key); err != nil {
				return generateError(err, "Error preparing layer")
			}
			if res, err := Session.CreateContent(layer, replaceWritable); err != nil {
				return generateError(err, "Error adding content layer")
			} else {
				return prettyPrint(res)
			}
		},
	}
	create.Flags().BoolVar(&replaceWritable, "replaceWritable", false, "Replace identically named writable objects")
	create.Flags().StringVar(&key, "key", "", "Location of key to use for embedded secure parameters")
	content.AddCommand(create)
	update := &cobra.Command{
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
			if err := decryptForUpload(layer, key); err != nil {
				return generateError(err, "Error preparing layer")
			}
			if id != layer.Meta.Name {
				return fmt.Errorf("Passed ID %s does not match layer ID %s", id, layer.Meta.Name)
			}
			if res, err := Session.ReplaceContent(layer, replaceWritable); err != nil {
				return generateError(err, "Error replacing content layer")
			} else {
				return prettyPrint(res)
			}
		},
	}
	update.Flags().BoolVar(&replaceWritable, "replaceWritable", false, "Replace identically named writable objects")
	update.Flags().StringVar(&key, "key", "", "Location of key to use for embedded secure parameters")
	content.AddCommand(update)
	upload := &cobra.Command{
		Use:   "upload [json]",
		Short: "Upload a content layer into the system, replacing the earlier one if needed.",
		Args: func(c *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("%v requires 1 arguments", c.UseLine())
			}
			return nil
		},
		RunE: func(c *cobra.Command, args []string) error {
			return replaceContent(args[0], key, replaceWritable)
		},
	}
	upload.Flags().BoolVar(&replaceWritable, "replaceWritable", false, "Replace identically named writable objects")
	upload.Flags().StringVar(&key, "key", "", "Location of key to use for embedded secure parameters")
	content.AddCommand(upload)
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
			if err := Session.DeleteContent(args[0]); err != nil {
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
			cc := api.DisconnectedClient()
			if err := cc.BundleContent(".", s, params); err != nil {
				return fmt.Errorf("Failed to load: %v", err)
			}
			s.Close()

			if ext == ".go" {
				if contents, err := ioutil.ReadFile(target + ".tmp"); err != nil {
					return fmt.Errorf("Failed to readfile: %v\n", err)
				} else {
					if err := outputGoBuffer(target, contents); err != nil {
						return fmt.Errorf("Failed to write file: %v\n", err)
					}
				}
			} else {
				os.Rename(target+".tmp", target)
			}
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
			codecString := "json"
			if format == "yaml" || format == "yml" {
				codecString = "yaml"
			}
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
			s, _ := store.Open("memory:///?codec=" + codecString)
			if err := content.ToStore(s); err != nil {
				return fmt.Errorf("Failed to open store %s: %v", src, err)
			}
			defer s.Close()
			cc := api.DisconnectedClient()
			return cc.UnbundleContent(s, ".")
		},
	})

	// Bundlize - takes a list of objects and makes them a bundle - deleting them optionaly.- interactive.
	var remove = false
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
			content := &models.Content{
				Meta: models.ContentMetaData{
					Name:             api.FindOrFake("", "Name", params),
					Description:      api.FindOrFake("", "Description", params),
					Documentation:    api.FindOrFake("", "Documentation", params),
					RequiredFeatures: api.FindOrFake("", "RequiredFeatures", params),
					Version:          api.FindOrFake("", "Version", params),
					Source:           api.FindOrFake("", "Source", params),
					Type:             api.FindOrFake("", "Type", params),
					Prerequisites:    api.FindOrFake("", "Prerequisites", params),
				},
				Sections: map[string]models.Section{},
			}
			if len(objects) == 0 {
				// interactive mode??
				return fmt.Errorf("No object specified")
			}

			// Get objects
			deleteObjects := map[string][]string{}
			for prefix, list := range objects {
				content.Sections[prefix] = map[string]interface{}{}
				for _, lookupKey := range list {
					if strings.ContainsAny(lookupKey, "=") {
						listArgs := strings.SplitN(lookupKey, "=", 2)
						listArgs = append(listArgs, "decode", "true")
						items, err := Session.ListModel(prefix, listArgs...)
						if err != nil {
							return fmt.Errorf("Failed to list: %s: %v", lookupKey, err)
						}
						for _, item := range items {
							content.Sections[prefix][item.Key()] = item
							deleteObjects[prefix] = append(deleteObjects[prefix], item.Key())
						}
					} else {
						item, _ := models.New(prefix)
						if err := Session.Req().UrlFor(prefix, lookupKey).Params("decode", "true").Do(&item); err != nil {
							return fmt.Errorf("Failed to get: %s: %v", lookupKey, err)
						}
						content.Sections[prefix][item.Key()] = item
						deleteObjects[prefix] = append(deleteObjects[prefix], item.Key())
					}
				}
			}
			genKey, err := encryptAfterDownload(content)
			if err != nil {
				return generateError(err, "Failed to postprocess content")
			}
			if len(genKey) > 0 {
				if key == "" {
					return fmt.Errorf("Content has secure parameters, but cannot save key!  Use --key=path/to/keyfile.json")
				}
				if err = ioutil.WriteFile(key, genKey, 0600); err != nil {
					return generateError(err, "Error saving key for secure params")
				}
			}
			storeURI := fmt.Sprintf("file:%s?codec=%s", target, codec)
			s, err := store.Open(storeURI)
			if err != nil {
				return fmt.Errorf("Failed to open store %s: %v", target, err)
			}
			if err := content.ToStore(s); err != nil {
				return fmt.Errorf("Failed to save content layer: %v", err)
			}
			s.Close()

			// delete the objects
			if remove {
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
						if _, err := Session.DeleteModel(prefix, lookupKey); err != nil {
							return fmt.Errorf("Failed to delete %s:%s: %v", prefix, lookupKey, err)
						}
					}
				}
				if reload {
					if err := replaceContent(target, key, remove); err != nil {
						return err
					}
				}
			}
			return nil
		},
	}
	bundlize.Flags().BoolVar(&remove, "delete", false, "Delete bundlized content")
	bundlize.Flags().BoolVar(&reload, "reload", false, "Load the bundle as a content package (requires delete)")
	bundlize.Flags().StringVar(&key, "key", "", "Location to save key for embedded secure parameters")
	content.AddCommand(bundlize)

	// Convert - load a yaml content as read=write objects.
	content.AddCommand(&cobra.Command{
		Use:   "convert [file]",
		Short: "Expand the content bundle [file or - for stdin] into DRP as read-write objects",
		Args: func(c *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("Must provide a file or stdin")
			}
			return nil
		},
		RunE: func(c *cobra.Command, args []string) error {
			src := args[0]
			var buf []byte
			var err error
			if src == "-" {
				buf, err = ioutil.ReadAll(os.Stdin)
			} else {
				ext := path.Ext(src)
				switch ext {
				case ".yaml", ".yml", ".json":
				default:
					return fmt.Errorf("Unknown store extension %s", ext)
				}
				buf, err = ioutil.ReadFile(src)
			}
			if err != nil {
				return fmt.Errorf("Failed to open store %s: %v", src, err)
			}

			content := &models.Content{}
			if yerr := api.DecodeYaml(buf, content); yerr != nil {
				return fmt.Errorf("Failed to unmarshal store content: %v", err)
			}

			for prefix, vals := range content.Sections {
				for _, v := range vals {
					item, _ := models.New(prefix)
					if err := models.Remarshal(v, item); err != nil {
						return fmt.Errorf("Failed to remarshal %s:%v: %v", prefix, v, err)
					}
					if err := Session.CreateModel(item); err != nil {
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
				RefName:       strings.ReplaceAll(content.Meta.Name, "-", "_"),
				DisplayName:   content.Meta.DisplayName,
				Full:          fmt.Sprintf("%s - %s", content.Meta.Name, content.Meta.DisplayName),
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
	RefName       string
	DisplayName   string
	Full          string
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

.. _rs_cp_{{.RefName}}:

{{ $topdot := . }}
{{.Full}}
{{ .LengthString .Full  "~" }}

The following documentation is for {{.DisplayName}} ({{.Name}}) content package at version {{.Version}}.

{{.Documentation}}

Object Specific Documentation
-----------------------------

{{ range .Objects }}
{{ $b := index . 0 }}
{{ $b := $b.Prefix }}

{{$b}}
{{$topdot.LengthString $b "="}}

The content package provides the following {{$b}}.

{{ range . }}

{{.Name}}
{{$topdot.LengthString .Name "+"}}

{{.GetDocumentation}}

{{ end }}

{{ end }}

`
