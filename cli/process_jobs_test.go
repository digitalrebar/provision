package cli

import (
	"io/ioutil"
	"os"
	"testing"
)

var processJobsJustineCreateTaskString = `{
  "Name": "justine",
  "Templates": [
    {
      "Contents": "test.txt Content\n\nHere\n",
      "Name": "part 1 - copy file",
      "Path": "test.txt"
    },
    {
      "Contents": "#!/bin/bash\n\ndata\necho\nexit 0\n",
      "Name": "part 2 - Print Date - success",
      "Path": ""
    },
    {
      "Contents": "#!/bin/bash\n\nif [ -e incomplete.txt ] ; then\ntouch failed.txt\nrm incomplete.txt\necho \"Return failed\"\nexit 4\nelif [ -e failed.txt ] ; then\necho \"Return success\"\nexit 0\nfi\n\ntouch incomplete.txt\necho \"Return incomplete\"\nexit 2\n\n",
      "Name": "part 3 - Test Return Codes",
      "Path": ""
    }
  ]
}
`
var processJobsJustineCreateOutputString string = `{
  "Available": true,
  "Errors": [],
  "Name": "justine",
  "OptionalParams": null,
  "ReadOnly": false,
  "RequiredParams": null,
  "Templates": [
    {
      "Contents": "test.txt Content\n\nHere\n",
      "Name": "part 1 - copy file",
      "Path": "test.txt"
    },
    {
      "Contents": "#!/bin/bash\n\ndata\necho\nexit 0\n",
      "Name": "part 2 - Print Date - success",
      "Path": ""
    },
    {
      "Contents": "#!/bin/bash\n\nif [ -e incomplete.txt ] ; then\ntouch failed.txt\nrm incomplete.txt\necho \"Return failed\"\nexit 4\nelif [ -e failed.txt ] ; then\necho \"Return success\"\nexit 0\nfi\n\ntouch incomplete.txt\necho \"Return incomplete\"\nexit 2\n\n",
      "Name": "part 3 - Test Return Codes",
      "Path": ""
    }
  ],
  "Validated": true
}
`

var processJobsSetMachineToLocalOutputString = `{
  "Address": "192.168.100.110",
  "Available": true,
  "BootEnv": "local",
  "CurrentTask": -1,
  "Errors": [],
  "Name": "john",
  "Profile": {
    "Available": false,
    "Errors": null,
    "Name": "",
    "ReadOnly": false,
    "Validated": false
  },
  "Profiles": null,
  "ReadOnly": false,
  "Runnable": true,
  "Secret": "secret1",
  "Stage": "stage1",
  "Tasks": [
    "jamie",
    "justine"
  ],
  "Uuid": "3e7031fe-3062-45f1-835c-92541bc9cbd3",
  "Validated": true
}
`
var processJobsNoArgsString = "Error: drpcli machines processjobs [id] [flags] requires at least 1 argument\n"
var processJobsTooManyArgsString = "Error: drpcli machines processjobs [id] [flags] requires at most 1 arguments\n"
var processJobsMissingMachineString = "Error: machines GET: p1: Not Found\n\n"
var processJobsNoJobsNoWait = "Processing jobs for 3e7031fe-3062-45f1-835c-92541bc9cbd3\n"

var processJobsErrorSuccessString = "Error: Task failed, exiting ...\n\n\n"
var processJobsOutputSuccessString = `RE:
Processing jobs for 3e7031fe-3062-45f1-835c-92541bc9cbd3
Starting Task: jamie \([\S\s]*\)
Task: jamie finished
Starting Task: justine \([\S\s]*\)
Putting Content in place for Task Template: part 1 - copy file
Task Template: part 1 - copy file - Copied contents to test.txt successfully
Task Template , part 1 - copy file, finished
Running Task Template: part 2 - Print Date - success
Command part 2 - Print Date - success succeeded
Task Template , part 2 - Print Date - success, finished
Running Task Template: part 3 - Test Return Codes
Command part 3 - Test Return Codes incomplete
Task Template , part 3 - Test Return Codes, incomplete
Task: justine incomplete
Starting Task: justine \([\S\s]*\)
Putting Content in place for Task Template: part 1 - copy file
Task Template: part 1 - copy file - Copied contents to test.txt successfully
Task Template , part 1 - copy file, finished
Running Task Template: part 2 - Print Date - success
Command part 2 - Print Date - success succeeded
Task Template , part 2 - Print Date - success, finished
Running Task Template: part 3 - Test Return Codes
Command part 3 - Test Return Codes failed
Task Template , part 3 - Test Return Codes, failed
Task: justine failed
`

