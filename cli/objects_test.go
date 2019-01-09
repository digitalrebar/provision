package cli

import "testing"

func TestObjectCli(t *testing.T) {
	// Since this data is dynamic, we will test errors here.
	cliTest(true, false, "objects").run(t)
	cliTest(false, false, "objects", "list").run(t)
}
