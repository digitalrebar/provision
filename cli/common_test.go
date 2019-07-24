package cli

import (
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strings"
	"sync"
	"testing"

	"github.com/digitalrebar/provision/api"
	"github.com/digitalrebar/provision/models"
)

var (
	running bool
)

var noErrorString string = ``
var noContentString string = ``
var noStdinString string = ``
var licenseLayer = `---
meta:
  Description: License Data for RackN
  Name: rackn-license
  Overwritable: false
  Source: RackN
  Type: dynamic
  Version: v0.0.1-tip-galthaus-dev-24-3a9ac1523e674e05ee91d1184feb007fba7cbe96
  Writable: false
sections:
  profiles:
    rackn-license:
      Available: false
      Description: RackN License Holder
      Errors: []
      Meta:
        color: blue
        icon: key
        title: RackN Content
      Name: rackn-license
      Params:
        rackn/license: zWbsv/HDxynQsaDQ6SPfb5TKz3GjdgVi1q77sEsnjoxdvjQHr1BIlYxy3hHCqXCZyj1/gw1gtxm7fYol7Ps7CXsKICAiQ29udGFjdCI6ICJVbml0IFRlc3QiLAogICJDb250YWN0RW1haWwiOiAic3VwcG9ydEByYWNrbi5jb20iLAogICJDb250YWN0SWQiOiAiMmE0NGU4MDctZmNjYS00NzJjLTg4MzItZTQ0ZDMzZjVjMTQwIiwKICAiR2VuZXJhdGlvblZlcnNpb24iOiAidjAuMC4xLXRpcC02LTYzNjE2MTQ4YmU5ZTJmZWFlMGFhNDAwNjczNGU3Mzc1ZjgwZWU3ZmYiLAogICJHcmFudG9yIjogIlJhY2tOIEluYyIsCiAgIkdyYW50b3JFbWFpbCI6ICJzdXBwb3J0QHJhY2tuLmNvbSIsCiAgIkxpY2Vuc2VzIjogWwogICAgewogICAgICAiRGF0YSI6IFtdLAogICAgICAiSGFyZEV4cGlyZURhdGUiOiAiMjAyMy0wNi0xNVQyMDo1MDo1Mi42OFoiLAogICAgICAiTG9uZ0xpY2Vuc2UiOiAiUmFja04gSW50ZXJuYWwgVXNlIExpY2Vuc2UiLAogICAgICAiTmFtZSI6ICJyYmFjIiwKICAgICAgIlB1cmNoYXNlRGF0ZSI6ICIyMDE4LTA1LTE1VDIwOjUwOjUyLjY4WiIsCiAgICAgICJTaG9ydExpY2Vuc2UiOiAiUmFja05fSW50ZXJuYWwiLAogICAgICAiU29mdEV4cGlyZURhdGUiOiAiMjAyMy0wNi0wNVQyMDo1MDo1Mi42OFoiLAogICAgICAiU3RhcnREYXRlIjogIjIwMTgtMDUtMTVUMjA6NTA6NTIuNjhaIiwKICAgICAgIlZlcnNpb24iOiAidjIuKiIKICAgIH0sCiAgICB7CiAgICAgICJEYXRhIjogW10sCiAgICAgICJIYXJkRXhwaXJlRGF0ZSI6ICIyMDIzLTA2LTIwVDIxOjIwOjIwLjYwMVoiLAogICAgICAiTG9uZ0xpY2Vuc2UiOiAiUmFja04gSW50ZXJuYWwgVXNlIExpY2Vuc2UiLAogICAgICAiTmFtZSI6ICJzZWN1cmUtcGFyYW1zIiwKICAgICAgIlB1cmNoYXNlRGF0ZSI6ICIyMDE4LTA1LTIwVDIxOjIwOjIwLjYwMVoiLAogICAgICAiU2hvcnRMaWNlbnNlIjogIlJhY2tOX0ludGVybmFsIiwKICAgICAgIlNvZnRFeHBpcmVEYXRlIjogIjIwMjMtMDYtMTBUMjE6MjA6MjAuNjAxWiIsCiAgICAgICJTdGFydERhdGUiOiAiMjAxOC0wNS0yMFQyMToyMDoyMC42MDFaIiwKICAgICAgIlZlcnNpb24iOiAidjIuKiIKICAgIH0KICBdLAogICJPd25lciI6ICJSYWNrTiBUZWFtIiwKICAiT3duZXJFbWFpbCI6ICJzdXBwb3J0QHJhY2tuLmNvbSIsCiAgIk93bmVySWQiOiAicmFja24iCn0K
        rackn/license-object:
          Contact: Unit Test
          ContactEmail: support@rackn.com
          ContactId: 2a44e807-fcca-472c-8832-e44d33f5c140
          GenerationVersion: v0.0.1-tip-galthaus-dev-24-3a9ac1523e674e05ee91d1184feb007fba7cbe96
          Grantor: RackN Inc
          GrantorEmail: support@rackn.com
          Licenses:
          - Data: []
            HardExpireDate: 2023-06-15T20:50:52.68Z
            LongLicense: RackN Internal Use License
            Name: rbac
            PurchaseDate: 2018-05-15T20:50:52.68Z
            ShortLicense: RackN_Internal
            SoftExpireDate: 2023-06-05T20:50:52.68Z
            StartDate: 2018-05-15T20:50:52.68Z
            Version: v2.*
          - Data: []
            HardExpireDate: 2023-06-20T21:20:20.601Z
            LongLicense: RackN Internal Use License
            Name: secure-params
            PurchaseDate: 2018-05-20T21:20:20.601Z
            ShortLicense: RackN_Internal
            SoftExpireDate: 2023-06-10T21:20:20.601Z
            StartDate: 2018-05-20T21:20:20.601Z
            Version: v2.*
          Owner: RackN Team
          OwnerEmail: support@rackn.com
          OwnerId: rackn
      ReadOnly: false
      Validated: false
`

