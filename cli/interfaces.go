package cli

import (
	"github.com/digitalrebar/provision/models"
	"github.com/spf13/cobra"
)

func init() {
	addRegistrar(registerInterface)
}

func registerInterface(app *cobra.Command) {
	op := &ops{
		name:       "interfaces",
		singleName: "interface",
		example:    func() models.Model { return &models.Interface{} },
		noCreate:   true,
		noUpdate:   true,
		noDestroy:  true,
	}
	op.command(app)
}
