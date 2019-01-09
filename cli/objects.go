package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	addRegistrar(registerObject)
}

func registerObject(app *cobra.Command) {
	tree := addObjectCommands()
	app.AddCommand(tree)
}

func addObjectCommands() (res *cobra.Command) {
	singularName := "object"
	name := "objects"
	res = &cobra.Command{
		Use:   name,
		Short: fmt.Sprintf("Access CLI commands relating to %v", name),
	}
	res.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "List object types in DRP",
		Long:  `A helper function to return object types in DRP`,
		Args:  cobra.NoArgs,
		RunE: func(c *cobra.Command, args []string) error {
			d, err := session.Objects()
			if err != nil {
				return generateError(err, "Failed to fetch info on %v", singularName)
			}
			return prettyPrint(d)
		},
	})
	return res
}
