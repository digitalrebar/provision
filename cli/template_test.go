package cli

import (
	"encoding/json"
	"io/ioutil"
	"strings"
	"testing"
)

var templateEmptyListString string = "[]\n"

var templateDefaultListString string = `[
  {
    "Contents": "etc\n",
    "Description": "A test template for LocalStore testing",
    "ID": "etc"
  },
  {
    "Contents": "exit\n",
    "ID": "local-elilo.tmpl"
  },
  {
    "Contents": "#!ipxe\nexit\n",
    "ID": "local-ipxe.tmpl"
  },
  {
    "Contents": "DEFAULT local\nPROMPT 0\nTIMEOUT 10\nLABEL local\nlocalboot 0\n",
    "ID": "local-pxelinux.tmpl"
  },
  {
    "Contents": "usrshare\n",
    "Description": "A test template for DefaultStore testing",
    "ID": "usrshare"
  }
]
`

var templateShowNoArgErrorString string = "Error: drpcli templates show [id] requires 1 argument\n"
var templateShowTooManyArgErrorString string = "Error: drpcli templates show [id] requires 1 argument\n"
var templateShowMissingArgErrorString string = "Error: templates GET: ignore: Not Found\n\n"
var templateShowJohnString string = `{
  "Contents": "John Rules",
  "ID": "john"
}
`

var templateExistsNoArgErrorString string = "Error: drpcli templates exists [id] requires 1 argument"
var templateExistsTooManyArgErrorString string = "Error: drpcli templates exists [id] requires 1 argument"
var templateExistsIgnoreString string = ""
var templateExistsMissingIgnoreString string = "Error: templates GET: ignore: Not Found\n\n"

var templateCreateNoArgErrorString string = "Error: drpcli templates create [json] requires 1 argument\n"
var templateCreateTooManyArgErrorString string = "Error: drpcli templates create [json] requires 1 argument\n"
var templateCreateBadJSONString = "asdgasdg"
var templateCreateBadJSONErrorString = "Error: Unable to create new template: Invalid type passed to template create\n\n"
var templateCreateInputString string = `{
  "Contents": "John Rules",
  "ID": "john"
}
`
var templateCreateJohnString string = `{
  "Contents": "John Rules",
  "ID": "john"
}
`
var templateCreateDuplicateErrorString = "Error: dataTracker create templates: john already exists\n\n"

var templateListJohnOnlyString = `[
  {
    "Contents": "John Rules",
    "ID": "john"
  }
]
`
var templateListBothEnvsString = `[
  {
    "Contents": "etc\n",
    "Description": "A test template for LocalStore testing",
    "ID": "etc"
  },
  {
    "Contents": "John Rules",
    "ID": "john"
  },
  {
    "Contents": "exit\n",
    "ID": "local-elilo.tmpl"
  },
  {
    "Contents": "#!ipxe\nexit\n",
    "ID": "local-ipxe.tmpl"
  },
  {
    "Contents": "DEFAULT local\nPROMPT 0\nTIMEOUT 10\nLABEL local\nlocalboot 0\n",
    "ID": "local-pxelinux.tmpl"
  },
  {
    "Contents": "usrshare\n",
    "Description": "A test template for DefaultStore testing",
    "ID": "usrshare"
  }
]
`

var templateUpdateNoArgErrorString string = "Error: drpcli templates update [id] [json] requires 2 arguments"
var templateUpdateTooManyArgErrorString string = "Error: drpcli templates update [id] [json] requires 2 arguments"
var templateUpdateBadJSONString = "asdgasdg"
var templateUpdateBadJSONErrorString = "Error: Unable to merge objects: json: cannot unmarshal string into Go value of type map[string]interface {}\n\n\n"
var templateUpdateInputString string = `{
  "Description": "NewStrat"
}
`
var templateUpdateJohnString string = `{
  "Contents": "John Rules",
  "Description": "NewStrat",
  "ID": "john"
}
`
var templateUpdateJohnMissingErrorString string = "Error: templates GET: john2: Not Found\n\n"

