package cli

import (
	"io/ioutil"
	"os"
	"testing"
)

var processJobsNoArgsString = "Error: drpcli machines processjobs [id] [flags] requires at least 1 argument\n"
var processJobsTooManyArgsString = "Error: drpcli machines processjobs [id] [flags] requires at most 1 arguments\n"
var processJobsMissingMachineString = "Error: GET: machines/p1: Not Found\n\n"
var processJobsStageMissingString = "Error: ValidationError: Stage fred does not exist\n\n"
var processYakovErrorSuccessString = "Error: Task failed, exiting ...\n\n\n"

var processJobsJustineCreateTaskString = `
---
Meta:
  feature-flags: sane-exit-codes
Name: justine
Templates:
  - Contents: |
      test.txt Content

      Here
    Name: "part 1 - copy file"
    Path: /tmp/test.txt
  - Contents: |
      #!/usr/bin/env bash
      date
      echo
      exit 0
    Name: "part 2 - Print Date - success"
  - Contents: |
      #!/usr/bin/env bash
      . helper
      if [ -e /tmp/incomplete.txt ]; then
         touch /tmp/failed.txt
         rm /tmp/incomplete.txt
         echo "Return failed"
         exit 1
      elif [ -e /tmp/stop.txt ]; then
         echo "Final success"
         exit 0
      elif [ -e /tmp/failed.txt ]; then
         touch /tmp/stop.txt
         rm /tmp/failed.txt
         echo "Return stop"
         exit_stop
      fi
      touch /tmp/incomplete.txt
      echo "Return incomplete"
      exit_incomplete
    Name: "part 3 - Test Return Codes"
  - Contents: |
      #!/usr/bin/env bash
      date
      echo
      exit 0
    Name: "part 4 - Final - success"
`
var processJobsJustineCreateOutputString string = `{
  "Available": true,
  "Errors": [],
  "Meta": {
    "feature-flags": "sane-exit-codes"
  },
  "Name": "justine",
  "OptionalParams": [],
  "ReadOnly": false,
  "RequiredParams": [],
  "Templates": [
    {
      "Contents": "test.txt Content\n\nHere\n",
      "Name": "part 1 - copy file",
      "Path": "/tmp/test.txt"
    },
    {
      "Contents": "#!/usr/bin/env bash\ndate\necho\nexit 0\n",
      "Name": "part 2 - Print Date - success",
      "Path": ""
    },
    {
      "Contents": "#!/usr/bin/env bash\n. helper\nif [ -e /tmp/incomplete.txt ]; then\n   touch /tmp/failed.txt\n   rm /tmp/incomplete.txt\n   echo \"Return failed\"\n   exit 1\nelif [ -e /tmp/stop.txt ]; then\n   echo \"Final success\"\n   exit 0\nelif [ -e /tmp/failed.txt ]; then\n   touch /tmp/stop.txt\n   rm /tmp/failed.txt\n   echo \"Return stop\"\n   exit_stop\nfi\ntouch /tmp/incomplete.txt\necho \"Return incomplete\"\nexit_incomplete\n",
      "Name": "part 3 - Test Return Codes",
      "Path": ""
    },
    {
      "Contents": "#!/usr/bin/env bash\ndate\necho\nexit 0\n",
      "Name": "part 4 - Final - success",
      "Path": ""
    }
  ],
  "Validated": true
}
`

