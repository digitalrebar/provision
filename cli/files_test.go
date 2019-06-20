package cli

import "testing"

func TestFilesCli(t *testing.T) {
	// TODO: Add GetAs
	cliTest(true, false, "files").run(t)
	cliTest(false, false, "files", "list").run(t)
	cliTest(true, true, "files", "upload").run(t)
	cliTest(false, true, "files", "upload", "asg").run(t)
	cliTest(false, true, "files", "upload", "asg", "two", "three", "four").run(t)
	cliTest(false, true, "files", "upload", "greg", "as", "greg").run(t)
	cliTest(false, false, "files", "upload", "files.go", "as", "greg").run(t)
	cliTest(false, false, "files", "upload", "common.go", "as", "greg").run(t)
	cliTest(false, false, "files", "upload", "files.go", "as", "greg").run(t)
	cliTest(true, true, "files", "exists").run(t)
	cliTest(false, true, "files", "exists", "cow", "flka").run(t)
	cliTest(false, false, "files", "exists", "greg").run(t)
	cliTest(false, true, "files", "exists", "greg2").run(t)
	cliTest(false, true, "files", "upload", "files.go", "as", "greg/greg").run(t)
	cliTest(false, false, "files", "list").run(t)
	cliTest(true, true, "files", "destroy").run(t)
	cliTest(true, true, "files", "destroy", "asdg", "asgs").run(t)
	cliTest(false, false, "files", "destroy", "greg").run(t)
	cliTest(false, true, "files", "destroy", "fred").run(t)
	cliTest(false, false, "files", "list").run(t)
}