var templatePatchNoArgErrorString string = "Error: drpcli templates patch [objectJson] [changesJson] requires 2 arguments"
var templatePatchTooManyArgErrorString string = "Error: drpcli templates patch [objectJson] [changesJson] requires 2 arguments"
var templatePatchBadPatchJSONString = "asdgasdg"
var templatePatchBadPatchJSONErrorString = "Error: Unable to parse drpcli templates patch [objectJson] [changesJson] JSON asdgasdg\nError: error unmarshaling JSON: json: cannot unmarshal string into Go value of type models.Template\n\n"
var templatePatchBadBaseJSONString = "asdgasdg"
var templatePatchBadBaseJSONErrorString = "Error: Unable to parse drpcli templates patch [objectJson] [changesJson] JSON asdgasdg\nError: error unmarshaling JSON: json: cannot unmarshal string into Go value of type models.Template\n\n"
var templatePatchBaseString string = `{
  "Contents": "John Rules",
  "Description": "NewStrat",
  "ID": "john"
}
`
var templatePatchInputString string = `{
  "Description": "bootx64.efi"
}
`
var templatePatchJohnString string = `{
  "Contents": "John Rules",
  "Description": "bootx64.efi",
  "ID": "john"
}
`
var templatePatchMissingBaseString string = `{
  "Contents": "John Rules",
  "Description": "NewStrat",
  "ID": "john2"
}
`
var templatePatchJohnMissingErrorString string = "Error: templates: PATCH john2: Not Found\n\n"

var templateDestroyNoArgErrorString string = "Error: drpcli templates destroy [id] requires 1 argument"
var templateDestroyTooManyArgErrorString string = "Error: drpcli templates destroy [id] requires 1 argument"
var templateDestroyJohnString string = "Deleted template john\n"
var templateDestroyMissingJohnString string = "Error: templates: DELETE john: Not Found\n\n"

var templatesUploadNoArgsErrorString string = "Error: Wrong number of args: expected 3, got 0\n"
var templatesUploadOneArgsErrorString string = "Error: Wrong number of args: expected 3, got 1\n"
var templatesUploadFourArgsErrorString string = "Error: Wrong number of args: expected 3, got 4\n"
var templatesUploadMissingFileErrorString string = "Error: Failed to open greg: open greg: no such file or directory\n\n"
var templatesUploadSuccessString string = `{
  "Contents": *REPLACE_WITH_TEMPLATE_GO_CONTENT*,
  "ID": "greg"
}
`
var templatesUploadReplaceSuccessString string = `{
  "Contents": *REPLACE_WITH_LEASE_GO_CONTENT*,
  "ID": "greg"
}
`
var templateDestroyGregString string = "Deleted template greg\n"

