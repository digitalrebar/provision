package cli

import "testing"

func TestIsosCli(t *testing.T) {
	// TODO: Add GetAs
	cliTest(true, false, "isos").run(t)
	cliTest(false, false, "isos", "list").run(t)
	cliTest(true, true, "isos", "upload").run(t)
	cliTest(false, true, "isos", "upload", "asg").run(t)
	cliTest(false, true, "isos", "upload", "asg", "two", "three", "four").run(t)
	cliTest(false, true, "isos", "upload", "greg", "as", "greg").run(t)
	cliTest(false, false, "isos", "upload", "isos.go", "as", "greg").run(t)
	cliTest(false, false, "isos", "upload", "common.go", "as", "greg").run(t)
	cliTest(false, false, "isos", "upload", "isos.go", "as", "greg").run(t)
	cliTest(true, true, "isos", "exists").run(t)
	cliTest(false, true, "isos", "exists", "cow", "ted").run(t)
	cliTest(false, false, "isos", "exists", "greg").run(t)
	cliTest(false, true, "isos", "exists", "greg2").run(t)
	cliTest(false, false, "isos", "list").run(t)
	cliTest(false, true, "isos", "destroy").run(t)
	cliTest(false, true, "isos", "destroy", "asdg", "asgs").run(t)
	cliTest(false, false, "isos", "destroy", "greg").run(t)
	cliTest(false, true, "isos", "destroy", "fred").run(t)
	cliTest(false, false, "isos", "list").run(t)
}
