package cli

import (
	"testing"
)

func TestContentCli(t *testing.T) {
	var (
		contentCreateBadJSONString         = "{asdgasdg"
		contentCreateBadJSON2String        = "[asdgasdg]"
		contentCreateInputString    string = `{
  "meta": {
    "Name": "john"
  }
}
`

		contentUpdateBadJSONString         = "asdgasdg"
		contentUpdateBadInputString string = `{
  "meta": {
    "Name": "john2"
  }
}
`
		contentUpdateInputString string = `{
  "meta": {
    "Description": "Fred Rules",
    "Name": "john"
  }
}
`
	)

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
