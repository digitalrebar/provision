package cli

import "testing"

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
