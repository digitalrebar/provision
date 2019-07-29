package cli

import (
	"github.com/digitalrebar/provision/v4/models"
	"github.com/spf13/cobra"
)

func init() {
	addRegistrar(registerLease)
}

func registerLease(app *cobra.Command) {
	op := &ops{
		name:       "leases",
		singleName: "lease",
		example:    func() models.Model { return &models.Lease{} },
		noCreate:   true,
		noUpdate:   true,
	}
	op.command(app)
}
