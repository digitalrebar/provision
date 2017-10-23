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

	tests := []CliTest{
		CliTest{true, false, []string{"stages"}, noStdinString, "Access CLI commands relating to stages\n", ""},
		CliTest{false, false, []string{"stages", "list"}, noStdinString, stageDefaultListString, noErrorString},

		CliTest{true, true, []string{"stages", "create"}, noStdinString, noContentString, stageCreateNoArgErrorString},
		CliTest{true, true, []string{"stages", "create", "john", "john2"}, noStdinString, noContentString, stageCreateTooManyArgErrorString},
		CliTest{false, true, []string{"stages", "create", stageCreateBadJSONString}, noStdinString, noContentString, stageCreateBadJSONErrorString},
		CliTest{false, true, []string{"stages", "create", stageCreateBadJSON2String}, noStdinString, noContentString, stageCreateBadJSON2ErrorString},
		CliTest{false, false, []string{"stages", "create", stageCreateInputString}, noStdinString, stageCreateJohnString, noErrorString},
		CliTest{false, true, []string{"stages", "create", stageCreateInputString}, noStdinString, noContentString, stageCreateDuplicateErrorString},
		CliTest{false, false, []string{"stages", "list"}, noStdinString, stageListStagesString, noErrorString},
		CliTest{false, false, []string{"stages", "list", "Name=fred"}, noStdinString, stageEmptyListString, noErrorString},
		CliTest{false, false, []string{"stages", "list", "Name=john"}, noStdinString, stageListJohnOnlyString, noErrorString},
		CliTest{false, false, []string{"stages", "list", "BootEnv=fred"}, noStdinString, stageEmptyListString, noErrorString},
		CliTest{false, false, []string{"stages", "list", "BootEnv=local"}, noStdinString, stageListJohnOnlyString, noErrorString},
		CliTest{false, false, []string{"stages", "list", "Reboot=true"}, noStdinString, stageEmptyListString, noErrorString},
		CliTest{false, false, []string{"stages", "list", "Reboot=false"}, noStdinString, stageListStagesString, noErrorString},
		CliTest{false, true, []string{"stages", "list", "Reboot=fred"}, noStdinString, noContentString, "Error: GET: stages: Reboot must be true or false\n\n"},
		CliTest{true, true, []string{"stages", "show"}, noStdinString, noContentString, stageShowNoArgErrorString},
		CliTest{true, true, []string{"stages", "show", "john", "john2"}, noStdinString, noContentString, stageShowTooManyArgErrorString},
		CliTest{false, true, []string{"stages", "show", "john2"}, noStdinString, noContentString, stageShowMissingArgErrorString},
		CliTest{false, false, []string{"stages", "show", "john"}, noStdinString, stageShowStageString, noErrorString},

		CliTest{true, true, []string{"stages", "exists"}, noStdinString, noContentString, stageExistsNoArgErrorString},
		CliTest{true, true, []string{"stages", "exists", "john", "john2"}, noStdinString, noContentString, stageExistsTooManyArgErrorString},
		CliTest{false, true, []string{"stages", "exists", "john2"}, noStdinString, noContentString, stageExistsMissingJohnString},
		CliTest{true, true, []string{"stages", "exists", "john", "john2"}, noStdinString, noContentString, stageExistsTooManyArgErrorString},
		CliTest{false, false, []string{"stages", "exists", "john"}, noStdinString, stageExistsStageString, noErrorString},

		CliTest{true, true, []string{"stages", "update"}, noStdinString, noContentString, stageUpdateNoArgErrorString},
		CliTest{true, true, []string{"stages", "update", "john", "john2", "john3"}, noStdinString, noContentString, stageUpdateTooManyArgErrorString},
		CliTest{false, true, []string{"stages", "update", "john", stageUpdateBadJSONString}, noStdinString, noContentString, stageUpdateBadJSONErrorString},
		CliTest{false, false, []string{"stages", "update", "john", stageUpdateInputString}, noStdinString, stageUpdateJohnString, noErrorString},
		CliTest{false, true, []string{"stages", "update", "john2", stageUpdateInputString}, noStdinString, noContentString, stageUpdateJohnMissingErrorString},
		CliTest{false, false, []string{"stages", "show", "john"}, noStdinString, stageUpdateJohnString, noErrorString},

		CliTest{true, true, []string{"stages", "patch"}, noStdinString, noContentString, stagePatchNoArgErrorString},
		CliTest{true, true, []string{"stages", "patch", "john", "john2", "john3"}, noStdinString, noContentString, stagePatchTooManyArgErrorString},
		CliTest{false, true, []string{"stages", "patch", stagePatchBaseString, stagePatchBadPatchJSONString}, noStdinString, noContentString, stagePatchBadPatchJSONErrorString},
		CliTest{false, true, []string{"stages", "patch", stagePatchBadBaseJSONString, stagePatchInputString}, noStdinString, noContentString, stagePatchBadBaseJSONErrorString},
		CliTest{false, false, []string{"stages", "patch", stagePatchBaseString, stagePatchInputString}, noStdinString, stagePatchJohnString, noErrorString},
		CliTest{false, true, []string{"stages", "patch", stagePatchMissingBaseString, stagePatchInputString}, noStdinString, noContentString, stagePatchJohnMissingErrorString},
		CliTest{false, false, []string{"stages", "show", "john"}, noStdinString, stagePatchJohnString, noErrorString},

		CliTest{true, true, []string{"stages", "destroy"}, noStdinString, noContentString, stageDestroyNoArgErrorString},
		CliTest{true, true, []string{"stages", "destroy", "john", "june"}, noStdinString, noContentString, stageDestroyTooManyArgErrorString},
		CliTest{false, false, []string{"stages", "destroy", "john"}, noStdinString, stageDestroyJohnString, noErrorString},
		CliTest{false, true, []string{"stages", "destroy", "john"}, noStdinString, noContentString, stageDestroyMissingJohnString},
		CliTest{false, false, []string{"stages", "list"}, noStdinString, stageDefaultListString, noErrorString},

		CliTest{false, false, []string{"stages", "create", "-"}, stageCreateInputString + "\n", stageCreateJohnString, noErrorString},
		CliTest{false, false, []string{"stages", "list"}, noStdinString, stageListStagesString, noErrorString},
		CliTest{false, false, []string{"stages", "update", "john", "-"}, stageUpdateInputString + "\n", stageUpdateJohnString, noErrorString},
		CliTest{false, false, []string{"stages", "show", "john"}, noStdinString, stageUpdateJohnString, noErrorString},

		CliTest{false, false, []string{"stages", "destroy", "john"}, noStdinString, stageDestroyJohnString, noErrorString},
		CliTest{false, false, []string{"stages", "list"}, noStdinString, stageDefaultListString, noErrorString},
	}

	for _, test := range tests {
		testCli(t, test)
	}

}
