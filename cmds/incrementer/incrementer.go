package main

//go:generate sh -c "cd content ; drpcli contents bundle ../content.go"

// Using go generate to package a content bundle into a file we can pull in
// is a fairly common pattern we use while building plugin providers.  It
// is an easy method to use to ensure that a plugin provider and the content
// designed to integrate with it stay in sync.

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/VictorLowther/jsonpatch2/utils"
	"github.com/digitalrebar/logger"
	v4 "github.com/digitalrebar/provision/v4"
	"github.com/digitalrebar/provision/v4/api"
	"github.com/digitalrebar/provision/v4/models"
	"github.com/digitalrebar/provision/v4/plugin"
)

var (
	// Every plugin provider needs to define its version.
	// A global variable is as good a place as any.

	version = v4.RSVersion

	// Every plugin provider also needs to be able to define itself
	// with dr-provision.  models.PluginProvider
	// is the struct that provides that definition.
	def = models.PluginProvider{
		Name:          "incrementer",
		Version:       version,
		PluginVersion: 4,
		AutoStart:     false,
		HasPublish:    false,
		AvailableActions: []models.AvailableAction{
			{
				Command:        "increment",
				Model:          "machines",
				OptionalParams: []string{"incrementer/step", "incrementer/parameter"},
			},
			{
				Command:        "reset_count",
				Model:          "machines",
				RequiredParams: []string{"incrementer/touched"},
			},
			{
				Command: "incrstatus",
			},
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
// By convention, it is named Plugin, although it can be anything.
// It should hold all the plugin-specific information needed for the plugin
// to do its job.  In this case, that is just a reference to the API client.
type Plugin struct {
	session *api.Client
}

// Config is the plugin's configuration entrypoint.  It is responsible
// for handling any and all configuration changes over the lifecycle
// of a running plugin, including initialization.
// You can rely on Config being the first method called.
//
// For incrementer, the only thing it has to do is to save a reference to the
// api client that gets passed in for later.
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

// Action is the plugin's action entrypoint.  It is responsible for handling
// any actions that the provider declares it can handle via the AvailableActions field
// in the PluginProvider definition.
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

// Unpack is the plugin's unpack entrypoint.  It is responsible for
// unpacking any extra content the plugin provider may require to do its
// job.
func (p *Plugin) Unpack(thelog logger.Logger, dir string) error {
	return ioutil.WriteFile(path.Join(dir, "testFile"), []byte("ImaFile"), 0644)
}

// The main function of a plugin provider should do the bare minimum to run plugin.InitApp()
// and then run plugin.App.Execute().  The plugin package provides App as a global variable
// and will arrange for all the necessary command line and protocol definition interfaces.
func main() {
	plugin.InitApp("incrementer", "Increments a parameter on a machine", version, &def, &Plugin{})
	err := plugin.App.Execute()
	if err != nil {
		os.Exit(1)
	}
}
