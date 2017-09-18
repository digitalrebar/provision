package cli

import (
	"testing"
)

var userEmptyListString string = "[]\n"
var userDefaultListString string = `[
  {
    "Available": true,
    "Errors": null,
    "Name": "rocketskates",
    "PasswordHash": null,
    "ReadOnly": false,
    "Validated": true
  }
]
`

var userShowNoArgErrorString string = "Error: drpcli users show [id] [flags] requires 1 argument\n"
var userShowTooManyArgErrorString string = "Error: drpcli users show [id] [flags] requires 1 argument\n"
var userShowMissingArgErrorString string = "Error: users GET: ignore: Not Found\n\n"
var userShowJohnString string = `{
  "Available": true,
  "Errors": null,
  "Name": "john",
  "PasswordHash": null,
  "ReadOnly": false,
  "Validated": true
}
`

var userExistsNoArgErrorString string = "Error: drpcli users exists [id] [flags] requires 1 argument"
var userExistsTooManyArgErrorString string = "Error: drpcli users exists [id] [flags] requires 1 argument"
var userExistsIgnoreString string = ""
var userExistsMissingIgnoreString string = "Error: users GET: ignore: Not Found\n\n"

var userCreateNoArgErrorString string = "Error: drpcli users create [json] [flags] requires 1 argument\n"
var userCreateTooManyArgErrorString string = "Error: drpcli users create [json] [flags] requires 1 argument\n"
var userCreateBadJSONString = "{asdgasdg"
var userCreateBadJSONErrorString = "Error: Invalid user object: error converting YAML to JSON: yaml: line 1: did not find expected ',' or '}' and error converting YAML to JSON: yaml: line 1: did not find expected ',' or '}'\n\n"
var userCreateBadJSON2String = "[asdgasdg]"
var userCreateBadJSON2ErrorString = "Error: Unable to create new user: Invalid type passed to user create\n\n"
var userCreateInputString string = `{
  "Name": "john",
  "PasswordHash": null
}
`
var userCreateJohnString string = `{
  "Available": true,
  "Errors": null,
  "Name": "john",
  "PasswordHash": null,
  "ReadOnly": false,
  "Validated": true
}
`
var userCreateFredInputString string = `fred`
var userCreateFredString string = `{
  "Available": true,
  "Errors": null,
  "Name": "fred",
  "PasswordHash": null,
  "ReadOnly": false,
  "Validated": true
}
`
var userDestroyFredString string = "Deleted user fred\n"
var userCreateDuplicateErrorString = "Error: dataTracker create users: john already exists\n\n"

var userListJohnOnlyString = `[
  {
    "Available": true,
    "Errors": null,
    "Name": "john",
    "PasswordHash": null,
    "ReadOnly": false,
    "Validated": true
  }
]
`
var userListBothEnvsString = `[
  {
    "Available": true,
    "Errors": null,
    "Name": "john",
    "PasswordHash": null,
    "ReadOnly": false,
    "Validated": true
  },
  {
    "Available": true,
    "Errors": null,
    "Name": "rocketskates",
    "PasswordHash": null,
    "ReadOnly": false,
    "Validated": true
  }
]
`

var userUpdateNoArgErrorString string = "Error: drpcli users update [id] [json] [flags] requires 2 arguments"
var userUpdateTooManyArgErrorString string = "Error: drpcli users update [id] [json] [flags] requires 2 arguments"
var userUpdateBadJSONString = "asdgasdg"
var userUpdateBadJSONErrorString = "Error: Unable to merge objects: json: cannot unmarshal string into Go value of type map[string]interface {}\n\n\n"
var userUpdateInputString string = `{
  "PasswordHash": "NewStrat"
}
`
var userUpdateJohnString string = `{
  "Available": true,
  "Errors": null,
  "Name": "john",
  "PasswordHash": null,
  "ReadOnly": false,
  "Validated": true
}
`
var userUpdateJohnMissingErrorString string = "Error: users GET: john2: Not Found\n\n"

