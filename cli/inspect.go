package cli

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/digitalrebar/provision/v4/models"

	"github.com/spf13/cobra"
)

// drpcli machines inspect tasks <MachineID> <taskname> <index back optional>
// drpcli machines inspect jobs <MachineID> <index back optional>

func lookupMachineID(id string) (string, error) {
	m := &models.Machine{}
	if err := Session.FillModel(m, id); err != nil {
		return "", err
	}
	return m.UUID(), nil
}

func printJobs(data []*models.Job, full bool) error {
	if format != "text" {
		return prettyPrint(data)
	}

	formatString := "%3s: %40s %15s %15s %15s %15s\n"
	fmt.Printf(formatString, "Idx", "UUID", "State", "Workflow", "Stage", "Task")
	for i, j := range data {
		fmt.Printf(formatString, fmt.Sprintf("%d", i), j.Uuid.String(), j.State, j.Workflow, j.Stage, j.Task)
		if full {
			fmt.Printf(">>>>>>Log<<<<<<:\n%s\n", j.Meta["log"])

			if aa, ok := j.Meta["actions"]; ok {
				fmt.Printf(">>>>>>Actions<<<<<<:\n%s\n", aa)
			}
			if aa, ok := j.Meta["actionCount"]; ok {
				ai, _ := strconv.Atoi(aa)
				for i := 0; i < ai; i++ {
					fmt.Printf(">>>>>>Action %d<<<<<<:\n%s\n", i, j.Meta[fmt.Sprintf("action.%d", i)])
				}
			}
		}
	}
	return nil
}

func inspectCommands() *cobra.Command {
	full := false

	inspectCmds := &cobra.Command{
		Use:   "inspect",
		Short: "Commands to inspect tasks and jobs on machines",
	}
	inspectCmds.PersistentFlags().BoolVarP(&full, "full", "A", false, "Should the command return full information or summary")

	inspectCmds.AddCommand(&cobra.Command{
		Use:   "tasks [Machine] [Task Name] ([task index of previous tasks])",
		Short: "Shows the job info for the named task on the named machine",
		Long: `Lists the job information about the specified task on the specified machine.
Additionally, an optional index can be used to specify the nth previous instance.
Specifying --full will also show the log contents with the file.
`,
		Args: func(c *cobra.Command, args []string) error {
			if len(args) != 2 && len(args) != 3 {
				return fmt.Errorf("inspect tasks requires 2 or 3 arguments")
			}
			return nil
		},
		RunE: func(c *cobra.Command, args []string) error {
			machineID := args[0]
			taskName := args[1]
			index := -1
			if len(args) == 3 {
				var e error
				index, e = strconv.Atoi(args[2])
				if e != nil {
					return fmt.Errorf("%s is not a valid number: %v", args[2], e)
				}
			}

			mu, err := lookupMachineID(machineID)
			if err != nil {
				return err
			}

			req := Session.Req().List("jobs")
			ps := []string{}
			if index != -1 {
				ps = append(ps, "offset", args[2])
				ps = append(ps, "limit", "1")
			}
			ps = append(ps, "Machine", mu)
			ps = append(ps, "Task", taskName)

			req.Params(ps...)
			data := []*models.Job{}
			if err := req.Do(&data); err != nil {
				return err
			}

			if full {
				for _, j := range data {
					uuid := j.Uuid.String()
					var buf bytes.Buffer
					if err := Session.Req().UrlFor("jobs", uuid, "log").Do(&buf); err != nil {
						j.Meta["log"] = fmt.Sprintf("Failed to get log: %v", err)
					} else {
						j.Meta["log"] = buf.String()
					}

					res := models.JobActions{}
					if err := Session.Req().UrlFor("jobs", uuid, "actions").Do(&res); err != nil {
						j.Meta["actions"] = fmt.Sprintf("Failed to get actions for job: %v", err)
						continue
					}
					j.Meta["actionCount"] = fmt.Sprintf("%d", len(res))
					for i, a := range res {
						j.Meta[fmt.Sprintf("action.%d", i)] = a.Content
					}
				}
			}

			return printJobs(data, full)
		},
	})

	inspectCmds.AddCommand(&cobra.Command{
		Use:   "jobs [Machine] ([job index of previous jobs])",
		Short: "Shows the job info for the named task on the named machine",
		Long: `Lists the jobs about on the specified machine.
Additionally, an optional index can be used to specify the nth previous job.
Specifying --full will also show the log contents with the file.
`,
		Args: func(c *cobra.Command, args []string) error {
			if len(args) != 1 && len(args) != 2 {
				return fmt.Errorf("inspect task requires 1 or 2 arguments")
			}
			return nil
		},
		RunE: func(c *cobra.Command, args []string) error {
			machineID := args[0]
			index := -1
			if len(args) == 2 {
				var e error
				index, e = strconv.Atoi(args[1])
				if e != nil {
					return fmt.Errorf("%s is not a valid number: %v", args[1], e)
				}
			}

			mu, err := lookupMachineID(machineID)
			if err != nil {
				return err
			}

			req := Session.Req().List("jobs")
			ps := []string{}
			if index != -1 {
				ps = append(ps, "offset", args[1])
				ps = append(ps, "limit", "1")
			}
			ps = append(ps, "Machine", mu)

			req.Params(ps...)
			data := []*models.Job{}
			if err := req.Do(&data); err != nil {
				return err
			}

			if full {
				for _, j := range data {
					uuid := j.Uuid.String()
					var buf bytes.Buffer
					if err := Session.Req().UrlFor("jobs", uuid, "log").Do(&buf); err != nil {
						j.Meta["log"] = fmt.Sprintf("Failed to get log: %v", err)
					} else {
						j.Meta["log"] = buf.String()
					}

					res := models.JobActions{}
					if err := Session.Req().UrlFor("jobs", uuid, "actions").Do(&res); err != nil {
						j.Meta["actions"] = fmt.Sprintf("Failed to get actions for job: %v", err)
						continue
					}
					j.Meta["actionCount"] = fmt.Sprintf("%d", len(res))
					for i, a := range res {
						j.Meta[fmt.Sprintf("action.%d", i)] = a.Content
					}
				}
			}

			return printJobs(data, full)
		},
	})

	return inspectCmds
}
