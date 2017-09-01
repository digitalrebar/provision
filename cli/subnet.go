package cli

import (
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/digitalrebar/provision/backend"
	"github.com/digitalrebar/provision/client/subnets"
	models "github.com/digitalrebar/provision/genmodels"
	"github.com/go-openapi/strfmt"
	"github.com/spf13/cobra"
)

type SubnetOps struct{ CommonOps }

func (be SubnetOps) GetType() interface{} {
	return &models.Subnet{}
}

func (be SubnetOps) GetId(obj interface{}) (string, error) {
	subnet, ok := obj.(*models.Subnet)
	if !ok {
		return "", fmt.Errorf("Invalid type passed to subnet create")
	}
	return *subnet.Name, nil
}

func (be SubnetOps) GetIndexes() map[string]string {
	b := &backend.Subnet{}
	ans := map[string]string{}
	for k, v := range b.Indexes() {
		ans[k] = v.Type
	}
	return ans
}

func (be SubnetOps) List(parms map[string]string) (interface{}, error) {
	params := subnets.NewListSubnetsParams()
	if listLimit != -1 {
		t1 := int64(listLimit)
		params = params.WithLimit(&t1)
	}
	if listOffset != -1 {
		t1 := int64(listOffset)
		params = params.WithOffset(&t1)
	}
	for k, v := range parms {
		switch k {
		case "Available":
			params = params.WithAvailable(&v)
		case "Valid":
			params = params.WithValid(&v)
		case "Name":
			params = params.WithName(&v)
		case "Enabled":
			params = params.WithEnabled(&v)
		case "Subnet":
			params = params.WithSubnet(&v)
		case "Strategy":
			params = params.WithStrategy(&v)
		case "NextServer":
			params = params.WithNextServer(&v)
		}
	}

	d, e := session.Subnets.ListSubnets(params, basicAuth)
	if e != nil {
		return nil, e
	}
	return d.Payload, nil
}

func (be SubnetOps) Get(id string) (interface{}, error) {
	d, e := session.Subnets.GetSubnet(subnets.NewGetSubnetParams().WithName(id), basicAuth)
	if e != nil {
		return nil, e
	}
	return d.Payload, nil
}

func (be SubnetOps) Create(obj interface{}) (interface{}, error) {
	subnet, ok := obj.(*models.Subnet)
	if !ok {
		return nil, fmt.Errorf("Invalid type passed to subnet create")
	}
	d, e := session.Subnets.CreateSubnet(subnets.NewCreateSubnetParams().WithBody(subnet), basicAuth)
	if e != nil {
		return nil, e
	}
	return d.Payload, nil
}

func (be SubnetOps) Patch(id string, obj interface{}) (interface{}, error) {
	data, ok := obj.(models.Patch)
	if !ok {
		return nil, fmt.Errorf("Invalid type passed to subnet patch")
	}
	d, e := session.Subnets.PatchSubnet(subnets.NewPatchSubnetParams().WithName(id).WithBody(data), basicAuth)
	if e != nil {
		return nil, e
	}
	return d.Payload, nil
}

func (be SubnetOps) Delete(id string) (interface{}, error) {
	d, e := session.Subnets.DeleteSubnet(subnets.NewDeleteSubnetParams().WithName(id), basicAuth)
	if e != nil {
		return nil, e
	}
	return d.Payload, nil
}

func init() {
	tree := addSubnetCommands()
	App.AddCommand(tree)
}

