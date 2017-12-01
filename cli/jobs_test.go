package cli

import (
	"testing"
)

var jobEmptyListString string = "[]\n"
var jobDefaultListString string = "[]\n"

var jobTask1Create string = `{
  "Available": true,
  "Errors": [],
  "Meta": {
    "feature-flags": "original-exit-codes"
  },
  "Name": "task1",
  "OptionalParams": [],
  "ReadOnly": false,
  "RequiredParams": [],
  "Templates": [],
  "Validated": true
}
`
var jobTask2Create string = `{
  "Available": true,
  "Errors": [],
  "Meta": {
    "feature-flags": "original-exit-codes"
  },
  "Name": "task2",
  "OptionalParams": [],
  "ReadOnly": false,
  "RequiredParams": [],
  "Templates": [
    {
      "Contents": "Fred rules",
      "Name": "part 1",
      "Path": ""
    }
  ],
  "Validated": true
}
`
var jobTask3Create string = `{
  "Available": true,
  "Errors": [],
  "Meta": {
    "feature-flags": "original-exit-codes"
  },
  "Name": "task3",
  "OptionalParams": [],
  "ReadOnly": false,
  "RequiredParams": [],
  "Templates": [],
  "Validated": true
}
`

var jobLocal2Create string = `{
  "Available": true,
  "BootEnv": "local",
  "Errors": [],
  "Name": "stage3",
  "OptionalParams": [],
  "Profiles": [],
  "ReadOnly": false,
  "RequiredParams": [],
  "Tasks": [
    "task3",
    "task2",
    "task1"
  ],
  "Templates": [],
  "Validated": true
}
`
var jobLocal2CreateInput string = `{
  "Name": "stage3",
  "Tasks": [
    "task3",
    "task2",
    "task1"
  ],
  "BootEnv": "local"
}
`

var jobCreateMachineInputString string = `{
  "Address": "192.168.100.110",
  "Name": "john",
  "Secret": "secret1",
  "Uuid": "3e7031fe-3062-45f1-835c-92541bc9cbd3",
  "Stage": "stage3"
}
`
var jobCreateMachineJohnString string = `{
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
  "Stage": "stage3",
  "Tasks": [
    "task1",
    "task2",
    "task3"
  ],
  "Uuid": "3e7031fe-3062-45f1-835c-92541bc9cbd3",
  "Validated": true
}
`

var jobLocalUpdateInput string = `{
  "Tasks": ["task1","task2","task3"]
}
`
var jobLocalUpdateString string = `{
  "Available": true,
  "BootEnv": "local",
  "Errors": [],
  "Name": "stage3",
  "OptionalParams": [],
  "Profiles": [],
  "ReadOnly": false,
  "RequiredParams": [],
  "Tasks": [
    "task1",
    "task2",
    "task3"
  ],
  "Templates": [],
  "Validated": true
}
`

var jobCreateNoArgErrorString string = "Error: drpcli jobs create [json] [flags] requires 1 argument"
var jobCreateTooManyArgErrorString string = "Error: drpcli jobs create [json] [flags] requires 1 argument"
var jobCreateBadJSONErrorString string = "Error: Invalid job object: error converting YAML to JSON: yaml: line 1: did not find expected ',' or '}' and error converting YAML to JSON: yaml: line 1: did not find expected ',' or '}'\n\n"
var jobCreateBadJSON2ErrorString string = "Error: Unable to create new job: Invalid type passed to job create\n\n"

