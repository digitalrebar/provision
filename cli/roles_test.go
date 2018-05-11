package cli

import (
	"encoding/json"
	"testing"

	"github.com/digitalrebar/provision/models"
)

func mkr(name string, args ...string) string {
	r := models.MakeRole(name, args...)
	buf, _ := json.Marshal(r)
	return string(buf)
}

func TestRoleCLI(t *testing.T) {
	cliTest(false, false, "contents", "upload", "-").Stdin(licenseLayer).run(t)
	cliTest(true, false, "roles").run(t)
	cliTest(false, false, "roles", "list").run(t)
	cliTest(true, true, "roles", "create").run(t)
	cliTest(true, true, "roles", "create", "john", "john2").run(t)
	cliTest(false, true, "roles", "create", "{foo").run(t)
	cliTest(false, true, "roles", "create", "[foo]").run(t)
	cliTest(false, true, "roles", "create", mkr("superuser", "*", "*", "*")).run(t)
	cliTest(false, true, "roles", "create", mkr("noScope", "", "", "")).run(t)
	cliTest(false, true, "roles", "create", mkr("noAction", "machines", "", "")).run(t)
	cliTest(false, true, "roles", "create", mkr("", "machines", "list", "noName")).run(t)
	cliTest(false, true, "roles", "create", mkr("badScope", "bar", "", "")).run(t)
	cliTest(false, true, "roles", "create", mkr("badAction", "machines", "bar", "")).run(t)
	cliTest(false, false, "roles", "create", mkr("validButUseless", "*", "*", "")).run(t)
	cliTest(false, false, "roles", "show", "superuser").run(t)
	cliTest(false, false, "roles", "show", "validButUseless").run(t)
	cliTest(false, true, "roles", "show", "badAction").run(t)
	cliTest(false, false, "roles", "list").run(t)
	cliTest(false, true, "roles", "destroy", "badAction").run(t)
	cliTest(false, true, "roles", "destroy", "superuser").run(t)
	cliTest(false, false, "roles", "destroy", "validButUseless").run(t)
	cliTest(false, false, "roles", "list").run(t)
	cliTest(false, false, "contents", "destroy", "rackn-license").run(t)
	verifyClean(t)
}
