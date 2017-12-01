package cli

// TODO: Add validations around templates and content checks.

import (
	"testing"
)

var taskAddrErrorString string = "Error: Invalid Address: fred\n\n"
var taskExpireTimeErrorString string = "Error: Invalid Address: false\n\n"
var taskShowMissingArgErrorString string = "Error: GET: tasks/jill: Not Found\n\n"
var taskExistsMissingString string = "Error: GET: tasks/jill: Not Found\n\n"
var taskCreateBadJSONErrorString = "Error: CREATE: tasks: Empty key not allowed\n\n"
var taskCreateDuplicateErrorString = "Error: CREATE: tasks/john: already exists\n\n"
var taskUpdateJohnMissingErrorString string = "Error: GET: tasks/jill: Not Found\n\n"
var taskPatchOldBaseErrorString = `Error: PATCH: tasks/john
  Patch error at line 0: Test op failed.
  Patch line: {"op":"test","path":"/OptionalParams","from":"","value":[]}

`
var taskPatchJohnMissingErrorString string = "Error: PATCH: tasks/jill: Not Found\n\n"
var taskPatchBadBaseErrorString string = "Error: Cannot get key for obj: Invalid type passed to task create\n\n"
var taskDestroyMissingJohnString string = "Error: DELETE: tasks/jill: Not Found\n\n"

var taskDefaultListString string = "[]\n"
var taskEmptyListString string = "[]\n"

var taskShowNoArgErrorString string = "Error: drpcli tasks show [id] [flags] requires 1 argument\n"
var taskShowTooManyArgErrorString string = "Error: drpcli tasks show [id] [flags] requires 1 argument\n"

var taskShowJohnString string = `{
  "Available": true,
  "Errors": [],
  "Meta": {
    "feature-flags": "original-exit-codes"
  },
  "Name": "john",
  "OptionalParams": [],
  "ReadOnly": false,
  "RequiredParams": [],
  "Templates": [],
  "Validated": true
}
`

var taskExistsNoArgErrorString string = "Error: drpcli tasks exists [id] [flags] requires 1 argument"
var taskExistsTooManyArgErrorString string = "Error: drpcli tasks exists [id] [flags] requires 1 argument"
var taskExistsIgnoreString string = ""

var taskCreateNoArgErrorString string = "Error: drpcli tasks create [json] [flags] requires 1 argument\n"
var taskCreateTooManyArgErrorString string = "Error: drpcli tasks create [json] [flags] requires 1 argument\n"
var taskCreateBadJSONString = "{asdgasdg}"

var taskCreateInputString string = `{
  "Name": "john",
  "OptionalParams": [],
  "RequiredParams": [],
  "Templates": []
}
`
var taskCreateJohnString string = `{
  "Available": true,
  "Errors": [],
  "Meta": {
    "feature-flags": "original-exit-codes"
  },
  "Name": "john",
  "OptionalParams": [],
  "ReadOnly": false,
  "RequiredParams": [],
  "Templates": [],
  "Validated": true
}
`

var taskListTasksString = `[
  {
    "Available": true,
    "Errors": [],
    "Meta": {
      "feature-flags": "original-exit-codes"
    },
    "Name": "john",
    "OptionalParams": [],
    "ReadOnly": false,
    "RequiredParams": [],
    "Templates": [],
    "Validated": true
  }
]
`
var taskListBothEnvsString = `[
  {
    "Available": true,
    "Errors": [],
    "Meta": {
      "feature-flags": "original-exit-codes"
    },
    "Name": "john",
    "OptionalParams": [],
    "ReadOnly": false,
    "RequiredParams": [],
    "Templates": [],
    "Validated": true
  }
]
`

var taskUpdateNoArgErrorString string = "Error: drpcli tasks update [id] [json] [flags] requires 2 arguments"
var taskUpdateTooManyArgErrorString string = "Error: drpcli tasks update [id] [json] [flags] requires 2 arguments"
var taskUpdateBadJSONString = "asdgasdg"
var taskUpdateBadJSONErrorString = "Error: Unable to merge objects: json: cannot unmarshal string into Go value of type map[string]interface {}\n\n\n"
var taskUpdateInputString string = `{
  "OptionalParams": [ "jillparam" ]
}
`
var taskUpdateJohnString string = `{
  "Available": true,
  "Errors": [],
  "Meta": {
    "feature-flags": "original-exit-codes"
  },
  "Name": "john",
  "OptionalParams": [
    "jillparam"
  ],
  "ReadOnly": false,
  "RequiredParams": [],
  "Templates": [],
  "Validated": true
}
`

