package cli

import "testing"

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
      "change-stage-map",
      "job-exit-states",
      "package-repository-handling"
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
	cliTest(true, false, "users").run(t)
	cliTest(false, false, "users", "list").run(t)
	cliTest(true, true, "users", "create").run(t)
	cliTest(true, true, "users", "create", "john", "john2").run(t)
	cliTest(false, true, "users", "create", userCreateBadJSONString).run(t)
	cliTest(false, true, "users", "create", userCreateBadJSON2String).run(t)
	cliTest(false, false, "users", "create", userCreateInputString).run(t)
	cliTest(false, false, "users", "create", userCreateFredInputString).run(t)
	cliTest(false, false, "users", "destroy", userCreateFredInputString).run(t)
	cliTest(false, true, "users", "create", userCreateInputString).run(t)
	cliTest(false, false, "users", "list").run(t)
	cliTest(false, false, "users", "list", "Name=fred").run(t)
	cliTest(false, false, "users", "list", "Name=john").run(t)
	cliTest(true, true, "users", "show").run(t)
	cliTest(true, true, "users", "show", "john", "john2").run(t)
	cliTest(false, true, "users", "show", "ignore").run(t)
	cliTest(false, false, "users", "show", "john").run(t)
	cliTest(true, true, "users", "exists").run(t)
	cliTest(true, true, "users", "exists", "john", "john2").run(t)
	cliTest(false, false, "users", "exists", "john").run(t)
	cliTest(false, true, "users", "exists", "ignore").run(t)
	cliTest(true, true, "users", "update").run(t)
	cliTest(true, true, "users", "update", "john", "john2", "john3").run(t)
	cliTest(false, true, "users", "update", "john", userUpdateBadJSONString).run(t)
	cliTest(false, true, "users", "update", "john2", userUpdateInputString).run(t)
	cliTest(false, false, "users", "show", "john").run(t)
	cliTest(false, false, "users", "show", "john").run(t)
	cliTest(true, true, "users", "password").run(t)
	cliTest(true, true, "users", "password", "one").run(t)
	cliTest(true, true, "users", "password", "one", "two", "three").run(t)
	cliTest(false, true, "users", "password", "jill", "june").run(t)
	cliTest(false, false, "users", "password", "john", "june").run(t)
	cliTest(true, true, "users", "destroy").run(t)
	cliTest(true, true, "users", "destroy", "john", "june").run(t)
	cliTest(false, false, "users", "destroy", "john").run(t)
	cliTest(false, true, "users", "destroy", "john").run(t)
	cliTest(false, false, "users", "list").run(t)
	cliTest(false, false, "users", "create", "-").Stdin(userCreateInputString + "\n").run(t)
	cliTest(false, false, "users", "list").run(t)
	cliTest(false, false, "users", "show", "john").run(t)
	cliTest(false, false, "users", "destroy", "john").run(t)
	cliTest(false, false, "users", "list").run(t)
	cliTest(true, true, "users", "token").run(t)
	cliTest(true, true, "users", "token", "greg", "greg2").run(t)
	cliTest(true, true, "users", "token", "greg", "greg2", "greg3").run(t)
	cliTest(false, true, "users", "token", "greg").run(t)
	cliTest(false, false, "users", "token", "rocketskates").run(t)
	cliTest(false, false, "users", "token", "rocketskates", "scope", "all", "ttl", "330", "action", "list", "specific", "asdgag").run(t)
	cliTest(false, true, "users", "token", "rocketskates", "scope", "all", "ttl", "cow", "action", "list", "specific", "asdgag").run(t)
}