var jobCreateBadJSONString string = "{asdgasdg"
var jobCreateBadJSON2String string = "[asdgasdg]"
var jobCreateInputString string = `{
  "Uuid":    "00000000-0000-0000-0000-000000000001",
  "Machine": "3e7031fe-3062-45f1-835c-92541bc9cbd3"
}
`
var jobCreateNextInputString string = `{
  "Uuid":    "00000000-0000-0000-0000-000000000002",
  "Machine": "3e7031fe-3062-45f1-835c-92541bc9cbd3"
}
`
var jobCreateNextInput3String string = `{
  "Uuid":    "00000000-0000-0000-0000-000000000003",
  "Machine": "3e7031fe-3062-45f1-835c-92541bc9cbd3"
}
`
var jobCreateNextInput4String string = `{
  "Uuid":    "00000000-0000-0000-0000-000000000004",
  "Machine": "3e7031fe-3062-45f1-835c-92541bc9cbd3"
}
`
var jobCreateNextInput5String string = `{
  "Uuid":    "00000000-0000-0000-0000-000000000005",
  "Machine": "3e7031fe-3062-45f1-835c-92541bc9cbd3"
}
`
var jobCreateJohnString string = `RE:
{
  "Archived": false,
  "Available": true,
  "Current": true,
  "EndTime": "0001-01-01T00:00:00Z",
  "Errors": [],
  "LogPath": "[\S\s]+/job-logs/00000000-0000-0000-0000-000000000001",
  "Machine": "3e7031fe-3062-45f1-835c-92541bc9cbd3",
  "Previous": "00000000-0000-0000-0000-000000000000",
  "ReadOnly": false,
  "Stage": "stage3",
  "StartTime": "0001-01-01T00:00:00Z",
  "State": "created",
  "Task": "task1",
  "Uuid": "00000000-0000-0000-0000-000000000001",
  "Validated": true
}
`
var jobListJobsString string = `RE:
\[
  {
    "Archived": false,
    "Available": true,
    "Current": true,
    "EndTime": "0001-01-01T00:00:00Z",
    "Errors": [],
    "LogPath": "[\S\s]*/job-logs/00000000-0000-0000-0000-000000000001",
    "Machine": "3e7031fe-3062-45f1-835c-92541bc9cbd3",
    "Previous": "00000000-0000-0000-0000-000000000000",
    "ReadOnly": false,
    "Stage": "stage3",
    "StartTime": "0001-01-01T00:00:00Z",
    "State": "created",
    "Task": "task1",
    "Uuid": "00000000-0000-0000-0000-000000000001",
    "Validated": true
  }
\]
`

var jobCreateJobAlreadyRunningErrorString string = "Error: Conflict: Machine 3e7031fe-3062-45f1-835c-92541bc9cbd3 already has running or created job\n\n"

var jobShowNoArgErrorString string = "Error: drpcli jobs show [id] [flags] requires 1 argument"
var jobShowTooManyArgErrorString string = "Error: drpcli jobs show [id] [flags] requires 1 argument"
var jobShowMissingArgErrorString string = "Error: GET: jobs/john: Not Found\n\n"
var jobShowJobString string = `RE:
{
  "Archived": false,
  "Available": true,
  "Current": true,
  "EndTime": "0001-01-01T00:00:00Z",
  "Errors": [],
  "LogPath": "[\S\s]+/job-logs/00000000-0000-0000-0000-000000000001",
  "Machine": "3e7031fe-3062-45f1-835c-92541bc9cbd3",
  "Previous": "00000000-0000-0000-0000-000000000000",
  "ReadOnly": false,
  "Stage": "stage3",
  "StartTime": "0001-01-01T00:00:00Z",
  "State": "created",
  "Task": "task1",
  "Uuid": "00000000-0000-0000-0000-000000000001",
  "Validated": true
}
`
var jobExistsNoArgErrorString string = "Error: drpcli jobs exists [id] [flags] requires 1 argument"
var jobExistsTooManyArgErrorString string = "Error: drpcli jobs exists [id] [flags] requires 1 argument"
var jobExistsJobString string = ""
var jobExistsMissingJohnString string = "Error: GET: jobs/john: Not Found\n\n"

var jobExpireTimeErrorString string = "Error: GET: jobs: Invalid UUID: false\n\n"
var jobDestroyBadString string = "Error: ValidationError: jobs/00000000-0000-0000-0000-000000000001: Jobs 00000000-0000-0000-0000-000000000001 is not in a deletable state: created\n\n"
var jobBadTimeFormatString string = `Error: GET: jobs: parsing time "fred" as "2006-01-02T15:04:05Z07:00": cannot parse "fred" as "2006"

`
var jobCreateJobInvalidMachineNameErrorString string = "Error: Unable to create new job: Invalid machine name passed to job create: james\n\n"

var jobUpdateNoArgErrorString string = "Error: drpcli jobs update [id] [json] [flags] requires 2 arguments\n"
var jobUpdateTooManyArgErrorString string = "Error: drpcli jobs update [id] [json] [flags] requires 2 arguments\n"

var jobShowMachineJohnString string = `{
  "Address": "192.168.100.110",
  "Available": true,
  "BootEnv": "local",
  "CurrentJob": "00000000-0000-0000-0000-000000000001",
  "CurrentTask": 0,
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
  "Stage": "stage3",
  "Tasks": [
    "task1",
    "task2",
    "task3"
  ],
  "Uuid": "3e7031fe-3062-45f1-835c-92541bc9cbd3",
  "Validated": true
}
`

