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
  "Name": "justine",
  "OptionalParams": null,
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
  ]
}
`
var processJobsJillCreateOutputString string = `{
  "Name": "jill",
  "Tasks": null
}
`

var processJobsAddProfileJillOutputString = `{
  "Address": "192.168.100.110",
  "BootEnv": "local",
  "CurrentTask": 0,
  "Errors": null,
  "Name": "john",
  "Profile": {
    "Name": "",
    "Tasks": null
  },
  "Profiles": [
    "jill"
  ],
  "Runnable": true,
  "Tasks": [],
  "Uuid": "3e7031fe-3062-45f1-835c-92541bc9cbd3"
}
`

var processJobsLocal2UpdateOutputString = `{
  "Available": true,
  "BootParams": "",
  "Errors": null,
  "Initrds": null,
  "Kernel": "",
  "Name": "local2",
  "OS": {
    "Name": "local2"
  },
  "OnlyUnknown": false,
  "OptionalParams": null,
  "RequiredParams": null,
  "Tasks": [
    "jamie"
  ],
  "Templates": [
    {
      "ID": "local-pxelinux.tmpl",
      "Name": "pxelinux",
      "Path": "pxelinux.cfg/{{.Machine.HexAddress}}"
    },
    {
      "ID": "local-elilo.tmpl",
      "Name": "elilo",
      "Path": "{{.Machine.HexAddress}}.conf"
    },
    {
      "ID": "local-ipxe.tmpl",
      "Name": "ipxe",
      "Path": "{{.Machine.Address}}.ipxe"
    }
  ],
  "Validated": true
}
`

var processJobsSetMachineToLocalOutputString = `{
  "Address": "192.168.100.110",
  "BootEnv": "local",
  "CurrentTask": 0,
  "Errors": null,
  "Name": "john",
  "Profile": {
    "Name": "",
    "Tasks": null
  },
  "Profiles": [
    "jill"
  ],
  "Runnable": true,
  "Tasks": [],
  "Uuid": "3e7031fe-3062-45f1-835c-92541bc9cbd3"
}
`
var processJobsSetMachineToLocal2OutputString = `{
  "Address": "192.168.100.110",
  "BootEnv": "local2",
  "CurrentTask": -1,
  "Errors": null,
  "Name": "john",
  "Profile": {
    "Name": "",
    "Tasks": null
  },
  "Profiles": [
    "jill"
  ],
  "Runnable": true,
  "Tasks": [
    "jamie",
    "justine"
  ],
  "Uuid": "3e7031fe-3062-45f1-835c-92541bc9cbd3"
}
`

var processJobsNoArgsString = "Error: drpcli machines processjobs [id] [wait] requires at least 1 argument\n"
var processJobsTooManyArgsString = "Error: drpcli machines processjobs [id] [wait] requires at most 2 arguments\n"
var processJobsMissingMachineString = "Error: machines GET: p1: Not Found\n\n"
var processJobsBadBooleanString = "Error: Error reading wait argument: strconv.ParseBool: parsing \"asga\": invalid syntax\n\n"
var processJobsNoJobsNoWait = "Processing jobs for 3e7031fe-3062-45f1-835c-92541bc9cbd3 (will not wait for new jobs)\n"

var processJobsErrorSuccessString = "Error: Task failed, exiting ...\n\n\n"
var processJobsOutputSuccessString = `RE:
Processing jobs for 3e7031fe-3062-45f1-835c-92541bc9cbd3 \(will not wait for new jobs\)
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
  "BootEnv": "local2",
  "CurrentJob": "[\S\s]*",
  "CurrentTask": 2,
  "Errors": null,
  "Name": "john",
  "Profile": {
    "Name": "",
    "Tasks": null
  },
  "Profiles": null,
  "Runnable": true,
  "Tasks": \[
    "jamie",
    "justine"
  \],
  "Uuid": "3e7031fe-3062-45f1-835c-92541bc9cbd3"
}
`

var processJobsResetToLocalSuccessString = `RE:
{
  "Address": "192.168.100.110",
  "BootEnv": "local",
  "CurrentJob": "[\S\s]*",
  "CurrentTask": 0,
  "Errors": null,
  "Name": "john",
  "Profile": {
    "Name": "",
    "Tasks": null
  },
  "Profiles": null,
  "Runnable": true,
  "Tasks": \[\],
  "Uuid": "3e7031fe-3062-45f1-835c-92541bc9cbd3"
}
`

var processJobsShowFailedMachineString string = `RE:
{
  "Address": "192.168.100.110",
  "BootEnv": "local2",
  "CurrentJob": "[\S\s]*",
  "CurrentTask": 1,
  "Errors": null,
  "Name": "john",
  "Profile": {
    "Name": "",
    "Tasks": null
  },
  "Profiles": \[
    "jill"
  \],
  "Runnable": false,
  "Tasks": \[
    "jamie",
    "justine"
  \],
  "Uuid": "3e7031fe-3062-45f1-835c-92541bc9cbd3"
}
`

var processJobsRunnableString = `{ "Runnable": true }`
var processJobsShowRunnableMachineString = `RE:
{
  "Address": "192.168.100.110",
  "BootEnv": "local2",
  "CurrentJob": "[\S\s]*",
  "CurrentTask": 1,
  "Errors": null,
  "Name": "john",
  "Profile": {
    "Name": "",
    "Tasks": null
  },
  "Profiles": \[
    "jill"
  \],
  "Runnable": true,
  "Tasks": \[
    "jamie",
    "justine"
  \],
  "Uuid": "3e7031fe-3062-45f1-835c-92541bc9cbd3"
}
`

var processJobsOutputSecondPassSuccessString = `RE:
Processing jobs for 3e7031fe-3062-45f1-835c-92541bc9cbd3 \(will not wait for new jobs\)
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

