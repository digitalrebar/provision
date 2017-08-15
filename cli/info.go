package cli

import (
	"fmt"

	"github.com/digitalrebar/provision/client/info"
	"github.com/spf13/cobra"
)

func init() {
	tree := addInfoCommands()
	App.AddCommand(tree)
}

func addInfoCommands() (res *cobra.Command) {
	singularName := "info"
	name := "info"
	d("Making command tree for %v\n", name)
	res = &cobra.Command{
		Use:   name,
		Short: fmt.Sprintf("Access CLI commands relating to %v", name),
	}

	command := &cobra.Command{
		Use:   "get",
		Short: fmt.Sprintf("Get info about DRP"),
		Long:  `A helper function to return information about DRP`,
		RunE: func(c *cobra.Command, args []string) error {
			if len(args) != 0 {
				return fmt.Errorf("%v requires no arguments", c.UseLine())
			}
			dumpUsage = false

			d, err := session.Info.GetInfo(info.NewGetInfoParams(), basicAuth)
			if err != nil {
				return generateError(err, "Failed to fetch info %v", singularName)
			}
			return prettyPrint(d.Payload)
		},
	}

	res.AddCommand(command)
	return res
}
