package cli

import (
	"fmt"
	"path"

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

	consensus := &cobra.Command{
		Use:   "ha",
		Short: "Access CLI commands to get the state of high availability",
	}

	consensus.AddCommand(&cobra.Command{
		Use:   "id",
		Short: "Get the machine ID of this endpoint in the consensus system",
		Args:  cobra.NoArgs,
		RunE: func(c *cobra.Command, args []string) error {
			var res interface{}
			if err := Session.Req().UrlFor("system", "consensus", "id").Do(&res); err != nil {
				return err
			}
			return prettyPrint(res)
		},
	})
	consensus.AddCommand(&cobra.Command{
		Use:   "leader",
		Short: "Get the machine ID of the leader in the consensus system",
		Args:  cobra.NoArgs,
		RunE: func(c *cobra.Command, args []string) error {
			var res interface{}
			if err := Session.Req().UrlFor("system", "consensus", "leader").Do(&res); err != nil {
				return err
			}
			return prettyPrint(res)
		},
	})
	consensus.AddCommand(&cobra.Command{
		Use:   "active",
		Short: "Get the machine ID of the current active node in the consensus system",
		Args:  cobra.NoArgs,
		RunE: func(c *cobra.Command, args []string) error {
			var res interface{}
			if err := Session.Req().UrlFor("system", "consensus", "active").Do(&res); err != nil {
				return err
			}
			return prettyPrint(res)
		},
	})
	consensus.AddCommand(&cobra.Command{
		Use:   "peers",
		Short: "Get basic info on all members of the consensus system",
		Args:  cobra.NoArgs,
		RunE: func(c *cobra.Command, args []string) error {
			var res interface{}
			if err := Session.Req().UrlFor("system", "consensus", "peers").Do(&res); err != nil {
				return err
			}
			return prettyPrint(res)
		},
	})
	consensus.AddCommand(&cobra.Command{
		Use:   "dump",
		Short: "Dump the detailed state of the consensus system.",
		Args:  cobra.NoArgs,
		RunE: func(c *cobra.Command, args []string) error {
			var res interface{}
			if err := Session.Req().UrlFor("system", "consensus", "state").Do(&res); err != nil {
				return err
			}
			return prettyPrint(res)
		},
	})

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
			if info, err := Session.PostBlob(fi, "system", "upgrade"); err != nil {
				return generateError(err, "Failed to post upgrade: %v", filePath)
			} else {
				return prettyPrint(info)
			}
		},
	})
	res.AddCommand(&cobra.Command{
		Use:   "passive",
		Short: "Switch DRP to HA Passive State",
		Args: func(c *cobra.Command, args []string) error {
			if len(args) == 0 {
				return nil
			}
			return fmt.Errorf("%v requires 0 argument", c.UseLine())
		},
		RunE: func(c *cobra.Command, args []string) error {
			r := Session.Req().Post(nil).UrlFor(path.Join("/", "system", "passive"))
			var info interface{}
			if err := r.Do(&info); err != nil {
				return generateError(err, "Failed to set passive state")
			} else {
				return prettyPrint(info)
			}
		},
	})
	res.AddCommand(&cobra.Command{
		Use:   "active",
		Short: "Switch DRP to HA Active State",
		Args: func(c *cobra.Command, args []string) error {
			if len(args) == 0 {
				return nil
			}
			return fmt.Errorf("%v requires 0 argument", c.UseLine())
		},
		RunE: func(c *cobra.Command, args []string) error {
			r := Session.Req().Post(nil).UrlFor(path.Join("/", "system", "active"))
			var info interface{}
			if err := r.Do(&info); err != nil {
				return generateError(err, "Failed to set active state")
			} else {
				return prettyPrint(info)
			}
		},
	})
	res.AddCommand(&cobra.Command{
		Use:   "signurl [URL]",
		Short: "Generate a RackN Signed URL for download",
		Args: func(c *cobra.Command, args []string) error {
			if len(args) == 1 {
				return nil
			}
			return fmt.Errorf("%v requires 1 argument", c.UseLine())
		},
		RunE: func(c *cobra.Command, args []string) error {
			if newurl, err := signRackNUrl(args[0]); err != nil {
				return generateError(err, "Failed to sign url")
			} else {
				fmt.Println(newurl)
				return nil
			}
		},
	})
	res.AddCommand(consensus)

	return res
}
