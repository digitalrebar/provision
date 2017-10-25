package api

// Come back to processJobs later

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"
	"syscall"
	"time"

	"github.com/VictorLowther/jsonpatch2"
	"github.com/digitalrebar/provision/models"
)

// JobLog gets the log for a specific Job and writes it to the passed
// io.Writer
func (c *Client) JobLog(j *models.Job, dst io.Writer) error {
	return c.Req().UrlFor("jobs", j.Key(), "log").Do(dst)
}

// JobActions returns the expanded list of templates that should be
// written or executed for a specific Job.
func (c *Client) JobActions(j *models.Job) ([]*models.JobAction, error) {
	res := []*models.JobAction{}
	return res, c.Req().UrlFor("jobs", j.Key(), "actions").Do(&res)
}

// TaskRunner is responsible for expanding templates and running
// scripts for a single task.
type TaskRunner struct {
	// Status codes that may be returned when a script exits.
	failed, incomplete, reboot, poweroff, stop bool
	// Client that the TaskRunner will use to communicate with the API
	c *Client
	// The Job that the TaskRunner will log to and update the status of.
	j *models.Job
	// The machine the TaskRunner is running on.
	m *models.Machine
	// The machine's current Task.
	t *models.Task
	// The io.Writer that all logging output goes to.
	// It writes to stderr and to the Job on the server.
	in io.Writer
	// The write side of the pipe that communicates to the servver.
	// Closing this will flush any data left in the pipe.
	pipeWriter       *io.PipeWriter
	agentDir, jobDir string
}

// NewTaskRunner creates a new TaskRunner for the passed-in machine.
// It creates the matching Job (or resumes the previous incomplete
// one), and handles making sure that all relavent output is written
// to the job log as well as local stderr
func NewTaskRunner(c *Client, m *models.Machine, agentDir string) (*TaskRunner, error) {
	res := &TaskRunner{
		c:        c,
		m:        m,
		agentDir: agentDir,
	}
	job := &models.Job{Machine: m.Uuid}
	if err := c.CreateModel(job); err != nil {
		return nil, err
	}
	if job.State == "" {
		// Nothing to do.  Not an error
		return nil, nil
	}
	if job.State != "created" && job.State != "incomplete" {
		err := &models.Error{
			Type:  "CLIENT_ERROR",
			Model: job.Prefix(),
			Key:   job.Key(),
		}
		err.Errorf("Invalid job state returned: %v", job.State)
		return nil, err
	}
	t := &models.Task{Name: job.Task}
	if err := c.FillModel(t, job.Task); err != nil {
		return nil, err
	}
	res.j = job
	res.t = t
	return res, nil
}

// Close() shuts down the writer side of the logging pipe.
// This will also flush any remaining data to stderr
func (r *TaskRunner) Close() {
	if r.pipeWriter != nil {
		r.pipeWriter.Close()
	}
	os.Stderr.Sync()
}

// Log writes the string (with a timestamp) to stderr and to the
// server-side log for the current job.
func (r *TaskRunner) Log(s string, items ...interface{}) {
	fmt.Fprintf(r.in, time.Now().String()+": "+s+"\n", items...)
}

// Expand a writes a file template to the appropriate location.
func (r *TaskRunner) Expand(action *models.JobAction) error {
	// Write the Contents of this template to the passed Path
	if !strings.HasPrefix(action.Path, "/") {
		action.Path = path.Join(r.agentDir, path.Clean(action.Path))
	}
	r.Log("%s: Writing %s to %s", time.Now(), action.Name, action.Path)
	if err := ioutil.WriteFile(action.Path, []byte(action.Content), 0644); err != nil {
		r.Log("Unable to write to %s: %v", action.Path, err)
		return err
	}
	return nil
}

