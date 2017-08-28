package models

import (
	"time"

	"github.com/pborman/uuid"
)

// Job Action is something that job runner will need to do.
// If path is specified, then the runner will place the contents into that location.
// If path is not specified, then the runner will attempt to bash exec the contents.
// swagger:model
type JobAction struct {
	// required: true
	Name string
	// required: true
	Path string
	// required: true
	Content string
}

// swagger:model
type Job struct {
	Validation
	Access
	// The UUID of the job.  The primary key.
	// required: true
	// swagger:strfmt uuid
	Uuid uuid.UUID
	// The UUID of the previous job to run on this machine.
	// swagger:strfmt uuid
	Previous uuid.UUID
	// The machine the job was created for.  This field must be the UUID of the machine.
	// required: true
	// swagger:strfmt uuid
	Machine uuid.UUID
	// The task the job was created for.  This will be the name of the task.
	// read only: true
	Task string
	// The boot environment that the task was created in.
	// read only: true
	BootEnv string
	// The state the job is in.  Must be one of "created", "running", "failed", "finished", "incomplete"
	// required: true
	State string
	// The time the job entered running.
	StartTime time.Time
	// The time the job entered failed or finished.
	EndTime time.Time
	// required: true
	Archived bool
	// DRP Filesystem path to the log for this job
	// read only: true
	LogPath string
}

func (j *Job) Prefix() string {
	return "jobs"
}

func (j *Job) Key() string {
	return j.Uuid.String()
}

func (j *Job) AuthKey() string {
	return j.Machine.String()
}