var processJobsRemoveProfileSuccessString = `RE:
{
  "Address": "192.168.100.110",
  "Available": true,
  "BootEnv": "local",
  "CurrentJob": "[\S\s]*",
  "CurrentTask": 2,
  "Errors": \[\],
  "Name": "john",
  "Profile": {
    "Available": false,
    "Errors": null,
    "Name": "",
    "ReadOnly": false,
    "Validated": false
  },
  "Profiles": null,
  "ReadOnly": false,
  "Runnable": false,
  "Secret": "secret1",
  "Stage": "stage1",
  "Tasks": \[
    "jamie",
    "justine"
  \],
  "Uuid": "3e7031fe-3062-45f1-835c-92541bc9cbd3",
  "Validated": true
}
`

var processJobsResetToLocalSuccessString = `RE:
{
  "Address": "192.168.100.110",
  "Available": true,
  "BootEnv": "local",
  "CurrentJob": "[\S\s]*",
  "CurrentTask": 0,
  "Errors": \[\],
  "Name": "john",
  "Profile": {
    "Available": false,
    "Errors": null,
    "Name": "",
    "ReadOnly": false,
    "Validated": false
  },
  "Profiles": null,
  "ReadOnly": false,
  "Runnable": true,
  "Secret": "secret1",
  "Stage": "none",
  "Tasks": \[\],
  "Uuid": "3e7031fe-3062-45f1-835c-92541bc9cbd3",
  "Validated": true
}
`

var processJobsShowFailedMachineString string = `RE:
{
  "Address": "192.168.100.110",
  "Available": true,
  "BootEnv": "local",
  "CurrentJob": "[\S\s]*",
  "CurrentTask": 1,
  "Errors": \[\],
  "Name": "john",
  "Profile": {
    "Available": false,
    "Errors": null,
    "Name": "",
    "ReadOnly": false,
    "Validated": false
  },
  "Profiles": null,
  "ReadOnly": false,
  "Runnable": false,
  "Secret": "secret1",
  "Stage": "stage1",
  "Tasks": \[
    "jamie",
    "justine"
  \],
  "Uuid": "3e7031fe-3062-45f1-835c-92541bc9cbd3",
  "Validated": true
}
`

var processJobsRunnableString = `{ "Runnable": true }`
var processJobsShowRunnableMachineString = `RE:
{
  "Address": "192.168.100.110",
  "Available": true,
  "BootEnv": "local",
  "CurrentJob": "[\S\s]*",
  "CurrentTask": 1,
  "Errors": \[\],
  "Name": "john",
  "Profile": {
    "Available": false,
    "Errors": null,
    "Name": "",
    "ReadOnly": false,
    "Validated": false
  },
  "Profiles": null,
  "ReadOnly": false,
  "Runnable": true,
  "Secret": "secret1",
  "Stage": "stage1",
  "Tasks": \[
    "jamie",
    "justine"
  \],
  "Uuid": "3e7031fe-3062-45f1-835c-92541bc9cbd3",
  "Validated": true
}
`

var processJobsOutputSecondPassSuccessString = `RE:
Processing jobs for 3e7031fe-3062-45f1-835c-92541bc9cbd3
Starting Task: justine \([\S\s]*\)
Putting Content in place for Task Template: part 1 - copy file
Task Template: part 1 - copy file - Copied contents to test.txt successfully
Task Template , part 1 - copy file, finished
Running Task Template: part 2 - Print Date - success
Command part 2 - Print Date - success succeeded
Task Template , part 2 - Print Date - success, finished
Running Task Template: part 3 - Test Return Codes
Command part 3 - Test Return Codes succeeded
Task Template , part 3 - Test Return Codes, finished
Task: justine finished
Jobs finished
`

