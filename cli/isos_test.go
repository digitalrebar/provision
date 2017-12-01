package cli

import (
	"os"
	"strconv"
	"strings"
	"testing"
)

var isosDefaultListString string = "[]\n"

var isosUploadNoArgsErrorString = "Error: Wrong number of args: expected 3, got 0\n"
var isosUploadOneArgsErrorString = "Error: Wrong number of args: expected 3, got 1\n"
var isosUploadFourArgsErrorString = "Error: Wrong number of args: expected 3, got 4\n"
var isosUploadMissingIsoErrorString = "Error: Failed to open greg: open greg: no such file or directory\n\n"

var isosUploadSuccessString = `{
  "Path": "greg",
  "Size": *REPLACE_WITH_SIZE*
}
`
var isosUploadCommonSuccessString = `{
  "Path": "greg",
  "Size": *REPLACE_WITH_SIZE*
}
`
var isosGregListString string = `[
  "greg"
]
`
var isosDestroyNoArgsErrorString = "Error: drpcli isos destroy [id] [flags] requires 1 argument\n"
var isosDestroyTwoArgsErrorString = "Error: drpcli isos destroy [id] [flags] requires 1 argument\n"
var isosDestroyGregSuccessString = "Deleted iso greg\n"
var isosDestroyFredErrorString = "Error: DELETE: isos/fred: no such iso\n\n"

func TestIsosCli(t *testing.T) {
	fi, _ := os.Stat("common.go")
	common_size := fi.Size()
	fi, _ = os.Stat("isos.go")
	isos_size := fi.Size()

	isosUploadSuccessString = strings.Replace(isosUploadSuccessString, "*REPLACE_WITH_SIZE*", strconv.FormatInt(isos_size, 10), -1)
	isosUploadCommonSuccessString = strings.Replace(isosUploadCommonSuccessString, "*REPLACE_WITH_SIZE*", strconv.FormatInt(common_size, 10), -1)

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
	cliTest(false, false, "isos", "list").run(t)
	cliTest(false, true, "isos", "destroy").run(t)
	cliTest(false, true, "isos", "destroy", "asdg", "asgs").run(t)
	cliTest(false, false, "isos", "destroy", "greg").run(t)
	cliTest(false, true, "isos", "destroy", "fred").run(t)
	cliTest(false, false, "isos", "list").run(t)
}
