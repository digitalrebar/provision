package cli

import (
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/digitalrebar/provision/api"
	"github.com/digitalrebar/provision/backend"
	"github.com/digitalrebar/provision/embedded"
	"github.com/digitalrebar/provision/midlayer"
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
func runCliCommand(t *testing.T, args []string, stdin, realOut, realErr string) error {
	t.Helper()
	oldOut := os.Stdout
	oldErr := os.Stderr
	oldIn := os.Stdin
	ri, wi, _ := os.Pipe()
	os.Stdin = ri

	io.WriteString(wi, stdin)
	wi.Close()
	defer func() {
		os.Stdout = oldOut
		os.Stderr = oldErr
		os.Stdin = oldIn
		ri.Close()
	}()
	var err error
	var so, se *os.File
	so, err = os.Create(realOut)
	if err != nil {
		return err
	} else {
		defer so.Close()
	}
	se, err = os.Create(realErr)
	if err != nil {
		return err
	} else {
		defer se.Close()
	}
	os.Stdout = so
	os.Stderr = se

	app := NewApp()
	app.SilenceUsage = false
	app.SetArgs(args)
	return app.Execute()
}

var tnm = map[string]int{}
var tnmMux = &sync.Mutex{}

type CliTest struct {
	dumpUsage      bool
	genError       bool
	args           []string
	stdin          string
	expectedStdOut string
	expectedStdErr string
	trace          string
}

func (c CliTest) run(t *testing.T) {
	t.Helper()
	testCli(t, c)
}
func (c *CliTest) Stdin(s string) *CliTest {
	c.stdin = s
	return c
}

func (c *CliTest) Trace(s string) *CliTest {
	c.trace = s
	return c
}

func (c *CliTest) loc(prefix string) string {
	res := []string{}
	sum := md5.New()
	haveSum := false
	for i, arg := range c.args {
		if len(arg) == 0 {
			continue
		}
		switch arg[0] {
		case '[', '{', '-':
			haveSum = true
			str := strings.Join(c.args[i:], ".")
			sum.Write([]byte(str))
		default:
			res = append(res, strings.Replace(arg, "\n", "nl", -1))
		}
		if haveSum {
			break
		}
	}
	if c.stdin != "" {
		haveSum = true
		sum.Write([]byte(c.stdin))
	}
	if haveSum {
		res = append(res, fmt.Sprintf("%x", sum.Sum(nil)))
	}
	ret := path.Join(prefix, strings.Join(res, "."))
	tnmMux.Lock()
	defer tnmMux.Unlock()
	if idx, ok := tnm[ret]; !ok {
		tnm[ret] = 2
		return ret
	} else {
		tnm[ret] += 1
		return fmt.Sprintf("%s.%d", ret, idx)
	}
}

func diff(a, b string) (string, error) {
	f1, err := os.Open(a)
	if err != nil {
		return "", err
	}
	defer f1.Close()
	f2, err := os.Open(b)
	if err != nil {
		return "", err
	}
	defer f2.Close()
	cmd := exec.Command("diff", "-u", f1.Name(), f2.Name())
	res, err := cmd.CombinedOutput()
	return string(res), err
}

func reTest(t *testing.T, expected, actual string) bool {
	res := []*regexp.Regexp{}
	for _, line := range strings.Split(expected[4:], "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		re, err := regexp.Compile(`^[\s]*(` + line + `)[\s]*$`)
		if err != nil {
			t.Errorf("Regex compile error: %v", err)
		} else {
			res = append(res, re)
		}
	}
	soLines := strings.Split(string(actual), "\n")
	soIdx := 0
	reMatches := []bool{}
	for _, re := range res {
		for {
			matched := false
			if re.MatchString(soLines[soIdx]) {
				reMatches = append(reMatches, true)
				matched = true
			}
			soIdx++
			if soIdx == len(soLines) || matched {
				break
			}
		}
		if soIdx == len(soLines) {
			break
		}
	}
	return len(res) == len(reMatches)
}

func testCli(t *testing.T, test CliTest) {
	t.Helper()
	var err error
	loc := test.loc(t.Name())
	testPath := path.Join("test-data", "output", loc)
	expectOut := path.Join(testPath, "stdout.expect")
	expectErr := path.Join(testPath, "stderr.expect")
	realOut := path.Join(testPath, "stdout.actual")
	realErr := path.Join(testPath, "stderr.actual")
	t.Logf("Testing: %v (stdin: %s)\n", test.args, test.stdin)
	t.Logf("Test path: %s", testPath)
	if err := os.MkdirAll(testPath, 0755); err != nil {
		t.Fatalf("Failed to make test input path %s: %v", testPath, err)
		return
	}
	var sob, seb []byte
	sob, err = ioutil.ReadFile(expectOut)
	if err != nil {
		if test.expectedStdOut == "" {
			if !(test.genError || test.dumpUsage) {
				t.Logf("Missing stdout at %s", expectOut)
			}
		} else {
			t.Logf("Saving to %s", expectOut)
			if err := ioutil.WriteFile(expectOut, []byte(test.expectedStdOut), 0644); err != nil {
				t.Fatalf("Error saving to %s: %v", expectOut, err)
				return
			}
			sob = []byte(test.expectedStdOut)
		}
	} else if test.expectedStdOut != "" {
		t.Logf("Expected stdout overridden by %s", expectOut)
	}
	err = nil
	seb, err = ioutil.ReadFile(expectErr)
	if err != nil {
		if test.expectedStdErr == "" {
			if test.genError || test.dumpUsage {
				t.Logf("Missing stderr at %s", expectErr)
			}
		} else {
			t.Logf("Saving to %s", expectErr)
			if err := ioutil.WriteFile(expectErr, []byte(test.expectedStdErr), 0644); err != nil {
				t.Fatalf("Error saving to %s: %v", expectErr, err)
				return
			}
			seb = []byte(test.expectedStdErr)
		}
	} else if test.expectedStdErr != "" {
		t.Logf("Expected stderr overridden by %s", expectErr)
	}
	test.expectedStdOut, test.expectedStdErr = string(sob), string(seb)
	os.Remove(path.Join(testPath, "untouched"))
	err = nil
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
	}
	if test.trace != "" {
		args = append(args, "--trace", test.trace, "--traceToken", loc)
	}
	if session != nil {
		session.Close()
		session = nil
	}

	err = runCliCommand(t, args, test.stdin, realOut, realErr)
	var so, se string
	if buf, err := ioutil.ReadFile(realOut); err == nil {
		so = string(buf)
	}
	if buf, err := ioutil.ReadFile(realErr); err == nil {
		se = string(buf)
	}

	if test.genError && err == nil {
		t.Errorf("FAIL: Expected Error: but none\n")
	}
	if !test.genError && err != nil {
		t.Errorf("FAIL: Expected No Error: but got: %v\n", err)
	}

	// if we are not dumping usage, expect exact/regexp matches
	// If we are dumping usage and there is an error, expect out to match exact and error to prefix match
	// If we are dumping usage and there is not an error, expect err to match exact and out to prefix match
	patchSO, _ := diff(realOut, expectOut)
	patchSE, _ := diff(realErr, expectErr)
	if !test.dumpUsage {
		if strings.HasPrefix(test.expectedStdOut, "RE:\n") {
			if !reTest(t, test.expectedStdOut, so) {
				t.Errorf("FAIL: Expected StdOut:\n%s", test.expectedStdOut)
				t.Errorf("FAIL: Diff from expected:\n%s", patchSO)
			}
		} else {
			if so != test.expectedStdOut {
				t.Errorf("FAIL: Stdout Diff from expected:\n%s", patchSO)
			}
		}
		if strings.HasPrefix(test.expectedStdErr, "RE:\n") {
			if !reTest(t, test.expectedStdErr, se) {
				t.Errorf("FAIL: Expected StdErr:\n%s", test.expectedStdErr)
				t.Errorf("FAIL: Diff from expected:\n%s", patchSE)
			}
		} else {
			if se != test.expectedStdErr {
				t.Errorf("FAIL: Stderr Diff from expected:\n%s", patchSE)
			}
		}
	} else {
		os.Create(path.Join(testPath, "want-usage"))
		if test.genError {
			if !strings.HasPrefix(se, test.expectedStdErr) {
				t.Errorf("FAIL: Expected StdErr to start with: aa%saa, but got: aa%saa\n", test.expectedStdErr, se)
			}
			if so != test.expectedStdOut {
				t.Errorf("FAIL: Stdout Diff from expected:\n%s", patchSO)
			}
			if !strings.Contains(se, "Usage:") {
				t.Errorf("FAIL: Expected StdErr to have Usage, but didn't: %s", se)
			}
		} else {
			if se != test.expectedStdErr {
				t.Errorf("FAIL: Stderr Diff from expected:\n%s", patchSE)
			}
			if !strings.HasPrefix(so, test.expectedStdOut) {
				t.Errorf("FAIL: Expected StdOut to start with: aa%saa, but got: aa%saa\n", test.expectedStdOut, so)
			}
			if !strings.Contains(so, "Usage:") {
				t.Errorf("FAIL: Expected StdOut to have Usage, but didn't")
			}
		}
	}
	t.Logf("Test finished: %v", test.args)
}

