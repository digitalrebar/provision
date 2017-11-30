package cli

import (
	"testing"
)

var stageShowMissingArgErrorString string = "Error: GET: stages/john2: Not Found\n\n"
var stageExistsMissingJohnString string = "Error: GET: stages/john2: Not Found\n\n"
var stageCreateDuplicateErrorString = "Error: CREATE: stages/john: already exists\n\n"
var stageUpdateJohnMissingErrorString string = "Error: GET: stages/john2: Not Found\n\n"
var stagePatchJohnMissingErrorString string = "Error: PATCH: stages/john2: Not Found\n\n"
var stageDestroyMissingJohnString string = "Error: DELETE: stages/john: Not Found\n\n"
var stageBootEnvMissingStageErrorString string = "Error: stages GET: john: Not Found\n\n"

var stageDefaultListString string = `[
  {
    "Available": true,
    "BootEnv": "",
    "Description": "Stage to boot into the local BootEnv.",
    "Errors": [],
    "Meta": {
      "color": "green",
      "icon": "radio",
      "title": "Digital Rebar Provision"
    },
    "Name": "local",
    "OptionalParams": [],
    "Profiles": [],
    "ReadOnly": true,
    "RequiredParams": [],
    "Tasks": [],
    "Templates": [],
    "Validated": true
  },
  {
    "Available": true,
    "BootEnv": "",
    "Description": "Noop / Nothing stage",
    "Errors": [],
    "Meta": {
      "color": "green",
      "icon": "circle thin",
      "title": "Digital Rebar Provision"
    },
    "Name": "none",
    "OptionalParams": [],
    "Profiles": [],
    "ReadOnly": true,
    "RequiredParams": [],
    "Tasks": [],
    "Templates": [],
    "Validated": true
  }
]
`
var stageEmptyListString string = "[]\n"

var stageShowNoArgErrorString string = "Error: drpcli stages show [id] [flags] requires 1 argument\n"
var stageShowTooManyArgErrorString string = "Error: drpcli stages show [id] [flags] requires 1 argument\n"

var stageShowStageString string = `{
  "Available": true,
  "BootEnv": "local",
  "Errors": [],
  "Name": "john",
  "OptionalParams": [],
  "Profiles": [],
  "ReadOnly": false,
  "RequiredParams": [],
  "Tasks": [],
  "Templates": [],
  "Validated": true
}
`

var stageExistsNoArgErrorString string = "Error: drpcli stages exists [id] [flags] requires 1 argument"
var stageExistsTooManyArgErrorString string = "Error: drpcli stages exists [id] [flags] requires 1 argument"
var stageExistsStageString string = ""

var stageCreateNoArgErrorString string = "Error: drpcli stages create [json] [flags] requires 1 argument\n"
var stageCreateTooManyArgErrorString string = "Error: drpcli stages create [json] [flags] requires 1 argument\n"
var stageCreateBadJSONString = "{asdgasdg"
var stageCreateBadJSONErrorString = "Error: Invalid stage object: error converting YAML to JSON: yaml: line 1: did not find expected ',' or '}' and error converting YAML to JSON: yaml: line 1: did not find expected ',' or '}'\n\n"
var stageCreateBadJSON2String = "[asdgasdg]"
var stageCreateBadJSON2ErrorString = "Error: Unable to create new stage: Invalid type passed to stage create\n\n"
var stageCreateInputString string = `{
  "Name": "john",
  "BootEnv": "local"
}
`
var stageCreateJohnString string = `{
  "Available": true,
  "BootEnv": "local",
  "Errors": [],
  "Name": "john",
  "OptionalParams": [],
  "Profiles": [],
  "ReadOnly": false,
  "RequiredParams": [],
  "Tasks": [],
  "Templates": [],
  "Validated": true
}
`

