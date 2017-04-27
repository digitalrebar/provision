package cli

import (
	"testing"
)

var profileDefaultListString string = `[
  {
    "Name": "global"
  }
]
`

var profileShowNoArgErrorString string = "Error: drpcli profiles show [id] requires 1 argument\n"
var profileShowTooManyArgErrorString string = "Error: drpcli profiles show [id] requires 1 argument\n"
var profileShowMissingArgErrorString string = "Error: profiles GET: john2: Not Found\n\n"
var profileShowProfileString string = `{
  "Name": "john",
  "Params": {
    "FRED": "GREG"
  }
}
`

var profileExistsNoArgErrorString string = "Error: drpcli profiles exists [id] requires 1 argument"
var profileExistsTooManyArgErrorString string = "Error: drpcli profiles exists [id] requires 1 argument"
var profileExistsProfileString string = ""
var profileExistsMissingJohnString string = "Error: profiles GET: john2: Not Found\n\n"

var profileCreateNoArgErrorString string = "Error: drpcli profiles create [json] requires 1 argument\n"
var profileCreateTooManyArgErrorString string = "Error: drpcli profiles create [json] requires 1 argument\n"
var profileCreateBadJSONString = "asdgasdg"
var profileCreateBadJSONErrorString = "Error: Invalid profile object: error unmarshaling JSON: json: cannot unmarshal string into Go value of type models.Profile\n\n"
var profileCreateInputString string = `{
  "Name": "john",
  "Params": {
    "FRED": "GREG"
  }
}
`
var profileCreateJohnString string = `{
  "Name": "john",
  "Params": {
    "FRED": "GREG"
  }
}
`
var profileCreateDuplicateErrorString = "Error: dataTracker create profiles: john already exists\n\n"

var profileListProfilesString = `[
  {
    "Name": "global"
  },
  {
    "Name": "john",
    "Params": {
      "FRED": "GREG"
    }
  }
]
`

var profileUpdateNoArgErrorString string = "Error: drpcli profiles update [id] [json] requires 2 arguments"
var profileUpdateTooManyArgErrorString string = "Error: drpcli profiles update [id] [json] requires 2 arguments"
var profileUpdateBadJSONString = "asdgasdg"
var profileUpdateBadJSONErrorString = "Error: Unable to merge objects: json: cannot unmarshal string into Go value of type map[string]interface {}\n\n\n"
var profileUpdateInputString string = `{
  "Params": {
    "JESSIE": "JAMES"
  }
}
`
var profileUpdateJohnString string = `{
  "Name": "john",
  "Params": {
    "FRED": "GREG",
    "JESSIE": "JAMES"
  }
}
`
var profileUpdateJohnMissingErrorString string = "Error: profiles GET: john2: Not Found\n\n"

var profilePatchNoArgErrorString string = "Error: drpcli profiles patch [objectJson] [changesJson] requires 2 arguments"
var profilePatchTooManyArgErrorString string = "Error: drpcli profiles patch [objectJson] [changesJson] requires 2 arguments"
var profilePatchBadPatchJSONString = "asdgasdg"
var profilePatchBadPatchJSONErrorString = "Error: Unable to parse drpcli profiles patch [objectJson] [changesJson] JSON asdgasdg\nError: error unmarshaling JSON: json: cannot unmarshal string into Go value of type models.Profile\n\n"
var profilePatchBadBaseJSONString = "asdgasdg"
var profilePatchBadBaseJSONErrorString = "Error: Unable to parse drpcli profiles patch [objectJson] [changesJson] JSON asdgasdg\nError: error unmarshaling JSON: json: cannot unmarshal string into Go value of type models.Profile\n\n"
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
  "Name": "john",
  "Params": {
    "FRED": "LYNN",
    "JESSIE": "HAUG",
    "JOHN": "StClaire"
  }
}
`
var profilePatchMissingBaseString string = `{
  "Name": "john2",
  "Params": {
    "Name": ""
  }
}
`
var profilePatchJohnMissingErrorString string = "Error: profiles: PATCH john2: Not Found\n\n"

var profileDestroyNoArgErrorString string = "Error: drpcli profiles destroy [id] requires 1 argument"
var profileDestroyTooManyArgErrorString string = "Error: drpcli profiles destroy [id] requires 1 argument"
var profileDestroyJohnString string = "Deleted profile john\n"
var profileDestroyMissingJohnString string = "Error: profiles: DELETE john: Not Found\n\n"

var profileBootEnvNoArgErrorString string = "Error: drpcli profiles bootenv [id] [bootenv] requires 2 arguments"
var profileBootEnvMissingProfileErrorString string = "Error: profiles GET: john: Not Found\n\n"

var profileGetNoArgErrorString string = "Error: drpcli profiles get [id] param [key] requires 3 arguments"
var profileGetMissingProfileErrorString string = "Error: profiles GET Params: john2: Not Found\n\n"

var profileSetNoArgErrorString string = "Error: drpcli profiles set [id] param [key] to [json blob] requires 5 arguments"
var profileSetMissingProfileErrorString string = "Error: profiles GET Params: john2: Not Found\n\n"

var profileParamsNoArgErrorString string = "Error: drpcli profiles params [id] [json] requires 1 or 2 arguments\n"
var profileParamsMissingProfileErrorString string = "Error: profiles GET Params: john2: Not Found\n\n"
var profilesParamsSetMissingProfileString string = "Error: profiles SET Params: john2: Not Found\n\n"

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
  "Name": "john",
  "Params": {
    "jj": 3
  }
}
`

