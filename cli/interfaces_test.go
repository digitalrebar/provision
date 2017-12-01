package cli

import "testing"

var interfaceShowNoArgErrorString string = "Error: drpcli interfaces show [id] [flags] requires 1 argument\n"
var interfaceShowTooManyArgErrorString string = "Error: drpcli interfaces show [id] [flags] requires 1 argument\n"
var interfaceShowMissingArgErrorString string = "Error: GET: interfaces/john: No interface\n\n"

var interfaceExistsNoArgErrorString string = "Error: drpcli interfaces exists [id] [flags] requires 1 argument"
var interfaceExistsTooManyArgErrorString string = "Error: drpcli interfaces exists [id] [flags] requires 1 argument"
var interfaceExistsIgnoreString string = ""
var interfaceExistsMissingJohnString string = "Error: GET: interfaces/john: No interface\n\n"

func TestInterfaceCli(t *testing.T) {
	// Since this data is dynamic, we will test errors here.
	cliTest(true, false, "interfaces").run(t)

	cliTest(true, true, "interfaces", "show").run(t)
	cliTest(true, true, "interfaces", "show", "john", "john2").run(t)
	cliTest(false, true, "interfaces", "show", "john").run(t)

	cliTest(true, true, "interfaces", "exists").run(t)
	cliTest(true, true, "interfaces", "exists", "john", "john2").run(t)
	cliTest(false, true, "interfaces", "exists", "john").run(t)
	cliTest(false, true, "interfaces", "exists", "john", "john2").run(t)
}
