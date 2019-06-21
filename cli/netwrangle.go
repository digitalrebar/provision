// +build linux

package cli

import (
	"fmt"

	gnet "github.com/rackn/gohai/plugins/net"
	"github.com/rackn/netwrangler"
	"github.com/spf13/cobra"
)

func init() {
	addRegistrar(registerNet)
}

func registerNet(app *cobra.Command) {
	net := &cobra.Command{
		Use:   "net",
		Short: "Command for local network management",
	}

	net.AddCommand(&cobra.Command{
		Use:   "phys",
		Short: "Get the physical network interfaces present on the system",
		Args:  cobra.NoArgs,
		RunE: func(c *cobra.Command, args []string) error {
			phys, err := netwrangler.GatherPhys()
			if err != nil {
				return generateError(err, "Failed to fetch phys")
			}
			return prettyPrint(phys)
		},
	})

	bindMac := true
	phyLoc := ""

	wrangle := &cobra.Command{
		Use:   "compile [plan] to [format] at [dest]",
		Short: "Compile the netplan at [plan] into a final configuration of type [format] at [dest]",
		Long: fmt.Sprintf(`[plan] must be a YAML file in netplan.io format
[format] must be one of: %v
[dest] must be the final location for the network configuration in [format]`,
			netwrangler.DestFormats),
		Args: func(c *cobra.Command, args []string) error {
			if len(args) != 5 {
				return fmt.Errorf("%v requires 3 arguments", c.UseLine())
			}
			return nil
		},
		RunE: func(c *cobra.Command, args []string) error {
			var (
				phys []gnet.Interface
				err  error
			)
			if phyLoc == "" {
				phys, err = netwrangler.GatherPhys()
			} else {
				phys, err = netwrangler.GatherPhysFromFile(phyLoc)
			}
			if err != nil {
				return generateError(err, "Failed to fetch phys")
			}
			err = netwrangler.Compile(phys, "netplan", args[2], args[0], args[4], bindMac)
			if err == nil {
				fmt.Printf("Plan %s compiled to format %s at %s", args[0], args[2], args[4])
				return nil
			}
			return generateError(err, "Failed to compile plan")
		},
	}
	wrangle.Flags().StringVar(&phyLoc, "phys", "", "Location for phy definitions.  If not specified, use the live ones from the kernel")
	wrangle.Flags().BoolVar(&bindMac, "bindMac", true, "Bind all base nic definitions to their mac address. Defaults to true, otherwise name bindings will be used.")
	net.AddCommand(wrangle)
	app.AddCommand(net)
}