// Perform runs a single script action.
func (r *TaskRunner) Perform(action *models.JobAction, taskDir string) error {
	taskFile := path.Join(taskDir, r.j.Task+"-"+action.Name)
	if err := ioutil.WriteFile(taskFile, []byte(action.Content), 0700); err != nil {
		r.Log("Unable to write to script %s: %v", taskFile, err)
		return err
	}
	cmd := exec.Command("./" + path.Base(taskFile))
	cmd.Dir = taskDir
	cmd.Env = append(os.Environ(), "RS_TASK_DIR="+taskDir, "RS_RUNNER_DIR="+r.agentDir)
	cmd.Stdout = r.in
	cmd.Stderr = r.in
	r.Log("%s: Starting command %s\n\n", time.Now(), cmd.Path)
	if err := cmd.Start(); err != nil {
		r.Log("%s: Command failed to start")
	}
	// Wait on the process, not the command to exit.
	// We don't want to auto-close stdout and stderr,
	// as we will continue to use them.
	pState, _ := cmd.Process.Wait()
	status := pState.Sys().(syscall.WaitStatus)
	sane := r.t.HasFeature("sane-exit-codes")
	if !sane {
		st, err := os.Stat(path.Join(taskDir, ".sane-exit-codes"))
		sane = err == nil && st.Mode().IsRegular()
	}
	code := uint(status.ExitStatus())
	if sane {
		// codes can be between 0 and 255
		// if the low bits are not 0, the command failed.
		r.failed = code&^240 > uint(0)
		// If the high bit is set, the command was incomplete and the
		// the current task pointer should not be advanced.
		r.incomplete = code&128 > uint(0)
		// If we need a reboot, set bit 6
		r.reboot = code&64 > uint(0)
		// If we need to poweroff, set bit 5.  Reboot wins if it is set.
		r.poweroff = code&32 > uint(0)
		// If we need to stop, set bit 4
		r.stop = code&16 > uint(0)
	} else {
		switch code {
		case 0:
		case 1:
			r.reboot = true
		case 2:
			r.incomplete = true
		case 3:
			r.incomplete = true
			r.reboot = true
		default:
			r.failed = true
		}
	}
	return nil
}

// Run loops over all of the actions for a particular job,
// placing files and executing scripts as appropriate.
// It also arranges for all logging output for the actions
// to go to the right places.
func (r *TaskRunner) Run() error {
	finalErr := &models.Error{
		Type:  "RUNNER_ERR",
		Model: r.j.Prefix(),
		Key:   r.j.Key(),
	}
	// Arrange to log everything to the job log and stderr at the same time.
	// Due to how io.Pipe works, this should wind up being fairly synchronous.
	reader, writer := io.Pipe()
	r.in = io.MultiWriter(writer, os.Stderr)
	r.pipeWriter = writer
	helperWritten := false
	go func() {
		defer reader.Close()
		buf := bytes.NewBuffer(make([]byte, 64*1024))
		for {
			count, err := io.CopyN(buf, reader, 64*1024)
			if count > 0 {
				if r.c.Req().Put(buf).UrlFor("jobs", r.j.Key(), "log").Do(nil) != nil {
					return
				}
				buf.Reset()
			}
			if err != nil {
				return
			}
		}
	}()
	// We are responsible for going from created to running.
	// If this patch fails, we cannot do it
	patch := jsonpatch2.Patch{
		{Op: "test", Path: "/State", Value: "created"},
		{Op: "replace", Path: "/State", Value: "running"},
	}
	finalState := "failed"
	taskDir, err := ioutil.TempDir(r.agentDir, r.j.Task+"-")
	if err != nil {
		r.Log("Failed to create local tmpdir: %v", err)
		finalErr.AddError(err)
		return finalErr
	}
	// No matter how the function exits, we will try to patch the Job
	// to an appropriate final state.
	defer os.RemoveAll(taskDir)
	defer func() {
		finalPatch := jsonpatch2.Patch{
			{Op: "test", Path: "/State", Value: "running"},
			{Op: "replace", Path: "/State", Value: finalState},
		}
		obj, err := r.c.PatchModel(r.j.Prefix(), r.j.Key(), finalPatch)
		if err != nil {
			r.j = obj.(*models.Job)
		}
	}()
	obj, err := r.c.PatchModel(r.j.Prefix(), r.j.Key(), patch)
	if err != nil {
		finalErr.AddError(err)
		return finalErr
	}
	r.j = obj.(*models.Job)
	r.Log("Starting task %s on %s", r.j.Task, r.m.Uuid)
	// At this point, we are running.
	actions, err := r.c.JobActions(r.j)
	if err != nil {
		r.Log("Failed to render actions: %v", err)
		finalErr.AddError(err)
		return finalErr
	}
	for _, action := range actions {
		r.failed = false
		r.incomplete = false
		r.poweroff = false
		r.reboot = false
		r.stop = false
		var err error
		if action.Path != "" {
			err = r.Expand(action)
		} else {
			if !helperWritten {
				err := ioutil.WriteFile(path.Join(taskDir, "helper"), cmdHelper, 0600)
				if err != nil {
					finalErr.AddError(err)
					return finalErr
				}
				helperWritten = true
			}
			err = r.Perform(action, taskDir)
			// Contents is a script to run, run it.
		}
		if err != nil {
			finalErr.AddError(err)
			return finalErr
		}
		if r.failed {
			finalState = "failed"
		} else if r.incomplete {
			finalState = "incomplete"
		}
		if r.failed || r.incomplete || r.reboot || r.poweroff || r.stop {
			break
		}
	}
	if finalState == "running" {
		finalState = "finished"
	}
	return nil
}

