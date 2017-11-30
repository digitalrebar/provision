package cli

import (
	"testing"
)

var profileShowNoArgErrorString string = "Error: drpcli profiles show [id] [flags] requires 1 argument\n"
var profileShowTooManyArgErrorString string = "Error: drpcli profiles show [id] [flags] requires 1 argument\n"
var profileShowMissingArgErrorString string = "Error: GET: profiles/john2: Not Found\n\n"
var profileExistsNoArgErrorString string = "Error: drpcli profiles exists [id] [flags] requires 1 argument"
var profileExistsTooManyArgErrorString string = "Error: drpcli profiles exists [id] [flags] requires 1 argument"
var profileExistsMissingJohnString string = "Error: GET: profiles/john2: Not Found\n\n"
var profileCreateNoArgErrorString string = "Error: drpcli profiles create [json] [flags] requires 1 argument\n"
var profileCreateTooManyArgErrorString string = "Error: drpcli profiles create [json] [flags] requires 1 argument\n"
var profileCreateBadJSONErrorString = "Error: Invalid profile object: error converting YAML to JSON: yaml: line 1: did not find expected ',' or '}' and error converting YAML to JSON: yaml: line 1: did not find expected ',' or '}'\n\n"
var profileCreateBadJSON2ErrorString = "Error: Unable to create new profile: Invalid type passed to profile create\n\n"
var profileCreateDuplicateErrorString = "Error: CREATE: profiles/john: already exists\n\n"
var profileUpdateNoArgErrorString string = "Error: drpcli profiles update [id] [json] [flags] requires 2 arguments"
var profileUpdateTooManyArgErrorString string = "Error: drpcli profiles update [id] [json] [flags] requires 2 arguments"
var profileUpdateBadJSONErrorString = "Error: Unable to merge objects: json: cannot unmarshal string into Go value of type map[string]interface {}\n\n\n"
var profileUpdateJohnMissingErrorString string = "Error: GET: profiles/john2: Not Found\n\n"
var profilePatchNoArgErrorString string = "Error: drpcli profiles patch [objectJson] [changesJson] [flags] requires 2 arguments"
var profilePatchTooManyArgErrorString string = "Error: drpcli profiles patch [objectJson] [changesJson] [flags] requires 2 arguments"
var profilePatchBadPatchJSONErrorString = "Error: Unable to parse drpcli profiles patch [objectJson] [changesJson] [flags] JSON asdgasdg\nError: error unmarshaling JSON: json: cannot unmarshal string into Go value of type genmodels.Profile\n\n"
var profilePatchBadBaseJSONErrorString = "Error: Unable to parse drpcli profiles patch [objectJson] [changesJson] [flags] JSON asdgasdg\nError: error unmarshaling JSON: json: cannot unmarshal string into Go value of type genmodels.Profile\n\n"
var profilePatchJohnMissingErrorString string = "Error: PATCH: profiles/john2: Not Found\n\n"
var profileDestroyNoArgErrorString string = "Error: drpcli profiles destroy [id] [flags] requires 1 argument"
var profileDestroyTooManyArgErrorString string = "Error: drpcli profiles destroy [id] [flags] requires 1 argument"
var profileDestroyMissingJohnString string = "Error: DELETE: profiles/john: Not Found\n\n"
var profileBootEnvNoArgErrorString string = "Error: drpcli profiles bootenv [id] [bootenv] [flags] requires 2 arguments"
var profileBootEnvMissingProfileErrorString string = "Error: profiles GET: john: Not Found\n\n"
var profileGetNoArgErrorString string = "Error: drpcli profiles get [id] param [key] [flags] requires 3 arguments"
var profileGetMissingProfileErrorString string = "Error: GET: profiles/john2: Not Found\n\n"
var profileSetNoArgErrorString string = "Error: drpcli profiles set [id] param [key] to [json blob] [flags] requires 5 arguments"
var profileSetMissingProfileErrorString string = "Error: GET: profiles/john2: Not Found\n\n"
var profileParamsNoArgErrorString string = "Error: drpcli profiles params [id] [json] [flags] requires 1 or 2 arguments\n"
var profileParamsMissingProfileErrorString string = "Error: GET: profiles/john2: Not Found\n\n"
var profilesParamsSetMissingProfileString string = "Error: POST: profiles/john2: Not Found\n\n"

