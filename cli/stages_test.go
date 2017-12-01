package cli

import (
	"testing"
)

func TestStageCli(t *testing.T) {

	var stageCreateBadJSONString = "{asdgasdg"

	var stageCreateBadJSON2String = "[asdgasdg]"
	var stageCreateInputString string = `{
  "Name": "john",
  "BootEnv": "local"
}
`
	var stageUpdateBadJSONString = "asdgasdg"
	var stageUpdateInputString string = `{
  "Description": "Awesome sauce"
}
`
	cliTest(true, false, "stages").run(t)
	cliTest(false, false, "stages", "list").run(t)
	cliTest(true, true, "stages", "create").run(t)
	cliTest(true, true, "stages", "create", "john", "john2").run(t)
	cliTest(false, true, "stages", "create", stageCreateBadJSONString).run(t)
	cliTest(false, true, "stages", "create", stageCreateBadJSON2String).run(t)
	cliTest(false, false, "stages", "create", stageCreateInputString).run(t)
	cliTest(false, true, "stages", "create", stageCreateInputString).run(t)
	cliTest(false, false, "stages", "list").run(t)
	cliTest(false, false, "stages", "list", "Name=fred").run(t)
	cliTest(false, false, "stages", "list", "Name=john").run(t)
	cliTest(false, false, "stages", "list", "BootEnv=fred").run(t)
	cliTest(false, false, "stages", "list", "BootEnv=local").run(t)
	cliTest(false, false, "stages", "list", "Reboot=true").run(t)
	cliTest(false, false, "stages", "list", "Reboot=false").run(t)
	cliTest(false, true, "stages", "list", "Reboot=fred").run(t)
	cliTest(true, true, "stages", "show").run(t)
	cliTest(true, true, "stages", "show", "john", "john2").run(t)
	cliTest(false, true, "stages", "show", "john2").run(t)
	cliTest(false, false, "stages", "show", "john").run(t)
	cliTest(true, true, "stages", "exists").run(t)
	cliTest(true, true, "stages", "exists", "john", "john2").run(t)
	cliTest(false, true, "stages", "exists", "john2").run(t)
	cliTest(true, true, "stages", "exists", "john", "john2").run(t)
	cliTest(false, false, "stages", "exists", "john").run(t)
	cliTest(true, true, "stages", "update").run(t)
	cliTest(true, true, "stages", "update", "john", "john2", "john3").run(t)
	cliTest(false, true, "stages", "update", "john", stageUpdateBadJSONString).run(t)
	cliTest(false, false, "stages", "update", "john", stageUpdateInputString).run(t)
	cliTest(false, true, "stages", "update", "john2", stageUpdateInputString).run(t)
	cliTest(false, false, "stages", "show", "john").run(t)
	cliTest(false, false, "stages", "show", "john").run(t)
	cliTest(true, true, "stages", "destroy").run(t)
	cliTest(true, true, "stages", "destroy", "john", "june").run(t)
	cliTest(false, false, "stages", "destroy", "john").run(t)
	cliTest(false, true, "stages", "destroy", "john").run(t)
	cliTest(false, false, "stages", "list").run(t)
	cliTest(false, false, "stages", "create", "-").Stdin(stageCreateInputString + "\n").run(t)
	cliTest(false, false, "stages", "list").run(t)
	cliTest(false, false, "stages", "update", "john", "-").Stdin(stageUpdateInputString + "\n").run(t)
	cliTest(false, false, "stages", "show", "john").run(t)
	cliTest(false, false, "stages", "destroy", "john").run(t)
	cliTest(false, false, "stages", "list").run(t)
}
