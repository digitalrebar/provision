package cli

import (
	"github.com/digitalrebar/provision/v4/models"
	"github.com/spf13/cobra"
)

func init() {
	addRegistrar(registerAsyncActionTemplate)
}

func registerAsyncActionTemplate(app *cobra.Command) {
	op := &ops{
		name:       "async_action_templates",
		singleName: "async_action_template",
		example:    func() models.Model { return &models.AsyncActionTemplate{} },
	}
	op.command(app)
}