var taskPatchNoArgErrorString string = "Error: drpcli tasks patch [objectJson] [changesJson] [flags] requires 2 arguments"
var taskPatchTooManyArgErrorString string = "Error: drpcli tasks patch [objectJson] [changesJson] [flags] requires 2 arguments"
var taskPatchBadPatchJSONString = "asdgasdg"
var taskPatchBadPatchJSONErrorString = "Error: Unable to parse drpcli tasks patch [objectJson] [changesJson] [flags] JSON asdgasdg\nError: error unmarshaling JSON: json: cannot unmarshal string into Go value of type genmodels.Task\n\n"
var taskPatchBadBaseJSONString = "asdgasdg"
var taskPatchBadBaseJSONErrorString = "Error: Unable to parse drpcli tasks patch [objectJson] [changesJson] [flags] JSON asdgasdg\nError: error unmarshaling JSON: json: cannot unmarshal string into Go value of type genmodels.Task\n\n"
var taskPatchOldBaseString string = `{
  "Name": "john",
  "OptionalParams": [],
  "RequiredParams": [],
  "Templates": []
}
`

var taskPatchBaseString string = `{
  "Name": "john",
  "OptionalParams": [ "jillparam" ],
  "RequiredParams": [],
  "Templates": []
}
`
var taskPatchInputString string = `{
  "OptionalParams": [ "joan" ]
}
`
var taskPatchJohnString string = `{
  "Available": true,
  "Errors": [],
  "Meta": {
    "feature-flags": "original-exit-codes"
  },
  "Name": "john",
  "OptionalParams": [
    "joan"
  ],
  "ReadOnly": false,
  "RequiredParams": [],
  "Templates": [],
  "Validated": true
}
`
var taskPatchMissingBaseString string = `{
  "Name": "jill",
  "OptionalParams": [],
  "RequiredParams": [],
  "Templates": []
}
`

var taskPatchBadBaseString string = `{
  "Addr": "jill",
  "NextServer": "2.2.2.2",
  "Strategy": "NewStrat",
  "Token": "john"
}
`

var taskDestroyNoArgErrorString string = "Error: drpcli tasks destroy [id] [flags] requires 1 argument"
var taskDestroyTooManyArgErrorString string = "Error: drpcli tasks destroy [id] [flags] requires 1 argument"
var taskDestroyJohnString string = "Deleted task john\n"

func TestTaskCli(t *testing.T) {
	cliTest(true, false, "tasks").run(t)
	cliTest(false, false, "tasks", "list").run(t)
	cliTest(true, true, "tasks", "create").run(t)
	cliTest(true, true, "tasks", "create", "john", "john2").run(t)
	cliTest(false, true, "tasks", "create", taskCreateBadJSONString).run(t)
	cliTest(false, false, "tasks", "create", taskCreateInputString).run(t)
	cliTest(false, true, "tasks", "create", taskCreateInputString).run(t)
	cliTest(false, false, "tasks", "list").run(t)
	cliTest(false, false, "tasks", "list", "Name=fred").run(t)
	cliTest(false, false, "tasks", "list", "Name=john").run(t)
	cliTest(true, true, "tasks", "show").run(t)
	cliTest(true, true, "tasks", "show", "john", "john2").run(t)
	cliTest(false, true, "tasks", "show", "jill").run(t)
	cliTest(false, false, "tasks", "show", "john").run(t)
	cliTest(true, true, "tasks", "exists").run(t)
	cliTest(true, true, "tasks", "exists", "john", "john2").run(t)
	cliTest(false, true, "tasks", "exists", "jill").run(t)
	cliTest(false, false, "tasks", "exists", "john").run(t)
	cliTest(true, true, "tasks", "update").run(t)
	cliTest(true, true, "tasks", "update", "john", "john2", "john3").run(t)
	cliTest(false, true, "tasks", "update", "john", taskUpdateBadJSONString).run(t)
	cliTest(false, true, "tasks", "update", "jill", taskUpdateInputString).run(t)
	cliTest(false, false, "tasks", "update", "john", taskUpdateInputString).run(t)
	cliTest(false, false, "tasks", "show", "john").run(t)
	cliTest(false, false, "tasks", "show", "john").run(t)
	cliTest(true, true, "tasks", "destroy").run(t)
	cliTest(true, true, "tasks", "destroy", "john", "june").run(t)
	cliTest(false, false, "tasks", "destroy", "john").run(t)
	cliTest(false, true, "tasks", "destroy", "jill").run(t)
	cliTest(false, false, "tasks", "list").run(t)
	cliTest(false, false, "tasks", "create", "-").Stdin(taskCreateInputString + "\n").run(t)
	cliTest(false, false, "tasks", "list").run(t)
	cliTest(false, false, "tasks", "update", "john", "-").Stdin(taskUpdateInputString + "\n").run(t)
	cliTest(false, false, "tasks", "show", "john").run(t)
	cliTest(false, false, "tasks", "destroy", "john").run(t)
	cliTest(false, false, "tasks", "list").run(t)
}
