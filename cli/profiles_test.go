package cli

import (
	"fmt"
	"testing"
)

func TestProfileCli(t *testing.T) {

	var profileCreateBadJSONString = "{asdgasdg"

	var profileCreateBadJSON2String = "[asdgasdg]"
	var profileCreateInputString string = `{
  "Name": "john",
  "Params": {
    "FRED": "GREG"
  }
}
`
	var profileUpdateBadJSONString = "asdgasdg"

	var profileUpdateInputString string = `{
  "Params": {
    "JESSIE": "JAMES"
  }
}
`
	var profilesParamsNextString string = `{
  "jj": 3
}
`

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
	cliTest(true, true, "profiles", "add").run(t)
	cliTest(true, true, "profiles", "add", "john2").run(t)
	cliTest(true, true, "profiles", "add", "john2", "extra").run(t)
	cliTest(false, false, "profiles", "add", "john", "param", "newparam", "to", "toast").run(t)
	cliTest(false, true, "profiles", "add", "john", "param", "newparam", "to", "toast").run(t)
	cliTest(false, false, "profiles", "get", "john", "param", "newparam").run(t)
	cliTest(false, false, "profiles", "add", "john", "param", "newparam2", "to", "toast").run(t)
	cliTest(true, true, "profiles", "remove").run(t)
	cliTest(true, true, "profiles", "remove", "john2").run(t)
	cliTest(true, true, "profiles", "remove", "john2", "extra").run(t)
	cliTest(false, false, "profiles", "remove", "john", "param", "newparam").run(t)
	cliTest(false, false, "profiles", "params", "john").run(t)
	cliTest(true, true, "profiles", "params", "john", "--params").run(t)
	cliTest(false, false, "profiles", "params", "john", "--params", "FRED,JESSIE,JOHN").run(t)
	cliTest(false, true, "profiles", "remove", "john", "param", "newparam").run(t)
	cliTest(false, true, "profiles", "remove", "john", "param", "newparam2", "--ref", "bagel").run(t)
	cliTest(false, false, "profiles", "remove", "john", "param", "newparam2", "--ref", "toast").run(t)
	cliTest(false, false, "profiles", "params", "john").run(t)
	cliTest(true, true, "profiles", "set").run(t)
	cliTest(false, true, "profiles", "set", "john2", "param", "john2", "to", "cow").run(t)
	cliTest(false, true, "profiles", "set", "john", "param", "john2", "to", "cow", "--ref", "fred").run(t)
	cliTest(false, false, "profiles", "set", "john", "param", "john2", "to", "cow").run(t)
	cliTest(false, false, "profiles", "get", "john", "param", "john2").run(t)
	cliTest(false, false, "profiles", "set", "john", "param", "john2", "to", "cow", "--ref", "cow").run(t)
	cliTest(false, false, "profiles", "get", "john", "param", "john2").run(t)
	cliTest(false, false, "profiles", "set", "john", "param", "john2", "to", "sow", "--ref", "cow").run(t)
	cliTest(false, false, "profiles", "get", "john", "param", "john2").run(t)
	cliTest(false, false, "profiles", "set", "john", "param", "john2", "to", "cow", "--ref", "sow").run(t)
	cliTest(false, false, "profiles", "get", "john", "param", "john2").run(t)
	cliTest(false, false, "profiles", "set", "john", "param", "john2", "to", "3").run(t)
	cliTest(false, false, "profiles", "set", "john", "param", "john3", "to", "4").run(t)
	cliTest(false, false, "profiles", "get", "john", "param", "john2").run(t)
	cliTest(false, false, "profiles", "get", "john", "param", "john3").run(t)
	cliTest(false, false, "profiles", "set", "john", "param", "john2", "to", "null").run(t)
	cliTest(false, false, "profiles", "get", "john", "param", "john2").run(t)
	cliTest(false, false, "profiles", "get", "john", "param", "john3").run(t)
	cliTest(false, false, "profiles", "show", "john", "--slim", "params,meta").run(t)
	cliTest(false, false, "profiles", "list", "--slim", "params,meta").run(t)
	cliTest(true, true, "profiles", "show", "john", "--slim", "params,meta", "--params").run(t)
	cliTest(false, false, "profiles", "show", "john", "--slim", "params,meta", "--params", "FREDJESSIEJOHN").run(t)
	cliTest(false, false, "profiles", "show", "john", "--slim", "params,meta", "--params", "FRED,JESSIE,JOHN").run(t)
	cliTest(false, false, "profiles", "show", "john", "--slim", "meta", "--params", "FRED,JESSIE,JOHN").run(t)
	cliTest(false, false, "profiles", "show", "john", "--params", "FRED,JESSIE,JOHN").run(t)
	cliTest(true, true, "profiles", "list", "--slim", "params,meta", "--params").run(t)
	cliTest(false, false, "profiles", "list", "--slim", "params,meta", "--params", "FREDJESSIEJOHN").run(t)
	cliTest(false, false, "profiles", "list", "--slim", "params,meta", "--params", "FRED,JESSIE,JOHN").run(t)
	cliTest(false, false, "profiles", "list", "--slim", "meta", "--params", "FRED,JESSIE,JOHN").run(t)
	cliTest(false, false, "profiles", "list", "--params", "FRED,JESSIE,JOHN").run(t)
	cliTest(true, true, "profiles", "params").run(t)
	cliTest(false, true, "profiles", "params", "john2").run(t)
	cliTest(false, false, "profiles", "params", "john").run(t)
	cliTest(false, true, "profiles", "params", "john2", profilesParamsNextString).run(t)
	cliTest(false, true, "profiles", "params", "john", profilesParamsNextString, "--ref", profilesParamsNextString).run(t)
	cliTest(false, false, "profiles", "params", "john", profilesParamsNextString).run(t)
	cliTest(false, false, "profiles", "params", "john").run(t)
	cliTest(false, false, "profiles", "params", "john", profilesParamsNextString, "--ref", profilesParamsNextString).run(t)
	cliTest(false, false, "profiles", "params", "john").run(t)
	cliTest(false, false, "profiles", "params", "john", "{}", "--ref", profilesParamsNextString).run(t)
	cliTest(false, false, "profiles", "params", "john").run(t)
	cliTest(false, false, "profiles", "params", "john", profilesParamsNextString).run(t)
	cliTest(false, false, "profiles", "params", "john").run(t)
	cliTest(false, false, "profiles", "show", "john").run(t)
	cliTest(false, false, "profiles", "destroy", "john").run(t)
	cliTest(false, false, "profiles", "list").run(t)
	verifyClean(t)
}

