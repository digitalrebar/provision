package cli

import (
	"github.com/digitalrebar/provision/v4/models"
	"github.com/spf13/cobra"
)

func init() {
	addRegistrar(registerProfile)
}

func registerProfile(app *cobra.Command) {
	op := &ops{
		name:       "profiles",
		singleName: "profile",
		example:    func() models.Model { return &models.Profile{} },
	}
	op.command(app)
}
