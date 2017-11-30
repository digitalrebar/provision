package cli

import (
	"encoding/json"
	"io/ioutil"
	"strings"
	"testing"
)

var templateShowMissingArgErrorString string = "Error: GET: templates/ignore: Not Found\n\n"
var templateExistsMissingIgnoreString string = "Error: GET: templates/ignore: Not Found\n\n"
var templateCreateDuplicateErrorString = "Error: CREATE: templates/john: already exists\n\n"
var templateUpdateJohnMissingErrorString string = "Error: GET: templates/john2: Not Found\n\n"
var templatePatchJohnMissingErrorString string = "Error: PATCH: templates/john2: Not Found\n\n"
var templateDestroyMissingJohnString string = "Error: DELETE: templates/john: Not Found\n\n"

var templateEmptyListString string = "[]\n"

var templateDefaultListString string = `[
  {
    "Available": true,
    "Contents": "etc\n",
    "Description": "A test template for LocalStore testing",
    "Errors": [],
    "ID": "etc",
    "ReadOnly": true,
    "Validated": true
  },
  {
    "Available": true,
    "Contents": "usrshare\n",
    "Description": "A test template for DefaultStore testing",
    "Errors": [],
    "ID": "usrshare",
    "ReadOnly": true,
    "Validated": true
  }
]
`

var templateShowNoArgErrorString string = "Error: drpcli templates show [id] [flags] requires 1 argument\n"
var templateShowTooManyArgErrorString string = "Error: drpcli templates show [id] [flags] requires 1 argument\n"

var templateShowJohnString string = `{
  "Available": true,
  "Contents": "John Rules",
  "Errors": [],
  "ID": "john",
  "ReadOnly": false,
  "Validated": true
}
`

var templateExistsNoArgErrorString string = "Error: drpcli templates exists [id] [flags] requires 1 argument"
var templateExistsTooManyArgErrorString string = "Error: drpcli templates exists [id] [flags] requires 1 argument"
var templateExistsIgnoreString string = ""

var templateCreateNoArgErrorString string = "Error: drpcli templates create [json] [flags] requires 1 argument\n"
var templateCreateTooManyArgErrorString string = "Error: drpcli templates create [json] [flags] requires 1 argument\n"
var templateCreateBadJSONString = "asdgasdg"
var templateCreateBadJSONErrorString = "Error: Unable to create new template: Invalid type passed to template create\n\n"
var templateCreateInputString string = `{
  "Contents": "John Rules",
  "ID": "john"
}
`
var templateCreateJohnString string = `{
  "Available": true,
  "Contents": "John Rules",
  "Errors": [],
  "ID": "john",
  "ReadOnly": false,
  "Validated": true
}
`

var templateListJohnOnlyString = `[
  {
    "Available": true,
    "Contents": "John Rules",
    "Errors": [],
    "ID": "john",
    "ReadOnly": false,
    "Validated": true
  }
]
`
var templateListBothEnvsString = `[
  {
    "Available": true,
    "Contents": "etc\n",
    "Description": "A test template for LocalStore testing",
    "Errors": [],
    "ID": "etc",
    "ReadOnly": true,
    "Validated": true
  },
  {
    "Available": true,
    "Contents": "John Rules",
    "Errors": [],
    "ID": "john",
    "ReadOnly": false,
    "Validated": true
  },
  {
    "Available": true,
    "Contents": "usrshare\n",
    "Description": "A test template for DefaultStore testing",
    "Errors": [],
    "ID": "usrshare",
    "ReadOnly": true,
    "Validated": true
  }
]
`

var templateUpdateNoArgErrorString string = "Error: drpcli templates update [id] [json] [flags] requires 2 arguments"
var templateUpdateTooManyArgErrorString string = "Error: drpcli templates update [id] [json] [flags] requires 2 arguments"
var templateUpdateBadJSONString = "asdgasdg"
var templateUpdateBadJSONErrorString = "Error: Unable to merge objects: json: cannot unmarshal string into Go value of type map[string]interface {}\n\n\n"
var templateUpdateInputString string = `{
  "Description": "NewStrat"
}
`
var templateUpdateJohnString string = `{
  "Available": true,
  "Contents": "John Rules",
  "Description": "NewStrat",
  "Errors": [],
  "ID": "john",
  "ReadOnly": false,
  "Validated": true
}
`

var templatePatchNoArgErrorString string = "Error: drpcli templates patch [objectJson] [changesJson] [flags] requires 2 arguments"
var templatePatchTooManyArgErrorString string = "Error: drpcli templates patch [objectJson] [changesJson] [flags] requires 2 arguments"
var templatePatchBadPatchJSONString = "asdgasdg"
var templatePatchBadPatchJSONErrorString = "Error: Unable to parse drpcli templates patch [objectJson] [changesJson] [flags] JSON asdgasdg\nError: error unmarshaling JSON: json: cannot unmarshal string into Go value of type genmodels.Template\n\n"
var templatePatchBadBaseJSONString = "asdgasdg"
var templatePatchBadBaseJSONErrorString = "Error: Unable to parse drpcli templates patch [objectJson] [changesJson] [flags] JSON asdgasdg\nError: error unmarshaling JSON: json: cannot unmarshal string into Go value of type genmodels.Template\n\n"
var templatePatchBaseString string = `{
  "Available": true,
  "Contents": "John Rules",
  "Description": "NewStrat",
  "Errors": [],
  "ID": "john",
  "Validated": true
}
`
var templatePatchInputString string = `{
  "Description": "bootx64.efi"
}
`
var templatePatchJohnString string = `{
  "Available": true,
  "Contents": "John Rules",
  "Description": "bootx64.efi",
  "Errors": [],
  "ID": "john",
  "ReadOnly": false,
  "Validated": true
}
`
var templatePatchMissingBaseString string = `{
  "Contents": "John Rules",
  "Description": "NewStrat",
  "ID": "john2"
}
`

var templateDestroyNoArgErrorString string = "Error: drpcli templates destroy [id] [flags] requires 1 argument"
var templateDestroyTooManyArgErrorString string = "Error: drpcli templates destroy [id] [flags] requires 1 argument"
var templateDestroyJohnString string = "Deleted template john\n"

var templatesUploadNoArgsErrorString string = "Error: Wrong number of args: expected 3, got 0\n"
var templatesUploadOneArgsErrorString string = "Error: Wrong number of args: expected 3, got 1\n"
var templatesUploadFourArgsErrorString string = "Error: Wrong number of args: expected 3, got 4\n"
var templatesUploadMissingFileErrorString string = "Error: Failed to open greg: open greg: no such file or directory\n\n"
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
var templateDestroyGregString string = "Deleted template greg\n"

var templateReadOnlyTrueString string = `[
  {
    "Available": true,
    "Contents": "etc\n",
    "Description": "A test template for LocalStore testing",
    "Errors": [],
    "ID": "etc",
    "ReadOnly": true,
    "Validated": true
  },
  {
    "Available": true,
    "Contents": "usrshare\n",
    "Description": "A test template for DefaultStore testing",
    "Errors": [],
    "ID": "usrshare",
    "ReadOnly": true,
    "Validated": true
  }
]
`

var templateReadOnlyFalseString string = `[
  {
    "Available": true,
    "Contents": "John Rules",
    "Errors": [],
    "ID": "john",
    "ReadOnly": false,
    "Validated": true
  }
]
`

func TestTemplateCli(t *testing.T) {
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
}
