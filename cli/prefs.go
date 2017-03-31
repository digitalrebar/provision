package cli

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/ghodss/yaml"
	"github.com/rackn/rocket-skates/client/prefs"
	"github.com/spf13/cobra"
)

func init() {
	tree := addPrefCommands()
	App.AddCommand(tree)
}

func addPrefCommands() (res *cobra.Command) {
	res = &cobra.Command{
		Use:   "prefs",
		Short: "List and set RocketSkates operation preferences",
	}

	commands := []*cobra.Command{}
	commands = append(commands, &cobra.Command{
		Use:   "list",
		Short: "List all preferences",
		RunE: func(c *cobra.Command, args []string) error {
			dumpUsage = false
			if resp, err := session.Prefs.ListPrefs(prefs.NewListPrefsParams()); err != nil {
				return generateError(err, "Error listing prefs")
			} else {
				return prettyPrint(resp.Payload)
			}
		},
	})
	commands = append(commands, &cobra.Command{
		Use:   "set [- | JSON or YAML Map of strings | pairs of string args]",
		Short: "Set preferences",
		RunE: func(c *cobra.Command, args []string) error {
			prefsMap := map[string]string{}
			if len(args) == 1 {
				var buf []byte
				var err error
				if args[0] == `-` {
					buf, err = ioutil.ReadAll(os.Stdin)
					if err != nil {
						dumpUsage = false
						return fmt.Errorf("Error reading from stdin: %v", err)
					}
				} else {
					buf = []byte(args[0])
				}
				err = yaml.Unmarshal(buf, &prefsMap)
				if err != nil {
					dumpUsage = false
					return fmt.Errorf("Invalid prefs: %v\n", err)
				}
			} else if len(args) != 0 && len(args)%2 == 0 {
				for i := 0; i < len(args); i += 2 {
					prefsMap[args[i]] = args[i+1]
				}
			} else {
				return fmt.Errorf("prefs set either takes a single argument or a multiple of two, not %d", len(args))
			}
			dumpUsage = false
			if resp, err := session.Prefs.SetPrefs(prefs.NewSetPrefsParams().WithBody(prefsMap)); err != nil {
				return generateError(err, "Error setting prefs")
			} else {
				return prettyPrint(resp.Payload)
			}
		},
	})
	res.AddCommand(commands...)
	return res
}
