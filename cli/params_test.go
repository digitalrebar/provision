package cli

import (
	"testing"
)

var paramShowNoArgErrorString string = "Error: drpcli params show [id] [flags] requires 1 argument\n"
var paramShowTooManyArgErrorString string = "Error: drpcli params show [id] [flags] requires 1 argument\n"
var paramShowMissingArgErrorString string = "Error: GET: params/john2: Not Found\n\n"
var paramExistsNoArgErrorString string = "Error: drpcli params exists [id] [flags] requires 1 argument"
var paramExistsTooManyArgErrorString string = "Error: drpcli params exists [id] [flags] requires 1 argument"
var paramExistsMissingJohnString string = "Error: GET: params/john2: Not Found\n\n"
var paramCreateNoArgErrorString string = "Error: drpcli params create [json] [flags] requires 1 argument\n"
var paramCreateTooManyArgErrorString string = "Error: drpcli params create [json] [flags] requires 1 argument\n"
var paramCreateBadJSONErrorString = "Error: Invalid param object: error converting YAML to JSON: yaml: line 1: did not find expected ',' or '}' and error converting YAML to JSON: yaml: line 1: did not find expected ',' or '}'\n\n"
var paramCreateBadJSON2ErrorString = "Error: Unable to create new param: Invalid type passed to param create\n\n"
var paramCreateDuplicateErrorString = "Error: CREATE: params/john: already exists\n\n"
var paramUpdateNoArgErrorString string = "Error: drpcli params update [id] [json] [flags] requires 2 arguments"
var paramUpdateTooManyArgErrorString string = "Error: drpcli params update [id] [json] [flags] requires 2 arguments"
var paramUpdateBadJSONErrorString = "Error: Unable to merge objects: json: cannot unmarshal string into Go value of type map[string]interface {}\n\n\n"
var paramUpdateJohnMissingErrorString string = "Error: GET: params/john2: Not Found\n\n"
var paramPatchNoArgErrorString string = "Error: drpcli params patch [objectJson] [changesJson] [flags] requires 2 arguments"
var paramPatchTooManyArgErrorString string = "Error: drpcli params patch [objectJson] [changesJson] [flags] requires 2 arguments"
var paramPatchBadPatchJSONErrorString = "Error: Unable to parse drpcli params patch [objectJson] [changesJson] [flags] JSON asdgasdg\nError: error unmarshaling JSON: json: cannot unmarshal string into Go value of type genmodels.Param\n\n"
var paramPatchBadBaseJSONErrorString = "Error: Unable to parse drpcli params patch [objectJson] [changesJson] [flags] JSON asdgasdg\nError: error unmarshaling JSON: json: cannot unmarshal string into Go value of type genmodels.Param\n\n"
var paramPatchJohnMissingErrorString string = "Error: PATCH: params/john2: Not Found\n\n"
var paramDestroyNoArgErrorString string = "Error: drpcli params destroy [id] [flags] requires 1 argument"
var paramDestroyTooManyArgErrorString string = "Error: drpcli params destroy [id] [flags] requires 1 argument"
var paramDestroyMissingJohnString string = "Error: DELETE: params/john: Not Found\n\n"
var paramBootEnvNoArgErrorString string = "Error: drpcli params bootenv [id] [bootenv] [flags] requires 2 arguments"
var paramBootEnvMissingParamErrorString string = "Error: params GET: john: Not Found\n\n"
var paramGetNoArgErrorString string = "Error: drpcli params get [id] param [key] [flags] requires 3 arguments"
var paramGetMissingParamErrorString string = "Error: params GET Params: john2: Not Found\n\n"

var paramDefaultListString string = `[
  {
    "Available": true,
    "Errors": [],
    "Name": "incrementer/parameter",
    "ReadOnly": true,
    "Schema": {
      "type": "string"
    },
    "Validated": true
  },
  {
    "Available": true,
    "Errors": [],
    "Name": "incrementer/step",
    "ReadOnly": true,
    "Schema": {
      "type": "integer"
    },
    "Validated": true
  },
  {
    "Available": true,
    "Errors": [],
    "Name": "incrementer/touched",
    "ReadOnly": true,
    "Schema": {
      "type": "integer"
    },
    "Validated": true
  }
]
`

var paramEmptyListString string = "[]\n"

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

var paramExistsParamString string = ""
var paramCreateBadJSONString = "{asdgasdg"
var paramCreateBadJSON2String = "[asdgasdg]"
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

var paramListParamsString = `[
  {
    "Available": true,
    "Errors": [],
    "Name": "incrementer/parameter",
    "ReadOnly": true,
    "Schema": {
      "type": "string"
    },
    "Validated": true
  },
  {
    "Available": true,
    "Errors": [],
    "Name": "incrementer/step",
    "ReadOnly": true,
    "Schema": {
      "type": "integer"
    },
    "Validated": true
  },
  {
    "Available": true,
    "Errors": [],
    "Name": "incrementer/touched",
    "ReadOnly": true,
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

var paramUpdateBadJSONString = "asdgasdg"

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

var paramPatchBadPatchJSONString = "asdgasdg"

var paramPatchBadBaseJSONString = "asdgasdg"

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

var paramDestroyJohnString string = "Deleted param john\n"

func TestParamCli(t *testing.T) {
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