var jobUpdateBadJSONString string = "{asgasdg"
var jobUpdateBadJSONErrorString string = "Error: Unable to unmarshal input stream: error converting YAML to JSON: yaml: line 1: did not find expected ',' or '}'\n\n\n"
var jobUpdateBadJSON2String string = "[ \"gadsg\" ]"
var jobUpdateBadJSON2ErrorString string = "Error: Unable to merge objects: json: cannot unmarshal array into Go value of type map[string]interface {}\n\n\n"
var jobUpdateInputString string = "{ \"State\": \"incomplete\" }"
var jobUpdateBadInputString string = "{ \"State\": \"fred\" }"
var jobUpdateBadInputErrorString string = "Error: ValidationError: jobs/00000000-0000-0000-0000-000000000001: State fred is not valid\n\n"
var jobUpdateJohnString string = `RE:
{
  "Archived": false,
  "Available": true,
  "Current": true,
  "EndTime": "0001-01-01T00:00:00Z",
  "Errors": [],
  "LogPath": "[\S\s]+/job-logs/00000000-0000-0000-0000-000000000001",
  "Machine": "3e7031fe-3062-45f1-835c-92541bc9cbd3",
  "Previous": "00000000-0000-0000-0000-000000000000",
  "ReadOnly": false,
  "Stage": "stage3",
  "StartTime": "0001-01-01T00:00:00Z",
  "State": "incomplete",
  "Task": "task1",
  "Uuid": "00000000-0000-0000-0000-000000000001",
  "Validated": true
}
`

var jobUpdateJohnMissingErrorString string = "Error: GET: jobs/john2: Not Found\n\n"

var jobPatchNoArgErrorString string = "Error: drpcli jobs patch [objectJson] [changesJson] [flags] requires 2 arguments\n"
var jobPatchTooManyArgErrorString = "Error: drpcli jobs patch [objectJson] [changesJson] [flags] requires 2 arguments\n"
var jobPatchBaseString = `{
  "Archived": false,
  "Available": true,
  "Errors": [],
  "Machine": "3e7031fe-3062-45f1-835c-92541bc9cbd3",
  "Previous": "00000000-0000-0000-0000-000000000000",
  "Stage": "stage3",
  "State": "incomplete",
  "Task": "task1",
  "Uuid": "00000000-0000-0000-0000-000000000001",
  "Validated": true
}
`
var jobPatchBase2String = `{
  "Archived": false,
  "Available": true,
  "Errors": [],
  "Machine": "3e7031fe-3062-45f1-835c-92541bc9cbd3",
  "Previous": "00000000-0000-0000-0000-000000000000",
  "Stage": "stage3",
  "State": "running",
  "Task": "task1",
  "Uuid": "00000000-0000-0000-0000-000000000001",
  "Validated": true
}
`
var jobPatchBadPatchJSONString = "{asdgasdg"
var jobPatchBadPatchJSONErrorString = "Error: Unable to parse drpcli jobs patch [objectJson] [changesJson] [flags] JSON {asdgasdg\nError: error converting YAML to JSON: yaml: line 1: did not find expected ',' or '}'\n\n"

var jobPatchBadPatchJSON2String = "[ \"asdgasdg\" ]"
var jobPatchBadPatchJSON2ErrorString = "Error: Unable to parse drpcli jobs patch [objectJson] [changesJson] [flags] JSON [ \"asdgasdg\" ]\nError: error unmarshaling JSON: json: cannot unmarshal array into Go value of type genmodels.Job\n\n"

var jobPatchBadBaseJSONString = "{ badbase"
var jobPatchBadBaseJSONErrorString = "Error: Unable to parse drpcli jobs patch [objectJson] [changesJson] [flags] JSON { badbase\nError: error converting YAML to JSON: yaml: line 1: did not find expected ',' or '}'\n\n"

