package cli

import (
	"encoding/json"
	"fmt"

	"github.com/VictorLowther/jsonpatch2"
	"github.com/VictorLowther/jsonpatch2/utils"
	"github.com/digitalrebar/provision/client/machines"
	"github.com/digitalrebar/provision/models"
	"github.com/ghodss/yaml"
	"github.com/go-openapi/strfmt"
	"github.com/spf13/cobra"
)

type MachineOps struct{}

func (be MachineOps) GetType() interface{} {
	return &models.Machine{}
}

func (be MachineOps) GetId(obj interface{}) (string, error) {
	machine, ok := obj.(*models.Machine)
	if !ok || machine.UUID == nil {
		return "", fmt.Errorf("Invalid type passed to machine create")
	}
	return machine.UUID.String(), nil
}

func (be MachineOps) List() (interface{}, error) {
	d, e := session.Machines.ListMachines(machines.NewListMachinesParams(), basicAuth)
	if e != nil {
		return nil, e
	}
	return d.Payload, nil
}

func (be MachineOps) Get(id string) (interface{}, error) {
	d, e := session.Machines.GetMachine(machines.NewGetMachineParams().WithUUID(strfmt.UUID(id)), basicAuth)
	if e != nil {
		return nil, e
	}
	return d.Payload, nil
}

func (be MachineOps) Create(obj interface{}) (interface{}, error) {
	machine, ok := obj.(*models.Machine)
	if !ok {
		return nil, fmt.Errorf("Invalid type passed to machine create")
	}
	d, e := session.Machines.CreateMachine(machines.NewCreateMachineParams().WithBody(machine), basicAuth)
	if e != nil {
		return nil, e
	}
	return d.Payload, nil
}

func (be MachineOps) Patch(id string, obj interface{}) (interface{}, error) {
	data, ok := obj.(models.Patch)
	if !ok {
		return nil, fmt.Errorf("Invalid type passed to machine patch")
	}
	d, e := session.Machines.PatchMachine(machines.NewPatchMachineParams().WithUUID(strfmt.UUID(id)).WithBody(data), basicAuth)
	if e != nil {
		return nil, e
	}
	return d.Payload, nil
}

func (be MachineOps) Delete(id string) (interface{}, error) {
	d, e := session.Machines.DeleteMachine(machines.NewDeleteMachineParams().WithUUID(strfmt.UUID(id)), basicAuth)
	if e != nil {
		return nil, e
	}
	return d.Payload, nil
}

func init() {
	tree := addMachineCommands()
	App.AddCommand(tree)
}

