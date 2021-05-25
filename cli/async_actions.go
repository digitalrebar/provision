package cli

import (
	"fmt"

	"github.com/digitalrebar/provision/v4/models"
	"github.com/spf13/cobra"
)

func init() {
	addRegistrar(registerAsyncAction)
}

func registerAsyncAction(app *cobra.Command) {
	op := &ops{
		name:       "async_actions",
		singleName: "async_action",
		example:    func() models.Model { return &models.AsyncAction{} },
		noCreate:   true,
	}
	op.addCommand(&cobra.Command{
		Use:   "create [json]",
		Short: fmt.Sprintf("Create a new %v with the passed-in JSON or string key", op.singleName),
		Long: `
As a useful shortcut, '-' can be passed to indicate that the JSON should
be read from stdin.

You may also pass in a machine UUID or Name to create a new job on that Name.
`,
		Args: func(c *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("%v requires 1 argument", c.UseLine())
			}
			return nil
		},
		RunE: func(c *cobra.Command, args []string) error {
			ref := &models.AsyncAction{}
			if err := into(args[0], ref); err != nil {
				if args[0] != "-" {
					m := &models.Machine{}
					if err := Session.FillModel(m, args[0]); err != nil {
						if err := Session.FillModel(m, "Name:"+args[0]); err != nil {
							return fmt.Errorf("Unable to create new AsyncAction: Invalid machine %s", args[0])
						}
					}
					ref.Machine = m.Uuid
				}
			}
			if err := Session.CreateModel(ref); err != nil {
				return generateError(err, "Unable to create new %v", op.singleName)
			}
			return prettyPrint(ref)
		},
	})
	op.addCommand(&cobra.Command{
		Use:   "purge",
		Short: "Purge action_actions in excess of the action_action retention preferences",
		Args:  cobra.NoArgs,
		RunE: func(c *cobra.Command, args []string) error {
			var res interface{}
			if err := Session.Req().Meth("DELETE").UrlFor("async_actions").Do(&res); err != nil {
				return err
			}
			return prettyPrint(res)
		},
	})
	op.command(app)
}