var jobPatchBadInputString = "{ \"State\": \"fred\"}"
var jobPatchBadInputErrorString = "Error: ValidationError: jobs/00000000-0000-0000-0000-000000000001: State fred is not valid\n\n"
var jobPatchInputString = "{ \"State\": \"running\"}"
var jobPatchInputReplyString = `RE:
{
  "Archived": false,
  "Available": true,
  "Current": true,
  "EndTime": "0001-01-01T00:00:00Z",
  "Errors": [],
  "LogPath": "[\S\s]+/job-logs/00000000-0000-0000-0000-000000000001",
  "Machine": "3e7031fe-3062-45f1-835c-92541bc9cbd3",
  "Previous": "00000000-0000-0000-0000-000000000000",
  "ReadOnly": false,
  "Stage": "stage3",
  "StartTime": "20[\s\S]*",
  "State": "running",
  "Task": "task1",
  "Uuid": "00000000-0000-0000-0000-000000000001",
  "Validated": true
}
`
var jobPatchInput2String = "{ \"State\": \"incomplete\"}"
var jobPatchJohnString = `RE:
{
  "Archived": false,
  "Available": true,
  "Current": true,
  "EndTime": "0001-01-01T00:00:00Z",
  "Errors": [],
  "LogPath": "[\S\s]+/job-logs/00000000-0000-0000-0000-000000000001",
  "Machine": "3e7031fe-3062-45f1-835c-92541bc9cbd3",
  "Previous": "00000000-0000-0000-0000-000000000000",
  "ReadOnly": false,
  "Stage": "stage3",
  "StartTime": "20[\s\S]*",
  "State": "incomplete",
  "Task": "task1",
  "Uuid": "00000000-0000-0000-0000-000000000001",
  "Validated": true
}
`
var jobPatchMissingBaseString = `{
  "Archived": false,
  "Available": true,
  "Errors": [],
  "Machine": "3e7031fe-3062-45f1-835c-92541bc9cbd3",
  "Previous": "00000000-0000-0000-0000-000000000000",
  "Stage": "stage3",
  "State": "incomplete",
  "Task": "task1",
  "Uuid": "10000000-0000-0000-0000-000000000001",
  "Validated": true
}
`
var jobPatchJohnMissingErrorString = "Error: PATCH: jobs/10000000-0000-0000-0000-000000000001: Not Found\n\n"

