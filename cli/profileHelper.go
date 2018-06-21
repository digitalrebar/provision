package cli

import (
	"fmt"

	"github.com/digitalrebar/provision/models"
	"github.com/spf13/cobra"
)

func (o *ops) profiles() {
	o.addCommand(&cobra.Command{
		Use:   "addprofile [id] [profile]",
		Short: fmt.Sprintf("Add profile to the machine's profile list"),
		Long:  `Helper function to add a profile to the machine's profile list.`,
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
			ex := models.Clone(data).(models.Profiler)
			ex.SetProfiles(append(ex.GetProfiles(), args[1]))
			res, err := session.PatchToFull(data, ex, ref != "")
			if err != nil {
				return generateError(err, "Unable to update %s: %v", o.singleName, args[0])
			}
			return prettyPrint(res)
		},
	})
	o.addCommand(&cobra.Command{
		Use:   "removeprofile [id] [profile]",
		Short: fmt.Sprintf("Remove a profile from the machine's list"),
		Long:  `Helper function to update the machine's profile list by removing one.`,
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
			ex := models.Clone(data).(models.Profiler)
			newProfiles := []string{}
			for _, s := range ex.GetProfiles() {
				if s == args[1] {
					continue
				}
				newProfiles = append(newProfiles, s)
			}
			ex.SetProfiles(newProfiles)
			res, err := session.PatchToFull(data, ex, ref != "")
			if err != nil {
				return generateError(err, "Unable to update %s: %v", o.singleName, args[0])
			}
			return prettyPrint(res)
		},
	})

}
