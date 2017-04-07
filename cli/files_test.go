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

var filesUploadMkdirErrorString = "Error: upload: unable to create directory /greg\n\n"

var filesGregListString string = `[
  "drpcli.amd64.linux",
  "greg",
  "jq"
]
`
var filesDestroyNoArgsErrorString = "Error: drpcli files destroy [id] requires 1 argument\n"
var filesDestroyTwoArgsErrorString = "Error: drpcli files destroy [id] requires 1 argument\n"
var filesDestroyGregSuccessString = "Deleted file greg\n"
var filesDestroyFredErrorString = "Error: delete: unable to delete /fred\n\n"

func TestFilesCli(t *testing.T) {
	fi, _ := os.Stat("common.go")
	common_size := fi.Size()
	fi, _ = os.Stat("files.go")
	files_size := fi.Size()

	filesUploadSuccessString = strings.Replace(filesUploadSuccessString, "*REPLACE_WITH_SIZE*", strconv.FormatInt(files_size, 10), -1)
	filesUploadCommonSuccessString = strings.Replace(filesUploadCommonSuccessString, "*REPLACE_WITH_SIZE*", strconv.FormatInt(common_size, 10), -1)

	// TODO: Add GetAs

	tests := []CliTest{
		CliTest{true, false, []string{"files"}, noStdinString, "Commands to manage files on the provisioner\n", ""},
		CliTest{false, false, []string{"files", "list"}, noStdinString, filesDefaultListString, noErrorString},

		CliTest{true, true, []string{"files", "upload"}, noStdinString, noContentString, filesUploadNoArgsErrorString},
		CliTest{true, true, []string{"files", "upload", "asg"}, noStdinString, noContentString, filesUploadOneArgsErrorString},
		CliTest{true, true, []string{"files", "upload", "asg", "two", "three", "four"}, noStdinString, noContentString, filesUploadFourArgsErrorString},
		CliTest{false, true, []string{"files", "upload", "greg", "as", "greg"}, noStdinString, noContentString, filesUploadMissingFileErrorString},
		CliTest{false, false, []string{"files", "upload", "files.go", "as", "greg"}, noStdinString, filesUploadSuccessString, noErrorString},
		CliTest{false, false, []string{"files", "upload", "common.go", "as", "greg"}, noStdinString, filesUploadCommonSuccessString, noErrorString},
		CliTest{false, false, []string{"files", "upload", "files.go", "as", "greg"}, noStdinString, filesUploadSuccessString, noErrorString},
		CliTest{false, true, []string{"files", "upload", "files.go", "as", "greg/greg"}, noStdinString, noContentString, filesUploadMkdirErrorString},
		CliTest{false, false, []string{"files", "list"}, noStdinString, filesGregListString, noErrorString},

		CliTest{true, true, []string{"files", "destroy"}, noStdinString, noContentString, filesDestroyNoArgsErrorString},
		CliTest{true, true, []string{"files", "destroy", "asdg", "asgs"}, noStdinString, noContentString, filesDestroyTwoArgsErrorString},
		CliTest{false, false, []string{"files", "destroy", "greg"}, noStdinString, filesDestroyGregSuccessString, noErrorString},
		CliTest{false, true, []string{"files", "destroy", "fred"}, noStdinString, noContentString, filesDestroyFredErrorString},
		CliTest{false, false, []string{"files", "list"}, noStdinString, filesDefaultListString, noErrorString},
	}

	for _, test := range tests {
		testCli(t, test)
	}

}
