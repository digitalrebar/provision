package cli

import (
	"testing"
)

var contentDefaultListString string = `[
  {
    "Counts": {
      "bootenvs": 1,
      "jobs": 0,
      "leases": 0,
      "machines": 0,
      "params": 3,
      "plugins": 0,
      "preferences": 0,
      "profiles": 1,
      "reservations": 0,
      "subnets": 0,
      "tasks": 0,
      "templates": 0,
      "users": 1
    },
    "Description": "Writable backing store",
    "Name": "BackingStore",
    "Version": "user"
  },
  {
    "Counts": {
      "templates": 1
    },
    "Description": "Local Override Store",
    "Name": "LocalStore",
    "Version": "user"
  },
  {
    "Counts": {
      "templates": 1
    },
    "Description": "Initial Default Content",
    "Name": "DefaultStore",
    "Version": "user"
  }
]
`

var contentEmptyListString string = "[]\n"
var contentListBadFlagErrorString string = "Error: unknown flag: --limit\n"
var contentListBadFilterErrorString string = "Error: Filter argument requires an '=' separator: Cow\n\n"
var contentListFilterNotSupportedErrorString string = "Error: listing contents: Does not support filtering\n\n"

var contentShowNoArgErrorString string = "Error: drpcli contents show [id] [flags] requires 1 argument\n"
var contentShowTooManyArgErrorString string = "Error: drpcli contents show [id] [flags] requires 1 argument\n"
var contentShowMissingArgErrorString string = "Error: content get: not found: john2\n\n"
var contentShowContentString string = `{
  "Name": "john"
}
`

var contentExistsNoArgErrorString string = "Error: drpcli contents exists [id] [flags] requires 1 argument"
var contentExistsTooManyArgErrorString string = "Error: drpcli contents exists [id] [flags] requires 1 argument"
var contentExistsContentString string = ""
var contentExistsMissingJohnString string = "Error: content get: not found: john2\n\n"

var contentCreateNoArgErrorString string = "Error: drpcli contents create [json] [flags] requires 1 argument\n"
var contentCreateTooManyArgErrorString string = "Error: drpcli contents create [json] [flags] requires 1 argument\n"
var contentCreateBadJSONString = "{asdgasdg"
var contentCreateBadJSONErrorString = "Error: Invalid content object: error converting YAML to JSON: yaml: line 1: did not find expected ',' or '}' and error converting YAML to JSON: yaml: line 1: did not find expected ',' or '}'\n\n"
var contentCreateBadJSON2String = "[asdgasdg]"
var contentCreateBadJSON2ErrorString = "Error: Unable to create new content: Invalid type passed to content create\n\n"
var contentCreateInputString string = `{
  "Name": "john"
}
`
var contentCreateJohnString string = `{
  "Name": "john"
}
`
var contentCreateDuplicateErrorString = "Error: content post: already exists: john\n\n"

var contentListContentsString = `[
  {
    "Counts": {
      "bootenvs": 1,
      "jobs": 0,
      "leases": 0,
      "machines": 0,
      "params": 3,
      "plugins": 0,
      "preferences": 0,
      "profiles": 1,
      "reservations": 0,
      "subnets": 0,
      "tasks": 0,
      "templates": 0,
      "users": 1
    },
    "Description": "Writable backing store",
    "Name": "BackingStore",
    "Version": "user"
  },
  {
    "Counts": {
      "templates": 1
    },
    "Description": "Local Override Store",
    "Name": "LocalStore",
    "Version": "user"
  },
  {
    "Name": "john"
  },
  {
    "Counts": {
      "templates": 1
    },
    "Description": "Initial Default Content",
    "Name": "DefaultStore",
    "Version": "user"
  }
]
`

var contentUpdateNoArgErrorString string = "Error: drpcli contents update [id] [json] [flags] requires 2 arguments"
var contentUpdateTooManyArgErrorString string = "Error: drpcli contents update [id] [json] [flags] requires 2 arguments"
var contentUpdateBadJSONString = "asdgasdg"
var contentUpdateBadJSONErrorString = "Error: Unable to unmarshal merged input stream: error unmarshaling JSON: json: cannot unmarshal string into Go value of type genmodels.Content\n\n\n"
var contentUpdateBadInputString string = `{
  "Name": "john2"
}
`
var contentUpdateBadInputErrorString string = "Error: Name must match: john2 != john\n\n\n"
var contentUpdateInputString string = `{
  "Description": "Fred Rules",
  "Name": "john"
}
`
var contentUpdateJohnString string = `{
  "Description": "Fred Rules",
  "Name": "john"
}
`
var contentUpdateJohnMissingErrorString string = "Error: content get: not found: john2\n\n"

var contentDestroyNoArgErrorString string = "Error: drpcli contents destroy [id] [flags] requires 1 argument"
var contentDestroyTooManyArgErrorString string = "Error: drpcli contents destroy [id] [flags] requires 1 argument"
var contentDestroyJohnString string = "Deleted content john\n"
var contentDestroyMissingJohnString string = "Error: content get: not found: john\n\n"

