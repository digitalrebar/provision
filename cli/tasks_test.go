package cli

// TODO: Add validations around templates and content checks.

import (
	"testing"
)

func TestTaskCli(t *testing.T) {
	var taskCreateBadJSONString = "{asdgasdg}"

	var taskCreateInputString string = `{
  "Name": "john",
  "OptionalParams": [],
  "RequiredParams": [],
  "Templates": []
}
`
	var taskUpdateBadJSONString = "asdgasdg"
	var taskUpdateInputString string = `{
  "OptionalParams": [ "jillparam" ]
}
`

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
