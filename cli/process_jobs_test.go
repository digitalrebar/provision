package cli

import (
	"testing"
)

var processJobsNoArgsString = "Error: drpcli machines processjobs [id] [flags] requires at least 1 argument\n"
var processJobsTooManyArgsString = "Error: drpcli machines processjobs [id] [flags] requires at most 1 arguments\n"
var processJobsMissingMachineString = "Error: GET: machines/p1: Not Found\n\n"
var processJobsStageMissingString = "Error: ValidationError: machines/3e7031fe-3062-45f1-835c-92541bc9cbd3: Stage fred does not exist\n\n"
var processYakovErrorSuccessString = "Error: Task failed, exiting ...\n\n\n"

var processJobsNoJobsNoWait = `RE:
Processing jobs for [\S\s]*: .*
`

func TestProcessJobsCli(t *testing.T) {
	actuallyPowerThings = false

	tests := []CliTest{
		// Setup
		CliTest{false, false, []string{"machines", "create", machineCreateInputString}, noStdinString, machineCreateJohnString, noErrorString},

		// Test basic process jobs cli
		CliTest{true, true, []string{"machines", "processjobs"}, noStdinString, noContentString, processJobsNoArgsString},
		CliTest{true, true, []string{"machines", "processjobs", "p1", "p2", "p3"}, noStdinString, noContentString, processJobsTooManyArgsString},
		CliTest{false, true, []string{"machines", "processjobs", "p1"}, noStdinString, noContentString, processJobsMissingMachineString},
		CliTest{false, false, []string{"machines", "processjobs", "3e7031fe-3062-45f1-835c-92541bc9cbd3"}, noStdinString, processJobsNoJobsNoWait, noErrorString},
		CliTest{false, false, []string{"machines", "destroy", "3e7031fe-3062-45f1-835c-92541bc9cbd3"}, noStdinString, machineDestroyJohnString, noErrorString},
	}

	for _, test := range tests {
		testCli(t, test)
	}
}
