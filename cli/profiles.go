package cli

import (
	"fmt"

	"github.com/digitalrebar/provision/client/profiles"
	"github.com/digitalrebar/provision/models"
	"github.com/ghodss/yaml"
	"github.com/spf13/cobra"
)

type ProfileOps struct{}

func (be ProfileOps) GetType() interface{} {
	return &models.Profile{}
}

func (be ProfileOps) GetId(obj interface{}) (string, error) {
	profile, ok := obj.(*models.Profile)
	if !ok || profile.Name == nil {
		return "", fmt.Errorf("Invalid type passed to profile create")
	}
	return *profile.Name, nil
}

func (be ProfileOps) List(parms map[string]string) (interface{}, error) {
	params := profiles.NewListProfilesParams()
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
	d, e := session.Profiles.ListProfiles(params, basicAuth)
	if e != nil {
		return nil, e
	}
	return d.Payload, nil
}

func (be ProfileOps) Get(id string) (interface{}, error) {
	d, e := session.Profiles.GetProfile(profiles.NewGetProfileParams().WithName(id), basicAuth)
	if e != nil {
		return nil, e
	}
	return d.Payload, nil
}

func (be ProfileOps) Create(obj interface{}) (interface{}, error) {
	profile, ok := obj.(*models.Profile)
	if !ok {
		return nil, fmt.Errorf("Invalid type passed to profile create")
	}
	d, e := session.Profiles.CreateProfile(profiles.NewCreateProfileParams().WithBody(profile), basicAuth)
	if e != nil {
		return nil, e
	}
	return d.Payload, nil
}

func (be ProfileOps) Patch(id string, obj interface{}) (interface{}, error) {
	data, ok := obj.(models.Patch)
	if !ok {
		return nil, fmt.Errorf("Invalid type passed to profile patch")
	}
	d, e := session.Profiles.PatchProfile(profiles.NewPatchProfileParams().WithName(id).WithBody(data), basicAuth)
	if e != nil {
		return nil, e
	}
	return d.Payload, nil
}

func (be ProfileOps) Delete(id string) (interface{}, error) {
	d, e := session.Profiles.DeleteProfile(profiles.NewDeleteProfileParams().WithName(id), basicAuth)
	if e != nil {
		return nil, e
	}
	return d.Payload, nil
}

func init() {
	tree := addProfileCommands()
	App.AddCommand(tree)
}

func addProfileCommands() (res *cobra.Command) {
	singularName := "profile"
	name := "profiles"
	d("Making command tree for %v\n", name)
	res = &cobra.Command{
		Use:   name,
		Short: fmt.Sprintf("Access CLI commands relating to %v", name),
	}

	mo := &ProfileOps{}

	commands := commonOps(singularName, name, mo)

	commands = append(commands, &cobra.Command{
		Use:   "params [id] [json]",
		Short: fmt.Sprintf("Gets/sets all parameters for the profile"),
		Long:  `A helper function to return all or set all the parameters on the profile`,
		RunE: func(c *cobra.Command, args []string) error {
			if len(args) != 1 && len(args) != 2 {
				return fmt.Errorf("%v requires 1 or 2 arguments", c.UseLine())
			}
			dumpUsage = false
			name := args[0]
			if len(args) == 1 {
				d, err := session.Profiles.GetProfileParams(profiles.NewGetProfileParamsParams().WithName(name), basicAuth)
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
				d, err := session.Profiles.PostProfileParams(profiles.NewPostProfileParamsParams().WithName(name).WithBody(value), basicAuth)
				if err != nil {
					return generateError(err, "Failed to fetch params %v: %v", singularName, name)
				}
				return prettyPrint(d.Payload)
			}
		},
	})

	commands = append(commands, &cobra.Command{
		Use:   "get [id] param [key]",
		Short: fmt.Sprintf("Get a parameter from the profile"),
		Long:  `A helper function to return the value of the parameter on the profile`,
		RunE: func(c *cobra.Command, args []string) error {
			if len(args) != 3 {
				return fmt.Errorf("%v requires 3 arguments", c.UseLine())
			}
			dumpUsage = false
			name := args[0]
			// at = args[1]
			key := args[2]

			d, err := session.Profiles.GetProfileParams(profiles.NewGetProfileParamsParams().WithName(name), basicAuth)
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
		Short: fmt.Sprintf("Set the profile's param <key> to <blob>"),
		Long:  `Helper function to update the profile's parameters.`,
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

			d, err := session.Profiles.GetProfileParams(profiles.NewGetProfileParamsParams().WithName(name), basicAuth)
			if err != nil {
				return generateError(err, "Failed to fetch params %v: %v", singularName, name)
			}
			pp := d.Payload
			if value == nil {
				delete(pp, key)
			} else {
				pp[key] = value
			}
			_, err = session.Profiles.PostProfileParams(profiles.NewPostProfileParamsParams().WithName(name).WithBody(pp), basicAuth)
			if err != nil {
				return generateError(err, "Failed to post params %v: %v", singularName, name)
			}
			return prettyPrint(value)
		},
	})

	res.AddCommand(commands...)
	return res
}