var processJobsCreateStage1InputString = `{
  "Name": "stage1",
  "BootEnv": "local",
  "Tasks": [ "jamie", "justine" ]
}`
var processJobsCreateStage1SuccessString = `{
  "Available": true,
  "BootEnv": "local",
  "Errors": [],
  "Name": "stage1",
  "OptionalParams": null,
  "Profiles": [],
  "ReadOnly": false,
  "RequiredParams": null,
  "Tasks": [
    "jamie",
    "justine"
  ],
  "Templates": [],
  "Validated": true
}
`

var processJobsStageMissingString = "Error: Stage fred does not exist\n\n"

func TestProcessJobsCli(t *testing.T) {

	tests := []CliTest{
		// Setup
		CliTest{false, false, []string{"tasks", "create", "jamie"}, noStdinString, machineJamieCreate, noErrorString},
		CliTest{false, false, []string{"tasks", "create", "-"}, processJobsJustineCreateTaskString, processJobsJustineCreateOutputString, noErrorString},
		CliTest{false, false, []string{"machines", "create", machineCreateInputString}, noStdinString, machineCreateJohnString, noErrorString},
		CliTest{false, false, []string{"stages", "create", processJobsCreateStage1InputString}, noStdinString, processJobsCreateStage1SuccessString, noErrorString},

		// Test basic process jobs cli
		CliTest{true, true, []string{"machines", "processjobs"}, noStdinString, noContentString, processJobsNoArgsString},
		CliTest{true, true, []string{"machines", "processjobs", "p1", "p2", "p3"}, noStdinString, noContentString, processJobsTooManyArgsString},
		CliTest{false, true, []string{"machines", "processjobs", "p1"}, noStdinString, noContentString, processJobsMissingMachineString},
		CliTest{false, false, []string{"machines", "processjobs", "3e7031fe-3062-45f1-835c-92541bc9cbd3"}, noStdinString, processJobsNoJobsNoWait, noErrorString},

		// Run a stage
		CliTest{false, false, []string{"machines", "stage", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "stage1"}, noStdinString, processJobsSetMachineToLocalOutputString, noErrorString},
		CliTest{false, true, []string{"machines", "processjobs", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "--exit-on-failure"}, noStdinString, processJobsOutputSuccessString, processJobsErrorSuccessString},
		CliTest{false, false, []string{"machines", "show", "3e7031fe-3062-45f1-835c-92541bc9cbd3"}, noStdinString, processJobsShowFailedMachineString, noErrorString},
		CliTest{false, false, []string{"machines", "update", "3e7031fe-3062-45f1-835c-92541bc9cbd3", processJobsRunnableString}, noStdinString, processJobsShowRunnableMachineString, noErrorString},
		CliTest{false, false, []string{"machines", "processjobs", "3e7031fe-3062-45f1-835c-92541bc9cbd3"}, noStdinString, processJobsOutputSecondPassSuccessString, noErrorString},

		// Test some other clean up actions
		CliTest{false, true, []string{"machines", "stage", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "fred"}, noStdinString, noContentString, processJobsStageMissingString},

		// Clean Up
		CliTest{false, false, []string{"machines", "stage", "3e7031fe-3062-45f1-835c-92541bc9cbd3", ""}, noStdinString, processJobsResetToLocalSuccessString, noErrorString},
		CliTest{false, false, []string{"machines", "destroy", "3e7031fe-3062-45f1-835c-92541bc9cbd3"}, noStdinString, machineDestroyJohnString, noErrorString},
		CliTest{false, false, []string{"stages", "destroy", "stage1"}, noStdinString, "Deleted stage stage1\n", noErrorString},
		CliTest{false, false, []string{"tasks", "destroy", "jamie"}, noStdinString, "Deleted task jamie\n", noErrorString},
		CliTest{false, false, []string{"tasks", "destroy", "justine"}, noStdinString, "Deleted task justine\n", noErrorString},
	}

	for _, test := range tests {
		testCli(t, test)
	}

	if bs, err := ioutil.ReadFile("test.txt"); err != nil {
		t.Errorf("Failed to read test.txt: %v\n", err)
	} else {
		s := string(bs)

		if s != "test.txt Content\n\nHere\n" {
			t.Errorf("Contents of test.txt don't match: %s %s\n", s, "test.txt Content\n\nHere\n")
		}
	}

	os.Remove("test.txt")
	os.Remove("failed.txt")
	os.Remove("incomplete.txt")
}
