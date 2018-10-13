package api

// Come back to processJobs later

import (
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
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
func (c *Client) JobActions(j *models.Job, targetOS string) (models.JobActions, error) {
	res := models.JobActions{}
	req := c.Req().UrlFor("jobs", j.Key(), "actions")
	if targetOS != "" {
		req.Params("os", targetOS)
	}
	return res, req.Do(&res)
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
	pipeWriter       net.Conn
	agentDir, jobDir string
	logger           io.Writer
}

// NewTaskRunner creates a new TaskRunner for the passed-in machine.
// It creates the matching Job (or resumes the previous incomplete
// one), and handles making sure that all relevant output is written
// to the job log as well as local stderr
func NewTaskRunner(c *Client, m *models.Machine, agentDir string, logger io.Writer) (*TaskRunner, error) {
	if logger == nil {
		logger = ioutil.Discard
	}
	res := &TaskRunner{
		c:        c,
		m:        m,
		agentDir: agentDir,
		logger:   logger,
	}
	job := &models.Job{Machine: m.Uuid}
	if err := c.CreateModel(job); err != nil && err != io.EOF {
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
		err.Errorf("Job: %#v", job)
		return nil, err
	}
	t := &models.Task{Name: job.Task}
	if err := c.Req().Fill(t); err != nil {
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
	type flusher interface {
		io.Writer
		Flush() error
	}
	type syncer interface {
		io.Writer
		Sync() error
	}
	switch o := r.logger.(type) {
	case flusher:
		o.Flush()
	case syncer:
		o.Sync()
	}
}

// Log writes the string (with a timestamp) to stderr and to the
// server-side log for the current job.
func (r *TaskRunner) Log(s string, items ...interface{}) {
	fmt.Fprintf(r.in, s+"\n", items...)
}

// Expand a writes a file template to the appropriate location.
func (r *TaskRunner) Expand(action *models.JobAction, taskDir string) error {
	// Write the Contents of this template to the passed Path
	if !strings.HasPrefix(action.Path, "/") {
		action.Path = path.Join(taskDir, path.Clean(action.Path))
	}
	r.Log("%s: Writing %s to %s", time.Now(), action.Name, action.Path)
	if err := os.MkdirAll(filepath.Dir(action.Path), os.ModePerm); err != nil {
		r.Log("Unable to mkdirs for %s: %v", action.Path, err)
		return err
	}
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

	cmdArray := []string{}
	if interp, ok := action.Meta["Interpreter"]; ok {
		// This is probably usually not required anywhere but Windows,
		// as basically all Unix shell scripts should start with #!,
		// and even on Windows we will try to guess based on the extension.
		cmdArray = append(cmdArray, interp)
	} else if strings.HasSuffix(taskFile, "ps1") {
		cmdArray = append(cmdArray, "powershell.exe")
		cmdArray = append(cmdArray, "-File")
	}
	cmdArray = append(cmdArray, "./"+path.Base(taskFile))
	cmd := exec.Command(cmdArray[0], cmdArray[1:]...)

	cmd.Dir = taskDir
	cmd.Env = append(os.Environ(), "RS_TASK_DIR="+taskDir, "RS_RUNNER_DIR="+r.agentDir)
	for _, e := range []string{"RS_UUID", "RS_ENDPOINT", "RS_TOKEN"} {
		if os.Getenv(e) == "" {
			cmd.Env = append(cmd.Env, e+"="+r.c.token.Token)
		}
	}
	cmd.Stdout = r.in
	cmd.Stderr = r.in
	r.Log("Starting command %s\n\n", cmd.Path)
	if err := cmd.Start(); err != nil {
		r.Log("Command failed to start: %v", err)
		return err
	}
	// Wait on the process, not the command to exit.
	// We don't want to auto-close stdout and stderr,
	// as we will continue to use them.
	r.Log("Command running")
	pState, _ := cmd.Process.Wait()
	status := pState.Sys().(syscall.WaitStatus)
	sane := r.t.HasFeature("sane-exit-codes")
	if !sane {
		st, err := os.Stat(path.Join(taskDir, ".sane-exit-codes"))
		sane = err == nil && st.Mode().IsRegular()
	}
	code := uint(status.ExitStatus())
	r.Log("Command exited with status %d", code)
	if sane {
		switch code {
		case 0:
		case 16:
			r.stop = true
		case 32:
			r.poweroff = true
		case 64:
			r.reboot = true
		case 128:
			r.incomplete = true
		case 144:
			r.stop = true
			r.incomplete = true
		case 160:
			r.incomplete = true
			r.poweroff = true
		case 192:
			r.incomplete = true
			r.reboot = true
		default:
			r.failed = true
		}
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
	jKey := r.j.Key()
	// Arrange to log everything to the job log and stderr at the same time.
	// Due to how io.Pipe works, this should wind up being fairly synchronous.
	reader, writer := net.Pipe()

	r.in = io.MultiWriter(writer, r.logger)
	r.pipeWriter = writer
	helperWritten := false

	go func() {
		defer reader.Close()
		buf := make([]byte, 1<<16)
		reader.SetReadDeadline(time.Now().Add(1 * time.Second))
		pos := 0
		for {
			count, err := reader.Read(buf[pos:])
			pos += count
			if pos < len(buf) && err == nil {
				continue
			}
			if pos > 0 {
				if r.c.Req().Put(buf[:pos]).UrlFor("jobs", jKey, "log").Do(nil) != nil {
					return
				}
				pos = 0
			}
			if err != nil {
				if os.IsTimeout(err) {
					reader.SetReadDeadline(time.Now().Add(1 * time.Second))
					continue
				}
				return
			}
		}
	}()
	// We are responsible for going from created to running.
	// If this patch fails, we cannot do it
	patch := jsonpatch2.Patch{
		{Op: "test", Path: "/State", Value: r.j.State},
		{Op: "replace", Path: "/State", Value: "running"},
	}
	finalState := "incomplete"
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
		if r.failed || r.reboot || r.stop || r.poweroff || r.incomplete {
			newM := models.Clone(r.m).(*models.Machine)
			newM.Runnable = false
			if err := r.c.Req().PatchTo(r.m, newM).Do(&newM); err == nil {
				r.Log("Marked machine %s as not runnable", r.m.Name)
				r.m = newM
			} else {
				r.Log("Failed to mark machine %s as not runnable: %v", r.m.Name, err)
			}
		}
		exitState := "complete"
		if finalState == "failed" {
			exitState = "failed"
		}
		if r.reboot {
			exitState = "reboot"
		} else if r.poweroff {
			exitState = "poweroff"
		} else if r.stop {
			exitState = "stop"
		}
		finalPatch := jsonpatch2.Patch{
			{Op: "test", Path: "/State", Value: "running"},
			{Op: "replace", Path: "/State", Value: finalState},
			{Op: "replace", Path: "/ExitState", Value: exitState},
		}
		if err := r.c.Req().Patch(finalPatch).UrlForM(r.j).Do(&r.j); err != nil {
			r.Log("Failed to update job %s:%s:%s to its final state %s", r.j.Workflow, r.j.Stage, r.j.Task, finalState)
		} else {
			r.Log("Updated job for %s:%s:%s to %s", r.j.Workflow, r.j.Stage, r.j.Task, finalState)
		}
	}()
	obj, err := r.c.PatchModel(r.j.Prefix(), r.j.Key(), patch)
	if err != nil {
		finalErr.AddError(err)
		return finalErr
	}
	r.j = obj.(*models.Job)
	r.Log("Starting task %s:%s:%s on %s", r.j.Workflow, r.j.Stage, r.j.Task, r.m.Name)
	// At this point, we are running.
	var actions models.JobActions
	if allActions, err := r.c.JobActions(r.j, runtime.GOOS); err != nil {
		r.Log("Failed to render actions: %v", err)
		finalErr.AddError(err)
		return finalErr
	} else {
		actions = allActions.FilterOS(runtime.GOOS)
	}
	for i, action := range actions {
		final := len(actions)-1 == i
		r.failed = false
		r.incomplete = false
		r.poweroff = false
		r.reboot = false
		r.stop = false
		var err error
		if action.Path != "" {
			err = r.Expand(action, taskDir)
		} else {
			if !helperWritten {
				err = ioutil.WriteFile(path.Join(taskDir, "helper"), cmdHelper, 0600)
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
			r.failed = true
			finalState = "failed"
			finalErr.AddError(err)
			r.Log("Task %s %s", r.j.Task, finalState)
			return finalErr
		}
		r.Log("Action %s finished", action.Name)
		// If a non-final action sets the incomplete flag, it actually
		// means early success and stop processing actions for this task.
		// This allows actions to be structured in an "early exit"
		// fashion.
		//
		// Only the final action can actually set things as incomplete.
		if !final && r.incomplete {
			r.incomplete = !r.incomplete
			break
		}
		if r.failed {
			finalState = "failed"
			break
		}
		if r.reboot || r.poweroff || r.stop {
			r.incomplete = !final
			break
		}
	}
	if !r.failed && !r.incomplete {
		finalState = "finished"
	}
	r.Log("Task %s %s", r.j.Task, finalState)
	return nil
}

// Agent runs the machine Agent on the current machine.
// It assumes there is only one Agent, which is not actually a safe assumption.
// We should make it safe someday.
func (c *Client) Agent(m *models.Machine, exitOnNotRunnable, exitOnFailure, actuallyPowerThings bool, logger io.Writer) error {
	a, err := c.NewAgent(m, exitOnNotRunnable, exitOnFailure, actuallyPowerThings, logger)
	if err != nil {
		return err
	}
	return a.Run()
}
