package api

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

	"github.com/digitalrebar/provision/v4/models"
	"github.com/digitalrebar/provision/v4/test"
	"github.com/ghodss/yaml"
)

var (
	session *Client
	tmpDir  string
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

func TestMain(m *testing.M) {
	var err error
	tmpDir, err = ioutil.TempDir("", "api-")
	if err != nil {
		log.Printf("Creating temp dir for file root failed: %v", err)
		os.Exit(1)
	}
	defer os.RemoveAll(tmpDir)
	if err := test.StartServer(tmpDir, 10011); err != nil {
		log.Printf("Error starting dr-provision: %v", err)
		os.RemoveAll(tmpDir)
		os.Exit(1)
	}
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
		err = fmt.Errorf("Failed to create UserSession: %v", err)
	}
	if err != nil {
		log.Printf("Error starting test run: %v", err)
		test.StopServer()
		os.RemoveAll(tmpDir)
		os.Exit(1)
	}
	if err := session.MakeProxy(path.Join(tmpDir, ".socket")); err != nil {
		log.Printf("failed to create local proxy socket: %v", err)
		test.StopServer()
		os.RemoveAll(tmpDir)
		os.Exit(1)
	}
	ret := m.Run()
	if err := test.StopServer(); err != nil {
		log.Printf("Error stopping dr-provision: %v", err)
	}
	os.RemoveAll(tmpDir)
	os.Exit(ret)
}
