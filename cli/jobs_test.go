package cli

import (
	"testing"
)

func TestJobCli(t *testing.T) {
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
      "Meta": {
        "foo":"bar"
      },
      "Name": "part 1",
      "Path": ""
    }
  ],
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
	var jobLocalUpdateInput string = `{
  "Tasks": ["task1","task2","task3"]
}
`
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
	var jobUpdateBadJSONString string = "{asgasdg"
	var jobUpdateBadJSON2String string = "[ \"gadsg\" ]"
	var jobUpdateInputString string = "{ \"State\": \"incomplete\" }"
	var jobUpdateBadInputString string = "{ \"State\": \"fred\" }"
	var jobUpdateFailedJobInputString = "{ \"State\": \"failed\" }"
	var jobUpdateFinishedJobInputString = "{ \"State\": \"finished\" }"
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
	verifyClean(t)
}

func TestJobOsFilter(t *testing.T) {
	cliTest(false, false, "tasks", "create", "-").Stdin(`---
Name: task1
Templates:
  - Name: t1
    Contents: 't1'
    Meta:
      OS: linux,darwin
  - Name: t2
    Contents: 't2'
    Meta:
      OS: darwin
  - Name: t3
    Contents: 't3'
    Meta:
      OS: linux
`).run(t)
	cliTest(false, false, "stages", "create", "-").Stdin(`---
Name: stage1
Tasks:
  - task1`).run(t)
	cliTest(false, false, "machines", "create", "-").Stdin(`---
Name: fred
Uuid: "3e7031fe-3062-45f1-835c-92541bc9cbd3"
Stage: stage1`).run(t)
	cliTest(false, false, "jobs", "create", "-").Stdin(`---
Uuid: "00000000-0000-0000-0000-000000000001"
Machine: "3e7031fe-3062-45f1-835c-92541bc9cbd3"
`).run(t)
	cliTest(false, false, "jobs", "actions", "00000000-0000-0000-0000-000000000001").run(t)
	cliTest(false, false, "jobs", "actions", "00000000-0000-0000-0000-000000000001", "--for-os", "").run(t)
	cliTest(false, false, "jobs", "actions", "00000000-0000-0000-0000-000000000001", "--for-os", "linux").run(t)
	cliTest(false, false, "jobs", "actions", "00000000-0000-0000-0000-000000000001", "--for-os", "darwin").run(t)
	cliTest(false, false, "jobs", "actions", "00000000-0000-0000-0000-000000000001", "--for-os", "windows").run(t)
	cliTest(false, false, "jobs", "update", "00000000-0000-0000-0000-000000000001", `{"State":"failed"}`).run(t)
	cliTest(false, false, "jobs", "destroy", "00000000-0000-0000-0000-000000000001").run(t)
	cliTest(false, false, "machines", "destroy", "Name:fred").run(t)
	cliTest(false, false, "stages", "destroy", "stage1").run(t)
	cliTest(false, false, "tasks", "destroy", "task1").run(t)
	verifyClean(t)
}