var jobUpdateFailedJobInputString = "{ \"State\": \"failed\" }"
var jobUpdateFailedJobUpdateString = `RE:
{
  "Archived": false,
  "Available": true,
  "Current": true,
  "EndTime": "20[\s\S]*",
  "Errors": [],
  "LogPath": "[\S\s]+/job-logs/00000000-0000-0000-0000-000000000001",
  "Machine": "3e7031fe-3062-45f1-835c-92541bc9cbd3",
  "Previous": "00000000-0000-0000-0000-000000000000",
  "ReadOnly": false,
  "Stage": "stage3",
  "StartTime": "20[\s\S]*",
  "State": "failed",
  "Task": "task1",
  "Uuid": "00000000-0000-0000-0000-000000000001",
  "Validated": true
}
`
var jobCreateMachineNotRunningErrorString = "Error: Conflict: Machine 3e7031fe-3062-45f1-835c-92541bc9cbd3 is not runnable\n\n"
var jobUpdateMachineRunnableString = `{
  "Address": "192.168.100.110",
  "Available": true,
  "BootEnv": "local",
  "CurrentJob": "00000000-0000-0000-0000-000000000001",
  "CurrentTask": 0,
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
  "Stage": "stage3",
  "Tasks": [
    "task1",
    "task2",
    "task3"
  ],
  "Uuid": "3e7031fe-3062-45f1-835c-92541bc9cbd3",
  "Validated": true
}
`
var jobCreateNextString = `RE:
{
  "Archived": false,
  "Available": true,
  "Current": true,
  "EndTime": "0001-01-01T00:00:00Z",
  "Errors": [],
  "LogPath": "[\S\s]+/job-logs/00000000-0000-0000-0000-000000000002",
  "Machine": "3e7031fe-3062-45f1-835c-92541bc9cbd3",
  "Previous": "00000000-0000-0000-0000-000000000001",
  "ReadOnly": false,
  "Stage": "stage3",
  "StartTime": "0001-01-01T00:00:00Z",
  "State": "created",
  "Task": "task1",
  "Uuid": "00000000-0000-0000-0000-000000000002",
  "Validated": true
}
`
var jobUpdateFinishedJobInputString = "{ \"State\": \"finished\" }"
var jobUpdateFinishedJob2UpdateString = `RE:
{
  "Archived": false,
  "Available": true,
  "Current": true,
  "EndTime": "20[\s\S]*",
  "Errors": [],
  "LogPath": "[\S\s]+/job-logs/00000000-0000-0000-0000-000000000002",
  "Machine": "3e7031fe-3062-45f1-835c-92541bc9cbd3",
  "Previous": "00000000-0000-0000-0000-000000000001",
  "ReadOnly": false,
  "Stage": "stage3",
  "StartTime": "0001-01-01T00:00:00Z",
  "State": "finished",
  "Task": "task1",
  "Uuid": "00000000-0000-0000-0000-000000000002",
  "Validated": true
}
`
var jobCreateNext3String = `RE:
{
  "Archived": false,
  "Available": true,
  "Current": true,
  "EndTime": "0001-01-01T00:00:00Z",
  "Errors": [],
  "LogPath": "[\S\s]+/job-logs/00000000-0000-0000-0000-000000000003",
  "Machine": "3e7031fe-3062-45f1-835c-92541bc9cbd3",
  "Previous": "00000000-0000-0000-0000-000000000002",
  "ReadOnly": false,
  "Stage": "stage3",
  "StartTime": "0001-01-01T00:00:00Z",
  "State": "created",
  "Task": "task2",
  "Uuid": "00000000-0000-0000-0000-000000000003",
  "Validated": true
}
`
var jobUpdateFinishedJob3UpdateString = `RE:
{
  "Archived": false,
  "Available": true,
  "Current": true,
  "EndTime": "20[\s\S]*",
  "Errors": [],
  "LogPath": "[\S\s]+/job-logs/00000000-0000-0000-0000-000000000003",
  "Machine": "3e7031fe-3062-45f1-835c-92541bc9cbd3",
  "Previous": "00000000-0000-0000-0000-000000000002",
  "ReadOnly": false,
  "Stage": "stage3",
  "StartTime": "0001-01-01T00:00:00Z",
  "State": "finished",
  "Task": "task2",
  "Uuid": "00000000-0000-0000-0000-000000000003",
  "Validated": true
}
`
var jobCreateNext4String = `RE:
{
  "Archived": false,
  "Available": true,
  "Current": true,
  "EndTime": "0001-01-01T00:00:00Z",
  "Errors": [],
  "LogPath": "[\S\s]+/job-logs/00000000-0000-0000-0000-000000000004",
  "Machine": "3e7031fe-3062-45f1-835c-92541bc9cbd3",
  "Previous": "00000000-0000-0000-0000-000000000003",
  "ReadOnly": false,
  "Stage": "stage3",
  "StartTime": "0001-01-01T00:00:00Z",
  "State": "created",
  "Task": "task3",
  "Uuid": "00000000-0000-0000-0000-000000000004",
  "Validated": true
}
`
var jobUpdateFinishedJob4UpdateString = `RE:
{
  "Archived": false,
  "Available": true,
  "Current": true,
  "EndTime": "20[\s\S]*",
  "Errors": [],
  "LogPath": "[\S\s]+/job-logs/00000000-0000-0000-0000-000000000004",
  "Machine": "3e7031fe-3062-45f1-835c-92541bc9cbd3",
  "Previous": "00000000-0000-0000-0000-000000000003",
  "ReadOnly": false,
  "Stage": "stage3",
  "StartTime": "0001-01-01T00:00:00Z",
  "State": "finished",
  "Task": "task3",
  "Uuid": "00000000-0000-0000-0000-000000000004",
  "Validated": true
}
`
var jobFullListString = `RE:
[
  {
    "Archived": false,
    "Available": true,
    "Current": false,
    "EndTime": "20[\s\S]*",
    "Errors": [],
    "LogPath": "[\S\s]+/job-logs/00000000-0000-0000-0000-000000000001",
    "Machine": "3e7031fe-3062-45f1-835c-92541bc9cbd3",
    "Previous": "00000000-0000-0000-0000-000000000000",
    "ReadOnly": false,
    "Stage": "stage3",
    "StartTime": "20[\s\S]*",
    "State": "failed",
    "Task": "task1",
    "Uuid": "00000000-0000-0000-0000-000000000001",
    "Validated": true
  },
  {
    "Archived": false,
    "Available": true,
    "Current": false,
    "EndTime": "20[\s\S]*",
    "Errors": [],
    "LogPath": "[\S\s]+/job-logs/00000000-0000-0000-0000-000000000002",
    "Machine": "3e7031fe-3062-45f1-835c-92541bc9cbd3",
    "Previous": "00000000-0000-0000-0000-000000000001",
    "ReadOnly": false,
    "Stage": "stage3",
    "StartTime": "0001-01-01T00:00:00Z",
    "State": "finished",
    "Task": "task1",
    "Uuid": "00000000-0000-0000-0000-000000000002",
    "Validated": true
  },
  {
    "Archived": false,
    "Available": true,
    "Current": false,
    "EndTime": "20[\s\S]*",
    "Errors": [],
    "LogPath": "[\S\s]+/job-logs/00000000-0000-0000-0000-000000000003",
    "Machine": "3e7031fe-3062-45f1-835c-92541bc9cbd3",
    "Previous": "00000000-0000-0000-0000-000000000002",
    "ReadOnly": false,
    "Stage": "stage3",
    "StartTime": "0001-01-01T00:00:00Z",
    "State": "finished",
    "Task": "task2",
    "Uuid": "00000000-0000-0000-0000-000000000003",
    "Validated": true
  },
  {
    "Archived": false,
    "Available": true,
    "Current": true,
    "EndTime": "20[\s\S]*",
    "Errors": [],
    "LogPath": "[\S\s]+/job-logs/00000000-0000-0000-0000-000000000004",
    "Machine": "3e7031fe-3062-45f1-835c-92541bc9cbd3",
    "Previous": "00000000-0000-0000-0000-000000000003",
    "ReadOnly": false,
    "Stage": "stage3",
    "StartTime": "0001-01-01T00:00:00Z",
    "State": "finished",
    "Task": "task3",
    "Uuid": "00000000-0000-0000-0000-000000000004",
    "Validated": true
  }
]
`
var jobDestroyNoArgErrorString = "Error: drpcli jobs destroy [id] [flags] requires 1 argument\n"
var jobDestroyTooManyArgErrorString = "Error: drpcli jobs destroy [id] [flags] requires 1 argument\n"
var jobDestroyMissingJohnString = "Error: DELETE: jobs/3e7031fe-3062-45f1-835c-92541bc9cbd3: Not Found\n\n"
var jobDestroy001String = "Deleted job 00000000-0000-0000-0000-000000000001\n"
var jobDestroy002String = "Deleted job 00000000-0000-0000-0000-000000000002\n"
var jobDestroy003String = "Deleted job 00000000-0000-0000-0000-000000000003\n"
var jobDestroy004String = "Deleted job 00000000-0000-0000-0000-000000000004\n"

