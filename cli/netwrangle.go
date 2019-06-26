// +build linux

package cli

import (
	"fmt"

	"github.com/rackn/netwrangler"
	"github.com/rackn/netwrangler/netplan"
	"github.com/rackn/netwrangler/util"
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
	bootMac := ""

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
				phys []util.Phy
				err  error
			)
			netwrangler.BootMac(bootMac)
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
	wrangle.Flags().StringVar(&bootMac, "bootmac", "", "MAC address of the interface the system booted from.")
	net.AddCommand(wrangle)
	generate := &cobra.Command{
		Use:   "generate [format] for [machine-id] at [dest]",
		Short: "Generate network configuration in [format] for [machine-id] at [dest]",
		Long: fmt.Sprintf(`[format] must be one of: %v
[machine-id] must me the UUID for the machine in question
[dest] must be the final location for the network configuration in [format]`, netwrangler.DestFormats),
		Args: func(c *cobra.Command, args []string) error {
			if len(args) != 5 {
				return fmt.Errorf("%v requires 3 arguments", c.UseLine())
			}
			return nil
		},
		RunE: func(c *cobra.Command, args []string) error {
			topology := &netplan.Netplan{}
			addressing := map[string]interface{}{}
			if err := session.Req().
				UrlFor("machines", args[2], "params", "net/interface-topology").
				Params("aggregate", "true").Do(topology); err != nil {
				return err
			}
			if err := session.Req().
				UrlFor("machines", args[2], "params", "net/interface-config").
				Params("aggregate", "true").Do(&addressing); err != nil {
				return err
			}
			if err := session.Req().
				UrlFor("machines", args[2], "params", "last-boot-macaddr").
				Do(&bootMac); err != nil {
				return err
			}
			netwrangler.BootMac(bootMac)
			phys, err := netwrangler.GatherPhys()
			if err != nil {
				return err
			}
			if topology.Network.Vlans != nil {
				for k := range topology.Network.Vlans {
					v, ok := addressing[k]
					if !ok {
						continue
					}
					delete(addressing, k)
					from, fromOk := v.(map[string]interface{})
					to, toOk := topology.Network.Vlans[k].(map[string]interface{})
					if !fromOk || !toOk {
						continue
					}
					for fk := range from {
						to[fk] = from[fk]
					}
					topology.Network.Vlans[k] = to
				}
			}
			if topology.Network.Bonds != nil {
				for k := range topology.Network.Bonds {
					v, ok := addressing[k]
					if !ok {
						continue
					}
					delete(addressing, k)
					from, fromOk := v.(map[string]interface{})
					to, toOk := topology.Network.Bonds[k].(map[string]interface{})
					if !fromOk || !toOk {
						continue
					}
					for fk := range from {
						to[fk] = from[fk]
					}
					topology.Network.Bonds[k] = to
				}
			}
			if topology.Network.Bridges != nil {
				for k := range topology.Network.Bridges {
					v, ok := addressing[k]
					if !ok {
						continue
					}
					delete(addressing, k)
					from, fromOk := v.(map[string]interface{})
					to, toOk := topology.Network.Bridges[k].(map[string]interface{})
					if !fromOk || !toOk {
						continue
					}
					for fk := range from {
						to[fk] = from[fk]
					}
					topology.Network.Bridges[k] = to
				}
			}
			if topology.Network.Ethernets == nil {
				topology.Network.Ethernets = map[string]interface{}{}
			}
			for k := range topology.Network.Ethernets {
				v, ok := addressing[k]
				if !ok {
					continue
				}
				delete(addressing, k)
				from, fromOk := v.(map[string]interface{})
				to, toOk := topology.Network.Ethernets[k].(map[string]interface{})
				if !fromOk || !toOk {
					continue
				}
				for fk := range from {
					to[fk] = from[fk]
				}
				topology.Network.Ethernets[k] = to
			}
			for k := range addressing {
				topology.Network.Ethernets[k] = addressing[k]
			}
			layout, err := topology.Compile(phys)
			if err != nil {
				return err
			}
			return netwrangler.Write(layout, args[0], args[4], bindMac)
		},
	}
	generate.Flags().BoolVar(&bindMac, "bindMac", true, "Bind all base nic definitions to their mac address. Defaults to true, otherwise name bindings will be used.")
	net.AddCommand(generate)
	app.AddCommand(net)
}
