package cli

import (
	"github.com/digitalrebar/provision/models"
	"github.com/spf13/cobra"
)

func init() {
	addRegistrar(registerTask)
}

func registerTask(app *cobra.Command) {
	op := &ops{
		name:       "tasks",
		singleName: "task",
		example:    func() models.Model { return &models.Task{} },
	}
	op.command(app)
}
