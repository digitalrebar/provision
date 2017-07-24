package main

import (
	"bytes"
	"fmt"
	"io"
	"net/rpc"
	"os"
	"strconv"
	"strings"

	"github.com/digitalrebar/provision"
	"github.com/digitalrebar/provision/backend"
	"github.com/digitalrebar/provision/cli"
	"github.com/digitalrebar/provision/midlayer"
	"github.com/digitalrebar/provision/plugin"
)

var (
	version = provision.RS_VERSION
	def     = midlayer.PluginProvider{
		Name:       "incrementer",
		Version:    version,
		HasPublish: false,
		AvailableActions: []*midlayer.AvailableAction{
			&midlayer.AvailableAction{Command: "increment",
				RequiredParams: []string{"incrementer.parameter"},
				OptionalParams: []string{"incrementer.step"},
			},
		},
		Parameters: []*backend.Param{
			&backend.Param{Name: "incrementer.parameter", Schema: map[string]interface{}{"type": "string"}},
			&backend.Param{Name: "incrementer.step", Schema: map[string]interface{}{"type": "integer"}},
		},
	}
)

type Plugin struct {
}

func (p *Plugin) Config(config map[string]interface{}, err *backend.Error) error {
	*err = backend.Error{Code: 0, Model: "plugin", Key: "incrementer", Type: "rpc", Messages: []string{}}
	plugin.Log("Config: %v\n", config)
	return nil
}

func executeDrpCliCommand(args ...string) (string, error) {
	old := os.Stdout // keep backup of the real stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	outC := make(chan string)
	// copy the output in a separate goroutine so printing can't block indefinitely
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		outC <- buf.String()
	}()

	cli.App.SetArgs(args)
	cli.App.SetOutput(os.Stderr)
	err2 := cli.App.Execute()

	// back to normal state
	w.Close()
	os.Stdout = old // restoring the real stdout
	out := <-outC

	plugin.Log("DrpCli: %s\nerr: %v\n", out, err2)

	return out, err2
}

func (p *Plugin) Action(ma midlayer.MachineAction, err *backend.Error) error {
	plugin.Log("Action: %v\n", ma)
	if ma.Command == "increment" {
		parameter, ok := ma.Params["incrementer.parameter"].(string)
		if !ok {
			*err = backend.Error{Code: 404,
				Model:    "plugin",
				Key:      "incrementer",
				Type:     "rpc",
				Messages: []string{fmt.Sprintf("Parameter is not specified: %s\n", ma.Command)}}
			return nil
		}

		step := 1
		if pstep, ok := ma.Params["incrementer.step"]; ok {
			if fstep, ok := pstep.(float64); ok {
				step = int(fstep)
			}
			if istep, ok := pstep.(int64); ok {
				step = int(istep)
			}
			if istep, ok := pstep.(int); ok {
				step = istep
			}
		}

		out, err2 := executeDrpCliCommand("machines", "get", ma.Uuid.String(), "param", parameter)
		if err2 != nil {
			*err = backend.Error{Code: 409,
				Model:    "plugin",
				Key:      "incrementer",
				Type:     "rpc",
				Messages: []string{fmt.Sprintf("Finding parameter failed: %s\n", err2.Error())}}
			return nil
		}

		if strings.TrimSpace(out) == "null" {
			_, err2 = executeDrpCliCommand("machines", "set", ma.Uuid.String(), "param", parameter, "to", fmt.Sprintf("%d", step))
			if err2 != nil {
				*err = backend.Error{Code: 409,
					Model:    "plugin",
					Key:      "incrementer",
					Type:     "rpc",
					Messages: []string{fmt.Sprintf("Failed to set an int to 0: %s\n", err2.Error())}}
				return nil
			}
		} else {
			i, err2 := strconv.ParseInt(strings.TrimSpace(out), 10, 64)
			if err2 != nil {
				*err = backend.Error{Code: 409,
					Model:    "plugin",
					Key:      "incrementer",
					Type:     "rpc",
					Messages: []string{fmt.Sprintf("Retrieved something not an int: %s\n", err2.Error())}}
				return nil
			}

			_, err2 = executeDrpCliCommand("machines", "set", ma.Uuid.String(), "param", parameter, "to", fmt.Sprintf("%d", i+int64(step)))
			if err2 != nil {
				*err = backend.Error{Code: 409,
					Model:    "plugin",
					Key:      "incrementer",
					Type:     "rpc",
					Messages: []string{fmt.Sprintf("Failed to set an int: %s\n", err2.Error())}}
				return nil
			}
		}
		*err = backend.Error{Code: 0,
			Model:    "plugin",
			Key:      "incrementer",
			Type:     "rpc",
			Messages: []string{}}
	} else {
		*err = backend.Error{Code: 404,
			Model:    "plugin",
			Key:      "incrementer",
			Type:     "rpc",
			Messages: []string{fmt.Sprintf("Unknown command: %s\n", ma.Command)}}
	}
	return nil
}

func main() {
	plugin.InitApp("incrementer", "Increments a parameter on a machine", version, &def)

	rpc.Register(&Plugin{})

	err := plugin.App.Execute()
	if err != nil {
		os.Exit(1)
	}
}
