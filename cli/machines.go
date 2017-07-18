package cli

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/VictorLowther/jsonpatch2"
	"github.com/VictorLowther/jsonpatch2/utils"
	"github.com/digitalrebar/provision/backend"
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

func (be MachineOps) GetIndexes() map[string]string {
	b := &backend.Machine{}
	ans := map[string]string{}
	for k, v := range b.Indexes() {
		ans[k] = v.Type
	}
	return ans
}

func (be MachineOps) List(parms map[string]string) (interface{}, error) {
	params := machines.NewListMachinesParams()
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
		case "Name":
			params = params.WithName(&v)
		case "BootEnv":
			params = params.WithBootEnv(&v)
		case "UUID":
			params = params.WithUUID(&v)
		case "Address":
			params = params.WithAddress(&v)
		}
	}
	d, e := session.Machines.ListMachines(params, basicAuth)
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
		name, ok := obj.(string)
		if !ok {
			return nil, fmt.Errorf("Invalid type passed to machine create")
		}
		hostname := strfmt.Hostname(name)
		machine = &models.Machine{Name: &hostname}
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
		Use:   "addprofile [id] [profile]",
		Short: fmt.Sprintf("Add profile to the machine's profile list"),
		Long:  `Helper function to add a profile to the machine's profile list.`,
		RunE: func(c *cobra.Command, args []string) error {
			if len(args) != 2 {
				return fmt.Errorf("%v requires 2 arguments", c.UseLine())
			}
			dumpUsage = false
			data, err := mo.Get(args[0])
			if err != nil {
				return generateError(err, "Failed to fetch %v: %v", singularName, args[0])
			}

			baseObj, err := json.Marshal(data)
			if err != nil {
				return fmt.Errorf("Unable to marshal object: %v\n", err)
			}

			machine, _ := data.(*models.Machine)
			machine.Profiles = append(machine.Profiles, args[1])
			merged, err := json.Marshal(machine)
			if err != nil {
				return fmt.Errorf("Unable to marshal input stream: %v\n", err)
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
		Use:   "removeprofile [id] [profile]",
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

			baseObj, err := json.Marshal(data)
			if err != nil {
				return fmt.Errorf("Unable to marshal object: %v\n", err)
			}

			changed := false
			machine, _ := data.(*models.Machine)
			newProfiles := []string{}
			for _, s := range machine.Profiles {
				if s == args[1] {
					changed = true
					continue
				}
				newProfiles = append(newProfiles, s)
			}
			machine.Profiles = newProfiles
			if len(newProfiles) == 0 {
				machine.Profiles = nil
			}

			if !changed {
				return prettyPrint(data)
			}

			merged, err := json.Marshal(data)
			if err != nil {
				return fmt.Errorf("Unable to marshal input stream: %v\n", err)
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

	commands = append(commands, &cobra.Command{
		Use:   "actions [id]",
		Short: fmt.Sprintf("Display actions for this machine"),
		Long:  `Helper function to display the machine's actions.`,
		RunE: func(c *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("%v requires 1 argument", c.UseLine())
			}
			uuid := args[0]
			dumpUsage = false

			d, err := session.Machines.GetMachineActions(machines.NewGetMachineActionsParams().WithUUID(strfmt.UUID(uuid)), basicAuth)
			if err != nil {
				return generateError(err, "Failed to fetch actions %v: %v", singularName, uuid)
			}
			pp := d.Payload
			return prettyPrint(pp)
		},
	})
	commands = append(commands, &cobra.Command{
		Use:   "action [id] [action]",
		Short: fmt.Sprintf("Display the action for this machine"),
		Long:  `Helper function to display the machine's action.`,
		RunE: func(c *cobra.Command, args []string) error {
			if len(args) != 2 {
				return fmt.Errorf("%v requires 1 argument", c.UseLine())
			}
			uuid := args[0]
			action := args[1]
			dumpUsage = false

			d, err := session.Machines.GetMachineAction(machines.NewGetMachineActionParams().WithUUID(strfmt.UUID(uuid)).WithName(action), basicAuth)
			if err != nil {
				return generateError(err, "Failed to fetch action %v: %v %v", singularName, uuid, action)
			}
			pp := d.Payload
			return prettyPrint(pp)
		},
	})

	commands = append(commands, &cobra.Command{
		Use:   "runaction [id] [command] [- | JSON or YAML Map of objects | pairs of string objects]",
		Short: "Set preferences",
		RunE: func(c *cobra.Command, args []string) error {
			actionParams := map[string]interface{}{}
			if len(args) == 3 {
				var buf []byte
				var err error
				if args[2] == `-` {
					buf, err = ioutil.ReadAll(os.Stdin)
					if err != nil {
						dumpUsage = false
						return fmt.Errorf("Error reading from stdin: %v", err)
					}
				} else {
					buf = []byte(args[2])
				}
				err = yaml.Unmarshal(buf, &actionParams)
				if err != nil {
					dumpUsage = false
					return fmt.Errorf("Invalid parameters: %v\n", err)
				}
			} else if len(args) > 3 && len(args)%2 == 0 {
				for i := 2; i < len(args); i += 2 {
					var obj interface{}
					err := yaml.Unmarshal([]byte(args[i+1]), &obj)
					if err != nil {
						dumpUsage = false
						return fmt.Errorf("Invalid parameters: %s %v\n", args[i+1], err)
					}
					actionParams[args[i]] = obj
				}
			} else if len(args) < 2 || len(args)%2 == 1 {
				return fmt.Errorf("runaction either takes a single argument or a multiple of two, not %d", len(args))
			}
			uuid := args[0]
			command := args[1]
			dumpUsage = false
			if resp, err := session.Machines.PostMachineAction(machines.NewPostMachineActionParams().WithBody(actionParams).WithUUID(strfmt.UUID(uuid)).WithName(command), basicAuth); err != nil {
				return generateError(err, "Error running action")
			} else {
				return prettyPrint(resp)
			}
		},
	})

	res.AddCommand(commands...)
	return res
}
