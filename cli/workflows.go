package cli

import (
	"github.com/digitalrebar/provision/models"
	"github.com/spf13/cobra"
)

func init() {
	addRegistrar(registerWorkflow)
}

func registerWorkflow(app *cobra.Command) {
	op := &ops{
		name:       "workflows",
		singleName: "workflow",
		example:    func() models.Model { return &models.Workflow{} },
	}
	op.command(app)
}
