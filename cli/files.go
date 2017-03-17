package main

import (
	"fmt"
	"log"
	"os"

	"github.com/rackn/rocket-skates/client/files"
	"github.com/spf13/cobra"
)

func init() {
	tree := addFileCommands()
	app.AddCommand(tree)
}

func addFileCommands() (cmds *cobra.Command) {
	cmds = &cobra.Command{
		Use:   "files",
		Short: "Commands to manage files on the provisioner",
	}
	cmds.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "List all uploaded files",
		Run: func(c *cobra.Command, args []string) {
			if resp, err := session.Files.ListFiles(files.NewListFilesParams()); err != nil {
				log.Fatalf("Error listing files: %v", err)
			} else {
				fmt.Println(pretty(resp.Payload))
			}
		},
	})
	cmds.AddCommand(&cobra.Command{
		Use:   "upload [file] as [name]",
		Short: "Upload a local file to RocketSkates",
		Run: func(c *cobra.Command, args []string) {
			if len(args) != 3 {
				log.Fatalf("Wrong number of args: expected 3, got %d", len(args))
			}
			params := files.NewUploadFileParams()
			params.Path = args[2]
			f, err := os.Open(args[0])
			if err != nil {
				log.Fatalf("Failed to open %s: %v", args[0], err)
			}
			defer f.Close()
			params.Body = f
			if resp, err := session.Files.UploadFile(params); err != nil {
				log.Fatalf("Error uploading: %v", err)
			} else {
				fmt.Println(pretty(resp.Payload))
			}
		},
	})
	return
}
