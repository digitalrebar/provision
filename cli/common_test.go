package cli

import (
	"bytes"
	"crypto/tls"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/jessevdk/go-flags"
	"github.com/rackn/rocket-skates/server"
)

var (
	tmpDir  string
	running bool
)

// Runs the args against a server and return stdout and stderr.
func runCliCommand(t *testing.T, args []string) (string, string) {
	createTestServer(t)

	App.SetArgs(args)

	var b bytes.Buffer
	App.SetOutput(&b)

	var c bytes.Buffer
	log.SetOutput(&c)

	App.Execute()

	return c.String(), b.String()
}

type CliTest struct {
	args           []string
	expectedStdOut string
	expectedStdErr string
}

func testCli(t *testing.T, test CliTest) {
	t.Logf("Testing: %v\n", test.args)

	so, se := runCliCommand(t, test.args)

	if so != test.expectedStdOut {
		t.Errorf("Expected StdOut: %s, but got: %s\n", test.expectedStdOut, so)
	}
	if se != test.expectedStdErr {
		t.Errorf("Expected StdErr: %s, but got: %s\n", test.expectedStdErr, se)
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

func createTestServer(t *testing.T) {
	if running {
		return
	}
	running = true

	var err error
	tmpDir, err = ioutil.TempDir("", "cli-")
	if err != nil {
		log.Printf("Creating temp dir for file root failed: %v", err)
		os.Exit(1)
	}
	defer os.RemoveAll(tmpDir)

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
	_, apierr := client.Get("https://127.0.0.1:10001/api/v3/subnets")
	count := 0
	for apierr != nil && count < 15 {
		t.Logf("Failed to get file: %v", apierr)
		time.Sleep(1 * time.Second)
		count++
		_, apierr = client.Get("https://127.0.0.1:10001/api/v3/subnets")
	}
	if count == 15 {
		t.Errorf("Server failed to start in time allowed")
	}
}
