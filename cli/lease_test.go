package cli

import (
	"testing"
)

var leaseDefaultListString string = "[]\n"

func TestLeaseCli(t *testing.T) {
	tests := []CliTest{
		CliTest{true, false, []string{"leases"}, noStdinString, "Access CLI commands relating to leases\n", ""},
		CliTest{false, false, []string{"leases", "list"}, noStdinString, leaseDefaultListString, noErrorString},
	}

	for _, test := range tests {
		testCli(t, test)
	}

}
