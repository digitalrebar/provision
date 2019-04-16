package main

//go:generate drbundler content content.go

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/VictorLowther/jsonpatch2/utils"
	"github.com/digitalrebar/logger"
	"github.com/digitalrebar/provision"
	"github.com/digitalrebar/provision/api"
	"github.com/digitalrebar/provision/models"
	"github.com/digitalrebar/provision/plugin"
)

var (
	version = provision.RSVersion
	def     = models.PluginProvider{
		Name:          "incrementer",
		Version:       version,
		PluginVersion: 2,
		HasPublish:    true,
		AvailableActions: []models.AvailableAction{
			{Command: "increment",
				Model:          "machines",
				OptionalParams: []string{"incrementer/step", "incrementer/parameter"},
			},
			{Command: "reset_count",
				Model:          "machines",
				RequiredParams: []string{"incrementer/touched"},
			},
			{Command: "incrstatus"},
		},
		StoreObjects: map[string]interface{}{
			"cows": map[string]interface{}{},
			"typed-cows": map[string]interface{}{
				"Spotted": map[string]interface{}{
					"type":       "boolean",
					"isrequired": true,
				},
				"CanMilk": map[string]interface{}{
					"type": "boolean",
				},
				"Location": map[string]interface{}{
					"type": "string",
				},
			},
		},
		Content: contentYamlString,
	}
)

// Plugin is the base structure for the plugin.
type Plugin struct {
	session *api.Client
}

// Config is the plugin's configuration entrypoint
func (p *Plugin) Config(thelog logger.Logger, session *api.Client, config map[string]interface{}) *models.Error {
	thelog.Infof("Config: %v\n", config)
	p.session = session
	return nil
}

func (p *Plugin) updateOrCreateParameter(uuid, parameter string, step int) (interface{}, *models.Error) {
	e := &models.Error{Code: 400,
		Model:    "plugin",
		Key:      "incrementer",
		Type:     "plugin",
		Messages: []string{}}
	var res interface{}
	if err := p.session.Req().UrlFor("machines", uuid, "params", parameter).Do(&res); err != nil {
		e.AddError(err)
		return nil, e
	}
	i := int64(step)
	if res != nil {
		i += int64(res.(float64))
	}
	var params interface{}
	if err := p.session.Req().Post(i).UrlFor("machines", uuid, "params", parameter).Do(&params); err != nil {
		e.AddError(err)
		return nil, e
	}
	return i, nil
}

func (p *Plugin) removeParameter(uuid, parameter string) *models.Error {
	var param interface{}
	if err := p.session.Req().Del().UrlFor("machines", uuid, "params", parameter).Do(&param); err != nil {
		e := &models.Error{Code: 400,
			Model:    "plugin",
			Key:      "incrementer",
			Type:     "plugin",
			Messages: []string{}}
		e.AddError(err)
		return e
	}
	return nil
}

// Action is the plugin's action entrypoint
func (p *Plugin) Action(thelog logger.Logger, ma *models.Action) (interface{}, *models.Error) {
	thelog.Infof("Action: %v\n", ma)
	var machine models.Machine
	switch ma.Command {
	case "increment":
		if err := utils.Remarshal(ma.Model, &machine); err != nil {
			return nil, &models.Error{Code: 409,
				Model:    "plugin",
				Key:      "incrementer",
				Type:     "plugin",
				Messages: []string{fmt.Sprintf("%v", err)}}
		}
		parameter, ok := ma.Params["incrementer/parameter"].(string)
		if !ok {
			parameter = "incrementer/default"
		}
		step := 1
		if pstep, ok := ma.Params["incrementer/step"]; ok {
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
		val, err := p.updateOrCreateParameter(machine.UUID(), parameter, step)
		if err == nil {
			_, err = p.updateOrCreateParameter(machine.UUID(), "incrementer/touched", 1)
		}
		return val, err
	case "reset_count":
		var machine models.Machine
		if err := utils.Remarshal(ma.Model, &machine); err != nil {
			return nil, &models.Error{Code: 409,
				Model:    "plugin",
				Key:      "incrementer",
				Type:     "plugin",
				Messages: []string{fmt.Sprintf("%v", err)}}
		}
		e := p.removeParameter(machine.UUID(), "incrementer/touched")
		return "Success", e
	case "incrstatus":
		return "Running", nil
	}

	return nil, &models.Error{Code: 404,
		Model:    "plugin",
		Key:      "incrementer",
		Type:     "plugin",
		Messages: []string{fmt.Sprintf("Unknown command: %s\n", ma.Command)}}
}

// Publish is the plugin's event publisher entrypoint
func (p *Plugin) Publish(thelog logger.Logger, e *models.Event) *models.Error {
	thelog.Debugf("Plugin received: %v\n", e)
	// Just eat the publish messages.
	return nil
}

// Unpack is the plugin's unpack entrypoint
func (p *Plugin) Unpack(thelog logger.Logger, dir string) error {
	return ioutil.WriteFile(path.Join(dir, "testFile"), []byte("ImaFile"), 0644)
}

func main() {
	plugin.InitApp("incrementer", "Increments a parameter on a machine", version, &def, &Plugin{})
	err := plugin.App.Execute()
	if err != nil {
		os.Exit(1)
	}
}