var profileDefaultListString string = `[
  {
    "Available": true,
    "Description": "Global profile attached automatically to all machines.",
    "Errors": [],
    "Meta": {
      "color": "blue",
      "icon": "world",
      "title": "Digital Rebar Provision"
    },
    "Name": "global",
    "ReadOnly": false,
    "Validated": true
  }
]
`

var profileEmptyListString string = "[]\n"
var profileShowProfileString string = `{
  "Available": true,
  "Errors": [],
  "Name": "john",
  "Params": {
    "FRED": "GREG"
  },
  "ReadOnly": false,
  "Validated": true
}
`

var profileExistsProfileString string = ""

var profileCreateBadJSONString = "{asdgasdg"

var profileCreateBadJSON2String = "[asdgasdg]"
var profileCreateInputString string = `{
  "Name": "john",
  "Params": {
    "FRED": "GREG"
  }
}
`
var profileCreateJohnString string = `{
  "Available": true,
  "Errors": [],
  "Name": "john",
  "Params": {
    "FRED": "GREG"
  },
  "ReadOnly": false,
  "Validated": true
}
`

var profileListProfilesString = `[
  {
    "Available": true,
    "Description": "Global profile attached automatically to all machines.",
    "Errors": [],
    "Meta": {
      "color": "blue",
      "icon": "world",
      "title": "Digital Rebar Provision"
    },
    "Name": "global",
    "ReadOnly": false,
    "Validated": true
  },
  {
    "Available": true,
    "Errors": [],
    "Name": "john",
    "Params": {
      "FRED": "GREG"
    },
    "ReadOnly": false,
    "Validated": true
  }
]
`
var profileListJohnOnlyString = `[
  {
    "Available": true,
    "Errors": [],
    "Name": "john",
    "Params": {
      "FRED": "GREG"
    },
    "ReadOnly": false,
    "Validated": true
  }
]
`

var profileUpdateBadJSONString = "asdgasdg"

var profileUpdateInputString string = `{
  "Params": {
    "JESSIE": "JAMES"
  }
}
`
var profileUpdateJohnString string = `{
  "Available": true,
  "Errors": [],
  "Name": "john",
  "Params": {
    "FRED": "GREG",
    "JESSIE": "JAMES"
  },
  "ReadOnly": false,
  "Validated": true
}
`

var profilePatchBadPatchJSONString = "asdgasdg"

var profilePatchBadBaseJSONString = "asdgasdg"

var profilePatchBaseString string = `{
  "Name": "john",
  "Params": {
    "FRED": "GREG",
    "JESSIE": "JAMES"
  }
}
`
var profilePatchInputString string = `{
  "Params": {
    "JOHN": "StClaire",
    "JESSIE": "HAUG",
    "FRED": "LYNN"
  }
}
`
var profilePatchJohnString string = `{
  "Available": true,
  "Errors": [],
  "Name": "john",
  "Params": {
    "FRED": "LYNN",
    "JESSIE": "HAUG",
    "JOHN": "StClaire"
  },
  "ReadOnly": false,
  "Validated": true
}
`
var profilePatchMissingBaseString string = `{
  "Name": "john2",
  "Params": {
    "Name": ""
  }
}
`

var profileDestroyJohnString string = "Deleted profile john\n"

var profileParamsStartingString string = `{
  "FRED": "GREG",
  "JESSIE": "JAMES",
  "john3": 4
}
`
var profilesParamsNextString string = `{
  "jj": 3
}
`
var profileUpdateJohnWithParamsString string = `{
  "Available": true,
  "Errors": [],
  "Name": "john",
  "Params": {
    "jj": 3
  },
  "ReadOnly": false,
  "Validated": true
}
`

