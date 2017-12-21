package api

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"reflect"
	"syscall"
	"testing"
	"time"

	"github.com/digitalrebar/provision/models"
	"github.com/digitalrebar/provision/server"
	"github.com/ghodss/yaml"
	flags "github.com/jessevdk/go-flags"
)

type crudTest struct {
	name      string
	expectRes interface{}
	expectErr error
	op        func() (interface{}, error)
	clean     func()
}

func (l crudTest) run(t *testing.T) {
	t.Helper()
	t.Logf("Testing %s", l.name)
	session.TraceToken(l.name)
	if l.clean != nil {
		defer l.clean()
	}
	res, err := l.op()
	if l.expectErr == nil {
		if err == nil {
			if equal, delta := diff(res, l.expectRes); !equal {
				t.Errorf("ERROR: Unexpected result:\n%s\n\nDiff:%s",
					pretty(res),
					delta)
			} else {
				t.Logf("Got expected results")
			}
		} else {
			t.Errorf("ERROR: Got unexpected error: %#v", err)
		}
	} else {
		if err == nil {
			t.Errorf("ERROR: Did not get expected error %v", l.expectErr)
			t.Errorf("Got result: %v", pretty(res))
		} else if !reflect.DeepEqual(err, l.expectErr) {
			t.Errorf("ERROR: Expected error %#v", l.expectErr)
			t.Errorf("Got error %#v", err)
		} else {
			t.Logf("Got expected error %v", err)
		}
	}
	session.TraceToken("")
}

func rt(t *testing.T,
	name string,
	res interface{},
	err error,
	op func() (interface{}, error),
	clean func()) {
	t.Helper()
	ct := crudTest{
		name:      name,
		expectRes: res,
		expectErr: err,
		op:        op,
		clean:     clean,
	}
	ct.run(t)
}

var session *Client
var tmpDir string

func testFill(m models.Model) {
	if f, ok := m.(models.Filler); ok {
		f.Fill()
	}
	if v, ok := m.(models.ValidateSetter); ok {
		v.SetValid()
		v.SetAvailable()
	}
}

func mustDecode(ref interface{}, obj string) interface{} {
	if err := DecodeYaml([]byte(obj), ref); err != nil {
		log.Panicf("Failed to decode: %v", err)
	}
	if tgt, ok := ref.(models.Model); ok {
		testFill(tgt)
	}
	return ref
}

func pretty(i interface{}) string {
	if s, k := i.(string); k {
		return s
	}
	buf, err := yaml.Marshal(i)
	if err != nil {
		log.Panicf("Error unmarshalling: %v", err)
	}
	return string(buf)
}

func diff(expected, got interface{}) (bool, string) {
	a, b := pretty(expected), pretty(got)
	f1, err := ioutil.TempFile("", "expected-api")
	if err != nil {
		log.Panicf("Failed to create tempfile1: %v", err)
	}
	defer f1.Close()
	defer os.Remove(f1.Name())
	f2, err := ioutil.TempFile("", "got-api")
	if err != nil {
		log.Panicf("Failed to create tempfile2: %v", err)
	}
	defer f2.Close()
	defer os.Remove(f2.Name())
	if _, err := io.WriteString(f1, a); err != nil {
		log.Panicf("Failed to write tempfile1: %v", err)
	}
	if _, err := io.WriteString(f2, b); err != nil {
		log.Panicf("Failed to write tempfile2: %v", err)
	}
	cmd := exec.Command("diff", "-u", f1.Name(), f2.Name())
	res, err := cmd.CombinedOutput()
	if err == nil {
		return true, string(res)
	}
	if es, ok := err.(*exec.ExitError); ok {
		if ec, ok := es.Sys().(syscall.WaitStatus); ok {
			if ec.ExitStatus() == 1 {
				return false, string(res)
			}
		}
	}
	log.Panicf("diff encountered an error: %v", err)
	return cmd.ProcessState.Success(), string(res)
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

func TestMain(m *testing.M) {
	var err error
	tmpDir, err = ioutil.TempDir("", "api-")
	if err != nil {
		log.Printf("Creating temp dir for file root failed: %v", err)
		os.Exit(1)
	}

	testArgs := []string{
		"--base-root", tmpDir,
		"--tls-key", tmpDir + "/server.key",
		"--tls-cert", tmpDir + "/server.crt",
		"--api-port", "10011",
		"--static-port", "10012",
		"--tftp-port", "10013",
		"--dhcp-port", "10014",
		"--binl-port", "10015",
		"--fake-pinger",
		"--drp-id", "Fred",
		"--backend", "memory:///",
		"--local-content", "directory:../test-data/etc/dr-provision?codec=yaml",
		"--default-content", "file:../test-data/usr/share/dr-provision/default.yaml?codec=yaml",
	}

	err = os.MkdirAll(tmpDir+"/plugins", 0755)
	if err != nil {
		log.Printf("Error creating required directory %s: %v", tmpDir+"/plugins", err)
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

	c_opts := generateArgs(testArgs)
	go server.Server(c_opts)

	count := 0
	for count < 30 {
		session, err = UserSession("https://127.0.0.1:10011", "rocketskates", "r0cketsk8ts")
		if err == nil {
			break
		}
		time.Sleep(1 * time.Second)
		count++
	}
	if session == nil {
		log.Printf("Failed to create UserSession: %v", err)
		os.RemoveAll(tmpDir)
		os.Exit(1)
	}
	if err != nil {
		log.Printf("Server failed to start in time allowed")
		os.RemoveAll(tmpDir)
		os.Exit(1)
	}
	ret := m.Run()
	os.RemoveAll(tmpDir)
	os.Exit(ret)
}
