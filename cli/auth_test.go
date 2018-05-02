package cli

import "testing"

func TestAuth(t *testing.T) {
	cliTest(false, false, "roles", "create", `{"Name":"stage","Claims":[{"Scope":"stages", "Action":"*","Specific":"*"}]}`).run(t)
	cliTest(false, false, "roles", "create", `{"Name":"task","Claims":[{"Scope":"tasks", "Action":"*","Specific":"*"}]}`).run(t)
	uMap := map[string]string{
		"t1-0": `{"Roles":["stage"]}`,
		"t1-1": `{"Roles":["stage","task"]}`,
		"t1-2": `{"Roles":["task"]}`,
		"t2-0": `{"Roles":["stage"]}`,
		"t2-1": `{"Roles":["stage","task"]}`,
		"t2-2": `{"Roles":["task"]}`,
	}
	for u, r := range uMap {
		cliTest(false, false, "users", "create", u).run(t)
		cliTest(false, false, "users", "password", u, "foo").run(t)
		cliTest(false, false, "users", "update", u, r).run(t)
	}
	cliTest(false, false, "roles", "list").run(t)
	cliTest(false, false, "users", "list").run(t)
	// user list is not in a role, so no dice.
	for u := range uMap {
		cliTest(false, true, "users", "list", "-T", "", "-U", u, "-P", "foo").run(t)
	}
	// Make some stages and tasks
	cliTest(false, false, "tasks", "create", "task1").run(t)
	cliTest(false, false, "tasks", "create", "task2").run(t)
	cliTest(false, false, "tasks", "create", "task3").run(t)
	cliTest(false, false, "stages", "create", "stage1").run(t)
	cliTest(false, false, "stages", "create", "stage2").run(t)
	cliTest(false, false, "stages", "create", "stage3").run(t)
	// Test to make sure auth restrictions are parsing properly
	cliTest(false, false, "stages", "list", "-T", "", "-U", "t1-0", "-P", "foo").run(t)
	cliTest(false, true, "tasks", "list", "-T", "", "-U", "t1-0", "-P", "foo").run(t)
	cliTest(false, false, "stages", "list", "-T", "", "-U", "t1-1", "-P", "foo").run(t)
	cliTest(false, false, "tasks", "list", "-T", "", "-U", "t1-1", "-P", "foo").run(t)
	cliTest(false, true, "stages", "list", "-T", "", "-U", "t1-2", "-P", "foo").run(t)
	cliTest(false, false, "tasks", "list", "-T", "", "-U", "t1-2", "-P", "foo").run(t)
	cliTest(false, false, "stages", "list", "-T", "", "-U", "t2-0", "-P", "foo").run(t)
	cliTest(false, true, "tasks", "list", "-T", "", "-U", "t2-0", "-P", "foo").run(t)
	cliTest(false, false, "stages", "list", "-T", "", "-U", "t2-1", "-P", "foo").run(t)
	cliTest(false, false, "tasks", "list", "-T", "", "-U", "t2-1", "-P", "foo").run(t)
	cliTest(false, true, "stages", "list", "-T", "", "-U", "t2-2", "-P", "foo").run(t)
	cliTest(false, false, "tasks", "list", "-T", "", "-U", "t2-2", "-P", "foo").run(t)
	// Make a couple of tenants with the existing data
	cliTest(false, false, "tenants", "create", `
Name: tenant1
Members:
  stages: [stage1, stage3]
  tasks: [task1, task3]
Users: [t1-0, t1-1, t1-2]
`).run(t)
	cliTest(false, false, "tenants", "create", `
Name: tenant2
Members:
  stages: [stage2, stage3]
  tasks: [task2, task3]
Users: [t2-0, t2-1, t2-2]
`).run(t)
	cliTest(false, false, "tenants", "list").run(t)
	// Test to make sure users in tenants can only see what they are allowed to see
	cliTest(false, false, "stages", "list", "-T", "", "-U", "t1-0", "-P", "foo").run(t)
	cliTest(false, true, "tasks", "list", "-T", "", "-U", "t1-0", "-P", "foo").run(t)
	cliTest(false, false, "stages", "list", "-T", "", "-U", "t1-1", "-P", "foo").run(t)
	cliTest(false, false, "tasks", "list", "-T", "", "-U", "t1-1", "-P", "foo").run(t)
	cliTest(false, true, "stages", "list", "-T", "", "-U", "t1-2", "-P", "foo").run(t)
	cliTest(false, false, "tasks", "list", "-T", "", "-U", "t1-2", "-P", "foo").run(t)
	cliTest(false, false, "stages", "list", "-T", "", "-U", "t2-0", "-P", "foo").run(t)
	cliTest(false, true, "tasks", "list", "-T", "", "-U", "t2-0", "-P", "foo").run(t)
	cliTest(false, false, "stages", "list", "-T", "", "-U", "t2-1", "-P", "foo").run(t)
	cliTest(false, false, "tasks", "list", "-T", "", "-U", "t2-1", "-P", "foo").run(t)
	cliTest(false, true, "stages", "list", "-T", "", "-U", "t2-2", "-P", "foo").run(t)
	cliTest(false, false, "tasks", "list", "-T", "", "-U", "t2-2", "-P", "foo").run(t)
	// Refuse to delete tenants with occupants
	cliTest(false, true, "tenants", "destroy", "tenant1").run(t)
	cliTest(false, true, "tenants", "destroy", "tenant2").run(t)
	// Delete and recreate objects, make sure they wind up in the right tenants.
	cliTest(false, false, "stages", "destroy", "stage3", "-T", "", "-U", "t1-0", "-P", "foo").run(t)
	cliTest(false, false, "stages", "create", "stage3", "-T", "", "-U", "t1-0", "-P", "foo").run(t)
	cliTest(false, true, "stages", "destroy", "stage2", "-T", "", "-U", "t1-0", "-P", "foo").run(t)
	cliTest(false, true, "stages", "create", "stage2", "-T", "", "-U", "t1-0", "-P", "foo").run(t)
	cliTest(false, false, "tasks", "destroy", "task3", "-T", "", "-U", "t2-2", "-P", "foo").run(t)
	cliTest(false, false, "tasks", "create", "task3", "-T", "", "-U", "t2-2", "-P", "foo").run(t)
	cliTest(false, true, "tasks", "destroy", "task1", "-T", "", "-U", "t2-2", "-P", "foo").run(t)
	cliTest(false, true, "tasks", "create", "task1", "-T", "", "-U", "t2-2", "-P", "foo").run(t)
	cliTest(false, false, "tasks", "list", "-T", "", "-U", "t1-1", "-P", "foo").run(t)
	cliTest(false, false, "tasks", "list", "-T", "", "-U", "t2-1", "-P", "foo").run(t)
	cliTest(false, false, "stages", "list", "-T", "", "-U", "t1-1", "-P", "foo").run(t)
	cliTest(false, false, "stages", "list", "-T", "", "-U", "t2-1", "-P", "foo").run(t)
	cliTest(false, false, "tenants", "list").run(t)
	// Refuse to remove roles that a user is using
	cliTest(false, true, "roles", "destroy", "task").run(t)
	cliTest(false, true, "roles", "destroy", "stage").run(t)
	// Clean up
	for u := range uMap {
		cliTest(false, false, "users", "destroy", u).run(t)
	}
	cliTest(false, false, "tenants", "destroy", "tenant1").run(t)
	cliTest(false, false, "tenants", "destroy", "tenant2").run(t)
	cliTest(false, false, "stages", "destroy", "stage3").run(t)
	cliTest(false, false, "stages", "destroy", "stage2").run(t)
	cliTest(false, false, "stages", "destroy", "stage1").run(t)
	cliTest(false, false, "tasks", "destroy", "task3").run(t)
	cliTest(false, false, "tasks", "destroy", "task2").run(t)
	cliTest(false, false, "tasks", "destroy", "task1").run(t)
	cliTest(false, false, "roles", "destroy", "task").run(t)
	cliTest(false, false, "roles", "destroy", "stage").run(t)
	verifyClean(t)
}
