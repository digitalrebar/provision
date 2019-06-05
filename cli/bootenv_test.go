package cli

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/digitalrebar/tftp"
)

func TestBootEnvCli(t *testing.T) {
	var (
		bootEnvCreateBadJSONString        = "{asdgasdg}"
		bootEnvCreateInputString   string = `{
  "name": "john"
}
`
		bootEnvCreateFredInputString string = `fred`
		bootEnvUpdateBadJSONString          = "asdgasdg"

		bootEnvUpdateInputString string = `{
  "Kernel": "lpxelinux.0",
  "OS": {"Name":"johann"}
}
`
	)

	cliTest(true, false, "bootenvs").run(t)
	cliTest(false, false, "bootenvs", "list").run(t)
	cliTest(false, false, "bootenvs", "list", "--limit=0").run(t)
	cliTest(false, false, "bootenvs", "list", "--limit=10", "--offset=0").run(t)
	cliTest(false, false, "bootenvs", "list", "--limit=10", "--offset=10").run(t)
	cliTest(false, true, "bootenvs", "list", "--limit=-10", "--offset=0").run(t)
	cliTest(false, true, "bootenvs", "list", "--limit=10", "--offset=-10").run(t)
	cliTest(false, false, "bootenvs", "list", "--limit=-1", "--offset=-1").run(t)
	cliTest(false, false, "bootenvs", "list", "Name=fred").run(t)
	cliTest(false, false, "bootenvs", "list", "Name=ignore").run(t)
	cliTest(false, false, "bootenvs", "list", "OnlyUnknown=true").run(t)
	cliTest(false, false, "bootenvs", "list", "OnlyUnknown=false").run(t)
	cliTest(false, false, "bootenvs", "list", "Available=true").run(t)
	cliTest(false, false, "bootenvs", "list", "Available=false").run(t)
	cliTest(false, true, "bootenvs", "list", "Available=fred").run(t)
	cliTest(false, false, "bootenvs", "list", "Valid=true").run(t)
	cliTest(false, false, "bootenvs", "list", "Valid=false").run(t)
	cliTest(false, true, "bootenvs", "list", "Valid=fred").run(t)
	cliTest(false, false, "bootenvs", "list", "ReadOnly=true").run(t)
	cliTest(false, false, "bootenvs", "list", "ReadOnly=false").run(t)
	cliTest(false, true, "bootenvs", "list", "ReadOnly=fred").run(t)
	cliTest(false, true, "bootenvs", "list", "sort=OnlyUnknown").run(t)
	cliTest(false, false, "bootenvs", "list", "sort=Name").run(t)
	cliTest(false, false, "bootenvs", "list", "sort=Name", "reverse=true").run(t)

	cliTest(true, true, "bootenvs", "show").run(t)
	cliTest(true, true, "bootenvs", "show", "john", "john2").run(t)
	cliTest(false, true, "bootenvs", "show", "john").run(t)
	cliTest(false, false, "bootenvs", "show", "ignore").run(t)

	cliTest(true, true, "bootenvs", "exists").run(t)
	cliTest(true, true, "bootenvs", "exists", "john", "john2").run(t)
	cliTest(false, false, "bootenvs", "exists", "ignore").run(t)
	cliTest(false, true, "bootenvs", "exists", "john").run(t)

	cliTest(true, true, "bootenvs", "uploadiso").run(t)
	cliTest(true, true, "bootenvs", "uploadiso", "john", "john2").run(t)
	cliTest(false, true, "bootenvs", "uploadiso", "ignore").run(t)
	cliTest(false, true, "bootenvs", "uploadiso", "john").run(t)

	cliTest(true, true, "bootenvs", "create").run(t)
	cliTest(true, true, "bootenvs", "create", "john", "john2").run(t)
	cliTest(false, true, "bootenvs", "create", bootEnvCreateBadJSONString).run(t)
	cliTest(false, false, "bootenvs", "create", bootEnvCreateInputString).run(t)
	cliTest(false, false, "bootenvs", "create", bootEnvCreateFredInputString).run(t)
	cliTest(false, false, "bootenvs", "destroy", bootEnvCreateFredInputString).run(t)
	cliTest(false, true, "bootenvs", "create", bootEnvCreateInputString).run(t)
	cliTest(false, false, "bootenvs", "list").run(t)

	cliTest(true, true, "bootenvs", "update").run(t)
	cliTest(true, true, "bootenvs", "update", "john", "john2", "john3").run(t)
	cliTest(false, true, "bootenvs", "update", "john", bootEnvUpdateBadJSONString).run(t)
	cliTest(false, false, "bootenvs", "update", "john", bootEnvUpdateInputString).run(t)
	cliTest(false, true, "bootenvs", "update", "john2", bootEnvUpdateInputString).run(t)
	cliTest(false, false, "bootenvs", "show", "john").run(t)
	cliTest(false, false, "bootenvs", "show", "john").run(t)

	cliTest(false, true, "bootenvs", "destroy").run(t)
	cliTest(false, true, "bootenvs", "destroy", "john", "june").run(t)
	cliTest(false, false, "bootenvs", "destroy", "john").run(t)
	cliTest(false, true, "bootenvs", "destroy", "john").run(t)
	cliTest(false, false, "bootenvs", "list").run(t)

	cliTest(false, false, "bootenvs", "create", "-").Stdin(bootEnvCreateInputString + "\n").run(t)
	cliTest(false, false, "bootenvs", "list").run(t)
	cliTest(false, false, "bootenvs", "update", "john", "-").Stdin(bootEnvUpdateInputString + "\n").run(t)
	cliTest(false, false, "bootenvs", "show", "john").run(t)
	cliTest(false, false, "bootenvs", "destroy", "john").run(t)
	cliTest(false, false, "bootenvs", "list").run(t)
	verifyClean(t)
}

