package cli

import (
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/digitalrebar/provision/models"
	"github.com/spf13/cobra"
)

func init() {
	addRegistrar(registerSubnet)
}

func registerSubnet(app *cobra.Command) {
	op := &ops{
		name:       "subnets",
		singleName: "subnet",
		example:    func() models.Model { return &models.Subnet{} },
	}
	op.addCommand(&cobra.Command{
		Use:   "range [subnetName] [startIP] [endIP]",
		Short: fmt.Sprintf("set the range of a subnet"),
		Long:  `Helper function to set the range of a given subnet.`,
		Args: func(c *cobra.Command, args []string) error {
			if len(args) != 3 {
				return fmt.Errorf("%s requires 3 arguments", c.UseLine())
			}
			return nil
		},
		RunE: func(c *cobra.Command, args []string) error {
			ipFirst, ipLast := net.ParseIP(args[1]), net.ParseIP(args[2])
			if ipFirst == nil || ipFirst.To4() == nil {
				return fmt.Errorf("%s is not a valid IPv4", args[1])
			}
			if ipLast == nil || ipLast.To4() == nil {
				return fmt.Errorf("%s is not a valid IPv4", args[1])
			}
			return PatchWithFunction(args[0], op, func(data models.Model) (models.Model, bool) {
				sub := data.(*models.Subnet)
				sub.ActiveStart = ipFirst
				sub.ActiveEnd = ipLast
				return sub, true
			})
		},
	})
	op.addCommand(&cobra.Command{
		Use:        "subnet [subnetName] [subnet CIDR]",
		Short:      fmt.Sprintf("Set the CIDR network address"),
		Deprecated: "Changing the subnet CIDR address is not supported.",
		Long:       `Helper function to set the CIDR of a given subnet.`,
		Args: func(c *cobra.Command, args []string) error {
			if len(args) != 2 {
				return fmt.Errorf("%s requires 2 arguments", c.UseLine())
			}
			return nil
		},
		RunE: func(c *cobra.Command, args []string) error {
			cidr := args[1]
			if _, _, e2 := net.ParseCIDR(cidr); e2 != nil {
				return fmt.Errorf("%s is not a valid subnet CIDR", cidr)

			}
			return PatchWithString(args[0], "{\"Subnet\": \""+cidr+"\"}", op)
		},
	})
	/* Save for when we have extra strategies other than MAC */
	/*
		op.addCommand(&cobra.Command{
			Use:   "strategy [subnetName] [MAC]",
			Short: fmt.Sprintf("Set Subnet strategy"),
			Long:  `Helper function to set the strategy of a given subnet.`,
			Args: func(c *cobra.Command, args []string) error {
				if len(args) != 2 {
					return fmt.Errorf("%s requires 2 arguments", c.UseLine())
				}
				return nil
			},
			RunE: func(c *cobra.Command, args []string) error {
				return PatchWithString(args[0], "{\"Strategy\": \""+args[1]+"\"}", op)
			},
		})
	*/
	op.addCommand(&cobra.Command{
		Use:   "pickers [subnetName] [list]",
		Short: fmt.Sprintf("assigns IP allocation methods to a subnet"),
		Long:  `Helper function that accepts a string of methods to allocate IP addresses separated by commas`,
		Args: func(c *cobra.Command, args []string) error {
			if len(args) != 2 {
				return fmt.Errorf("%v requires 2 arguments", c.UseLine())
			}
			return nil
		},
		RunE: func(c *cobra.Command, args []string) error {
			pickerString := args[1]
			return PatchWithFunction(args[0], op, func(data models.Model) (models.Model, bool) {
				sub := data.(*models.Subnet)
				sub.Pickers = strings.Split(pickerString, ",")
				return sub, true
			})
		},
	})

	op.addCommand(&cobra.Command{
		Use:   "nextserver [subnetName] [IP]",
		Short: fmt.Sprintf("Set next non-reserved IP"),
		Long:  `Helper function to set the first non-reserved IP of a subnet.`,
		Args: func(c *cobra.Command, args []string) error {
			if len(args) != 2 {
				return fmt.Errorf("%v requires 2 arguments", c.UseLine())
			}
			return nil
		},
		RunE: func(c *cobra.Command, args []string) error {
			addr := net.ParseIP(args[1])
			if addr == nil || addr.To4() == nil {
				return fmt.Errorf("%s is not a valid IPv4", args[1])
			}
			return PatchWithFunction(args[0], op, func(data models.Model) (models.Model, bool) {
				sub := data.(*models.Subnet)
				sub.NextServer = addr
				return sub, true
			})
		},
	})

	op.addCommand(&cobra.Command{
		Use:   "leasetimes [subnetName] [active] [reserved]",
		Short: fmt.Sprintf("Set the leasetimes of a subnet"),
		Long:  `Helper function to get the range of a given subnet.`,
		Args: func(c *cobra.Command, args []string) error {
			if len(args) != 3 {
				return fmt.Errorf("%v requires 3 arguments", c.UseLine())
			}
			return nil
		},
		RunE: func(c *cobra.Command, args []string) error {
			activeTimeString := args[1]
			reservedTimeString := args[2]

			activeTime64Int, e := strconv.ParseInt(activeTimeString, 10, 32)
			if e != nil {
				return fmt.Errorf("%v could not be read as a number", activeTimeString)
			}
			activeTime := int32(activeTime64Int)

			reservedTime64Int, e := strconv.ParseInt(reservedTimeString, 10, 32)
			if e != nil {
				return fmt.Errorf("%v could not be read as a number", reservedTimeString)
			}
			reservedTime := int32(reservedTime64Int)

			return PatchWithFunction(args[0], op, func(data models.Model) (models.Model, bool) {
				sub := data.(*models.Subnet)
				sub.ActiveLeaseTime = activeTime
				sub.ReservedLeaseTime = reservedTime
				return sub, true
			})
		},
	})

	op.addCommand(&cobra.Command{
		Use:   "set [subnetName] option [number] to [value]",
		Short: fmt.Sprintf("Set the given subnet's dhcpOption to a value"),
		Long:  `Helper function that sets the specified dhcpOption from a given subnet to a value. If an option does not exist yet, it adds a new option`,
		Args: func(c *cobra.Command, args []string) error {
			if len(args) != 5 {
				return fmt.Errorf("%v requires 5 arguments", c.UseLine())
			}
			return nil
		},
		RunE: func(c *cobra.Command, args []string) error {
			optionNumberString := args[2]
			newVal := args[4]

			gv, err := strconv.Atoi(optionNumberString)
			if err != nil {
				return fmt.Errorf("%v could not be read as a number", optionNumberString)
			}
			optionNumber := byte(gv)

			return PatchWithFunction(args[0], op, func(data models.Model) (models.Model, bool) {
				sub := data.(*models.Subnet)
				found := false
				if sub.Options == nil {
					sub.Options = []models.DhcpOption{}
				}
				idx := -1
				for ii, do := range sub.Options {
					if do.Code == optionNumber {
						if newVal == "null" {
							idx = ii
						} else {
							sub.Options[ii].Value = newVal
						}
						found = true
						break
					}
				}
				if idx != -1 {
					sub.Options = append(sub.Options[:idx], sub.Options[idx+1:]...)
				}
				if !found {
					newOption := models.DhcpOption{Code: optionNumber, Value: newVal}
					sub.Options = append(sub.Options, newOption)
				}
				return sub, true
			})
		},
	})

	op.addCommand(&cobra.Command{
		Use:   "get [subnetName] option [number]",
		Short: fmt.Sprintf("Get dhcpOption [number]"),
		Long:  `Helper function that gets the specified dhcpOption from a given subnet.`,
		Args: func(c *cobra.Command, args []string) error {
			if len(args) != 3 {
				return fmt.Errorf("%v requires 3 arguments", c.UseLine())
			}
			return nil
		},
		RunE: func(c *cobra.Command, args []string) error {
			subName := args[0]
			gettingVal := args[2]

			gv, err := strconv.Atoi(gettingVal)
			if err != nil {
				return fmt.Errorf("%v could not be read as a number", gettingVal)
			}
			getVal := byte(gv)
			sub := &models.Subnet{}
			if e := session.FillModel(sub, subName); e != nil {
				return e
			}

			for _, do := range sub.Options {
				if do.Code == getVal {
					fmt.Printf("Option %v: %v\n", getVal, do.Value)
					return nil
				}
			}

			return fmt.Errorf("option %v does not exist", getVal)
		},
	})
	op.command(app)
}
