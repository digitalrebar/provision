package cli

import "testing"

func TestWorkflowCli(t *testing.T) {
	stageCreateInput := `{
  "Name": "john",
  "BootEnv": "local"
}
`
	bootEnvCreateInput := `
Name: Fred
Kernel: lpxelinux.0
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
