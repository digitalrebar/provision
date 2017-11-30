package cli

import (
	"os"
	"testing"
)

var prefsDefaultListString = `{
  "baseTokenSecret": "token-secret-token-secret-token1",
  "debugBootEnv": "0",
  "debugDhcp": "0",
  "debugFrontend": "0",
  "debugPlugins": "0",
  "debugRenderer": "0",
  "defaultBootEnv": "local",
  "defaultStage": "none",
  "knownTokenTimeout": "3600",
  "systemGrantorSecret": "system-grantor-secret",
  "unknownBootEnv": "ignore",
  "unknownTokenTimeout": "600"
}
`

var prefsSetNoArgsErrorString string = "Error: prefs set either takes a single argument or a multiple of two, not 0"
var prefsSetOddArgsErrorString string = "Error: prefs set either takes a single argument or a multiple of two, not 3"
var prefsSetBadJSONErrorString string = "Error: Invalid prefs: error unmarshaling JSON: json: cannot unmarshal string into Go value of type map[string]string\n\n\n"

var prefsSetEmptyJSONString string = "{}"
var prefsSetJSONResponseString string = `{
  "baseTokenSecret": "token-secret-token-secret-token1",
  "debugBootEnv": "0",
  "debugDhcp": "0",
  "debugFrontend": "0",
  "debugPlugins": "0",
  "debugRenderer": "0",
  "defaultBootEnv": "local3",
  "defaultStage": "none",
  "knownTokenTimeout": "3600",
  "systemGrantorSecret": "system-grantor-secret",
  "unknownBootEnv": "ignore",
  "unknownTokenTimeout": "600"
}
`
var prefsSetIllegalJSONResponseString string = "Error: POST: prefs: defaultBootEnv: Bootenv illegal does not exist\n\n"
var prefsSetInvalidPrefResponseString string = "Error: POST: prefs: Unknown Preference greg\n\n"

var prefsChangedListString = `{
  "baseTokenSecret": "token-secret-token-secret-token1",
  "debugBootEnv": "0",
  "debugDhcp": "0",
  "debugFrontend": "0",
  "debugPlugins": "0",
  "debugRenderer": "0",
  "defaultBootEnv": "local3",
  "defaultStage": "none",
  "knownTokenTimeout": "3600",
  "systemGrantorSecret": "system-grantor-secret",
  "unknownBootEnv": "ignore",
  "unknownTokenTimeout": "600"
}
`

var prefsSetStdinBadJSONString = "fred\n"
var prefsSetStdinBadJSONErrorString = ""
var prefsSetStdinJSONString = `{
  "defaultBootEnv": "local3",
  "unknownBootEnv": "ignore"
}
`

var prefsSetBadKnownTokenTimeoutErrorString = "Error: POST: prefs: knownTokenTimeout: strconv.Atoi: parsing \"illegal\": invalid syntax\n\n"
var prefsSetBadUnknownTokenTimeoutErrorString = "Error: POST: prefs: unknownTokenTimeout: strconv.Atoi: parsing \"illegal\": invalid syntax\n\n"

var prefsKnownChangedListString = `{
  "baseTokenSecret": "token-secret-token-secret-token1",
  "debugBootEnv": "0",
  "debugDhcp": "0",
  "debugFrontend": "0",
  "debugPlugins": "0",
  "debugRenderer": "0",
  "defaultBootEnv": "local3",
  "defaultStage": "none",
  "knownTokenTimeout": "5000",
  "systemGrantorSecret": "system-grantor-secret",
  "unknownBootEnv": "ignore",
  "unknownTokenTimeout": "600"
}
`
var prefsBothPreDebugChangedListString = `{
  "baseTokenSecret": "token-secret-token-secret-token1",
  "debugBootEnv": "0",
  "debugDhcp": "0",
  "debugFrontend": "0",
  "debugPlugins": "0",
  "debugRenderer": "0",
  "defaultBootEnv": "local3",
  "defaultStage": "none",
  "knownTokenTimeout": "5000",
  "systemGrantorSecret": "system-grantor-secret",
  "unknownBootEnv": "ignore",
  "unknownTokenTimeout": "7000"
}
`
var prefsBothChangedListString = `{
  "baseTokenSecret": "token-secret-token-secret-token1",
  "debugBootEnv": "1",
  "debugDhcp": "2",
  "debugFrontend": "0",
  "debugPlugins": "0",
  "debugRenderer": "1",
  "defaultBootEnv": "local3",
  "defaultStage": "none",
  "knownTokenTimeout": "5000",
  "systemGrantorSecret": "system-grantor-secret",
  "unknownBootEnv": "ignore",
  "unknownTokenTimeout": "7000"
}
`

var prefsFinalListString = `{
  "baseTokenSecret": "token-secret-token-secret-token1",
  "debugBootEnv": "0",
  "debugDhcp": "0",
  "debugFrontend": "0",
  "debugPlugins": "0",
  "debugRenderer": "0",
  "defaultBootEnv": "local",
  "defaultStage": "none",
  "knownTokenTimeout": "5000",
  "systemGrantorSecret": "system-grantor-secret",
  "unknownBootEnv": "ignore",
  "unknownTokenTimeout": "7000"
}
`
var prefsFailedToDeleteBootenvErrorString = "Error: StillInUseError: bootenvs/local3: BootEnv local3 is the active defaultBootEnv, cannot remove it\n\n"
var prefSetBadBaseTokenSecretBadLengthErrorString = "Error: POST: prefs: baseTokenSecret: Must be 32 bytes long\n\n"

func TestPrefsCli(t *testing.T) {
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
