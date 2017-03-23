package cli

import (
	"testing"
)

func TestMachineCli(t *testing.T) {
	createTestServer(t)

	tests := []CliTest{
		CliTest{[]string{"-E", "https://127.0.0.1:10001", "machines", "list"}, "[]\n", ""},
		CliTest{[]string{"-E", "https://127.0.0.1:10001", "machines", "show"}, "", "rscli machines show [id] requires 1 argument\n"},
		CliTest{[]string{"-E", "https://127.0.0.1:10001", "machines", "show", "john"}, "", "Failed to fetch machine: john\n[GET /machines/{uuid}][404] getMachineNotFound  &{Key:john Messages:[machines GET: john: Not Found] Model:machines Type:API_ERROR}\n"},
	}

	for _, test := range tests {
		testCli(t, test)
	}

}
