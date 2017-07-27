package cli

import (
	"fmt"

	"github.com/digitalrebar/provision/backend"
	"github.com/digitalrebar/provision/client/jobs"
	"github.com/digitalrebar/provision/models"
	"github.com/go-openapi/strfmt"
	"github.com/spf13/cobra"
)

type JobOps struct{}

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
		// GREG: Add helper here. string uuid for machine
		return nil, fmt.Errorf("Invalid type passed to job create")
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

	mo := &JobOps{}
	commands := commonOps(singularName, name, mo)
	res.AddCommand(commands...)
	return res
}
