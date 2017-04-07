package cli

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/digitalrebar/provision/server"
	"github.com/jessevdk/go-flags"
)

var (
	tmpDir  string
	running bool
	myToken string
)

var noErrorString string = ``
var noContentString string = ``
var noStdinString string = ``

// Runs the args against a server and return stdout and stderr.
func runCliCommand(t *testing.T, args []string, stdin string) (string, string, error) {
	oldOut := os.Stdout
	oldErr := os.Stderr
	oldIn := os.Stdin

	ro, wo, _ := os.Pipe()
	os.Stdout = wo
	re, we, _ := os.Pipe()
	os.Stderr = we
	ri, wi, _ := os.Pipe()
	os.Stdin = ri

	io.WriteString(wi, stdin)
	wi.Close()

	// Can cause stdin read error here by: ri.Close(), but is it worth it for two lines of coverage.

	outC := make(chan string)
	errC := make(chan string)

	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, ro)
		outC <- buf.String()
	}()
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, re)
		errC <- buf.String()
	}()

	dumpUsage = true

	App.SetArgs(args)
	err := App.Execute()

	wo.Close()
	we.Close()

	os.Stdout = oldOut
	os.Stderr = oldErr
	os.Stdin = oldIn

	outS := <-outC
	errS := <-errC

	ri.Close()
	ro.Close()
	re.Close()

	return outS, errS, err
}

type CliTest struct {
	dumpUsage      bool
	genError       bool
	args           []string
	stdin          string
	expectedStdOut string
	expectedStdErr string
}

func testCli(t *testing.T, test CliTest) {
	t.Logf("Testing: %v (stdin: %s)\n", test.args, test.stdin)

	hasE := false
	// Add access args
	for _, a := range test.args {
		if a == "-E" {
			hasE = true
			break
		}
	}

	// Add access args
	args := test.args
	if !hasE {
		args = []string{"-E", "https://127.0.0.1:10001", "-T", myToken}
		for _, a := range test.args {
			args = append(args, a)
		}
	} else {
		session = nil
	}

	so, se, err := runCliCommand(t, args, test.stdin)

	if test.genError && err == nil {
		t.Errorf("Expected Error: but none\n")
	}
	if !test.genError && err != nil {
		t.Errorf("Expected No Error: but got: %v\n", err)
	}

	// if we are not dumping usage, expect exact/regexp matches
	// If we are dumping usage and there is an error, expect out to match exact and error to prefix match
	// If we are dumping usage and there is not an error, expect err to match exact and out to prefix match
	if !test.dumpUsage {
		if strings.HasPrefix(test.expectedStdOut, "RE:\n") {
			if matched, err := regexp.MatchString(test.expectedStdOut[4:], so); err != nil || !matched {
				if err != nil {
					t.Errorf("Expected StdOut: regexp fail: %v\n", err)
				}
				t.Errorf("Expected StdOut: aa%saa, but got: aa%saa\n", test.expectedStdOut[4:], so)
			}
		} else {
			if so != test.expectedStdOut {
				t.Errorf("Expected StdOut: aa%saa, but got: aa%saa\n", test.expectedStdOut, so)
			}
		}
		if se != test.expectedStdErr {
			t.Errorf("Expected StdErr: aa%saa, but got: aa%saa\n", test.expectedStdErr, se)
		}
	} else {
		if test.genError {
			if !strings.HasPrefix(se, test.expectedStdErr) {
				t.Errorf("Expected StdErr to start with: aa%saa, but got: aa%saa\n", test.expectedStdErr, se)
			}
			if so != test.expectedStdOut {
				t.Errorf("Expected StdOut: aa%saa, but got: aa%saa\n", test.expectedStdOut, so)
			}
			if !strings.Contains(se, "Usage:") {
				t.Errorf("Expected StdErr to have Usage, but didn't")
			}
		} else {
			if se != test.expectedStdErr {
				t.Errorf("Expected StdErr: aa%saa, but got: aa%saa\n", test.expectedStdErr, se)
			}
			if !strings.HasPrefix(so, test.expectedStdOut) {
				t.Errorf("Expected StdOut to start with: aa%saa, but got: aa%saa\n", test.expectedStdOut, so)
			}
			if !strings.Contains(so, "Usage:") {
				t.Errorf("Expected StdOut to have Usage, but didn't")
			}
		}
	}
}

func generateArgs(args []string) *server.ProgOpts {
	var c_opts server.ProgOpts

	parser := flags.NewParser(&c_opts, flags.Default)
	if _, err := parser.ParseArgs(args); err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	}

	return &c_opts
}

