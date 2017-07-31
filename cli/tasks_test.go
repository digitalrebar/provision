package cli

// TODO: Add validations around templates and content checks.

import (
	"testing"
)

var taskAddrErrorString string = "Error: Invalid Address: fred\n\n"
var taskExpireTimeErrorString string = "Error: Invalid Address: false\n\n"

var taskDefaultListString string = "[]\n"
var taskEmptyListString string = "[]\n"

var taskShowNoArgErrorString string = "Error: drpcli tasks show [id] requires 1 argument\n"
var taskShowTooManyArgErrorString string = "Error: drpcli tasks show [id] requires 1 argument\n"
var taskShowMissingArgErrorString string = "Error: tasks GET: jill: Not Found\n\n"
var taskShowJohnString string = `{
  "Name": "john",
  "OptionalParams": null,
  "RequiredParams": null,
  "Templates": null
}
`

var taskExistsNoArgErrorString string = "Error: drpcli tasks exists [id] requires 1 argument"
var taskExistsTooManyArgErrorString string = "Error: drpcli tasks exists [id] requires 1 argument"
var taskExistsIgnoreString string = ""
var taskExistsMissingString string = "Error: tasks GET: jill: Not Found\n\n"

var taskCreateNoArgErrorString string = "Error: drpcli tasks create [json] requires 1 argument\n"
var taskCreateTooManyArgErrorString string = "Error: drpcli tasks create [json] requires 1 argument\n"
var taskCreateBadJSONString = "{asdgasdg}"
var taskCreateBadJSONErrorString = "Error: dataTracker create tasks: Empty key not allowed\n\n"
var taskCreateInputString string = `{
  "Name": "john",
  "OptionalParams": null,
  "RequiredParams": null,
  "Templates": null
}
`
var taskCreateJohnString string = `{
  "Name": "john",
  "OptionalParams": null,
  "RequiredParams": null,
  "Templates": null
}
`
var taskCreateDuplicateErrorString = "Error: dataTracker create tasks: john already exists\n\n"

var taskListTasksString = `[
  {
    "Name": "john",
    "OptionalParams": null,
    "RequiredParams": null,
    "Templates": null
  }
]
`
var taskListBothEnvsString = `[
  {
    "Name": "john",
    "OptionalParams": null,
    "RequiredParams": null,
    "Templates": null
  }
]
`

var taskUpdateNoArgErrorString string = "Error: drpcli tasks update [id] [json] requires 2 arguments"
var taskUpdateTooManyArgErrorString string = "Error: drpcli tasks update [id] [json] requires 2 arguments"
var taskUpdateBadJSONString = "asdgasdg"
var taskUpdateBadJSONErrorString = "Error: Unable to merge objects: json: cannot unmarshal string into Go value of type map[string]interface {}\n\n\n"
var taskUpdateInputString string = `{
  "OptionalParams": [ "jillparam" ]
}
`
var taskUpdateJohnString string = `{
  "Name": "john",
  "OptionalParams": [
    "jillparam"
  ],
  "RequiredParams": null,
  "Templates": null
}
`
var taskUpdateJohnMissingErrorString string = "Error: tasks GET: jill: Not Found\n\n"

var taskPatchNoArgErrorString string = "Error: drpcli tasks patch [objectJson] [changesJson] requires 2 arguments"
var taskPatchTooManyArgErrorString string = "Error: drpcli tasks patch [objectJson] [changesJson] requires 2 arguments"
var taskPatchBadPatchJSONString = "asdgasdg"
var taskPatchBadPatchJSONErrorString = "Error: Unable to parse drpcli tasks patch [objectJson] [changesJson] JSON asdgasdg\nError: error unmarshaling JSON: json: cannot unmarshal string into Go value of type models.Task\n\n"
var taskPatchBadBaseJSONString = "asdgasdg"
var taskPatchBadBaseJSONErrorString = "Error: Unable to parse drpcli tasks patch [objectJson] [changesJson] JSON asdgasdg\nError: error unmarshaling JSON: json: cannot unmarshal string into Go value of type models.Task\n\n"
var taskPatchOldBaseString string = `{
  "Name": "john",
  "OptionalParams": null,
  "RequiredParams": null,
  "Templates": null
}
`
var taskPatchOldBaseErrorString = "Error: Patch error at line 0: Test op failed.\n\n"
var taskPatchBaseString string = `{
  "Name": "john",
  "OptionalParams": [ "jillparam" ],
  "RequiredParams": null,
  "Templates": null
}
`
var taskPatchInputString string = `{
  "OptionalParams": [ "joan" ]
}
`
var taskPatchJohnString string = `{
  "Name": "john",
  "OptionalParams": [
    "joan"
  ],
  "RequiredParams": null,
  "Templates": null
}
`
var taskPatchMissingBaseString string = `{
  "Name": "jill",
  "OptionalParams": null,
  "RequiredParams": null,
  "Templates": null
}
`
var taskPatchJohnMissingErrorString string = "Error: tasks: PATCH jill: Not Found\n\n"
var taskPatchBadBaseString string = `{
  "Addr": "jill",
  "NextServer": "2.2.2.2",
  "Strategy": "NewStrat",
  "Token": "john"
}
`
var taskPatchBadBaseErrorString string = "Error: Cannot get key for obj: Invalid type passed to task create\n\n"

