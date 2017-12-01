package cli

import (
	"os"
	"strconv"
	"strings"
	"testing"
)

var filesDefaultListString string = `[
  "drpcli.amd64.linux",
  "jq"
]
`

var filesUploadNoArgsErrorString = "Error: Wrong number of args: expected 3, got 0\n"
var filesUploadOneArgsErrorString = "Error: Wrong number of args: expected 3, got 1\n"
var filesUploadFourArgsErrorString = "Error: Wrong number of args: expected 3, got 4\n"
var filesUploadMissingFileErrorString = "Error: Failed to open greg: open greg: no such file or directory\n\n"

var filesUploadSuccessString = `{
  "Path": "/greg",
  "Size": *REPLACE_WITH_SIZE*
}
`
var filesUploadCommonSuccessString = `{
  "Path": "/greg",
  "Size": *REPLACE_WITH_SIZE*
}
`

var filesUploadMkdirErrorString = "Error: POST: files/greg/greg: Cannot create directory /greg\n\n"

var filesGregListString string = `[
  "drpcli.amd64.linux",
  "greg",
  "jq"
]
`
var filesDestroyNoArgsErrorString = "Error: drpcli files destroy [id] [flags] requires 1 argument\n"
var filesDestroyTwoArgsErrorString = "Error: drpcli files destroy [id] [flags] requires 1 argument\n"
var filesDestroyGregSuccessString = "Deleted file greg\n"
var filesDestroyFredErrorString = "Error: DELETE: files/fred: Unable to delete\n\n"

func TestFilesCli(t *testing.T) {
	fi, _ := os.Stat("common.go")
	common_size := fi.Size()
	fi, _ = os.Stat("files.go")
	files_size := fi.Size()

	filesUploadSuccessString = strings.Replace(filesUploadSuccessString, "*REPLACE_WITH_SIZE*", strconv.FormatInt(files_size, 10), -1)
	filesUploadCommonSuccessString = strings.Replace(filesUploadCommonSuccessString, "*REPLACE_WITH_SIZE*", strconv.FormatInt(common_size, 10), -1)

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
	cliTest(false, true, "files", "upload", "files.go", "as", "greg/greg").run(t)
	cliTest(false, false, "files", "list").run(t)
	cliTest(true, true, "files", "destroy").run(t)
	cliTest(true, true, "files", "destroy", "asdg", "asgs").run(t)
	cliTest(false, false, "files", "destroy", "greg").run(t)
	cliTest(false, true, "files", "destroy", "fred").run(t)
	cliTest(false, false, "files", "list").run(t)
}
