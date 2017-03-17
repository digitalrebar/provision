package main

import (
	"fmt"
	"log"
	"os"

	"github.com/rackn/rocket-skates/client/isos"
	"github.com/spf13/cobra"
)

func init() {
	tree := addIsoCommands()
	app.AddCommand(tree)
}

func addIsoCommands() (cmds *cobra.Command) {
	cmds = &cobra.Command{
		Use:   "isos",
		Short: "Commands to manage ISO files on the provisioner",
	}
	cmds.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "List all uploaded ISOs",
		Run: func(c *cobra.Command, args []string) {
			if resp, err := session.Isos.ListIsos(isos.NewListIsosParams()); err != nil {
				log.Fatalf("Error listing isos: %v", err)
			} else {
				fmt.Println(pretty(resp.Payload))
			}
		},
	})
	cmds.AddCommand(&cobra.Command{
		Use:   "upload [iso] as [name]",
		Short: "Upload a local iso to RocketSkates",
		Run: func(c *cobra.Command, args []string) {
			if len(args) != 3 {
				log.Fatalf("Wrong number of args: expected 3, got %d", len(args))
			}
			params := isos.NewUploadIsoParams()
			params.Path = args[2]
			f, err := os.Open(args[0])
			if err != nil {
				log.Fatalf("Failed to open %s: %v", args[0], err)
			}
			defer f.Close()
			params.Body = f
			if resp, err := session.Isos.UploadIso(params); err != nil {
				log.Fatalf("Error uploading: %v", err)
			} else {
				fmt.Println(pretty(resp.Payload))
			}
		},
	})
	return
}
