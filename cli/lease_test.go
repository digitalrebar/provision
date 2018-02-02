package cli

import (
	"testing"
)

func TestLeaseCli(t *testing.T) {
	cliTest(true, false, "leases").run(t)
	cliTest(false, false, "leases", "list").run(t)
	cliTest(false, true, "leases", "list", "Addr=fred").run(t)
	cliTest(false, false, "leases", "list", "Addr=1.1.1.1").run(t)
	cliTest(false, false, "leases", "list", "Token=1.1.1.1").run(t)
	cliTest(false, false, "leases", "list", "Token=11:22:33:44:55:66").run(t)
	cliTest(false, false, "leases", "list", "Strategy=MAC").run(t)
	cliTest(false, false, "leases", "list", "Strategy=COW").run(t)
	cliTest(false, false, "leases", "list", "State=sleep").run(t)
	cliTest(false, true, "leases", "list", "ExpireTime=fred").run(t)
	cliTest(false, false, "leases", "list", "ExpireTime=2006-01-02T15:04:05-07:00").run(t)
}
