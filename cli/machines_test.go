package cli

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/digitalrebar/provision/api"
	"github.com/digitalrebar/provision/models"
	"github.com/pborman/uuid"
)

var machineCreateInputString = `{
  "Address": "192.168.100.110",
  "name": "john",
  "Secret": "secret1",
  "Uuid": "3e7031fe-3062-45f1-835c-92541bc9cbd3",
  "bootenv": "local"
}
`
var machineCreateJohnString = `{
  "Address": "192.168.100.110",
  "Available": true,
  "BootEnv": "local",
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
  "Stage": "none",
  "Tasks": [],
  "Uuid": "3e7031fe-3062-45f1-835c-92541bc9cbd3",
  "Validated": true
}
`

var machineDestroyJohnString = "Deleted machine 3e7031fe-3062-45f1-835c-92541bc9cbd3\n"

var machinePluginCreateString = `{
  "Available": true,
  "Errors": [],
  "Name": "incr",
  "PluginErrors": [],
  "Provider": "incrementer",
  "ReadOnly": false,
  "Validated": true
}
`

func TestMachineCli(t *testing.T) {
	var machineCreateBadJSONString = "{asdgasdg"
	var machineCreateBadJSON2String = "[asdgasdg]"
	var machineUpdateBadJSONString = "asdgasdg"
	var machineUpdateInputString = `{
  "Description": "lpxelinux.0"
}
`
	var machinesParamsNextString = `{
  "jj": 3
}
`
	var machineRunActionMissingParameterStdinString = "{}"
	var machineRunActionGoodStdinString = `{
	"incrementer/parameter": "parm5",
	"incrementer/step": 10
}
`
	var machineStage1CreateString = `{
	"Name": "stage1",
	"BootEnv": "local",
	"Tasks": [ "jamie", "justine" ]
}
`
	var machineStage2CreateString = `{
  "Name": "stage2",
  "BootEnv": "local",
  "Templates": [
    {
      "Contents": "{{.Param \"sp-param\"}}",
      "Name": "test",
      "Path": "{{.Machine.Path}}/file"
    }
  ]
}
`
	var machineWorkflow1SetGood = `{
	"Name": "Workflow1Good",
	"Stages": [
		"none"
	]
}
`
	var machineWorkflow2SetBad = `{
	"Name": "Workflow2Bad",
	"Stages": [
		"nonexistent-stage"
	]
}
`
	cliTest(false, false, "profiles", "create", "jill").run(t)
	cliTest(false, false, "profiles", "create", "jean").run(t)
	cliTest(false, false, "profiles", "create", "stage-prof").run(t)
	cliTest(false, false, "tasks", "create", "jamie").run(t)
	cliTest(false, false, "tasks", "create", "justine").run(t)
	cliTest(false, false, "stages", "create", machineStage1CreateString).run(t)
	cliTest(false, false, "stages", "create", machineStage2CreateString).run(t)
	cliTest(false, false, "plugins", "create", machinePluginCreateString).run(t)
	cliTest(true, false, "machines").run(t)
	cliTest(false, false, "machines", "list").run(t)
	cliTest(true, true, "machines", "create").run(t)
	cliTest(true, true, "machines", "create", "john", "john2").run(t)
	cliTest(false, true, "machines", "create", machineCreateBadJSONString).run(t)
	cliTest(false, true, "machines", "create", machineCreateBadJSON2String).run(t)
	cliTest(false, false, "machines", "create", machineCreateInputString).run(t)
	cliTest(false, true, "machines", "create", machineCreateInputString).run(t)
	cliTest(false, false, "machines", "list").run(t)
	cliTest(false, false, "machines", "list", "Name=fred").run(t)
	cliTest(false, false, "machines", "list", "Name=john").run(t)
	cliTest(false, false, "machines", "list", "BootEnv=local").run(t)
	cliTest(false, false, "machines", "list", "BootEnv=false").run(t)
	cliTest(false, false, "machines", "list", "Address=192.168.100.110").run(t)
	cliTest(false, false, "machines", "list", "Address=1.1.1.1").run(t)
	cliTest(false, true, "machines", "list", "Address=fred").run(t)
	cliTest(false, false, "machines", "list", "Uuid=4e7031fe-3062-45f1-835c-92541bc9cbd3").run(t)
	cliTest(false, false, "machines", "list", "Uuid=3e7031fe-3062-45f1-835c-92541bc9cbd3").run(t)
	cliTest(false, true, "machines", "list", "Uuid=false").run(t)
	cliTest(false, false, "machines", "list", "Runnable=true").run(t)
	cliTest(false, false, "machines", "list", "Runnable=false").run(t)
	cliTest(false, true, "machines", "list", "Runnable=fred").run(t)
	cliTest(true, true, "machines", "show").run(t)
	cliTest(true, true, "machines", "show", "john", "john2").run(t)
	cliTest(false, true, "machines", "show", "john").run(t)
	cliTest(false, false, "machines", "show", "3e7031fe-3062-45f1-835c-92541bc9cbd3").run(t)
	cliTest(false, false, "machines", "show", "Key:3e7031fe-3062-45f1-835c-92541bc9cbd3").run(t)
	cliTest(false, false, "machines", "show", "Uuid:3e7031fe-3062-45f1-835c-92541bc9cbd3").run(t)
	cliTest(false, false, "machines", "show", "Name:john").run(t)
	cliTest(true, true, "machines", "exists").run(t)
	cliTest(true, true, "machines", "exists", "john", "john2").run(t)
	cliTest(false, false, "machines", "exists", "3e7031fe-3062-45f1-835c-92541bc9cbd3").run(t)
	cliTest(false, true, "machines", "exists", "john").run(t)
	cliTest(true, true, "machines", "exists", "john", "john2").run(t)
	cliTest(true, true, "machines", "update").run(t)
	cliTest(true, true, "machines", "update", "john", "john2", "john3").run(t)
	cliTest(false, true, "machines", "update", "3e7031fe-3062-45f1-835c-92541bc9cbd3", machineUpdateBadJSONString).run(t)
	cliTest(false, false, "machines", "update", "3e7031fe-3062-45f1-835c-92541bc9cbd3", machineUpdateInputString).run(t)
	cliTest(false, true, "machines", "update", "john2", machineUpdateInputString).run(t)
	cliTest(false, false, "machines", "show", "3e7031fe-3062-45f1-835c-92541bc9cbd3").run(t)
	cliTest(false, false, "machines", "show", "3e7031fe-3062-45f1-835c-92541bc9cbd3").run(t)
	cliTest(true, true, "machines", "destroy").run(t)
	cliTest(true, true, "machines", "destroy", "john", "june").run(t)
	cliTest(false, false, "machines", "destroy", "3e7031fe-3062-45f1-835c-92541bc9cbd3").run(t)
	cliTest(false, true, "machines", "destroy", "3e7031fe-3062-45f1-835c-92541bc9cbd3").run(t)
	cliTest(false, false, "machines", "list").run(t)
	cliTest(false, false, "machines", "create", "-").Stdin(machineCreateInputString + "\n").run(t)
	cliTest(false, false, "machines", "list").run(t)
	cliTest(false, false, "machines", "update", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "-").Stdin(machineUpdateInputString + "\n").run(t)
	cliTest(false, false, "machines", "show", "3e7031fe-3062-45f1-835c-92541bc9cbd3").run(t)
	// bootenv tests
	cliTest(true, true, "machines", "bootenv").run(t)
	cliTest(false, true, "machines", "bootenv", "john", "john2").run(t)
	cliTest(false, true, "machines", "bootenv", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "john2").run(t)
	// stage tests
	cliTest(true, true, "machines", "stage").run(t)
	cliTest(false, true, "machines", "stage", "john", "john2").run(t)
	cliTest(false, true, "machines", "stage", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "john2").run(t)
	cliTest(false, false, "machines", "stage", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "stage1").run(t)
	cliTest(false, false, "machines", "update", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "{ \"Runnable\": true }").run(t)
	cliTest(false, true, "machines", "stage", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "stage2").run(t)
	cliTest(false, false, "machines", "stage", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "stage2", "--force").run(t)
	cliTest(false, false, "machines", "stage", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "", "--force").run(t)
	// workflow tests
	cliTest(true, true, "machines", "workflow").run(t)
	cliTest(false, false, "workflows", "create", "-").Stdin(machineWorkflow1SetGood).run(t)
	cliTest(false, false, "workflows", "create", "-").Stdin(machineWorkflow2SetBad).run(t)
	cliTest(false, false, "machines", "workflow", "Name:john", "Workflow1Good").run(t)
	cliTest(false, true, "machines", "workflow", "Name:john", "Workflow2Bad").run(t)
	cliTest(false, false, "machines", "workflow", "Name:john", "").run(t)
	cliTest(false, false, "workflows", "destroy", "Workflow1Good").run(t)
	cliTest(false, false, "workflows", "destroy", "Workflow2Bad").run(t)
	// Add/Remove Profile tests
	cliTest(true, true, "machines", "addprofile").run(t)
	cliTest(false, false, "machines", "addprofile", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "jill").run(t)
	cliTest(false, false, "machines", "addprofile", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "jean").run(t)
	cliTest(false, true, "machines", "addprofile", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "jill").run(t)
	cliTest(false, false, "profiles", "set", "jill", "param", "jill-param", "to", "janga").run(t)
	cliTest(false, false, "profiles", "set", "stage-prof", "param", "sp-param", "to", "val").run(t)
	cliTest(false, false, "machines", "stage", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "stage2", "--force").run(t)
	cliTest(false, false, "stages", "addprofile", "stage2", "stage-prof").run(t)
	cliTest(false, false, "stages", "set", "stage2", "param", "sp-direct-param", "to", "val2").run(t)

	cliTest(false, false, "machines", "params", "3e7031fe-3062-45f1-835c-92541bc9cbd3").run(t)
	cliTest(false, false, "machines", "params", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "--aggregate").run(t)
	tr := &http.Transport{}
	client := &http.Client{Transport: tr}
	req, _ := http.NewRequest("GET", "http://127.0.0.1:10002/machines/3e7031fe-3062-45f1-835c-92541bc9cbd3/file", nil)
	req.SetBasicAuth("rocketskates", "r0cketsk8ts")
	rsp, apierr := client.Do(req)
	if apierr != nil {
		t.Errorf("FAIL: Failed to query machine file: %s", apierr)
	} else {
		defer rsp.Body.Close()
		body, err := ioutil.ReadAll(rsp.Body)
		if err != nil {
			t.Errorf("FAIL: Failed to read all: %s", err)
		}
		if string(body) != "val" {
			t.Errorf("FAIL: Body was: AA%sAA expected %s", string(body), "val")
		}
	}

	cliTest(false, false, "machines", "stage", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "", "--force").run(t)
	cliTest(true, true, "machines", "removeprofile").run(t)
	cliTest(false, false, "machines", "removeprofile", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "justine").run(t)
	cliTest(false, false, "machines", "removeprofile", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "jill").run(t)
	cliTest(false, false, "machines", "removeprofile", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "jean").run(t)
	cliTest(true, true, "machines", "get").run(t)
	cliTest(false, true, "machines", "get", "john", "param", "john2").run(t)
	cliTest(false, false, "machines", "get", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "param", "john2").run(t)
	cliTest(true, true, "machines", "set").run(t)
	cliTest(false, true, "machines", "set", "john", "param", "john2", "to", "cow").run(t)
	cliTest(false, false, "machines", "set", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "param", "john2", "to", "cow").run(t)
	cliTest(false, false, "machines", "get", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "param", "john2").run(t)
	cliTest(false, false, "machines", "set", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "param", "john2", "to", "3").run(t)
	cliTest(false, false, "machines", "set", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "param", "john3", "to", "4").run(t)
	cliTest(false, false, "machines", "get", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "param", "john2").run(t)
	cliTest(false, false, "machines", "get", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "param", "john3").run(t)
	cliTest(false, false, "machines", "set", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "param", "john2", "to", "null").run(t)
	cliTest(false, false, "machines", "get", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "param", "john2").run(t)
	cliTest(false, false, "machines", "get", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "param", "john3").run(t)
	cliTest(false, false, "machines", "set", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "param", "john4", "to", "-").Stdin("5").run(t)
	cliTest(false, false, "machines", "get", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "param", "john4").run(t)
	cliTest(false, false, "machines", "remove", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "param", "john4").run(t)
	cliTest(false, false, "machines", "add", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "param", "john5", "to", "-").Stdin("6").run(t)
	cliTest(false, false, "machines", "get", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "param", "john5").run(t)
	cliTest(false, false, "machines", "remove", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "param", "john5").run(t)
	cliTest(true, true, "machines", "actions").run(t)
	cliTest(false, true, "machines", "actions", "john").run(t)
	cliTest(false, false, "machines", "actions", "3e7031fe-3062-45f1-835c-92541bc9cbd3").run(t)
	cliTest(true, true, "machines", "action").run(t)
	cliTest(true, true, "machines", "action", "john").run(t)
	cliTest(false, true, "machines", "action", "john", "command").run(t)
	cliTest(false, true, "machines", "action", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "command").run(t)
	cliTest(false, false, "machines", "action", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "increment").run(t)
	cliTest(true, true, "machines", "runaction").run(t)
	cliTest(true, true, "machines", "runaction", "fred").run(t)
	cliTest(false, true, "machines", "runaction", "fred", "command").run(t)
	cliTest(false, true, "machines", "runaction", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "command").run(t)
	cliTest(false, false, "machines", "runaction", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "increment").run(t)
	cliTest(false, true, "machines", "runaction", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "increment", "fred").run(t)

	cliTest(false, false, "machines", "actions", "3e7031fe-3062-45f1-835c-92541bc9cbd3").run(t)
	cliTest(false, false, "machines", "action", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "reset_count").run(t)
	cliTest(false, false, "machines", "runaction", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "reset_count").run(t)
	cliTest(false, false, "machines", "actions", "3e7031fe-3062-45f1-835c-92541bc9cbd3").run(t)
	cliTest(false, true, "machines", "action", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "reset_count").run(t)
	cliTest(false, true, "machines", "runaction", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "reset_count").run(t)
	cliTest(false, false, "machines", "runaction", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "increment", "incrementer/parameter", "asgdasdg").run(t)
	cliTest(false, false, "machines", "runaction", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "increment", "incrementer/parameter", "parm1", "extra", "10").run(t)
	cliTest(false, false, "machines", "get", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "param", "parm1").run(t)
	cliTest(false, true, "machines", "runaction", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "increment", "incrementer/parameter", "parm2", "incrementer/step", "asgdasdg").run(t)
	cliTest(false, false, "machines", "get", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "param", "parm2").run(t)
	cliTest(false, false, "machines", "runaction", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "increment", "incrementer/parameter", "parm2", "incrementer/step", "10").run(t)
	cliTest(false, false, "machines", "get", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "param", "parm2").run(t)
	cliTest(false, true, "machines", "runaction", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "increment", "-").Stdin("fred").run(t)
	cliTest(false, false, "machines", "runaction", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "reset_count", "-").Stdin(machineRunActionMissingParameterStdinString).run(t)
	cliTest(false, true, "machines", "runaction", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "reset_count", "-").Stdin(machineRunActionMissingParameterStdinString).run(t)
	cliTest(false, false, "machines", "runaction", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "increment", "-").Stdin(machineRunActionMissingParameterStdinString).run(t)
	cliTest(false, false, "machines", "runaction", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "increment", "-").Stdin(machineRunActionGoodStdinString).run(t)
	cliTest(false, false, "machines", "runaction", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "increment", "-").Stdin(machineRunActionGoodStdinString).run(t)
	cliTest(false, false, "machines", "get", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "param", "parm5").run(t)
	cliTest(true, true, "machines", "wait").run(t)
	cliTest(true, true, "machines", "wait", "jk").run(t)
	cliTest(true, true, "machines", "wait", "jk", "jk").run(t)
	cliTest(true, true, "machines", "wait", "jk", "jk", "jk", "jk", "jk").run(t)
	cliTest(false, true, "machines", "wait", "jk", "jk", "jk", "jk").run(t)
	cliTest(false, true, "machines", "wait", "jk", "jk", "jk").run(t)
	cliTest(false, false, "machines", "wait", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "jk", "jk", "1").run(t)
	cliTest(false, false, "machines", "wait", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "BootEnv", "local", "1").run(t)
	cliTest(false, false, "machines", "wait", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "Runnable", "fred", "1").run(t)
	cliTest(true, true, "machines", "params").run(t)
	cliTest(false, true, "machines", "params", "john2").run(t)
	cliTest(false, false, "machines", "params", "3e7031fe-3062-45f1-835c-92541bc9cbd3").run(t)
	cliTest(false, true, "machines", "params", "john2", machinesParamsNextString).run(t)
	cliTest(false, false, "machines", "params", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "-").Stdin(machinesParamsNextString).run(t)
	cliTest(false, false, "machines", "params", "3e7031fe-3062-45f1-835c-92541bc9cbd3").run(t)
	cliTest(false, false, "machines", "params", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "{}").run(t)
	cliTest(false, false, "machines", "params", "3e7031fe-3062-45f1-835c-92541bc9cbd3").run(t)
	cliTest(false, false, "machines", "params", "3e7031fe-3062-45f1-835c-92541bc9cbd3", machinesParamsNextString).run(t)
	cliTest(false, false, "machines", "params", "3e7031fe-3062-45f1-835c-92541bc9cbd3").run(t)
	cliTest(false, false, "machines", "show", "3e7031fe-3062-45f1-835c-92541bc9cbd3").run(t)
	cliTest(false, false, "machines", "destroy", "3e7031fe-3062-45f1-835c-92541bc9cbd3").run(t)
	cliTest(false, false, "machines", "list").run(t)
	cliTest(false, false, "prefs", "set", "defaultStage", "stage1").run(t)
	cliTest(false, false, "machines", "create", machineCreateInputString).run(t)
	cliTest(false, false, "machines", "destroy", "3e7031fe-3062-45f1-835c-92541bc9cbd3").run(t)
	cliTest(false, false, "machines", "create", "1name").run(t)
	cliTest(false, false, "machines", "destroy", "Name:1name").run(t)
	cliTest(false, false, "machines", "list").run(t)
	cliTest(false, false, "prefs", "set", "defaultStage", "none").run(t)
	cliTest(false, false, "plugins", "destroy", "incr").run(t)
	cliTest(false, false, "stages", "destroy", "stage1").run(t)
	cliTest(false, false, "stages", "destroy", "stage2").run(t)
	cliTest(false, false, "profiles", "destroy", "jill").run(t)
	cliTest(false, false, "profiles", "destroy", "jean").run(t)
	cliTest(false, false, "profiles", "destroy", "stage-prof").run(t)
	cliTest(false, false, "tasks", "destroy", "jamie").run(t)
	cliTest(false, false, "tasks", "destroy", "justine").run(t)
	verifyClean(t)
}

