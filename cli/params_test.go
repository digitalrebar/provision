package cli

import (
	"testing"
)

var paramDefaultListString string = `[
  {
    "Available": true,
    "Errors": [],
    "Name": "incrementer/parameter",
    "ReadOnly": false,
    "Schema": {
      "type": "string"
    },
    "Validated": true
  },
  {
    "Available": true,
    "Errors": [],
    "Name": "incrementer/step",
    "ReadOnly": false,
    "Schema": {
      "type": "integer"
    },
    "Validated": true
  },
  {
    "Available": true,
    "Errors": [],
    "Name": "incrementer/touched",
    "ReadOnly": false,
    "Schema": {
      "type": "integer"
    },
    "Validated": true
  }
]
`

var paramEmptyListString string = "[]\n"

var paramShowNoArgErrorString string = "Error: drpcli params show [id] [flags] requires 1 argument\n"
var paramShowTooManyArgErrorString string = "Error: drpcli params show [id] [flags] requires 1 argument\n"
var paramShowMissingArgErrorString string = "Error: params GET: john2: Not Found\n\n"
var paramShowParamString string = `{
  "Available": true,
  "Errors": [],
  "Name": "john",
  "ReadOnly": false,
  "Schema": {
    "type": "string"
  },
  "Validated": true
}
`

var paramExistsNoArgErrorString string = "Error: drpcli params exists [id] [flags] requires 1 argument"
var paramExistsTooManyArgErrorString string = "Error: drpcli params exists [id] [flags] requires 1 argument"
var paramExistsParamString string = ""
var paramExistsMissingJohnString string = "Error: params GET: john2: Not Found\n\n"

var paramCreateNoArgErrorString string = "Error: drpcli params create [json] [flags] requires 1 argument\n"
var paramCreateTooManyArgErrorString string = "Error: drpcli params create [json] [flags] requires 1 argument\n"
var paramCreateBadJSONString = "{asdgasdg"
var paramCreateBadJSONErrorString = "Error: Invalid param object: error converting YAML to JSON: yaml: line 1: did not find expected ',' or '}' and error converting YAML to JSON: yaml: line 1: did not find expected ',' or '}'\n\n"
var paramCreateBadJSON2String = "[asdgasdg]"
var paramCreateBadJSON2ErrorString = "Error: Unable to create new param: Invalid type passed to param create\n\n"
var paramCreateInputString string = `{
  "Name": "john",
  "Schema": {
    "type": "string"
  }
}
`
var paramCreateJohnString string = `{
  "Available": true,
  "Errors": [],
  "Name": "john",
  "ReadOnly": false,
  "Schema": {
    "type": "string"
  },
  "Validated": true
}
`
var paramCreateDuplicateErrorString = "Error: dataTracker create params: john already exists\n\n"

var paramListParamsString = `[
  {
    "Available": true,
    "Errors": [],
    "Name": "incrementer/parameter",
    "ReadOnly": false,
    "Schema": {
      "type": "string"
    },
    "Validated": true
  },
  {
    "Available": true,
    "Errors": [],
    "Name": "incrementer/step",
    "ReadOnly": false,
    "Schema": {
      "type": "integer"
    },
    "Validated": true
  },
  {
    "Available": true,
    "Errors": [],
    "Name": "incrementer/touched",
    "ReadOnly": false,
    "Schema": {
      "type": "integer"
    },
    "Validated": true
  },
  {
    "Available": true,
    "Errors": [],
    "Name": "john",
    "ReadOnly": false,
    "Schema": {
      "type": "string"
    },
    "Validated": true
  }
]
`
var paramListJohnOnlyString = `[
  {
    "Available": true,
    "Errors": [],
    "Name": "john",
    "ReadOnly": false,
    "Schema": {
      "type": "string"
    },
    "Validated": true
  }
]
`

var paramUpdateNoArgErrorString string = "Error: drpcli params update [id] [json] [flags] requires 2 arguments"
var paramUpdateTooManyArgErrorString string = "Error: drpcli params update [id] [json] [flags] requires 2 arguments"
var paramUpdateBadJSONString = "asdgasdg"
var paramUpdateBadJSONErrorString = "Error: Unable to merge objects: json: cannot unmarshal string into Go value of type map[string]interface {}\n\n\n"
var paramUpdateInputString string = `{
  "Schema": {
    "type": "string"
  }
}
`
var paramUpdateJohnString string = `{
  "Available": true,
  "Errors": [],
  "Name": "john",
  "ReadOnly": false,
  "Schema": {
    "type": "string"
  },
  "Validated": true
}
`
var paramUpdateJohnMissingErrorString string = "Error: params GET: john2: Not Found\n\n"