var userPatchNoArgErrorString string = "Error: drpcli users patch [objectJson] [changesJson] [flags] requires 2 arguments"
var userPatchTooManyArgErrorString string = "Error: drpcli users patch [objectJson] [changesJson] [flags] requires 2 arguments"
var userPatchBadPatchJSONString = "asdgasdg"
var userPatchBadPatchJSONErrorString = "Error: Unable to parse drpcli users patch [objectJson] [changesJson] [flags] JSON asdgasdg\nError: error unmarshaling JSON: json: cannot unmarshal string into Go value of type genmodels.User\n\n"
var userPatchBadBaseJSONString = "asdgasdg"
var userPatchBadBaseJSONErrorString = "Error: Unable to parse drpcli users patch [objectJson] [changesJson] [flags] JSON asdgasdg\nError: error unmarshaling JSON: json: cannot unmarshal string into Go value of type genmodels.User\n\n"
var userPatchBaseString string = `{
  "Name": "john"
}
`
var userPatchInputString string = `{
  "PasswordHash": "Strat2n1"
}
`
var userPatchJohnString string = `{
  "Available": true,
  "Errors": null,
  "Name": "john",
  "PasswordHash": null,
  "ReadOnly": false,
  "Validated": true
}
`
var userPatchMissingBaseString string = `{
  "Name": "john2",
  "PasswordHash": null
}
`
var userPatchJohnMissingErrorString string = "Error: users: PATCH john2: Not Found\n\n"

var userDestroyNoArgErrorString string = "Error: drpcli users destroy [id] [flags] requires 1 argument"
var userDestroyTooManyArgErrorString string = "Error: drpcli users destroy [id] [flags] requires 1 argument"
var userDestroyJohnString string = "Deleted user john\n"
var userDestroyMissingJohnString string = "Error: users: DELETE john: Not Found\n\n"

var userTokenNoArgErrorString string = "Error: drpcli users token [id] [ttl [ttl]] [scope [scope]] [action [action]] [specific [specific]] [flags] needs at least 1 arg\n"
var userTokenTooManyArgErrorString string = "Error: drpcli users token [id] [ttl [ttl]] [scope [scope]] [action [action]] [specific [specific]] [flags] needs at least 1 and pairs arg\n"
var userTokenUnknownPairErrorString string = "Error: drpcli users token [id] [ttl [ttl]] [scope [scope]] [action [action]] [specific [specific]] [flags] does not support greg2\n"
var userTokenUserNotFoundErrorString string = "Error: User GET: greg: Not Found\n\n"
var userTokenTTLNotNumberErrorString string = "Error: ttl should be a number: strconv.ParseInt: parsing \"cow\": invalid syntax\n\n"
var userTokenSuccessString string = `RE:
{
  "Info": {
    "api_port": 10001,
    "arch": "[\s\S]*",
    "dhcp_enabled": false,
    "file_port": 10002,
    "id": "Fred",
    "os": "[\s\S]*",
    "prov_enabled": true,
    "stats": \[
      {
        "count": 0,
        "name": "machines.count"
      },
      {
        "count": 0,
        "name": "subnets.count"
      }
    \],
    "tftp_enabled": true,
    "version": "[\s\S]*"
  },
  "Token": "[\s\S]*"
}
`

var userPasswordNoArgsErrorString string = "Error: drpcli users password [id] [password] [flags] needs 2 args\n"
var userPasswordNotFoundErrorString string = "Error: User GET: jill: Not Found\n\n"