var jobActionsNoArgErrorString = "Error: drpcli jobs actions [id] [flags] requires 1 argument\n"
var jobActionsTooManyArgErrorString = "Error: drpcli jobs actions [id] [flags] requires 1 argument\n"
var jobActionsMissingJobErrorString = "Error: ValidationError: Job john does not exist\n\n"
var jobActionsRenderedTask1String = "[]\n"
var jobActionsRenderedTask2String = `[
  {
    "Content": "Fred rules",
    "Name": "part 1",
    "Path": ""
  }
]
`
var jobActionsMissingMachineRenderErrorString = "Error: ValidationError: Machine 3e7031fe-3062-45f1-835c-92541bc9cbd3 does not exist\n\n"
var jobActionsMissingTaskRenderErrorString = "Error: ValidationError: Task task2 does not exist\n\n"

var jobLogNoArgErrorString = "Error: drpcli jobs log [id] [- or string] [flags] requires at least 1 argument\n"
var jobLogTooManyArgsErrorString = "Error: drpcli jobs log [id] [- or string] [flags] requires at most 2 arguments\n"
var jobLogUnknownJobErrorString = "Error: ValidationError: Job john does not exist\n\n"

func TestJobCli(t *testing.T) {

	cliTest(false, false, "tasks", "create", "task1").run(t)
	cliTest(false, false, "tasks", "create", "-").Stdin(jobTask2Create).run(t)
	cliTest(false, false, "tasks", "create", "task3").run(t)

	cliTest(false, false, "stages", "create", jobLocal2CreateInput).run(t)
	cliTest(false, false, "stages", "update", "stage3", jobLocalUpdateInput).run(t)

	cliTest(false, false, "machines", "create", jobCreateMachineInputString).run(t)

	cliTest(true, false, "jobs").run(t)
	cliTest(false, false, "jobs", "list").run(t)

	cliTest(true, true, "jobs", "create").run(t)
	cliTest(true, true, "jobs", "create", "john", "john2").run(t)
	cliTest(false, true, "jobs", "create", jobCreateBadJSONString).run(t)
	cliTest(false, true, "jobs", "create", jobCreateBadJSON2String).run(t)
	cliTest(false, false, "jobs", "create", jobCreateInputString).run(t)
	cliTest(false, true, "jobs", "create", jobCreateInputString).run(t)

	cliTest(false, true, "jobs", "create", "3e7031fe-3062-45f1-835c-92541bc9cbd3").run(t)
	cliTest(false, true, "jobs", "create", "john").run(t)
	cliTest(false, true, "jobs", "create", "james").run(t)

	cliTest(true, true, "jobs", "show").run(t)
	cliTest(true, true, "jobs", "show", "john", "john2").run(t)
	cliTest(false, true, "jobs", "show", "john").run(t)
	cliTest(false, false, "jobs", "show", "00000000-0000-0000-0000-000000000001").run(t)
	cliTest(false, false, "jobs", "show", "Key:00000000-0000-0000-0000-000000000001").run(t)
	cliTest(false, false, "jobs", "show", "Uuid:00000000-0000-0000-0000-000000000001").run(t)

	cliTest(true, true, "jobs", "exists").run(t)
	cliTest(true, true, "jobs", "exists", "john", "john2").run(t)
	cliTest(false, false, "jobs", "exists", "00000000-0000-0000-0000-000000000001").run(t)
	cliTest(false, true, "jobs", "exists", "john").run(t)
	cliTest(false, false, "jobs", "exists", "Uuid:00000000-0000-0000-0000-000000000001").run(t)
	cliTest(false, false, "machines", "show", "3e7031fe-3062-45f1-835c-92541bc9cbd3").run(t)

	cliTest(false, false, "jobs", "list").run(t)
	cliTest(false, false, "jobs", "list", "Stage=stage3").run(t)
	cliTest(false, false, "jobs", "list", "Stage=false").run(t)
	cliTest(false, false, "jobs", "list", "Task=task1").run(t)
	cliTest(false, false, "jobs", "list", "Task=false").run(t)
	cliTest(false, false, "jobs", "list", "State=created").run(t)
	cliTest(false, false, "jobs", "list", "State=false").run(t)
	cliTest(false, false, "jobs", "list", "Machine=3e7031fe-3062-45f1-835c-92541bc9cbd3").run(t)
	cliTest(false, false, "jobs", "list", "Machine=4e7031fe-3062-45f1-835c-92541bc9cbd3").run(t)
	cliTest(false, true, "jobs", "list", "Machine=false").run(t)
	cliTest(false, false, "jobs", "list", "Archived=false").run(t)
	cliTest(false, false, "jobs", "list", "Archived=true").run(t)
	cliTest(false, false, "jobs", "list", "StartTime=0001-01-01T00:00:00Z").run(t)
	cliTest(false, false, "jobs", "list", "StartTime=2001-01-01T00:00:00Z").run(t)
	cliTest(false, true, "jobs", "list", "StartTime=fred").run(t)
	cliTest(false, false, "jobs", "list", "EndTime=0001-01-01T00:00:00Z").run(t)
	cliTest(false, false, "jobs", "list", "EndTime=2001-01-01T00:00:00Z").run(t)
	cliTest(false, true, "jobs", "list", "EndTime=fred").run(t)
	cliTest(false, false, "jobs", "list", "Uuid=4e7031fe-3062-45f1-835c-92541bc9cbd3").run(t)
	cliTest(false, false, "jobs", "list", "Uuid=00000000-0000-0000-0000-000000000001").run(t)
	cliTest(false, true, "jobs", "list", "Uuid=false").run(t)
	cliTest(false, true, "jobs", "destroy", "00000000-0000-0000-0000-000000000001").run(t)

	cliTest(true, true, "jobs", "log").run(t)
	cliTest(true, true, "jobs", "log", "john", "john2", "john3").run(t)
	cliTest(false, true, "jobs", "log", "john").run(t)
	cliTest(false, false, "jobs", "log", "00000000-0000-0000-0000-000000000001").run(t)
	cliTest(false, false, "jobs", "log", "00000000-0000-0000-0000-000000000001", "Fred\n").run(t)
	cliTest(false, false, "jobs", "log", "00000000-0000-0000-0000-000000000001").run(t)
	cliTest(false, false, "jobs", "log", "00000000-0000-0000-0000-000000000001", "-").Stdin("Freddy\n").run(t)
	cliTest(false, false, "jobs", "log", "00000000-0000-0000-0000-000000000001").run(t)

	cliTest(true, true, "jobs", "update").run(t)
	cliTest(true, true, "jobs", "update", "john", "john2", "john3").run(t)
	cliTest(false, true, "jobs", "update", "00000000-0000-0000-0000-000000000001", jobUpdateBadJSONString).run(t)
	cliTest(false, true, "jobs", "update", "00000000-0000-0000-0000-000000000001", jobUpdateBadJSON2String).run(t)
	cliTest(false, true, "jobs", "update", "00000000-0000-0000-0000-000000000001", jobUpdateBadInputString).run(t)
	cliTest(false, false, "jobs", "update", "00000000-0000-0000-0000-000000000001", jobUpdateInputString).run(t)
	cliTest(false, true, "jobs", "update", "john2", jobUpdateInputString).run(t)
	cliTest(false, false, "jobs", "show", "00000000-0000-0000-0000-000000000001").run(t)
	// This tests that incomplete jobs come back.
	cliTest(false, false, "jobs", "create", "john").run(t)
	/* No patch tests for now
	cliTest(true, true, "jobs", "patch").run(t)
	cliTest(true, true, "jobs", "patch", "john", "john2", "john3").run(t)
	cliTest(false, true, "jobs", "patch", jobPatchBaseString, jobPatchBadPatchJSONString).run(t)
	cliTest(false, true, "jobs", "patch", jobPatchBaseString, jobPatchBadPatchJSON2String).run(t)
	cliTest(false, true, "jobs", "patch", jobPatchBadBaseJSONString, jobPatchInputString).run(t)
	cliTest(false, true, "jobs", "patch", jobPatchBaseString, jobPatchBadInputString).run(t)
	cliTest(false, false, "jobs", "patch", jobPatchBaseString, jobPatchInputString).run(t)
	cliTest(false, false, "jobs", "patch", jobPatchBase2String, jobPatchInput2String).run(t)
	cliTest(false, true, "jobs", "patch", jobPatchMissingBaseString, jobPatchInputString).run(t)
	*/
	cliTest(false, false, "jobs", "show", "00000000-0000-0000-0000-000000000001").run(t)

	// This tests that incomplet jobs come back again
	cliTest(false, false, "jobs", "create", "-").Stdin(jobCreateNextInputString + "\n").run(t)
	cliTest(false, false, "jobs", "update", "00000000-0000-0000-0000-000000000001", jobUpdateFailedJobInputString).run(t)
	cliTest(false, true, "jobs", "create", "-").Stdin(jobCreateNextInputString + "\n").run(t)
	cliTest(false, false, "machines", "update", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "{ \"Runnable\": true }").run(t)
	cliTest(false, false, "jobs", "create", "-").Stdin(jobCreateNextInputString + "\n").run(t)
	cliTest(false, false, "jobs", "update", "00000000-0000-0000-0000-000000000002", jobUpdateFinishedJobInputString).run(t)
	cliTest(false, false, "jobs", "create", "-").Stdin(jobCreateNextInput3String).run(t)
	cliTest(false, false, "jobs", "update", "00000000-0000-0000-0000-000000000003", jobUpdateFinishedJobInputString).run(t)
	cliTest(false, false, "jobs", "create", "-").Stdin(jobCreateNextInput4String + "\n").run(t)
	cliTest(false, false, "jobs", "update", "00000000-0000-0000-0000-000000000004", jobUpdateFinishedJobInputString).run(t)
	cliTest(false, false, "jobs", "create", "-").Stdin(jobCreateNextInput5String + "\n").run(t)

	cliTest(false, false, "jobs", "list").run(t)

	cliTest(true, true, "jobs", "actions").run(t)
	cliTest(true, true, "jobs", "actions", "john", "june").run(t)
	cliTest(false, true, "jobs", "actions", "john").run(t)
	cliTest(false, false, "jobs", "actions", "00000000-0000-0000-0000-000000000001").run(t)
	cliTest(false, false, "jobs", "actions", "00000000-0000-0000-0000-000000000003").run(t)
	cliTest(false, false, "machines", "destroy", "3e7031fe-3062-45f1-835c-92541bc9cbd3").run(t)
	cliTest(false, true, "jobs", "actions", "00000000-0000-0000-0000-000000000003").run(t)
	cliTest(false, false, "stages", "destroy", "stage3").run(t)
	cliTest(false, false, "tasks", "destroy", "task1").run(t)
	cliTest(false, false, "tasks", "destroy", "task2").run(t)
	cliTest(false, false, "tasks", "destroy", "task3").run(t)
	cliTest(false, true, "jobs", "actions", "00000000-0000-0000-0000-000000000003").run(t)

	cliTest(true, true, "jobs", "destroy").run(t)
	cliTest(true, true, "jobs", "destroy", "john", "june").run(t)
	cliTest(false, true, "jobs", "destroy", "3e7031fe-3062-45f1-835c-92541bc9cbd3").run(t)
	cliTest(false, false, "jobs", "destroy", "00000000-0000-0000-0000-000000000001").run(t)
	cliTest(false, false, "jobs", "destroy", "00000000-0000-0000-0000-000000000002").run(t)
	cliTest(false, false, "jobs", "destroy", "00000000-0000-0000-0000-000000000003").run(t)
	cliTest(false, false, "jobs", "destroy", "00000000-0000-0000-0000-000000000004").run(t)
	cliTest(false, false, "jobs", "list").run(t)

}