var paramPatchNoArgErrorString string = "Error: drpcli params patch [objectJson] [changesJson] [flags] requires 2 arguments"
var paramPatchTooManyArgErrorString string = "Error: drpcli params patch [objectJson] [changesJson] [flags] requires 2 arguments"
var paramPatchBadPatchJSONString = "asdgasdg"
var paramPatchBadPatchJSONErrorString = "Error: Unable to parse drpcli params patch [objectJson] [changesJson] [flags] JSON asdgasdg\nError: error unmarshaling JSON: json: cannot unmarshal string into Go value of type genmodels.Param\n\n"
var paramPatchBadBaseJSONString = "asdgasdg"
var paramPatchBadBaseJSONErrorString = "Error: Unable to parse drpcli params patch [objectJson] [changesJson] [flags] JSON asdgasdg\nError: error unmarshaling JSON: json: cannot unmarshal string into Go value of type genmodels.Param\n\n"
var paramPatchBaseString string = `{
  "Available": true,
  "Errors": [],
  "Name": "john",
  "Schema": {
    "type": "string"
  },
  "Validated": true
}
`
var paramPatchInputString string = `{
  "Description": "Foo"
}
`
var paramPatchJohnString string = `{
  "Available": true,
  "Description": "Foo",
  "Errors": [],
  "Name": "john",
  "ReadOnly": false,
  "Schema": {
    "type": "string"
  },
  "Validated": true
}
`
var paramPatchMissingBaseString string = `{
  "Name": "john2",
  "Schema": {}
}
`
var paramPatchJohnMissingErrorString string = "Error: params: PATCH john2: Not Found\n\n"

var paramDestroyNoArgErrorString string = "Error: drpcli params destroy [id] [flags] requires 1 argument"
var paramDestroyTooManyArgErrorString string = "Error: drpcli params destroy [id] [flags] requires 1 argument"
var paramDestroyJohnString string = "Deleted param john\n"
var paramDestroyMissingJohnString string = "Error: params: DELETE john: Not Found\n\n"

var paramBootEnvNoArgErrorString string = "Error: drpcli params bootenv [id] [bootenv] [flags] requires 2 arguments"
var paramBootEnvMissingParamErrorString string = "Error: params GET: john: Not Found\n\n"

var paramGetNoArgErrorString string = "Error: drpcli params get [id] param [key] [flags] requires 3 arguments"
var paramGetMissingParamErrorString string = "Error: params GET Params: john2: Not Found\n\n"