var taskDestroyNoArgErrorString string = "Error: drpcli tasks destroy [id] requires 1 argument"
var taskDestroyTooManyArgErrorString string = "Error: drpcli tasks destroy [id] requires 1 argument"
var taskDestroyJohnString string = "Deleted task john\n"
var taskDestroyMissingJohnString string = "Error: tasks: DELETE jill: Not Found\n\n"

func TestTaskCli(t *testing.T) {
	tests := []CliTest{
		CliTest{true, false, []string{"tasks"}, noStdinString, "Access CLI commands relating to tasks\n", ""},
		CliTest{false, false, []string{"tasks", "list"}, noStdinString, taskDefaultListString, noErrorString},

		CliTest{true, true, []string{"tasks", "create"}, noStdinString, noContentString, taskCreateNoArgErrorString},
		CliTest{true, true, []string{"tasks", "create", "john", "john2"}, noStdinString, noContentString, taskCreateTooManyArgErrorString},
		CliTest{false, true, []string{"tasks", "create", taskCreateBadJSONString}, noStdinString, noContentString, taskCreateBadJSONErrorString},
		CliTest{false, false, []string{"tasks", "create", taskCreateInputString}, noStdinString, taskCreateJohnString, noErrorString},
		CliTest{false, true, []string{"tasks", "create", taskCreateInputString}, noStdinString, noContentString, taskCreateDuplicateErrorString},
		CliTest{false, false, []string{"tasks", "list"}, noStdinString, taskListBothEnvsString, noErrorString},

		CliTest{false, false, []string{"tasks", "list", "--limit=0"}, noStdinString, taskEmptyListString, noErrorString},
		CliTest{false, false, []string{"tasks", "list", "--limit=10", "--offset=0"}, noStdinString, taskListTasksString, noErrorString},
		CliTest{false, false, []string{"tasks", "list", "--limit=10", "--offset=10"}, noStdinString, taskEmptyListString, noErrorString},
		CliTest{false, true, []string{"tasks", "list", "--limit=-10", "--offset=0"}, noStdinString, noContentString, limitNegativeError},
		CliTest{false, true, []string{"tasks", "list", "--limit=10", "--offset=-10"}, noStdinString, noContentString, offsetNegativeError},
		CliTest{false, false, []string{"tasks", "list", "--limit=-1", "--offset=-1"}, noStdinString, taskListTasksString, noErrorString},
		CliTest{false, false, []string{"tasks", "list", "Name=fred"}, noStdinString, taskEmptyListString, noErrorString},
		CliTest{false, false, []string{"tasks", "list", "Name=john"}, noStdinString, taskListTasksString, noErrorString},

		CliTest{true, true, []string{"tasks", "show"}, noStdinString, noContentString, taskShowNoArgErrorString},
		CliTest{true, true, []string{"tasks", "show", "john", "john2"}, noStdinString, noContentString, taskShowTooManyArgErrorString},
		CliTest{false, true, []string{"tasks", "show", "jill"}, noStdinString, noContentString, taskShowMissingArgErrorString},
		CliTest{false, false, []string{"tasks", "show", "john"}, noStdinString, taskShowJohnString, noErrorString},

		CliTest{true, true, []string{"tasks", "exists"}, noStdinString, noContentString, taskExistsNoArgErrorString},
		CliTest{true, true, []string{"tasks", "exists", "john", "john2"}, noStdinString, noContentString, taskExistsTooManyArgErrorString},
		CliTest{false, true, []string{"tasks", "exists", "jill"}, noStdinString, noContentString, taskExistsMissingString},
		CliTest{false, false, []string{"tasks", "exists", "john"}, noStdinString, taskExistsIgnoreString, noErrorString},

		CliTest{true, true, []string{"tasks", "update"}, noStdinString, noContentString, taskUpdateNoArgErrorString},
		CliTest{true, true, []string{"tasks", "update", "john", "john2", "john3"}, noStdinString, noContentString, taskUpdateTooManyArgErrorString},
		CliTest{false, true, []string{"tasks", "update", "john", taskUpdateBadJSONString}, noStdinString, noContentString, taskUpdateBadJSONErrorString},
		CliTest{false, true, []string{"tasks", "update", "jill", taskUpdateInputString}, noStdinString, noContentString, taskUpdateJohnMissingErrorString},
		CliTest{false, false, []string{"tasks", "update", "john", taskUpdateInputString}, noStdinString, taskUpdateJohnString, noErrorString},
		CliTest{false, false, []string{"tasks", "show", "john"}, noStdinString, taskUpdateJohnString, noErrorString},

		CliTest{true, true, []string{"tasks", "patch"}, noStdinString, noContentString, taskPatchNoArgErrorString},
		CliTest{true, true, []string{"tasks", "patch", "john", "john2", "john3"}, noStdinString, noContentString, taskPatchTooManyArgErrorString},
		CliTest{false, true, []string{"tasks", "patch", taskPatchBaseString, taskPatchBadPatchJSONString}, noStdinString, noContentString, taskPatchBadPatchJSONErrorString},
		CliTest{false, true, []string{"tasks", "patch", taskPatchBadBaseJSONString, taskPatchInputString}, noStdinString, noContentString, taskPatchBadBaseJSONErrorString},
		CliTest{false, true, []string{"tasks", "patch", taskPatchMissingBaseString, taskPatchInputString}, noStdinString, noContentString, taskPatchJohnMissingErrorString},
		CliTest{false, true, []string{"tasks", "patch", taskPatchBadBaseString, taskPatchInputString}, noStdinString, noContentString, taskPatchBadBaseErrorString},
		CliTest{false, true, []string{"tasks", "patch", taskPatchOldBaseString, taskPatchInputString}, noStdinString, noContentString, taskPatchOldBaseErrorString},
		CliTest{false, false, []string{"tasks", "patch", taskPatchBaseString, taskPatchInputString}, noStdinString, taskPatchJohnString, noErrorString},
		CliTest{false, false, []string{"tasks", "show", "john"}, noStdinString, taskPatchJohnString, noErrorString},

		CliTest{true, true, []string{"tasks", "destroy"}, noStdinString, noContentString, taskDestroyNoArgErrorString},
		CliTest{true, true, []string{"tasks", "destroy", "john", "june"}, noStdinString, noContentString, taskDestroyTooManyArgErrorString},
		CliTest{false, false, []string{"tasks", "destroy", "john"}, noStdinString, taskDestroyJohnString, noErrorString},
		CliTest{false, true, []string{"tasks", "destroy", "jill"}, noStdinString, noContentString, taskDestroyMissingJohnString},
		CliTest{false, false, []string{"tasks", "list"}, noStdinString, taskDefaultListString, noErrorString},

		CliTest{false, false, []string{"tasks", "create", "-"}, taskCreateInputString + "\n", taskCreateJohnString, noErrorString},
		CliTest{false, false, []string{"tasks", "list"}, noStdinString, taskListBothEnvsString, noErrorString},
		CliTest{false, false, []string{"tasks", "update", "john", "-"}, taskUpdateInputString + "\n", taskUpdateJohnString, noErrorString},
		CliTest{false, false, []string{"tasks", "show", "john"}, noStdinString, taskUpdateJohnString, noErrorString},

		CliTest{false, false, []string{"tasks", "destroy", "john"}, noStdinString, taskDestroyJohnString, noErrorString},
		CliTest{false, false, []string{"tasks", "list"}, noStdinString, taskDefaultListString, noErrorString},
	}

	for _, test := range tests {
		testCli(t, test)
	}

}
