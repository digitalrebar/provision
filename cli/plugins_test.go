package cli

import (
	"testing"
)

var pluginEmptyListString string = "[]\n"
var pluginDefaultListString string = "[]\n"

var pluginShowNoArgErrorString string = "Error: drpcli plugins show [id] [flags] requires 1 argument\n"
var pluginShowTooManyArgErrorString string = "Error: drpcli plugins show [id] [flags] requires 1 argument\n"
var pluginShowMissingArgErrorString string = "Error: plugins GET: john: Not Found\n\n"
var pluginShowPluginString string = `{
  "Available": true,
  "Errors": [],
  "Name": "i-woman",
  "PluginErrors": [],
  "Provider": "incrementer",
  "ReadOnly": false,
  "Validated": true
}
`

var pluginExistsNoArgErrorString string = "Error: drpcli plugins exists [id] [flags] requires 1 argument"
var pluginExistsTooManyArgErrorString string = "Error: drpcli plugins exists [id] [flags] requires 1 argument"
var pluginExistsPluginString string = ""
var pluginExistsMissingJohnString string = "Error: plugins GET: john: Not Found\n\n"

var pluginCreateNoArgErrorString string = "Error: drpcli plugins create [json] [flags] requires 1 argument\n"
var pluginCreateTooManyArgErrorString string = "Error: drpcli plugins create [json] [flags] requires 1 argument\n"
var pluginCreateBadJSONString = "{asdgasdg"
var pluginCreateBadJSONErrorString = "Error: Invalid plugin object: error converting YAML to JSON: yaml: line 1: did not find expected ',' or '}' and error converting YAML to JSON: yaml: line 1: did not find expected ',' or '}'\n\n"
var pluginCreateBadJSON2String = "[asdgasdg]"
var pluginCreateBadJSON2ErrorString = "Error: Unable to create new plugin: Invalid type passed to plugin create\n\n"
var pluginCreateMissingProviderInputString string = `{
  "Name": "i-woman"
}
`
var pluginCreateMissingProviderErrorString string = "Error: Plugin i-woman must have a provider\n\n"
var pluginCreateMissingAllInputString string = `{
  "Description": "i-woman's plugin"
}
`
var pluginCreateMissingAllErrorString string = "Error: dataTracker create plugins: Empty key not allowed\n\n"
var pluginCreateInputString string = `{
  "Name": "i-woman",
  "Provider": "incrementer"
}
`
var pluginCreateJohnString string = `{
  "Available": true,
  "Errors": [],
  "Name": "i-woman",
  "PluginErrors": [],
  "Provider": "incrementer",
  "ReadOnly": false,
  "Validated": true
}
`

var pluginCreateDuplicateErrorString = "Error: dataTracker create plugins: i-woman already exists\n\n"

var pluginListPluginsString = `[
  {
    "Available": true,
    "Errors": [],
    "Name": "i-woman",
    "PluginErrors": [],
    "Provider": "incrementer",
    "ReadOnly": false,
    "Validated": true
  }
]
`

var pluginUpdateNoArgErrorString string = "Error: drpcli plugins update [id] [json] [flags] requires 2 arguments"
var pluginUpdateTooManyArgErrorString string = "Error: drpcli plugins update [id] [json] [flags] requires 2 arguments"
var pluginUpdateBadJSONString = "asdgasdg"
var pluginUpdateBadJSONErrorString = "Error: Unable to merge objects: json: cannot unmarshal string into Go value of type map[string]interface {}\n\n\n"
var pluginUpdateInputString string = `{
  "Description": "lpxelinux.0"
}
`
var pluginUpdateJohnString string = `{
  "Available": true,
  "Description": "lpxelinux.0",
  "Errors": [],
  "Name": "i-woman",
  "PluginErrors": [],
  "Provider": "incrementer",
  "ReadOnly": false,
  "Validated": true
}
`
var pluginUpdateJohnMissingErrorString string = "Error: plugins GET: john2: Not Found\n\n"