func cliTest(expectUsage, expectErr bool, args ...string) *CliTest {
	return &CliTest{
		dumpUsage: expectUsage,
		genError:  expectErr,
		args:      args,
	}
}

func verifyClean(t *testing.T) {
	t.Helper()
	layer, err := session.GetContentItem("BackingStore")
	if err != nil {
		t.Fatalf("Error getting BackingStore: %v", err)
		return
	}
	for k, v := range layer.Sections {
		switch k {
		case "preferences":
			continue
		case "users", "profiles":
			if len(v) == 1 {
				continue
			}
		}
		if len(v) != 0 {
			t.Errorf("BackingStore layer %s was not cleaned up!", k)
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

func TestCorePieces(t *testing.T) {
	cliTest(false, false, "-E", "https://127.0.0.1:10001", "-U", "rocketskates", "-P", "r0cketsk8ts", "version").run(t)
	cliTest(false, true, "-F", "cow", "bootenvs", "list").run(t)
	cliTest(false, false, "-F", "yaml", "bootenvs", "list").run(t)
	cliTest(false, false, "-F", "json", "bootenvs", "list").run(t)
	verifyClean(t)
}

func TestMain(m *testing.M) {
	var err error

	tmpDir, err = ioutil.TempDir("", "cli-")
	if err != nil {
		log.Printf("Creating temp dir for file root failed: %v", err)
		os.Exit(1)
	}

	testArgs := []string{
		"--base-root", tmpDir,
		"--tls-key", tmpDir + "/server.key",
		"--tls-cert", tmpDir + "/server.crt",
		"--api-port", "10001",
		"--static-port", "10002",
		"--tftp-port", "10003",
		"--dhcp-port", "10004",
		"--binl-port", "10005",
		"--fake-pinger",
		"--drp-id", "Fred",
		"--backend", "memory:///",
		"--local-content", "directory:../test-data/etc/dr-provision?codec=yaml",
		"--default-content", "file:../test-data/usr/share/dr-provision/default.yaml?codec=yaml",
		"--base-token-secret", "token-secret-token-secret-token1",
		"--system-grantor-secret", "system-grantor-secret",
	}

	err = os.MkdirAll(tmpDir+"/plugins", 0755)
	if err != nil {
		log.Printf("Error creating required directory %s: %v", d, err)
		os.Exit(1)
	}

	out, err := exec.Command("go", "generate", "../cmds/incrementer/incrementer.go").CombinedOutput()
	if err != nil {
		log.Printf("Failed to generate incrementer plugin: %v, %s", err, string(out))
		os.Exit(1)
	}

	out, err = exec.Command("go", "build", "-o", tmpDir+"/plugins/incrementer", "../cmds/incrementer/incrementer.go", "../cmds/incrementer/content.go").CombinedOutput()
	if err != nil {
		log.Printf("Failed to build incrementer plugin: %v, %s", err, string(out))
		os.Exit(1)
	}

	embedded.IncludeMeFunction()

	c_opts := generateArgs(testArgs)
	go server.Server(c_opts)
	count := 0
	for count < 30 {
		var apierr error
		session, apierr = api.UserSession("https://127.0.0.1:10001", "rocketskates", "r0cketsk8ts")
		if apierr == nil {
			break
		}
		count++
		time.Sleep(1 * time.Second)
	}
	ret := 1
	if session == nil {
		log.Printf("Server failed to start in time allowed")
	} else {
		log.Printf("Server started after %d seconds", count)
		myToken = session.Token()
		session.Close()
		session = nil
		midlayer.ServeStatic("127.0.0.1:10003", backend.NewFS("test-data", nil), nil, backend.NewPublishers(nil))
		ret = m.Run()
	}

	err = os.RemoveAll(tmpDir)
	if err != nil {
		log.Printf("Removing temp dir for file root failed: %v", err)
		os.Exit(1)
	}
	os.Exit(ret)
}