// Agent runs the machine Agent on the current machine.
// It assumes there is only one Agent, which is not actually a safe assumption.
// We should make it safe someday.
func (c *Client) Agent(m *models.Machine, exitOnFailure, actuallyPowerThings bool) error {
	// Clear the current running job, if any
	currentJob := &models.Job{}
	if c.FillModel(currentJob, m.CurrentJob.String()) == nil {
		if currentJob.State == "running" || currentJob.State == "created" {
			currentJob.State = "failed"
			if err := c.PutModel(currentJob); err != nil {
				return err
			}
		}
	}
	runnerDir, err := ioutil.TempDir("", "runner-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(runnerDir)
	events, err := c.Events()
	if err != nil {
		return err
	}
	defer events.Close()
	newM := models.Clone(m).(*models.Machine)
	newM.Runnable = true
	obj, err := c.PatchTo(m, newM)
	if err != nil {
		return err
	}
	m = obj.(*models.Machine)
	var runner *TaskRunner
	for {
		if runner != nil {
			runner.Close()
			runner = nil
		}
		found, err := events.WaitFor(m, TestItem("Runnable", "true"), 1*time.Hour)
		if err != nil {
			res := &models.Error{
				Type:  "AGENT_WAIT",
				Model: m.Prefix(),
				Key:   m.Key(),
			}
			res.Errorf("Event wait failed: %s", found)
			res.AddError(err)
			return res
		}
		switch found {
		case "timeout":
			continue
		case "interrupt":
			break
		case "complete":
		default:
			res := &models.Error{
				Type:  "AGENT_WAIT",
				Model: m.Prefix(),
				Key:   m.Key(),
			}
			res.Errorf("Unexpected return from WaitFor: %s", found)
			return res
		}
		runner, err = NewTaskRunner(c, m, runnerDir)
		if err != nil {
			return err
		}
		if err := runner.Run(); err != nil {
			return err
		}
		if runner.reboot ||
			runner.poweroff ||
			runner.stop ||
			runner.incomplete ||
			(exitOnFailure && runner.failed) {
			break
		}
	}
	if runner == nil {
		return nil
	}
	defer runner.Close()
	if runner.reboot {
		if actuallyPowerThings {
			_, err := exec.Command("reboot").Output()
			if err != nil {
				runner.Log("Failed to issue reboot: %v", err)
			}
		} else {
			runner.Log("Would have rebooted")
		}
	}
	if runner.poweroff {
		if actuallyPowerThings {
			_, err := exec.Command("poweroff").Output()
			if err != nil {
				runner.Log("Failed to issue poweroff: %v", err)
			}
		} else {
			runner.Log("Would have powered down")
		}
	}
	// If we failed, should we exit
	if runner.failed {
		runner.Log("Task failed")
	}
	// Task asked to stop
	if runner.stop {
		runner.Log("Task signalled runner to stop")
	}
	return nil
}
