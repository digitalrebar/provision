package cli

import (
	"github.com/digitalrebar/provision/v4/models"
	"github.com/spf13/cobra"
)

func init() {
	addRegistrar(registerContext)
}

func registerContext(app *cobra.Command) {
	op := &ops{
		name:       "contexts",
		singleName: "context",
		example:    func() models.Model { return &models.Context{} },
	}
	op.command(app)
}