func TestParamCli(t *testing.T) {

	tests := []CliTest{
		CliTest{true, false, []string{"params"}, noStdinString, "Access CLI commands relating to params\n", ""},
		CliTest{false, false, []string{"params", "list"}, noStdinString, paramDefaultListString, noErrorString},

		CliTest{true, true, []string{"params", "create"}, noStdinString, noContentString, paramCreateNoArgErrorString},
		CliTest{true, true, []string{"params", "create", "john", "john2"}, noStdinString, noContentString, paramCreateTooManyArgErrorString},
		CliTest{false, true, []string{"params", "create", paramCreateBadJSONString}, noStdinString, noContentString, paramCreateBadJSONErrorString},
		CliTest{false, true, []string{"params", "create", paramCreateBadJSON2String}, noStdinString, noContentString, paramCreateBadJSON2ErrorString},
		CliTest{false, false, []string{"params", "create", paramCreateInputString}, noStdinString, paramCreateJohnString, noErrorString},
		CliTest{false, true, []string{"params", "create", paramCreateInputString}, noStdinString, noContentString, paramCreateDuplicateErrorString},
		CliTest{false, false, []string{"params", "list"}, noStdinString, paramListParamsString, noErrorString},
		CliTest{false, false, []string{"params", "list", "--limit=0"}, noStdinString, paramEmptyListString, noErrorString},
		CliTest{false, false, []string{"params", "list", "--limit=10", "--offset=0"}, noStdinString, paramListParamsString, noErrorString},
		CliTest{false, false, []string{"params", "list", "--limit=10", "--offset=10"}, noStdinString, paramEmptyListString, noErrorString},
		CliTest{false, true, []string{"params", "list", "--limit=-10", "--offset=0"}, noStdinString, noContentString, limitNegativeError},
		CliTest{false, true, []string{"params", "list", "--limit=10", "--offset=-10"}, noStdinString, noContentString, offsetNegativeError},
		CliTest{false, false, []string{"params", "list", "--limit=-1", "--offset=-1"}, noStdinString, paramListParamsString, noErrorString},
		CliTest{false, false, []string{"params", "list", "Name=fred"}, noStdinString, paramEmptyListString, noErrorString},
		CliTest{false, false, []string{"params", "list", "Name=john"}, noStdinString, paramListJohnOnlyString, noErrorString},
		CliTest{false, false, []string{"params", "list", "Available=true"}, noStdinString, paramListParamsString, noErrorString},
		CliTest{false, false, []string{"params", "list", "Available=false"}, noStdinString, paramEmptyListString, noErrorString},
		CliTest{false, true, []string{"params", "list", "Available=fred"}, noStdinString, noContentString, bootEnvBadAvailableString},
		CliTest{false, false, []string{"params", "list", "Valid=true"}, noStdinString, paramListParamsString, noErrorString},
		CliTest{false, false, []string{"params", "list", "Valid=false"}, noStdinString, paramEmptyListString, noErrorString},
		CliTest{false, true, []string{"params", "list", "Valid=fred"}, noStdinString, noContentString, bootEnvBadValidString},
		CliTest{false, false, []string{"params", "list", "ReadOnly=true"}, noStdinString, paramEmptyListString, noErrorString},
		CliTest{false, false, []string{"params", "list", "ReadOnly=false"}, noStdinString, paramListParamsString, noErrorString},
		CliTest{false, true, []string{"params", "list", "ReadOnly=fred"}, noStdinString, noContentString, bootEnvBadReadOnlyString},

		CliTest{true, true, []string{"params", "show"}, noStdinString, noContentString, paramShowNoArgErrorString},
		CliTest{true, true, []string{"params", "show", "john", "john2"}, noStdinString, noContentString, paramShowTooManyArgErrorString},
		CliTest{false, true, []string{"params", "show", "john2"}, noStdinString, noContentString, paramShowMissingArgErrorString},
		CliTest{false, false, []string{"params", "show", "john"}, noStdinString, paramShowParamString, noErrorString},

		CliTest{true, true, []string{"params", "exists"}, noStdinString, noContentString, paramExistsNoArgErrorString},
		CliTest{true, true, []string{"params", "exists", "john", "john2"}, noStdinString, noContentString, paramExistsTooManyArgErrorString},
		CliTest{false, false, []string{"params", "exists", "john"}, noStdinString, paramExistsParamString, noErrorString},
		CliTest{false, true, []string{"params", "exists", "john2"}, noStdinString, noContentString, paramExistsMissingJohnString},
		CliTest{true, true, []string{"params", "exists", "john", "john2"}, noStdinString, noContentString, paramExistsTooManyArgErrorString},

		CliTest{true, true, []string{"params", "update"}, noStdinString, noContentString, paramUpdateNoArgErrorString},
		CliTest{true, true, []string{"params", "update", "john", "john2", "john3"}, noStdinString, noContentString, paramUpdateTooManyArgErrorString},
		CliTest{false, true, []string{"params", "update", "john", paramUpdateBadJSONString}, noStdinString, noContentString, paramUpdateBadJSONErrorString},
		CliTest{false, false, []string{"params", "update", "john", paramUpdateInputString}, noStdinString, paramUpdateJohnString, noErrorString},
		CliTest{false, true, []string{"params", "update", "john2", paramUpdateInputString}, noStdinString, noContentString, paramUpdateJohnMissingErrorString},
		CliTest{false, false, []string{"params", "show", "john"}, noStdinString, paramUpdateJohnString, noErrorString},

		CliTest{true, true, []string{"params", "patch"}, noStdinString, noContentString, paramPatchNoArgErrorString},
		CliTest{true, true, []string{"params", "patch", "john", "john2", "john3"}, noStdinString, noContentString, paramPatchTooManyArgErrorString},
		CliTest{false, true, []string{"params", "patch", paramPatchBaseString, paramPatchBadPatchJSONString}, noStdinString, noContentString, paramPatchBadPatchJSONErrorString},
		CliTest{false, true, []string{"params", "patch", paramPatchBadBaseJSONString, paramPatchInputString}, noStdinString, noContentString, paramPatchBadBaseJSONErrorString},
		CliTest{false, false, []string{"params", "patch", paramPatchBaseString, paramPatchInputString}, noStdinString, paramPatchJohnString, noErrorString},
		CliTest{false, true, []string{"params", "patch", paramPatchMissingBaseString, paramPatchInputString}, noStdinString, noContentString, paramPatchJohnMissingErrorString},
		CliTest{false, false, []string{"params", "show", "john"}, noStdinString, paramPatchJohnString, noErrorString},

		CliTest{true, true, []string{"params", "destroy"}, noStdinString, noContentString, paramDestroyNoArgErrorString},
		CliTest{true, true, []string{"params", "destroy", "john", "june"}, noStdinString, noContentString, paramDestroyTooManyArgErrorString},
		CliTest{false, false, []string{"params", "destroy", "john"}, noStdinString, paramDestroyJohnString, noErrorString},
		CliTest{false, true, []string{"params", "destroy", "john"}, noStdinString, noContentString, paramDestroyMissingJohnString},
		CliTest{false, false, []string{"params", "list"}, noStdinString, paramDefaultListString, noErrorString},

		CliTest{false, false, []string{"params", "create", "-"}, paramCreateInputString + "\n", paramCreateJohnString, noErrorString},
		CliTest{false, false, []string{"params", "list"}, noStdinString, paramListParamsString, noErrorString},
		CliTest{false, false, []string{"params", "update", "john", "-"}, paramUpdateInputString + "\n", paramUpdateJohnString, noErrorString},
		CliTest{false, false, []string{"params", "show", "john"}, noStdinString, paramUpdateJohnString, noErrorString},

		CliTest{false, false, []string{"params", "destroy", "john"}, noStdinString, paramDestroyJohnString, noErrorString},
		CliTest{false, false, []string{"params", "list"}, noStdinString, paramDefaultListString, noErrorString},
	}

	for _, test := range tests {
		testCli(t, test)
	}

}
