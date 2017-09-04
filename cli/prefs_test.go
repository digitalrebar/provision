package cli

import (
	"os"
	"testing"
)

var prefsDefaultListString = `{
  "debugBootEnv": "0",
  "debugDhcp": "0",
  "debugFrontend": "0",
  "debugPlugins": "0",
  "debugRenderer": "0",
  "defaultBootEnv": "local",
  "defaultStage": "",
  "knownTokenTimeout": "3600",
  "unknownBootEnv": "ignore",
  "unknownTokenTimeout": "600"
}
`

var prefsSetNoArgsErrorString string = "Error: prefs set either takes a single argument or a multiple of two, not 0"
var prefsSetOddArgsErrorString string = "Error: prefs set either takes a single argument or a multiple of two, not 3"
var prefsSetBadJSONErrorString string = "Error: Invalid prefs: error unmarshaling JSON: json: cannot unmarshal string into Go value of type map[string]string\n\n\n"

var prefsSetEmptyJSONString string = "{}"
var prefsSetJSONResponseString string = `{
  "debugBootEnv": "0",
  "debugDhcp": "0",
  "debugFrontend": "0",
  "debugPlugins": "0",
  "debugRenderer": "0",
  "defaultBootEnv": "local3",
  "defaultStage": "",
  "knownTokenTimeout": "3600",
  "unknownBootEnv": "ignore",
  "unknownTokenTimeout": "600"
}
`
var prefsSetIllegalJSONResponseString string = "Error: defaultBootEnv: Bootenv illegal does not exist\n\n"
var prefsSetInvalidPrefResponseString string = "Error: Unknown Preference greg\n\n"

