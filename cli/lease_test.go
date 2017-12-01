package cli

import (
	"testing"
)

func TestLeaseCli(t *testing.T) {
	cliTest(true, false, "leases").run(t)
	cliTest(false, false, "leases", "list").run(t)
}
