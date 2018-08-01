package cli

import (
	"testing"
)

func TestProcessJobsCli(t *testing.T) {
	actuallyPowerThings = false

	cliTest(false, false, "machines", "create", machineCreateInputString).run(t)

	// Test basic process jobs cli
	cliTest(true, true, "machines", "processjobs", "--oneshot").run(t)
	cliTest(true, true, "machines", "processjobs", "p1", "p2", "p3", "--oneshot").run(t)
	cliTest(false, true, "machines", "processjobs", "p1", "--oneshot").run(t)
	cliTest(false, false, "machines", "processjobs", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "--oneshot").run(t)
	cliTest(false, false, "machines", "show", "3e7031fe-3062-45f1-835c-92541bc9cbd3").run(t)
	cliTest(false, false, "machines", "destroy", "3e7031fe-3062-45f1-835c-92541bc9cbd3").run(t)
	verifyClean(t)
}
