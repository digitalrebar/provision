package cli

import (
	"github.com/digitalrebar/provision/v4/models"
	"github.com/spf13/cobra"
)

func init() {
	addRegistrar(registerVersionSet)
}

func registerVersionSet(app *cobra.Command) {
	op := &ops{
		name:       "version_sets",
		singleName: "version_set",
		example:    func() models.Model { return &models.VersionSet{} },
	}
	op.command(app)
}
