package cli

import (
	"github.com/digitalrebar/provision/models"
	"github.com/spf13/cobra"
)

func init() {
	addRegistrar(registerStage)
}

func registerStage(app *cobra.Command) {
	op := &ops{
		name:       "stages",
		singleName: "stage",
		example:    func() models.Model { return &models.Stage{} },
	}
	op.command(app)
}
