package cli

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"time"

	"github.com/digitalrebar/provision/agent"
	"github.com/digitalrebar/provision/models"
	"github.com/spf13/cobra"
)

var (
	actuallyPowerThings = true
	defaultStateLoc     string
)

func init() {
	if defaultStateLoc == "" {
		switch runtime.GOOS {
		case "windows":
			defaultStateLoc = os.ExpandEnv("${APPDATA}/drp-agent")
		default:
			defaultStateLoc = "/var/lib/drp-agent"
		}
	}
	addRegistrar(registerMachine)
}

func registerMachine(app *cobra.Command) {
	op := &ops{
		name:       "machines",
		singleName: "machine",
		example:    func() models.Model { return &models.Machine{} },
	}
	op.addCommand(&cobra.Command{
		Use:   "workflow [id] [workflow]",
		Short: fmt.Sprintf("Set the machine's workflow"),
		Long:  `Helper function to update the machine's workflow.`,
		Args: func(c *cobra.Command, args []string) error {
			if len(args) != 2 {
				return fmt.Errorf("%v requires 2 arguments", c.UseLine())
			}
			return nil
		},
		RunE: func(c *cobra.Command, args []string) error {
			m, err := op.refOrFill(args[0])
			if err != nil {
				return generateError(err, "Failed to fetch %v: %v", op.singleName, args[0])
			}
			clone := models.Clone(m).(*models.Machine)
			clone.Workflow = args[1]
			req := session.Req().ParanoidPatch().PatchTo(m, clone)
			if force {
				req.Params("force", "true")
			}
			if err := req.Do(&clone); err != nil {
				return err
			}
			return prettyPrint(clone)
		},
	})
	op.addCommand(&cobra.Command{
		Use:   "stage [id] [stage]",
		Short: fmt.Sprintf("Set the machine's stage"),
		Long:  `Helper function to update the machine's stage.`,
		Args: func(c *cobra.Command, args []string) error {
			if len(args) != 2 {
				return fmt.Errorf("%v requires 2 arguments", c.UseLine())
			}
			return nil
		},
		RunE: func(c *cobra.Command, args []string) error {
			m, err := op.refOrFill(args[0])
			if err != nil {
				return generateError(err, "Failed to fetch %v: %v", op.singleName, args[0])
			}
			clone := models.Clone(m).(*models.Machine)
			clone.Stage = args[1]
			req := session.Req().ParanoidPatch().PatchTo(m, clone)
			if force {
				req.Params("force", "true")
			}
			if err := req.Do(&clone); err != nil {
				return err
			}
			return prettyPrint(clone)
		},
	})
	jobs := &cobra.Command{
		Use:   "jobs",
		Short: "Access commands for manipulating the current job",
	}
	jobs.AddCommand(&cobra.Command{
		Use:   "create [id]",
		Short: "Create a job for the current task on machine [id]",
		Args: func(c *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("%v requires 1 argument", c.UseLine())
			}
			return nil
		},
		RunE: func(c *cobra.Command, args []string) error {
			id := args[0]
			m, err := op.refOrFill(id)
			if err != nil {
				return generateError(err, "Failed to fetch %v: %v", op.singleName, id)
			}
			machine := m.(*models.Machine)
			job := &models.Job{}
			j2 := &models.Job{}
			job.Machine = machine.Uuid
			if err := session.Req().Post(job).UrlFor("jobs").Do(j2); err != nil {
				return generateError(err, "Failed to create job for %v: %v", op.singleName, id)
			}
			return prettyPrint(j2)
		},
	})
	jobs.AddCommand(&cobra.Command{
		Use:   "current [id]",
		Short: "Get the current job on machine [id]",
		Args: func(c *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("%v requires 1 argument", c.UseLine())
			}
			return nil
		},
		RunE: func(c *cobra.Command, args []string) error {
			id := args[0]
			m, err := op.refOrFill(id)
			if err != nil {
				return generateError(err, "Failed to fetch %v: %v", op.singleName, id)
			}
			machine := m.(*models.Machine)
			if machine.CurrentJob == nil || len(machine.CurrentJob) == 0 {
				return fmt.Errorf("No current job on machine %v", m.Key())
			}
			job := &models.Job{}
			if err := session.Req().UrlFor("jobs", machine.CurrentJob.String()).Do(job); err != nil {
				return generateError(err, "Failed to fetch current job for %v: %v", op.singleName, id)
			}
			return prettyPrint(job)
		},
	})
	jobs.AddCommand(&cobra.Command{
		Use:   "state [id] to [state]",
		Short: "Set the current job on machine [id] to [state]",
		Args: func(c *cobra.Command, args []string) error {
			if len(args) != 3 {
				return fmt.Errorf("%v requires 2 arguments", c.UseLine())
			}
			return nil
		},
		RunE: func(c *cobra.Command, args []string) error {
			id := args[0]
			state := args[2]
			m, err := op.refOrFill(id)
			if err != nil {
				return generateError(err, "Failed to fetch %v: %v", op.singleName, id)
			}
			machine := m.(*models.Machine)
			if machine.CurrentJob == nil || len(machine.CurrentJob) == 0 {
				return fmt.Errorf("No current job on machine %v", m.Key())
			}
			job := &models.Job{}
			if err := session.Req().UrlFor("jobs", machine.CurrentJob.String()).Do(job); err != nil {
				return generateError(err, "Failed to fetch current job for %v: %v", op.singleName, id)
			}
			j2 := models.Clone(job).(*models.Job)
			j2.State = state
			j3, err := session.PatchTo(job, j2)
			if err != nil {
				return generateError(err, "Failed to mark job %s as %s", job.Uuid, state)
			}
			return prettyPrint(j3)
		},
	})
	op.addCommand(jobs)
	op.addCommand(&cobra.Command{
		Use:   "currentlog [id]",
		Short: "Get the log for the most recent job run on the machine",
		Args: func(c *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("%v requires 1 argument", c.UseLine())
			}
			return nil
		},
		RunE: func(c *cobra.Command, args []string) error {
			m, err := op.refOrFill(args[0])
			if err != nil {
				return generateError(err, "Failed to fetch %v: %v", op.singleName, args[0])
			}
			return session.Req().UrlFor("jobs", m.(*models.Machine).CurrentJob.String(), "log").Do(os.Stdout)
		},
	})
	op.addCommand(&cobra.Command{
		Use:   "deletejobs [id]",
		Short: "Delete all jobs associated with machine",
		Args: func(c *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("%v requires 1 argument", c.UseLine())
			}
			return nil
		},
		RunE: func(c *cobra.Command, args []string) error {
			m, err := op.refOrFill(args[0])
			if err != nil {
				return generateError(err, "Failed to fetch %v: %v", op.singleName, args[0])
			}
			jobs := []*models.Job{}
			if err := session.Req().Filter("jobs",
				"Machine", "Eq", m.Key(),
				"sort", "StartTime",
				"reverse").Do(&jobs); err != nil {
				return generateError(err, "Failed to fetch jobs for %s: %v", op.singleName, args[0])
			}
			for _, job := range jobs {
				if _, err := session.DeleteModel("jobs", job.Key()); err != nil {
					return generateError(err, "Failed to delete Job %s", job.Key())
				} else {
					fmt.Printf("Deleted Job %s", job.Key())
				}
			}
			return nil
		},
	})
	tasks := &cobra.Command{
		Use:   "tasks",
		Short: "Access task manipulation for machines",
	}
	tasks.AddCommand(&cobra.Command{
		Use:   "add [id] [at [offset]] [task...]",
		Short: "Add tasks to the task list for [id]",
		Long: `You may omit "at [offset]" to indicate that the task should be appended to the
end of the task list.  Otherwise, [offset] 0 indicates that the tasks
should be inserted immediately after the current task. Negative numbers
are not accepted.`,
		Args: func(c *cobra.Command, args []string) error {
			if len(args) < 2 {
				return fmt.Errorf("%v requires at least an id and one task", c.UseLine())
			}
			if args[1] == "at" {
				if len(args) < 4 {
					return fmt.Errorf("%v requires at least 3 arguments when specifying an offset", c.UseLine())
				}
				if _, err := strconv.Atoi(args[2]); err != nil {
					return fmt.Errorf("%v: %s is not a number", c.UseLine(), args[2])
				}
			}
			return nil
		},
		RunE: func(c *cobra.Command, args []string) error {
			id := args[0]
			var offset = -1
			var tasks []string
			if args[1] == "at" {
				offset, _ = strconv.Atoi(args[2])
				tasks = args[3:]
			} else {
				tasks = args[1:]
			}
			obj, err := op.refOrFill(id)
			if err != nil {
				return err
			}
			m := models.Clone(obj).(*models.Machine)
			if err := m.AddTasks(offset, tasks...); err != nil {
				generateError(err, "Cannot add tasks")
			}
			if err := session.Req().PatchTo(obj, m).Do(&m); err != nil {
				return err
			}
			return prettyPrint(m)
		},
	})
	tasks.AddCommand(&cobra.Command{
		Use:   "del [task...]",
		Short: "Remove tasks from the mutable part of the task list",
		Long: `Each entry in [task...] will remove at most one instance of it from the
machine task list.  If you want to remove more than one, you need to
pass in more than one task.`,
		Args: func(c *cobra.Command, args []string) error {
			if len(args) < 2 {
				return fmt.Errorf("%v requires at least an id and one task", c.UseLine())
			}
			return nil
		},
		RunE: func(c *cobra.Command, args []string) error {
			id := args[0]
			tasks := args[1:]
			obj, err := op.refOrFill(id)
			if err != nil {
				return err
			}
			m := models.Clone(obj).(*models.Machine)
			m.DelTasks(tasks...)
			if err := session.Req().PatchTo(obj, m).Do(&m); err != nil {
				return err
			}
			return prettyPrint(m)
		},
	})
	op.addCommand(tasks)
	var exitOnFailure = false
	var oneShot = false
	var runStateLoc string
	processJobs := &cobra.Command{
		Use:   "processjobs [id]",
		Short: "For the given machine, process pending jobs until done.",
		Long: `
For the provided machine, identified by UUID, process the task list on
that machine until an error occurs or all jobs are complete.  Upon
completion, optionally wait for additional jobs as specified by
the stage runner wait flag.
`,
		Args: func(c *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("%v requires 1 argument", c.UseLine())
			}
			return nil
		},
		RunE: func(c *cobra.Command, args []string) error {
			uuid := args[0]
			m := &models.Machine{}
			if err := session.FillModel(m, uuid); err != nil {
				return err
			}
			if runStateLoc == "" {
				runStateLoc = defaultStateLoc
			}
			if runStateLoc != "" {
				if err := os.MkdirAll(runStateLoc, 0700); err != nil {
					return fmt.Errorf("Unable to create state directory %s: %v", runStateLoc, err)
				}
			}
			agent, err := agent.New(session, m, oneShot, exitOnFailure, actuallyPowerThings, os.Stdout)
			if err != nil {
				return err
			}
			if oneShot {
				agent = agent.Timeout(time.Second)
			}
			return agent.StateLoc(runStateLoc).Run()
		},
	}
	processJobs.Flags().BoolVar(&exitOnFailure, "exit-on-failure", false, "Exit on failure of a task")
	processJobs.Flags().BoolVar(&oneShot, "oneshot", false, "Do not wait for additional tasks to appear")
	processJobs.Flags().StringVar(&runStateLoc, "stateDir", "", "Location to save agent runtime state")
	op.addCommand(processJobs)
	op.command(app)
}
