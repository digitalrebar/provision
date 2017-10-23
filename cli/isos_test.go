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

	tests := []CliTest{
		CliTest{true, false, []string{"isos"}, noStdinString, "Commands to manage isos on the provisioner\n", ""},
		CliTest{false, false, []string{"isos", "list"}, noStdinString, isosDefaultListString, noErrorString},

		CliTest{true, true, []string{"isos", "upload"}, noStdinString, noContentString, isosUploadNoArgsErrorString},
		CliTest{true, true, []string{"isos", "upload", "asg"}, noStdinString, noContentString, isosUploadOneArgsErrorString},
		CliTest{true, true, []string{"isos", "upload", "asg", "two", "three", "four"}, noStdinString, noContentString, isosUploadFourArgsErrorString},
		CliTest{false, true, []string{"isos", "upload", "greg", "as", "greg"}, noStdinString, noContentString, isosUploadMissingIsoErrorString},
		CliTest{false, false, []string{"isos", "upload", "isos.go", "as", "greg"}, noStdinString, isosUploadSuccessString, noErrorString},
		CliTest{false, false, []string{"isos", "upload", "common.go", "as", "greg"}, noStdinString, isosUploadCommonSuccessString, noErrorString},
		CliTest{false, false, []string{"isos", "upload", "isos.go", "as", "greg"}, noStdinString, isosUploadSuccessString, noErrorString},
		CliTest{false, false, []string{"isos", "list"}, noStdinString, isosGregListString, noErrorString},

		CliTest{true, true, []string{"isos", "destroy"}, noStdinString, noContentString, isosDestroyNoArgsErrorString},
		CliTest{true, true, []string{"isos", "destroy", "asdg", "asgs"}, noStdinString, noContentString, isosDestroyTwoArgsErrorString},
		CliTest{false, false, []string{"isos", "destroy", "greg"}, noStdinString, isosDestroyGregSuccessString, noErrorString},
		CliTest{false, true, []string{"isos", "destroy", "fred"}, noStdinString, noContentString, isosDestroyFredErrorString},
		CliTest{false, false, []string{"isos", "list"}, noStdinString, isosDefaultListString, noErrorString},
	}

	for _, test := range tests {
		testCli(t, test)
	}

}
