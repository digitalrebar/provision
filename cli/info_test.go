package cli

import "testing"

func TestInfoCli(t *testing.T) {
	// Since this data is dynamic, we will test errors here.
	cliTest(true, false, "info").run(t)
	cliTest(false, false, "info", "check").run(t)
	cliTest(true, true, "info", "get", "john2").run(t)
	cliTest(false, false, "info", "status").run(t)
}
