package cli

import (
	"testing"
)

func TestPluginCli(t *testing.T) {
	var pluginCreateBadJSONString = "{asdgasdg"
	var pluginCreateBadJSON2String = "[asdgasdg]"
	var pluginCreateMissingProviderInputString string = `{
  "Name": "i-woman"
}
`
	var pluginCreateMissingAllInputString string = `{
  "Description": "i-woman's plugin"
}
`
	var pluginCreateInputString string = `{
  "Name": "i-woman",
  "Provider": "incrementer"
}
`
	var pluginUpdateBadJSONString = "asdgasdg"

	var pluginUpdateInputString string = `{
  "Description": "lpxelinux.0"
}
`
	var pluginsParamsNextString string = `{
  "jj": 3
}
`

	cliTest(true, false, "plugins").run(t)
	cliTest(false, false, "plugins", "list").run(t)
	cliTest(true, true, "plugins", "create").run(t)
	cliTest(true, true, "plugins", "create", "john", "john2").run(t)
	cliTest(false, true, "plugins", "create", pluginCreateBadJSONString).run(t)
	cliTest(false, true, "plugins", "create", pluginCreateBadJSON2String).run(t)
	cliTest(false, true, "plugins", "create", pluginCreateMissingAllInputString).run(t)
	cliTest(false, true, "plugins", "create", pluginCreateMissingProviderInputString).run(t)
	cliTest(false, false, "plugins", "create", pluginCreateInputString).run(t)
	cliTest(false, true, "plugins", "create", pluginCreateInputString).run(t)
	cliTest(false, false, "plugins", "list").run(t)
	cliTest(false, false, "plugins", "list", "Name=fred").run(t)
	cliTest(false, false, "plugins", "list", "Name=i-woman").run(t)
	cliTest(false, false, "plugins", "list", "Provider=local").run(t)
	cliTest(false, false, "plugins", "list", "Provider=incrementer").run(t)
	cliTest(true, true, "plugins", "show").run(t)
	cliTest(true, true, "plugins", "show", "john", "john2").run(t)
	cliTest(false, true, "plugins", "show", "john").run(t)
	cliTest(false, false, "plugins", "show", "i-woman").run(t)
	cliTest(false, false, "plugins", "show", "Key:i-woman").run(t)
	cliTest(false, false, "plugins", "show", "Name:i-woman").run(t)
	cliTest(true, true, "plugins", "exists").run(t)
	cliTest(true, true, "plugins", "exists", "john", "john2").run(t)
	cliTest(false, false, "plugins", "exists", "i-woman").run(t)
	cliTest(false, true, "plugins", "exists", "john").run(t)
	cliTest(true, true, "plugins", "update").run(t)
	cliTest(true, true, "plugins", "update", "john", "john2", "john3").run(t)
	cliTest(false, true, "plugins", "update", "i-woman", pluginUpdateBadJSONString).run(t)
	cliTest(false, false, "plugins", "update", "i-woman", pluginUpdateInputString).run(t)
	cliTest(false, true, "plugins", "update", "john2", pluginUpdateInputString).run(t)
	cliTest(false, false, "plugins", "show", "i-woman").run(t)
	cliTest(false, false, "plugins", "show", "i-woman").run(t)
	cliTest(true, true, "plugins", "destroy").run(t)
	cliTest(true, true, "plugins", "destroy", "john", "june").run(t)
	cliTest(false, false, "plugins", "destroy", "i-woman").run(t)
	cliTest(false, true, "plugins", "destroy", "i-woman").run(t)
	cliTest(false, false, "plugins", "list").run(t)
	cliTest(false, false, "plugins", "create", "-").Stdin(pluginCreateInputString + "\n").run(t)
	cliTest(false, false, "plugins", "list").run(t)
	cliTest(false, false, "plugins", "update", "i-woman", "-").Stdin(pluginUpdateInputString + "\n").run(t)
	cliTest(false, false, "plugins", "show", "i-woman").run(t)
	cliTest(true, true, "plugins", "get").run(t)
	cliTest(false, true, "plugins", "get", "john", "param", "john2").run(t)
	cliTest(false, false, "plugins", "get", "i-woman", "param", "john2").run(t)
	cliTest(true, true, "plugins", "set").run(t)
	cliTest(false, true, "plugins", "set", "john", "param", "john2", "to", "cow").run(t)
	cliTest(false, false, "plugins", "set", "i-woman", "param", "john2", "to", "cow").run(t)
	cliTest(false, false, "plugins", "get", "i-woman", "param", "john2").run(t)
	cliTest(false, false, "plugins", "set", "i-woman", "param", "john2", "to", "3").run(t)
	cliTest(false, false, "plugins", "set", "i-woman", "param", "john3", "to", "4").run(t)
	cliTest(false, false, "plugins", "get", "i-woman", "param", "john2").run(t)
	cliTest(false, false, "plugins", "get", "i-woman", "param", "john3").run(t)
	cliTest(false, false, "plugins", "set", "i-woman", "param", "john2", "to", "null").run(t)
	cliTest(false, false, "plugins", "get", "i-woman", "param", "john2").run(t)
	cliTest(false, false, "plugins", "get", "i-woman", "param", "john3").run(t)
	cliTest(true, true, "plugins", "params").run(t)
	cliTest(false, true, "plugins", "params", "john2").run(t)
	cliTest(false, false, "plugins", "params", "i-woman").run(t)
	cliTest(false, true, "plugins", "params", "john2", pluginsParamsNextString).run(t)
	cliTest(false, false, "plugins", "params", "i-woman", pluginsParamsNextString).run(t)
	cliTest(false, false, "plugins", "params", "i-woman").run(t)
	cliTest(false, false, "plugins", "show", "i-woman").run(t)
	cliTest(false, false, "plugins", "destroy", "i-woman").run(t)
	cliTest(false, false, "plugins", "list").run(t)
}

