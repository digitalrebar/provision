package cli

import (
	"os"
	"os/exec"
	"testing"
	"time"
)

var plugin_providerShowNoArgErrorString string = "Error: drpcli plugin_providers show [id] [flags] requires 1 argument\n"
var plugin_providerShowTooManyArgErrorString string = "Error: drpcli plugin_providers show [id] [flags] requires 1 argument\n"
var plugin_providerShowMissingArgErrorString string = "Error: GET: plugin_providers/john: Not Found\n\n"

var plugin_providerExistsNoArgErrorString string = "Error: drpcli plugin_providers exists [id] [flags] requires 1 argument"
var plugin_providerExistsTooManyArgErrorString string = "Error: drpcli plugin_providers exists [id] [flags] requires 1 argument"
var plugin_providerExistsIgnoreString string = ""
var plugin_providerExistsMissingJohnString string = "Error: GET: plugin_providers/john: Not Found\n\n"

var plugin_providerListString string = `[
  {
    "AvailableActions": [
      {
        "Command": "increment",
        "OptionalParams": [
          "incrementer/step",
          "incrementer/parameter"
        ],
        "Provider": "incrementer",
        "RequiredParams": []
      },
      {
        "Command": "reset_count",
        "OptionalParams": [],
        "Provider": "incrementer",
        "RequiredParams": [
          "incrementer/touched"
        ]
      }
    ],
    "Content": "meta:\n  Description: Test Plugin for DRP\n  Name: incrementer\n  Source: Digital Rebar\n  Type: plugin\n  Version: Internal\nsections:\n  params:\n    incrementer/touched:\n      Available: false\n      Description: \"\"\n      Documentation: \"\"\n      Errors: []\n      Meta: {}\n      Name: incrementer/touched\n      ReadOnly: false\n      Schema:\n        type: integer\n      Validated: false\n",
    "HasPublish": true,
    "Name": "incrementer",
    "OptionalParams": [],
    "Parameters": [
      {
        "Available": false,
        "Errors": [],
        "Name": "incrementer/parameter",
        "ReadOnly": false,
        "Schema": {
          "type": "string"
        },
        "Validated": false
      },
      {
        "Available": false,
        "Errors": [],
        "Name": "incrementer/step",
        "ReadOnly": false,
        "Schema": {
          "type": "integer"
        },
        "Validated": false
      }
    ],
    "RequiredParams": [],
    "Version": "v3.0.2-pre-alpha-NotSet"
  }
]
`
var plugin_providerShowIncrementerString string = `{
  "AvailableActions": [
    {
      "Command": "increment",
      "OptionalParams": [
        "incrementer/step",
        "incrementer/parameter"
      ],
      "Provider": "incrementer",
      "RequiredParams": []
    },
    {
      "Command": "reset_count",
      "OptionalParams": [],
      "Provider": "incrementer",
      "RequiredParams": [
        "incrementer/touched"
      ]
    }
  ],
  "Content": "meta:\n  Description: Test Plugin for DRP\n  Name: incrementer\n  Source: Digital Rebar\n  Type: plugin\n  Version: Internal\nsections:\n  params:\n    incrementer/touched:\n      Available: false\n      Description: \"\"\n      Documentation: \"\"\n      Errors: []\n      Meta: {}\n      Name: incrementer/touched\n      ReadOnly: false\n      Schema:\n        type: integer\n      Validated: false\n",
  "HasPublish": true,
  "Name": "incrementer",
  "OptionalParams": [],
  "Parameters": [
    {
      "Available": false,
      "Errors": [],
      "Name": "incrementer/parameter",
      "ReadOnly": false,
      "Schema": {
        "type": "string"
      },
      "Validated": false
    },
    {
      "Available": false,
      "Errors": [],
      "Name": "incrementer/step",
      "ReadOnly": false,
      "Schema": {
        "type": "integer"
      },
      "Validated": false
    }
  ],
  "RequiredParams": [],
  "Version": "v3.0.2-pre-alpha-NotSet"
}
`

var plugin_providerUploadSuccessString = `RE:
{
  "path": "incrementer",
  "size": [\d]+
}
`

var plugin_providerDestroyNoArgErrorString = "Error: drpcli plugin_providers destroy [id] [flags] requires 1 argument\n"
var plugin_providerDestroyTooManyArgErrorString = "Error: drpcli plugin_providers destroy [id] [flags] requires 1 argument\n"
var plugin_providerDestroyMissingArgErrorString = `Error: DELETE: plugin_providers/john: Not Found

`
var plugin_providerUploadNoArgErrorString = "Error: Wrong number of args: expected 3, got 0\n"
var plugin_providerUploadTooFewArgErrorString = "Error: Wrong number of args: expected 3, got 1\n"
var plugin_providerUploadTooManyArgErrorString = "Error: Wrong number of args: expected 3, got 4\n"
var plugin_providerUploadMissingArgErrorString = "Error: Failed to open john: open john: no such file or directory\n\n"

var plugin_providerDestroySuccesString = "Deleted plugin_provider incrementer\n"

var plugin_providerParamString = `{
  "Available": true,
  "Errors": [],
  "Name": "incrementer/parameter",
  "ReadOnly": true,
  "Schema": {
    "type": "string"
  },
  "Validated": true
}
`
var plugin_providerMissingParamString = "Error: GET: params/incrementer/parameter: Not Found\n\n"

func TestPluginProviderCli(t *testing.T) {

	srcFolder := tmpDir + "/plugins/incrementer"
	cpCmd := exec.Command("cp", "-rf", srcFolder, "incrementer")
	err := cpCmd.Run()
	if err != nil {
		t.Errorf("Failed to copy incrementer: %v\n", err)
	}

	cliTest(true, false, "plugin_providers").run(t)
	cliTest(true, true, "plugin_providers", "show").run(t)
	cliTest(true, true, "plugin_providers", "show", "john", "john2").run(t)
	cliTest(false, true, "plugin_providers", "show", "john").run(t)
	cliTest(false, false, "plugin_providers", "show", "incrementer").run(t)
	cliTest(true, true, "plugin_providers", "exists").run(t)
	cliTest(true, true, "plugin_providers", "exists", "john", "john2").run(t)
	cliTest(false, true, "plugin_providers", "exists", "john").run(t)
	cliTest(false, false, "plugin_providers", "exists", "incrementer").run(t)
	cliTest(false, false, "plugin_providers", "list").run(t)
	cliTest(true, true, "plugin_providers", "destroy").run(t)
	cliTest(true, true, "plugin_providers", "destroy", "john", "john2").run(t)
	cliTest(false, true, "plugin_providers", "destroy", "john").run(t)
	cliTest(false, false, "params", "show", "incrementer/parameter").run(t)
	cliTest(false, false, "plugin_providers", "destroy", "incrementer").run(t)
	cliTest(false, true, "params", "show", "incrementer/parameter").run(t)
	time.Sleep(3 * time.Second)
	cliTest(false, false, "plugin_providers", "list").run(t)
	cliTest(true, true, "plugin_providers", "upload").run(t)
	cliTest(true, true, "plugin_providers", "upload", "john").run(t)
	cliTest(true, true, "plugin_providers", "upload", "john", "as", "john2", "asdga").run(t)
	cliTest(false, true, "plugin_providers", "upload", "john", "as", "john").run(t)
	cliTest(false, false, "plugin_providers", "upload", "incrementer", "as", "incrementer").run(t)
	time.Sleep(3 * time.Second)
	cliTest(false, false, "plugin_providers", "list").run(t)
	os.Remove("incrementer")
}
