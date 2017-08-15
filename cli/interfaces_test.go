package cli

import (
	"encoding/json"
	"testing"

	models "github.com/digitalrebar/provision/genmodels"
)

var interfaceShowNoArgErrorString string = "Error: drpcli interfaces show [id] [flags] requires 1 argument\n"
var interfaceShowTooManyArgErrorString string = "Error: drpcli interfaces show [id] [flags] requires 1 argument\n"
var interfaceShowMissingArgErrorString string = "Error: interface get: not found: john\n\n"

var interfaceExistsNoArgErrorString string = "Error: drpcli interfaces exists [id] [flags] requires 1 argument"
var interfaceExistsTooManyArgErrorString string = "Error: drpcli interfaces exists [id] [flags] requires 1 argument"
var interfaceExistsIgnoreString string = ""
var interfaceExistsMissingJohnString string = "Error: interface get: not found: john\n\n"

func TestInterfaceCli(t *testing.T) {
	// Since this data is dynamic, we will test errors here.
	tests := []CliTest{
		CliTest{true, false, []string{"interfaces"}, noStdinString, "Access CLI commands relating to interfaces\n", ""},

		CliTest{true, true, []string{"interfaces", "show"}, noStdinString, noContentString, interfaceShowNoArgErrorString},
		CliTest{true, true, []string{"interfaces", "show", "john", "john2"}, noStdinString, noContentString, interfaceShowTooManyArgErrorString},
		CliTest{false, true, []string{"interfaces", "show", "john"}, noStdinString, noContentString, interfaceShowMissingArgErrorString},

		CliTest{true, true, []string{"interfaces", "exists"}, noStdinString, noContentString, interfaceExistsNoArgErrorString},
		CliTest{true, true, []string{"interfaces", "exists", "john", "john2"}, noStdinString, noContentString, interfaceExistsTooManyArgErrorString},
		CliTest{false, true, []string{"interfaces", "exists", "john"}, noStdinString, noContentString, interfaceExistsMissingJohnString},
		CliTest{true, true, []string{"interfaces", "exists", "john", "john2"}, noStdinString, noContentString, interfaceExistsTooManyArgErrorString},
	}

	for _, test := range tests {
		testCli(t, test)
	}

	sout, serr, err := runCliCommand(t, []string{"interfaces", "list"}, noStdinString)
	if err != nil {
		t.Errorf("Expected error to be nil for interfaces list: %v\n", err)
	}
	if serr != "" {
		t.Errorf("Expected StdErr to be empty for interfaces list: %s\n", serr)
	}
	var intfs []*models.Interface
	if err := json.Unmarshal([]byte(sout), &intfs); err != nil {
		t.Errorf("Failed to unmarshal sout: %s\n%v\n", sout, err)
	}

	if len(intfs) > 0 {
		name := *intfs[0].Name

		sout, serr, err = runCliCommand(t, []string{"interfaces", "show", name}, noStdinString)
		if err != nil {
			t.Errorf("Expected error to be nil for interfaces show: %v\n", err)
		}
		if serr != "" {
			t.Errorf("Expected StdErr to be empty for interfaces show: %s\n", serr)
		}
		var intf models.Interface
		if err := json.Unmarshal([]byte(sout), &intf); err != nil {
			t.Errorf("Failed to unmarshal sout: %s\n%v\n", sout, err)
		}

		if *intf.Name != name {
			t.Errorf("Expected to get interface name: %s, but got %s\n", name, intf.Name)
		}

		sout, serr, err = runCliCommand(t, []string{"interfaces", "exists", name}, noStdinString)
		if err != nil {
			t.Errorf("Expected error to be nil for interfaces exists: %v\n", err)
		}
		if serr != "" {
			t.Errorf("Expected StdErr to be empty for interfaces exists: %s\n", serr)
		}
		if sout != "" {
			t.Errorf("Expected StdOut to be empty for interfaces exists: %s\n", sout)
		}
	} else {
		t.Errorf("No interfaces found!\n")
	}
}