// Runs the args against a server and return stdout and stderr.
func runCliCommand(t *testing.T, args []string, stdin, realOut, realErr string, curl bool) error {
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
	if !curl {
		app := NewApp()
		app.SilenceUsage = false
		app.SetArgs(args)
		err := app.Execute()
		if session != nil {
			session.Close()
			session = nil
		}
		return err
	} else {
		uri := "http://127.0.0.1:10002/" + strings.Join(args, "/")
		resp, err := http.Get(uri)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v", err)
			return err
		}
		defer resp.Body.Close()
		io.Copy(os.Stdout, resp.Body)
		return nil
	}
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
	curl           bool
}

func (c CliTest) run(t *testing.T) {
	t.Helper()
	testCli(t, c)
}

func (c CliTest) get(t *testing.T) {
	c.curl = true
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
			res = append(res, arg)
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
	ret = strings.Replace(ret, "\n", "nl", -1)
	ret = strings.Replace(ret, ":", ".", -1)
	ret = strings.Replace(ret, "|", "pipe", -1)
	ret = strings.Replace(ret, "\"", "dquote", -1)
	ret = strings.Replace(ret, "'", "squote", -1)
	ret = regexp.MustCompile(`\.+`).ReplaceAllString(ret, `.`)
	ret = regexp.MustCompile(`\.?/\.?`).ReplaceAllString(ret, `/`)
	ret = strings.Trim(ret, ".")
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
	hasE := false
	// Add access args
	args := test.args
	if !test.curl {
		for _, a := range test.args {
			if a == "-E" {
				hasE = true
				break
			}
		}
		// Add access args
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
	}

	err = runCliCommand(t, args, test.stdin, realOut, realErr, test.curl)
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
	lsession, apierr := api.UserSession("https://127.0.0.1:10001", "rocketskates", "r0cketsk8ts")
	if apierr != nil {
		t.Fatalf("Error getting session: %v", apierr)
		return
	}
	defer lsession.Close()
	layer, err := lsession.GetContentItem("BackingStore")
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

func TestCorePieces(t *testing.T) {
	cliTest(false, false, "-E", "https://127.0.0.1:10001", "-U", "rocketskates", "-P", "r0cketsk8ts", "version").run(t)
	cliTest(false, false, "gohai", "--help").run(t)
	cliTest(false, true, "-F", "cow", "bootenvs", "list").run(t)
	cliTest(false, false, "-F", "yaml", "bootenvs", "list").run(t)
	cliTest(false, false, "-F", "json", "bootenvs", "list").run(t)
	for _, p := range models.AllPrefixes() {
		switch p {
		case "interfaces", "plugin_providers", "preferences":
			continue
		default:
			cliTest(false, false, p, "indexes").run(t)
		}
	}
	verifyClean(t)
}

func TestMain(m *testing.M) {
	var err error
	actuallyPowerThings = false
	defaultStateLoc = ""

	tmpDir, err = ioutil.TempDir("", "cli-")
	if err != nil {
		log.Printf("Creating temp dir for file root failed: %v", err)
		os.Exit(1)
	}

	ret := m.Run()

	err = os.RemoveAll(tmpDir)
	if err != nil {
		log.Printf("Removing temp dir for file root failed: %v", err)
		os.Exit(1)
	}
	os.Exit(ret)
}