func TestContentCli(t *testing.T) {

	tests := []CliTest{
		CliTest{true, false, []string{"contents"}, noStdinString, "Access CLI commands relating to contents\n", ""},
		CliTest{false, false, []string{"contents", "list"}, noStdinString, contentDefaultListString, noErrorString},

		CliTest{true, true, []string{"contents", "create"}, noStdinString, noContentString, contentCreateNoArgErrorString},
		CliTest{true, true, []string{"contents", "create", "john", "john2"}, noStdinString, noContentString, contentCreateTooManyArgErrorString},
		CliTest{false, true, []string{"contents", "create", contentCreateBadJSONString}, noStdinString, noContentString, contentCreateBadJSONErrorString},
		CliTest{false, true, []string{"contents", "create", contentCreateBadJSON2String}, noStdinString, noContentString, contentCreateBadJSON2ErrorString},
		CliTest{false, false, []string{"contents", "create", contentCreateInputString}, noStdinString, contentCreateJohnString, noErrorString},
		CliTest{false, true, []string{"contents", "create", contentCreateInputString}, noStdinString, noContentString, contentCreateDuplicateErrorString},
		CliTest{false, false, []string{"contents", "list"}, noStdinString, contentListContentsString, noErrorString},
		CliTest{true, true, []string{"contents", "list", "--limit=-1", "--offset=-1"}, noStdinString, noContentString, contentListBadFlagErrorString},
		CliTest{false, true, []string{"contents", "list", "Cow"}, noStdinString, noContentString, contentListBadFilterErrorString},
		CliTest{false, true, []string{"contents", "list", "Cow=john"}, noStdinString, noContentString, contentListFilterNotSupportedErrorString},

		CliTest{true, true, []string{"contents", "show"}, noStdinString, noContentString, contentShowNoArgErrorString},
		CliTest{true, true, []string{"contents", "show", "john", "john2"}, noStdinString, noContentString, contentShowTooManyArgErrorString},
		CliTest{false, true, []string{"contents", "show", "john2"}, noStdinString, noContentString, contentShowMissingArgErrorString},
		CliTest{false, false, []string{"contents", "show", "john"}, noStdinString, contentShowContentString, noErrorString},

		CliTest{true, true, []string{"contents", "exists"}, noStdinString, noContentString, contentExistsNoArgErrorString},
		CliTest{true, true, []string{"contents", "exists", "john", "john2"}, noStdinString, noContentString, contentExistsTooManyArgErrorString},
		CliTest{false, false, []string{"contents", "exists", "john"}, noStdinString, contentExistsContentString, noErrorString},
		CliTest{false, true, []string{"contents", "exists", "john2"}, noStdinString, noContentString, contentExistsMissingJohnString},
		CliTest{true, true, []string{"contents", "exists", "john", "john2"}, noStdinString, noContentString, contentExistsTooManyArgErrorString},

		CliTest{true, true, []string{"contents", "update"}, noStdinString, noContentString, contentUpdateNoArgErrorString},
		CliTest{true, true, []string{"contents", "update", "john", "john2", "john3"}, noStdinString, noContentString, contentUpdateTooManyArgErrorString},
		CliTest{false, true, []string{"contents", "update", "john", contentUpdateBadJSONString}, noStdinString, noContentString, contentUpdateBadJSONErrorString},
		CliTest{false, true, []string{"contents", "update", "john", contentUpdateBadInputString}, noStdinString, noContentString, contentUpdateBadInputErrorString},
		CliTest{false, false, []string{"contents", "update", "john", contentUpdateInputString}, noStdinString, contentUpdateJohnString, noErrorString},
		CliTest{false, true, []string{"contents", "update", "john2", contentUpdateInputString}, noStdinString, noContentString, contentUpdateJohnMissingErrorString},
		CliTest{false, false, []string{"contents", "show", "john"}, noStdinString, contentUpdateJohnString, noErrorString},

		CliTest{true, true, []string{"contents", "destroy"}, noStdinString, noContentString, contentDestroyNoArgErrorString},
		CliTest{true, true, []string{"contents", "destroy", "john", "june"}, noStdinString, noContentString, contentDestroyTooManyArgErrorString},
		CliTest{false, false, []string{"contents", "destroy", "john"}, noStdinString, contentDestroyJohnString, noErrorString},
		CliTest{false, true, []string{"contents", "destroy", "john"}, noStdinString, noContentString, contentDestroyMissingJohnString},
		CliTest{false, false, []string{"contents", "list"}, noStdinString, contentDefaultListString, noErrorString},

		CliTest{false, false, []string{"contents", "create", "-"}, contentCreateInputString + "\n", contentCreateJohnString, noErrorString},
		CliTest{false, false, []string{"contents", "list"}, noStdinString, contentListContentsString, noErrorString},
		CliTest{false, false, []string{"contents", "update", "john", "-"}, contentUpdateInputString + "\n", contentUpdateJohnString, noErrorString},
		CliTest{false, false, []string{"contents", "show", "john"}, noStdinString, contentUpdateJohnString, noErrorString},

		CliTest{false, false, []string{"contents", "destroy", "john"}, noStdinString, contentDestroyJohnString, noErrorString},
		CliTest{false, false, []string{"contents", "list"}, noStdinString, contentDefaultListString, noErrorString},
	}

	for _, test := range tests {
		testCli(t, test)
	}

}
