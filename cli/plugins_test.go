package cli

import (
	"testing"
)

var pluginEmptyListString string = "[]\n"
var pluginDefaultListString string = "[]\n"

var pluginShowNoArgErrorString string = "Error: drpcli plugins show [id] [flags] requires 1 argument\n"
var pluginShowTooManyArgErrorString string = "Error: drpcli plugins show [id] [flags] requires 1 argument\n"
var pluginShowMissingArgErrorString string = "Error: GET: plugins/john: Not Found\n\n"
var pluginExistsNoArgErrorString string = "Error: drpcli plugins exists [id] [flags] requires 1 argument"
var pluginExistsTooManyArgErrorString string = "Error: drpcli plugins exists [id] [flags] requires 1 argument"
var pluginExistsMissingJohnString string = "Error: GET: plugins/john: Not Found\n\n"
var pluginCreateNoArgErrorString string = "Error: drpcli plugins create [json] [flags] requires 1 argument\n"
var pluginCreateTooManyArgErrorString string = "Error: drpcli plugins create [json] [flags] requires 1 argument\n"
var pluginCreateBadJSONErrorString = "Error: Invalid plugin object: error converting YAML to JSON: yaml: line 1: did not find expected ',' or '}' and error converting YAML to JSON: yaml: line 1: did not find expected ',' or '}'\n\n"
var pluginCreateBadJSON2ErrorString = "Error: Unable to create new plugin: Invalid type passed to plugin create\n\n"
var pluginCreateMissingProviderErrorString string = "Error: ValidationError: plugins/i-woman: Missing provider\n\n"
var pluginCreateMissingAllErrorString string = "Error: CREATE: plugins: Empty key not allowed\n\n"
var pluginCreateDuplicateErrorString = "Error: CREATE: plugins/i-woman: already exists\n\n"
var pluginUpdateNoArgErrorString string = "Error: drpcli plugins update [id] [json] [flags] requires 2 arguments"
var pluginUpdateTooManyArgErrorString string = "Error: drpcli plugins update [id] [json] [flags] requires 2 arguments"
var pluginUpdateBadJSONErrorString = "Error: Unable to merge objects: json: cannot unmarshal string into Go value of type map[string]interface {}\n\n\n"
var pluginUpdateJohnMissingErrorString string = "Error: GET: plugins/john2: Not Found\n\n"
var pluginPatchNoArgErrorString string = "Error: drpcli plugins patch [objectJson] [changesJson] [flags] requires 2 arguments"
var pluginPatchTooManyArgErrorString string = "Error: drpcli plugins patch [objectJson] [changesJson] [flags] requires 2 arguments"
var pluginPatchBadPatchJSONErrorString = "Error: Unable to parse drpcli plugins patch [objectJson] [changesJson] [flags] JSON asdgasdg\nError: error unmarshaling JSON: json: cannot unmarshal string into Go value of type genmodels.Plugin\n\n"
var pluginPatchBadBaseJSONErrorString = "Error: Unable to parse drpcli plugins patch [objectJson] [changesJson] [flags] JSON asdgasdg\nError: error unmarshaling JSON: json: cannot unmarshal string into Go value of type genmodels.Plugin\n\n"
var pluginPatchJohnMissingErrorString string = "Error: PATCH: plugins/spider-woman: Not Found\n\n"
var pluginDestroyNoArgErrorString string = "Error: drpcli plugins destroy [id] [flags] requires 1 argument"
var pluginDestroyTooManyArgErrorString string = "Error: drpcli plugins destroy [id] [flags] requires 1 argument"
var pluginDestroyMissingJohnString string = "Error: DELETE: plugins/i-woman: Not Found\n\n"
var pluginBootEnvNoArgErrorString string = "Error: drpcli plugins bootenv [id] [bootenv] [flags] requires 2 arguments"
var pluginBootEnvMissingPluginErrorString string = "Error: plugins GET: john: Not Found\n\n"
var pluginBootEnvErrorBootEnvString string = "Error: Bootenv john2 does not exist\n\n"
var pluginGetNoArgErrorString string = "Error: drpcli plugins get [id] param [key] [flags] requires 3 arguments"
var pluginGetMissingPluginErrorString string = "Error: GET: plugins/john: Not Found\n\n"
var pluginSetNoArgErrorString string = "Error: drpcli plugins set [id] param [key] to [json blob] [flags] requires 5 arguments"
var pluginSetMissingPluginErrorString string = "Error: GET: plugins/john: Not Found\n\n"
var pluginParamsNoArgErrorString string = "Error: drpcli plugins params [id] [json] [flags] requires 1 or 2 arguments\n"
var pluginParamsMissingPluginErrorString string = "Error: GET: plugins/john2: Not Found\n\n"
var pluginsParamsSetMissingPluginString string = "Error: POST: plugins/john2: Not Found\n\n"

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

var pluginExistsPluginString string = ""
var pluginCreateBadJSONString = "{asdgasdg"

var pluginCreateBadJSON2String = "[asdgasdg]"

var pluginCreateMissingProviderInputString string = `{
  "Name": "i-woman"
}
`

var pluginCreateMissingAllInputString string = `{
  "Description": "i-woman's plugin"
}
`

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

var pluginUpdateBadJSONString = "asdgasdg"

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

var pluginPatchBadPatchJSONString = "asdgasdg"