var processJobsYakovCreateTaskString = `{
  "Meta": {
    "feature-flags": "original-exit-codes"
  },
  "Name": "yakov",
  "Templates": [
    {
      "Contents": "test.txt Content\n\nHere\n",
      "Name": "part 1 - copy file",
      "Path": "/tmp/test.txt"
    },
    {
      "Contents": "#!/usr/bin/env bash\n\ndata\necho\nexit 0\n",
      "Name": "part 2 - Print Date - success",
      "Path": ""
    },
    {
      "Contents": "#!/usr/bin/env bash\n\nif [ -e /tmp/incomplete.txt ] ; then\ntouch /tmp/failed.txt\nrm /tmp/incomplete.txt\necho \"Return failed\"\nexit 4\nelif [ -e /tmp/failed.txt ] ; then\necho \"Return success\"\nexit 0\nfi\n\ntouch /tmp/incomplete.txt\necho \"Return incomplete\"\nexit 2\n\n",
      "Name": "part 3 - Test Return Codes",
      "Path": ""
    }
  ]
}
`
var processJobsYakovCreateOutputString string = `{
  "Available": true,
  "Errors": [],
  "Meta": {
    "feature-flags": "original-exit-codes"
  },
  "Name": "yakov",
  "OptionalParams": [],
  "ReadOnly": false,
  "RequiredParams": [],
  "Templates": [
    {
      "Contents": "test.txt Content\n\nHere\n",
      "Name": "part 1 - copy file",
      "Path": "/tmp/test.txt"
    },
    {
      "Contents": "#!/usr/bin/env bash\n\ndata\necho\nexit 0\n",
      "Name": "part 2 - Print Date - success",
      "Path": ""
    },
    {
      "Contents": "#!/usr/bin/env bash\n\nif [ -e /tmp/incomplete.txt ] ; then\ntouch /tmp/failed.txt\nrm /tmp/incomplete.txt\necho \"Return failed\"\nexit 4\nelif [ -e /tmp/failed.txt ] ; then\necho \"Return success\"\nexit 0\nfi\n\ntouch /tmp/incomplete.txt\necho \"Return incomplete\"\nexit 2\n\n",
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
    "Errors": [],
    "Name": "",
    "ReadOnly": false,
    "Validated": false
  },
  "Profiles": [],
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

var yakovCreateInputString string = `{
  "Address": "192.168.100.111",
  "name": "yakov",
  "Uuid": "3e7031fe-3062-45f1-835c-92541bc9cbd4",
  "Bootenv": "local",
  "Secret": "secret2",
  "Tasks": [
    "yakov"
  ],
}
`

var processYakovSetMachineToLocalOutputString = `{
  "Address": "192.168.100.111",
  "Available": true,
  "BootEnv": "local",
  "CurrentTask": -1,
  "Errors": [],
  "Name": "yakov",
  "Profile": {
    "Available": false,
    "Errors": [],
    "Name": "",
    "ReadOnly": false,
    "Validated": false
  },
  "Profiles": [],
  "ReadOnly": false,
  "Runnable": true,
  "Secret": "secret2",
  "Stage": "none",
  "Tasks": [
    "yakov"
  ],
  "Uuid": "3e7031fe-3062-45f1-835c-92541bc9cbd4",
  "Validated": true
}
`

var (
	processYakovOutputSuccessString = `RE:
Processing jobs for 3e7031fe-3062-45f1-835c-92541bc9cbd4
Starting Task: yakov \([\S\s]*\)
Putting Content in place for Task Template: part 1 - copy file
Task Template: part 1 - copy file - Copied contents to /tmp/test.txt successfully
Task Template , part 1 - copy file, finished
Running Task Template: part 2 - Print Date - success
Command part 2 - Print Date - success: failed: false, incomplete: false, reboot: false, poweroff: false, stop: false
Task Template , part 2 - Print Date - success, finished
Running Task Template: part 3 - Test Return Codes
Command part 3 - Test Return Codes: failed: false, incomplete: true, reboot: false, poweroff: false, stop: false
Task Template , part 3 - Test Return Codes, incomplete
Task: yakov incomplete
Starting Task: yakov \([\S\s]*\)
Putting Content in place for Task Template: part 1 - copy file
Task Template: part 1 - copy file - Copied contents to /tmp/test.txt successfully
Task Template , part 1 - copy file, finished
Running Task Template: part 2 - Print Date - success
Command part 2 - Print Date - success: failed: false, incomplete: false, reboot: false, poweroff: false, stop: false
Task Template , part 2 - Print Date - success, finished
Running Task Template: part 3 - Test Return Codes
Command part 3 - Test Return Codes: failed: true, incomplete: false, reboot: false, poweroff: false, stop: false
Task Template , part 3 - Test Return Codes, failed
Task: yakov failed
`

	processYakovShowFailedMachineString = `RE:
{
  "Address": "192.168.100.111",
  "Available": true,
  "BootEnv": "local",
  "CurrentJob": "[\S\s]*",
  "CurrentTask": 0,
  "Errors": \[\],
  "Name": "yakov",
  "Profile": {
    "Available": false,
    "Errors": \[\],
    "Name": "",
    "ReadOnly": false,
    "Validated": false
  },
  "Profiles": \[\],
  "ReadOnly": false,
  "Runnable": false,
  "Secret": "secret2",
  "Stage": "none",
  "Tasks": \[
    "yakov"
  \],
  "Uuid": "3e7031fe-3062-45f1-835c-92541bc9cbd4",
  "Validated": true
}
`
	processYakovRunnableString            = `{ "Runnable": true }`
	processYakovShowRunnableMachineString = `RE:
{
  "Address": "192.168.100.111",
  "Available": true,
  "BootEnv": "local",
  "CurrentJob": "[\S\s]*",
  "CurrentTask": 0,
  "Errors": \[\],
  "Name": "yakov",
  "Profile": {
    "Available": false,
    "Errors": \[\],
    "Name": "",
    "ReadOnly": false,
    "Validated": false
  },
  "Profiles": \[\],
  "ReadOnly": false,
  "Runnable": true,
  "Secret": "secret2",
  "Stage": "none",
  "Tasks": \[
    "yakov"
  \],
  "Uuid": "3e7031fe-3062-45f1-835c-92541bc9cbd4",
  "Validated": true
}
`
	processYakovOutputSecondPassSuccessString = `RE:
