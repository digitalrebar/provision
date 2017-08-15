package cli

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/digitalrebar/provision/client/events"
	models "github.com/digitalrebar/provision/genmodels"
	"github.com/ghodss/yaml"
	"github.com/spf13/cobra"
)

func init() {
	tree := addEventCommands()
	App.AddCommand(tree)
}

func addEventCommands() (res *cobra.Command) {
	res = &cobra.Command{
		Use:   "events",
		Short: "DigitalRebar Provision Event Commands",
	}

	commands := []*cobra.Command{}
	commands = append(commands, &cobra.Command{
		Use:   "post [- | JSON or YAML Event]",
		Short: "Post an event",
		RunE: func(c *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("%v requires 1 argument", c.UseLine())
			}
			dumpUsage = false

			var buf []byte
			var err error
			if args[0] == `-` {
				buf, err = ioutil.ReadAll(os.Stdin)
				if err != nil {
					return fmt.Errorf("Error reading from stdin: %v", err)
				}
			} else {
				buf = []byte(args[0])
			}
			event := models.Event{}
			err = yaml.Unmarshal(buf, &event)
			if err != nil {
				return fmt.Errorf("Invalid event: %v\n", err)
			}

			if _, err := session.Events.PostEvent(events.NewPostEventParams().WithBody(&event), basicAuth); err != nil {
				return generateError(err, "Error posting event")
			} else {
				return nil
			}
		},
	})
	res.AddCommand(commands...)
	return res
}
