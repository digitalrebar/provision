package cli

import "github.com/spf13/cobra"

func init() {
	addRegistrar(registerIso)
}

func registerIso(app *cobra.Command) {
	app.AddCommand(blobCommands("isos"))
}
