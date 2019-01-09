package cli

import (
	"runtime"
	"testing"
)

func TestWorkflowCli(t *testing.T) {
	stageCreateInput := `{
  "Name": "john",
  "BootEnv": "local"
}
`
	bootEnvCreateInput := `
Name: Fred
Templates:
  - Name: ipxe
    Path: /ipxe
    Contents: 'foo'
  - Name: ipxe-mac
    Path: /ipxe-mac
    Contents: 'bar'
`

	stage2CreateInput := `{
  "Name": "james",
  "BootEnv": "Fred"
}
`

	workflow1CreateInput := `
Name: wf1
Stages: [john, james]
`
	workflow2CreateInput := `
Name: wf2
Stages: [james, john]
`
	workflow3CreateInput := `
Name: wf3
Stages: [james, local]
`
	workflow4CreateInput := `
Name: wf4
Stages: [missing]
`
	m0 := `
Name: m0
`
	m1 := `
Name: m1
Workflow: wf1
`
	m2 := `
Name: m2
Workflow: wf2
`
	m3 := `
Name: m3
Workflow: wf3
`
	m4 := `
Name: m4
Workflow: wf4
`
	cliTest(true, false, "workflows").run(t)
	cliTest(false, false, "workflows", "list").run(t)
	cliTest(true, true, "workflows", "create").run(t)
	cliTest(true, true, "workflows", "create", "john", "john2").run(t)
	cliTest(false, false, "bootenvs", "create", bootEnvCreateInput).run(t)
	cliTest(false, false, "stages", "create", stageCreateInput).run(t)
	cliTest(false, false, "stages", "create", stage2CreateInput).run(t)
	cliTest(false, true, "workflows", "create", `{"asdg"`).run(t)
	cliTest(false, true, "workflows", "create", "{}").run(t)
	cliTest(false, false, "workflows", "create", workflow1CreateInput).run(t)
	cliTest(false, false, "workflows", "create", workflow2CreateInput).run(t)
	cliTest(false, false, "workflows", "create", workflow3CreateInput).run(t)
	cliTest(false, true, "workflows", "create", workflow3CreateInput).run(t)
	cliTest(false, false, "workflows", "create", workflow4CreateInput).run(t)
	cliTest(false, true, "prefs", "set", "defaultWorkflow", "foo").run(t)
	cliTest(false, false, "machines", "create", m0).run(t)
	cliTest(false, false, "machines", "create", m1).run(t)
	cliTest(false, false, "machines", "create", m2).run(t)
	cliTest(false, false, "machines", "create", m3).run(t)
	cliTest(false, true, "machines", "create", m4).run(t)
	cliTest(false, false, "prefs", "set", "defaultWorkflow", "wf3").run(t)
	cliTest(false, false, "machines", "create", "m4").run(t)
	cliTest(false, false, "machines", "destroy", "Name:m4").run(t)
	cliTest(false, false, "prefs", "set", "defaultWorkflow", "").run(t)
	cliTest(false, false, "machines", "create", "m4").run(t)
	cliTest(false, false, "machines", "destroy", "Name:m4").run(t)
	cliTest(false, false, "bootenvs", "list", "sort", "Name").run(t)
	cliTest(false, false, "stages", "list", "sort", "Name").run(t)
	cliTest(false, false, "workflows", "list", "sort", "Name").run(t)
	cliTest(false, false, "machines", "list", "sort", "Name").run(t)
	cliTest(false, false, "machines", "update", "Name:m0", `{"Workflow":"wf1"}`).run(t)
	cliTest(false, false, "machines", "update", "Name:m0", `{"Workflow":"wf2"}`).run(t)
	cliTest(false, false, "machines", "update", "Name:m0", `{"Workflow":"wf3"}`).run(t)
	cliTest(false, true, "machines", "update", "Name:m0", `{"Workflow":"wf4"}`).run(t)
	cliTest(false, false, "machines", "update", "Name:m0", `{"Workflow":""}`).run(t)

	// Clean up
	cliTest(false, false, "machines", "destroy", "Name:m3").run(t)
	cliTest(false, false, "machines", "destroy", "Name:m2").run(t)
	cliTest(false, false, "machines", "destroy", "Name:m1").run(t)
	cliTest(false, false, "machines", "destroy", "Name:m0").run(t)
	cliTest(false, false, "workflows", "destroy", "wf4").run(t)
	cliTest(false, false, "workflows", "destroy", "wf3").run(t)
	cliTest(false, false, "workflows", "destroy", "wf2").run(t)
	cliTest(false, false, "workflows", "destroy", "wf1").run(t)
	cliTest(false, false, "stages", "destroy", "james").run(t)
	cliTest(false, false, "stages", "destroy", "john").run(t)
	cliTest(false, false, "bootenvs", "destroy", "Fred").run(t)
	verifyClean(t)
}

