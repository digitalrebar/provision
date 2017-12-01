package cli

import (
	"os"
	"testing"
)

func TestPrefsCli(t *testing.T) {

	var prefsSetEmptyJSONString string = "{}"

	var prefsSetStdinBadJSONString = "fred\n"
	var prefsSetStdinJSONString = `{
  "defaultBootEnv": "local3",
  "unknownBootEnv": "ignore"
}
`

	if err := os.MkdirAll("bootenvs", 0755); err != nil {
		t.Errorf("FAIL: Failed to create bootenvs dir: %v\n", err)
	}
	if err := os.Symlink("../test-data/local3.yml", "bootenvs/local3.yml"); err != nil {
		t.Errorf("FAIL: Failed to create link to local3.yml: %v\n", err)
	}

	if err := os.MkdirAll("templates", 0755); err != nil {
		t.Errorf("FAIL: Failed to create templates dir: %v\n", err)
	}
	tmpls := []string{"local3-pxelinux.tmpl", "local3-elilo.tmpl", "local3-ipxe.tmpl"}
	for _, tmpl := range tmpls {
		if err := os.Symlink("../test-data/"+tmpl, "templates/"+tmpl); err != nil {
			t.Errorf("FAIL: Failed to create link to %s: %v\n", tmpl, err)
		}
	}

	cliTest(false, false, "bootenvs", "install", "bootenvs/local3.yml").run(t)
	cliTest(true, false, "prefs").run(t)
	cliTest(false, false, "prefs", "list").run(t)
	cliTest(true, true, "prefs", "set").run(t)
	cliTest(true, true, "prefs", "set", "john", "john2", "john3").run(t)
	cliTest(false, true, "prefs", "set", "john").run(t)
	// Set empty hash - should result in no changes
	cliTest(false, false, "prefs", "set", prefsSetEmptyJSONString).run(t)
	cliTest(false, false, "prefs", "list").run(t)
	cliTest(false, false, "prefs", "set", "defaultBootEnv", "local3").run(t)
	cliTest(false, false, "prefs", "list").run(t)
	cliTest(false, true, "prefs", "set", "defaultBootEnv", "illegal").run(t)
	cliTest(false, false, "prefs", "list").run(t)
	cliTest(false, true, "prefs", "set", "baseTokenSecret", "illegal").run(t)
	cliTest(false, true, "prefs", "set", "baseTokenSecret", "illegalillegalillegalillegalillegal").run(t)
	cliTest(false, true, "prefs", "set", "knownTokenTimeout", "illegal").run(t)
	cliTest(false, true, "prefs", "set", "unknownTokenTimeout", "illegal").run(t)
	cliTest(false, false, "prefs", "set", "knownTokenTimeout", "5000").run(t)
	cliTest(false, false, "prefs", "set", "unknownTokenTimeout", "7000").run(t)
	cliTest(false, false, "prefs", "set", "debugRenderer", "1", "debugDhcp", "2", "debugBootEnv", "1").run(t)
	cliTest(false, false, "prefs", "list").run(t)
	cliTest(false, true, "prefs", "set", "greg", "ignore").run(t)
	cliTest(false, false, "prefs", "list").run(t)
	cliTest(false, true, "prefs", "set", "-").Stdin(prefsSetStdinBadJSONString).run(t)
	cliTest(false, false, "prefs", "list").run(t)
	cliTest(false, false, "prefs", "set", "-").Stdin(prefsSetStdinJSONString).run(t)
	cliTest(false, false, "prefs", "list").run(t)
	cliTest(false, true, "bootenvs", "destroy", "local3").run(t)
	cliTest(false, false, "prefs", "set", "defaultBootEnv", "local", "debugDhcp", "0", "debugRenderer", "0", "debugBootEnv", "0").run(t)
	cliTest(false, false, "bootenvs", "destroy", "local3").run(t)
	cliTest(false, false, "templates", "destroy", "local3-pxelinux.tmpl").run(t)
	cliTest(false, false, "templates", "destroy", "local3-elilo.tmpl").run(t)
	cliTest(false, false, "templates", "destroy", "local3-ipxe.tmpl").run(t)

	os.RemoveAll("bootenvs")
	os.RemoveAll("templates")
	os.RemoveAll("isos")
	os.RemoveAll("ic")
}
