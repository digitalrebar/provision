package cli

import (
	"github.com/digitalrebar/provision/models"
	"github.com/spf13/cobra"
)

func init() {
	addRegistrar(registerRole)
}

func registerRole(app *cobra.Command) {
	op := &ops{
		name:       "roles",
		singleName: "role",
		example:    func() models.Model { return &models.Role{} },
	}
	op.command(app)
}
