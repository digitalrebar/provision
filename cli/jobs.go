package cli

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/digitalrebar/provision/v4/models"
	"github.com/spf13/cobra"
)

func init() {
	addRegistrar(registerJob)
}

func registerJob(app *cobra.Command) {
	op := &ops{
		name:       "jobs",
		singleName: "job",
		example:    func() models.Model { return &models.Job{} },
		noCreate:   true,
		actionName: "plugin_action",
	}
	op.addCommand(&cobra.Command{
		Use:   "create [json]",
		Short: fmt.Sprintf("Create a new %v with the passed-in JSON or string key", op.singleName),
		Long: `
As a useful shortcut, '-' can be passed to indicate that the JSON should
be read from stdin.

You may also pass in a machine UUID or Name to create a new job on that Name.
`,
		Args: func(c *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("%v requires 1 argument", c.UseLine())
			}
			return nil
		},
		RunE: func(c *cobra.Command, args []string) error {
			ref := &models.Job{}
			if err := into(args[0], ref); err != nil {
				if args[0] != "-" {
					m := &models.Machine{}
					if err := Session.FillModel(m, args[0]); err != nil {
						if err := Session.FillModel(m, "Name:"+args[0]); err != nil {
							return fmt.Errorf("Unable to create new Job: Invalid machine %s", args[0])
						}
					}
					ref.Machine = m.Uuid
				}
			}
			if err := Session.CreateModel(ref); err != nil {
				return generateError(err, "Unable to create new %v", op.singleName)
			}
			return prettyPrint(ref)
		},
	})
	purgeCmd := &cobra.Command{
		Use:   "purge",
		Short: "Purge jobs in excess of the job retention preferences",
		Args:  cobra.NoArgs,
	}
	purgeDryRun := purgeCmd.Flags().Bool("dry-run", false, "Report what would have been purged")
	purgeMachine := purgeCmd.Flags().String("machine", "",
		"Purge for a specific machine.  Use 'missing' to purge jobs for deleted machines")
	purgejobsToKeep := purgeCmd.Flags().String("jobs-to-keep", "",
		"Per-machine jobs to keep.  Overrides the jobsToKeep preference.")
	purgefailedAfter := purgeCmd.Flags().String("failed-after", "", "Purge failed jobs older than this duration.  "+
		"Overrides failedJobsPurgedAfter preference")
	purgeCompleteAfter := purgeCmd.Flags().String("complete-after", "", "Purge complete jobs older than this duration. "+
		" "+
		"Overrides completeJobsPurgedAfter preference")
	purgeCmd.RunE = func(c *cobra.Command, args []string) error {
		var res interface{}
		params := []string{}
		if *purgeDryRun {
			params = append(params, "dryrun", "true")
		}
		if *purgeMachine != "" {
			params = append(params, "machine", *purgeMachine)
		}
		if *purgejobsToKeep != "" {
			params = append(params, "jobsToKeep", *purgejobsToKeep)
		}
		if *purgefailedAfter != "" {
			params = append(params, "failedJobsPurgedAfter", *purgefailedAfter)
		}
		if *purgeCompleteAfter != "" {
			params = append(params, "completeJobsPurgedAfter", *purgeCompleteAfter)
		}
		info, err := Session.Info()
		if err != nil {
			return err
		}
		if len(params) > 0 && !info.HasFeature("selective-job-purge") {
			return fmt.Errorf("Server does not support purge flags.  Command ignored.")
		}
		if err = Session.Req().Meth("DELETE").UrlFor("jobs").Params(params...).Do(&res); err != nil {
			return err
		}
		return prettyPrint(res)
	}
	op.addCommand(purgeCmd)
	actionsFor := ""
	actionsCmd := &cobra.Command{
		Use:   "actions [id]",
		Short: "Get the actions for this job",
		Args: func(c *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("%v requires 1 argument", c.UseLine())
			}
			return nil
		},
		RunE: func(c *cobra.Command, args []string) error {
			uuid := args[0]
			res := models.JobActions{}
			if err := Session.Req().UrlFor("jobs", uuid, "actions").
				Params("os", actionsFor).Do(&res); err != nil {
				return generateError(err, "Error running action")
			}
			return prettyPrint(res)
		},
	}
	actionsCmd.Flags().StringVar(&actionsFor, "for-os", "", "OS to fetch actions for.  Defaults to fetching all actions")
	op.addCommand(actionsCmd)
	op.addCommand(&cobra.Command{
		Use:   "log [id] [- or string]",
		Short: "Gets the log or appends to the log if a second argument or stream is given",
		Args: func(c *cobra.Command, args []string) error {
			if len(args) < 1 {
				return fmt.Errorf("%v requires at least 1 argument", c.UseLine())
			}
			if len(args) > 2 {
				return fmt.Errorf("%v requires at most 2 arguments", c.UseLine())
			}
			return nil
		},
		RunE: func(c *cobra.Command, args []string) error {
			uuid := args[0]
			if len(args) == 1 {
				if err := Session.Req().UrlFor("jobs", uuid, "log").Do(os.Stdout); err != nil {
					return generateError(err, "Error getting log")
				}
				return nil
			}
			var src io.Reader
			if args[1] == "-" {
				src = os.Stdin
			} else {
				src = bytes.NewBufferString(args[1])
			}
			if err := Session.Req().Put(src).UrlFor("jobs", uuid, "log").Do(nil); err != nil {
				return generateError(err, "Error appending to log")
			}
			return nil
		},
	})
	op.command(app)
}
