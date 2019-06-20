package cli

import (
	"fmt"
	"io"
	"os"
	"path"

	"github.com/spf13/cobra"
)

func blobCommands(bt string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   bt,
		Short: fmt.Sprintf("Access CLI commands relating to %v", bt),
	}
	cmd.AddCommand(&cobra.Command{
		Use:   "list [path]",
		Short: fmt.Sprintf("List all %v", bt),
		Long:  fmt.Sprintf("You can pass an optional path parameter to show just part of the %s", bt),
		Args: func(c *cobra.Command, args []string) error {
			if len(args) <= 1 {
				return nil
			}
			return fmt.Errorf("%v: Expected 0 or 1 argument", c.UseLine())
		},
		RunE: func(c *cobra.Command, args []string) error {
			req := session.Req().List(bt)
			if len(args) == 1 {
				req.Params("path", args[0])
			}
			data := []interface{}{}
			err := req.Do(&data)
			if err != nil {
				return generateError(err, "listing %v", bt)
			} else {
				return prettyPrint(data)
			}
		},
	})
	cmd.AddCommand(&cobra.Command{
		Use:     "download [item] to [dest]",
		Aliases: []string{"show", "get"},
		Short:   fmt.Sprintf("Download the %v named [item] to [dest]", bt),
		Args: func(c *cobra.Command, args []string) error {
			if len(args) == 1 || len(args) == 3 {
				return nil
			}
			return fmt.Errorf("%v requires 1 or 2 arguments", c.UseLine())
		},
		RunE: func(c *cobra.Command, args []string) error {
			dest := os.Stdout
			if len(args) == 2 && args[1] != "-" {
				var err error
				dest, err = os.Create(args[1])
				if err != nil {
					return fmt.Errorf("Error opening dest file %s: %v", args[1], err)
				}
			}
			if err := session.GetBlob(dest, bt, args[0]); err != nil {
				return generateError(err, "Failed to fetch %v: %v", bt, args[0])
			}
			return nil
		},
	})
	cmd.AddCommand(&cobra.Command{
		Use:   "exists [item]",
		Short: fmt.Sprintf("Checks to see if [item] %s exists and prints its checksum", bt),
		Args: func(c *cobra.Command, args []string) error {
			if len(args) == 1 {
				return nil
			}
			return fmt.Errorf("%v requires 1", c.UseLine())
		},
		RunE: func(c *cobra.Command, args []string) error {
			sum, err := session.GetBlobSum(bt, args[0])
			if err != nil {
				return generateError(err, "Failed to exists %v: %v", bt, args[0])
			}
			fmt.Printf("%s: %s\n", args[0], sum)
			return nil
		},
	})
	cmd.AddCommand(&cobra.Command{
		Use:    "static [item]",
		Hidden: true,
		Short:  "Download [item] from the static file server. They will always go to stdout.",
		Args: func(c *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("%v requires 1 argument", c.UseLine())
			}
			return nil
		},
		RunE: func(c *cobra.Command, args []string) error {
			rd, err := session.File(args[0])
			if rd != nil {
				defer rd.Close()
			}
			if err != nil {
				return err
			}
			_, err = io.Copy(os.Stdout, rd)
			return err
		},
	})
	explode := false
	upload := &cobra.Command{
		Use:   "upload [src] as [dest]",
		Short: fmt.Sprintf("Upload the %v [src] as [dest]", bt),
		Args: func(c *cobra.Command, args []string) error {
			if len(args) == 1 || len(args) == 3 {
				return nil
			}
			return fmt.Errorf("%v requires 1 or 2 arguments", c.UseLine())
		},
		RunE: func(c *cobra.Command, args []string) error {
			item := args[0]
			dest := path.Base(item)
			if len(args) == 3 {
				dest = args[2]
			}
			data, err := urlOrFileAsReadCloser(item)
			if err != nil {
				return fmt.Errorf("Error opening src file %s: %v", item, err)
			}
			defer data.Close()
			if info, err := session.PostBlobExplode(data, explode, bt, dest); err != nil {
				return generateError(err, "Failed to post %v: %v", bt, dest)
			} else {
				return prettyPrint(info)
			}
		},
	}
	upload.Flags().BoolVar(&explode, "explode", false, "Should the upload file be untarred")
	cmd.AddCommand(upload)

	cmd.AddCommand(&cobra.Command{
		Use:   "destroy [item]",
		Short: fmt.Sprintf("Delete the %v [item] on the DRP server", bt),
		Args: func(c *cobra.Command, args []string) error {
			if len(args) == 1 {
				return nil
			}
			return fmt.Errorf("%v requires 1 argument", c.UseLine())
		},
		RunE: func(c *cobra.Command, args []string) error {
			if err := session.DeleteBlob(bt, args[0]); err != nil {
				return generateError(err, "Failed to delete %v: %v", bt, args[0])
			}
			fmt.Printf("Deleted %s", args[0])
			return nil
		},
	})
	return cmd
}

func init() {
	addRegistrar(registerFile)
}

func registerFile(app *cobra.Command) {
	app.AddCommand(blobCommands("files"))
}
