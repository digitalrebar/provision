package cli

import (
	"testing"
)

var userDefaultListString string = `[
  {
    "Name": "rocketskates"
  }
]
`

var userShowNoArgErrorString string = "Error: drpcli users show [id] requires 1 argument\n"
var userShowTooManyArgErrorString string = "Error: drpcli users show [id] requires 1 argument\n"
var userShowMissingArgErrorString string = "Error: users GET: ignore: Not Found\n\n"
var userShowJohnString string = `{
  "Name": "john"
}
`

var userExistsNoArgErrorString string = "Error: drpcli users exists [id] requires 1 argument"
var userExistsTooManyArgErrorString string = "Error: drpcli users exists [id] requires 1 argument"
var userExistsIgnoreString string = ""
var userExistsMissingIgnoreString string = "Error: users GET: ignore: Not Found\n\n"

var userCreateNoArgErrorString string = "Error: drpcli users create [json] requires 1 argument\n"
var userCreateTooManyArgErrorString string = "Error: drpcli users create [json] requires 1 argument\n"
var userCreateBadJSONString = "asdgasdg"
var userCreateBadJSONErrorString = "Error: Invalid user object: error unmarshaling JSON: json: cannot unmarshal string into Go value of type models.User\n\n"
var userCreateInputString string = `{
  "Name": "john"
}
`
var userCreateJohnString string = `{
  "Name": "john"
}
`
var userCreateDuplicateErrorString = "Error: dataTracker create users: john already exists\n\n"

var userListBothEnvsString = `[
  {
    "Name": "john"
  },
  {
    "Name": "rocketskates"
  }
]
`

var userUpdateNoArgErrorString string = "Error: drpcli users update [id] [json] requires 2 arguments"
var userUpdateTooManyArgErrorString string = "Error: drpcli users update [id] [json] requires 2 arguments"
var userUpdateBadJSONString = "asdgasdg"
var userUpdateBadJSONErrorString = "Error: Unable to merge objects: json: cannot unmarshal string into Go value of type map[string]interface {}\n\n\n"
var userUpdateInputString string = `{
  "PasswordHash": "NewStrat"
}
`
var userUpdateJohnString string = `{
  "Name": "john"
}
`
var userUpdateJohnMissingErrorString string = "Error: users GET: john2: Not Found\n\n"

var userPatchNoArgErrorString string = "Error: drpcli users patch [objectJson] [changesJson] requires 2 arguments"
var userPatchTooManyArgErrorString string = "Error: drpcli users patch [objectJson] [changesJson] requires 2 arguments"
var userPatchBadPatchJSONString = "asdgasdg"
var userPatchBadPatchJSONErrorString = "Error: Unable to parse drpcli users patch [objectJson] [changesJson] JSON asdgasdg\nError: error unmarshaling JSON: json: cannot unmarshal string into Go value of type models.User\n\n"
var userPatchBadBaseJSONString = "asdgasdg"
var userPatchBadBaseJSONErrorString = "Error: Unable to parse drpcli users patch [objectJson] [changesJson] JSON asdgasdg\nError: error unmarshaling JSON: json: cannot unmarshal string into Go value of type models.User\n\n"
var userPatchBaseString string = `{
  "Name": "john"
}
`
var userPatchInputString string = `{
  "PasswordHash": "Strat2n1"
}
`
var userPatchJohnString string = `{
  "Name": "john"
}
`
var userPatchMissingBaseString string = `{
  "Name": "john2"
}
`
var userPatchJohnMissingErrorString string = "Error: users: PATCH john2: Not Found\n\n"

var userDestroyNoArgErrorString string = "Error: drpcli users destroy [id] requires 1 argument"
var userDestroyTooManyArgErrorString string = "Error: drpcli users destroy [id] requires 1 argument"
var userDestroyJohnString string = "Deleted user john\n"
var userDestroyMissingJohnString string = "Error: users: DELETE john: Not Found\n\n"

var userTokenNoArgErrorString string = "Error: drpcli users token [id] [ttl [ttl]] [scope [scope]] [action [action]] [specific [specific]] needs at least 1 arg\n"
var userTokenTooManyArgErrorString string = "Error: drpcli users token [id] [ttl [ttl]] [scope [scope]] [action [action]] [specific [specific]] needs at least 1 and pairs arg\n"
var userTokenUnknownPairErrorString string = "Error: drpcli users token [id] [ttl [ttl]] [scope [scope]] [action [action]] [specific [specific]] does not support greg2\n"
var userTokenUserNotFoundErrorString string = "Error: User GET: greg: Not Found\n\n"
var userTokenTTLNotNumberErrorString string = "Error: ttl should be a number: strconv.ParseInt: parsing \"cow\": invalid syntax\n\n"
var userTokenSuccessString string = `RE:
{
  "Token": "[\s\S]*"
}
`

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

		CliTest{true, true, []string{"users", "token"}, noStdinString, noContentString, userTokenNoArgErrorString},
		CliTest{true, true, []string{"users", "token", "greg", "greg2"}, noStdinString, noContentString, userTokenTooManyArgErrorString},
		CliTest{true, true, []string{"users", "token", "greg", "greg2", "greg3"}, noStdinString, noContentString, userTokenUnknownPairErrorString},
		CliTest{false, true, []string{"users", "token", "greg"}, noStdinString, noContentString, userTokenUserNotFoundErrorString},
		CliTest{false, false, []string{"users", "token", "rocketskates"}, noStdinString, userTokenSuccessString, noErrorString},
		CliTest{false, false, []string{"users", "token", "rocketskates", "scope", "all", "ttl", "330", "action", "list", "specific", "asdgag"}, noStdinString, userTokenSuccessString, noErrorString},
		CliTest{false, true, []string{"users", "token", "rocketskates", "scope", "all", "ttl", "cow", "action", "list", "specific", "asdgag"}, noStdinString, noContentString, userTokenTTLNotNumberErrorString},
	}

	for _, test := range tests {
		testCli(t, test)
	}

}
