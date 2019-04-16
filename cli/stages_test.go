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
	var stagesParamsNextString string = `{
  "jj": 3
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

	cliTest(true, true, "stages", "get").run(t)
	cliTest(false, true, "stages", "get", "john2", "param", "john2").run(t)
	cliTest(false, false, "stages", "get", "john", "param", "john2").run(t)
	cliTest(true, true, "stages", "add").run(t)
	cliTest(true, true, "stages", "add", "john2").run(t)
	cliTest(true, true, "stages", "add", "john2", "extra").run(t)
	cliTest(false, false, "stages", "add", "john", "param", "newparam", "to", "toast").run(t)
	cliTest(false, true, "stages", "add", "john", "param", "newparam", "to", "toast").run(t)
	cliTest(false, false, "stages", "get", "john", "param", "newparam").run(t)
	cliTest(false, false, "stages", "add", "john", "param", "newparam2", "to", "toast").run(t)
	cliTest(true, true, "stages", "remove").run(t)
	cliTest(true, true, "stages", "remove", "john2").run(t)
	cliTest(true, true, "stages", "remove", "john2", "extra").run(t)
	cliTest(false, false, "stages", "remove", "john", "param", "newparam").run(t)
	cliTest(false, false, "stages", "params", "john").run(t)
	cliTest(false, true, "stages", "remove", "john", "param", "newparam").run(t)
	cliTest(false, true, "stages", "remove", "john", "param", "newparam2", "--ref", "bagel").run(t)
	cliTest(false, false, "stages", "remove", "john", "param", "newparam2", "--ref", "toast").run(t)
	cliTest(false, false, "stages", "params", "john").run(t)
	cliTest(true, true, "stages", "set").run(t)
	cliTest(false, true, "stages", "set", "john2", "param", "john2", "to", "cow").run(t)
	cliTest(false, true, "stages", "set", "john", "param", "john2", "to", "cow", "--ref", "fred").run(t)
	cliTest(false, false, "stages", "set", "john", "param", "john2", "to", "cow").run(t)
	cliTest(false, false, "stages", "get", "john", "param", "john2").run(t)
	cliTest(false, false, "stages", "set", "john", "param", "john2", "to", "cow", "--ref", "cow").run(t)
	cliTest(false, false, "stages", "get", "john", "param", "john2").run(t)
	cliTest(false, false, "stages", "set", "john", "param", "john2", "to", "sow", "--ref", "cow").run(t)
	cliTest(false, false, "stages", "get", "john", "param", "john2").run(t)
	cliTest(false, false, "stages", "set", "john", "param", "john2", "to", "cow", "--ref", "sow").run(t)
	cliTest(false, false, "stages", "get", "john", "param", "john2").run(t)
	cliTest(false, false, "stages", "set", "john", "param", "john2", "to", "3").run(t)
	cliTest(false, false, "stages", "set", "john", "param", "john3", "to", "4").run(t)
	cliTest(false, false, "stages", "get", "john", "param", "john2").run(t)
	cliTest(false, false, "stages", "get", "john", "param", "john3").run(t)
	cliTest(false, false, "stages", "set", "john", "param", "john2", "to", "null").run(t)
	cliTest(false, false, "stages", "get", "john", "param", "john2").run(t)
	cliTest(false, false, "stages", "get", "john", "param", "john3").run(t)
	cliTest(false, false, "stages", "show", "john", "--slim", "params,meta").run(t)
	cliTest(false, false, "stages", "list", "--slim", "params,meta").run(t)
	cliTest(true, true, "stages", "params").run(t)
	cliTest(false, true, "stages", "params", "john2").run(t)
	cliTest(false, false, "stages", "params", "john").run(t)
	cliTest(false, true, "stages", "params", "john2", stagesParamsNextString).run(t)
	cliTest(false, true, "stages", "params", "john", stagesParamsNextString, "--ref", stagesParamsNextString).run(t)
	cliTest(false, false, "stages", "params", "john", stagesParamsNextString).run(t)
	cliTest(false, false, "stages", "params", "john").run(t)
	cliTest(false, false, "stages", "params", "john", stagesParamsNextString, "--ref", stagesParamsNextString).run(t)
	cliTest(false, false, "stages", "params", "john").run(t)
	cliTest(false, false, "stages", "params", "john", "{}", "--ref", stagesParamsNextString).run(t)
	cliTest(false, false, "stages", "params", "john").run(t)
	cliTest(false, false, "stages", "params", "john", stagesParamsNextString).run(t)
	cliTest(false, false, "stages", "params", "john").run(t)

	cliTest(false, false, "stages", "show", "john").run(t)
	cliTest(false, false, "stages", "destroy", "john").run(t)
	cliTest(false, false, "stages", "list").run(t)
	verifyClean(t)
}
