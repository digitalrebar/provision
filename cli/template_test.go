package cli

import (
	"encoding/json"
	"io/ioutil"
	"strings"
	"testing"
)

func TestTemplateCli(t *testing.T) {

	var templateCreateBadJSONString = "asdgasdg"
	var templateCreateInputString string = `{
  "Contents": "John Rules",
  "ID": "john"
}
`

	var templateUpdateBadJSONString = "asdgasdg"
	var templateUpdateInputString string = `{
  "Description": "NewStrat"
}
`

	var templatesUploadSuccessString string = `{
  "Available": true,
  "Contents": *REPLACE_WITH_TEMPLATE_GO_CONTENT*,
  "Errors": [],
  "ID": "greg",
  "ReadOnly": false,
  "Validated": true
}
`
	var templatesUploadReplaceSuccessString string = `{
  "Available": true,
  "Contents": *REPLACE_WITH_LEASE_GO_CONTENT*,
  "Errors": [],
  "ID": "greg",
  "ReadOnly": false,
  "Validated": true
}
`

	templateContent, _ := ioutil.ReadFile("template.go")
	sb, _ := json.Marshal(string(templateContent))
	templatesUploadSuccessString = strings.Replace(templatesUploadSuccessString, "*REPLACE_WITH_TEMPLATE_GO_CONTENT*", string(sb), 1)

	templateContent, _ = ioutil.ReadFile("lease.go")
	sb, _ = json.Marshal(string(templateContent))
	templatesUploadReplaceSuccessString = strings.Replace(templatesUploadReplaceSuccessString, "*REPLACE_WITH_LEASE_GO_CONTENT*", string(sb), 1)

	cliTest(true, false, "templates").run(t)
	cliTest(false, false, "templates", "list").run(t)
	cliTest(true, true, "templates", "create").run(t)
	cliTest(true, true, "templates", "create", "john", "john2").run(t)
	cliTest(false, true, "templates", "create", templateCreateBadJSONString).run(t)
	cliTest(false, false, "templates", "create", templateCreateInputString).run(t)
	cliTest(false, true, "templates", "create", templateCreateInputString).run(t)
	cliTest(false, false, "templates", "list").run(t)
	cliTest(false, false, "templates", "list", "ID=fred").run(t)
	cliTest(false, false, "templates", "list", "ID=john").run(t)
	cliTest(true, true, "templates", "show").run(t)
	cliTest(true, true, "templates", "show", "john", "john2").run(t)
	cliTest(false, true, "templates", "show", "ignore").run(t)
	cliTest(false, false, "templates", "show", "john").run(t)
	cliTest(true, true, "templates", "exists").run(t)
	cliTest(true, true, "templates", "exists", "john", "john2").run(t)
	cliTest(false, false, "templates", "exists", "john").run(t)
	cliTest(false, true, "templates", "exists", "ignore").run(t)
	cliTest(true, true, "templates", "exists", "john", "john2").run(t)
	cliTest(true, true, "templates", "update").run(t)
	cliTest(true, true, "templates", "update", "john", "john2", "john3").run(t)
	cliTest(false, true, "templates", "update", "john", templateUpdateBadJSONString).run(t)
	cliTest(false, false, "templates", "update", "john", templateUpdateInputString).run(t)
	cliTest(false, true, "templates", "update", "john2", templateUpdateInputString).run(t)
	cliTest(false, false, "templates", "show", "john").run(t)
	cliTest(false, false, "templates", "show", "john").run(t)
	cliTest(true, true, "templates", "destroy").run(t)
	cliTest(true, true, "templates", "destroy", "john", "june").run(t)
	cliTest(false, false, "templates", "destroy", "john").run(t)
	cliTest(false, true, "templates", "destroy", "john").run(t)
	cliTest(false, false, "templates", "list").run(t)
	cliTest(false, false, "templates", "create", "-").Stdin(templateCreateInputString + "\n").run(t)
	cliTest(false, false, "templates", "list").run(t)
	cliTest(false, false, "templates", "update", "john", "-").Stdin(templateUpdateInputString + "\n").run(t)
	cliTest(false, false, "templates", "show", "john").run(t)
	cliTest(false, false, "templates", "destroy", "john").run(t)
	cliTest(false, false, "templates", "list").run(t)
	cliTest(true, true, "templates", "upload").run(t)
	cliTest(true, true, "templates", "upload", "asg").run(t)
	cliTest(true, true, "templates", "upload", "asg", "two", "three", "four").run(t)
	cliTest(false, true, "templates", "upload", "greg", "as", "greg").run(t)
	cliTest(false, false, "templates", "upload", "template.go", "as", "greg").run(t)
	cliTest(false, false, "templates", "upload", "template.go", "as", "greg").run(t)
	cliTest(false, false, "templates", "upload", "lease.go", "as", "greg").run(t)
	cliTest(false, false, "templates", "destroy", "greg").run(t)
	cliTest(false, false, "templates", "exists", "etc").run(t)
	cliTest(false, false, "templates", "exists", "usrshare").run(t)
	cliTest(false, true, "templates", "destroy", "etc").run(t)
	cliTest(false, true, "templates", "destroy", "usrshare").run(t)
	verifyClean(t)
}
