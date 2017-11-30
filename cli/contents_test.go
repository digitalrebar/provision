package cli

import (
	"testing"
)

var contentDefaultListString string = `[
  {
    "Counts": {
      "bootenvs": 0,
      "jobs": 0,
      "leases": 0,
      "machines": 0,
      "params": 0,
      "plugins": 0,
      "preferences": 0,
      "profiles": 1,
      "reservations": 0,
      "stages": 0,
      "subnets": 0,
      "tasks": 0,
      "templates": 0,
      "users": 1
    },
    "Warnings": [],
    "meta": {
      "Description": "Writable backing store",
      "Name": "BackingStore",
      "Type": "writable",
      "Version": "user",
      "Writable": true
    }
  },
  {
    "Counts": {
      "templates": 1
    },
    "Warnings": [],
    "meta": {
      "Description": "Local Override Store",
      "Name": "LocalStore",
      "Type": "local",
      "Version": "user"
    }
  },
  {
    "Counts": {
      "templates": 1
    },
    "Warnings": [],
    "meta": {
      "Description": "Initial Default Content",
      "Name": "DefaultStore",
      "Overwritable": true,
      "Source": "Unspecified",
      "Type": "default",
      "Version": "user"
    }
  },
  {
    "Counts": {
      "params": 3
    },
    "Warnings": [],
    "meta": {
      "Description": "Test Plugin for DRP",
      "Name": "incrementer",
      "Source": "Digital Rebar",
      "Type": "plugin",
      "Version": "Internal"
    }
  },
  {
    "Counts": {
      "bootenvs": 2,
      "stages": 2
    },
    "Warnings": [],
    "meta": {
      "Description": "Default objects that must be present",
      "Name": "BasicStore",
      "Type": "basic",
      "Version": "Unversioned"
    }
  }
]
`

var contentEmptyListString string = "[]\n"
var contentListBadFlagErrorString string = "Error: unknown flag: --limit\n"
var contentListBadFilterErrorString string = "Error: Filter argument requires an '=' separator: Cow\n\n"
var contentListFilterNotSupportedErrorString string = "Error: listing contents: Does not support filtering\n\n"

var contentShowNoArgErrorString string = "Error: drpcli contents show [id] [flags] requires 1 argument\n"
var contentShowTooManyArgErrorString string = "Error: drpcli contents show [id] [flags] requires 1 argument\n"
var contentShowMissingArgErrorString string = "Error: GET: contents/john2: No such content store\n\n"
var contentShowContentString string = `{
  "meta": {
    "Name": "john",
    "Type": "dynamic"
  }
}
`

var contentExistsNoArgErrorString string = "Error: drpcli contents exists [id] [flags] requires 1 argument"
var contentExistsTooManyArgErrorString string = "Error: drpcli contents exists [id] [flags] requires 1 argument"
var contentExistsContentString string = ""
var contentExistsMissingJohnString string = "Error: GET: contents/john2: No such content store\n\n"

var contentCreateNoArgErrorString string = "Error: drpcli contents create [json] [flags] requires 1 argument\n"
var contentCreateTooManyArgErrorString string = "Error: drpcli contents create [json] [flags] requires 1 argument\n"
var contentCreateBadJSONString = "{asdgasdg"
var contentCreateBadJSONErrorString = "Error: Invalid content object: error converting YAML to JSON: yaml: line 1: did not find expected ',' or '}' and error converting YAML to JSON: yaml: line 1: did not find expected ',' or '}'\n\n"
var contentCreateBadJSON2String = "[asdgasdg]"
var contentCreateBadJSON2ErrorString = "Error: Unable to create new content: Invalid type passed to content create\n\n"
var contentCreateInputString string = `{
  "meta": {
    "Name": "john"
  }
}
`
var contentCreateJohnString string = `{
  "Warnings": [],
  "meta": {
    "Name": "john",
    "Type": "dynamic"
  }
}
`
var contentCreateDuplicateErrorString = "Error: POST: contents/john: Content john already exists\n\n"

var contentListContentsString = `[
  {
    "Counts": {
      "bootenvs": 0,
      "jobs": 0,
      "leases": 0,
      "machines": 0,
      "params": 0,
      "plugins": 0,
      "preferences": 0,
      "profiles": 1,
      "reservations": 0,
      "stages": 0,
      "subnets": 0,
      "tasks": 0,
      "templates": 0,
      "users": 1
    },
    "Warnings": [],
    "meta": {
      "Description": "Writable backing store",
      "Name": "BackingStore",
      "Type": "writable",
      "Version": "user",
      "Writable": true
    }
  },
  {
    "Counts": {
      "templates": 1
    },
    "Warnings": [],
    "meta": {
      "Description": "Local Override Store",
      "Name": "LocalStore",
      "Type": "local",
      "Version": "user"
    }
  },
  {
    "Warnings": [],
    "meta": {
      "Name": "john",
      "Type": "dynamic"
    }
  },
  {
    "Counts": {
      "templates": 1
    },
    "Warnings": [],
    "meta": {
      "Description": "Initial Default Content",
      "Name": "DefaultStore",
      "Overwritable": true,
      "Source": "Unspecified",
      "Type": "default",
      "Version": "user"
    }
  },
  {
    "Counts": {
      "params": 3
    },
    "Warnings": [],
    "meta": {
      "Description": "Test Plugin for DRP",
      "Name": "incrementer",
      "Source": "Digital Rebar",
      "Type": "plugin",
      "Version": "Internal"
    }
  },
  {
    "Counts": {
      "bootenvs": 2,
      "stages": 2
    },
    "Warnings": [],
    "meta": {
      "Description": "Default objects that must be present",
      "Name": "BasicStore",
      "Type": "basic",
      "Version": "Unversioned"
    }
  }
]
`

