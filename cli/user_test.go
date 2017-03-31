package cli

import (
	"testing"
)

var userDefaultListString string = "[]\n"

var userShowNoArgErrorString string = "Error: rscli users show [id] requires 1 argument\n"
var userShowTooManyArgErrorString string = "Error: rscli users show [id] requires 1 argument\n"
var userShowMissingArgErrorString string = "Error: users GET: ignore: Not Found\n\n"
var userShowJohnString string = `{
  "Name": "john",
  "PasswordHash": "asdg"
}
`

var userExistsNoArgErrorString string = "Error: rscli users exists [id] requires 1 argument"
var userExistsTooManyArgErrorString string = "Error: rscli users exists [id] requires 1 argument"
var userExistsIgnoreString string = ""
var userExistsMissingIgnoreString string = "Error: users GET: ignore: Not Found\n\n"

var userCreateNoArgErrorString string = "Error: rscli users create [json] requires 1 argument\n"
var userCreateTooManyArgErrorString string = "Error: rscli users create [json] requires 1 argument\n"
var userCreateBadJSONString = "asdgasdg"
var userCreateBadJSONErrorString = "Error: Invalid user object: error unmarshaling JSON: json: cannot unmarshal string into Go value of type models.User\n\n"
var userCreateInputString string = `{
  "Name": "john",
  "PasswordHash": "asdg"
}
`
var userCreateJohnString string = `{
  "Name": "john",
  "PasswordHash": "asdg"
}
`
var userCreateDuplicateErrorString = "Error: dataTracker create users: john already exists\n\n"

var userListBothEnvsString = `[
  {
    "Name": "john",
    "PasswordHash": "asdg"
  }
]
`

var userUpdateNoArgErrorString string = "Error: rscli users update [id] [json] requires 2 arguments"
var userUpdateTooManyArgErrorString string = "Error: rscli users update [id] [json] requires 2 arguments"
var userUpdateBadJSONString = "asdgasdg"
var userUpdateBadJSONErrorString = "Error: Unable to merge objects: json: cannot unmarshal string into Go value of type map[string]interface {}\n\n\n"
var userUpdateInputString string = `{
  "PasswordHash": "NewStrat"
}
`
var userUpdateJohnString string = `{
  "Name": "john",
  "PasswordHash": "NewStrat"
}
`
var userUpdateJohnMissingErrorString string = "Error: users GET: john2: Not Found\n\n"

var userPatchNoArgErrorString string = "Error: rscli users patch [objectJson] [changesJson] requires 2 arguments"
var userPatchTooManyArgErrorString string = "Error: rscli users patch [objectJson] [changesJson] requires 2 arguments"
var userPatchBadPatchJSONString = "asdgasdg"
var userPatchBadPatchJSONErrorString = "Error: Unable to parse rscli users patch [objectJson] [changesJson] JSON asdgasdg\nError: error unmarshaling JSON: json: cannot unmarshal string into Go value of type models.User\n\n"
var userPatchBadBaseJSONString = "asdgasdg"
var userPatchBadBaseJSONErrorString = "Error: Unable to parse rscli users patch [objectJson] [changesJson] JSON asdgasdg\nError: error unmarshaling JSON: json: cannot unmarshal string into Go value of type models.User\n\n"
var userPatchBaseString string = `{
  "Name": "john",
  "PasswordHash": "NewStrat"
}
`
var userPatchInputString string = `{
  "PasswordHash": "Strat2n1"
}
`
var userPatchJohnString string = `{
  "Name": "john",
  "PasswordHash": "Strat2n1"
}
`
var userPatchMissingBaseString string = `{
  "Name": "john2",
  "PasswordHash": "Strat2n1"
}
`
var userPatchJohnMissingErrorString string = "Error: users: PATCH john2: Not Found\n\n"

var userDestroyNoArgErrorString string = "Error: rscli users destroy [id] requires 1 argument"
var userDestroyTooManyArgErrorString string = "Error: rscli users destroy [id] requires 1 argument"
var userDestroyJohnString string = "Deleted user john\n"
var userDestroyMissingJohnString string = "Error: users: DELETE john: Not Found\n\n"