func mta(usage, err bool, tasks ...string) *CliTest {
	args := []string{"machines", "tasks", "add", "3e7031fe-3062-45f1-835c-92541bc9cbd3"}
	args = append(args, tasks...)
	return cliTest(usage, err, args...)
}

func rta(usage, err bool, tasks ...string) *CliTest {
	args := []string{"machines", "tasks", "del", "3e7031fe-3062-45f1-835c-92541bc9cbd3"}
	args = append(args, tasks...)
	return cliTest(usage, err, args...)
}

func fakeJob(t *testing.T, mUuid, state string) {
	t.Helper()
	j := &models.Job{Machine: uuid.Parse(mUuid)}
	lsession, apierr := api.UserSession("https://127.0.0.1:10001", "rocketskates", "r0cketsk8ts")
	if apierr != nil {
		t.Fatalf("Error getting session: %v", apierr)
		return
	}
	defer lsession.Close()
	if err := lsession.CreateModel(j); err != nil {
		t.Errorf("Error creating job :%v", err)
		return
	}
	j.State = state
	if err := lsession.PutModel(j); err != nil {
		t.Errorf("Error updating state to %s: %v", state, err)
		return
	}
	cliTest(false, false, "machines", "show", mUuid).run(t)
}

func TestMachineTaskCli(t *testing.T) {
	mUUID := "3e7031fe-3062-45f1-835c-92541bc9cbd3"
	cliTest(false, false, "machines", "create", machineCreateInputString).run(t)
	tasks := []string{"task1", "task2", "task3", "task4"}
	for _, task := range tasks {
		cliTest(false, false, "tasks", "create", task).run(t)
	}
	// 4 tasks - 1 2 3 4
	mta(false, false, tasks...).run(t)
	// Idempotent add -- 4 tasks -- 1 2 3 4
	mta(false, false, tasks...).run(t)
	// 2 tasks - 1 3
	rta(false, false, "task2", "task4").run(t)
	// 0 tasks
	rta(false, false, "task1", "task3").run(t)
	// 4 tasks - 1 2 3 4
	mta(false, false, tasks...).run(t)
	// 8 tasks - 4 3 2 1 1 2 3 4
	mta(false, false, "at", "0", "task4", "task3", "task2", "task1").run(t)
	// still 8 tasks - 4 3 2 1 1 2 3 4
	mta(false, false, "at", "0", "task4", "task3", "task2", "task1").run(t)
	// 6 tasks - 4 3 2 2 3 4
	rta(false, false, "task1", "task1").run(t)
	fakeJob(t, mUUID, "finished")
	// 7 tasks - 4 3 1 2 2 3 4
	mta(false, false, "at", "1", "task1").run(t)
	// still 7 tasks - 4 3 1 2 2 3 4
	mta(false, false, "at", "1", "task1").run(t)
	// 4 tasks - 4 3 2 4
	rta(false, false, "task1", "task2", "task3").run(t)
	fakeJob(t, mUUID, "finished")
	cliTest(false, false, "machines", "destroy", mUUID).run(t)
	for _, task := range tasks {
		cliTest(false, false, "tasks", "destroy", task).run(t)
	}
	jobs := []*models.Job{}
	lsession, apierr := api.UserSession("https://127.0.0.1:10001", "rocketskates", "r0cketsk8ts")
	if apierr != nil {
		t.Fatalf("Error getting session: %v", apierr)
		return
	}
	defer lsession.Close()
	if err := lsession.Req().List("jobs").Do(&jobs); err == nil {
		for _, j := range jobs {
			lsession.DeleteModel("jobs", j.Uuid.String())
		}
	}
	verifyClean(t)
}

