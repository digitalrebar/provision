package cli

// TODO: Add validations around templates and content checks.

import (
	"testing"
)

func TestTaskCli(t *testing.T) {
	var taskCreateBadJSONString = "{asdgasdg}"

	var taskCreateInputString string = `{
  "Name": "john",
  "OptionalParams": [],
  "RequiredParams": [],
  "Templates": []
}
`
	var taskUpdateBadJSONString = "asdgasdg"
	var taskUpdateInputString string = `{
  "OptionalParams": [ "jillparam" ]
}
`

	cliTest(true, false, "tasks").run(t)
	cliTest(false, false, "tasks", "list").run(t)
	cliTest(true, true, "tasks", "create").run(t)
	cliTest(true, true, "tasks", "create", "john", "john2").run(t)
	cliTest(false, true, "tasks", "create", taskCreateBadJSONString).run(t)
	cliTest(false, false, "tasks", "create", taskCreateInputString).run(t)
	cliTest(false, true, "tasks", "create", taskCreateInputString).run(t)
	cliTest(false, false, "tasks", "list").run(t)
	cliTest(false, false, "tasks", "list", "Name=fred").run(t)
	cliTest(false, false, "tasks", "list", "Name=john").run(t)
	cliTest(true, true, "tasks", "show").run(t)
	cliTest(true, true, "tasks", "show", "john", "john2").run(t)
	cliTest(false, true, "tasks", "show", "jill").run(t)
	cliTest(false, false, "tasks", "show", "john").run(t)
	cliTest(true, true, "tasks", "exists").run(t)
	cliTest(true, true, "tasks", "exists", "john", "john2").run(t)
	cliTest(false, true, "tasks", "exists", "jill").run(t)
	cliTest(false, false, "tasks", "exists", "john").run(t)
	cliTest(true, true, "tasks", "update").run(t)
	cliTest(true, true, "tasks", "update", "john", "john2", "john3").run(t)
	cliTest(false, true, "tasks", "update", "john", taskUpdateBadJSONString).run(t)
	cliTest(false, true, "tasks", "update", "jill", taskUpdateInputString).run(t)
	cliTest(false, false, "tasks", "update", "john", taskUpdateInputString).run(t)
	cliTest(false, false, "tasks", "show", "john").run(t)
	cliTest(false, false, "tasks", "show", "john").run(t)
	cliTest(true, true, "tasks", "destroy").run(t)
	cliTest(true, true, "tasks", "destroy", "john", "june").run(t)
	cliTest(false, false, "tasks", "destroy", "john").run(t)
	cliTest(false, true, "tasks", "destroy", "jill").run(t)
	cliTest(false, false, "tasks", "list").run(t)
	cliTest(false, false, "tasks", "create", "-").Stdin(taskCreateInputString + "\n").run(t)
	cliTest(false, false, "tasks", "list").run(t)
	cliTest(false, false, "tasks", "update", "john", "-").Stdin(taskUpdateInputString + "\n").run(t)
	cliTest(false, false, "tasks", "show", "john").run(t)
	cliTest(false, false, "tasks", "destroy", "john").run(t)
	cliTest(false, true, "tasks", "create", "-").Stdin(`---
Name: mixedOSes
Templates:
  - Name: t1
    Contents: '1'
  - Name: t2
    Contents: '2'
    Meta:
      OS: any`).run(t)
	cliTest(false, true, "tasks", "create", "-").Stdin(`---
Name: badOS
Templates:
  - Name: t2
    Contents: '2'
    Meta:
      OS: sithisOS`).run(t)
	cliTest(false, false, "tasks", "create", "-").Stdin(`---
Name: multiGoodOS
Templates:
  - Name: t1
    Contents: '1'
    Meta:
      OS: darwin
  - Name: t2
    Contents: '2'
    Meta:
      OS: linux`).run(t)
	cliTest(false, true, "tasks", "create", "-").Stdin(`---
Name: multiBadCommaOS
Templates:
  - Name: t1
    Contents: '1'
    Meta:
      OS: darwin,netbsd
  - Name: t2
    Contents: '2'
    Meta:
      OS: linux,sithisOS`).run(t)
	cliTest(false, false, "tasks", "create", "-").Stdin(`---
Name: multiNoOS
Templates:
  - Name: t1
    Contents: '1'
  - Name: t2
    Contents: '2'`).run(t)
	cliTest(false, false, "tasks", "destroy", "multiGoodOS").run(t)
	cliTest(false, false, "tasks", "destroy", "multiNoOS").run(t)
	cliTest(false, false, "tasks", "list").run(t)
	verifyClean(t)
}

