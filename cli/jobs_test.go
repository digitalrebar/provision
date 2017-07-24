package cli

import (
	"testing"
)

var jobAddrErrorString string = "Error: Invalid Address: fred\n\n"
var jobExpireTimeErrorString string = "Error: Invalid Address: false\n\n"

var jobDefaultListString string = "[]\n"
var jobEmptyListString string = "[]\n"

var jobShowNoArgErrorString string = "Error: drpcli jobs show [id] requires 1 argument\n"
var jobShowTooManyArgErrorString string = "Error: drpcli jobs show [id] requires 1 argument\n"
var jobShowMissingArgErrorString string = "Error: jobs GET: C0A86467: Not Found\n\n"
var jobShowJohnString string = `{
  "Addr": "192.168.100.100",
  "NextServer": "2.2.2.2",
  "Options": null,
  "Strategy": "MAC",
  "Token": "john"
}
`

var jobExistsNoArgErrorString string = "Error: drpcli jobs exists [id] requires 1 argument"
var jobExistsTooManyArgErrorString string = "Error: drpcli jobs exists [id] requires 1 argument"
var jobExistsIgnoreString string = ""
var jobExistsMissingIgnoreString string = "Error: job get: address not valid: ignore\n\n"

var jobCreateNoArgErrorString string = "Error: drpcli jobs create [json] requires 1 argument\n"
var jobCreateTooManyArgErrorString string = "Error: drpcli jobs create [json] requires 1 argument\n"
var jobCreateBadJSONString = "asdgasdg"
var jobCreateBadJSONErrorString = "Error: Unable to create new job: Invalid type passed to job create\n\n"
var jobCreateInputString string = `{
  "Addr": "192.168.100.100",
  "NextServer": "2.2.2.2",
  "Strategy": "MAC",
  "Token": "john"
}
`
var jobCreateJohnString string = `{
  "Addr": "192.168.100.100",
  "NextServer": "2.2.2.2",
  "Options": null,
  "Strategy": "MAC",
  "Token": "john"
}
`
var jobCreateDuplicateErrorString = "Error: dataTracker create jobs: C0A86464 already exists\n\n"

var jobListJobsString = `[
  {
    "Addr": "192.168.100.100",
    "NextServer": "2.2.2.2",
    "Options": null,
    "Strategy": "MAC",
    "Token": "john"
  }
]
`
var jobListBothEnvsString = `[
  {
    "Addr": "192.168.100.100",
    "NextServer": "2.2.2.2",
    "Options": null,
    "Strategy": "MAC",
    "Token": "john"
  }
]
`

var jobUpdateNoArgErrorString string = "Error: drpcli jobs update [id] [json] requires 2 arguments"
var jobUpdateTooManyArgErrorString string = "Error: drpcli jobs update [id] [json] requires 2 arguments"
var jobUpdateBadJSONString = "asdgasdg"
var jobUpdateBadJSONErrorString = "Error: Unable to merge objects: json: cannot unmarshal string into Go value of type map[string]interface {}\n\n\n"
var jobUpdateInputString string = `{
  "Options": [ { "Code": 3, "Value": "1.1.1.1" } ]
}
`
var jobUpdateJohnString string = `{
  "Addr": "192.168.100.100",
  "NextServer": "2.2.2.2",
  "Options": [
    {
      "Code": 3,
      "Value": "1.1.1.1"
    }
  ],
  "Strategy": "MAC",
  "Token": "john"
}
`
var jobUpdateJohnMissingErrorString string = "Error: jobs GET: C0A86467: Not Found\n\n"

var jobPatchNoArgErrorString string = "Error: drpcli jobs patch [objectJson] [changesJson] requires 2 arguments"
var jobPatchTooManyArgErrorString string = "Error: drpcli jobs patch [objectJson] [changesJson] requires 2 arguments"
var jobPatchBadPatchJSONString = "asdgasdg"
var jobPatchBadPatchJSONErrorString = "Error: Unable to parse drpcli jobs patch [objectJson] [changesJson] JSON asdgasdg\nError: error unmarshaling JSON: json: cannot unmarshal string into Go value of type models.Job\n\n"
var jobPatchBadBaseJSONString = "asdgasdg"
var jobPatchBadBaseJSONErrorString = "Error: Unable to parse drpcli jobs patch [objectJson] [changesJson] JSON asdgasdg\nError: error unmarshaling JSON: json: cannot unmarshal string into Go value of type models.Job\n\n"
var jobPatchBaseString string = `{
  "Addr": "192.168.100.100",
  "Strategy": "MAC",
  "Token": "john"
}
`
var jobPatchInputString string = `{
  "Options": [ { "Code": 3, "Value": "1.1.3.1" } ]
}
`
var jobPatchJohnString string = `{
  "Addr": "192.168.100.100",
  "NextServer": "2.2.2.2",
  "Options": [
    {
      "Code": 3,
      "Value": "1.1.3.1"
    }
  ],
  "Strategy": "MAC",
  "Token": "john"
}
`
var jobPatchMissingBaseString string = `{
  "Addr": "193.168.100.100",
  "NextServer": "2.2.2.2",
  "Strategy": "NewStrat",
  "Token": "john"
}
`
var jobPatchJohnMissingErrorString string = "Error: jobs: PATCH C1A86464: Not Found\n\n"