var pluginPatchNoArgErrorString string = "Error: drpcli plugins patch [objectJson] [changesJson] [flags] requires 2 arguments"
var pluginPatchTooManyArgErrorString string = "Error: drpcli plugins patch [objectJson] [changesJson] [flags] requires 2 arguments"
var pluginPatchBadPatchJSONString = "asdgasdg"
var pluginPatchBadPatchJSONErrorString = "Error: Unable to parse drpcli plugins patch [objectJson] [changesJson] [flags] JSON asdgasdg\nError: error unmarshaling JSON: json: cannot unmarshal string into Go value of type genmodels.Plugin\n\n"
var pluginPatchBadBaseJSONString = "asdgasdg"
var pluginPatchBadBaseJSONErrorString = "Error: Unable to parse drpcli plugins patch [objectJson] [changesJson] [flags] JSON asdgasdg\nError: error unmarshaling JSON: json: cannot unmarshal string into Go value of type genmodels.Plugin\n\n"
var pluginPatchBaseString string = `{
  "Available": true,
  "Description": "lpxelinux.0",
  "Errors": [],
  "Name": "i-woman",
  "PluginErrors": [],
  "Provider": "incrementer",
  "Validated": true
}
`
var pluginPatchInputString string = `{
  "Description": "bootx64.efi"
}
`
var pluginPatchJohnString string = `{
  "Available": true,
  "Description": "bootx64.efi",
  "Errors": [],
  "Name": "i-woman",
  "PluginErrors": [],
  "Provider": "incrementer",
  "ReadOnly": false,
  "Validated": true
}
`
var pluginPatchMissingBaseString string = `{
  "Description": "bootx64.efi",
  "Errors": [],
  "Name": "spider-woman",
  "PluginErrors": [],
  "Provider": "incrementer"
}
`
var pluginPatchJohnMissingErrorString string = "Error: plugins: PATCH spider-woman: Not Found\n\n"

var pluginDestroyNoArgErrorString string = "Error: drpcli plugins destroy [id] [flags] requires 1 argument"
var pluginDestroyTooManyArgErrorString string = "Error: drpcli plugins destroy [id] [flags] requires 1 argument"
var pluginDestroyJohnString string = "Deleted plugin i-woman\n"
var pluginDestroyMissingJohnString string = "Error: plugins: DELETE i-woman: Not Found\n\n"

var pluginBootEnvNoArgErrorString string = "Error: drpcli plugins bootenv [id] [bootenv] [flags] requires 2 arguments"
var pluginBootEnvMissingPluginErrorString string = "Error: plugins GET: john: Not Found\n\n"
var pluginBootEnvErrorBootEnvString string = "Error: Bootenv john2 does not exist\n\n"

var pluginGetNoArgErrorString string = "Error: drpcli plugins get [id] param [key] [flags] requires 3 arguments"
var pluginGetMissingPluginErrorString string = "Error: plugins GET Params: john: Not Found\n\n"

var pluginSetNoArgErrorString string = "Error: drpcli plugins set [id] param [key] to [json blob] [flags] requires 5 arguments"
var pluginSetMissingPluginErrorString string = "Error: plugins GET Params: john: Not Found\n\n"

var pluginParamsNoArgErrorString string = "Error: drpcli plugins params [id] [json] [flags] requires 1 or 2 arguments\n"
var pluginParamsMissingPluginErrorString string = "Error: plugins GET Params: john2: Not Found\n\n"
var pluginsParamsSetMissingPluginString string = "Error: plugins SET Params: john2: Not Found\n\n"

var pluginParamsStartingString string = `{
  "john3": 4
}
`
var pluginsParamsNextString string = `{
  "jj": 3
}
`
var pluginUpdateJohnWithParamsString string = `{
  "Available": true,
  "Description": "lpxelinux.0",
  "Errors": [],
  "Name": "i-woman",
  "Params": {
    "jj": 3
  },
  "PluginErrors": [],
  "Provider": "incrementer",
  "ReadOnly": false,
  "Validated": true
}
`

