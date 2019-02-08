package cli

import "testing"

func TestSystemCli(t *testing.T) {
	// Since this data is dynamic, we will test errors here.
	cliTest(false, false, "plugins", "create", machinePluginCreateString).run(t)
	cliTest(true, false, "system").run(t)
	cliTest(false, true, "system", "upgrade").run(t)
	cliTest(false, true, "system", "upgrade", "cows").run(t)
	cliTest(true, false, "system", "get", "john2").run(t)
	cliTest(false, false, "system", "actions").run(t)
	cliTest(false, true, "system", "action", "command").run(t)
	cliTest(false, false, "system", "action", "incrstatus").run(t)
	cliTest(false, true, "system", "runaction").run(t)
	cliTest(false, true, "system", "runaction", "command").run(t)
	cliTest(false, false, "system", "runaction", "incrstatus").run(t)
	cliTest(false, false, "plugins", "destroy", "incr").run(t)
}
