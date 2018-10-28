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
	cliTest(false, true, "files", "get", "/plugin_providers/incrementer/noFile", "to", "-").run(t)
	cliTest(false, false, "files", "get", "/plugin_providers/incrementer/testFile", "to", "-").run(t)
	os.Remove("incrementer")

	// Test extended here because it is related to the cows.

	cliTest(true, false, "extended").run(t)
	cliTest(true, true, "extended", "show").run(t)
	cliTest(true, true, "extended", "show", "john", "john2").run(t)
	cliTest(false, true, "extended", "show", "john").run(t)
	cliTest(true, true, "extended", "exists", "john", "john2").run(t)
	cliTest(false, true, "extended", "exists", "john").run(t)
	cliTest(false, true, "extended", "-l", "freds", "show", "john").run(t)
	cliTest(false, true, "extended", "-l", "freds", "exists", "john").run(t)
	cliTest(false, true, "extended", "-l", "cows", "show", "john").run(t)
	cliTest(false, true, "extended", "-l", "cows", "exists", "john").run(t)
	cliTest(true, true, "extended", "-l", "cows", "create").run(t)
	cliTest(false, true, "extended", "-l", "freds", "create", "{ \"Type\": \"freds\", \"Id\": \"fred1\" }").run(t)
	cliTest(false, true, "extended", "-l", "cows", "create", "{ \"Type\": \"freds\", \"Id\": \"fred1\" }").run(t)
	cliTest(true, true, "extended", "-l", "cows", "update").run(t)
	cliTest(false, true, "extended", "-l", "cows", "update", "john", "john1").run(t)
	cliTest(false, true, "extended", "-l", "cows", "update", "john").run(t)
	cliTest(true, true, "extended", "-l", "cows", "destroy").run(t)
	cliTest(true, true, "extended", "-l", "cows", "destroy", "john", "john1").run(t)
	cliTest(false, true, "extended", "-l", "cows", "destroy", "john").run(t)

	cliTest(true, true, "extended", "-l", "cows", "params").run(t)
	cliTest(false, true, "extended", "-l", "cows", "params", "john").run(t)

	// yes - this one is weird.  It should fail, but not.
	cliTest(false, false, "extended", "-l", "freds", "indexes").run(t)

	cliTest(false, false, "extended", "-l", "cows", "indexes").run(t)
	cliTest(false, false, "extended", "-l", "typed-cows", "indexes").run(t)

	cliTest(false, false, "extended", "-l", "cows", "create", "{ \"Type\": \"cows\", \"Id\": \"fred1\" }").run(t)
	cliTest(false, false, "extended", "-l", "cows", "show", "fred1").run(t)
	cliTest(false, false, "extended", "-l", "cows", "exists", "fred1").run(t)

	cliTest(false, false, "extended", "-l", "cows", "update", "fred1", "{ \"Tree\": \"neat\" }").run(t)
	cliTest(false, false, "extended", "-l", "cows", "show", "fred1").run(t)
	cliTest(false, false, "extended", "-l", "cows", "params", "fred1").run(t)

	cliTest(false, false, "extended", "-l", "cows", "get", "fred1", "param", "cow-size").run(t)
	cliTest(false, false, "extended", "-l", "cows", "set", "fred1", "param", "cow-size", "to", "big").run(t)
	cliTest(false, false, "extended", "-l", "cows", "get", "fred1", "param", "cow-size").run(t)
	cliTest(false, false, "extended", "-l", "cows", "show", "fred1").run(t)

	cliTest(false, false, "extended", "-l", "cows", "add", "fred1", "param", "cow2-size", "to", "big").run(t)
	cliTest(false, false, "extended", "-l", "cows", "show", "fred1").run(t)
	cliTest(false, false, "extended", "-l", "cows", "remove", "fred1", "param", "cow2-size").run(t)
	cliTest(false, false, "extended", "-l", "cows", "show", "fred1").run(t)

	cliTest(false, false, "extended", "-l", "cows", "actions", "fred1").run(t)
	cliTest(false, false, "extended", "-l", "cows", "meta", "fred1").run(t)

	cliTest(false, false, "extended", "-l", "cows", "destroy", "fred1").run(t)

	cliTest(false, true, "extended", "-l", "typed-cows", "create", "{ \"Type\": \"typed-cows\", \"Id\": \"fred1\" }").run(t)
	cliTest(false, true, "extended", "-l", "typed-cows", "create", "{ \"Type\": \"typed-cows\", \"Id\": \"fred1\", \"Spotted\": \"true\" }").run(t)
	cliTest(false, false, "extended", "-l", "typed-cows", "create", "{ \"Type\": \"typed-cows\", \"Id\": \"fred1\", \"Spotted\": true }").run(t)
	cliTest(false, false, "extended", "-l", "typed-cows", "show", "fred1").run(t)

	cliTest(false, true, "extended", "-l", "typed-cows", "update", "fred1", "{ \"Spotted\": \"greg\" }").run(t)
	cliTest(false, false, "extended", "-l", "typed-cows", "show", "fred1").run(t)
	cliTest(false, false, "extended", "-l", "typed-cows", "update", "fred1", "{ \"Spotted\": false }").run(t)
	cliTest(false, false, "extended", "-l", "typed-cows", "show", "fred1").run(t)

	cliTest(false, false, "extended", "-l", "typed-cows", "destroy", "fred1").run(t)
}
