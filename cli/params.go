package cli

import (
	"github.com/digitalrebar/provision/models"
	"github.com/spf13/cobra"
)

func init() {
	addRegistrar(registerParam)
}

func registerParam(app *cobra.Command) {
	op := &ops{
		name:       "params",
		singleName: "param",
		example:    func() models.Model { return &models.Param{} },
	}
	op.command(app)
}