var contentUpdateNoArgErrorString string = "Error: drpcli contents update [id] [json] [flags] requires 2 arguments"
var contentUpdateTooManyArgErrorString string = "Error: drpcli contents update [id] [json] [flags] requires 2 arguments"
var contentUpdateBadJSONString = "asdgasdg"
var contentUpdateBadJSONErrorString = "Error: Unable to unmarshal merged input stream: error unmarshaling JSON: json: cannot unmarshal string into Go value of type genmodels.Content\n\n\n"
var contentUpdateBadInputString string = `{
  "meta": {
    "Name": "john2"
  }
}
`
var contentUpdateBadInputErrorString string = "Error: PUT: contents/john: Cannot change name from john to john2\n\n"
var contentUpdateInputString string = `{
  "meta": {
    "Description": "Fred Rules",
    "Name": "john"
  }
}
`
var contentUpdateJohnString string = `{
  "Warnings": [],
  "meta": {
    "Description": "Fred Rules",
    "Name": "john",
    "Type": "dynamic"
  }
}
`
var contentShowJohnString string = `{
  "meta": {
    "Description": "Fred Rules",
    "Name": "john",
    "Type": "dynamic"
  }
}
`
var contentUpdateJohnMissingErrorString string = "Error: GET: contents/john2: No such content store\n\n"

var contentDestroyNoArgErrorString string = "Error: drpcli contents destroy [id] [flags] requires 1 argument"
var contentDestroyTooManyArgErrorString string = "Error: drpcli contents destroy [id] [flags] requires 1 argument"
var contentDestroyJohnString string = "Deleted content john\n"
var contentDestroyMissingJohnString string = "Error: DELETE: contents/john: No such content store\n\n"

func TestContentCli(t *testing.T) {

	cliTest(true, false, "contents").run(t)
	cliTest(false, false, "contents", "list").run(t)

	cliTest(false, true, "contents", "create").run(t)
	cliTest(false, true, "contents", "create", "john", "john2").run(t)
	cliTest(false, true, "contents", "create", contentCreateBadJSONString).run(t)
	cliTest(false, true, "contents", "create", contentCreateBadJSON2String).run(t)
	cliTest(false, false, "contents", "create", contentCreateInputString).run(t)
	cliTest(false, true, "contents", "create", contentCreateInputString).run(t)
	cliTest(false, false, "contents", "list").run(t)
	cliTest(false, true, "contents", "list", "--limit=-1", "--offset=-1").run(t)
	cliTest(false, true, "contents", "list", "Cow").run(t)
	cliTest(false, true, "contents", "list", "Cow=john").run(t)

	cliTest(true, true, "contents", "show").run(t)
	cliTest(true, true, "contents", "show", "john", "john2").run(t)
	cliTest(false, true, "contents", "show", "john2").run(t)
	cliTest(false, false, "contents", "show", "john").run(t)

	cliTest(false, true, "contents", "exists").run(t)
	cliTest(false, true, "contents", "exists", "john", "john2").run(t)
	cliTest(false, false, "contents", "exists", "john").run(t)
	cliTest(false, true, "contents", "exists", "john2").run(t)
	cliTest(true, true, "contents", "exists", "john", "john2").run(t)

	cliTest(false, true, "contents", "update").run(t)
	cliTest(false, true, "contents", "update", "john", "john2", "john3").run(t)
	cliTest(false, true, "contents", "update", "john", contentUpdateBadJSONString).run(t)
	cliTest(false, true, "contents", "update", "john", contentUpdateBadInputString).run(t)
	cliTest(false, false, "contents", "update", "john", contentUpdateInputString).run(t)
	cliTest(false, true, "contents", "update", "john2", contentUpdateInputString).run(t)
	cliTest(false, false, "contents", "show", "john").run(t)

	cliTest(false, true, "contents", "destroy").run(t)
	cliTest(false, true, "contents", "destroy", "john", "june").run(t)
	cliTest(false, false, "contents", "destroy", "john").run(t)
	cliTest(false, true, "contents", "destroy", "john").run(t)
	cliTest(false, false, "contents", "list").run(t)

	cliTest(false, false, "contents", "create", "-").Stdin(contentCreateInputString + "\n").run(t)
	cliTest(false, false, "contents", "list").run(t)
	cliTest(false, false, "contents", "update", "john", "-").Stdin(contentUpdateInputString + "\n").run(t)
	cliTest(false, false, "contents", "show", "john").run(t)

	cliTest(false, false, "contents", "destroy", "john").run(t)
	cliTest(false, false, "contents", "list").run(t)
}
