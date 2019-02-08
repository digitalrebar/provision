package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	addRegistrar(systemInfo)
}

func systemInfo(app *cobra.Command) {
	tree := addSystemCommands()
	app.AddCommand(tree)
}

func addSystemCommands() (res *cobra.Command) {
	singularName := "system"
	name := "system"
	res = &cobra.Command{
		Use:   name,
		Short: fmt.Sprintf("Access CLI commands relating to %v", name),
	}

	op := &ops{
		name:       name,
		singleName: singularName,
	}
	op.actions()
	res.AddCommand(op.extraCommands...)

	res.AddCommand(&cobra.Command{
		Use:   "upgrade [zip file]",
		Short: "Upgrade DRP with the provided file",
		Args: func(c *cobra.Command, args []string) error {
			if len(args) == 1 {
				return nil
			}
			return fmt.Errorf("%v requires 1 argument", c.UseLine())
		},
		RunE: func(c *cobra.Command, args []string) error {
			filePath := args[0]
			fi, err := urlOrFileAsReadCloser(filePath)
			if err != nil {
				return fmt.Errorf("Error opening %s: %v", filePath, err)
			}
			defer fi.Close()
			if info, err := session.PostBlob(fi, "system", "upgrade"); err != nil {
				return generateError(err, "Failed to post upgrade: %v", filePath)
			} else {
				return prettyPrint(info)
			}
		},
	})

	return res
}
