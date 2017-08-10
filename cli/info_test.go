package cli

import (
	"encoding/json"
	"runtime"
	"testing"

	"github.com/digitalrebar/provision"
	"github.com/digitalrebar/provision/models"
)

var infoGetTooManyArgsErrorString string = "Error: drpcli info get [flags] requires no arguments"

func TestInfoCli(t *testing.T) {
	// Since this data is dynamic, we will test errors here.
	tests := []CliTest{
		CliTest{true, false, []string{"info"}, noStdinString, "Access CLI commands relating to info\n", ""},
		CliTest{true, true, []string{"info", "get", "john2"}, noStdinString, noContentString, infoGetTooManyArgsErrorString},
	}

	for _, test := range tests {
		testCli(t, test)
	}

	sout, serr, err := runCliCommand(t, []string{"info", "get"}, noStdinString)
	if err != nil {
		t.Errorf("Expected error to be nil for info get: %v\n", err)
	}
	if serr != "" {
		t.Errorf("Expected StdErr to be empty for info get: %s\n", serr)
	}
	var info models.Info
	if err := json.Unmarshal([]byte(sout), &info); err != nil {
		t.Errorf("Failed to unmarshal sout: %s\n%v\n", sout, err)
	}

	if *info.Arch != runtime.GOARCH {
		t.Errorf("Expected matching arch: %s %s\n", info.Arch, runtime.GOARCH)
	}
	if *info.Os != runtime.GOOS {
		t.Errorf("Expected matching os: %s %s\n", info.Os, runtime.GOOS)
	}
	if *info.Version != provision.RS_VERSION {
		t.Errorf("Expected matching os: %s %s\n", info.Version, provision.RS_VERSION)
	}
	if *info.ID != "Fred" {
		t.Errorf("Expected matching id: %s %s\n", info.ID, "Fred")
	}
}
