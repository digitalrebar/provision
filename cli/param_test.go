package cli

import (
	"testing"
)

var paramDefaultListString string = "[]\n"

var paramShowNoArgErrorString string = "Error: rscli params show [id] requires 1 argument\n"
var paramShowTooManyArgErrorString string = "Error: rscli params show [id] requires 1 argument\n"
var paramShowMissingArgErrorString string = "Error: parameters GET: ignore: Not Found\n\n"
var paramShowJohnString string = `{
  "Name": "john",
  "Value": "asdg"
}
`

var paramExistsNoArgErrorString string = "Error: rscli params exists [id] requires 1 argument"
var paramExistsTooManyArgErrorString string = "Error: rscli params exists [id] requires 1 argument"
var paramExistsIgnoreString string = ""
var paramExistsMissingIgnoreString string = "Error: parameters GET: ignore: Not Found\n\n"

var paramCreateNoArgErrorString string = "Error: rscli params create [json] requires 1 argument\n"
var paramCreateTooManyArgErrorString string = "Error: rscli params create [json] requires 1 argument\n"
var paramCreateBadJSONString = "asdgasdg"
var paramCreateBadJSONErrorString = "Error: Invalid param object: error unmarshaling JSON: json: cannot unmarshal string into Go value of type models.Param\n\n"
var paramCreateInputString string = `{
  "Name": "john",
  "Value": "asdg"
}
`
var paramCreateJohnString string = `{
  "Name": "john",
  "Value": "asdg"
}
`
var paramCreateDuplicateErrorString = "Error: dataTracker create parameters: john already exists\n\n"

var paramListBothEnvsString = `[
  {
    "Name": "john",
    "Value": "asdg"
  }
]
`

var paramUpdateNoArgErrorString string = "Error: rscli params update [id] [json] requires 2 arguments"
var paramUpdateTooManyArgErrorString string = "Error: rscli params update [id] [json] requires 2 arguments"
var paramUpdateBadJSONString = "asdgasdg"
var paramUpdateBadJSONErrorString = "Error: Unable to merge objects: json: cannot unmarshal string into Go value of type map[string]interface {}\n\n\n"
var paramUpdateInputString string = `{
  "Value": "NewStrat"
}
`
var paramUpdateJohnString string = `{
  "Name": "john",
  "Value": "NewStrat"
}
`
var paramUpdateJohnMissingErrorString string = "Error: parameters GET: john2: Not Found\n\n"

var paramPatchNoArgErrorString string = "Error: rscli params patch [objectJson] [changesJson] requires 2 arguments"
var paramPatchTooManyArgErrorString string = "Error: rscli params patch [objectJson] [changesJson] requires 2 arguments"
var paramPatchBadPatchJSONString = "asdgasdg"
var paramPatchBadPatchJSONErrorString = "Error: Unable to parse rscli params patch [objectJson] [changesJson] JSON asdgasdg\nError: error unmarshaling JSON: json: cannot unmarshal string into Go value of type models.Param\n\n"
var paramPatchBadBaseJSONString = "asdgasdg"
var paramPatchBadBaseJSONErrorString = "Error: Unable to parse rscli params patch [objectJson] [changesJson] JSON asdgasdg\nError: error unmarshaling JSON: json: cannot unmarshal string into Go value of type models.Param\n\n"
var paramPatchBaseString string = `{
  "Name": "john",
  "Value": "NewStrat"
}
`
var paramPatchInputString string = `{
  "Value": "Strat2n1"
}
`
var paramPatchJohnString string = `{
  "Name": "john",
  "Value": "Strat2n1"
}
`
var paramPatchMissingBaseString string = `{
  "Name": "john2",
  "Value": "Strat2n1"
}
`
var paramPatchJohnMissingErrorString string = "Error: parameters: PATCH john2: Not Found\n\n"

var paramDestroyNoArgErrorString string = "Error: rscli params destroy [id] requires 1 argument"
var paramDestroyTooManyArgErrorString string = "Error: rscli params destroy [id] requires 1 argument"
var paramDestroyJohnString string = "Deleted param john\n"
var paramDestroyMissingJohnString string = "Error: parameters: DELETE john: Not Found\n\n"

func TestParamCli(t *testing.T) {
	tests := []CliTest{
		CliTest{true, false, []string{"params"}, noStdinString, "Access CLI commands relating to params\n", ""},
		CliTest{false, false, []string{"params", "list"}, noStdinString, paramDefaultListString, noErrorString},

		CliTest{true, true, []string{"params", "create"}, noStdinString, noContentString, paramCreateNoArgErrorString},
		CliTest{true, true, []string{"params", "create", "john", "john2"}, noStdinString, noContentString, paramCreateTooManyArgErrorString},
		CliTest{false, true, []string{"params", "create", paramCreateBadJSONString}, noStdinString, noContentString, paramCreateBadJSONErrorString},
		CliTest{false, false, []string{"params", "create", paramCreateInputString}, noStdinString, paramCreateJohnString, noErrorString},
		CliTest{false, true, []string{"params", "create", paramCreateInputString}, noStdinString, noContentString, paramCreateDuplicateErrorString},
		CliTest{false, false, []string{"params", "list"}, noStdinString, paramListBothEnvsString, noErrorString},

		CliTest{true, true, []string{"params", "show"}, noStdinString, noContentString, paramShowNoArgErrorString},
		CliTest{true, true, []string{"params", "show", "john", "john2"}, noStdinString, noContentString, paramShowTooManyArgErrorString},
		CliTest{false, true, []string{"params", "show", "ignore"}, noStdinString, noContentString, paramShowMissingArgErrorString},
		CliTest{false, false, []string{"params", "show", "john"}, noStdinString, paramShowJohnString, noErrorString},

		CliTest{true, true, []string{"params", "exists"}, noStdinString, noContentString, paramExistsNoArgErrorString},
		CliTest{true, true, []string{"params", "exists", "john", "john2"}, noStdinString, noContentString, paramExistsTooManyArgErrorString},
		CliTest{false, false, []string{"params", "exists", "john"}, noStdinString, paramExistsIgnoreString, noErrorString},
		CliTest{false, true, []string{"params", "exists", "ignore"}, noStdinString, noContentString, paramExistsMissingIgnoreString},
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
		CliTest{false, false, []string{"params", "list"}, noStdinString, paramListBothEnvsString, noErrorString},
		CliTest{false, false, []string{"params", "update", "john", "-"}, paramUpdateInputString + "\n", paramUpdateJohnString, noErrorString},
		CliTest{false, false, []string{"params", "show", "john"}, noStdinString, paramUpdateJohnString, noErrorString},

		CliTest{false, false, []string{"params", "destroy", "john"}, noStdinString, paramDestroyJohnString, noErrorString},
		CliTest{false, false, []string{"params", "list"}, noStdinString, paramDefaultListString, noErrorString},
	}

	for _, test := range tests {
		testCli(t, test)
	}

}