func addMachineCommands() (res *cobra.Command) {
	singularName := "machine"
	name := "machines"
	d("Making command tree for %v\n", name)
	res = &cobra.Command{
		Use:   name,
		Short: fmt.Sprintf("Access CLI commands relating to %v", name),
	}

	mo := &MachineOps{}

	commands := commonOps(singularName, name, mo)

	commands = append(commands, &cobra.Command{
		Use:   "bootenv [id] [bootenv]",
		Short: fmt.Sprintf("Set the machine's bootenv"),
		Long:  `Helper function to update the machine's bootenv.`,
		RunE: func(c *cobra.Command, args []string) error {
			if len(args) != 2 {
				return fmt.Errorf("%v requires 2 arguments", c.UseLine())
			}
			dumpUsage = false
			data, err := mo.Get(args[0])
			if err != nil {
				return generateError(err, "Failed to fetch %v: %v", singularName, args[0])
			}
			var buf []byte

			baseObj, err := json.Marshal(data)
			if err != nil {
				return fmt.Errorf("Unable to marshal object: %v\n", err)
			}
			buf = []byte("{ \"BootEnv\": \"" + args[1] + "\" }")
			var intermediate interface{}
			err = yaml.Unmarshal(buf, &intermediate)
			if err != nil {
				return fmt.Errorf("Unable to unmarshal input stream: %v\n", err)
			}
			updateObj, err := json.Marshal(intermediate)
			if err != nil {
				return fmt.Errorf("Unable to marshal input stream: %v\n", err)
			}
			merged, err := safeMergeJSON(data, updateObj)
			if err != nil {
				return fmt.Errorf("Unable to merge objects: %v\n", err)
			}
			patch, err := jsonpatch2.Generate(baseObj, merged, true)
			if err != nil {
				return fmt.Errorf("Error generating patch: %v", err)
			}
			p := models.Patch{}
			if err := utils.Remarshal(&patch, &p); err != nil {
				return fmt.Errorf("Error translating patch for bootenv: %v", err)
			}

			if data, err := mo.Patch(args[0], p); err != nil {
				return generateError(err, "Unable to update bootenv %v", args[0])
			} else {
				return prettyPrint(data)
			}
		},
	})

	commands = append(commands, &cobra.Command{
		Use:   "params [id] [json]",
		Short: fmt.Sprintf("Gets/sets all parameters for the machine"),
		Long:  `A helper function to return all or set all the parameters on the machine`,
		RunE: func(c *cobra.Command, args []string) error {
			if len(args) != 1 && len(args) != 2 {
				return fmt.Errorf("%v requires 1 or 2 arguments", c.UseLine())
			}
			dumpUsage = false
			uuid := args[0]
			if len(args) == 1 {
				d, err := session.Machines.GetMachineParams(machines.NewGetMachineParamsParams().WithUUID(strfmt.UUID(uuid)), basicAuth)
				if err != nil {
					return generateError(err, "Failed to fetch params %v: %v", singularName, uuid)
				}
				return prettyPrint(d.Payload)
			} else {
				newValue := args[1]
				var value map[string]interface{}
				err := yaml.Unmarshal([]byte(newValue), &value)
				if err != nil {
					return fmt.Errorf("Unable to unmarshal input stream: %v\n", err)
				}
				d, err := session.Machines.PostMachineParams(machines.NewPostMachineParamsParams().WithUUID(strfmt.UUID(uuid)).WithBody(value), basicAuth)
				if err != nil {
					return generateError(err, "Failed to fetch params %v: %v", singularName, uuid)
				}
				return prettyPrint(d.Payload)
			}
		},
	})

	commands = append(commands, &cobra.Command{
		Use:   "get [id] param [key]",
		Short: fmt.Sprintf("Get a parameter from the machine"),
		Long:  `A helper function to return the value of the parameter on the machine`,
		RunE: func(c *cobra.Command, args []string) error {
			if len(args) != 3 {
				return fmt.Errorf("%v requires 3 arguments", c.UseLine())
			}
			dumpUsage = false
			uuid := args[0]
			// at = args[1]
			key := args[2]

			d, err := session.Machines.GetMachineParams(machines.NewGetMachineParamsParams().WithUUID(strfmt.UUID(uuid)), basicAuth)
			if err != nil {
				return generateError(err, "Failed to fetch params %v: %v", singularName, uuid)
			}
			pp := d.Payload
			if pp == nil {
				return prettyPrint(pp)
			}

			if val, found := pp[key]; found {
				return prettyPrint(val)
			} else {
				return prettyPrint(nil)
			}
		},
	})

	commands = append(commands, &cobra.Command{
		Use:   "set [id] param [key] to [json blob]",
		Short: fmt.Sprintf("Set the machine's param <key> to <blob>"),
		Long:  `Helper function to update the machine's parameters.`,
		RunE: func(c *cobra.Command, args []string) error {
			if len(args) != 5 {
				return fmt.Errorf("%v requires 5 arguments", c.UseLine())
			}
			uuid := args[0]
			key := args[2]
			newValue := args[4]
			dumpUsage = false

			var value interface{}
			err := yaml.Unmarshal([]byte(newValue), &value)
			if err != nil {
				return fmt.Errorf("Unable to unmarshal input stream: %v\n", err)
			}

			d, err := session.Machines.GetMachineParams(machines.NewGetMachineParamsParams().WithUUID(strfmt.UUID(uuid)), basicAuth)
			if err != nil {
				return generateError(err, "Failed to fetch params %v: %v", singularName, uuid)
			}
			pp := d.Payload
			if value == nil {
				delete(pp, key)
			} else {
				pp[key] = value
			}
			_, err = session.Machines.PostMachineParams(machines.NewPostMachineParamsParams().WithUUID(strfmt.UUID(uuid)).WithBody(pp), basicAuth)
			if err != nil {
				return generateError(err, "Failed to post params %v: %v", singularName, uuid)
			}
			return prettyPrint(value)
		},
	})

	res.AddCommand(commands...)
	return res
}
