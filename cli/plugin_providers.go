package cli

import (
	"fmt"

	"github.com/digitalrebar/provision/models"
	"github.com/spf13/cobra"
)

func init() {
	addRegistrar(registerPluginProvider)
}

func registerPluginProvider(app *cobra.Command) {
	op := &ops{
		name:       "plugin_providers",
		singleName: "plugin_provider",
		example:    func() models.Model { return &models.PluginProvider{} },
		noCreate:   true,
		noUpdate:   true,
	}
	op.addCommand(&cobra.Command{
		Use:   "upload [name] from [file]",
		Short: "Upload a program to act as a plugin_provider",
		Args: func(c *cobra.Command, args []string) error {
			if len(args) != 3 {
				return fmt.Errorf("%v requires 2 arguments", c.UseLine())
			}
			return nil
		},
		RunE: func(c *cobra.Command, args []string) error {
			name := args[0]
			filePath := args[2]
			fi, err := urlOrFileAsReadCloser(filePath)
			if err != nil {
				return fmt.Errorf("Error opening %s: %v", filePath, err)
			}
			defer fi.Close()
			res := &models.PluginProviderUploadInfo{}
			if err := session.Req().Post(fi).UrlFor(op.name, name).Do(res); err != nil {
				return err
			}
			return prettyPrint(res)
		},
	})
	op.command(app)
}
