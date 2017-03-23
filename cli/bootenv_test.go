package cli

import (
	"testing"
)

func TestBootEnvCli(t *testing.T) {
	createTestServer(t)

	cliArgs := []string{
		"-E", "https://127.0.0.1:10001",
		"bootenvs", "ferd",
	}

	App.SetArgs(cliArgs)
	App.Execute()

}
