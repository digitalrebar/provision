package cli

import (
	"testing"
)

var prefsDefaultListString = `{
  "defaultBootEnv": "sledgehammer",
  "unknownBootEnv": "ignore"
}
`

var prefsSetNoArgsErrorString string = "Error: prefs set either takes a single argument or a multiple of two, not 0"
var prefsSetOddArgsErrorString string = "Error: prefs set either takes a single argument or a multiple of two, not 3"
var prefsSetBadJSONErrorString string = "Error: Invalid prefs: error unmarshaling JSON: json: cannot unmarshal string into Go value of type map[string]string\n\n\n"

var prefsSetEmptyJSONString string = "{}"
var prefsSetEmptyJSONResponseString string = "{}\n"
var prefsSetJSONResponseString string = `{
  "defaultBootEnv": "ignore"
}
`
var prefsSetIllegalJSONResponseString string = "Error: defaultBootEnv: Bootenv illegal does not exist\n\n"
var prefsSetInvalidPrefResponseString string = "Error: Unknown Preference greg\n\n"

var prefsChangedListString = `{
  "defaultBootEnv": "ignore",
  "unknownBootEnv": "ignore"
}
`

var prefsSetStdinBadJSONString = "fred\n"
var prefsSetStdinBadJSONErrorString = ""
var prefsSetStdinJSONString = "{}\n"

func TestPrefsCli(t *testing.T) {
	tests := []CliTest{
		CliTest{true, false, []string{"prefs"}, noStdinString, "List and set RocketSkates operation preferences\n", noErrorString},
		CliTest{false, false, []string{"prefs", "list"}, noStdinString, prefsDefaultListString, noErrorString},

		CliTest{true, true, []string{"prefs", "set"}, noStdinString, noContentString, prefsSetNoArgsErrorString},
		CliTest{true, true, []string{"prefs", "set", "john", "john2", "john3"}, noStdinString, noContentString, prefsSetOddArgsErrorString},
		CliTest{false, true, []string{"prefs", "set", "john"}, noStdinString, noContentString, prefsSetBadJSONErrorString},

		// Set empty hash - should result in no changes
		CliTest{false, false, []string{"prefs", "set", prefsSetEmptyJSONString}, noStdinString, prefsSetEmptyJSONResponseString, noErrorString},
		CliTest{false, false, []string{"prefs", "list"}, noStdinString, prefsDefaultListString, noErrorString},

		CliTest{false, false, []string{"prefs", "set", "defaultBootEnv", "ignore"}, noStdinString, prefsSetJSONResponseString, noErrorString},
		CliTest{false, false, []string{"prefs", "list"}, noStdinString, prefsChangedListString, noErrorString},

		CliTest{false, true, []string{"prefs", "set", "defaultBootEnv", "illegal"}, noStdinString, noContentString, prefsSetIllegalJSONResponseString},
		CliTest{false, false, []string{"prefs", "list"}, noStdinString, prefsChangedListString, noErrorString},

		CliTest{false, true, []string{"prefs", "set", "greg", "ignore"}, noStdinString, noContentString, prefsSetInvalidPrefResponseString},
		CliTest{false, false, []string{"prefs", "list"}, noStdinString, prefsChangedListString, noErrorString},

		CliTest{false, true, []string{"prefs", "set", "-"}, prefsSetStdinBadJSONString, noContentString, prefsSetBadJSONErrorString},
		CliTest{false, false, []string{"prefs", "list"}, noStdinString, prefsChangedListString, noErrorString},
		CliTest{false, false, []string{"prefs", "set", "-"}, prefsSetStdinJSONString, prefsSetEmptyJSONResponseString, noErrorString},
		CliTest{false, false, []string{"prefs", "list"}, noStdinString, prefsChangedListString, noErrorString},
	}
	for _, test := range tests {
		testCli(t, test)
	}
}