func addSubnetCommands() (res *cobra.Command) {
	singularName := "subnet"
	name := "subnets"
	d("Making command tree for %v\n", name)
	res = &cobra.Command{
		Use:   name,
		Short: fmt.Sprintf("Access CLI commands relating to %v", name),
	}
	op := &SubnetOps{CommonOps{Name: name, SingularName: singularName}}
	commands := commonOps(op)

	commands = append(commands, &cobra.Command{
		Use:   "range [subnetName] [startIP] [endIP]",
		Short: fmt.Sprintf("set the range of a subnet"),
		Long:  `Helper function to set the range of a given subnet.`,
		RunE: func(c *cobra.Command, args []string) error {
			if len(args) != 3 {
				return fmt.Errorf("%s requires 3 arguments", c.UseLine())
			}
			dumpUsage = false
			StartAddr := args[1]
			EndAddr := args[2]

			var IPfirst strfmt.IPv4
			if e := IPfirst.Scan(StartAddr); e != nil {
				return fmt.Errorf("%s is not a valid IPv4", StartAddr)
			}

			var IPlast strfmt.IPv4
			if e := IPlast.Scan(EndAddr); e != nil {
				return fmt.Errorf("%s is not a valid IPv4", EndAddr)
			}

			return PatchWithFunction(args[0], op, func(data interface{}) (interface{}, bool) {
				sub := data.(*models.Subnet)
				sub.ActiveStart = &IPfirst
				sub.ActiveEnd = &IPlast
				return sub, true
			})
		},
	})

	commands = append(commands, &cobra.Command{
		Use:   "subnet [subnetName] [subnet CIDR]",
		Short: fmt.Sprintf("Set the CIDR network address"),
		Long:  `Helper function to set the CIDR of a given subnet.`,
		RunE: func(c *cobra.Command, args []string) error {
			if len(args) != 2 {
				return fmt.Errorf("%s requires 2 arguments", c.UseLine())
			}
			dumpUsage = false
			CIDR := args[1]
			if _, _, e2 := net.ParseCIDR(CIDR); e2 != nil {
				return fmt.Errorf("%s is not a valid subnet CIDR", CIDR)

			}
			return PatchWithString(args[0], "{\"Subnet\": \""+CIDR+"\"}", op)
		},
	})

	commands = append(commands, &cobra.Command{
		Use:   "strategy [subnetName] [MAC]",
		Short: fmt.Sprintf("Set Subnet strategy"),
		Long:  `Helper function to set the strategy of a given subnet.`,
		RunE: func(c *cobra.Command, args []string) error {
			if len(args) != 2 {
				return fmt.Errorf("%s requires 2 arguments", c.UseLine())
			}
			dumpUsage = false
			MACAddress := args[1]
			if _, e := net.ParseMAC(MACAddress); e != nil {
				return fmt.Errorf("%s is not a valid MAC address", MACAddress)
			}
			return PatchWithString(args[0], "{\"Strategy\": \""+MACAddress+"\"}", op)
		},
	})

	commands = append(commands, &cobra.Command{
		Use:   "pickers [subnetName] [list]",
		Short: fmt.Sprintf("assigns IP allocation methods to a subnet"),
		Long:  `Helper function that accepts a string of methods to allocate IP addresses separated by commas`,
		RunE: func(c *cobra.Command, args []string) error {
			if len(args) != 2 {
				return fmt.Errorf("%v requires 2 arguments", c.UseLine())
			}
			dumpUsage = false
			pickerString := args[1]
			return PatchWithFunction(args[0], op, func(data interface{}) (interface{}, bool) {
				sub := data.(*models.Subnet)
				sub.Pickers = strings.Split(pickerString, ",")
				return sub, true
			})

		},
	})

	commands = append(commands, &cobra.Command{
		Use:   "nextserver [subnetName] [IP]",
		Short: fmt.Sprintf("Set next non-reserved IP"),
		Long:  `Helper function to set the first non-reserved IP of a subnet.`,
		RunE: func(c *cobra.Command, args []string) error {
			if len(args) != 2 {
				return fmt.Errorf("%v requires 2 arguments", c.UseLine())
			}
			dumpUsage = false
			IPAddr := args[1]

			var nextIP strfmt.IPv4
			if e := nextIP.Scan(IPAddr); e != nil {
				return fmt.Errorf("%v is not a valid IPv4", IPAddr)
			}

			return PatchWithFunction(args[0], op, func(data interface{}) (interface{}, bool) {
				sub := data.(*models.Subnet)
				sub.NextServer = &nextIP
				return sub, true
			})
		},
	})

	commands = append(commands, &cobra.Command{
		Use:   "leasetimes [subnetName] [active] [reserved]",
		Short: fmt.Sprintf("Set the leasetimes of a subnet"),
		Long:  `Helper function to get the range of a given subnet.`,
		RunE: func(c *cobra.Command, args []string) error {
			if len(args) != 3 {
				return fmt.Errorf("%v requires 3 arguments", c.UseLine())
			}
			dumpUsage = false
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

			return PatchWithFunction(args[0], op, func(data interface{}) (interface{}, bool) {
				sub := data.(*models.Subnet)
				sub.ActiveLeaseTime = &activeTime
				sub.ReservedLeaseTime = &reservedTime
				return sub, true
			})
		},
	})

	commands = append(commands, &cobra.Command{
		Use:   "set [subnetName] option [number] to [value]",
		Short: fmt.Sprintf("Set the given subnet's dhcpOption to a value"),
		Long:  `Helper function that sets the specified dhcpOption from a given subnet to a value. If an option does not exist yet, it adds a new option`,
		RunE: func(c *cobra.Command, args []string) error {
			if len(args) != 5 {
				return fmt.Errorf("%v requires 5 arguments", c.UseLine())
			}
			dumpUsage = false
			optionNumberString := args[2]
			newVal := args[4]

			gv, err := strconv.Atoi(optionNumberString)
			if err != nil {
				return fmt.Errorf("%v could not be read as a number", optionNumberString)
			}
			optionNumber := models.OptionCode(gv)

			return PatchWithFunction(args[0], op, func(data interface{}) (interface{}, bool) {
				sub := data.(*models.Subnet)
				found := false
				if sub.Options == nil {
					sub.Options = []*models.DhcpOption{}
				}
				idx := -1
				for ii, do := range sub.Options {
					if do.Code == optionNumber {
						if newVal == "null" {
							idx = ii
						} else {
							do.Value = &newVal
						}
						found = true
						break
					}
				}
				if idx != -1 {
					sub.Options = append(sub.Options[:idx], sub.Options[idx+1:]...)
				}
				if !found {
					newOption := &models.DhcpOption{Code: optionNumber, Value: &newVal}
					sub.Options = append(sub.Options, newOption)
				}
				return sub, true
			})
		},
	})

	commands = append(commands, &cobra.Command{
		Use:   "get [subnetName] option [number]",
		Short: fmt.Sprintf("Get dhcpOption [number]"),
		Long:  `Helper function that gets the specified dhcpOption from a given subnet.`,
		RunE: func(c *cobra.Command, args []string) error {
			if len(args) != 3 {
				return fmt.Errorf("%v requires 3 arguments", c.UseLine())
			}
			dumpUsage = false
			subName := args[0]
			gettingVal := args[2]

			gv, err := strconv.Atoi(gettingVal)
			if err != nil {
				return fmt.Errorf("%v could not be read as a number", gettingVal)
			}
			getVal := models.OptionCode(gv)

			d, e := session.Subnets.GetSubnet(subnets.NewGetSubnetParams().WithName(subName), basicAuth)
			if e != nil {
				return e
			}
			sub := d.Payload

			for _, do := range sub.Options {
				if do.Code == getVal {
					fmt.Printf("Option %v: %v\n", getVal, *do.Value)
					return nil
				}
			}

			return fmt.Errorf("option %v does not exist", getVal)
		},
	})

	res.AddCommand(commands...)
	return res
}