func TestWorkflowSwitch(t *testing.T) {
	cliTest(false, false, "tasks", "create", "task1").run(t)
	cliTest(false, false, "tasks", "create", "task2").run(t)
	cliTest(false, false, "tasks", "create", "task3").run(t)
	cliTest(false, false, "tasks", "create", "task4").run(t)
	cliTest(false, false, "stages", "create", "-").Stdin(`{"Name":"stage1","Tasks":["task1","task2","task3"]}`).run(t)
	cliTest(false, false, "stages", "create", "-").Stdin(`{"Name":"stage2","Tasks":["task2","task3","task4"]}`).run(t)
	cliTest(false, false, "stages", "create", "-").Stdin(`{"Name":"stage3","Tasks":["task3","task2","task1"]}`).run(t)
	cliTest(false, false, "stages", "create", "-").Stdin(`{"Name":"stage4","Tasks":["task4","task3","task2"]}`).run(t)
	cliTest(false, false, "workflows", "create", "-").Stdin(`{"Name":"wf1","Stages":["stage1","stage2"]}`).run(t)
	cliTest(false, false, "workflows", "create", "-").Stdin(`{"Name":"wf2","Stages":["stage3","stage4"]}`).run(t)
	cliTest(false, false, "machines", "create", "-").Stdin(`{"Name":"m1","Workflow":"wf1","Runnable":true}`).run(t)
	cliTest(false, false, "machines", "jobs", "create", "Name:m1").run(t)
	cliTest(false, false, "machines", "jobs", "state", "Name:m1", "to", "running").run(t)
	cliTest(false, false, "machines", "jobs", "state", "Name:m1", "to", "finished").run(t)
	cliTest(false, false, "machines", "jobs", "create", "Name:m1").run(t)
	cliTest(false, false, "machines", "jobs", "state", "Name:m1", "to", "running").run(t)
	cliTest(false, false, "machines", "workflow", "Name:m1", "wf2").run(t)
	cliTest(false, false, "machines", "jobs", "state", "Name:m1", "to", "finished").run(t)
	cliTest(false, false, "machines", "jobs", "create", "Name:m1").run(t)
	cliTest(false, false, "machines", "jobs", "state", "Name:m1", "to", "running").run(t)
	cliTest(false, false, "machines", "jobs", "state", "Name:m1", "to", "finished").run(t)
	cliTest(false, false, "machines", "jobs", "create", "Name:m1").run(t)
	cliTest(false, false, "machines", "jobs", "state", "Name:m1", "to", "running").run(t)
	cliTest(false, false, "machines", "workflow", "Name:m1", "wf1").run(t)
	cliTest(false, false, "machines", "jobs", "state", "Name:m1", "to", "finished").run(t)
	cliTest(false, false, "machines", "jobs", "create", "Name:m1").run(t)
	cliTest(false, false, "machines", "jobs", "state", "Name:m1", "to", "failed").run(t)
	cliTest(false, false, "machines", "show", "Name:m1").run(t)
	cliTest(false, false, "machines", "deletejobs", "Name:m1").run(t)
	cliTest(false, false, "machines", "destroy", "Name:m1").run(t)
	cliTest(false, false, "workflows", "destroy", "wf2").run(t)
	cliTest(false, false, "workflows", "destroy", "wf1").run(t)
	cliTest(false, false, "stages", "destroy", "stage1").run(t)
	cliTest(false, false, "stages", "destroy", "stage2").run(t)
	cliTest(false, false, "stages", "destroy", "stage3").run(t)
	cliTest(false, false, "stages", "destroy", "stage4").run(t)
	cliTest(false, false, "tasks", "destroy", "task1").run(t)
	cliTest(false, false, "tasks", "destroy", "task2").run(t)
	cliTest(false, false, "tasks", "destroy", "task3").run(t)
	cliTest(false, false, "tasks", "destroy", "task4").run(t)
	verifyClean(t)
}

