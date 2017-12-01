package cli

import (
	"fmt"

	"github.com/digitalrebar/provision/models"
	"github.com/spf13/cobra"
)

func init() {
	addRegistrar(registerTemplate)
}

func registerTemplate(app *cobra.Command) {
	op := &ops{
		name:       "templates",
		singleName: "template",
		example:    func() models.Model { return &models.Template{} },
	}
	op.addCommand(&cobra.Command{
		Use:   "upload [file] as [id]",
		Short: "Upload the template file [file] as template [id]",
		Args: func(c *cobra.Command, args []string) error {
			if len(args) != 3 {
				return fmt.Errorf("%v: expected 3 argument", c.UseLine())
			}
			return nil
		},
		RunE: func(c *cobra.Command, args []string) error {
			tmpl, err := session.InstallRawTemplateFromFileWithId(args[0], args[2])
			if err != nil {
				return err
			}
			return prettyPrint(tmpl)
		},
	})
	op.command(app)
}
