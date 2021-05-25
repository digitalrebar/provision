package cli

import (
	"github.com/digitalrebar/provision/v4/models"
	"github.com/spf13/cobra"
)

func init() {
	addRegistrar(registerAsyncActionCron)
}

func registerAsyncActionCron(app *cobra.Command) {
	op := &ops{
		name:       "async_action_crons",
		singleName: "async_action_cron",
		example:    func() models.Model { return &models.AsyncActionCron{} },
	}
	op.command(app)
}
