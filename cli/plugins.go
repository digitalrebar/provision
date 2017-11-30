package cli

import (
	"github.com/digitalrebar/provision/models"
	"github.com/spf13/cobra"
)

func init() {
	addRegistrar(registerPlugin)
}

func registerPlugin(app *cobra.Command) {
	op := &ops{
		name:       "plugins",
		singleName: "plugin",
		example:    func() models.Model { return &models.Plugin{} },
	}
	op.command(app)
}
