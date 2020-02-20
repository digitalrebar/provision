package cli

import (
	"github.com/digitalrebar/provision/v4/models"
	"github.com/spf13/cobra"
)

func init() {
	addRegistrar(registerEndpoint)
}

func registerEndpoint(app *cobra.Command) {
	op := &ops{
		name:       "endpoints",
		singleName: "endpoint",
		example:    func() models.Model { return &models.Endpoint{} },
	}
	op.command(app)
}
