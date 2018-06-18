package cli

import (
	"fmt"

	"github.com/digitalrebar/provision/models"
	"github.com/spf13/cobra"
)

func (o *ops) bootenv() {
	o.addCommand(&cobra.Command{
		Use:   "bootenv [id] [bootenv]",
		Short: fmt.Sprintf("Set the %s's bootenv", o.singleName),
		Long:  fmt.Sprintf(`Helper function to update the %s's bootenv.`, o.singleName),
		Args: func(c *cobra.Command, args []string) error {
			if len(args) != 2 {
				return fmt.Errorf("%v requires 2 arguments", c.UseLine())
			}
			return nil
		},
		RunE: func(c *cobra.Command, args []string) error {
			data, err := o.refOrFill(args[0])
			if err != nil {
				return err
			}
			ex := models.Clone(data).(models.BootEnver)
			ex.SetBootEnv(args[1])
			res, err := session.PatchToFull(data, ex, ref != "")
			if err != nil {
				return generateError(err, "Unable to update %s: %v", o.singleName, args[0])
			}
			return prettyPrint(res)
		},
	})
}
