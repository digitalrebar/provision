package cli

import (
	"testing"
)

var userShowMissingArgErrorString string = "Error: GET: users/ignore: Not Found\n\n"
var userExistsMissingIgnoreString string = "Error: GET: users/ignore: Not Found\n\n"
var userCreateDuplicateErrorString = "Error: CREATE: users/john: already exists\n\n"
var userUpdateJohnMissingErrorString string = "Error: GET: users/john2: Not Found\n\n"
var userPatchJohnMissingErrorString string = "Error: PATCH: users/john2: Not Found\n\n"
var userDestroyMissingJohnString string = "Error: DELETE: users/john: Not Found\n\n"
var userTokenUserNotFoundErrorString string = "Error: GET: users/greg: Not Found\n\n"
var userPasswordNotFoundErrorString string = "Error: PUT: users/jill: Not Found\n\n"

var userEmptyListString string = "[]\n"
var userDefaultListString string = `RE:
\[
  {
    "Available": true,
    "Errors": \[\],
    "Name": "rocketskates",
    "PasswordHash": null,
    "ReadOnly": false,
    "Secret": "[\s\S]*",
    "Validated": true
  }
\]
`

var userShowNoArgErrorString string = "Error: drpcli users show [id] [flags] requires 1 argument\n"
var userShowTooManyArgErrorString string = "Error: drpcli users show [id] [flags] requires 1 argument\n"

var userShowJohnString string = `RE:
{
  "Available": true,
  "Errors": \[\],
  "Name": "john",
  "PasswordHash": null,
  "ReadOnly": false,
  "Secret": "[\s\S]*",
  "Validated": true
}
`

var userExistsNoArgErrorString string = "Error: drpcli users exists [id] [flags] requires 1 argument"
var userExistsTooManyArgErrorString string = "Error: drpcli users exists [id] [flags] requires 1 argument"
var userExistsIgnoreString string = ""

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
var userCreateJohnString string = `RE:
{
  "Available": true,
  "Errors": \[\],
  "Name": "john",
  "PasswordHash": null,
  "ReadOnly": false,
  "Secret": "[\s\S]*",
  "Validated": true
}
`
var userCreateFredInputString string = `fred`
var userCreateFredString string = `RE:
{
  "Available": true,
  "Errors": \[\],
  "Name": "fred",
  "PasswordHash": null,
  "ReadOnly": false,
  "Secret": "[\s\S]*",
  "Validated": true
}
`
var userDestroyFredString string = "Deleted user fred\n"

var userListJohnOnlyString = `RE:
\[
  {
    "Available": true,
    "Errors": \[\],
    "Name": "john",
    "PasswordHash": null,
    "ReadOnly": false,
    "Secret": "[\s\S]*",
    "Validated": true
  }
\]
`
var userListBothEnvsString = `RE:
\[
  {
    "Available": true,
    "Errors": \[\],
    "Name": "john",
    "PasswordHash": null,
    "ReadOnly": false,
    "Secret": "[\s\S]*",
    "Validated": true
  },
  {
    "Available": true,
    "Errors": [],
    "Name": "rocketskates",
    "PasswordHash": null,
    "ReadOnly": false,
    "Secret": "[\s\S]*",
    "Validated": true
  }
\]
`

var userUpdateNoArgErrorString string = "Error: drpcli users update [id] [json] [flags] requires 2 arguments"
var userUpdateTooManyArgErrorString string = "Error: drpcli users update [id] [json] [flags] requires 2 arguments"
var userUpdateBadJSONString = "asdgasdg"
var userUpdateBadJSONErrorString = "Error: Unable to merge objects: json: cannot unmarshal string into Go value of type map[string]interface {}\n\n\n"
var userUpdateInputString string = `{
  "PasswordHash": "NewStrat"
}
`
var userUpdateJohnString string = `RE:
{
  "Available": true,
  "Errors": \[\],
  "Name": "john",
  "PasswordHash": null,
  "ReadOnly": false,
  "Secret": "[\s\S]*",
  "Validated": true
}
`

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
var userPatchJohnString string = `RE:
{
  "Available": true,
  "Errors": \[\],
  "Name": "john",
  "PasswordHash": null,
  "ReadOnly": false,
  "Secret": "[\s\S]*",
  "Validated": true
}
`
var userPatchMissingBaseString string = `{
  "Name": "john2",
  "PasswordHash": null
}
`

var userDestroyNoArgErrorString string = "Error: drpcli users destroy [id] [flags] requires 1 argument"
var userDestroyTooManyArgErrorString string = "Error: drpcli users destroy [id] [flags] requires 1 argument"
var userDestroyJohnString string = "Deleted user john\n"

var userTokenNoArgErrorString string = "Error: drpcli users token [id] [ttl [ttl]] [scope [scope]] [action [action]] [specific [specific]] [flags] needs at least 1 arg\n"
var userTokenTooManyArgErrorString string = "Error: drpcli users token [id] [ttl [ttl]] [scope [scope]] [action [action]] [specific [specific]] [flags] needs at least 1 and pairs arg\n"
var userTokenUnknownPairErrorString string = "Error: drpcli users token [id] [ttl [ttl]] [scope [scope]] [action [action]] [specific [specific]] [flags] does not support greg2\n"

var userTokenTTLNotNumberErrorString string = "Error: ttl should be a number: strconv.ParseInt: parsing \"cow\": invalid syntax\n\n"
var userTokenSuccessString string = `RE:
{
  "Info": {
    "api_port": 10001,
    "arch": "[\s\S]*",
    "dhcp_enabled": false,
    "features": \[
      "api-v3",
      "sane-exit-codes",
      "common-blob-size",
      "change-stage-map"
    \],
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
		CliTest{false, false, []string{"users", "list", "Name=fred"}, noStdinString, userEmptyListString, noErrorString},
		CliTest{false, false, []string{"users", "list", "Name=john"}, noStdinString, userListJohnOnlyString, noErrorString},
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