func TestMachineFileImport(t *testing.T) {
	prefix := "machines"
	yamlId := "a2d9b43a-b545-464b-8bc4-088daa7fa7c4"
	jsonId := "b2d9b43a-b545-464b-8bc4-088daa7fa7c4"

	yamlCreate := fmt.Sprintf("test-data/base/%s/create.yaml", prefix)
	jsonCreate := fmt.Sprintf("test-data/base/%s/create.json", prefix)
	yamlBad := fmt.Sprintf("test-data/base/%s/bad.yaml", prefix)
	jsonBad := fmt.Sprintf("test-data/base/%s/bad.json", prefix)
	yamlUpdate := fmt.Sprintf("test-data/base/%s/update.yaml", prefix)
	jsonUpdate := fmt.Sprintf("test-data/base/%s/update.json", prefix)

	cliTest(false, false, prefix, "create", yamlCreate).run(t)
	cliTest(false, false, prefix, "create", jsonCreate).run(t)
	cliTest(false, true, prefix, "create", yamlBad).run(t)
	cliTest(false, true, prefix, "create", jsonBad).run(t)
	cliTest(false, false, prefix, "update", yamlId, yamlUpdate).run(t)
	cliTest(false, false, prefix, "update", jsonId, jsonUpdate).run(t)
	cliTest(false, true, prefix, "update", yamlId, yamlBad).run(t)
	cliTest(false, true, prefix, "update", jsonId, jsonBad).run(t)
	cliTest(false, false, prefix, "destroy", yamlId).run(t)
	cliTest(false, false, prefix, "destroy", jsonId).run(t)
}

