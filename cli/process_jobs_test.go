package cli

import (
	"testing"
)

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

var processJobsOutputSuccessString = `RE:
Processing jobs for 3e7031fe-3062-45f1-835c-92541bc9cbd3 \(will not wait for new jobs\)
Starting Task: jamie \([\S\s]*\)
Task: jamie finished
Starting Task: justine \([\S\s]*\)
Task: justine finished
Jobs finished
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

func TestProcessJobsCli(t *testing.T) {
	// Assumes that local is present (from prefs)

	tests := []CliTest{
		CliTest{false, false, []string{"profiles", "create", "jill"}, noStdinString, machineJillCreate, noErrorString},
		CliTest{false, false, []string{"profiles", "create", "jean"}, noStdinString, machineJeanCreate, noErrorString},
		CliTest{false, false, []string{"tasks", "create", "jamie"}, noStdinString, machineJamieCreate, noErrorString},
		CliTest{false, false, []string{"tasks", "create", "justine"}, noStdinString, machineJustineCreate, noErrorString},
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
		CliTest{false, false, []string{"machines", "processjobs", "3e7031fe-3062-45f1-835c-92541bc9cbd3"}, noStdinString, processJobsOutputSuccessString, noErrorString},
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

}