Processing jobs for 3e7031fe-3062-45f1-835c-92541bc9cbd4
Starting Task: yakov \([\S\s]*\)
Putting Content in place for Task Template: part 1 - copy file
Task Template: part 1 - copy file - Copied contents to /tmp/test.txt successfully
Task Template , part 1 - copy file, finished
Running Task Template: part 2 - Print Date - success
Command part 2 - Print Date - success: failed: false, incomplete: false, reboot: false, poweroff: false, stop: false
Task Template , part 2 - Print Date - success, finished
Running Task Template: part 3 - Test Return Codes
Command part 3 - Test Return Codes: failed: false, incomplete: false, reboot: false, poweroff: false, stop: false
Task Template , part 3 - Test Return Codes, finished
Task: yakov finished
Jobs finished
`
)

var processJobsNoJobsNoWait = "Processing jobs for 3e7031fe-3062-45f1-835c-92541bc9cbd3\n"

var processJobsErrorSuccessString = "Error: Task failed, exiting ...\n\n\n"
var processJobsOutputSuccessString = `RE:
Processing jobs for 3e7031fe-3062-45f1-835c-92541bc9cbd3
Starting Task: jamie \([\S\s]*\)
Task: jamie finished
Starting Task: justine \([\S\s]*\)
Putting Content in place for Task Template: part 1 - copy file
Task Template: part 1 - copy file - Copied contents to /tmp/test.txt successfully
Task Template , part 1 - copy file, finished
Running Task Template: part 2 - Print Date - success
Command part 2 - Print Date - success: failed: false, incomplete: false, reboot: false, poweroff: false, stop: false
Task Template , part 2 - Print Date - success, finished
Running Task Template: part 3 - Test Return Codes
Command part 3 - Test Return Codes: failed: false, incomplete: true, reboot: false, poweroff: false, stop: false
Task Template , part 3 - Test Return Codes, incomplete
Task: justine incomplete
Starting Task: justine \([\S\s]*\)
Putting Content in place for Task Template: part 1 - copy file
Task Template: part 1 - copy file - Copied contents to /tmp/test.txt successfully
Task Template , part 1 - copy file, finished
Running Task Template: part 2 - Print Date - success
Command part 2 - Print Date - success: failed: false, incomplete: false, reboot: false, poweroff: false, stop: false
Task Template , part 2 - Print Date - success, finished
Running Task Template: part 3 - Test Return Codes
Command part 3 - Test Return Codes: failed: true, incomplete: false, reboot: false, poweroff: false, stop: false
Task Template , part 3 - Test Return Codes, failed
Task: justine failed
`

var processJobsOutputSecondPassSuccessString = `RE:
Processing jobs for 3e7031fe-3062-45f1-835c-92541bc9cbd3
Starting Task: justine \([\S\s]*\)
Putting Content in place for Task Template: part 1 - copy file
Task Template: part 1 - copy file - Copied contents to /tmp/test.txt successfully
Task Template , part 1 - copy file, finished
Running Task Template: part 2 - Print Date - success
Command part 2 - Print Date - success: failed: false, incomplete: false, reboot: false, poweroff: false, stop: false
Task Template , part 2 - Print Date - success, finished
Running Task Template: part 3 - Test Return Codes
Command part 3 - Test Return Codes: failed: false, incomplete: false, reboot: false, poweroff: false, stop: true
Task Template , part 3 - Test Return Codes, finished
`
var processJobsOutputThirdPassSuccessString = `RE:
Processing jobs for 3e7031fe-3062-45f1-835c-92541bc9cbd3
Starting Task: justine \([\S\s]*\)
Putting Content in place for Task Template: part 1 - copy file
Task Template: part 1 - copy file - Copied contents to /tmp/test.txt successfully
Task Template , part 1 - copy file, finished
Running Task Template: part 2 - Print Date - success
Command part 2 - Print Date - success: failed: false, incomplete: false, reboot: false, poweroff: false, stop: false
Task Template , part 2 - Print Date - success, finished
Running Task Template: part 3 - Test Return Codes
Command part 3 - Test Return Codes: failed: false, incomplete: false, reboot: false, poweroff: false, stop: false
Task Template , part 3 - Test Return Codes, finished
Running Task Template: part 4 - Final - success
Command part 4 - Final - success: failed: false, incomplete: false, reboot: false, poweroff: false, stop: false
Task Template , part 4 - Final - success, finished
Task: justine finished
Jobs finished
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
    "Errors": \[\],
    "Name": "",
    "ReadOnly": false,
    "Validated": false
  },
  "Profiles": \[\],
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
    "Errors": \[\],
    "Name": "",
    "ReadOnly": false,
    "Validated": false
  },
  "Profiles": \[\],
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
    "Errors": \[\],
    "Name": "",
    "ReadOnly": false,
    "Validated": false
  },
  "Profiles": \[\],
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
    "Errors": \[\],
    "Name": "",
    "ReadOnly": false,
    "Validated": false
  },
  "Profiles": \[\],
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
  "OptionalParams": [],
  "Profiles": [],
  "ReadOnly": false,
  "RequiredParams": [],
  "Tasks": [
    "jamie",
    "justine"
  ],
  "Templates": [],
  "Validated": true
}
`

func TestProcessJobsCli(t *testing.T) {
	actuallyPowerThings = false

	tests := []CliTest{
		// Setup
		CliTest{false, false, []string{"tasks", "create", "jamie"}, noStdinString, machineJamieCreate, noErrorString},
		CliTest{false, false, []string{"tasks", "create", "-"}, processJobsJustineCreateTaskString, processJobsJustineCreateOutputString, noErrorString},
		CliTest{false, false, []string{"tasks", "create", "-"}, processJobsYakovCreateTaskString, processJobsYakovCreateOutputString, noErrorString},
		CliTest{false, false, []string{"machines", "create", machineCreateInputString}, noStdinString, machineCreateJohnString, noErrorString},
		CliTest{false, false, []string{"machines", "create", yakovCreateInputString}, noStdinString, processYakovSetMachineToLocalOutputString, noErrorString},
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
		CliTest{false, false, []string{"machines", "processjobs", "3e7031fe-3062-45f1-835c-92541bc9cbd3"}, noStdinString, processJobsOutputThirdPassSuccessString, noErrorString},
	}

	for _, test := range tests {
		testCli(t, test)
	}
	// Run tasks on yakov

	os.Remove("/tmp/stop.txt")
	os.Remove("/tmp/test.txt")
	os.Remove("/tmp/failed.txt")
	os.Remove("/tmp/incomplete.txt")
	tests = []CliTest{
		CliTest{false, true, []string{"machines", "processjobs", "3e7031fe-3062-45f1-835c-92541bc9cbd4", "--exit-on-failure"}, noStdinString, processYakovOutputSuccessString, processYakovErrorSuccessString},
		CliTest{false, false, []string{"machines", "show", "3e7031fe-3062-45f1-835c-92541bc9cbd4"}, noStdinString, processYakovShowFailedMachineString, noErrorString},
		CliTest{false, false, []string{"machines", "update", "3e7031fe-3062-45f1-835c-92541bc9cbd4", processYakovRunnableString}, noStdinString, processYakovShowRunnableMachineString, noErrorString},
		CliTest{false, false, []string{"machines", "processjobs", "3e7031fe-3062-45f1-835c-92541bc9cbd4"}, noStdinString, processYakovOutputSecondPassSuccessString, noErrorString},

		// Test some other clean up actions
		CliTest{false, true, []string{"machines", "stage", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "fred"}, noStdinString, noContentString, processJobsStageMissingString},

		// Clean Up
		CliTest{false, false, []string{"machines", "stage", "3e7031fe-3062-45f1-835c-92541bc9cbd3", ""}, noStdinString, processJobsResetToLocalSuccessString, noErrorString},
		CliTest{false, false, []string{"machines", "destroy", "3e7031fe-3062-45f1-835c-92541bc9cbd3"}, noStdinString, machineDestroyJohnString, noErrorString},
		CliTest{false, false, []string{"machines", "destroy", "3e7031fe-3062-45f1-835c-92541bc9cbd4"}, noStdinString, "Deleted machine 3e7031fe-3062-45f1-835c-92541bc9cbd4\n", noErrorString},
		CliTest{false, false, []string{"stages", "destroy", "stage1"}, noStdinString, "Deleted stage stage1\n", noErrorString},
		CliTest{false, false, []string{"tasks", "destroy", "jamie"}, noStdinString, "Deleted task jamie\n", noErrorString},
		CliTest{false, false, []string{"tasks", "destroy", "justine"}, noStdinString, "Deleted task justine\n", noErrorString},
		CliTest{false, false, []string{"tasks", "destroy", "yakov"}, noStdinString, "Deleted task yakov\n", noErrorString},
	}

	for _, test := range tests {
		testCli(t, test)
	}

	if bs, err := ioutil.ReadFile("/tmp/test.txt"); err != nil {
		t.Errorf("Failed to read /tmp/test.txt: %v\n", err)
	} else {
		s := string(bs)

		if s != "test.txt Content\n\nHere\n" {
			t.Errorf("Contents of test.txt don't match: %s %s\n", s, "test.txt Content\n\nHere\n")
		}
	}

	os.Remove("/tmp/stop.txt")
	os.Remove("/tmp/test.txt")
	os.Remove("/tmp/failed.txt")
	os.Remove("/tmp/incomplete.txt")
}
