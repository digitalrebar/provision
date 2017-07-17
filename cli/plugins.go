package cli

import (
	"fmt"

	"github.com/digitalrebar/provision/backend"
	"github.com/digitalrebar/provision/client/plugins"
	"github.com/digitalrebar/provision/models"
	"github.com/ghodss/yaml"
	"github.com/spf13/cobra"
)

type PluginOps struct{}

func (be PluginOps) GetType() interface{} {
	return &models.Plugin{}
}

func (be PluginOps) GetId(obj interface{}) (string, error) {
	plugin, ok := obj.(*models.Plugin)
	if !ok {
		return "", fmt.Errorf("Invalid type passed to plugin create")
	}
	return *plugin.Name, nil
}

func (be PluginOps) GetIndexes() map[string]string {
	b := &backend.Plugin{}
	ans := map[string]string{}
	for k, v := range b.Indexes() {
		ans[k] = v.Type
	}
	return ans
}

func (be PluginOps) List(parms map[string]string) (interface{}, error) {
	params := plugins.NewListPluginsParams()
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
		}
	}
	d, e := session.Plugins.ListPlugins(params, basicAuth)
	if e != nil {
		return nil, e
	}
	return d.Payload, nil
}

func (be PluginOps) Get(id string) (interface{}, error) {
	d, e := session.Plugins.GetPlugin(plugins.NewGetPluginParams().WithName(id), basicAuth)
	if e != nil {
		return nil, e
	}
	return d.Payload, nil
}

func (be PluginOps) Create(obj interface{}) (interface{}, error) {
	plugin, ok := obj.(*models.Plugin)
	if !ok {
		name, ok := obj.(string)
		if !ok {
			return nil, fmt.Errorf("Invalid type passed to plugin create")
		}
		plugin = &models.Plugin{Name: &name}
	}
	d, e := session.Plugins.CreatePlugin(plugins.NewCreatePluginParams().WithBody(plugin), basicAuth)
	if e != nil {
		return nil, e
	}
	return d.Payload, nil
}

func (be PluginOps) Patch(id string, obj interface{}) (interface{}, error) {
	data, ok := obj.(models.Patch)
	if !ok {
		return nil, fmt.Errorf("Invalid type passed to plugin patch")
	}
	d, e := session.Plugins.PatchPlugin(plugins.NewPatchPluginParams().WithName(id).WithBody(data), basicAuth)
	if e != nil {
		return nil, e
	}
	return d.Payload, nil
}

func (be PluginOps) Delete(id string) (interface{}, error) {
	d, e := session.Plugins.DeletePlugin(plugins.NewDeletePluginParams().WithName(id), basicAuth)
	if e != nil {
		return nil, e
	}
	return d.Payload, nil
}

func init() {
	tree := addPluginCommands()
	App.AddCommand(tree)
}

func addPluginCommands() (res *cobra.Command) {
	singularName := "plugin"
	name := "plugins"
	d("Making command tree for %v\n", name)
	res = &cobra.Command{
		Use:   name,
		Short: fmt.Sprintf("Access CLI commands relating to %v", name),
	}

	mo := &PluginOps{}

	commands := commonOps(singularName, name, mo)

	commands = append(commands, &cobra.Command{
		Use:   "params [id] [json]",
		Short: fmt.Sprintf("Gets/sets all parameters for the plugin"),
		Long:  `A helper function to return all or set all the parameters on the plugin`,
		RunE: func(c *cobra.Command, args []string) error {
			if len(args) != 1 && len(args) != 2 {
				return fmt.Errorf("%v requires 1 or 2 arguments", c.UseLine())
			}
			dumpUsage = false
			name := args[0]
			if len(args) == 1 {
				d, err := session.Plugins.GetPluginParams(plugins.NewGetPluginParamsParams().WithName(name), basicAuth)
				if err != nil {
					return generateError(err, "Failed to fetch params %v: %v", singularName, name)
				}
				return prettyPrint(d.Payload)
			} else {
				newValue := args[1]
				var value map[string]interface{}
				err := yaml.Unmarshal([]byte(newValue), &value)
				if err != nil {
					return fmt.Errorf("Unable to unmarshal input stream: %v\n", err)
				}
				d, err := session.Plugins.PostPluginParams(plugins.NewPostPluginParamsParams().WithName(name).WithBody(value), basicAuth)
				if err != nil {
					return generateError(err, "Failed to fetch params %v: %v", singularName, name)
				}
				return prettyPrint(d.Payload)
			}
		},
	})

	commands = append(commands, &cobra.Command{
		Use:   "get [id] param [key]",
		Short: fmt.Sprintf("Get a parameter from the plugin"),
		Long:  `A helper function to return the value of the parameter on the plugin`,
		RunE: func(c *cobra.Command, args []string) error {
			if len(args) != 3 {
				return fmt.Errorf("%v requires 3 arguments", c.UseLine())
			}
			dumpUsage = false
			name := args[0]
			// at = args[1]
			key := args[2]

			d, err := session.Plugins.GetPluginParams(plugins.NewGetPluginParamsParams().WithName(name), basicAuth)
			if err != nil {
				return generateError(err, "Failed to fetch params %v: %v", singularName, name)
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
		Short: fmt.Sprintf("Set the plugin's param <key> to <blob>"),
		Long:  `Helper function to update the plugin's parameters.`,
		RunE: func(c *cobra.Command, args []string) error {
			if len(args) != 5 {
				return fmt.Errorf("%v requires 5 arguments", c.UseLine())
			}
			name := args[0]
			key := args[2]
			newValue := args[4]
			dumpUsage = false

			var value interface{}
			err := yaml.Unmarshal([]byte(newValue), &value)
			if err != nil {
				return fmt.Errorf("Unable to unmarshal input stream: %v\n", err)
			}

			d, err := session.Plugins.GetPluginParams(plugins.NewGetPluginParamsParams().WithName(name), basicAuth)
			if err != nil {
				return generateError(err, "Failed to fetch params %v: %v", singularName, name)
			}
			pp := d.Payload
			if value == nil {
				delete(pp, key)
			} else {
				pp[key] = value
			}
			_, err = session.Plugins.PostPluginParams(plugins.NewPostPluginParamsParams().WithName(name).WithBody(pp), basicAuth)
			if err != nil {
				return generateError(err, "Failed to post params %v: %v", singularName, name)
			}
			return prettyPrint(value)
		},
	})

	res.AddCommand(commands...)
	return res
}