var pluginPatchBadBaseJSONString = "asdgasdg"

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

var pluginDestroyJohnString string = "Deleted plugin i-woman\n"

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
	cliTest(true, false, "plugins").run(t)
	cliTest(false, false, "plugins", "list").run(t)
	cliTest(true, true, "plugins", "create").run(t)
	cliTest(true, true, "plugins", "create", "john", "john2").run(t)
	cliTest(false, true, "plugins", "create", pluginCreateBadJSONString).run(t)
	cliTest(false, true, "plugins", "create", pluginCreateBadJSON2String).run(t)
	cliTest(false, true, "plugins", "create", pluginCreateMissingAllInputString).run(t)
	cliTest(false, true, "plugins", "create", pluginCreateMissingProviderInputString).run(t)
	cliTest(false, false, "plugins", "create", pluginCreateInputString).run(t)
	cliTest(false, true, "plugins", "create", pluginCreateInputString).run(t)
	cliTest(false, false, "plugins", "list").run(t)
	cliTest(false, false, "plugins", "list", "Name=fred").run(t)
	cliTest(false, false, "plugins", "list", "Name=i-woman").run(t)
	cliTest(false, false, "plugins", "list", "Provider=local").run(t)
	cliTest(false, false, "plugins", "list", "Provider=incrementer").run(t)
	cliTest(true, true, "plugins", "show").run(t)
	cliTest(true, true, "plugins", "show", "john", "john2").run(t)
	cliTest(false, true, "plugins", "show", "john").run(t)
	cliTest(false, false, "plugins", "show", "i-woman").run(t)
	cliTest(false, false, "plugins", "show", "Key:i-woman").run(t)
	cliTest(false, false, "plugins", "show", "Name:i-woman").run(t)
	cliTest(true, true, "plugins", "exists").run(t)
	cliTest(true, true, "plugins", "exists", "john", "john2").run(t)
	cliTest(false, false, "plugins", "exists", "i-woman").run(t)
	cliTest(false, true, "plugins", "exists", "john").run(t)
	cliTest(true, true, "plugins", "update").run(t)
	cliTest(true, true, "plugins", "update", "john", "john2", "john3").run(t)
	cliTest(false, true, "plugins", "update", "i-woman", pluginUpdateBadJSONString).run(t)
	cliTest(false, false, "plugins", "update", "i-woman", pluginUpdateInputString).run(t)
	cliTest(false, true, "plugins", "update", "john2", pluginUpdateInputString).run(t)
	cliTest(false, false, "plugins", "show", "i-woman").run(t)
	cliTest(false, false, "plugins", "show", "i-woman").run(t)
	cliTest(true, true, "plugins", "destroy").run(t)
	cliTest(true, true, "plugins", "destroy", "john", "june").run(t)
	cliTest(false, false, "plugins", "destroy", "i-woman").run(t)
	cliTest(false, true, "plugins", "destroy", "i-woman").run(t)
	cliTest(false, false, "plugins", "list").run(t)
	cliTest(false, false, "plugins", "create", "-").Stdin(pluginCreateInputString + "\n").run(t)
	cliTest(false, false, "plugins", "list").run(t)
	cliTest(false, false, "plugins", "update", "i-woman", "-").Stdin(pluginUpdateInputString + "\n").run(t)
	cliTest(false, false, "plugins", "show", "i-woman").run(t)
	cliTest(true, true, "plugins", "get").run(t)
	cliTest(false, true, "plugins", "get", "john", "param", "john2").run(t)
	cliTest(false, false, "plugins", "get", "i-woman", "param", "john2").run(t)
	cliTest(true, true, "plugins", "set").run(t)
	cliTest(false, true, "plugins", "set", "john", "param", "john2", "to", "cow").run(t)
	cliTest(false, false, "plugins", "set", "i-woman", "param", "john2", "to", "cow").run(t)
	cliTest(false, false, "plugins", "get", "i-woman", "param", "john2").run(t)
	cliTest(false, false, "plugins", "set", "i-woman", "param", "john2", "to", "3").run(t)
	cliTest(false, false, "plugins", "set", "i-woman", "param", "john3", "to", "4").run(t)
	cliTest(false, false, "plugins", "get", "i-woman", "param", "john2").run(t)
	cliTest(false, false, "plugins", "get", "i-woman", "param", "john3").run(t)
	cliTest(false, false, "plugins", "set", "i-woman", "param", "john2", "to", "null").run(t)
	cliTest(false, false, "plugins", "get", "i-woman", "param", "john2").run(t)
	cliTest(false, false, "plugins", "get", "i-woman", "param", "john3").run(t)
	cliTest(true, true, "plugins", "params").run(t)
	cliTest(false, true, "plugins", "params", "john2").run(t)
	cliTest(false, false, "plugins", "params", "i-woman").run(t)
	cliTest(false, true, "plugins", "params", "john2", pluginsParamsNextString).run(t)
	cliTest(false, false, "plugins", "params", "i-woman", pluginsParamsNextString).run(t)
	cliTest(false, false, "plugins", "params", "i-woman").run(t)
	cliTest(false, false, "plugins", "show", "i-woman").run(t)
	cliTest(false, false, "plugins", "destroy", "i-woman").run(t)
	cliTest(false, false, "plugins", "list").run(t)
}