var prefsChangedListString = `{
  "debugBootEnv": "0",
  "debugDhcp": "0",
  "debugFrontend": "0",
  "debugPlugins": "0",
  "debugRenderer": "0",
  "defaultBootEnv": "local3",
  "defaultStage": "",
  "knownTokenTimeout": "3600",
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

var prefsSetBadKnownTokenTimeoutErrorString = "Error: Preference knownTokenTimeout: strconv.Atoi: parsing \"illegal\": invalid syntax\n\n"
var prefsSetBadUnknownTokenTimeoutErrorString = "Error: Preference unknownTokenTimeout: strconv.Atoi: parsing \"illegal\": invalid syntax\n\n"

var prefsKnownChangedListString = `{
  "debugBootEnv": "0",
  "debugDhcp": "0",
  "debugFrontend": "0",
  "debugPlugins": "0",
  "debugRenderer": "0",
  "defaultBootEnv": "local3",
  "defaultStage": "",
  "knownTokenTimeout": "5000",
  "unknownBootEnv": "ignore",
  "unknownTokenTimeout": "600"
}
`
var prefsBothPreDebugChangedListString = `{
  "debugBootEnv": "0",
  "debugDhcp": "0",
  "debugFrontend": "0",
  "debugPlugins": "0",
  "debugRenderer": "0",
  "defaultBootEnv": "local3",
  "defaultStage": "",
  "knownTokenTimeout": "5000",
  "unknownBootEnv": "ignore",
  "unknownTokenTimeout": "7000"
}
`
var prefsBothChangedListString = `{
  "debugBootEnv": "1",
  "debugDhcp": "2",
  "debugFrontend": "0",
  "debugPlugins": "0",
  "debugRenderer": "1",
  "defaultBootEnv": "local3",
  "defaultStage": "",
  "knownTokenTimeout": "5000",
  "unknownBootEnv": "ignore",
  "unknownTokenTimeout": "7000"
}
`

var prefsFinalListString = `{
  "debugBootEnv": "0",
  "debugDhcp": "0",
  "debugFrontend": "0",
  "debugPlugins": "0",
  "debugRenderer": "0",
  "defaultBootEnv": "local",
  "defaultStage": "",
  "knownTokenTimeout": "5000",
  "unknownBootEnv": "ignore",
  "unknownTokenTimeout": "7000"
}
`
var prefsFailedToDeleteBootenvErrorString = "Error: BootEnv local3 is the active defaultBootEnv, cannot remove it\n\n"

func TestPrefsCli(t *testing.T) {
	if err := os.MkdirAll("bootenvs", 0755); err != nil {
		t.Errorf("Failed to create bootenvs dir: %v\n", err)
	}
	if err := os.Symlink("../test-data/local3.yml", "bootenvs/local3.yml"); err != nil {
		t.Errorf("Failed to create link to local3.yml: %v\n", err)
	}

	if err := os.MkdirAll("templates", 0755); err != nil {
		t.Errorf("Failed to create templates dir: %v\n", err)
	}
	tmpls := []string{"local3-pxelinux.tmpl", "local3-elilo.tmpl", "local3-ipxe.tmpl"}
	for _, tmpl := range tmpls {
		if err := os.Symlink("../test-data/"+tmpl, "templates/"+tmpl); err != nil {
			t.Errorf("Failed to create link to %s: %v\n", tmpl, err)
		}
	}

	tests := []CliTest{
		CliTest{false, false, []string{"bootenvs", "install", "bootenvs/local3.yml"}, noStdinString, bootEnvInstallLocalSuccessString, bootEnvInstallLocal3ErrorString},

		CliTest{true, false, []string{"prefs"}, noStdinString, "List and set DigitalRebar Provision operational preferences\n", noErrorString},
		CliTest{false, false, []string{"prefs", "list"}, noStdinString, prefsDefaultListString, noErrorString},

		CliTest{true, true, []string{"prefs", "set"}, noStdinString, noContentString, prefsSetNoArgsErrorString},
		CliTest{true, true, []string{"prefs", "set", "john", "john2", "john3"}, noStdinString, noContentString, prefsSetOddArgsErrorString},
		CliTest{false, true, []string{"prefs", "set", "john"}, noStdinString, noContentString, prefsSetBadJSONErrorString},

		// Set empty hash - should result in no changes
		CliTest{false, false, []string{"prefs", "set", prefsSetEmptyJSONString}, noStdinString, prefsDefaultListString, noErrorString},
		CliTest{false, false, []string{"prefs", "list"}, noStdinString, prefsDefaultListString, noErrorString},

		CliTest{false, false, []string{"prefs", "set", "defaultBootEnv", "local3"}, noStdinString, prefsSetJSONResponseString, noErrorString},
		CliTest{false, false, []string{"prefs", "list"}, noStdinString, prefsChangedListString, noErrorString},

		CliTest{false, true, []string{"prefs", "set", "defaultBootEnv", "illegal"}, noStdinString, noContentString, prefsSetIllegalJSONResponseString},
		CliTest{false, false, []string{"prefs", "list"}, noStdinString, prefsChangedListString, noErrorString},

		CliTest{false, true, []string{"prefs", "set", "knownTokenTimeout", "illegal"}, noStdinString, noContentString, prefsSetBadKnownTokenTimeoutErrorString},
		CliTest{false, true, []string{"prefs", "set", "unknownTokenTimeout", "illegal"}, noStdinString, noContentString, prefsSetBadUnknownTokenTimeoutErrorString},
		CliTest{false, false, []string{"prefs", "set", "knownTokenTimeout", "5000"}, noStdinString, prefsKnownChangedListString, noErrorString},
		CliTest{false, false, []string{"prefs", "set", "unknownTokenTimeout", "7000"}, noStdinString, prefsBothPreDebugChangedListString, noErrorString},
		CliTest{false, false, []string{"prefs", "set", "debugRenderer", "1", "debugDhcp", "2", "debugBootEnv", "1"}, noStdinString, prefsBothChangedListString, noErrorString},
		CliTest{false, false, []string{"prefs", "list"}, noStdinString, prefsBothChangedListString, noErrorString},

		CliTest{false, true, []string{"prefs", "set", "greg", "ignore"}, noStdinString, noContentString, prefsSetInvalidPrefResponseString},
		CliTest{false, false, []string{"prefs", "list"}, noStdinString, prefsBothChangedListString, noErrorString},

		CliTest{false, true, []string{"prefs", "set", "-"}, prefsSetStdinBadJSONString, noContentString, prefsSetBadJSONErrorString},
		CliTest{false, false, []string{"prefs", "list"}, noStdinString, prefsBothChangedListString, noErrorString},
		CliTest{false, false, []string{"prefs", "set", "-"}, prefsSetStdinJSONString, prefsBothChangedListString, noErrorString},
		CliTest{false, false, []string{"prefs", "list"}, noStdinString, prefsBothChangedListString, noErrorString},

		CliTest{false, true, []string{"bootenvs", "destroy", "local3"}, noStdinString, noContentString, prefsFailedToDeleteBootenvErrorString},
		CliTest{false, false, []string{"prefs", "set", "defaultBootEnv", "local", "debugDhcp", "0", "debugRenderer", "0", "debugBootEnv", "0"}, noStdinString, prefsFinalListString, noErrorString},

		CliTest{false, false, []string{"bootenvs", "destroy", "local3"}, noStdinString, "Deleted bootenv local3\n", noErrorString},
		CliTest{false, false, []string{"templates", "destroy", "local3-pxelinux.tmpl"}, noStdinString, "Deleted template local3-pxelinux.tmpl\n", noErrorString},
		CliTest{false, false, []string{"templates", "destroy", "local3-elilo.tmpl"}, noStdinString, "Deleted template local3-elilo.tmpl\n", noErrorString},
		CliTest{false, false, []string{"templates", "destroy", "local3-ipxe.tmpl"}, noStdinString, "Deleted template local3-ipxe.tmpl\n", noErrorString},
	}
	for _, test := range tests {
		testCli(t, test)
	}

	os.RemoveAll("bootenvs")
	os.RemoveAll("templates")
	os.RemoveAll("isos")
	os.RemoveAll("ic")
}
