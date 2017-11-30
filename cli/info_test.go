package cli

import "testing"

var infoGetTooManyArgsErrorString string = "Error: drpcli info get [flags] requires no arguments"

func TestInfoCli(t *testing.T) {
	// Since this data is dynamic, we will test errors here.
	cliTest(true, false, "info").run(t)
	cliTest(true, true, "info", "get", "john2").run(t)
}
