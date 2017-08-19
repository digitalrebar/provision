package cli

import (
	"os"
	"testing"
)

var jobEmptyListString string = "[]\n"
var jobDefaultListString string = "[]\n"

var jobTask1Create string = `{
  "Available": true,
  "Errors": [],
  "Name": "task1",
  "OptionalParams": null,
  "RequiredParams": null,
  "Templates": null,
  "Validated": true
}
`
var jobTask2Create string = `{
  "Available": true,
  "Errors": [],
  "Name": "task2",
  "OptionalParams": null,
  "RequiredParams": null,
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
  "Name": "task3",
  "OptionalParams": null,
  "RequiredParams": null,
  "Templates": null,
  "Validated": true
}
`

var jobLocal2Create string = `{
  "Available": true,
  "BootParams": "",
  "Errors": [],
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
    "task3",
    "task2",
    "task1"
  ],
  "Templates": [
    {
      "ID": "local3-pxelinux.tmpl",
      "Name": "pxelinux",
      "Path": "pxelinux.cfg/{{.Machine.HexAddress}}"
    },
    {
      "ID": "local3-elilo.tmpl",
      "Name": "elilo",
      "Path": "{{.Machine.HexAddress}}.conf"
    },
    {
      "ID": "local3-ipxe.tmpl",
      "Name": "ipxe",
      "Path": "{{.Machine.Address}}.ipxe"
    }
  ],
  "Validated": true
}
`
var jobLocal2CreateInput string = `{
  "BootParams": "",
  "Errors": [],
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
    "task3",
    "task2",
    "task1"
  ],
  "Templates": [
    {
      "ID": "local3-pxelinux.tmpl",
      "Name": "pxelinux",
      "Path": "pxelinux.cfg/{{.Machine.HexAddress}}"
    },
    {
      "ID": "local3-elilo.tmpl",
      "Name": "elilo",
      "Path": "{{.Machine.HexAddress}}.conf"
    },
    {
      "ID": "local3-ipxe.tmpl",
      "Name": "ipxe",
      "Path": "{{.Machine.Address}}.ipxe"
    }
  ]
}
`

var jobCreateMachineInputString string = `{
  "Address": "192.168.100.110",
  "name": "john",
  "Uuid": "3e7031fe-3062-45f1-835c-92541bc9cbd3",
  "bootenv": "local3"
}
`
var jobCreateMachineJohnString string = `{
  "Address": "192.168.100.110",
  "Available": true,
  "BootEnv": "local3",
  "CurrentTask": -1,
  "Errors": [],
  "Name": "john",
  "Profile": {
    "Available": false,
    "Errors": null,
    "Name": "",
    "Tasks": null,
    "Validated": false
  },
  "Profiles": null,
  "Runnable": true,
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
  "BootParams": "",
  "Errors": [],
  "Initrds": null,
  "Kernel": "",
  "Name": "local3",
  "OS": {
    "Name": "local3"
  },
  "OnlyUnknown": false,
  "OptionalParams": null,
  "RequiredParams": null,
  "Tasks": [
    "task1",
    "task2",
    "task3"
  ],
  "Templates": [
    {
      "ID": "local3-pxelinux.tmpl",
      "Name": "pxelinux",
      "Path": "pxelinux.cfg/{{.Machine.HexAddress}}"
    },
    {
      "ID": "local3-elilo.tmpl",
      "Name": "elilo",
      "Path": "{{.Machine.HexAddress}}.conf"
    },
    {
      "ID": "local3-ipxe.tmpl",
      "Name": "ipxe",
      "Path": "{{.Machine.Address}}.ipxe"
    }
  ],
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
  "BootEnv": "local3",
  "EndTime": "0001-01-01T00:00:00Z",
  "Errors": [],
  "LogPath": "[\S\s]+/job-logs/00000000-0000-0000-0000-000000000001",
  "Machine": "3e7031fe-3062-45f1-835c-92541bc9cbd3",
  "Previous": "00000000-0000-0000-0000-000000000000",
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
    "BootEnv": "local3",
    "EndTime": "0001-01-01T00:00:00Z",
    "Errors": [],
    "LogPath": "[\S\s]*/job-logs/00000000-0000-0000-0000-000000000001",
    "Machine": "3e7031fe-3062-45f1-835c-92541bc9cbd3",
    "Previous": "00000000-0000-0000-0000-000000000000",
    "StartTime": "0001-01-01T00:00:00Z",
    "State": "created",
    "Task": "task1",
    "Uuid": "00000000-0000-0000-0000-000000000001",
    "Validated": true
  }
\]
`

var jobCreateJobAlreadyRunningErrorString string = "Error: Machine 3e7031fe-3062-45f1-835c-92541bc9cbd3 already has running or created job\n\n"