func TestPluginActionsInTaskList(t *testing.T) {
	// Make a noisy task that sleeps some to test log write coalescing
	cliTest(false, false, "plugins", "create", "-").Stdin(`---
Name: incr
Provider: incrementer
`).run(t)
	cliTest(false, false, "machines", "create", "-").Stdin(`---
Name: phred
Uuid: c9196b77-deef-4c8e-8130-299b3e3d9a10
Tasks:
  - action:increment
  - action:incr:increment
Params:
  incrementer/step: 2
Runnable: true`).run(t)
	cliTest(false, false, "machines", "processjobs", "c9196b77-deef-4c8e-8130-299b3e3d9a10", "--oneshot").run(t)
	cliTest(false, false, "machines", "show", "Name:phred").run(t)
	cliTest(false, false, "machines", "jobs", "current", "Name:phred").run(t)
	cliTest(false, false, "machines", "currentlog", "Name:phred").run(t)
	cliTest(false, false, "machines", "update", "Name:phred", "-").Stdin(`---
Tasks:
  - action:reset_count
CurrentTask: -1
Runnable: true
`).run(t)
	cliTest(false, false, "machines", "processjobs", "c9196b77-deef-4c8e-8130-299b3e3d9a10", "--oneshot").run(t)
	cliTest(false, false, "machines", "show", "Name:phred").run(t)
	cliTest(false, false, "machines", "jobs", "current", "Name:phred").run(t)
	cliTest(false, false, "machines", "currentlog", "Name:phred").run(t)
	cliTest(false, false, "machines", "update", "Name:phred", "-").Stdin(`---
Tasks:
  - action:explode
CurrentTask: -1
Runnable: true
`).run(t)

	cliTest(false, false, "machines", "processjobs", "c9196b77-deef-4c8e-8130-299b3e3d9a10", "--oneshot").run(t)
	cliTest(false, false, "machines", "show", "Name:phred").run(t)
	cliTest(false, false, "machines", "jobs", "current", "Name:phred").run(t)
	cliTest(false, false, "machines", "currentlog", "Name:phred").run(t)
	cliTest(false, false, "machines", "update", "Name:phred", "-").Stdin(`---
Tasks:
  - action:incr:explode
CurrentTask: -1
Runnable: true
`).run(t)
	cliTest(false, false, "machines", "processjobs", "c9196b77-deef-4c8e-8130-299b3e3d9a10", "--oneshot").run(t)
	cliTest(false, false, "machines", "show", "Name:phred").run(t)
	cliTest(false, false, "machines", "jobs", "current", "Name:phred").run(t)
	cliTest(false, false, "machines", "currentlog", "Name:phred").run(t)
	cliTest(false, false, "machines", "update", "Name:phred", "-").Stdin(`---
Tasks:
  - action:TROGDOR:BURNINATE
CurrentTask: -1
Runnable: true
`).run(t)
	cliTest(false, false, "machines", "processjobs", "c9196b77-deef-4c8e-8130-299b3e3d9a10", "--oneshot").run(t)
	cliTest(false, false, "machines", "show", "Name:phred").run(t)
	cliTest(false, false, "machines", "jobs", "current", "Name:phred").run(t)
	cliTest(false, false, "machines", "currentlog", "Name:phred").run(t)
	cliTest(false, false, "machines", "deletejobs", "Name:phred").run(t)
	cliTest(false, false, "machines", "destroy", "Name:phred").run(t)
	cliTest(false, false, "plugins", "destroy", "incr").run(t)
	verifyClean(t)
}