var yamlTestString = `- Available: true
  BootParams: ""
  Description: The boot environment you should use to have unknown machines boot off
    their local hard drive
  Errors: null
  Initrds: null
  Kernel: ""
  Name: ignore
  OS:
    Name: ignore
  OnlyUnknown: true
  OptionalParams: null
  RequiredParams: null
  Templates:
  - Contents: |
      DEFAULT local
      PROMPT 0
      TIMEOUT 10
      LABEL local
      localboot 0
    Name: pxelinux
    Path: pxelinux.cfg/default
  - Contents: exit
    Name: elilo
    Path: elilo.conf
  - Contents: |
      #!ipxe
      chain tftp://{{.ProvisionerAddress}}/${netX/ip}.ipxe || exit
    Name: ipxe
    Path: default.ipxe

`

var jsonTestString = `[
  {
    "Available": true,
    "BootParams": "",
    "Description": "The boot environment you should use to have unknown machines boot off their local hard drive",
    "Errors": null,
    "Initrds": null,
    "Kernel": "",
    "Name": "ignore",
    "OS": {
      "Name": "ignore"
    },
    "OnlyUnknown": true,
    "OptionalParams": null,
    "RequiredParams": null,
    "Templates": [
      {
        "Contents": "DEFAULT local\nPROMPT 0\nTIMEOUT 10\nLABEL local\nlocalboot 0\n",
        "Name": "pxelinux",
        "Path": "pxelinux.cfg/default"
      },
      {
        "Contents": "exit",
        "Name": "elilo",
        "Path": "elilo.conf"
      },
      {
        "Contents": "#!ipxe\nchain tftp://{{.ProvisionerAddress}}/${netX/ip}.ipxe || exit\n",
        "Name": "ipxe",
        "Path": "default.ipxe"
      }
    ]
  }
]
`

func TestCorePieces(t *testing.T) {
	tests := []CliTest{
		CliTest{false, true, []string{"-E", "khttps://1.1.1.2:325", "bootenvs", "list"}, noStdinString, noContentString, "Error: Error listing bootenvs: Get khttps://1.1.1.2:325/api/v3/bootenvs: unsupported protocol scheme \"khttps\"\n\n"},
		CliTest{false, false, []string{"-E", "https://127.0.0.1:10001", "-U", "rocketskates", "-P", "r0cketsk8ts", "version"}, noStdinString, "Version: " + version + "\n", noErrorString},
		CliTest{false, true, []string{"-F", "cow", "bootenvs", "list"}, noStdinString, noContentString, "Error: Unknown pretty format cow\n\n"},
		CliTest{false, false, []string{"-F", "yaml", "bootenvs", "list"}, noStdinString, yamlTestString, noErrorString},
		CliTest{false, false, []string{"-F", "json", "bootenvs", "list"}, noStdinString, jsonTestString, noErrorString},
	}
	for _, test := range tests {
		testCli(t, test)
	}
}

func TestMain(m *testing.M) {
	var err error
	tmpDir, err = ioutil.TempDir("", "cli-")
	if err != nil {
		log.Printf("Creating temp dir for file root failed: %v", err)
		os.Exit(1)
	}

	testArgs := []string{
		"--data-root", tmpDir + "/digitalrebar",
		"--file-root", tmpDir + "/tftpboot",
		"--tls-key", tmpDir + "/server.key",
		"--tls-cert", tmpDir + "/server.crt",
		"--api-port", "10001",
		"--static-port", "10002",
		"--tftp-port", "10003",
		"--disable-dhcp",
	}

	c_opts := generateArgs(testArgs)
	go server.Server(c_opts)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	req, _ := http.NewRequest("GET", "https://127.0.0.1:10001/api/v3/subnets", nil)
	req.SetBasicAuth("rocketskates", "r0cketsk8ts")
	_, apierr := client.Do(req)
	count := 0
	for apierr != nil && count < 30 {
		time.Sleep(1 * time.Second)
		count++
		req, _ = http.NewRequest("GET", "https://127.0.0.1:10001/api/v3/subnets", nil)
		req.SetBasicAuth("rocketskates", "r0cketsk8ts")
		_, apierr = client.Do(req)
	}
	ret := 1
	if count == 30 {
		log.Printf("Server failed to start in time allowed")
	} else {
		req, _ = http.NewRequest("GET", "https://127.0.0.1:10001/api/v3/users/rocketskates/token", nil)
		req.SetBasicAuth("rocketskates", "r0cketsk8ts")
		resp, _ := client.Do(req)
		var token map[string]string
		buf, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal(buf, &token)
		myToken, _ = token["Token"]

		ret = m.Run()
	}

	err = os.RemoveAll(tmpDir)
	if err != nil {
		log.Printf("Creating temp dir for file root failed: %v", err)
		os.Exit(1)
	}
	os.Exit(ret)
}