var stageListStagesString = `[
  {
    "Available": true,
    "BootEnv": "local",
    "Errors": [],
    "Name": "john",
    "OptionalParams": [],
    "Profiles": [],
    "ReadOnly": false,
    "RequiredParams": [],
    "Tasks": [],
    "Templates": [],
    "Validated": true
  },
  {
    "Available": true,
    "BootEnv": "",
    "Description": "Stage to boot into the local BootEnv.",
    "Errors": [],
    "Meta": {
      "color": "green",
      "icon": "radio",
      "title": "Digital Rebar Provision"
    },
    "Name": "local",
    "OptionalParams": [],
    "Profiles": [],
    "ReadOnly": true,
    "RequiredParams": [],
    "Tasks": [],
    "Templates": [],
    "Validated": true
  },
  {
    "Available": true,
    "BootEnv": "",
    "Description": "Noop / Nothing stage",
    "Errors": [],
    "Meta": {
      "color": "green",
      "icon": "circle thin",
      "title": "Digital Rebar Provision"
    },
    "Name": "none",
    "OptionalParams": [],
    "Profiles": [],
    "ReadOnly": true,
    "RequiredParams": [],
    "Tasks": [],
    "Templates": [],
    "Validated": true
  }
]
`
var stageListJohnOnlyString = `[
  {
    "Available": true,
    "BootEnv": "local",
    "Errors": [],
    "Name": "john",
    "OptionalParams": [],
    "Profiles": [],
    "ReadOnly": false,
    "RequiredParams": [],
    "Tasks": [],
    "Templates": [],
    "Validated": true
  }
]
`

var stageUpdateNoArgErrorString string = "Error: drpcli stages update [id] [json] [flags] requires 2 arguments"
var stageUpdateTooManyArgErrorString string = "Error: drpcli stages update [id] [json] [flags] requires 2 arguments"
var stageUpdateBadJSONString = "asdgasdg"
var stageUpdateBadJSONErrorString = "Error: Unable to merge objects: json: cannot unmarshal string into Go value of type map[string]interface {}\n\n\n"
var stageUpdateInputString string = `{
  "Description": "Awesome sauce"
}
`
var stageUpdateJohnString string = `{
  "Available": true,
  "BootEnv": "local",
  "Description": "Awesome sauce",
  "Errors": [],
  "Name": "john",
  "OptionalParams": [],
  "Profiles": [],
  "ReadOnly": false,
  "RequiredParams": [],
  "Tasks": [],
  "Templates": [],
  "Validated": true
}
`

var stagePatchNoArgErrorString string = "Error: drpcli stages patch [objectJson] [changesJson] [flags] requires 2 arguments"
var stagePatchTooManyArgErrorString string = "Error: drpcli stages patch [objectJson] [changesJson] [flags] requires 2 arguments"
var stagePatchBadPatchJSONString = "asdgasdg"
var stagePatchBadPatchJSONErrorString = "Error: Unable to parse drpcli stages patch [objectJson] [changesJson] [flags] JSON asdgasdg\nError: error unmarshaling JSON: json: cannot unmarshal string into Go value of type genmodels.Stage\n\n"
var stagePatchBadBaseJSONString = "asdgasdg"
var stagePatchBadBaseJSONErrorString = "Error: Unable to parse drpcli stages patch [objectJson] [changesJson] [flags] JSON asdgasdg\nError: error unmarshaling JSON: json: cannot unmarshal string into Go value of type genmodels.Stage\n\n"
var stagePatchBaseString string = `{
  "Name": "john",
  "BootEnv": "local"
}
`
var stagePatchInputString string = `{
  "Description": "No Really Awesome Sauce"
}
`
var stagePatchJohnString string = `{
  "Available": true,
  "BootEnv": "local",
  "Description": "No Really Awesome Sauce",
  "Errors": [],
  "Name": "john",
  "OptionalParams": [],
  "Profiles": [],
  "ReadOnly": false,
  "RequiredParams": [],
  "Tasks": [],
  "Templates": [],
  "Validated": true
}
`
var stagePatchMissingBaseString string = `{
  "Name": "john2",
  "Params": {
    "Name": ""
  }
}
`

var stageDestroyNoArgErrorString string = "Error: drpcli stages destroy [id] [flags] requires 1 argument"
var stageDestroyTooManyArgErrorString string = "Error: drpcli stages destroy [id] [flags] requires 1 argument"
var stageDestroyJohnString string = "Deleted stage john\n"

var stageBootEnvNoArgErrorString string = "Error: drpcli stages bootenv [id] [bootenv] [flags] requires 2 arguments"

func TestStageCli(t *testing.T) {
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