var jobDestroyNoArgErrorString string = "Error: drpcli jobs destroy [id] requires 1 argument"
var jobDestroyTooManyArgErrorString string = "Error: drpcli jobs destroy [id] requires 1 argument"
var jobDestroyJohnString string = "Deleted job 192.168.100.100\n"
var jobDestroyMissingJohnString string = "Error: jobs: DELETE C0A86464: Not Found\n\n"

func TestJobCli(t *testing.T) {
	tests := []CliTest{
		CliTest{true, false, []string{"jobs"}, noStdinString, "Access CLI commands relating to jobs\n", ""},
		CliTest{false, false, []string{"jobs", "list"}, noStdinString, jobDefaultListString, noErrorString},

		CliTest{true, true, []string{"jobs", "create"}, noStdinString, noContentString, jobCreateNoArgErrorString},
		CliTest{true, true, []string{"jobs", "create", "john", "john2"}, noStdinString, noContentString, jobCreateTooManyArgErrorString},
		CliTest{false, true, []string{"jobs", "create", jobCreateBadJSONString}, noStdinString, noContentString, jobCreateBadJSONErrorString},
		CliTest{false, false, []string{"jobs", "create", jobCreateInputString}, noStdinString, jobCreateJohnString, noErrorString},
		CliTest{false, true, []string{"jobs", "create", jobCreateInputString}, noStdinString, noContentString, jobCreateDuplicateErrorString},
		CliTest{false, false, []string{"jobs", "list"}, noStdinString, jobListBothEnvsString, noErrorString},

		CliTest{false, false, []string{"jobs", "list", "--limit=0"}, noStdinString, jobEmptyListString, noErrorString},
		CliTest{false, false, []string{"jobs", "list", "--limit=10", "--offset=0"}, noStdinString, jobListJobsString, noErrorString},
		CliTest{false, false, []string{"jobs", "list", "--limit=10", "--offset=10"}, noStdinString, jobEmptyListString, noErrorString},
		CliTest{false, true, []string{"jobs", "list", "--limit=-10", "--offset=0"}, noStdinString, noContentString, limitNegativeError},
		CliTest{false, true, []string{"jobs", "list", "--limit=10", "--offset=-10"}, noStdinString, noContentString, offsetNegativeError},
		CliTest{false, false, []string{"jobs", "list", "--limit=-1", "--offset=-1"}, noStdinString, jobListJobsString, noErrorString},
		CliTest{false, false, []string{"jobs", "list", "UUID=fred"}, noStdinString, jobEmptyListString, noErrorString},
		CliTest{false, false, []string{"jobs", "list", "UUID=MAC"}, noStdinString, jobListJobsString, noErrorString},
		CliTest{false, false, []string{"jobs", "list", "Key=john"}, noStdinString, jobListJobsString, noErrorString},
		CliTest{false, false, []string{"jobs", "list", "Key=false"}, noStdinString, jobEmptyListString, noErrorString},
		CliTest{false, false, []string{"jobs", "list", "BootEnv=192.168.100.100"}, noStdinString, jobListJobsString, noErrorString},
		CliTest{false, false, []string{"jobs", "list", "BootEnv=1.1.1.1"}, noStdinString, jobEmptyListString, noErrorString},
		CliTest{false, false, []string{"jobs", "list", "Task=3.3.3.3"}, noStdinString, jobEmptyListString, noErrorString},
		CliTest{false, false, []string{"jobs", "list", "Task=2.2.2.2"}, noStdinString, jobListJobsString, noErrorString},
		CliTest{false, false, []string{"jobs", "list", "State=3.3.3.3"}, noStdinString, jobEmptyListString, noErrorString},
		CliTest{false, false, []string{"jobs", "list", "State=2.2.2.2"}, noStdinString, jobListJobsString, noErrorString},
		CliTest{false, false, []string{"jobs", "list", "Machine=3.3.3.3"}, noStdinString, jobEmptyListString, noErrorString},
		CliTest{false, false, []string{"jobs", "list", "Machine=2.2.2.2"}, noStdinString, jobListJobsString, noErrorString},
		CliTest{false, false, []string{"jobs", "list", "Archived=3.3.3.3"}, noStdinString, jobEmptyListString, noErrorString},
		CliTest{false, false, []string{"jobs", "list", "Archived=2.2.2.2"}, noStdinString, jobListJobsString, noErrorString},
		CliTest{false, false, []string{"jobs", "list", "StartTime=3.3.3.3"}, noStdinString, jobEmptyListString, noErrorString},
		CliTest{false, false, []string{"jobs", "list", "StartTime=2.2.2.2"}, noStdinString, jobListJobsString, noErrorString},
		CliTest{false, false, []string{"jobs", "list", "EndTime=3.3.3.3"}, noStdinString, jobEmptyListString, noErrorString},
		CliTest{false, false, []string{"jobs", "list", "EndTime=2.2.2.2"}, noStdinString, jobListJobsString, noErrorString},

		CliTest{true, true, []string{"jobs", "show"}, noStdinString, noContentString, jobShowNoArgErrorString},
		CliTest{true, true, []string{"jobs", "show", "john", "john2"}, noStdinString, noContentString, jobShowTooManyArgErrorString},
		CliTest{false, true, []string{"jobs", "show", "192.168.100.103"}, noStdinString, noContentString, jobShowMissingArgErrorString},
		CliTest{false, false, []string{"jobs", "show", "192.168.100.100"}, noStdinString, jobShowJohnString, noErrorString},

		CliTest{true, true, []string{"jobs", "exists"}, noStdinString, noContentString, jobExistsNoArgErrorString},
		CliTest{true, true, []string{"jobs", "exists", "john", "john2"}, noStdinString, noContentString, jobExistsTooManyArgErrorString},
		CliTest{false, false, []string{"jobs", "exists", "192.168.100.100"}, noStdinString, jobExistsIgnoreString, noErrorString},
		CliTest{false, true, []string{"jobs", "exists", "ignore"}, noStdinString, noContentString, jobExistsMissingIgnoreString},
		CliTest{true, true, []string{"jobs", "exists", "john", "john2"}, noStdinString, noContentString, jobExistsTooManyArgErrorString},

		CliTest{true, true, []string{"jobs", "update"}, noStdinString, noContentString, jobUpdateNoArgErrorString},
		CliTest{true, true, []string{"jobs", "update", "john", "john2", "john3"}, noStdinString, noContentString, jobUpdateTooManyArgErrorString},
		CliTest{false, true, []string{"jobs", "update", "192.168.100.100", jobUpdateBadJSONString}, noStdinString, noContentString, jobUpdateBadJSONErrorString},
		CliTest{false, false, []string{"jobs", "update", "192.168.100.100", jobUpdateInputString}, noStdinString, jobUpdateJohnString, noErrorString},
		CliTest{false, true, []string{"jobs", "update", "192.168.100.103", jobUpdateInputString}, noStdinString, noContentString, jobUpdateJohnMissingErrorString},
		CliTest{false, false, []string{"jobs", "show", "192.168.100.100"}, noStdinString, jobUpdateJohnString, noErrorString},

		CliTest{true, true, []string{"jobs", "patch"}, noStdinString, noContentString, jobPatchNoArgErrorString},
		CliTest{true, true, []string{"jobs", "patch", "john", "john2", "john3"}, noStdinString, noContentString, jobPatchTooManyArgErrorString},
		CliTest{false, true, []string{"jobs", "patch", jobPatchBaseString, jobPatchBadPatchJSONString}, noStdinString, noContentString, jobPatchBadPatchJSONErrorString},
		CliTest{false, true, []string{"jobs", "patch", jobPatchBadBaseJSONString, jobPatchInputString}, noStdinString, noContentString, jobPatchBadBaseJSONErrorString},
		CliTest{false, false, []string{"jobs", "patch", jobPatchBaseString, jobPatchInputString}, noStdinString, jobPatchJohnString, noErrorString},
		CliTest{false, true, []string{"jobs", "patch", jobPatchMissingBaseString, jobPatchInputString}, noStdinString, noContentString, jobPatchJohnMissingErrorString},
		CliTest{false, false, []string{"jobs", "show", "192.168.100.100"}, noStdinString, jobPatchJohnString, noErrorString},

		CliTest{true, true, []string{"jobs", "destroy"}, noStdinString, noContentString, jobDestroyNoArgErrorString},
		CliTest{true, true, []string{"jobs", "destroy", "john", "june"}, noStdinString, noContentString, jobDestroyTooManyArgErrorString},
		CliTest{false, false, []string{"jobs", "destroy", "192.168.100.100"}, noStdinString, jobDestroyJohnString, noErrorString},
		CliTest{false, true, []string{"jobs", "destroy", "192.168.100.100"}, noStdinString, noContentString, jobDestroyMissingJohnString},
		CliTest{false, false, []string{"jobs", "list"}, noStdinString, jobDefaultListString, noErrorString},

		CliTest{false, false, []string{"jobs", "create", "-"}, jobCreateInputString + "\n", jobCreateJohnString, noErrorString},
		CliTest{false, false, []string{"jobs", "list"}, noStdinString, jobListBothEnvsString, noErrorString},
		CliTest{false, false, []string{"jobs", "update", "192.168.100.100", "-"}, jobUpdateInputString + "\n", jobUpdateJohnString, noErrorString},
		CliTest{false, false, []string{"jobs", "show", "192.168.100.100"}, noStdinString, jobUpdateJohnString, noErrorString},

		CliTest{false, false, []string{"jobs", "destroy", "192.168.100.100"}, noStdinString, jobDestroyJohnString, noErrorString},
		CliTest{false, false, []string{"jobs", "list"}, noStdinString, jobDefaultListString, noErrorString},
	}

	for _, test := range tests {
		testCli(t, test)
	}

}
