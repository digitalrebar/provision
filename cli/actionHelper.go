package cli

import (
	"fmt"

	"github.com/digitalrebar/provision/api"
	"github.com/digitalrebar/provision/models"
	"github.com/spf13/cobra"
)

func (o *ops) getPrefix() string {
	prefix := "system"
	if o.example != nil {
		prefix = o.example().Prefix()
	}
	if o.singleName == "extended" {
		prefix = o.name
	}
	return prefix
}

func (o *ops) actions() {
	actionName := o.actionName
	if actionName == "" {
		actionName = "action"
	}
	actionsName := fmt.Sprintf("%ss", actionName)
	idStr := ""
	argCount := 0
	evenCount := 1
	if o.example != nil {
		idStr = " [id]"
		argCount = 1
		evenCount = 0
	}
	plugin := ""
	actions := &cobra.Command{
		Use:   fmt.Sprintf("%s%s", actionsName, idStr),
		Short: fmt.Sprintf("Display actions for this %s", o.singleName),
		Long:  fmt.Sprintf("Helper function to display the %s's actions.", o.singleName),
		Args: func(c *cobra.Command, args []string) error {
			if len(args) != argCount {
				return fmt.Errorf("%v requires %d argument", c.UseLine(), argCount)
			}
			return nil
		},
		RunE: func(c *cobra.Command, args []string) error {
			prefix := o.getPrefix()
			res := []models.AvailableAction{}
			var req *api.R
			id := "system"
			if argCount == 1 {
				id = args[0]
				req = session.Req().UrlFor(prefix, id, actionsName)
			} else {
				req = session.Req().UrlFor(prefix, actionsName)
			}
			if plugin != "" {
				req = req.Params("plugin", plugin)
			}
			if err := req.Do(&res); err != nil {
				return generateError(err, "Failed to fetch actions %v: %v", o.singleName, id)
			}
			return prettyPrint(res)
		},
	}
	actions.Flags().StringVar(&plugin, "plugin", "", "Plugin to filter action search")
	o.addCommand(actions)
	action := &cobra.Command{
		Use:   fmt.Sprintf("%s%s [action]", actionName, idStr),
		Short: fmt.Sprintf("Display the action for this %s", o.singleName),
		Long:  fmt.Sprintf("Helper function to display the %s's action.", o.singleName),
		Args: func(c *cobra.Command, args []string) error {
			if len(args) != argCount+1 {
				return fmt.Errorf("%v requires %d arguments", c.UseLine(), argCount+1)
			}
			return nil
		},
		RunE: func(c *cobra.Command, args []string) error {
			prefix := o.getPrefix()
			action := args[argCount]
			res := &models.AvailableAction{}
			var req *api.R
			id := "system"
			if argCount == 1 {
				id = args[0]
				req = session.Req().UrlFor(prefix, id, actionsName, action)
			} else {
				req = session.Req().UrlFor(prefix, actionsName, action)
			}
			if plugin != "" {
				req = req.Params("plugin", plugin)
			}
			if err := req.Do(&res); err != nil {
				return generateError(err, "Failed to fetch action %v: %v", o.singleName, id)
			}
			return prettyPrint(res)
		},
	}
	action.Flags().StringVar(&plugin, "plugin", "", "Plugin to filter action search")
	o.addCommand(action)
	actionParams := map[string]interface{}{}
	runaction := &cobra.Command{
		Use:   fmt.Sprintf("run%s%s [command] [- | JSON or YAML Map of objects | pairs of string objects]", actionName, idStr),
		Short: "Run action on object from plugin",
		Args: func(c *cobra.Command, args []string) error {
			actionParams = map[string]interface{}{}
			if len(args) == argCount+2 {
				if err := into(args[argCount+1], &actionParams); err != nil {
					return err
				}
				return nil
			}
			if len(args) >= argCount+1 && len(args)%2 == evenCount {
				for i := argCount + 1; i < len(args); i += 2 {
					var obj interface{}
					if err := api.DecodeYaml([]byte(args[i+1]), &obj); err != nil {
						return fmt.Errorf("Invalid parameters: %s %v\n", args[i+1], err)
					}
					actionParams[args[i]] = obj
				}
				return nil
			}
			if argCount == 1 {
				return fmt.Errorf("runaction either takes three arguments or a multiple of two, not %d", len(args))
			}
			return fmt.Errorf("runaction either takes two arguments or one plus a multiple of two, not %d", len(args))
		},

		RunE: func(c *cobra.Command, args []string) error {
			prefix := o.getPrefix()
			command := args[argCount]
			var resp interface{}
			var req *api.R
			if argCount == 1 {
				id := args[0]
				req = session.Req().Post(actionParams).UrlFor(prefix, id, actionsName, command)
			} else {
				req = session.Req().Post(actionParams).UrlFor(prefix, actionsName, command)
			}
			if plugin != "" {
				req = req.Params("plugin", plugin)
			}
			if err := req.Do(&resp); err != nil {
				return generateError(err, "Error running action")
			}
			return prettyPrint(resp)
		},
	}
	runaction.Flags().StringVar(&plugin, "plugin", "", "Plugin to filter action search")
	o.addCommand(runaction)
}
