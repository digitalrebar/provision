package cli

import (
	"os/exec"
	"testing"
	"time"
)

var plugin_providerShowNoArgErrorString string = "Error: drpcli plugin_providers show [id] [flags] requires 1 argument\n"
var plugin_providerShowTooManyArgErrorString string = "Error: drpcli plugin_providers show [id] [flags] requires 1 argument\n"
var plugin_providerShowMissingArgErrorString string = "Error: plugin provider get: not found: john\n\n"

var plugin_providerExistsNoArgErrorString string = "Error: drpcli plugin_providers exists [id] [flags] requires 1 argument"
var plugin_providerExistsTooManyArgErrorString string = "Error: drpcli plugin_providers exists [id] [flags] requires 1 argument"
var plugin_providerExistsIgnoreString string = ""
var plugin_providerExistsMissingJohnString string = "Error: plugin provider get: not found: john\n\n"

var plugin_providerListString string = `[
  {
    "AvailableActions": [
      {
        "Command": "increment",
        "OptionalParams": [
          "incrementer.step",
          "incrementer.parameter"
        ],
        "Provider": "incrementer",
        "RequiredParams": null
      },
      {
        "Command": "reset_count",
        "OptionalParams": null,
        "Provider": "incrementer",
        "RequiredParams": [
          "incrementer.touched"
        ]
      }
    ],
    "HasPublish": true,
    "Name": "incrementer",
    "OptionalParams": null,
    "Parameters": [
      {
        "Available": true,
        "Errors": [],
        "Name": "incrementer.parameter",
        "ReadOnly": false,
        "Schema": {
          "type": "string"
        },
        "Validated": true
      },
      {
        "Available": true,
        "Errors": [],
        "Name": "incrementer.step",
        "ReadOnly": false,
        "Schema": {
          "type": "integer"
        },
        "Validated": true
      },
      {
        "Available": true,
        "Errors": [],
        "Name": "incrementer.touched",
        "ReadOnly": false,
        "Schema": {
          "type": "integer"
        },
        "Validated": true
      }
    ],
    "RequiredParams": null,
    "Version": "v3.0.2-pre-alpha-NotSet"
  }
]
`
var plugin_providerShowIncrementerString string = `{
  "AvailableActions": [
    {
      "Command": "increment",
      "OptionalParams": [
        "incrementer.step",
        "incrementer.parameter"
      ],
      "Provider": "incrementer",
      "RequiredParams": null
    },
    {
      "Command": "reset_count",
      "OptionalParams": null,
      "Provider": "incrementer",
      "RequiredParams": [
        "incrementer.touched"
      ]
    }
  ],
  "HasPublish": true,
  "Name": "incrementer",
  "OptionalParams": null,
  "Parameters": [
    {
      "Available": true,
      "Errors": [],
      "Name": "incrementer.parameter",
      "ReadOnly": false,
      "Schema": {
        "type": "string"
      },
      "Validated": true
    },
    {
      "Available": true,
      "Errors": [],
      "Name": "incrementer.step",
      "ReadOnly": false,
      "Schema": {
        "type": "integer"
      },
      "Validated": true
    },
    {
      "Available": true,
      "Errors": [],
      "Name": "incrementer.touched",
      "ReadOnly": false,
      "Schema": {
        "type": "integer"
      },
      "Validated": true
    }
  ],
  "RequiredParams": null,
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
var plugin_providerDestroyMissingArgErrorString = "Error: delete: unable to delete john\n\n"
var plugin_providerUploadNoArgErrorString = "Error: Wrong number of args: expected 3, got 0\n"
var plugin_providerUploadTooFewArgErrorString = "Error: Wrong number of args: expected 3, got 1\n"
var plugin_providerUploadTooManyArgErrorString = "Error: Wrong number of args: expected 3, got 4\n"
var plugin_providerUploadMissingArgErrorString = "Error: Failed to open john: open john: no such file or directory\n\n"

var plugin_providerDestroySuccesString = "Deleted plugin_provider incrementer\n"

func TestPluginProviderCli(t *testing.T) {

	srcFolder := tmpDir + "/plugins/incrementer"
	destFolder := tmpDir + "/incrementer"
	cpCmd := exec.Command("cp", "-rf", srcFolder, destFolder)
	err := cpCmd.Run()
	if err != nil {
		t.Errorf("Failed to copy incrementer: %v\n", err)
	}

	tests := []CliTest{
		CliTest{true, false, []string{"plugin_providers"}, noStdinString, "Access CLI commands relating to plugin_providers\n", ""},

		CliTest{true, true, []string{"plugin_providers", "show"}, noStdinString, noContentString, plugin_providerShowNoArgErrorString},
		CliTest{true, true, []string{"plugin_providers", "show", "john", "john2"}, noStdinString, noContentString, plugin_providerShowTooManyArgErrorString},
		CliTest{false, true, []string{"plugin_providers", "show", "john"}, noStdinString, noContentString, plugin_providerShowMissingArgErrorString},
		CliTest{false, false, []string{"plugin_providers", "show", "incrementer"}, noStdinString, plugin_providerShowIncrementerString, noErrorString},

		CliTest{true, true, []string{"plugin_providers", "exists"}, noStdinString, noContentString, plugin_providerExistsNoArgErrorString},
		CliTest{true, true, []string{"plugin_providers", "exists", "john", "john2"}, noStdinString, noContentString, plugin_providerExistsTooManyArgErrorString},
		CliTest{false, true, []string{"plugin_providers", "exists", "john"}, noStdinString, noContentString, plugin_providerExistsMissingJohnString},
		CliTest{false, false, []string{"plugin_providers", "exists", "incrementer"}, noStdinString, noContentString, noErrorString},

		CliTest{false, false, []string{"plugin_providers", "list"}, noStdinString, plugin_providerListString, noErrorString},

		CliTest{true, true, []string{"plugin_providers", "destroy"}, noStdinString, noContentString, plugin_providerDestroyNoArgErrorString},
		CliTest{true, true, []string{"plugin_providers", "destroy", "john", "john2"}, noStdinString, noContentString, plugin_providerDestroyTooManyArgErrorString},
		CliTest{false, true, []string{"plugin_providers", "destroy", "john"}, noStdinString, noContentString, plugin_providerDestroyMissingArgErrorString},
		CliTest{false, false, []string{"plugin_providers", "destroy", "incrementer"}, noStdinString, plugin_providerDestroySuccesString, noErrorString},
	}
	for _, test := range tests {
		testCli(t, test)
	}

	time.Sleep(3 * time.Second)

	tests2 := []CliTest{
		CliTest{false, false, []string{"plugin_providers", "list"}, noStdinString, "[]\n", noErrorString},
		CliTest{true, true, []string{"plugin_providers", "upload"}, noStdinString, noContentString, plugin_providerUploadNoArgErrorString},
		CliTest{true, true, []string{"plugin_providers", "upload", "john"}, noStdinString, noContentString, plugin_providerUploadTooFewArgErrorString},
		CliTest{true, true, []string{"plugin_providers", "upload", "john", "as", "john2", "asdga"}, noStdinString, noContentString, plugin_providerUploadTooManyArgErrorString},
		CliTest{false, true, []string{"plugin_providers", "upload", "john", "as", "john"}, noStdinString, noContentString, plugin_providerUploadMissingArgErrorString},
		CliTest{false, false, []string{"plugin_providers", "upload", destFolder, "as", "incrementer"}, noStdinString, plugin_providerUploadSuccessString, noErrorString},
	}

	for _, test := range tests2 {
		testCli(t, test)
	}

	time.Sleep(3 * time.Second)

	tests3 := []CliTest{
		CliTest{false, false, []string{"plugin_providers", "list"}, noStdinString, plugin_providerListString, noErrorString},
	}
	for _, test := range tests3 {
		testCli(t, test)
	}
}
