package cli

import "testing"

func TestLogsCli(t *testing.T) {
	// Since this data is dynamic, we will test errors here.
	cliTest(false, false, "logs").run(t)
	cliTest(false, false, "logs", "get").run(t)
}