func TestPluginCli(t *testing.T) {

	tests := []CliTest{
		CliTest{true, false, []string{"plugins"}, noStdinString, "Access CLI commands relating to plugins\n", ""},
		CliTest{false, false, []string{"plugins", "list"}, noStdinString, pluginDefaultListString, noErrorString},

		CliTest{true, true, []string{"plugins", "create"}, noStdinString, noContentString, pluginCreateNoArgErrorString},
		CliTest{true, true, []string{"plugins", "create", "john", "john2"}, noStdinString, noContentString, pluginCreateTooManyArgErrorString},
		CliTest{false, true, []string{"plugins", "create", pluginCreateBadJSONString}, noStdinString, noContentString, pluginCreateBadJSONErrorString},
		CliTest{false, true, []string{"plugins", "create", pluginCreateBadJSON2String}, noStdinString, noContentString, pluginCreateBadJSON2ErrorString},
		CliTest{false, true, []string{"plugins", "create", pluginCreateMissingAllInputString}, noStdinString, noContentString, pluginCreateMissingAllErrorString},
		CliTest{false, true, []string{"plugins", "create", pluginCreateMissingProviderInputString}, noStdinString, noContentString, pluginCreateMissingProviderErrorString},
		CliTest{false, false, []string{"plugins", "create", pluginCreateInputString}, noStdinString, pluginCreateJohnString, noErrorString},
		CliTest{false, true, []string{"plugins", "create", pluginCreateInputString}, noStdinString, noContentString, pluginCreateDuplicateErrorString},
		CliTest{false, false, []string{"plugins", "list"}, noStdinString, pluginListPluginsString, noErrorString},
		CliTest{false, false, []string{"plugins", "list", "--limit=0"}, noStdinString, pluginEmptyListString, noErrorString},
		CliTest{false, false, []string{"plugins", "list", "--limit=10", "--offset=0"}, noStdinString, pluginListPluginsString, noErrorString},
		CliTest{false, false, []string{"plugins", "list", "--limit=10", "--offset=10"}, noStdinString, pluginEmptyListString, noErrorString},
		CliTest{false, true, []string{"plugins", "list", "--limit=-10", "--offset=0"}, noStdinString, noContentString, limitNegativeError},
		CliTest{false, true, []string{"plugins", "list", "--limit=10", "--offset=-10"}, noStdinString, noContentString, offsetNegativeError},
		CliTest{false, false, []string{"plugins", "list", "--limit=-1", "--offset=-1"}, noStdinString, pluginListPluginsString, noErrorString},
		CliTest{false, false, []string{"plugins", "list", "Name=fred"}, noStdinString, pluginEmptyListString, noErrorString},
		CliTest{false, false, []string{"plugins", "list", "Name=i-woman"}, noStdinString, pluginListPluginsString, noErrorString},
		CliTest{false, false, []string{"plugins", "list", "Provider=local"}, noStdinString, pluginEmptyListString, noErrorString},
		CliTest{false, false, []string{"plugins", "list", "Provider=incrementer"}, noStdinString, pluginListPluginsString, noErrorString},
		CliTest{false, false, []string{"plugins", "list", "Available=true"}, noStdinString, pluginListPluginsString, noErrorString},
		CliTest{false, false, []string{"plugins", "list", "Available=false"}, noStdinString, pluginEmptyListString, noErrorString},
		CliTest{false, true, []string{"plugins", "list", "Available=fred"}, noStdinString, noContentString, bootEnvBadAvailableString},
		CliTest{false, false, []string{"plugins", "list", "Valid=true"}, noStdinString, pluginListPluginsString, noErrorString},
		CliTest{false, false, []string{"plugins", "list", "Valid=false"}, noStdinString, pluginEmptyListString, noErrorString},
		CliTest{false, true, []string{"plugins", "list", "Valid=fred"}, noStdinString, noContentString, bootEnvBadValidString},
		CliTest{false, false, []string{"plugins", "list", "ReadOnly=true"}, noStdinString, pluginEmptyListString, noErrorString},
		CliTest{false, false, []string{"plugins", "list", "ReadOnly=false"}, noStdinString, pluginListPluginsString, noErrorString},
		CliTest{false, true, []string{"plugins", "list", "ReadOnly=fred"}, noStdinString, noContentString, bootEnvBadReadOnlyString},

		CliTest{true, true, []string{"plugins", "show"}, noStdinString, noContentString, pluginShowNoArgErrorString},
		CliTest{true, true, []string{"plugins", "show", "john", "john2"}, noStdinString, noContentString, pluginShowTooManyArgErrorString},
		CliTest{false, true, []string{"plugins", "show", "john"}, noStdinString, noContentString, pluginShowMissingArgErrorString},
		CliTest{false, false, []string{"plugins", "show", "i-woman"}, noStdinString, pluginShowPluginString, noErrorString},
		CliTest{false, false, []string{"plugins", "show", "Key:i-woman"}, noStdinString, pluginShowPluginString, noErrorString},
		CliTest{false, false, []string{"plugins", "show", "Name:i-woman"}, noStdinString, pluginShowPluginString, noErrorString},

		CliTest{true, true, []string{"plugins", "exists"}, noStdinString, noContentString, pluginExistsNoArgErrorString},
		CliTest{true, true, []string{"plugins", "exists", "john", "john2"}, noStdinString, noContentString, pluginExistsTooManyArgErrorString},
		CliTest{false, false, []string{"plugins", "exists", "i-woman"}, noStdinString, pluginExistsPluginString, noErrorString},
		CliTest{false, true, []string{"plugins", "exists", "john"}, noStdinString, noContentString, pluginExistsMissingJohnString},

		CliTest{true, true, []string{"plugins", "update"}, noStdinString, noContentString, pluginUpdateNoArgErrorString},
		CliTest{true, true, []string{"plugins", "update", "john", "john2", "john3"}, noStdinString, noContentString, pluginUpdateTooManyArgErrorString},
		CliTest{false, true, []string{"plugins", "update", "i-woman", pluginUpdateBadJSONString}, noStdinString, noContentString, pluginUpdateBadJSONErrorString},
		CliTest{false, false, []string{"plugins", "update", "i-woman", pluginUpdateInputString}, noStdinString, pluginUpdateJohnString, noErrorString},
		CliTest{false, true, []string{"plugins", "update", "john2", pluginUpdateInputString}, noStdinString, noContentString, pluginUpdateJohnMissingErrorString},
		CliTest{false, false, []string{"plugins", "show", "i-woman"}, noStdinString, pluginUpdateJohnString, noErrorString},

		CliTest{true, true, []string{"plugins", "patch"}, noStdinString, noContentString, pluginPatchNoArgErrorString},
		CliTest{true, true, []string{"plugins", "patch", "john", "john2", "john3"}, noStdinString, noContentString, pluginPatchTooManyArgErrorString},
		CliTest{false, true, []string{"plugins", "patch", pluginPatchBaseString, pluginPatchBadPatchJSONString}, noStdinString, noContentString, pluginPatchBadPatchJSONErrorString},
		CliTest{false, true, []string{"plugins", "patch", pluginPatchBadBaseJSONString, pluginPatchInputString}, noStdinString, noContentString, pluginPatchBadBaseJSONErrorString},
		CliTest{false, false, []string{"plugins", "patch", pluginPatchBaseString, pluginPatchInputString}, noStdinString, pluginPatchJohnString, noErrorString},
		CliTest{false, true, []string{"plugins", "patch", pluginPatchMissingBaseString, pluginPatchInputString}, noStdinString, noContentString, pluginPatchJohnMissingErrorString},
		CliTest{false, false, []string{"plugins", "show", "i-woman"}, noStdinString, pluginPatchJohnString, noErrorString},

		CliTest{true, true, []string{"plugins", "destroy"}, noStdinString, noContentString, pluginDestroyNoArgErrorString},
		CliTest{true, true, []string{"plugins", "destroy", "john", "june"}, noStdinString, noContentString, pluginDestroyTooManyArgErrorString},
		CliTest{false, false, []string{"plugins", "destroy", "i-woman"}, noStdinString, pluginDestroyJohnString, noErrorString},
		CliTest{false, true, []string{"plugins", "destroy", "i-woman"}, noStdinString, noContentString, pluginDestroyMissingJohnString},
		CliTest{false, false, []string{"plugins", "list"}, noStdinString, pluginDefaultListString, noErrorString},

		CliTest{false, false, []string{"plugins", "create", "-"}, pluginCreateInputString + "\n", pluginCreateJohnString, noErrorString},
		CliTest{false, false, []string{"plugins", "list"}, noStdinString, pluginListPluginsString, noErrorString},
		CliTest{false, false, []string{"plugins", "update", "i-woman", "-"}, pluginUpdateInputString + "\n", pluginUpdateJohnString, noErrorString},
		CliTest{false, false, []string{"plugins", "show", "i-woman"}, noStdinString, pluginUpdateJohnString, noErrorString},

		CliTest{true, true, []string{"plugins", "get"}, noStdinString, noContentString, pluginGetNoArgErrorString},
		CliTest{false, true, []string{"plugins", "get", "john", "param", "john2"}, noStdinString, noContentString, pluginGetMissingPluginErrorString},
		CliTest{false, false, []string{"plugins", "get", "i-woman", "param", "john2"}, noStdinString, "null\n", noErrorString},

		CliTest{true, true, []string{"plugins", "set"}, noStdinString, noContentString, pluginSetNoArgErrorString},
		CliTest{false, true, []string{"plugins", "set", "john", "param", "john2", "to", "cow"}, noStdinString, noContentString, pluginSetMissingPluginErrorString},
		CliTest{false, false, []string{"plugins", "set", "i-woman", "param", "john2", "to", "cow"}, noStdinString, "\"cow\"\n", noErrorString},
		CliTest{false, false, []string{"plugins", "get", "i-woman", "param", "john2"}, noStdinString, "\"cow\"\n", noErrorString},
		CliTest{false, false, []string{"plugins", "set", "i-woman", "param", "john2", "to", "3"}, noStdinString, "3\n", noErrorString},
		CliTest{false, false, []string{"plugins", "set", "i-woman", "param", "john3", "to", "4"}, noStdinString, "4\n", noErrorString},
		CliTest{false, false, []string{"plugins", "get", "i-woman", "param", "john2"}, noStdinString, "3\n", noErrorString},
		CliTest{false, false, []string{"plugins", "get", "i-woman", "param", "john3"}, noStdinString, "4\n", noErrorString},
		CliTest{false, false, []string{"plugins", "set", "i-woman", "param", "john2", "to", "null"}, noStdinString, "null\n", noErrorString},
		CliTest{false, false, []string{"plugins", "get", "i-woman", "param", "john2"}, noStdinString, "null\n", noErrorString},
		CliTest{false, false, []string{"plugins", "get", "i-woman", "param", "john3"}, noStdinString, "4\n", noErrorString},

		CliTest{true, true, []string{"plugins", "params"}, noStdinString, noContentString, pluginParamsNoArgErrorString},
		CliTest{false, true, []string{"plugins", "params", "john2"}, noStdinString, noContentString, pluginParamsMissingPluginErrorString},
		CliTest{false, false, []string{"plugins", "params", "i-woman"}, noStdinString, pluginParamsStartingString, noErrorString},
		CliTest{false, true, []string{"plugins", "params", "john2", pluginsParamsNextString}, noStdinString, noContentString, pluginsParamsSetMissingPluginString},
		CliTest{false, false, []string{"plugins", "params", "i-woman", pluginsParamsNextString}, noStdinString, pluginsParamsNextString, noErrorString},
		CliTest{false, false, []string{"plugins", "params", "i-woman"}, noStdinString, pluginsParamsNextString, noErrorString},

		CliTest{false, false, []string{"plugins", "show", "i-woman"}, noStdinString, pluginUpdateJohnWithParamsString, noErrorString},

		CliTest{false, false, []string{"plugins", "destroy", "i-woman"}, noStdinString, pluginDestroyJohnString, noErrorString},
		CliTest{false, false, []string{"plugins", "list"}, noStdinString, pluginDefaultListString, noErrorString},
	}

	for _, test := range tests {
		testCli(t, test)
	}
}
