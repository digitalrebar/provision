package cli

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/digitalrebar/provision/backend"
	"github.com/digitalrebar/provision/client/jobs"
	models "github.com/digitalrebar/provision/genmodels"
	"github.com/go-openapi/strfmt"
	"github.com/pborman/uuid"
	"github.com/spf13/cobra"
)

type JobOps struct{ CommonOps }

func (be JobOps) GetType() interface{} {
	return &models.Job{}
}

func (be JobOps) GetId(obj interface{}) (string, error) {
	job, ok := obj.(*models.Job)
	if !ok || job.UUID == nil {
		return "", fmt.Errorf("Invalid type passed to job create")
	}
	return job.UUID.String(), nil
}

func (be JobOps) GetIndexes() map[string]string {
	b := &backend.Job{}
	ans := map[string]string{}
	for k, v := range b.Indexes() {
		ans[k] = v.Type
	}
	return ans
}

func (be JobOps) List(parms map[string]string) (interface{}, error) {
	params := jobs.NewListJobsParams()
	if listLimit != -1 {
		t1 := int64(listLimit)
		params = params.WithLimit(&t1)
	}
	if listOffset != -1 {
		t1 := int64(listOffset)
		params = params.WithOffset(&t1)
	}
	for k, v := range parms {
		switch k {
		case "UUID":
			params = params.WithUUID(&v)
		case "BootEnv":
			params = params.WithBootEnv(&v)
		case "Task":
			params = params.WithTask(&v)
		case "State":
			params = params.WithState(&v)
		case "Machine":
			params = params.WithMachine(&v)
		case "Archived":
			params = params.WithArchived(&v)
		case "StartTime":
			params = params.WithStartTime(&v)
		case "EndTime":
			params = params.WithEndTime(&v)
		}
	}
	d, e := session.Jobs.ListJobs(params, basicAuth)
	if e != nil {
		return nil, e
	}
	return d.Payload, nil
}

func (be JobOps) Get(id string) (interface{}, error) {
	d, e := session.Jobs.GetJob(jobs.NewGetJobParams().WithUUID(strfmt.UUID(id)), basicAuth)
	if e != nil {
		return nil, e
	}
	return d.Payload, nil
}

func (be JobOps) Create(obj interface{}) (interface{}, error) {
	job, ok := obj.(*models.Job)
	if !ok {
		if s, ok := obj.(string); ok {
			uu := uuid.Parse(s)
			if uu == nil {
				mo := &MachineOps{}
				if answer, err := mo.List(map[string]string{"Name": s}); err != nil {
					return nil, fmt.Errorf("List machine failed: %s", err)
				} else {
					list := answer.([]*models.Machine)
					if len(list) != 1 {
						return nil, fmt.Errorf("Invalid machine name passed to job create: %s", s)
					}
					m := list[0]

					job = &models.Job{}
					job.Machine = m.UUID
				}
			} else {
				job = &models.Job{}
				u := strfmt.UUID(s)
				job.Machine = &u
			}
		} else {
			return nil, fmt.Errorf("Invalid type passed to job create")
		}
	}
	newJob, oldJob, _, e := session.Jobs.CreateJob(jobs.NewCreateJobParams().WithBody(job), basicAuth)
	if e != nil {
		return nil, e
	}
	if newJob != nil {
		return newJob.Payload, nil
	}
	if oldJob != nil {
		return oldJob.Payload, nil
	}
	return nil, nil
}

func (be JobOps) Patch(id string, obj interface{}) (interface{}, error) {
	data, ok := obj.(models.Patch)
	if !ok {
		return nil, fmt.Errorf("Invalid type passed to job patch")
	}
	d, e := session.Jobs.PatchJob(jobs.NewPatchJobParams().WithUUID(strfmt.UUID(id)).WithBody(data), basicAuth)
	if e != nil {
		return nil, e
	}
	return d.Payload, nil
}

func (be JobOps) Delete(id string) (interface{}, error) {
	d, e := session.Jobs.DeleteJob(jobs.NewDeleteJobParams().WithUUID(strfmt.UUID(id)), basicAuth)
	if e != nil {
		return nil, e
	}
	return d.Payload, nil
}

func init() {
	tree := addJobCommands()
	App.AddCommand(tree)
}

func addJobCommands() (res *cobra.Command) {
	singularName := "job"
	name := "jobs"
	d("Making command tree for %v\n", name)
	res = &cobra.Command{
		Use:   name,
		Short: fmt.Sprintf("Access CLI commands relating to %v", name),
	}

	mo := &JobOps{CommonOps{Name: name, SingularName: singularName}}
	commands := commonOps(mo)

	commands = append(commands, &cobra.Command{
		Use:   "actions [id]",
		Short: "Get the actions for this job",
		RunE: func(c *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("%v requires 1 argument", c.UseLine())
			}
			uuid := args[0]
			dumpUsage = false
			if resp, err := session.Jobs.GetJobActions(jobs.NewGetJobActionsParams().WithUUID(strfmt.UUID(uuid)), basicAuth); err != nil {
				return generateError(err, "Error running action")
			} else {
				return prettyPrint(resp.Payload)
			}
		},
	})

	commands = append(commands, &cobra.Command{
		Use:   "log [id] [- or string]",
		Short: "Gets the log or appends to the log if a second argument or stream is given",
		RunE: func(c *cobra.Command, args []string) error {
			if len(args) < 1 {
				return fmt.Errorf("%v requires at least 1 argument", c.UseLine())
			}
			if len(args) > 2 {
				return fmt.Errorf("%v requires at most 2 arguments", c.UseLine())
			}
			uuid := strfmt.UUID(args[0])
			dumpUsage = false

			if len(args) == 2 {
				var src io.Reader
				if args[1] == "-" {
					src = os.Stdin
				} else {
					buf := bytes.NewBufferString(args[1])
					src = buf
				}
				if _, err := session.Jobs.PutJobLog(
					jobs.NewPutJobLogParams().
						WithUUID(strfmt.UUID(uuid)).
						WithBody(src), basicAuth); err != nil {
					return generateError(err, "Error appending log")
				} else {
					fmt.Println("Success")
					return nil
				}
			} else {
				b := bytes.NewBuffer(nil)
				if _, err := session.Jobs.GetJobLog(jobs.NewGetJobLogParams().WithUUID(uuid), basicAuth, b); err != nil {
					return generateError(err, "Error get log")
				} else {

					fmt.Printf("%s", string(b.Bytes()))
					return nil
				}
			}
		},
	})

	res.AddCommand(commands...)
	return res
}
