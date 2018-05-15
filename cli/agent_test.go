package cli

import "testing"

func TestAgent(t *testing.T) {
	// Make a noisy task that sleeps some to test log write coalescing
	cliTest(false, false, "tasks", "create", "-").Stdin(`---
Name: noisyTask
Templates:
  - Name: noisy
    Contents: |
      #!/usr/bin/env bash
      . ./helper
      # The internal buffer the logger uses is 64K, so make sure to overflow it a bit.
      for ((i=0;i<1026;i++)); do
         printf '%04d...........................................................\n' "$i"
      done
      echo "Pause"
      sleep 3
      for ((i=0;i<1026;i++)); do
         printf '%04d...........................................................\n' "$i"
      done
      sleep 3
      echo "Done"
      exit_stop
`).run(t)
	cliTest(false, false, "stages", "create", "-").Stdin(`---
Name: noisyStage
Tasks: [noisyTask]
`).run(t)
	cliTest(false, false, "machines", "create", "-").Stdin(`---
Name: phred
Uuid: c9196b77-deef-4c8e-8130-299b3e3d9a10
Stage: noisyStage
Runnable: true
`).run(t)
	// We need to log at debug level to make sure we catch the debug messages
	// indicating how much data was written to the log
	cliTest(false, false, "prefs", "set", "debugFrontend", "debug").run(t)
	cliTest(false, false, "machines", "processjobs", "c9196b77-deef-4c8e-8130-299b3e3d9a10", "--oneshot").run(t)
	cliTest(false, false, "prefs", "set", "debugFrontend", "warn").run(t)
	cliTest(false, false, "machines", "currentlog", "Name:phred").run(t)
	cliTest(false, false, "logs", "get").run(t)
	cliTest(false, false, "machines", "deletejobs", "Name:phred").run(t)
	cliTest(false, false, "machines", "destroy", "Name:phred").run(t)
	cliTest(false, false, "stages", "destroy", "noisyStage").run(t)
	cliTest(false, false, "tasks", "destroy", "noisyTask").run(t)
	verifyClean(t)
}