func TestMachineArch(t *testing.T) {
	validArches := []string{"amd64", "386", "arm", "arm64", "ppc64", "ppc64le", "mips64", "mips64le", "s390x", "mipsle", "mips"}
	for i, arch := range validArches {
		cliTest(false, false, "machines", "create", "-").Stdin(
			fmt.Sprintf(`---
Name: test-%d
Arch: %s`, i, arch)).run(t)
	}
	for i := range validArches {
		cliTest(false, false, "machines", "destroy", fmt.Sprintf("Name:test-%d", i)).run(t)
	}
	archAliases := []string{
		"x86_64",
		"486", "686", "i386", "i486", "i686",
		"armel", "armhfp",
		"aarch64",
		"power9",
		"mips64el",
		"mipsel",
	}
	for i, arch := range archAliases {
		cliTest(false, true, "machines", "create", "-").Stdin(
			fmt.Sprintf(`---
Name: test-%d
Arch: %s`, i, arch)).run(t)
	}
	for i, arch := range []string{"foo", "bar"} {
		cliTest(false, true, "machines", "create", "-").Stdin(
			fmt.Sprintf(`---
Name: test-%d
Arch: %s`, i, arch)).run(t)
	}
	verifyClean(t)
}

func TestMachineLocked(t *testing.T) {
	cliTest(false, false, "machines", "create", "-").Stdin(machineCreateInputString).run(t)
	cliTest(false, false, "machines", "update", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "-").Stdin(`{"Locked":true}`).run(t)
	cliTest(false, true, "machines", "update", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "-").Stdin(`{"Address":"192.168.124.20"}`).run(t)
	cliTest(false, false, "machines", "update", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "-").Stdin(`{"Locked":false}`).run(t)
	cliTest(false, false, "machines", "update", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "-").Stdin(`{"Address":"192.168.124.20"}`).run(t)
	cliTest(false, false, "machines", "update", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "-").Stdin(`{"Locked":true}`).run(t)
	cliTest(false, false, "machines", "update", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "-").Stdin(`{"Address":"192.168.124.30","Locked":false}`).run(t)
	cliTest(false, false, "machines", "update", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "-").Stdin(`{"Locked":true}`).run(t)
	cliTest(false, true, "machines", "destroy", "3e7031fe-3062-45f1-835c-92541bc9cbd3").run(t)
	cliTest(false, false, "machines", "update", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "-").Stdin(`{"Locked":false}`).run(t)
	cliTest(false, false, "machines", "destroy", "3e7031fe-3062-45f1-835c-92541bc9cbd3").run(t)
	verifyClean(t)
}

