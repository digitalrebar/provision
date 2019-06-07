package agent

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"reflect"
	"syscall"
	"testing"
	"time"

	"github.com/digitalrebar/logger"
	"github.com/digitalrebar/provision/api"
	"github.com/digitalrebar/provision/backend"
	"github.com/digitalrebar/provision/embedded"
	"github.com/digitalrebar/provision/midlayer"
	"github.com/digitalrebar/provision/models"
	"github.com/digitalrebar/provision/server"
	"github.com/digitalrebar/yaml"
	"github.com/jessevdk/go-flags"
)

var (
	tmpDir              string
	myToken             string
	session             *api.Client
	actuallyPowerThings = false
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
			if equal, delta := apiDiff(res, l.expectRes); !equal {
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
	if err := api.DecodeYaml([]byte(obj), ref); err != nil {
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

func apiDiff(expected, got interface{}) (bool, string) {
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

func fakeServer() error {
	var err error

	os.Setenv("RS_TOKEN_PATH", path.Join(tmpDir, "tokens"))
	os.Setenv("RS_ENDPOINT", "https://127.0.0.1:10001")

	testArgs := []string{
		"--base-root", tmpDir,
		"--tls-key", tmpDir + "/server.key",
		"--tls-cert", tmpDir + "/server.crt",
		"--api-port", "10001",
		"--static-port", "10002",
		"--tftp-port", "10003",
		"--dhcp-port", "10004",
		"--binl-port", "10005",
		"--metrics-port", "10006",
		"--static-ip", "127.0.0.1",
		"--fake-pinger",
		"--no-watcher",
		"--drp-id", "Fred",
		"--backend", "memory:///",
		"--plugin-comm-root", "/tmp",
		"--local-content", "directory:../test-data/etc/dr-provision?codec=yaml",
		"--default-content", "file:../test-data/usr/share/dr-provision/default.yaml?codec=yaml",
		"--base-token-secret", "token-secret-token-secret-token1",
		"--system-grantor-secret", "system-grantor-secret",
	}

	err = os.MkdirAll(tmpDir+"/plugins", 0755)
	if err != nil {
		log.Printf("Error creating required directory %s: %v", tmpDir, err)
		return err
	}

	out, err := exec.Command("go", "generate", "../cmds/incrementer/incrementer.go").CombinedOutput()
	if err != nil {
		log.Printf("Failed to generate incrementer plugin: %v, %s", err, string(out))
		return err
	}

	out, err = exec.Command("go", "build", "-o", tmpDir+"/plugins/incrementer", "../cmds/incrementer/incrementer.go", "../cmds/incrementer/content.go").CombinedOutput()
	if err != nil {
		log.Printf("Failed to build incrementer plugin: %v, %s", err, string(out))
		return err
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
	if session == nil {
		return fmt.Errorf("Server failed to start in time allowed")
	} else {
		log.Printf("Server started after %d seconds", count)
		myToken = session.Token()
		session.Close()
		session = nil
		midlayer.ServeStatic("127.0.0.1:10003",
			backend.NewFS("test-data", nil),
			logger.New(nil).Log(""),
			backend.NewPublishers(nil))
	}
	return nil
}

func TestMain(m *testing.M) {
	var err error
	actuallyPowerThings = false

	tmpDir, err = ioutil.TempDir("", "cli-")
	if err != nil {
		log.Printf("Creating temp dir for file root failed: %v", err)
		os.Exit(1)
	}
	err = fakeServer()
	if err != nil {
		log.Fatalf("Failed with error: %v", err)
	}

	ret := m.Run()

	err = os.RemoveAll(tmpDir)
	if err != nil {
		log.Printf("Removing temp dir for file root failed: %v", err)
		os.Exit(1)
	}
	os.Exit(ret)
}