func TestUserCli(t *testing.T) {
	tests := []CliTest{
		CliTest{true, false, []string{"users"}, noStdinString, "Access CLI commands relating to users\n", ""},
		CliTest{false, false, []string{"users", "list"}, noStdinString, userDefaultListString, noErrorString},

		CliTest{true, true, []string{"users", "create"}, noStdinString, noContentString, userCreateNoArgErrorString},
		CliTest{true, true, []string{"users", "create", "john", "john2"}, noStdinString, noContentString, userCreateTooManyArgErrorString},
		CliTest{false, true, []string{"users", "create", userCreateBadJSONString}, noStdinString, noContentString, userCreateBadJSONErrorString},
		CliTest{false, false, []string{"users", "create", userCreateInputString}, noStdinString, userCreateJohnString, noErrorString},
		CliTest{false, true, []string{"users", "create", userCreateInputString}, noStdinString, noContentString, userCreateDuplicateErrorString},
		CliTest{false, false, []string{"users", "list"}, noStdinString, userListBothEnvsString, noErrorString},

		CliTest{true, true, []string{"users", "show"}, noStdinString, noContentString, userShowNoArgErrorString},
		CliTest{true, true, []string{"users", "show", "john", "john2"}, noStdinString, noContentString, userShowTooManyArgErrorString},
		CliTest{false, true, []string{"users", "show", "ignore"}, noStdinString, noContentString, userShowMissingArgErrorString},
		CliTest{false, false, []string{"users", "show", "john"}, noStdinString, userShowJohnString, noErrorString},

		CliTest{true, true, []string{"users", "exists"}, noStdinString, noContentString, userExistsNoArgErrorString},
		CliTest{true, true, []string{"users", "exists", "john", "john2"}, noStdinString, noContentString, userExistsTooManyArgErrorString},
		CliTest{false, false, []string{"users", "exists", "john"}, noStdinString, userExistsIgnoreString, noErrorString},
		CliTest{false, true, []string{"users", "exists", "ignore"}, noStdinString, noContentString, userExistsMissingIgnoreString},
		CliTest{true, true, []string{"users", "exists", "john", "john2"}, noStdinString, noContentString, userExistsTooManyArgErrorString},

		CliTest{true, true, []string{"users", "update"}, noStdinString, noContentString, userUpdateNoArgErrorString},
		CliTest{true, true, []string{"users", "update", "john", "john2", "john3"}, noStdinString, noContentString, userUpdateTooManyArgErrorString},
		CliTest{false, true, []string{"users", "update", "john", userUpdateBadJSONString}, noStdinString, noContentString, userUpdateBadJSONErrorString},
		CliTest{false, false, []string{"users", "update", "john", userUpdateInputString}, noStdinString, userUpdateJohnString, noErrorString},
		CliTest{false, true, []string{"users", "update", "john2", userUpdateInputString}, noStdinString, noContentString, userUpdateJohnMissingErrorString},
		CliTest{false, false, []string{"users", "show", "john"}, noStdinString, userUpdateJohnString, noErrorString},

		CliTest{true, true, []string{"users", "patch"}, noStdinString, noContentString, userPatchNoArgErrorString},
		CliTest{true, true, []string{"users", "patch", "john", "john2", "john3"}, noStdinString, noContentString, userPatchTooManyArgErrorString},
		CliTest{false, true, []string{"users", "patch", userPatchBaseString, userPatchBadPatchJSONString}, noStdinString, noContentString, userPatchBadPatchJSONErrorString},
		CliTest{false, true, []string{"users", "patch", userPatchBadBaseJSONString, userPatchInputString}, noStdinString, noContentString, userPatchBadBaseJSONErrorString},
		CliTest{false, false, []string{"users", "patch", userPatchBaseString, userPatchInputString}, noStdinString, userPatchJohnString, noErrorString},
		CliTest{false, true, []string{"users", "patch", userPatchMissingBaseString, userPatchInputString}, noStdinString, noContentString, userPatchJohnMissingErrorString},
		CliTest{false, false, []string{"users", "show", "john"}, noStdinString, userPatchJohnString, noErrorString},

		CliTest{true, true, []string{"users", "destroy"}, noStdinString, noContentString, userDestroyNoArgErrorString},
		CliTest{true, true, []string{"users", "destroy", "john", "june"}, noStdinString, noContentString, userDestroyTooManyArgErrorString},
		CliTest{false, false, []string{"users", "destroy", "john"}, noStdinString, userDestroyJohnString, noErrorString},
		CliTest{false, true, []string{"users", "destroy", "john"}, noStdinString, noContentString, userDestroyMissingJohnString},
		CliTest{false, false, []string{"users", "list"}, noStdinString, userDefaultListString, noErrorString},

		CliTest{false, false, []string{"users", "create", "-"}, userCreateInputString + "\n", userCreateJohnString, noErrorString},
		CliTest{false, false, []string{"users", "list"}, noStdinString, userListBothEnvsString, noErrorString},
		CliTest{false, false, []string{"users", "update", "john", "-"}, userUpdateInputString + "\n", userUpdateJohnString, noErrorString},
		CliTest{false, false, []string{"users", "show", "john"}, noStdinString, userUpdateJohnString, noErrorString},

		CliTest{false, false, []string{"users", "destroy", "john"}, noStdinString, userDestroyJohnString, noErrorString},
		CliTest{false, false, []string{"users", "list"}, noStdinString, userDefaultListString, noErrorString},
	}

	for _, test := range tests {
		testCli(t, test)
	}

}