func TestUserCli(t *testing.T) {

	tests := []CliTest{
		CliTest{true, false, []string{"users"}, noStdinString, "Access CLI commands relating to users\n", ""},
		CliTest{false, false, []string{"users", "list"}, noStdinString, userDefaultListString, noErrorString},

		CliTest{true, true, []string{"users", "create"}, noStdinString, noContentString, userCreateNoArgErrorString},
		CliTest{true, true, []string{"users", "create", "john", "john2"}, noStdinString, noContentString, userCreateTooManyArgErrorString},
		CliTest{false, true, []string{"users", "create", userCreateBadJSONString}, noStdinString, noContentString, userCreateBadJSONErrorString},
		CliTest{false, true, []string{"users", "create", userCreateBadJSON2String}, noStdinString, noContentString, userCreateBadJSON2ErrorString},
		CliTest{false, false, []string{"users", "create", userCreateInputString}, noStdinString, userCreateJohnString, noErrorString},
		CliTest{false, false, []string{"users", "create", userCreateFredInputString}, noStdinString, userCreateFredString, noErrorString},
		CliTest{false, false, []string{"users", "destroy", userCreateFredInputString}, noStdinString, userDestroyFredString, noErrorString},
		CliTest{false, true, []string{"users", "create", userCreateInputString}, noStdinString, noContentString, userCreateDuplicateErrorString},
		CliTest{false, false, []string{"users", "list"}, noStdinString, userListBothEnvsString, noErrorString},

		CliTest{false, false, []string{"users", "list", "--limit=0"}, noStdinString, userEmptyListString, noErrorString},
		CliTest{false, false, []string{"users", "list", "--limit=10", "--offset=0"}, noStdinString, userListBothEnvsString, noErrorString},
		CliTest{false, false, []string{"users", "list", "--limit=10", "--offset=10"}, noStdinString, userEmptyListString, noErrorString},
		CliTest{false, true, []string{"users", "list", "--limit=-10", "--offset=0"}, noStdinString, noContentString, limitNegativeError},
		CliTest{false, true, []string{"users", "list", "--limit=10", "--offset=-10"}, noStdinString, noContentString, offsetNegativeError},
		CliTest{false, false, []string{"users", "list", "--limit=-1", "--offset=-1"}, noStdinString, userListBothEnvsString, noErrorString},
		CliTest{false, false, []string{"users", "list", "Name=fred"}, noStdinString, userEmptyListString, noErrorString},
		CliTest{false, false, []string{"users", "list", "Name=john"}, noStdinString, userListJohnOnlyString, noErrorString},
		CliTest{false, false, []string{"users", "list", "Available=true"}, noStdinString, userListBothEnvsString, noErrorString},
		CliTest{false, false, []string{"users", "list", "Available=false"}, noStdinString, userEmptyListString, noErrorString},
		CliTest{false, true, []string{"users", "list", "Available=fred"}, noStdinString, noContentString, bootEnvBadAvailableString},
		CliTest{false, false, []string{"users", "list", "Valid=true"}, noStdinString, userListBothEnvsString, noErrorString},
		CliTest{false, false, []string{"users", "list", "Valid=false"}, noStdinString, userEmptyListString, noErrorString},
		CliTest{false, true, []string{"users", "list", "Valid=fred"}, noStdinString, noContentString, bootEnvBadValidString},
		CliTest{false, false, []string{"users", "list", "ReadOnly=true"}, noStdinString, userEmptyListString, noErrorString},
		CliTest{false, false, []string{"users", "list", "ReadOnly=false"}, noStdinString, userListBothEnvsString, noErrorString},
		CliTest{false, true, []string{"users", "list", "ReadOnly=fred"}, noStdinString, noContentString, bootEnvBadReadOnlyString},

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
		CliTest{false, true, []string{"users", "update", "john2", userUpdateInputString}, noStdinString, noContentString, userUpdateJohnMissingErrorString},
		CliTest{false, false, []string{"users", "show", "john"}, noStdinString, userUpdateJohnString, noErrorString},

		CliTest{true, true, []string{"users", "patch"}, noStdinString, noContentString, userPatchNoArgErrorString},
		CliTest{true, true, []string{"users", "patch", "john", "john2", "john3"}, noStdinString, noContentString, userPatchTooManyArgErrorString},
		CliTest{false, true, []string{"users", "patch", userPatchBaseString, userPatchBadPatchJSONString}, noStdinString, noContentString, userPatchBadPatchJSONErrorString},
		CliTest{false, true, []string{"users", "patch", userPatchBadBaseJSONString, userPatchInputString}, noStdinString, noContentString, userPatchBadBaseJSONErrorString},
		CliTest{false, false, []string{"users", "patch", userPatchBaseString, userPatchInputString}, noStdinString, userPatchJohnString, noErrorString},
		CliTest{false, true, []string{"users", "patch", userPatchMissingBaseString, userPatchInputString}, noStdinString, noContentString, userPatchJohnMissingErrorString},
		CliTest{false, false, []string{"users", "show", "john"}, noStdinString, userPatchJohnString, noErrorString},

		CliTest{true, true, []string{"users", "password"}, noStdinString, noContentString, userPasswordNoArgsErrorString},
		CliTest{true, true, []string{"users", "password", "one"}, noStdinString, noContentString, userPasswordNoArgsErrorString},
		CliTest{true, true, []string{"users", "password", "one", "two", "three"}, noStdinString, noContentString, userPasswordNoArgsErrorString},
		CliTest{false, true, []string{"users", "password", "jill", "june"}, noStdinString, noContentString, userPasswordNotFoundErrorString},
		CliTest{false, false, []string{"users", "password", "john", "june"}, noStdinString, userPatchJohnString, noErrorString},

		CliTest{true, true, []string{"users", "destroy"}, noStdinString, noContentString, userDestroyNoArgErrorString},
		CliTest{true, true, []string{"users", "destroy", "john", "june"}, noStdinString, noContentString, userDestroyTooManyArgErrorString},
		CliTest{false, false, []string{"users", "destroy", "john"}, noStdinString, userDestroyJohnString, noErrorString},
		CliTest{false, true, []string{"users", "destroy", "john"}, noStdinString, noContentString, userDestroyMissingJohnString},
		CliTest{false, false, []string{"users", "list"}, noStdinString, userDefaultListString, noErrorString},

		CliTest{false, false, []string{"users", "create", "-"}, userCreateInputString + "\n", userCreateJohnString, noErrorString},
		CliTest{false, false, []string{"users", "list"}, noStdinString, userListBothEnvsString, noErrorString},
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
