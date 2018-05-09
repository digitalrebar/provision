package cli

import (
	"github.com/digitalrebar/provision/models"
	"github.com/spf13/cobra"
)

func init() {
	addRegistrar(registerTenant)
}

func registerTenant(app *cobra.Command) {
	op := &ops{
		name:       "tenants",
		singleName: "tenant",
		example:    func() models.Model { return &models.Tenant{} },
	}
	op.command(app)
}