func TestTemplateCli(t *testing.T) {
	templateContent, _ := ioutil.ReadFile("template.go")
	sb, _ := json.Marshal(string(templateContent))
	templatesUploadSuccessString = strings.Replace(templatesUploadSuccessString, "*REPLACE_WITH_TEMPLATE_GO_CONTENT*", string(sb), 1)

	templateContent, _ = ioutil.ReadFile("lease.go")
	sb, _ = json.Marshal(string(templateContent))
	templatesUploadReplaceSuccessString = strings.Replace(templatesUploadReplaceSuccessString, "*REPLACE_WITH_LEASE_GO_CONTENT*", string(sb), 1)

	tests := []CliTest{
		CliTest{true, false, []string{"templates"}, noStdinString, "Access CLI commands relating to templates\n", ""},
		CliTest{false, false, []string{"templates", "list"}, noStdinString, templateDefaultListString, noErrorString},

		CliTest{true, true, []string{"templates", "create"}, noStdinString, noContentString, templateCreateNoArgErrorString},
		CliTest{true, true, []string{"templates", "create", "john", "john2"}, noStdinString, noContentString, templateCreateTooManyArgErrorString},
		CliTest{false, true, []string{"templates", "create", templateCreateBadJSONString}, noStdinString, noContentString, templateCreateBadJSONErrorString},
		CliTest{false, false, []string{"templates", "create", templateCreateInputString}, noStdinString, templateCreateJohnString, noErrorString},
		CliTest{false, true, []string{"templates", "create", templateCreateInputString}, noStdinString, noContentString, templateCreateDuplicateErrorString},
		CliTest{false, false, []string{"templates", "list"}, noStdinString, templateListBothEnvsString, noErrorString},
		CliTest{false, false, []string{"templates", "list", "--limit=0"}, noStdinString, templateEmptyListString, noErrorString},
		CliTest{false, false, []string{"templates", "list", "--limit=10", "--offset=0"}, noStdinString, templateListBothEnvsString, noErrorString},
		CliTest{false, false, []string{"templates", "list", "--limit=10", "--offset=10"}, noStdinString, templateEmptyListString, noErrorString},
		CliTest{false, true, []string{"templates", "list", "--limit=-10", "--offset=0"}, noStdinString, noContentString, limitNegativeError},
		CliTest{false, true, []string{"templates", "list", "--limit=10", "--offset=-10"}, noStdinString, noContentString, offsetNegativeError},
		CliTest{false, false, []string{"templates", "list", "--limit=-1", "--offset=-1"}, noStdinString, templateListBothEnvsString, noErrorString},
		CliTest{false, false, []string{"templates", "list", "ID=fred"}, noStdinString, templateEmptyListString, noErrorString},
		CliTest{false, false, []string{"templates", "list", "ID=john"}, noStdinString, templateListJohnOnlyString, noErrorString},

		CliTest{true, true, []string{"templates", "show"}, noStdinString, noContentString, templateShowNoArgErrorString},
		CliTest{true, true, []string{"templates", "show", "john", "john2"}, noStdinString, noContentString, templateShowTooManyArgErrorString},
		CliTest{false, true, []string{"templates", "show", "ignore"}, noStdinString, noContentString, templateShowMissingArgErrorString},
		CliTest{false, false, []string{"templates", "show", "john"}, noStdinString, templateShowJohnString, noErrorString},

		CliTest{true, true, []string{"templates", "exists"}, noStdinString, noContentString, templateExistsNoArgErrorString},
		CliTest{true, true, []string{"templates", "exists", "john", "john2"}, noStdinString, noContentString, templateExistsTooManyArgErrorString},
		CliTest{false, false, []string{"templates", "exists", "john"}, noStdinString, templateExistsIgnoreString, noErrorString},
		CliTest{false, true, []string{"templates", "exists", "ignore"}, noStdinString, noContentString, templateExistsMissingIgnoreString},
		CliTest{true, true, []string{"templates", "exists", "john", "john2"}, noStdinString, noContentString, templateExistsTooManyArgErrorString},

		CliTest{true, true, []string{"templates", "update"}, noStdinString, noContentString, templateUpdateNoArgErrorString},
		CliTest{true, true, []string{"templates", "update", "john", "john2", "john3"}, noStdinString, noContentString, templateUpdateTooManyArgErrorString},
		CliTest{false, true, []string{"templates", "update", "john", templateUpdateBadJSONString}, noStdinString, noContentString, templateUpdateBadJSONErrorString},
		CliTest{false, false, []string{"templates", "update", "john", templateUpdateInputString}, noStdinString, templateUpdateJohnString, noErrorString},
		CliTest{false, true, []string{"templates", "update", "john2", templateUpdateInputString}, noStdinString, noContentString, templateUpdateJohnMissingErrorString},
		CliTest{false, false, []string{"templates", "show", "john"}, noStdinString, templateUpdateJohnString, noErrorString},

		CliTest{true, true, []string{"templates", "patch"}, noStdinString, noContentString, templatePatchNoArgErrorString},
		CliTest{true, true, []string{"templates", "patch", "john", "john2", "john3"}, noStdinString, noContentString, templatePatchTooManyArgErrorString},
		CliTest{false, true, []string{"templates", "patch", templatePatchBaseString, templatePatchBadPatchJSONString}, noStdinString, noContentString, templatePatchBadPatchJSONErrorString},
		CliTest{false, true, []string{"templates", "patch", templatePatchBadBaseJSONString, templatePatchInputString}, noStdinString, noContentString, templatePatchBadBaseJSONErrorString},
		CliTest{false, false, []string{"templates", "patch", templatePatchBaseString, templatePatchInputString}, noStdinString, templatePatchJohnString, noErrorString},
		CliTest{false, true, []string{"templates", "patch", templatePatchMissingBaseString, templatePatchInputString}, noStdinString, noContentString, templatePatchJohnMissingErrorString},
		CliTest{false, false, []string{"templates", "show", "john"}, noStdinString, templatePatchJohnString, noErrorString},

		CliTest{true, true, []string{"templates", "destroy"}, noStdinString, noContentString, templateDestroyNoArgErrorString},
		CliTest{true, true, []string{"templates", "destroy", "john", "june"}, noStdinString, noContentString, templateDestroyTooManyArgErrorString},
		CliTest{false, false, []string{"templates", "destroy", "john"}, noStdinString, templateDestroyJohnString, noErrorString},
		CliTest{false, true, []string{"templates", "destroy", "john"}, noStdinString, noContentString, templateDestroyMissingJohnString},
		CliTest{false, false, []string{"templates", "list"}, noStdinString, templateDefaultListString, noErrorString},

		CliTest{false, false, []string{"templates", "create", "-"}, templateCreateInputString + "\n", templateCreateJohnString, noErrorString},
		CliTest{false, false, []string{"templates", "list"}, noStdinString, templateListBothEnvsString, noErrorString},
		CliTest{false, false, []string{"templates", "update", "john", "-"}, templateUpdateInputString + "\n", templateUpdateJohnString, noErrorString},
		CliTest{false, false, []string{"templates", "show", "john"}, noStdinString, templateUpdateJohnString, noErrorString},

		CliTest{false, false, []string{"templates", "destroy", "john"}, noStdinString, templateDestroyJohnString, noErrorString},
		CliTest{false, false, []string{"templates", "list"}, noStdinString, templateDefaultListString, noErrorString},

		CliTest{true, true, []string{"templates", "upload"}, noStdinString, noContentString, templatesUploadNoArgsErrorString},
		CliTest{true, true, []string{"templates", "upload", "asg"}, noStdinString, noContentString, templatesUploadOneArgsErrorString},
		CliTest{true, true, []string{"templates", "upload", "asg", "two", "three", "four"}, noStdinString, noContentString, templatesUploadFourArgsErrorString},
		CliTest{false, true, []string{"templates", "upload", "greg", "as", "greg"}, noStdinString, noContentString, templatesUploadMissingFileErrorString},
		CliTest{false, false, []string{"templates", "upload", "template.go", "as", "greg"}, noStdinString, templatesUploadSuccessString, noErrorString},
		CliTest{false, false, []string{"templates", "upload", "template.go", "as", "greg"}, noStdinString, templatesUploadSuccessString, noErrorString},
		CliTest{false, false, []string{"templates", "upload", "lease.go", "as", "greg"}, noStdinString, templatesUploadReplaceSuccessString, noErrorString},
		CliTest{false, false, []string{"templates", "destroy", "greg"}, noStdinString, templateDestroyGregString, noErrorString},
		CliTest{false, false, []string{"templates", "exists", "etc"}, noStdinString, noContentString, noErrorString},
		CliTest{false, false, []string{"templates", "exists", "usrshare"}, noStdinString, noContentString, noErrorString},
		CliTest{false, true, []string{"templates", "destroy", "etc"}, noStdinString, noContentString, "Error: readonly: etc\n\n"},
		CliTest{false, true, []string{"templates", "destroy", "usrshare"}, noStdinString, noContentString, "Error: readonly: usrshare\n\n"},
	}

	for _, test := range tests {
		testCli(t, test)
	}
}