func TestBootenvStageHandling(t *testing.T) {
	cliTest(false, false, "stages", "create", "-").Stdin(`---
Name: fred
BootEnv: fred`).run(t)
	cliTest(false, false, "bootenvs", "create", "fred").run(t)
	cliTest(false, false, "stages", "show", "fred").run(t)
	cliTest(false, false, "bootenvs", "show", "fred").run(t)
	cliTest(false, false, "stages", "destroy", "fred").run(t)
	cliTest(false, false, "bootenvs", "destroy", "fred").run(t)
	verifyClean(t)
}

func TestBootEnvLookaside(t *testing.T) {
	testFile := "sledgehammer/708de8b878e3818b1c1bb598a56de968939f9d4b/vmlinuz0"
	cliTest(false, false, "profiles", "add", "global", "param", "package-repositories", "to", "-").Stdin(`
- tag: "sledgehammer-708de8b878e3818b1c1bb598a56de968939f9d4b"
  os:
    - "sledgehammer/708de8b878e3818b1c1bb598a56de968939f9d4b"
  installSource: true
  arch: amd64
  url: "http://127.0.0.1:10003/hammertime"
`).run(t)
	cliTest(false, false, "bootenvs", "install", "test-data/no-phredhammer.yml").run(t)
	cliTest(false, false, "bootenvs", "install", "test-data/phredhammer.yml").run(t)
	time.Sleep(5 * time.Second)
	expected := "GREG-vmlinuz0\n"
	cliTest(false, false, "bootenvs", "show", "phredhammer").run(t)
	testUrl := "http://127.0.0.1:10002/" + testFile
	resp, err := http.Get(testUrl)
	if err != nil {
		t.Errorf("http: Error %v looking for redirected phredhammer files", err)
	} else if resp.StatusCode != 200 {
		t.Errorf("http: Invalid status code looking for phredhammer files: %d", resp.StatusCode)
	} else if resp.ContentLength != 14 {
		t.Errorf("http: Expected size 14, not %d", resp.ContentLength)
	} else {
		body, _ := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if string(body) != expected {
			t.Errorf("http: Wanted body\n`%s`\nnot\n`%s`\n", expected, string(body))
		} else {
			t.Logf(`http: Lookaside from 
http://127.0.0.1:10002/sledgehammer/708de8b878e3818b1c1bb598a56de968939f9d4b/
to
http://127.0.0.1:10003/hammertime
worked`)
		}
	}
	c, err := tftp.NewClient("127.0.0.1:10003")
	if err == nil {
		c.RequestTSize(true)
		if src, err := c.Receive(testFile, ""); err != nil {
			t.Errorf("tftp: Error fetching: %v", err)
		} else if n, ok := src.(tftp.IncomingTransfer); !ok {
			t.Errorf("tftp: Expected to get a sized answer, but did not")
		} else if sz, _ := n.Size(); sz != 14 {
			t.Errorf("tftp: Expected size 14, got %d", sz)
		} else {
			buf := &bytes.Buffer{}
			src.WriteTo(buf)
			body := buf.String()
			if body != expected {
				t.Errorf("tftp: Wanted body\n`%s`\nnot\n`%s`\n", expected, body)
			} else {
				t.Logf(`tftp: Lookaside from 
sledgehammer/708de8b878e3818b1c1bb598a56de968939f9d4b/vmlinuz0
to
http://127.0.0.1:10003/hammertime/vmlinuz0
worked`)
			}
		}
	}
	cliTest(false, false, "profiles", "remove", "global", "param", "package-repositories").run(t)
	cliTest(false, false, "bootenvs", "destroy", "phredhammer").run(t)
	cliTest(false, false, "bootenvs", "destroy", "no-phredhammer").run(t)
	cliTest(false, false, "templates", "destroy", "local3-pxelinux.tmpl").run(t)
	cliTest(false, false, "templates", "destroy", "local3-elilo.tmpl").run(t)
	cliTest(false, false, "templates", "destroy", "local3-ipxe.tmpl").run(t)
	verifyClean(t)
}