func TestMachineProfilesAndParams(t *testing.T) {
	cliTest(false, false, "machines", "create", "bob").run(t)
	cliTest(false, false, "machines", "create", "fred").run(t)
	cliTest(false, false, "machines", "create", "julius").run(t)
	cliTest(false, false, "profiles", "create", "foo").run(t)
	cliTest(false, false, "profiles", "create", "bar").run(t)
	cliTest(false, false, "machines", "update", "Name:bob", "-").Stdin(`{"Profiles":["foo","bar"]}`).run(t)
	cliTest(false, false, "machines", "update", "Name:fred", "-").Stdin(`{"Profiles":["bar"]}`).run(t)
	cliTest(false, false, "machines", "set", "Name:julius", "param", "dog", "to", "bark").run(t)
	cliTest(false, false, "machines", "set", "Name:bob", "param", "dog", "to", "bark").run(t)
	cliTest(false, false, "machines", "set", "Name:fred", "param", "cat", "to", "meow").run(t)
	cliTest(false, false, "machines", "set", "Name:julius", "param", "cat", "to", "meow").run(t)
	cliTest(false, false, "machines", "set", "Name:bob", "param", "bird", "to", "tweet").run(t)
	cliTest(false, false, "machines", "set", "Name:fred", "param", "bird", "to", "tweet").run(t)
	for _, op := range []string{"Eq", "Ne"} {
		for _, profiles := range []string{"foo,bar", "bar"} {
			cliTest(false, false, "machines", "list", "Profiles", op, profiles, "sort", "Name").run(t)
		}
		for _, params := range []string{"dog", "cat", "bird", "dog,cat", "cat,bird", "bird,dog"} {
			cliTest(false, false, "machines", "list", "Params", op, params, "sort", "Name").run(t)
		}
	}
	cliTest(false, false, "machines", "list", "Profiles", "In", "foo,bar", "sort", "Name").run(t)
	cliTest(false, false, "machines", "list", "Profiles", "Nin", "foo,bar", "sort", "Name").run(t)
	cliTest(false, false, "machines", "list", "Name", "In", "fred,bob", "sort", "Name").run(t)
	cliTest(false, false, "machines", "list", "Name", "Nin", "fred,bob", "sort", "Name").run(t)
	cliTest(false, false, "machines", "destroy", "Name:bob").run(t)
	cliTest(false, false, "machines", "destroy", "Name:fred").run(t)
	cliTest(false, false, "machines", "destroy", "Name:julius").run(t)
	cliTest(false, false, "profiles", "destroy", "foo").run(t)
	cliTest(false, false, "profiles", "destroy", "bar").run(t)
	verifyClean(t)
}