func TestWorkflowAgent(t *testing.T) {
	if runtime.GOARCH != "amd64" {
		t.Logf("Agent tests only run on amd64")
		return
	}
	cliTest(false, false, "tasks", "create", "-").Stdin(`---
Name: task1
Templates:
  - Contents: |
      #!/usr/bin/env bash
      exit 0
    Name: task1`).run(t)
	cliTest(false, false, "tasks", "create", "-").Stdin(`---
Name: task2
Templates:
  - Contents: |
      #!/usr/bin/env bash
      if [[ $(uname -s) == Darwin ]] ; then
        LOS=darwin
      else
        LOS=linux
      fi
      DRPCLI="$GOPATH/src/github.com/digitalrebar/provision/bin/$LOS/amd64/drpcli"
      if [[ ! -x $DRPCLI ]]; then
         echo "Missing drpcli.  Please run tools/build.sh before running tests"
         exit 1
      fi
      "$DRPCLI" machines workflow Name:m1 wf2 &>/dev/null
    Name: task2`).run(t)
	cliTest(false, false, "tasks", "create", "-").Stdin(`---
Name: task3
Templates:
  - Contents: |
      #!/usr/bin/env bash
      echo "Shouldn't get here 0"
      exit 1
    Name: task3`).run(t)
	cliTest(false, false, "tasks", "create", "-").Stdin(`---
Name: task4
Templates:
  - Contents: |
      #!/usr/bin/env bash
      echo "Should exit here"
      sleep 2
      exit 1
    Name: task4`).run(t)
	cliTest(false, false, "tasks", "create", "-").Stdin(`---
Name: task5
Templates:
  - Contents: |
      #!/usr/bin/env bash
      echo "Should not get here 1"
      exit 1
    Name: task5`).run(t)
	cliTest(false, false, "tasks", "create", "-").Stdin(`---
Name: task6
Templates:
  - Contents: |
      #!/usr/bin/env bash
      echo "Should not get here 2"
      exit 1
    Name: task6`).run(t)
	cliTest(false, false, "stages", "create", "-").Stdin(`{"Name":"stage1","Tasks":["task1","task2","task3"]}`).run(t)
	cliTest(false, false, "stages", "create", "-").Stdin(`{"Name":"stage2","Tasks":["task4","task5","task6"]}`).run(t)
	cliTest(false, false, "workflows", "create", "-").Stdin(`{"Name":"wf1","Stages":["stage1"]}`).run(t)
	cliTest(false, false, "workflows", "create", "-").Stdin(`{"Name":"wf2","Stages":["stage2"]}`).run(t)
	cliTest(false, false, "machines", "create", "-").Stdin(`{"Name":"m1","Workflow":"wf1","Runnable":true}`).run(t)
	cliTest(false, false, "machines", "processjobs", "Name:m1", "--oneshot", "--exit-on-failure", "--").run(t)
	cliTest(false, false, "machines", "currentlog", "Name:m1").run(t)
	cliTest(false, false, "machines", "deletejobs", "Name:m1").run(t)
	cliTest(false, false, "machines", "destroy", "Name:m1").run(t)
	cliTest(false, false, "workflows", "destroy", "wf2").run(t)
	cliTest(false, false, "workflows", "destroy", "wf1").run(t)
	cliTest(false, false, "stages", "destroy", "stage1").run(t)
	cliTest(false, false, "stages", "destroy", "stage2").run(t)
	cliTest(false, false, "tasks", "destroy", "task1").run(t)
	cliTest(false, false, "tasks", "destroy", "task2").run(t)
	cliTest(false, false, "tasks", "destroy", "task3").run(t)
	cliTest(false, false, "tasks", "destroy", "task4").run(t)
	cliTest(false, false, "tasks", "destroy", "task5").run(t)
	cliTest(false, false, "tasks", "destroy", "task6").run(t)
	verifyClean(t)
}
