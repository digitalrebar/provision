package cli

import "testing"

func TestUserCli(t *testing.T) {

	var userCreateBadJSONString = "{asdgasdg"
	var userCreateBadJSON2String = "[asdgasdg]"

	var userCreateInputString string = `{
  "Name": "john",
  "PasswordHash": null
}
`
	var userCreateFredInputString string = `fred`

	var userUpdateBadJSONString = "asdgasdg"
	var userUpdateInputString string = `{
  "PasswordHash": "NewStrat"
}
`
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
	cliTest(false, true, "-U", "john", "-P", "jack", "-T", "", "users", "list").run(t)
	cliTest(false, false, "-U", "john", "-P", "june", "-T", "", "users", "list").run(t)
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
	cliTest(true, true, "users", "passwordhash").run(t)
	cliTest(false, false, "users", "passwordhash", "fred").run(t)
	verifyClean(t)
}
