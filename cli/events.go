package cli

import (
	"fmt"

	"github.com/digitalrebar/provision/models"
	"github.com/spf13/cobra"
)

func init() {
	addRegistrar(registerEvent)
}

func registerEvent(app *cobra.Command) {
	res := &cobra.Command{
		Use:   "events",
		Short: "DigitalRebar Provision Event Commands",
	}
	res.AddCommand(&cobra.Command{
		Use:   "post [- | JSON or YAML Event]",
		Short: "Post an event",
		Args: func(c *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("%v requires 1 argument", c.UseLine())
			}
			return nil
		},
		RunE: func(c *cobra.Command, args []string) error {
			evt := &models.Event{}
			if err := into(args[0], evt); err != nil {
				return fmt.Errorf("Invalid event: %v\n", err)
			}
			return session.PostEvent(evt)
		},
	})
	app.AddCommand(res)
}
