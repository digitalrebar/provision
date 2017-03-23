package main

import (
	"os"

	"github.com/rackn/rocket-skates/client/files"
	"github.com/spf13/cobra"
)

type FileOps struct{}

func (be FileOps) List() (interface{}, error) {
	d, e := session.Files.ListFiles(files.NewListFilesParams())
	if e != nil {
		return nil, e
	}
	return d.Payload, nil
}

func (be FileOps) Upload(path string, f *os.File) (interface{}, error) {
	d, e := session.Files.UploadFile(files.NewUploadFileParams().WithPath(path).WithBody(f))
	if e != nil {
		return nil, e
	}
	return d.Payload, nil
}

func init() {
	tree := addFileCommands()
	app.AddCommand(tree)
}

func addFileCommands() (res *cobra.Command) {
	singularName := "file"
	name := "files"
	d("Making command tree for %v\n", name)
	res = &cobra.Command{
		Use:   name,
		Short: "Commands to manage files on the provisioner",
	}
	commands := commonOps(singularName, name, &FileOps{})
	res.AddCommand(commands...)
	return res
}
