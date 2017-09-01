package cli

import (
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"

	"github.com/digitalrebar/provision/backend"
	"github.com/digitalrebar/provision/client/machines"
	models "github.com/digitalrebar/provision/genmodels"
	"github.com/ghodss/yaml"
	"github.com/go-openapi/strfmt"
	"github.com/gorilla/websocket"
	"github.com/spf13/cobra"
)

type MachineOps struct{ CommonOps }

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
		case "Available":
			params = params.WithAvailable(&v)
		case "Valid":
			params = params.WithValid(&v)
		case "Name":
			params = params.WithName(&v)
		case "BootEnv":
			params = params.WithBootEnv(&v)
		case "UUID":
			params = params.WithUUID(&v)
		case "Address":
			params = params.WithAddress(&v)
		case "Runnable":
			params = params.WithRunnable(&v)
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
	a := machines.NewPatchMachineParams().WithUUID(strfmt.UUID(id)).WithBody(data)
	if force {
		s := "true"
		a = a.WithForce(&s)
	}
	d, e := session.Machines.PatchMachine(a, basicAuth)
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

func (be MachineOps) DoWait(id, field, value string, timeout int64) (string, error) {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	defer signal.Reset(os.Interrupt)

	u := url.URL{Scheme: "wss", Host: strings.TrimPrefix(endpoint, "https://"), Path: "/api/v3/ws"}
	// Set up auth stuff.
	websocket.DefaultDialer.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	h := http.Header(make(map[string][]string))
	h["Authorization"] = []string{"Basic " + base64.StdEncoding.EncodeToString([]byte(username+":"+password))}
	if token != "" {
		h["Authorization"] = []string{"Bearer " + token}
	}

	// Get the socket.
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), h)
	if err != nil {
		return "", err
	}
	defer conn.Close()

	done := make(chan struct{})
	var machine *models.Machine
	answer := ""

	go func() {
		defer close(done)
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
					fmt.Println("read:", err)
				}
				return
			}

			data := map[string]interface{}{}
			err = json.Unmarshal(message, &data)
			if err != nil {
				fmt.Println("json error: ", err)
				return
			}

			if md, ok := data["Object"].(map[string]interface{}); ok {
				if m, e := testMachine(field, value, md); e == nil && m {
					answer = "complete"
					return
				}
			}
		}
	}()

	// register for events.
	do_wait := true
	err = conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("register machines.update.%s\n", id)))
	if err != nil {
		do_wait = false
	}
	err = conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("register machines.save.%s\n", id)))
	if err != nil {
		do_wait = false
	}

	data, err := Get(id, be)
	if err != nil {
		err = generateError(err, "Failed to fetch %v: %v", be.GetSingularName(), id)
		do_wait = false
	}
	machine, _ = data.(*models.Machine)

	if do_wait {
		var jstring []byte
		jstring, err = json.Marshal(machine)
		if err == nil {
			data := map[string]interface{}{}
			err = json.Unmarshal(jstring, &data)
			if err == nil {
				var matched bool
				if matched, err = testMachine(field, value, data); err == nil && matched {
					answer = "complete"
					do_wait = false
				} else if err != nil {
					do_wait = false
				}
			}
		}

	}

	if do_wait {
		timer := time.NewTimer(time.Second * time.Duration(timeout))
		// Wait for reader to close (error, closed connection, or matched field)
		// Wait for timeout
		// Wait for user interrupt
		select {
		case <-done:
		case <-timer.C:
			answer = "timeout"
		case <-interrupt:
			answer = "interrupt"
		}
		timer.Stop()
	}

	// To cleanly close a connection, a client should send a close
	// frame and wait for the server to close the connection.
	terr := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if terr != nil && err == nil {
		err = terr
	}
	timer := time.NewTimer(time.Second)
	select {
	case <-done:
	case <-timer.C:
	}
	timer.Stop()

	return answer, err
}

func init() {
	tree := addMachineCommands()
	App.AddCommand(tree)
}

func testMachine(field, value string, fields map[string]interface{}) (bool, error) {
	var err error
	matched := false

	if d, ok := fields[field]; ok {
		switch v := d.(type) {
		case bool:
			var bval bool
			bval, err = strconv.ParseBool(value)
			if err == nil {
				if v == bval {
					matched = true
				}
			}
		case string:
			if v == value {
				matched = true
			}
		case int:
			var ival int64
			ival, err = strconv.ParseInt(value, 10, 64)
			if err == nil {
				if int(ival) == v {
					matched = true
				}
			}
		default:
			err = fmt.Errorf("Unsupported field type: %T\n", d)
		}
	}
	return matched, err
}