func TestProcessJobsCli(t *testing.T) {
	// Assumes that local is present (from prefs)

	tests := []CliTest{
		CliTest{false, false, []string{"profiles", "create", "jill"}, noStdinString, processJobsJillCreateOutputString, noErrorString},
		CliTest{false, false, []string{"profiles", "create", "jean"}, noStdinString, machineJeanCreate, noErrorString},
		CliTest{false, false, []string{"tasks", "create", "jamie"}, noStdinString, machineJamieCreate, noErrorString},
		CliTest{false, false, []string{"tasks", "create", "-"}, processJobsJustineCreateTaskString, processJobsJustineCreateOutputString, noErrorString},
		CliTest{false, false, []string{"bootenvs", "create", machineLocal2CreateInput}, noStdinString, machineLocal2Create, noErrorString},
		CliTest{false, false, []string{"machines", "create", machineCreateInputString}, noStdinString, machineCreateJohnString, noErrorString},

		CliTest{true, true, []string{"machines", "processjobs"}, noStdinString, noContentString, processJobsNoArgsString},
		CliTest{true, true, []string{"machines", "processjobs", "p1", "p2", "p3"}, noStdinString, noContentString, processJobsTooManyArgsString},
		CliTest{false, true, []string{"machines", "processjobs", "p1"}, noStdinString, noContentString, processJobsMissingMachineString},
		CliTest{false, true, []string{"machines", "processjobs", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "asga"}, noStdinString, noContentString, processJobsBadBooleanString},
		CliTest{false, false, []string{"machines", "processjobs", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "false"}, noStdinString, processJobsNoJobsNoWait, noErrorString},
		CliTest{false, false, []string{"machines", "processjobs", "3e7031fe-3062-45f1-835c-92541bc9cbd3"}, noStdinString, processJobsNoJobsNoWait, noErrorString},

		CliTest{false, false, []string{"machines", "addprofile", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "jill"}, noStdinString, processJobsAddProfileJillOutputString, noErrorString},
		CliTest{false, false, []string{"profiles", "update", "jill", "{ \"Tasks\": [ \"justine\" ] }"}, noStdinString, machineProfileJamieUpdate, noErrorString},
		CliTest{false, false, []string{"bootenvs", "update", "local2", "{ \"Tasks\": [ \"jamie\" ] }"}, noStdinString, processJobsLocal2UpdateOutputString, noErrorString},
		CliTest{false, false, []string{"machines", "bootenv", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "local"}, noStdinString, processJobsSetMachineToLocalOutputString, noErrorString},
		CliTest{false, false, []string{"machines", "bootenv", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "local2"}, noStdinString, processJobsSetMachineToLocal2OutputString, noErrorString},
		CliTest{false, true, []string{"machines", "processjobs", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "--exit-on-failure"}, noStdinString, processJobsOutputSuccessString, processJobsErrorSuccessString},
		CliTest{false, false, []string{"machines", "show", "3e7031fe-3062-45f1-835c-92541bc9cbd3"}, noStdinString, processJobsShowFailedMachineString, noErrorString},
		CliTest{false, false, []string{"machines", "update", "3e7031fe-3062-45f1-835c-92541bc9cbd3", processJobsRunnableString}, noStdinString, processJobsShowRunnableMachineString, noErrorString},
		CliTest{false, false, []string{"machines", "processjobs", "3e7031fe-3062-45f1-835c-92541bc9cbd3"}, noStdinString, processJobsOutputSecondPassSuccessString, noErrorString},

		CliTest{false, false, []string{"machines", "removeprofile", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "jill"}, noStdinString, processJobsRemoveProfileSuccessString, noErrorString},
		CliTest{false, false, []string{"machines", "bootenv", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "local"}, noStdinString, processJobsResetToLocalSuccessString, noErrorString},

		CliTest{false, false, []string{"machines", "destroy", "3e7031fe-3062-45f1-835c-92541bc9cbd3"}, noStdinString, machineDestroyJohnString, noErrorString},
		CliTest{false, false, []string{"profiles", "destroy", "jill"}, noStdinString, "Deleted profile jill\n", noErrorString},
		CliTest{false, false, []string{"profiles", "destroy", "jean"}, noStdinString, "Deleted profile jean\n", noErrorString},
		CliTest{false, false, []string{"bootenvs", "destroy", "local2"}, noStdinString, "Deleted bootenv local2\n", noErrorString},
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
