package cli

import (
	"os"
	"testing"
)

var prefsDefaultListString = `{
  "defaultBootEnv": "sledgehammer",
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
  "defaultBootEnv": "local",
  "knownTokenTimeout": "3600",
  "unknownBootEnv": "ignore",
  "unknownTokenTimeout": "600"
}
`
var prefsSetIllegalJSONResponseString string = "Error: defaultBootEnv: Bootenv illegal does not exist\n\n"
var prefsSetInvalidPrefResponseString string = "Error: Unknown Preference greg\n\n"

var prefsChangedListString = `{
  "defaultBootEnv": "local",
  "knownTokenTimeout": "3600",
  "unknownBootEnv": "ignore",
  "unknownTokenTimeout": "600"
}
`

var prefsSetStdinBadJSONString = "fred\n"
var prefsSetStdinBadJSONErrorString = ""
var prefsSetStdinJSONString = `{
  "defaultBootEnv": "local",
  "unknownBootEnv": "ignore"
}
`

var prefsSetBadKnownTokenTimeoutErrorString = "Error: Preference knownTokenTimeout: strconv.Atoi: parsing \"illegal\": invalid syntax\n\n"
var prefsSetBadUnknownTokenTimeoutErrorString = "Error: Preference unknownTokenTimeout: strconv.Atoi: parsing \"illegal\": invalid syntax\n\n"

var prefsKnownChangedListString = `{
  "defaultBootEnv": "local",
  "knownTokenTimeout": "5000",
  "unknownBootEnv": "ignore",
  "unknownTokenTimeout": "600"
}
`
var prefsBothChangedListString = `{
  "defaultBootEnv": "local",
  "knownTokenTimeout": "5000",
  "unknownBootEnv": "ignore",
  "unknownTokenTimeout": "7000"
}
`

func TestPrefsCli(t *testing.T) {
	if err := os.MkdirAll("bootenvs", 0755); err != nil {
		t.Errorf("Failed to create bootenvs dir: %v\n", err)
	}
	if err := os.Symlink("../../assets/bootenvs/local.yml", "bootenvs/local.yml"); err != nil {
		t.Errorf("Failed to create link to local.yml: %v\n", err)
	}

	if err := os.MkdirAll("templates", 0755); err != nil {
		t.Errorf("Failed to create templates dir: %v\n", err)
	}
	tmpls := []string{"local-pxelinux.tmpl", "local-elilo.tmpl", "local-ipxe.tmpl"}
	for _, tmpl := range tmpls {
		if err := os.Symlink("../../assets/templates/"+tmpl, "templates/"+tmpl); err != nil {
			t.Errorf("Failed to create link to %s: %v\n", tmpl, err)
		}
	}

	tests := []CliTest{
		CliTest{false, false, []string{"bootenvs", "install", "bootenvs/local.yml"}, noStdinString, bootEnvInstallLocalSuccessString, noErrorString},

		CliTest{true, false, []string{"prefs"}, noStdinString, "List and set DigitalRebar Provision operational preferences\n", noErrorString},
		CliTest{false, false, []string{"prefs", "list"}, noStdinString, prefsDefaultListString, noErrorString},

		CliTest{true, true, []string{"prefs", "set"}, noStdinString, noContentString, prefsSetNoArgsErrorString},
		CliTest{true, true, []string{"prefs", "set", "john", "john2", "john3"}, noStdinString, noContentString, prefsSetOddArgsErrorString},
		CliTest{false, true, []string{"prefs", "set", "john"}, noStdinString, noContentString, prefsSetBadJSONErrorString},

		// Set empty hash - should result in no changes
		CliTest{false, false, []string{"prefs", "set", prefsSetEmptyJSONString}, noStdinString, prefsDefaultListString, noErrorString},
		CliTest{false, false, []string{"prefs", "list"}, noStdinString, prefsDefaultListString, noErrorString},

		CliTest{false, false, []string{"prefs", "set", "defaultBootEnv", "local"}, noStdinString, prefsSetJSONResponseString, noErrorString},
		CliTest{false, false, []string{"prefs", "list"}, noStdinString, prefsChangedListString, noErrorString},

		CliTest{false, true, []string{"prefs", "set", "defaultBootEnv", "illegal"}, noStdinString, noContentString, prefsSetIllegalJSONResponseString},
		CliTest{false, false, []string{"prefs", "list"}, noStdinString, prefsChangedListString, noErrorString},

		CliTest{false, true, []string{"prefs", "set", "knownTokenTimeout", "illegal"}, noStdinString, noContentString, prefsSetBadKnownTokenTimeoutErrorString},
		CliTest{false, true, []string{"prefs", "set", "unknownTokenTimeout", "illegal"}, noStdinString, noContentString, prefsSetBadUnknownTokenTimeoutErrorString},
		CliTest{false, false, []string{"prefs", "set", "knownTokenTimeout", "5000"}, noStdinString, prefsKnownChangedListString, noErrorString},
		CliTest{false, false, []string{"prefs", "set", "unknownTokenTimeout", "7000"}, noStdinString, prefsBothChangedListString, noErrorString},
		CliTest{false, false, []string{"prefs", "list"}, noStdinString, prefsBothChangedListString, noErrorString},

		CliTest{false, true, []string{"prefs", "set", "greg", "ignore"}, noStdinString, noContentString, prefsSetInvalidPrefResponseString},
		CliTest{false, false, []string{"prefs", "list"}, noStdinString, prefsBothChangedListString, noErrorString},

		CliTest{false, true, []string{"prefs", "set", "-"}, prefsSetStdinBadJSONString, noContentString, prefsSetBadJSONErrorString},
		CliTest{false, false, []string{"prefs", "list"}, noStdinString, prefsBothChangedListString, noErrorString},
		CliTest{false, false, []string{"prefs", "set", "-"}, prefsSetStdinJSONString, prefsBothChangedListString, noErrorString},
		CliTest{false, false, []string{"prefs", "list"}, noStdinString, prefsBothChangedListString, noErrorString},

		// Clean-up - can't happen now.
	}
	for _, test := range tests {
		testCli(t, test)
	}

	os.RemoveAll("bootenvs")
	os.RemoveAll("templates")
	os.RemoveAll("isos")
	os.RemoveAll("ic")
}