func addMachineCommands() (res *cobra.Command) {
	singularName := "machine"
	name := "machines"
	d("Making command tree for %v\n", name)
	res = &cobra.Command{
		Use:   name,
		Short: fmt.Sprintf("Access CLI commands relating to %v", name),
	}

	mo := &MachineOps{CommonOps{Name: name, SingularName: singularName}}

	commands := commonOps(mo)

	commands = append(commands, &cobra.Command{
		Use:   "wait [id] [field] [value] [timeout]",
		Short: fmt.Sprintf("Wait for a machine's field to become a value within a number of seconds"),
		Long: `
This function opens a web socket, registers for machine events, and waits for the value to become the new value.

Timeout is optional, defaults to indefinite, and is measured in seconds.

Returns the following strings:
  complete - field is equal to value
  interrupt - user interrupted the command
  timeout - timeout has exceeded
		`,
		RunE: func(c *cobra.Command, args []string) error {
			if len(args) < 3 {
				return fmt.Errorf("%v requires at least 3 arguments", c.UseLine())
			}
			if len(args) > 4 {
				return fmt.Errorf("%v requires at most 4 arguments", c.UseLine())
			}
			dumpUsage = false

			id := args[0]
			field := args[1]
			value := args[2]
			timeout := int64(100000000)
			if len(args) == 4 {
				var e error
				if timeout, e = strconv.ParseInt(args[3], 10, 64); e != nil {
					return e
				}
			}

			answer, err := mo.DoWait(id, field, value, timeout)
			if answer != "" {
				fmt.Println(answer)
			}
			return err
		},
	})

	commands = append(commands, &cobra.Command{
		Use:   "stage [id] [stage]",
		Short: fmt.Sprintf("Set the machine's stage"),
		Long:  `Helper function to update the machine's stage.`,
		RunE: func(c *cobra.Command, args []string) error {
			if len(args) != 2 {
				return fmt.Errorf("%v requires 2 arguments", c.UseLine())
			}
			dumpUsage = false

			return PatchWithString(args[0], "{ \"Stage\": \""+args[1]+"\" }", mo)
		},
	})

	commands = append(commands, &cobra.Command{
		Use:   "bootenv [id] [bootenv]",
		Short: fmt.Sprintf("Set the machine's bootenv"),
		Long:  `Helper function to update the machine's bootenv.`,
		RunE: func(c *cobra.Command, args []string) error {
			if len(args) != 2 {
				return fmt.Errorf("%v requires 2 arguments", c.UseLine())
			}
			dumpUsage = false
			return PatchWithString(args[0], "{ \"BootEnv\": \""+args[1]+"\" }", mo)
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
			return PatchWithFunction(args[0], mo, func(data interface{}) (interface{}, bool) {
				machine, _ := data.(*models.Machine)
				machine.Profiles = append(machine.Profiles, args[1])
				return machine, true
			})
		},
	})

	commands = append(commands, &cobra.Command{
		Use:   "removeprofile [id] [profile]",
		Short: fmt.Sprintf("Remove a profile from the machine's list"),
		Long:  `Helper function to update the machine's profile list by removing one.`,
		RunE: func(c *cobra.Command, args []string) error {
			if len(args) != 2 {
				return fmt.Errorf("%v requires 2 arguments", c.UseLine())
			}
			dumpUsage = false
			return PatchWithFunction(args[0], mo, func(data interface{}) (interface{}, bool) {
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
				return machine, changed
			})

		},
	})

	aggregate := false
	getParams := &cobra.Command{
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
				as := "false"
				if aggregate {
					as = "true"
				}
				d, err := session.Machines.GetMachineParams(machines.NewGetMachineParamsParams().WithAggregate(&as).WithUUID(strfmt.UUID(uuid)), basicAuth)
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
	}
	commands = append(commands, getParams)
	getParams.Flags().BoolVar(&aggregate, "aggregate", false, "Should machine return aggregated view")

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
				return fmt.Errorf("%v requires 2 arguments", c.UseLine())
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
				return fmt.Errorf("runaction either takes three arguments or a multiple of two, not %d", len(args))
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

	commands = append(commands, processJobsCommand())

	res.AddCommand(commands...)
	return res
}
