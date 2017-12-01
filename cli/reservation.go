package cli

import (
	"github.com/digitalrebar/provision/models"
	"github.com/spf13/cobra"
)

func init() {
	addRegistrar(registerReservation)
}

func registerReservation(app *cobra.Command) {
	op := &ops{
		name:       "reservations",
		singleName: "reservation",
		example:    func() models.Model { return &models.Reservation{} },
	}
	op.command(app)
}
