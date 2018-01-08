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
	"github.com/VictorLowther/jsonpatch2/utils"
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
	logger           io.Writer
}

// NewTaskRunner creates a new TaskRunner for the passed-in machine.
// It creates the matching Job (or resumes the previous incomplete
// one), and handles making sure that all relavent output is written
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
	r.in = io.MultiWriter(writer, r.logger)
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
				r.Log("Marked machine %s as not runnable", r.m.Key())
				r.m = newM
			} else {
				r.Log("Failed to mark machine %s as not runnable: %v", r.m.Key(), err)
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
			r.Log("Failed to update job %s to its final state %s", r.j.Key(), finalState)
		} else {
			r.Log("Updated job %s to %s", r.j.Key(), finalState)
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
			finalErr.AddError(err)
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

//
// changeStage takes a machine and attempts to change its stage and set return flags.
// It will issue reboot if needed and force flags to stop = true
//
func (c *Client) changeStage(im *models.Machine, actuallyPowerThings bool, logger io.Writer) (m *models.Machine, wait, stop, changed bool, err error) {
	m = im

	// Get the following data:
	// - stop - should we stop
	// - wait - should we wait for more task in the current stage.
	// - nextStage - what is the next stage we should go to.
	// - reboot - should the machine be rebooted on this stage change.
	//
	// Do this:
	//   if we have a nextStage, then:
	//     set the machine's stage to the next stage.
	//     if we reboot or the new stage says to reboot, then reboot the machine.
	//
	reboot := false
	nextStage := ""
	currentStage := m.Stage

	// Get current Stage
	cs := &models.Stage{Name: currentStage}
	if err = c.Req().Fill(cs); err != nil {
		return
	}
	wait = cs.RunnerWait

	var cmObj interface{}
	csMap := map[string]string{}
	if csErr := c.Req().Get().UrlForM(m, "params", "change-stage/map").Params("aggregate", "true").Do(&cmObj); csErr == nil {
		if err = utils.Remarshal(cmObj, &csMap); err != nil {
			return
		}
	}

	if ns, ok := csMap[currentStage]; ok {
		pieces := strings.Split(ns, ":")
		nextStage = pieces[0]
		if len(pieces) > 1 && pieces[1] == "Reboot" {
			reboot = true
		}
		if len(pieces) > 1 && pieces[1] == "Stop" {
			stop = true
		}
	} else {
		// if current stage ends in -install and no stage map entry, then we need to set the
		// next stage to local.  This makes the code work like the old ce-bootenvs did.
		if strings.HasSuffix(currentStage, "-install") {
			nextStage = "local"
			reboot = false
			wait = false
			stop = true
		} else {
			nextStage = ""
		}
	}

	// If no stage, then just return.
	if nextStage == "" {
		return
	}

	// Get the new stage
	ns := &models.Stage{Name: nextStage}
	if err = c.Req().Fill(ns); err != nil {
		return
	}
	if !reboot {
		reboot = ns.Reboot
	}

	// Change stage
	newM := models.Clone(m).(*models.Machine)
	newM.Stage = nextStage
	obj, err := c.PatchTo(m, newM)
	if err != nil {
		return
	}
	m = obj.(*models.Machine)
	changed = true

	// Reboot if needed
	if reboot {
		// Reboot implies stopping the runner
		wait = false
		stop = true
		if actuallyPowerThings {
			var actionObj interface{}
			if err = c.Req().Get().UrlForM(m, "actions", "nextbootpxe").Do(&actionObj); err == nil {
				emptyMap := map[string]interface{}{}
				var results interface{}
				if err = c.Req().Post(emptyMap).UrlForM(m, "actions", "nextbootpxe").Do(&results); err != nil {
					return
				}
			} else {
				err = nil
			}

			if _, err = exec.Command("reboot").Output(); err != nil {
				return
			}
		} else {
			fmt.Fprintf(logger, "Would have rebooted on stage change")
		}
	}

	return
}

// Agent runs the machine Agent on the current machine.
// It assumes there is only one Agent, which is not actually a safe assumption.
// We should make it safe someday.
func (c *Client) Agent(m *models.Machine, exitOnNotRunnable, exitOnFailure, actuallyPowerThings bool, logger io.Writer) error {
	fmt.Fprintf(logger, "Processing jobs for %s: %s\n", m.Key(), time.Now())

	// Clear the current running job, if any
	currentJob := &models.Job{Uuid: m.CurrentJob}
	if c.Req().Fill(currentJob) == nil {
		if currentJob.State == "running" || currentJob.State == "created" {
			cj := models.Clone(currentJob).(*models.Job)
			cj.State = "failed"
			if _, err := c.PatchTo(currentJob, cj); err != nil {
				return err
			}
		}
	}
	if logger == nil {
		logger = os.Stderr
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

	if m.HasFeature("original-change-stage") || !m.HasFeature("change-stage-v2") {
		newM := models.Clone(m).(*models.Machine)
		newM.Runnable = true
		if err := c.Req().PatchTo(m, newM).Do(&newM); err == nil {
			m = newM
		} else {
			res := &models.Error{
				Type:  "AGENT_WAIT",
				Model: m.Prefix(),
				Key:   m.Key(),
			}
			res.Errorf("Failed to mark machine runnable.")
			res.AddError(err)
			return res
		}
	}

	var runner *TaskRunner
	for {
		if runner != nil {
			runner.Close()
			runner = nil
		}

		found, err := events.WaitFor(m, EqualItem("Runnable", true), 1*time.Hour)
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

		runner, err = NewTaskRunner(c, m, runnerDir, logger)
		if err != nil {
			return err
		}
		if runner == nil {
			// changeStage may have changed stage and rebooted on us.
			// if it rebooted, it will set stop to true and we should just exit.
			//
			// if stop == true, change stage wants us to leave.
			// if changed == true and stop == false, then change_stage wants us to run more tasks regardless.
			// if all three are false, then we should leave with nothing to do.
			// else we are waiting for machine state to change more tasks.
			//
			if nm, wait, stop, changed, err := c.changeStage(m, actuallyPowerThings, logger); err != nil {
				return err
			} else if stop {
				break
			} else if changed {
				// Continue without waiting
				m = nm
				continue
			} else if !wait {
				break
			} else {
				// Wait
				m = nm
			}

			if exitOnNotRunnable {
				return nil
			}

			// wait here for the task list or the current task to change on the machine
			found, err := events.WaitFor(m,
				OrItems(NotItem(EqualItem("CurrentTask", m.CurrentTask)),
					NotItem(EqualItem("Tasks", m.Tasks))), 1*time.Hour)
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
			case "interrupt":
				break
			case "complete":
			case "timeout":
			default:
				res := &models.Error{
					Type:  "AGENT_WAIT",
					Model: m.Prefix(),
					Key:   m.Key(),
				}
				res.Errorf("Unexpected return from WaitFor: %s", found)
				return res
			}

			continue
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
