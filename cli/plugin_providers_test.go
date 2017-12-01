package cli

import (
	"os"
	"os/exec"
	"testing"
	"time"
)

func TestPluginProviderCli(t *testing.T) {
	srcFolder := tmpDir + "/plugins/incrementer"
	cpCmd := exec.Command("cp", "-rf", srcFolder, "incrementer")
	err := cpCmd.Run()
	if err != nil {
		t.Errorf("Failed to copy incrementer: %v\n", err)
	}

	cliTest(true, false, "plugin_providers").run(t)
	cliTest(true, true, "plugin_providers", "show").run(t)
	cliTest(true, true, "plugin_providers", "show", "john", "john2").run(t)
	cliTest(false, true, "plugin_providers", "show", "john").run(t)
	cliTest(false, false, "plugin_providers", "show", "incrementer").run(t)
	cliTest(true, true, "plugin_providers", "exists").run(t)
	cliTest(true, true, "plugin_providers", "exists", "john", "john2").run(t)
	cliTest(false, true, "plugin_providers", "exists", "john").run(t)
	cliTest(false, false, "plugin_providers", "exists", "incrementer").run(t)
	cliTest(false, false, "plugin_providers", "list").run(t)
	cliTest(true, true, "plugin_providers", "destroy").run(t)
	cliTest(true, true, "plugin_providers", "destroy", "john", "john2").run(t)
	cliTest(false, true, "plugin_providers", "destroy", "john").run(t)
	cliTest(false, false, "params", "show", "incrementer/parameter").run(t)
	cliTest(false, false, "plugin_providers", "destroy", "incrementer").run(t)
	cliTest(false, true, "params", "show", "incrementer/parameter").run(t)
	time.Sleep(3 * time.Second)
	cliTest(false, false, "plugin_providers", "list").run(t)
	cliTest(true, true, "plugin_providers", "upload").run(t)
	cliTest(true, true, "plugin_providers", "upload", "john").run(t)
	cliTest(true, true, "plugin_providers", "upload", "john", "as", "john2", "asdga").run(t)
	cliTest(false, true, "plugin_providers", "upload", "john", "as", "john").run(t)
	cliTest(false, false, "plugin_providers", "upload", "incrementer", "as", "incrementer").run(t)
	time.Sleep(3 * time.Second)
	cliTest(false, false, "plugin_providers", "list").run(t)
	os.Remove("incrementer")
}
