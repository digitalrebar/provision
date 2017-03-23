package cli

import (
	"fmt"
	"io/ioutil"
	"log"
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
		Run: func(c *cobra.Command, args []string) {
			if resp, err := Session.Prefs.ListPrefs(prefs.NewListPrefsParams()); err != nil {
				log.Fatalf("Error listing prefs: %v", err)
			} else {
				fmt.Println(pretty(resp.Payload))
			}
		},
	})
	commands = append(commands, &cobra.Command{
		Use:   "set",
		Short: "Set preferences",
		Run: func(c *cobra.Command, args []string) {
			prefsMap := map[string]string{}
			if len(args) == 1 {
				var buf []byte
				var err error
				if args[0] == `-` {
					buf, err = ioutil.ReadAll(os.Stdin)
					if err != nil {
						log.Fatalf("Error reading from stdin: %v", err)
					}
				} else {
					buf = []byte(args[0])
				}
				err = yaml.Unmarshal(buf, prefsMap)
				if err != nil {
					log.Fatalf("Invalid prefs: %v\n", err)
				}
			} else if len(args)%2 == 0 {
				for i := 0; i < len(args); i += 2 {
					prefsMap[args[i]] = args[i+1]
				}
			} else {
				log.Fatalf("prefs set either takes a single argument or a multiple of two, not %d", len(args))
			}
			if resp, err := Session.Prefs.SetPrefs(prefs.NewSetPrefsParams().WithBody(prefsMap)); err != nil {
				log.Fatalf("Error setting prefs: %v", err)
			} else {
				fmt.Println(pretty(resp.Payload))
			}
		},
	})
	res.AddCommand(commands...)
	return res
}