func TestProfileCli(t *testing.T) {

	tests := []CliTest{
		CliTest{true, false, []string{"profiles"}, noStdinString, "Access CLI commands relating to profiles\n", ""},
		CliTest{false, false, []string{"profiles", "list"}, noStdinString, profileDefaultListString, noErrorString},

		CliTest{true, true, []string{"profiles", "create"}, noStdinString, noContentString, profileCreateNoArgErrorString},
		CliTest{true, true, []string{"profiles", "create", "john", "john2"}, noStdinString, noContentString, profileCreateTooManyArgErrorString},
		CliTest{false, true, []string{"profiles", "create", profileCreateBadJSONString}, noStdinString, noContentString, profileCreateBadJSONErrorString},
		CliTest{false, false, []string{"profiles", "create", profileCreateInputString}, noStdinString, profileCreateJohnString, noErrorString},
		CliTest{false, true, []string{"profiles", "create", profileCreateInputString}, noStdinString, noContentString, profileCreateDuplicateErrorString},
		CliTest{false, false, []string{"profiles", "list"}, noStdinString, profileListProfilesString, noErrorString},

		CliTest{true, true, []string{"profiles", "show"}, noStdinString, noContentString, profileShowNoArgErrorString},
		CliTest{true, true, []string{"profiles", "show", "john", "john2"}, noStdinString, noContentString, profileShowTooManyArgErrorString},
		CliTest{false, true, []string{"profiles", "show", "john2"}, noStdinString, noContentString, profileShowMissingArgErrorString},
		CliTest{false, false, []string{"profiles", "show", "john"}, noStdinString, profileShowProfileString, noErrorString},

		CliTest{true, true, []string{"profiles", "exists"}, noStdinString, noContentString, profileExistsNoArgErrorString},
		CliTest{true, true, []string{"profiles", "exists", "john", "john2"}, noStdinString, noContentString, profileExistsTooManyArgErrorString},
		CliTest{false, false, []string{"profiles", "exists", "john"}, noStdinString, profileExistsProfileString, noErrorString},
		CliTest{false, true, []string{"profiles", "exists", "john2"}, noStdinString, noContentString, profileExistsMissingJohnString},
		CliTest{true, true, []string{"profiles", "exists", "john", "john2"}, noStdinString, noContentString, profileExistsTooManyArgErrorString},

		CliTest{true, true, []string{"profiles", "update"}, noStdinString, noContentString, profileUpdateNoArgErrorString},
		CliTest{true, true, []string{"profiles", "update", "john", "john2", "john3"}, noStdinString, noContentString, profileUpdateTooManyArgErrorString},
		CliTest{false, true, []string{"profiles", "update", "john", profileUpdateBadJSONString}, noStdinString, noContentString, profileUpdateBadJSONErrorString},
		CliTest{false, false, []string{"profiles", "update", "john", profileUpdateInputString}, noStdinString, profileUpdateJohnString, noErrorString},
		CliTest{false, true, []string{"profiles", "update", "john2", profileUpdateInputString}, noStdinString, noContentString, profileUpdateJohnMissingErrorString},
		CliTest{false, false, []string{"profiles", "show", "john"}, noStdinString, profileUpdateJohnString, noErrorString},

		CliTest{true, true, []string{"profiles", "patch"}, noStdinString, noContentString, profilePatchNoArgErrorString},
		CliTest{true, true, []string{"profiles", "patch", "john", "john2", "john3"}, noStdinString, noContentString, profilePatchTooManyArgErrorString},
		CliTest{false, true, []string{"profiles", "patch", profilePatchBaseString, profilePatchBadPatchJSONString}, noStdinString, noContentString, profilePatchBadPatchJSONErrorString},
		CliTest{false, true, []string{"profiles", "patch", profilePatchBadBaseJSONString, profilePatchInputString}, noStdinString, noContentString, profilePatchBadBaseJSONErrorString},
		CliTest{false, false, []string{"profiles", "patch", profilePatchBaseString, profilePatchInputString}, noStdinString, profilePatchJohnString, noErrorString},
		CliTest{false, true, []string{"profiles", "patch", profilePatchMissingBaseString, profilePatchInputString}, noStdinString, noContentString, profilePatchJohnMissingErrorString},
		CliTest{false, false, []string{"profiles", "show", "john"}, noStdinString, profilePatchJohnString, noErrorString},

		CliTest{true, true, []string{"profiles", "destroy"}, noStdinString, noContentString, profileDestroyNoArgErrorString},
		CliTest{true, true, []string{"profiles", "destroy", "john", "june"}, noStdinString, noContentString, profileDestroyTooManyArgErrorString},
		CliTest{false, false, []string{"profiles", "destroy", "john"}, noStdinString, profileDestroyJohnString, noErrorString},
		CliTest{false, true, []string{"profiles", "destroy", "john"}, noStdinString, noContentString, profileDestroyMissingJohnString},
		CliTest{false, false, []string{"profiles", "list"}, noStdinString, profileDefaultListString, noErrorString},

		CliTest{false, false, []string{"profiles", "create", "-"}, profileCreateInputString + "\n", profileCreateJohnString, noErrorString},
		CliTest{false, false, []string{"profiles", "list"}, noStdinString, profileListProfilesString, noErrorString},
		CliTest{false, false, []string{"profiles", "update", "john", "-"}, profileUpdateInputString + "\n", profileUpdateJohnString, noErrorString},
		CliTest{false, false, []string{"profiles", "show", "john"}, noStdinString, profileUpdateJohnString, noErrorString},

		CliTest{true, true, []string{"profiles", "get"}, noStdinString, noContentString, profileGetNoArgErrorString},
		CliTest{false, true, []string{"profiles", "get", "john2", "param", "john2"}, noStdinString, noContentString, profileGetMissingProfileErrorString},
		CliTest{false, false, []string{"profiles", "get", "john", "param", "john2"}, noStdinString, "null\n", noErrorString},

		CliTest{true, true, []string{"profiles", "set"}, noStdinString, noContentString, profileSetNoArgErrorString},
		CliTest{false, true, []string{"profiles", "set", "john2", "param", "john2", "to", "cow"}, noStdinString, noContentString, profileSetMissingProfileErrorString},
		CliTest{false, false, []string{"profiles", "set", "john", "param", "john2", "to", "cow"}, noStdinString, "\"cow\"\n", noErrorString},
		CliTest{false, false, []string{"profiles", "get", "john", "param", "john2"}, noStdinString, "\"cow\"\n", noErrorString},
		CliTest{false, false, []string{"profiles", "set", "john", "param", "john2", "to", "3"}, noStdinString, "3\n", noErrorString},
		CliTest{false, false, []string{"profiles", "set", "john", "param", "john3", "to", "4"}, noStdinString, "4\n", noErrorString},
		CliTest{false, false, []string{"profiles", "get", "john", "param", "john2"}, noStdinString, "3\n", noErrorString},
		CliTest{false, false, []string{"profiles", "get", "john", "param", "john3"}, noStdinString, "4\n", noErrorString},
		CliTest{false, false, []string{"profiles", "set", "john", "param", "john2", "to", "null"}, noStdinString, "null\n", noErrorString},
		CliTest{false, false, []string{"profiles", "get", "john", "param", "john2"}, noStdinString, "null\n", noErrorString},
		CliTest{false, false, []string{"profiles", "get", "john", "param", "john3"}, noStdinString, "4\n", noErrorString},

		CliTest{true, true, []string{"profiles", "params"}, noStdinString, noContentString, profileParamsNoArgErrorString},
		CliTest{false, true, []string{"profiles", "params", "john2"}, noStdinString, noContentString, profileParamsMissingProfileErrorString},
		CliTest{false, false, []string{"profiles", "params", "john"}, noStdinString, profileParamsStartingString, noErrorString},
		CliTest{false, true, []string{"profiles", "params", "john2", profilesParamsNextString}, noStdinString, noContentString, profilesParamsSetMissingProfileString},
		CliTest{false, false, []string{"profiles", "params", "john", profilesParamsNextString}, noStdinString, profilesParamsNextString, noErrorString},
		CliTest{false, false, []string{"profiles", "params", "john"}, noStdinString, profilesParamsNextString, noErrorString},

		CliTest{false, false, []string{"profiles", "show", "john"}, noStdinString, profileUpdateJohnWithParamsString, noErrorString},

		CliTest{false, false, []string{"profiles", "destroy", "john"}, noStdinString, profileDestroyJohnString, noErrorString},
		CliTest{false, false, []string{"profiles", "list"}, noStdinString, profileDefaultListString, noErrorString},
	}

	for _, test := range tests {
		testCli(t, test)
	}

}