func TestProfileFileImport(t *testing.T) {
	prefix := "profiles"
	yamlId := "yamltest"
	jsonId := "jsontest"

	yamlCreate := fmt.Sprintf("test-data/base/%s/create.yaml", prefix)
	jsonCreate := fmt.Sprintf("test-data/base/%s/create.json", prefix)
	yamlBad := fmt.Sprintf("test-data/base/%s/bad.yaml", prefix)
	jsonBad := fmt.Sprintf("test-data/base/%s/bad.json", prefix)
	yamlUpdate := fmt.Sprintf("test-data/base/%s/update.yaml", prefix)
	jsonUpdate := fmt.Sprintf("test-data/base/%s/update.json", prefix)

	cliTest(false, false, prefix, "create", yamlCreate).run(t)
	cliTest(false, false, prefix, "create", jsonCreate).run(t)
	cliTest(false, true, prefix, "create", yamlBad).run(t)
	cliTest(false, true, prefix, "create", jsonBad).run(t)
	cliTest(false, false, prefix, "update", yamlId, yamlUpdate).run(t)
	cliTest(false, false, prefix, "update", jsonId, jsonUpdate).run(t)
	cliTest(false, true, prefix, "update", yamlId, yamlBad).run(t)
	cliTest(false, true, prefix, "update", jsonId, jsonBad).run(t)
	cliTest(false, false, prefix, "destroy", yamlId).run(t)
	cliTest(false, false, prefix, "destroy", jsonId).run(t)
}