func TestTaskPrereqs(t *testing.T) {
	cliTest(false, false, "tasks", "create", "badfoo").run(t)
	cliTest(false, false, "tasks", "create", "-").Stdin(`---
Name: badbar
Prerequisites:
  - badfoo`).run(t)
	cliTest(false, false, "tasks", "create", "-").Stdin(`---
Name: badbaz
Prerequisites:
  - badbar`).run(t)
	cliTest(false, true, "tasks", "update", "badfoo", `{"Prerequisites":["badbaz"]}`).run(t)
	cliTest(false, true, "tasks", "update", "badbar", `{"Prerequisites":["badfoo","badbaz"]}`).run(t)
	cliTest(false, false, "tasks", "create", "foo1").run(t)
	cliTest(false, false, "tasks", "create", "foo2").run(t)
	cliTest(false, false, "tasks", "create", "foo3").run(t)
	cliTest(false, false, "tasks", "create", "foo4").run(t)
	cliTest(false, false, "tasks", "create", "foo5").run(t)
	cliTest(false, false, "tasks", "create", "foo6").run(t)
	cliTest(false, false, "tasks", "create", "foo7").run(t)
	cliTest(false, false, "tasks", "create", "foo8").run(t)
	cliTest(false, false, "tasks", "create", "foo9").run(t)
	cliTest(false, false, "tasks", "create", "-").Stdin(`
Name: bar1
Prerequisites:
  - foo1
  - foo2
  - foo3`).run(t)
	cliTest(false, false, "tasks", "create", "-").Stdin(`
Name: bar2
Prerequisites:
  - foo3
  - foo2
  - foo1`).run(t)
	cliTest(false, false, "tasks", "create", "-").Stdin(`
Name: bar3
Prerequisites:
  - foo4
  - foo5
  - foo6`).run(t)
	cliTest(false, false, "tasks", "create", "-").Stdin(`
Name: bar4
Prerequisites:
  - bar3
  - foo1
  - foo5
  - foo7`).run(t)
	cliTest(false, false, "stages", "create", "-").Stdin(`
Name: flat1
Tasks:
  - foo1
  - foo4
  - bar1
  - bar2
  - bar4`).run(t)
	cliTest(false, false, "machines", "create", "bob").run(t)
	// tasks should wind up with
	// foo1 foo4 foo2 foo3 bar1 bar2 foo5 foo6 bar3 foo7 bar4
	cliTest(false, false, "machines", "update", "Name:bob", `{"Stage":"flat1"}`).run(t)
	cliTest(false, false, "bootenvs", "create", "three").run(t)
	cliTest(false, false, "bootenvs", "create", "two").run(t)
	cliTest(false, false, "stages", "create", "-").Stdin(`
Name: flat2
BootEnv: two
Tasks:
  - foo1
  - foo4
  - bar1
  - bar2
  - bar4`).run(t)
	cliTest(false, false, "stages", "create", "-").Stdin(`
Name: flat3
BootEnv: three
Tasks:
  - foo1
  - foo4
  - bar1
  - bar2
  - bar4`).run(t)
	cliTest(false, false, "workflows", "create", "-").Stdin(`
Name: wfPrereqs
Stages:
  - flat2
  - flat1
  - flat3
`).run(t)
	cliTest(false, false, "machines", "update", "Name:bob", `{"Workflow":"wfPrereqs"}`).run(t)
	cliTest(false, false, "machines", "update", "Name:bob", `{"Workflow":"","Stage":"none"}`).run(t)
	cliTest(false, false, "tasks", "destroy", "foo2").run(t)
	cliTest(false, true, "machines", "update", "Name:bob", `{"Workflow":"","Stage":"flat2"}`).run(t)
	cliTest(false, false, "machines", "destroy", "Name:bob").run(t)
	cliTest(false, false, "workflows", "destroy", "wfPrereqs").run(t)
	cliTest(false, false, "stages", "destroy", "flat1").run(t)
	cliTest(false, false, "stages", "destroy", "flat2").run(t)
	cliTest(false, false, "stages", "destroy", "flat3").run(t)
	cliTest(false, false, "bootenvs", "destroy", "two").run(t)
	cliTest(false, false, "bootenvs", "destroy", "three").run(t)
	cliTest(false, false, "tasks", "destroy", "bar4").run(t)
	cliTest(false, false, "tasks", "destroy", "bar3").run(t)
	cliTest(false, false, "tasks", "destroy", "bar2").run(t)
	cliTest(false, false, "tasks", "destroy", "bar1").run(t)
	cliTest(false, false, "tasks", "destroy", "foo9").run(t)
	cliTest(false, false, "tasks", "destroy", "foo8").run(t)
	cliTest(false, false, "tasks", "destroy", "foo7").run(t)
	cliTest(false, false, "tasks", "destroy", "foo6").run(t)
	cliTest(false, false, "tasks", "destroy", "foo5").run(t)
	cliTest(false, false, "tasks", "destroy", "foo4").run(t)
	cliTest(false, false, "tasks", "destroy", "foo3").run(t)
	cliTest(false, false, "tasks", "destroy", "foo1").run(t)
	cliTest(false, false, "tasks", "destroy", "badbaz").run(t)
	cliTest(false, false, "tasks", "destroy", "badfoo").run(t)
	cliTest(false, false, "tasks", "destroy", "badbar").run(t)
	cliTest(false, false, "tasks", "list").run(t)
	verifyClean(t)
}

func TestTasksInContentBundles(t *testing.T) {
	cliTest(false, true, "contents", "upload", "-").Stdin(`
meta:
  Name: circle2
sections:
  tasks:
    t1:
      Name: t1
      Prerequisites:
        - t2
    t2:
      Name: t2
      Prerequisites:
       - t1`).run(t)
	cliTest(false, true, "contents", "upload", "-").Stdin(`
meta:
  Name: circle3
sections:
  tasks:
    t1:
      Name: t1
      Prerequisites:
        - t3
    t2:
      Name: t2
      Prerequisites:
       - t1
    t3:
      Name: t3
      Prerequisites:
       - t2`).run(t)
	cliTest(false, true, "contents", "upload", "-").Stdin(`
meta:
  Name: forker
sections:
  tasks:
    t1:
      Name: t1
      Prerequisites:
        - t2
        - t3
        - t4
    t2:
      Name: t2
      Prerequisites:
       - t3
    t3:
      Name: t3
      Prerequisites:
       - t2
    t4:
      Name: t4
      Prerequisites:
       - t4`).run(t)
	cliTest(false, true, "contents", "upload", "-").Stdin(`
meta:
  Name: missing
sections:
  tasks:
    t1:
      Name: t1
      Prerequisites:
        - f1`).run(t)

}
