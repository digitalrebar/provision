package cli

import (
	"github.com/digitalrebar/provision/models"
	"github.com/spf13/cobra"
)

func init() {
	addRegistrar(registerExtended)
}

func registerExtended(app *cobra.Command) {
	op := &ops{
		name:       "extended",
		singleName: "extended",
		example:    func() models.Model { return &models.RawModel{} },
	}
	op.command(app)
}
