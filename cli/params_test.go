package cli

import (
	"testing"
)

func TestParamCli(t *testing.T) {
	var paramCreateBadJSONString = "{asdgasdg"
	var paramCreateBadJSON2String = "[asdgasdg]"
	var paramCreateInputString string = `{
  "Name": "john",
  "Schema": {
    "type": "string"
  }
}
`
	var paramUpdateBadJSONString = "asdgasdg"

	var paramUpdateInputString string = `{
  "Schema": {
    "type": "string"
  }
}
`
	cliTest(true, false, "params").run(t)
	cliTest(false, false, "params", "list").run(t)
	cliTest(true, true, "params", "create").run(t)
	cliTest(true, true, "params", "create", "john", "john2").run(t)
	cliTest(false, true, "params", "create", paramCreateBadJSONString).run(t)
	cliTest(false, true, "params", "create", paramCreateBadJSON2String).run(t)
	cliTest(false, false, "params", "create", paramCreateInputString).run(t)
	cliTest(false, true, "params", "create", paramCreateInputString).run(t)
	cliTest(false, false, "params", "list").run(t)
	cliTest(false, false, "params", "list", "Name=fred").run(t)
	cliTest(false, false, "params", "list", "Name=john").run(t)
	cliTest(true, true, "params", "show").run(t)
	cliTest(true, true, "params", "show", "john", "john2").run(t)
	cliTest(false, true, "params", "show", "john2").run(t)
	cliTest(false, false, "params", "show", "john").run(t)
	cliTest(true, true, "params", "exists").run(t)
	cliTest(true, true, "params", "exists", "john", "john2").run(t)
	cliTest(false, false, "params", "exists", "john").run(t)
	cliTest(false, true, "params", "exists", "john2").run(t)
	cliTest(true, true, "params", "exists", "john", "john2").run(t)
	cliTest(true, true, "params", "update").run(t)
	cliTest(true, true, "params", "update", "john", "john2", "john3").run(t)
	cliTest(false, true, "params", "update", "john", paramUpdateBadJSONString).run(t)
	cliTest(false, false, "params", "update", "john", paramUpdateInputString).run(t)
	cliTest(false, true, "params", "update", "john2", paramUpdateInputString).run(t)
	cliTest(false, false, "params", "show", "john").run(t)
	cliTest(false, false, "params", "show", "john").run(t)
	cliTest(true, true, "params", "destroy").run(t)
	cliTest(true, true, "params", "destroy", "john", "june").run(t)
	cliTest(false, false, "params", "destroy", "john").run(t)
	cliTest(false, true, "params", "destroy", "john").run(t)
	cliTest(false, false, "params", "list").run(t)
	cliTest(false, false, "params", "create", "-").Stdin(paramCreateInputString + "\n").run(t)
	cliTest(false, false, "params", "list").run(t)
	cliTest(false, false, "params", "update", "john", "-").Stdin(paramUpdateInputString + "\n").run(t)
	cliTest(false, false, "params", "show", "john").run(t)
	cliTest(false, false, "params", "destroy", "john").run(t)
	cliTest(false, false, "params", "list").run(t)
}