var jobShowNoArgErrorString string = "Error: drpcli jobs show [id] [flags] requires 1 argument"
var jobShowTooManyArgErrorString string = "Error: drpcli jobs show [id] [flags] requires 1 argument"
var jobShowMissingArgErrorString string = "Error: jobs GET: john: Not Found\n\n"
var jobShowJobString string = `RE:
{
  "Archived": false,
  "Available": true,
  "BootEnv": "local3",
  "EndTime": "0001-01-01T00:00:00Z",
  "Errors": [],
  "LogPath": "[\S\s]+/job-logs/00000000-0000-0000-0000-000000000001",
  "Machine": "3e7031fe-3062-45f1-835c-92541bc9cbd3",
  "Previous": "00000000-0000-0000-0000-000000000000",
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
var jobExistsMissingJohnString string = "Error: jobs GET: john: Not Found\n\n"

var jobExpireTimeErrorString string = "Error: Invalid UUID: false\n\n"
var jobDestroyBadString string = "Error: Jobs 00000000-0000-0000-0000-000000000001 is not in a deletable state: created\n\n"
var jobBadTimeFormatString string = "Error: parsing time \"fred\" as \"2006-01-02T15:04:05Z07:00\": cannot parse \"fred\" as \"2006\"\n\n"
var jobCreateJobInvalidMachineNameErrorString string = "Error: Unable to create new job: Invalid machine name passed to job create: james\n\n"

var jobUpdateNoArgErrorString string = "Error: drpcli jobs update [id] [json] [flags] requires 2 arguments\n"
var jobUpdateTooManyArgErrorString string = "Error: drpcli jobs update [id] [json] [flags] requires 2 arguments\n"

var jobShowMachineJohnString string = `{
  "Address": "192.168.100.110",
  "Available": true,
  "BootEnv": "local3",
  "CurrentJob": "00000000-0000-0000-0000-000000000001",
  "CurrentTask": 0,
  "Errors": [],
  "Name": "john",
  "Profile": {
    "Available": false,
    "Errors": null,
    "Name": "",
    "Tasks": null,
    "Validated": false
  },
  "Profiles": null,
  "Runnable": true,
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
var jobUpdateBadInputErrorString string = "Error: Jobs 00000000-0000-0000-0000-000000000001 wants State fred, which is not valid\n\n"
var jobUpdateJohnString string = `RE:
{
  "Archived": false,
  "Available": true,
  "BootEnv": "local3",
  "EndTime": "0001-01-01T00:00:00Z",
  "Errors": [],
  "LogPath": "[\S\s]+/job-logs/00000000-0000-0000-0000-000000000001",
  "Machine": "3e7031fe-3062-45f1-835c-92541bc9cbd3",
  "Previous": "00000000-0000-0000-0000-000000000000",
  "StartTime": "0001-01-01T00:00:00Z",
  "State": "incomplete",
  "Task": "task1",
  "Uuid": "00000000-0000-0000-0000-000000000001",
  "Validated": true
}
`

var jobUpdateJohnMissingErrorString string = "Error: jobs GET: john2: Not Found\n\n"

var jobPatchNoArgErrorString string = "Error: drpcli jobs patch [objectJson] [changesJson] [flags] requires 2 arguments\n"
var jobPatchTooManyArgErrorString = "Error: drpcli jobs patch [objectJson] [changesJson] [flags] requires 2 arguments\n"
var jobPatchBaseString = `{
  "Archived": false,
  "Available": true,
  "BootEnv": "local3",
  "Errors": [],
  "Machine": "3e7031fe-3062-45f1-835c-92541bc9cbd3",
  "Previous": "00000000-0000-0000-0000-000000000000",
  "State": "incomplete",
  "Task": "task1",
  "Uuid": "00000000-0000-0000-0000-000000000001",
  "Validated": true
}
`
var jobPatchBase2String = `{
  "Archived": false,
  "Available": true,
  "BootEnv": "local3",
  "Errors": [],
  "Machine": "3e7031fe-3062-45f1-835c-92541bc9cbd3",
  "Previous": "00000000-0000-0000-0000-000000000000",
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
var jobPatchBadInputErrorString = "Error: Jobs 00000000-0000-0000-0000-000000000001 wants State fred, which is not valid\n\n"
var jobPatchInputString = "{ \"State\": \"running\"}"
var jobPatchInputReplyString = `RE:
{
  "Archived": false,
  "Available": true,
  "BootEnv": "local3",
  "EndTime": "0001-01-01T00:00:00Z",
  "Errors": [],
  "LogPath": "[\S\s]+/job-logs/00000000-0000-0000-0000-000000000001",
  "Machine": "3e7031fe-3062-45f1-835c-92541bc9cbd3",
  "Previous": "00000000-0000-0000-0000-000000000000",
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
  "BootEnv": "local3",
  "EndTime": "0001-01-01T00:00:00Z",
  "Errors": [],
  "LogPath": "[\S\s]+/job-logs/00000000-0000-0000-0000-000000000001",
  "Machine": "3e7031fe-3062-45f1-835c-92541bc9cbd3",
  "Previous": "00000000-0000-0000-0000-000000000000",
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
  "BootEnv": "local3",
  "Errors": [],
  "Machine": "3e7031fe-3062-45f1-835c-92541bc9cbd3",
  "Previous": "00000000-0000-0000-0000-000000000000",
  "State": "incomplete",
  "Task": "task1",
  "Uuid": "10000000-0000-0000-0000-000000000001",
  "Validated": true
}
`
var jobPatchJohnMissingErrorString = "Error: jobs: PATCH 10000000-0000-0000-0000-000000000001: Not Found\n\n"

var jobUpdateFailedJobInputString = "{ \"State\": \"failed\" }"
var jobUpdateFailedJobUpdateString = `RE:
{
  "Archived": false,
  "Available": true,
  "BootEnv": "local3",
  "EndTime": "20[\s\S]*",
  "Errors": [],
  "LogPath": "[\S\s]+/job-logs/00000000-0000-0000-0000-000000000001",
  "Machine": "3e7031fe-3062-45f1-835c-92541bc9cbd3",
  "Previous": "00000000-0000-0000-0000-000000000000",
  "StartTime": "20[\s\S]*",
  "State": "failed",
  "Task": "task1",
  "Uuid": "00000000-0000-0000-0000-000000000001",
  "Validated": true
}
`
var jobCreateMachineNotRunningErrorString = "Error: Machine 3e7031fe-3062-45f1-835c-92541bc9cbd3 is not runnable\n\n"
var jobUpdateMachineRunnableString = `{
  "Address": "192.168.100.110",
  "Available": true,
  "BootEnv": "local3",
  "CurrentJob": "00000000-0000-0000-0000-000000000001",
  "CurrentTask": 0,
  "Errors": [],
  "Name": "john",
  "Profile": {
    "Available": false,
    "Errors": null,
    "Name": "",
    "Tasks": null,
    "Validated": false
  },
  "Profiles": null,
  "Runnable": true,
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
  "BootEnv": "local3",
  "EndTime": "0001-01-01T00:00:00Z",
  "Errors": [],
  "LogPath": "[\S\s]+/job-logs/00000000-0000-0000-0000-000000000002",
  "Machine": "3e7031fe-3062-45f1-835c-92541bc9cbd3",
  "Previous": "00000000-0000-0000-0000-000000000001",
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
  "BootEnv": "local3",
  "EndTime": "20[\s\S]*",
  "Errors": [],
  "LogPath": "[\S\s]+/job-logs/00000000-0000-0000-0000-000000000002",
  "Machine": "3e7031fe-3062-45f1-835c-92541bc9cbd3",
  "Previous": "00000000-0000-0000-0000-000000000001",
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
  "BootEnv": "local3",
  "EndTime": "0001-01-01T00:00:00Z",
  "Errors": [],
  "LogPath": "[\S\s]+/job-logs/00000000-0000-0000-0000-000000000003",
  "Machine": "3e7031fe-3062-45f1-835c-92541bc9cbd3",
  "Previous": "00000000-0000-0000-0000-000000000002",
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
  "BootEnv": "local3",
  "EndTime": "20[\s\S]*",
  "Errors": [],
  "LogPath": "[\S\s]+/job-logs/00000000-0000-0000-0000-000000000003",
  "Machine": "3e7031fe-3062-45f1-835c-92541bc9cbd3",
  "Previous": "00000000-0000-0000-0000-000000000002",
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
  "BootEnv": "local3",
  "EndTime": "0001-01-01T00:00:00Z",
  "Errors": [],
  "LogPath": "[\S\s]+/job-logs/00000000-0000-0000-0000-000000000004",
  "Machine": "3e7031fe-3062-45f1-835c-92541bc9cbd3",
  "Previous": "00000000-0000-0000-0000-000000000003",
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
  "BootEnv": "local3",
  "EndTime": "20[\s\S]*",
  "Errors": [],
  "LogPath": "[\S\s]+/job-logs/00000000-0000-0000-0000-000000000004",
  "Machine": "3e7031fe-3062-45f1-835c-92541bc9cbd3",
  "Previous": "00000000-0000-0000-0000-000000000003",
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
    "BootEnv": "local3",
    "EndTime": "20[\s\S]*",
    "Errors": [],
    "LogPath": "[\S\s]+/job-logs/00000000-0000-0000-0000-000000000001",
    "Machine": "3e7031fe-3062-45f1-835c-92541bc9cbd3",
    "Previous": "00000000-0000-0000-0000-000000000000",
    "StartTime": "20[\s\S]*",
    "State": "failed",
    "Task": "task1",
    "Uuid": "00000000-0000-0000-0000-000000000001",
    "Validated": true
  },
  {
    "Archived": false,
    "Available": true,
    "BootEnv": "local3",
    "EndTime": "20[\s\S]*",
    "Errors": [],
    "LogPath": "[\S\s]+/job-logs/00000000-0000-0000-0000-000000000002",
    "Machine": "3e7031fe-3062-45f1-835c-92541bc9cbd3",
    "Previous": "00000000-0000-0000-0000-000000000001",
    "StartTime": "0001-01-01T00:00:00Z",
    "State": "finished",
    "Task": "task1",
    "Uuid": "00000000-0000-0000-0000-000000000002",
    "Validated": true
  },
  {
    "Archived": false,
    "Available": true,
    "BootEnv": "local3",
    "EndTime": "20[\s\S]*",
    "Errors": [],
    "LogPath": "[\S\s]+/job-logs/00000000-0000-0000-0000-000000000003",
    "Machine": "3e7031fe-3062-45f1-835c-92541bc9cbd3",
    "Previous": "00000000-0000-0000-0000-000000000002",
    "StartTime": "0001-01-01T00:00:00Z",
    "State": "finished",
    "Task": "task2",
    "Uuid": "00000000-0000-0000-0000-000000000003",
    "Validated": true
  },
  {
    "Archived": false,
    "Available": true,
    "BootEnv": "local3",
    "EndTime": "20[\s\S]*",
    "Errors": [],
    "LogPath": "[\S\s]+/job-logs/00000000-0000-0000-0000-000000000004",
    "Machine": "3e7031fe-3062-45f1-835c-92541bc9cbd3",
    "Previous": "00000000-0000-0000-0000-000000000003",
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
var jobDestroyMissingJohnString = "Error: jobs: DELETE 3e7031fe-3062-45f1-835c-92541bc9cbd3: Not Found\n\n"
var jobDestroy001String = "Deleted job 00000000-0000-0000-0000-000000000001\n"
var jobDestroy002String = "Deleted job 00000000-0000-0000-0000-000000000002\n"
var jobDestroy003String = "Deleted job 00000000-0000-0000-0000-000000000003\n"
var jobDestroy004String = "Deleted job 00000000-0000-0000-0000-000000000004\n"

var jobActionsNoArgErrorString = "Error: drpcli jobs actions [id] [flags] requires 1 argument\n"
var jobActionsTooManyArgErrorString = "Error: drpcli jobs actions [id] [flags] requires 1 argument\n"
var jobActionsMissingJobErrorString = "Error: Job john does not exist\n\n"
var jobActionsRenderedTask1String = "[]\n"
var jobActionsRenderedTask2String = `[
  {
    "Content": "Fred rules",
    "Name": "part 1",
    "Path": ""
  }
]
`
var jobActionsMissingMachineRenderErrorString = "Error: Machine 3e7031fe-3062-45f1-835c-92541bc9cbd3 does not exist\n\n"
var jobActionsMissingTaskRenderErrorString = "Error: Task task2 does not exist\n\n"

var jobLogNoArgErrorString = "Error: drpcli jobs log [id] [- or string] [flags] requires at least 1 argument\n"
var jobLogTooManyArgsErrorString = "Error: drpcli jobs log [id] [- or string] [flags] requires at most 2 arguments\n"
var jobLogUnknownJobErrorString = "Error: Job john does not exist\n\n"

func TestJobCli(t *testing.T) {
	if err := os.MkdirAll("bootenvs", 0755); err != nil {
		t.Errorf("Failed to create bootenvs dir: %v\n", err)
	}
	if err := os.Symlink("../test-data/local3.yml", "bootenvs/local3.yml"); err != nil {
		t.Errorf("Failed to create link to local.yml: %v\n", err)
	}

	if err := os.MkdirAll("templates", 0755); err != nil {
		t.Errorf("Failed to create templates dir: %v\n", err)
	}
	tmpls := []string{"local3-pxelinux.tmpl", "local3-elilo.tmpl", "local3-ipxe.tmpl"}
	for _, tmpl := range tmpls {
		if err := os.Symlink("../test-data/"+tmpl, "templates/"+tmpl); err != nil {
			t.Errorf("Failed to create link to %s: %v\n", tmpl, err)
		}
	}

	tests := []CliTest{
		CliTest{false, false, []string{"bootenvs", "install", "bootenvs/local3.yml"}, noStdinString, bootEnvInstallLocalSuccessString, bootEnvInstallLocal3ErrorString},
		CliTest{false, false, []string{"tasks", "create", "task1"}, noStdinString, jobTask1Create, noErrorString},
		CliTest{false, false, []string{"tasks", "create", "-"}, jobTask2Create, jobTask2Create, noErrorString},
		CliTest{false, false, []string{"tasks", "create", "task3"}, noStdinString, jobTask3Create, noErrorString},
		CliTest{false, false, []string{"bootenvs", "create", jobLocal2CreateInput}, noStdinString, jobLocal2Create, noErrorString},
		CliTest{false, false, []string{"bootenvs", "update", "local3", jobLocalUpdateInput}, noStdinString, jobLocalUpdateString, noErrorString},

		CliTest{false, false, []string{"machines", "create", jobCreateMachineInputString}, noStdinString, jobCreateMachineJohnString, noErrorString},

		CliTest{true, false, []string{"jobs"}, noStdinString, "Access CLI commands relating to jobs\n", ""},
		CliTest{false, false, []string{"jobs", "list"}, noStdinString, jobDefaultListString, noErrorString},

		CliTest{true, true, []string{"jobs", "create"}, noStdinString, noContentString, jobCreateNoArgErrorString},
		CliTest{true, true, []string{"jobs", "create", "john", "john2"}, noStdinString, noContentString, jobCreateTooManyArgErrorString},
		CliTest{false, true, []string{"jobs", "create", jobCreateBadJSONString}, noStdinString, noContentString, jobCreateBadJSONErrorString},
		CliTest{false, true, []string{"jobs", "create", jobCreateBadJSON2String}, noStdinString, noContentString, jobCreateBadJSON2ErrorString},
		CliTest{false, false, []string{"jobs", "create", jobCreateInputString}, noStdinString, jobCreateJohnString, noErrorString},
		CliTest{false, true, []string{"jobs", "create", jobCreateInputString}, noStdinString, noContentString, jobCreateJobAlreadyRunningErrorString},

		CliTest{false, true, []string{"jobs", "create", "3e7031fe-3062-45f1-835c-92541bc9cbd3"}, noStdinString, noContentString, jobCreateJobAlreadyRunningErrorString},
		CliTest{false, true, []string{"jobs", "create", "john"}, noStdinString, noContentString, jobCreateJobAlreadyRunningErrorString},
		CliTest{false, true, []string{"jobs", "create", "james"}, noStdinString, noContentString, jobCreateJobInvalidMachineNameErrorString},

		CliTest{true, true, []string{"jobs", "show"}, noStdinString, noContentString, jobShowNoArgErrorString},
		CliTest{true, true, []string{"jobs", "show", "john", "john2"}, noStdinString, noContentString, jobShowTooManyArgErrorString},
		CliTest{false, true, []string{"jobs", "show", "john"}, noStdinString, noContentString, jobShowMissingArgErrorString},
		CliTest{false, false, []string{"jobs", "show", "00000000-0000-0000-0000-000000000001"}, noStdinString, jobShowJobString, noErrorString},
		CliTest{false, false, []string{"jobs", "show", "Key:00000000-0000-0000-0000-000000000001"}, noStdinString, jobShowJobString, noErrorString},
		CliTest{false, false, []string{"jobs", "show", "Uuid:00000000-0000-0000-0000-000000000001"}, noStdinString, jobShowJobString, noErrorString},

		CliTest{true, true, []string{"jobs", "exists"}, noStdinString, noContentString, jobExistsNoArgErrorString},
		CliTest{true, true, []string{"jobs", "exists", "john", "john2"}, noStdinString, noContentString, jobExistsTooManyArgErrorString},
		CliTest{false, false, []string{"jobs", "exists", "00000000-0000-0000-0000-000000000001"}, noStdinString, jobExistsJobString, noErrorString},
		CliTest{false, true, []string{"jobs", "exists", "john"}, noStdinString, noContentString, jobExistsMissingJohnString},
		CliTest{false, false, []string{"jobs", "exists", "Uuid:00000000-0000-0000-0000-000000000001"}, noStdinString, jobExistsJobString, noErrorString},
		CliTest{false, false, []string{"machines", "show", "3e7031fe-3062-45f1-835c-92541bc9cbd3"}, noStdinString, jobShowMachineJohnString, noErrorString},

		CliTest{false, false, []string{"jobs", "list"}, noStdinString, jobListJobsString, noErrorString},
		CliTest{false, false, []string{"jobs", "list", "--limit=0"}, noStdinString, jobEmptyListString, noErrorString},
		CliTest{false, false, []string{"jobs", "list", "--limit=10", "--offset=0"}, noStdinString, jobListJobsString, noErrorString},
		CliTest{false, false, []string{"jobs", "list", "--limit=10", "--offset=10"}, noStdinString, jobEmptyListString, noErrorString},
		CliTest{false, true, []string{"jobs", "list", "--limit=-10", "--offset=0"}, noStdinString, noContentString, limitNegativeError},
		CliTest{false, true, []string{"jobs", "list", "--limit=10", "--offset=-10"}, noStdinString, noContentString, offsetNegativeError},
		CliTest{false, false, []string{"jobs", "list", "--limit=-1", "--offset=-1"}, noStdinString, jobListJobsString, noErrorString},
		CliTest{false, false, []string{"jobs", "list", "BootEnv=local3"}, noStdinString, jobListJobsString, noErrorString},
		CliTest{false, false, []string{"jobs", "list", "BootEnv=false"}, noStdinString, jobEmptyListString, noErrorString},
		CliTest{false, false, []string{"jobs", "list", "Task=task1"}, noStdinString, jobListJobsString, noErrorString},
		CliTest{false, false, []string{"jobs", "list", "Task=false"}, noStdinString, jobEmptyListString, noErrorString},
		CliTest{false, false, []string{"jobs", "list", "State=created"}, noStdinString, jobListJobsString, noErrorString},
		CliTest{false, false, []string{"jobs", "list", "State=false"}, noStdinString, jobEmptyListString, noErrorString},
		CliTest{false, false, []string{"jobs", "list", "Machine=3e7031fe-3062-45f1-835c-92541bc9cbd3"}, noStdinString, jobListJobsString, noErrorString},
		CliTest{false, false, []string{"jobs", "list", "Machine=4e7031fe-3062-45f1-835c-92541bc9cbd3"}, noStdinString, jobEmptyListString, noErrorString},
		CliTest{false, true, []string{"jobs", "list", "Machine=false"}, noStdinString, noContentString, jobExpireTimeErrorString},
		CliTest{false, false, []string{"jobs", "list", "Archived=false"}, noStdinString, jobListJobsString, noErrorString},
		CliTest{false, false, []string{"jobs", "list", "Archived=true"}, noStdinString, jobEmptyListString, noErrorString},
		CliTest{false, false, []string{"jobs", "list", "StartTime=0001-01-01T00:00:00Z"}, noStdinString, jobListJobsString, noErrorString},
		CliTest{false, false, []string{"jobs", "list", "StartTime=2001-01-01T00:00:00Z"}, noStdinString, jobEmptyListString, noErrorString},
		CliTest{false, true, []string{"jobs", "list", "StartTime=fred"}, noStdinString, noContentString, jobBadTimeFormatString},
		CliTest{false, false, []string{"jobs", "list", "EndTime=0001-01-01T00:00:00Z"}, noStdinString, jobListJobsString, noErrorString},
		CliTest{false, false, []string{"jobs", "list", "EndTime=2001-01-01T00:00:00Z"}, noStdinString, jobEmptyListString, noErrorString},
		CliTest{false, true, []string{"jobs", "list", "EndTime=fred"}, noStdinString, noContentString, jobBadTimeFormatString},
		CliTest{false, false, []string{"jobs", "list", "UUID=4e7031fe-3062-45f1-835c-92541bc9cbd3"}, noStdinString, jobEmptyListString, noErrorString},
		CliTest{false, false, []string{"jobs", "list", "UUID=00000000-0000-0000-0000-000000000001"}, noStdinString, jobListJobsString, noErrorString},
		CliTest{false, true, []string{"jobs", "list", "UUID=false"}, noStdinString, noContentString, jobExpireTimeErrorString},

		CliTest{false, true, []string{"jobs", "destroy", "00000000-0000-0000-0000-000000000001"}, noStdinString, noContentString, jobDestroyBadString},

		CliTest{true, true, []string{"jobs", "log"}, noStdinString, noContentString, jobLogNoArgErrorString},
		CliTest{true, true, []string{"jobs", "log", "john", "john2", "john3"}, noStdinString, noContentString, jobLogTooManyArgsErrorString},
		CliTest{false, true, []string{"jobs", "log", "john"}, noStdinString, noContentString, jobLogUnknownJobErrorString},
		CliTest{false, false, []string{"jobs", "log", "00000000-0000-0000-0000-000000000001"}, noStdinString, "Log for Job: 00000000-0000-0000-0000-000000000001\n", noErrorString},
		CliTest{false, false, []string{"jobs", "log", "00000000-0000-0000-0000-000000000001", "Fred\n"}, noStdinString, "Success\n", noErrorString},
		CliTest{false, false, []string{"jobs", "log", "00000000-0000-0000-0000-000000000001"}, noStdinString, "Log for Job: 00000000-0000-0000-0000-000000000001\nFred\n", noErrorString},
		CliTest{false, false, []string{"jobs", "log", "00000000-0000-0000-0000-000000000001", "-"}, "Freddy\n", "Success\n", noErrorString},
		CliTest{false, false, []string{"jobs", "log", "00000000-0000-0000-0000-000000000001"}, noStdinString, "Log for Job: 00000000-0000-0000-0000-000000000001\nFred\nFreddy\n", noErrorString},

		CliTest{true, true, []string{"jobs", "update"}, noStdinString, noContentString, jobUpdateNoArgErrorString},
		CliTest{true, true, []string{"jobs", "update", "john", "john2", "john3"}, noStdinString, noContentString, jobUpdateTooManyArgErrorString},
		CliTest{false, true, []string{"jobs", "update", "00000000-0000-0000-0000-000000000001", jobUpdateBadJSONString}, noStdinString, noContentString, jobUpdateBadJSONErrorString},
		CliTest{false, true, []string{"jobs", "update", "00000000-0000-0000-0000-000000000001", jobUpdateBadJSON2String}, noStdinString, noContentString, jobUpdateBadJSON2ErrorString},
		CliTest{false, true, []string{"jobs", "update", "00000000-0000-0000-0000-000000000001", jobUpdateBadInputString}, noStdinString, noContentString, jobUpdateBadInputErrorString},
		CliTest{false, false, []string{"jobs", "update", "00000000-0000-0000-0000-000000000001", jobUpdateInputString}, noStdinString, jobUpdateJohnString, noErrorString},
		CliTest{false, true, []string{"jobs", "update", "john2", jobUpdateInputString}, noStdinString, noContentString, jobUpdateJohnMissingErrorString},
		CliTest{false, false, []string{"jobs", "show", "00000000-0000-0000-0000-000000000001"}, noStdinString, jobUpdateJohnString, noErrorString},
		// This tests that incomplete jobs come back.
		CliTest{false, false, []string{"jobs", "create", "john"}, noStdinString, jobUpdateJohnString, noErrorString},

		CliTest{true, true, []string{"jobs", "patch"}, noStdinString, noContentString, jobPatchNoArgErrorString},
		CliTest{true, true, []string{"jobs", "patch", "john", "john2", "john3"}, noStdinString, noContentString, jobPatchTooManyArgErrorString},
		CliTest{false, true, []string{"jobs", "patch", jobPatchBaseString, jobPatchBadPatchJSONString}, noStdinString, noContentString, jobPatchBadPatchJSONErrorString},
		CliTest{false, true, []string{"jobs", "patch", jobPatchBaseString, jobPatchBadPatchJSON2String}, noStdinString, noContentString, jobPatchBadPatchJSON2ErrorString},
		CliTest{false, true, []string{"jobs", "patch", jobPatchBadBaseJSONString, jobPatchInputString}, noStdinString, noContentString, jobPatchBadBaseJSONErrorString},
		CliTest{false, true, []string{"jobs", "patch", jobPatchBaseString, jobPatchBadInputString}, noStdinString, noContentString, jobPatchBadInputErrorString},
		CliTest{false, false, []string{"jobs", "patch", jobPatchBaseString, jobPatchInputString}, noStdinString, jobPatchInputReplyString, noErrorString},
		CliTest{false, false, []string{"jobs", "patch", jobPatchBase2String, jobPatchInput2String}, noStdinString, jobPatchJohnString, noErrorString},
		CliTest{false, true, []string{"jobs", "patch", jobPatchMissingBaseString, jobPatchInputString}, noStdinString, noContentString, jobPatchJohnMissingErrorString},
		CliTest{false, false, []string{"jobs", "show", "00000000-0000-0000-0000-000000000001"}, noStdinString, jobPatchJohnString, noErrorString},

		// This tests that incomplet jobs come back again
		CliTest{false, false, []string{"jobs", "create", "-"}, jobCreateNextInputString + "\n", jobPatchJohnString, noErrorString},
		CliTest{false, false, []string{"jobs", "update", "00000000-0000-0000-0000-000000000001", jobUpdateFailedJobInputString}, noStdinString, jobUpdateFailedJobUpdateString, noErrorString},
		CliTest{false, true, []string{"jobs", "create", "-"}, jobCreateNextInputString + "\n", noContentString, jobCreateMachineNotRunningErrorString},
		CliTest{false, false, []string{"machines", "update", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "{ \"Runnable\": true }"}, noStdinString, jobUpdateMachineRunnableString, noErrorString},
		CliTest{false, false, []string{"jobs", "create", "-"}, jobCreateNextInputString + "\n", jobCreateNextString, noErrorString},
		CliTest{false, false, []string{"jobs", "update", "00000000-0000-0000-0000-000000000002", jobUpdateFinishedJobInputString}, noStdinString, jobUpdateFinishedJob2UpdateString, noErrorString},
		CliTest{false, false, []string{"jobs", "create", "-"}, jobCreateNextInput3String + "\n", jobCreateNext3String, noErrorString},
		CliTest{false, false, []string{"jobs", "update", "00000000-0000-0000-0000-000000000003", jobUpdateFinishedJobInputString}, noStdinString, jobUpdateFinishedJob3UpdateString, noErrorString},
		CliTest{false, false, []string{"jobs", "create", "-"}, jobCreateNextInput4String + "\n", jobCreateNext4String, noErrorString},
		CliTest{false, false, []string{"jobs", "update", "00000000-0000-0000-0000-000000000004", jobUpdateFinishedJobInputString}, noStdinString, jobUpdateFinishedJob4UpdateString, noErrorString},
		CliTest{false, false, []string{"jobs", "create", "-"}, jobCreateNextInput5String + "\n", "null\n", noErrorString},

		CliTest{false, false, []string{"jobs", "list"}, noStdinString, jobFullListString, noErrorString},

		CliTest{true, true, []string{"jobs", "actions"}, noStdinString, noContentString, jobActionsNoArgErrorString},
		CliTest{true, true, []string{"jobs", "actions", "john", "june"}, noStdinString, noContentString, jobActionsTooManyArgErrorString},
		CliTest{false, true, []string{"jobs", "actions", "john"}, noStdinString, noContentString, jobActionsMissingJobErrorString},
		CliTest{false, false, []string{"jobs", "actions", "00000000-0000-0000-0000-000000000001"}, noStdinString, jobActionsRenderedTask1String, noErrorString},
		CliTest{false, false, []string{"jobs", "actions", "00000000-0000-0000-0000-000000000003"}, noStdinString, jobActionsRenderedTask2String, noErrorString},
		CliTest{false, false, []string{"machines", "destroy", "3e7031fe-3062-45f1-835c-92541bc9cbd3"}, noStdinString, machineDestroyJohnString, noErrorString},
		CliTest{false, true, []string{"jobs", "actions", "00000000-0000-0000-0000-000000000003"}, noStdinString, noContentString, jobActionsMissingMachineRenderErrorString},
		CliTest{false, false, []string{"bootenvs", "destroy", "local3"}, noStdinString, "Deleted bootenv local3\n", noErrorString},
		CliTest{false, false, []string{"bootenvs", "destroy", "local2"}, noStdinString, "Deleted bootenv local2\n", noErrorString},
		CliTest{false, false, []string{"tasks", "destroy", "task1"}, noStdinString, "Deleted task task1\n", noErrorString},
		CliTest{false, false, []string{"tasks", "destroy", "task2"}, noStdinString, "Deleted task task2\n", noErrorString},
		CliTest{false, false, []string{"tasks", "destroy", "task3"}, noStdinString, "Deleted task task3\n", noErrorString},
		CliTest{false, true, []string{"jobs", "actions", "00000000-0000-0000-0000-000000000003"}, noStdinString, noContentString, jobActionsMissingTaskRenderErrorString},

		CliTest{true, true, []string{"jobs", "destroy"}, noStdinString, noContentString, jobDestroyNoArgErrorString},
		CliTest{true, true, []string{"jobs", "destroy", "john", "june"}, noStdinString, noContentString, jobDestroyTooManyArgErrorString},
		CliTest{false, true, []string{"jobs", "destroy", "3e7031fe-3062-45f1-835c-92541bc9cbd3"}, noStdinString, noContentString, jobDestroyMissingJohnString},
		CliTest{false, false, []string{"jobs", "destroy", "00000000-0000-0000-0000-000000000001"}, noStdinString, jobDestroy001String, noErrorString},
		CliTest{false, false, []string{"jobs", "destroy", "00000000-0000-0000-0000-000000000002"}, noStdinString, jobDestroy002String, noErrorString},
		CliTest{false, false, []string{"jobs", "destroy", "00000000-0000-0000-0000-000000000003"}, noStdinString, jobDestroy003String, noErrorString},
		CliTest{false, false, []string{"jobs", "destroy", "00000000-0000-0000-0000-000000000004"}, noStdinString, jobDestroy004String, noErrorString},
		CliTest{false, false, []string{"jobs", "list"}, noStdinString, jobDefaultListString, noErrorString},

		CliTest{false, false, []string{"templates", "destroy", "local3-pxelinux.tmpl"}, noStdinString, "Deleted template local3-pxelinux.tmpl\n", noErrorString},
		CliTest{false, false, []string{"templates", "destroy", "local3-elilo.tmpl"}, noStdinString, "Deleted template local3-elilo.tmpl\n", noErrorString},
		CliTest{false, false, []string{"templates", "destroy", "local3-ipxe.tmpl"}, noStdinString, "Deleted template local3-ipxe.tmpl\n", noErrorString},
	}

	for _, test := range tests {
		testCli(t, test)
	}

	os.RemoveAll("bootenvs")
	os.RemoveAll("templates")
	os.RemoveAll("isos")
	os.RemoveAll("ic")
}
