package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/ghodss/yaml"

	"github.com/VictorLowther/jsonpatch"
	"github.com/go-openapi/strfmt"
	"github.com/rackn/rocket-skates/client/leases"
	"github.com/rackn/rocket-skates/models"
	"github.com/spf13/cobra"
)

func init() {
	tree := addLeaseCommands()
	app.AddCommand(tree)
}

func addLeaseCommands() (res *cobra.Command) {
	singularName := "lease"
	name := "leases"
	d("Making command tree for %v\n", name)
	res = &cobra.Command{
		Use:   name,
		Short: fmt.Sprintf("Access CLI commands relating to %v", name),
	}
	commands := make([]*cobra.Command, 0, 0)
	commands = append(commands, &cobra.Command{
		Use:   "list",
		Short: fmt.Sprintf("List all %v", name),
		Run: func(c *cobra.Command, args []string) {
			if resp, err := session.Leases.ListLeases(leases.NewListLeasesParams()); err != nil {
				log.Fatalf("Error listing %v: %v", name, err)
			} else {
				fmt.Println(pretty(resp.Payload))
			}
		},
	})
	/* Match not supported today
	commands = append(commands, &cobra.Command{
		Use:   "match [json]",
		Short: fmt.Sprintf("List all %v that match the template in [json]", name),
		Run: func(c *cobra.Command, args []string) {
			if len(args) != 1 {
				log.Fatalf("%v requires 1 argument\n", c.UseLine())
			}
			objs := []interface{}{}
			vals := map[string]interface{}{}
			if err := json.Unmarshal([]byte(args[0]), &vals); err != nil {
				log.Fatalf("Matches not valid JSON\n%v", err)
			}
			if err := session.Match(session.UrlPath(maker()), vals, &objs); err != nil {
				log.Fatalf("Error getting matches for %v\nError:%v\n", singularName, err)
			}
			fmt.Println(prettyJSON(objs))
		},
	})
	*/
	commands = append(commands, &cobra.Command{
		Use:   "show [id]",
		Short: fmt.Sprintf("Show a single %v by id", singularName),
		Run: func(c *cobra.Command, args []string) {
			if len(args) != 1 {
				log.Fatalf("%v requires 1 argument\n", c.UseLine())
			}
			var s strfmt.IPv4
			err := s.Scan(args[0])
			if err != nil {
				log.Fatalf("%v is not a valid IPv4\n", args[0])
			}
			if resp, err := session.Leases.GetLease(leases.NewGetLeaseParams().WithAddress(s)); err != nil {
				log.Fatalf("Failed to fetch %v: %v\n%v\n", singularName, args[0], err)
			} else {
				fmt.Println(pretty(resp.Payload))
			}
		},
	})
	/* Sample not supported today
	commands = append(commands, &cobra.Command{
		Use:   "sample",
		Short: fmt.Sprintf("Get the default values for a %v", singularName),
		Run: func(c *cobra.Command, args []string) {
			if len(args) != 0 {
				log.Fatalf("%v takes no arguments", c.UseLine())
			}
			obj := maker()
			if err := session.Init(obj); err != nil {
				log.Fatalf("Unable to fetch defaults for %v: %v\n", singularName, err)
			}
			fmt.Println(prettyJSON(obj))
		},
	})
	*/
	commands = append(commands, &cobra.Command{
		Use:   "create [json]",
		Short: fmt.Sprintf("Create a new %v with the passed-in JSON", singularName),
		Long:  `As a useful shortcut, you can pass '-' to indicate that the JSON should be read from stdin`,
		Run: func(c *cobra.Command, args []string) {
			if len(args) != 1 {
				log.Fatalf("%v requires 1 argument\n", c.UseLine())
			}
			var buf []byte
			var err error
			if args[0] == "-" {
				buf, err = ioutil.ReadAll(os.Stdin)
				if err != nil {
					log.Fatalf("Error reading from stdin: %v", err)
				}
			} else {
				buf = []byte(args[0])
			}
			lease := &models.Lease{}
			err = yaml.Unmarshal(buf, lease)
			if err != nil {
				log.Fatalf("Invalid %v object: %v\n", singularName, err)
			}
			if resp, err := session.Leases.CreateLease(leases.NewCreateLeaseParams().WithBody(lease)); err != nil {
				log.Fatalf("Unable to create new %v: %v\n", singularName, err)
			} else {
				fmt.Println(pretty(resp.Payload))
			}
		},
	})
	commands = append(commands, &cobra.Command{
		Use:   "update [id] [json]",
		Short: fmt.Sprintf("Unsafely update %v by id with the passed-in JSON", singularName),
		Long:  `As a useful shortcut, you can pass '-' to indicate that the JSON should be read from stdin`,
		Run: func(c *cobra.Command, args []string) {
			if len(args) != 2 {
				log.Fatalf("%v requires 2 arguments\n", c.UseLine())
			}
			var s strfmt.IPv4
			err := s.Scan(args[0])
			if err != nil {
				log.Fatalf("%v is not a valid IPv4\n", args[0])
			}
			if resp, err := session.Leases.GetLease(leases.NewGetLeaseParams().WithAddress(s)); err != nil {
				log.Fatalf("Failed to fetch %v: %v\n%v\n", singularName, args[0], err)
			} else {
				var buf []byte
				var err error
				if args[1] == "-" {
					buf, err = ioutil.ReadAll(os.Stdin)
					if err != nil {
						log.Fatalf("Error reading from stdin: %v", err)
					}
				} else {
					buf = []byte(args[1])
				}
				lease := resp.Payload
				buf2, err := yaml.Marshal(lease)
				if err != nil {
					log.Fatalf("Unable to marshal object: %v\n", err)
				}

				merged, err := safeMergeJSON(buf2, buf)
				if err != nil {
					log.Fatalf("Unable to merge objects: %v\n", err)
				}

				lease = &models.Lease{}
				err = yaml.Unmarshal(merged, lease)
				if err != nil {
					log.Fatalf("Unable to unmarshal merged object: %v\n", err)
				}

				if resp, err := session.Leases.PutLease(leases.NewPutLeaseParams().WithAddress(s).WithBody(lease)); err != nil {
					log.Fatalf("Unable to patch %v\n%v\n", args[0], err)
				} else {
					fmt.Println(pretty(resp.Payload))
				}
			}
		},
	})
	commands = append(commands, &cobra.Command{
		Use:   "patch [objectJson] [changesJson]",
		Short: fmt.Sprintf("Patch %v with the passed-in JSON", singularName),
		Run: func(c *cobra.Command, args []string) {
			if len(args) != 2 {
				log.Fatalf("%v requires 2 arguments\n", c.UseLine())
			}
			obj := &models.Lease{}
			if err := yaml.Unmarshal([]byte(args[0]), obj); err != nil {
				log.Fatalf("Unable to parse %v JSON %v\nError: %v\n", c.UseLine(), args[0], err)
			}
			newObj := &models.Lease{}
			yaml.Unmarshal([]byte(args[0]), newObj)
			if err := yaml.Unmarshal([]byte(args[1]), newObj); err != nil {
				log.Fatalf("Unable to parse %v JSON %v\nError: %v\n", c.UseLine(), args[1], err)
			}
			newBuf, _ := yaml.Marshal(newObj)
			patch, err := jsonpatch.GenerateJSON([]byte(args[0]), newBuf, true)
			if err != nil {
				log.Fatalf("Cannot generate JSON Patch\n%v\n", err)
			}
			p := []*models.JSONPatchOperation{}
			err = yaml.Unmarshal(patch, p)
			if err != nil {
				log.Fatalf("Cannot generate JSON Patch Object\n%v\n", err)
			}
			if resp, err := session.Leases.PatchLease(leases.NewPatchLeaseParams().WithAddress(*obj.Addr).WithBody(p)); err != nil {
				log.Fatalf("Unable to patch %v\n%v\n", args[0], err)
			} else {
				fmt.Println(pretty(resp.Payload))
			}
		},
	})
	commands = append(commands, &cobra.Command{
		Use:   "destroy [id]",
		Short: fmt.Sprintf("Destroy %v by id", singularName),
		Run: func(c *cobra.Command, args []string) {
			if len(args) != 1 {
				log.Fatalf("%v requires 1 argument\n", c.UseLine())
			}
			var s strfmt.IPv4
			err := s.Scan(args[0])
			if err != nil {
				log.Fatalf("%v is not a valid IPv4\n", args[0])
			}
			if _, err := session.Leases.DeleteLease(leases.NewDeleteLeaseParams().WithAddress(s)); err != nil {
				log.Fatalf("Unable to destroy %v %v\nError: %v\n", singularName, args[0], err)
			} else {
				fmt.Printf("Deleted %v %v\n", singularName, args[0])
			}
		},
	})
	commands = append(commands, &cobra.Command{
		Use:   "exists [id]",
		Short: fmt.Sprintf("See if a %v exists by id", singularName),
		Run: func(c *cobra.Command, args []string) {
			if len(args) != 1 {
				log.Fatalf("%v requires 1 argument\n", c.UseLine())
			}
			var s strfmt.IPv4
			err := s.Scan(args[0])
			if err != nil {
				log.Fatalf("%v is not a valid IPv4\n", args[0])
			}
			if _, err := session.Leases.GetLease(leases.NewGetLeaseParams().WithAddress(s)); err != nil {
				log.Fatalf("Failed to fetch %v: %v\n%v\n", singularName, args[0], err)
			} else {
				os.Exit(0)
			}
		},
	})

	res.AddCommand(commands...)
	return res
}
