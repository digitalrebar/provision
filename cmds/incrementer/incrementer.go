package main

import (
	"fmt"
	"net/rpc"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/digitalrebar/provision"
	"github.com/digitalrebar/provision/backend"
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
			},
		},
		Parameters: []*backend.Param{
			&backend.Param{Name: "incrementer.parameter", Schema: map[string]interface{}{"type": "string"}},
		},
	}
)

type Plugin struct {
}

func (p *Plugin) Config(config map[string]interface{}, err *backend.Error) error {
	*err = backend.Error{Code: 0, Model: "plugin", Key: "incrementer", Type: "rpc", Messages: []string{}}
	return nil
}

func (p *Plugin) Action(ma midlayer.MachineAction, err *backend.Error) error {
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

		out, err2 := exec.Command("./drpcli",
			"machines",
			"get",
			ma.Uuid.String(),
			"param", parameter).Output()
		if err2 != nil {
			*err = backend.Error{Code: 409,
				Model:    "plugin",
				Key:      "incrementer",
				Type:     "rpc",
				Messages: []string{fmt.Sprintf("Finding parameter failed: %s\n", err2.Error())}}
			return nil
		}

		if strings.TrimSpace(string(out)) == "null" {
			_, err2 = exec.Command("./drpcli",
				"machines",
				"set",
				ma.Uuid.String(),
				"param",
				parameter,
				"to",
				"0").Output()
			if err2 != nil {
				*err = backend.Error{Code: 409,
					Model:    "plugin",
					Key:      "incrementer",
					Type:     "rpc",
					Messages: []string{fmt.Sprintf("Failed to set an int to 0: %s\n", err2.Error())}}
				return nil
			}
		} else {
			i, err2 := strconv.ParseInt(strings.TrimSpace(string(out)), 10, 64)
			if err2 != nil {
				*err = backend.Error{Code: 409,
					Model:    "plugin",
					Key:      "incrementer",
					Type:     "rpc",
					Messages: []string{fmt.Sprintf("Retrieved something not an int: %s\n", err2.Error())}}
				return nil
			}

			out, err2 = exec.Command("./drpcli",
				"machines",
				"set",
				ma.Uuid.String(),
				"param", parameter,
				"to",
				fmt.Sprintf("%d", i+1)).Output()
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

func (p *Plugin) Stop(dummmy int, err *backend.Error) error {
	*err = backend.Error{Code: 0, Model: "plugin", Key: "incrementer", Type: "rpc", Messages: []string{}}
	const delay = 5000 * time.Millisecond
	go func() {
		time.Sleep(delay)
		os.Exit(0)
	}()

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
