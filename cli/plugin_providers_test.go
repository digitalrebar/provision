package cli

import (
	"testing"
)

var plugin_providerShowNoArgErrorString string = "Error: drpcli plugin_providers show [id] requires 1 argument\n"
var plugin_providerShowTooManyArgErrorString string = "Error: drpcli plugin_providers show [id] requires 1 argument\n"
var plugin_providerShowMissingArgErrorString string = "Error: plugin provider get: not found: john\n\n"

var plugin_providerExistsNoArgErrorString string = "Error: drpcli plugin_providers exists [id] requires 1 argument"
var plugin_providerExistsTooManyArgErrorString string = "Error: drpcli plugin_providers exists [id] requires 1 argument"
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
        "Name": "incrementer.parameter",
        "Schema": {
          "type": "string"
        }
      },
      {
        "Name": "incrementer.step",
        "Schema": {
          "type": "integer"
        }
      },
      {
        "Name": "incrementer.touched",
        "Schema": {
          "type": "integer"
        }
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
      "Name": "incrementer.parameter",
      "Schema": {
        "type": "string"
      }
    },
    {
      "Name": "incrementer.step",
      "Schema": {
        "type": "integer"
      }
    },
    {
      "Name": "incrementer.touched",
      "Schema": {
        "type": "integer"
      }
    }
  ],
  "RequiredParams": null,
  "Version": "v3.0.2-pre-alpha-NotSet"
}
`

func TestPluginProviderCli(t *testing.T) {
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
	}

	for _, test := range tests {
		testCli(t, test)
	}
}
