package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/digitalrebar/provision"
	"github.com/digitalrebar/provision/cli"
	"github.com/digitalrebar/provision/models"
	"github.com/digitalrebar/provision/plugin"
)

var (
	version = provision.RS_VERSION
	def     = models.PluginProvider{
		Name:       "incrementer",
		Version:    version,
		HasPublish: true,
		AvailableActions: []*models.AvailableAction{
			&models.AvailableAction{Command: "increment",
				OptionalParams: []string{"incrementer.step", "incrementer.parameter"},
			},
			&models.AvailableAction{Command: "reset_count",
				RequiredParams: []string{"incrementer.touched"},
			},
		},
		Parameters: []*models.Param{
			&models.Param{Name: "incrementer.parameter", Schema: map[string]interface{}{"type": "string"}},
			&models.Param{Name: "incrementer.step", Schema: map[string]interface{}{"type": "integer"}},
			&models.Param{Name: "incrementer.touched", Schema: map[string]interface{}{"type": "integer"}},
		},
	}
	lock sync.Mutex
)

type Plugin struct {
}

func (p *Plugin) Config(config map[string]interface{}) *models.Error {
	err := &models.Error{Code: 0, Model: "plugin", Key: "incrementer", Type: "plugin", Messages: []string{}}
	plugin.Log("Config: %v\n", config)
	return err
}

func executeDrpCliCommand(args ...string) (string, error) {
	r, w, _ := os.Pipe()
	lock.Lock()
	savedOut := os.Stdout
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
	os.Stdout = savedOut
	lock.Unlock()
	out := <-outC

	return out, err2
}

func updateOrCreateParameter(uuid, parameter string, step int) *models.Error {
	out, err2 := executeDrpCliCommand("machines", "get", uuid, "param", parameter)
	if err2 != nil {
		return &models.Error{Code: 409,
			Model:    "plugin",
			Key:      "incrementer",
			Type:     "plugin",
			Messages: []string{fmt.Sprintf("Finding parameter failed: %s\n", err2.Error())}}
	}

	if strings.TrimSpace(out) == "null" {
		_, err2 = executeDrpCliCommand("machines", "set", uuid, "param", parameter, "to", fmt.Sprintf("%d", step))
		if err2 != nil {
			return &models.Error{Code: 409,
				Model:    "plugin",
				Key:      "incrementer",
				Type:     "plugin",
				Messages: []string{fmt.Sprintf("Failed to set an int to 0: %s\n", err2.Error())}}
		}
	} else {
		i, err2 := strconv.ParseInt(strings.TrimSpace(out), 10, 64)
		if err2 != nil {
			return &models.Error{Code: 409,
				Model:    "plugin",
				Key:      "incrementer",
				Type:     "plugin",
				Messages: []string{fmt.Sprintf("Retrieved something not an int: %s\n", err2.Error())}}
		}

		_, err2 = executeDrpCliCommand("machines", "set", uuid, "param", parameter, "to", fmt.Sprintf("%d", i+int64(step)))
		if err2 != nil {
			return &models.Error{Code: 409,
				Model:    "plugin",
				Key:      "incrementer",
				Type:     "plugin",
				Messages: []string{fmt.Sprintf("Failed to set an int: %s\n", err2.Error())}}
		}
	}

	return &models.Error{Code: 0,
		Model:    "plugin",
		Key:      "incrementer",
		Type:     "plugin",
		Messages: []string{}}
}

func removeParameter(uuid, parameter string) *models.Error {
	_, err2 := executeDrpCliCommand("machines", "set", uuid, "param", parameter, "to", "null")
	if err2 != nil {
		return &models.Error{Code: 409,
			Model:    "plugin",
			Key:      "incrementer",
			Type:     "plugin",
			Messages: []string{fmt.Sprintf("Failed to remove param %s: %s\n", parameter, err2.Error())}}
	}

	return &models.Error{Code: 0,
		Model:    "plugin",
		Key:      "incrementer",
		Type:     "plugin",
		Messages: []string{}}
}

func (p *Plugin) Action(ma *models.MachineAction) *models.Error {
	plugin.Log("Action: %v\n", ma)
	if ma.Command == "increment" {
		parameter, ok := ma.Params["incrementer.parameter"].(string)
		if !ok {
			parameter = "incrementer.default"
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

		err := updateOrCreateParameter(ma.Uuid.String(), parameter, step)
		if err.Code == 0 {
			updateOrCreateParameter(ma.Uuid.String(), "incrementer.touched", 1)
		}
		return err
	} else if ma.Command == "reset_count" {
		return removeParameter(ma.Uuid.String(), "incrementer.touched")
	}

	return &models.Error{Code: 404,
		Model:    "plugin",
		Key:      "incrementer",
		Type:     "plugin",
		Messages: []string{fmt.Sprintf("Unknown command: %s\n", ma.Command)}}
}

func (p *Plugin) Publish(e *models.Event) *models.Error {
	return &models.Error{Code: 0,
		Model:    "plugin",
		Type:     "publish",
		Key:      "incrementer",
		Messages: []string{}}
}

func main() {
	plugin.InitApp("incrementer", "Increments a parameter on a machine", version, &def, &Plugin{})
	err := plugin.App.Execute()
	if err != nil {
		os.Exit(1)
	}
}
