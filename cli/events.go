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
	res.AddCommand(&cobra.Command{
		Use:   "watch [filter]",
		Short: "Watch events as they come in real time. Optional filter can be specified.",
		Args: func(c *cobra.Command, args []string) error {
			if len(args) > 1 {
				return fmt.Errorf("%v requires 0 or 1 argument", c.UseLine())
			}
			return nil
		},
		RunE: func(c *cobra.Command, args []string) error {
			stream, err := session.Events()
			if err != nil {
				return err
			}
			filter := "*.*.*"
			if len(args) == 1 {
				filter = args[0]
			}
			handle, es, err := stream.Register(filter)
			if err != nil {
				return err
			}
			defer stream.Deregister(handle)
			for {
				evt := <-es
				if evt.Err != nil {
					return err
				}
				prettyPrint(evt.E)
			}
		},
	})
	app.AddCommand(res)
}
