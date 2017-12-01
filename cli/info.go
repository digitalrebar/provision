package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	addRegistrar(registerInfo)
}

func registerInfo(app *cobra.Command) {
	tree := addInfoCommands()
	app.AddCommand(tree)
}

func addInfoCommands() (res *cobra.Command) {
	singularName := "info"
	name := "info"
	res = &cobra.Command{
		Use:   name,
		Short: fmt.Sprintf("Access CLI commands relating to %v", name),
	}

	command := &cobra.Command{
		Use:   "get",
		Short: fmt.Sprintf("Get info about DRP"),
		Long:  `A helper function to return information about DRP`,
		Args:  cobra.NoArgs,
		RunE: func(c *cobra.Command, args []string) error {

			d, err := session.Info()
			if err != nil {
				return generateError(err, "Failed to fetch info %v", singularName)
			}
			return prettyPrint(d)
		},
	}

	res.AddCommand(command)
	return res
}