func TestProfileCli(t *testing.T) {
	cliTest(true, false, "profiles").run(t)
	cliTest(false, false, "profiles", "list").run(t)
	cliTest(true, true, "profiles", "create").run(t)
	cliTest(true, true, "profiles", "create", "john", "john2").run(t)
	cliTest(false, true, "profiles", "create", profileCreateBadJSONString).run(t)
	cliTest(false, true, "profiles", "create", profileCreateBadJSON2String).run(t)
	cliTest(false, false, "profiles", "create", profileCreateInputString).run(t)
	cliTest(false, true, "profiles", "create", profileCreateInputString).run(t)
	cliTest(false, false, "profiles", "list").run(t)
	cliTest(false, false, "profiles", "list", "Name=fred").run(t)
	cliTest(false, false, "profiles", "list", "Name=john").run(t)
	cliTest(true, true, "profiles", "show").run(t)
	cliTest(true, true, "profiles", "show", "john", "john2").run(t)
	cliTest(false, true, "profiles", "show", "john2").run(t)
	cliTest(false, false, "profiles", "show", "john").run(t)
	cliTest(true, true, "profiles", "exists").run(t)
	cliTest(true, true, "profiles", "exists", "john", "john2").run(t)
	cliTest(false, false, "profiles", "exists", "john").run(t)
	cliTest(false, true, "profiles", "exists", "john2").run(t)
	cliTest(true, true, "profiles", "exists", "john", "john2").run(t)
	cliTest(true, true, "profiles", "update").run(t)
	cliTest(true, true, "profiles", "update", "john", "john2", "john3").run(t)
	cliTest(false, true, "profiles", "update", "john", profileUpdateBadJSONString).run(t)
	cliTest(false, false, "profiles", "update", "john", profileUpdateInputString).run(t)
	cliTest(false, true, "profiles", "update", "john2", profileUpdateInputString).run(t)
	cliTest(false, false, "profiles", "show", "john").run(t)
	cliTest(false, false, "profiles", "show", "john").run(t)
	cliTest(true, true, "profiles", "destroy").run(t)
	cliTest(true, true, "profiles", "destroy", "john", "june").run(t)
	cliTest(false, false, "profiles", "destroy", "john").run(t)
	cliTest(false, true, "profiles", "destroy", "john").run(t)
	cliTest(false, false, "profiles", "list").run(t)
	cliTest(false, false, "profiles", "create", "-").Stdin(profileCreateInputString + "\n").run(t)
	cliTest(false, false, "profiles", "list").run(t)
	cliTest(false, false, "profiles", "update", "john", "-").Stdin(profileUpdateInputString + "\n").run(t)
	cliTest(false, false, "profiles", "show", "john").run(t)
	cliTest(true, true, "profiles", "get").run(t)
	cliTest(false, true, "profiles", "get", "john2", "param", "john2").run(t)
	cliTest(false, false, "profiles", "get", "john", "param", "john2").run(t)
	cliTest(true, true, "profiles", "set").run(t)
	cliTest(false, true, "profiles", "set", "john2", "param", "john2", "to", "cow").run(t)
	cliTest(false, false, "profiles", "set", "john", "param", "john2", "to", "cow").run(t)
	cliTest(false, false, "profiles", "get", "john", "param", "john2").run(t)
	cliTest(false, false, "profiles", "set", "john", "param", "john2", "to", "3").run(t)
	cliTest(false, false, "profiles", "set", "john", "param", "john3", "to", "4").run(t)
	cliTest(false, false, "profiles", "get", "john", "param", "john2").run(t)
	cliTest(false, false, "profiles", "get", "john", "param", "john3").run(t)
	cliTest(false, false, "profiles", "set", "john", "param", "john2", "to", "null").run(t)
	cliTest(false, false, "profiles", "get", "john", "param", "john2").run(t)
	cliTest(false, false, "profiles", "get", "john", "param", "john3").run(t)
	cliTest(true, true, "profiles", "params").run(t)
	cliTest(false, true, "profiles", "params", "john2").run(t)
	cliTest(false, false, "profiles", "params", "john").run(t)
	cliTest(false, true, "profiles", "params", "john2", profilesParamsNextString).run(t)
	cliTest(false, false, "profiles", "params", "john", profilesParamsNextString).run(t)
	cliTest(false, false, "profiles", "params", "john").run(t)
	cliTest(false, false, "profiles", "show", "john").run(t)
	cliTest(false, false, "profiles", "destroy", "john").run(t)
	cliTest(false, false, "profiles", "list").run(t)
}
