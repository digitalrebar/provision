package cli

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"testing"
	"time"

	"github.com/pin/tftp"
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

	cliTest(true, true, "bootenvs", "install").run(t)
	cliTest(true, true, "bootenvs", "install", "john", "john", "john2").run(t)
	cliTest(false, true, "bootenvs", "install", "fredhammer").run(t)

	if f, err := os.Create("bootenvs"); err != nil {
		t.Errorf("FAIL: Failed to create bootenvs file: %v\n", err)
	} else {
		f.Close()
	}

	cliTest(false, true, "bootenvs", "install", "bootenvs/fredhammer.yml").run(t)
	os.RemoveAll("bootenvs")
	if err := os.MkdirAll("bootenvs", 0755); err != nil {
		t.Errorf("FAIL: Failed to create bootenvs dir: %v\n", err)
	}

	cliTest(false, true, "bootenvs", "install", "bootenvs/fredhammer.yml").run(t)
	if err := ioutil.WriteFile("bootenvs/fredhammer.yml", []byte("TEST"), 0644); err != nil {
		t.Errorf("FAIL: Failed to create bootenvs file: %v\n", err)
	}

	cliTest(false, true, "bootenvs", "install", "bootenvs/fredhammer.yml").run(t)

	os.RemoveAll("bootenvs/fredhammer.yml")
	if err := os.MkdirAll("bootenvs", 0755); err != nil {
		t.Errorf("FAIL: Failed to create bootenvs dir: %v\n", err)
	}
	if err := os.Symlink("../test-data/fredhammer.yml", "bootenvs/fredhammer.yml"); err != nil {
		t.Errorf("FAIL: Failed to create link to fredhammer.yml: %v\n", err)
	}
	if err := os.Symlink("../test-data/local3.yml", "bootenvs/local3.yml"); err != nil {
		t.Errorf("FAIL: Failed to create link to local3.yml: %v\n", err)
	}

	cliTest(false, false, "bootenvs", "install", "--skip-download", "bootenvs/fredhammer.yml").run(t)
	cliTest(false, false, "bootenvs", "destroy", "fredhammer").run(t)

	installSkipDownloadIsos = false

	cliTest(false, false, "bootenvs", "install", "bootenvs/fredhammer.yml").run(t)
	cliTest(false, true, "bootenvs", "install", "bootenvs/local3.yml").run(t)

	if err := os.MkdirAll("templates", 0755); err != nil {
		t.Errorf("FAIL: Failed to create templates dir: %v\n", err)
	}
	tmpls := []string{"local3-pxelinux.tmpl", "local3-elilo.tmpl", "local3-ipxe.tmpl"}
	for _, tmpl := range tmpls {
		if err := os.Symlink("../test-data/"+tmpl, "templates/"+tmpl); err != nil {
			t.Errorf("FAIL: Failed to create link to %s: %v\n", tmpl, err)
		}
	}

	cliTest(false, false, "bootenvs", "install", "bootenvs/local3.yml", "ic").run(t)
	cliTest(false, false, "bootenvs", "destroy", "fredhammer").run(t)
	cliTest(false, false, "bootenvs", "install", "bootenvs/fredhammer.yml").run(t)
	cliTest(false, false, "bootenvs", "uploadiso", "fredhammer").run(t)
	cliTest(false, false, "bootenvs", "install", "test-data/no-fredhammer.yml").run(t)

	// Clean up
	cliTest(false, false, "bootenvs", "destroy", "fredhammer").run(t)
	cliTest(false, false, "bootenvs", "destroy", "no-fredhammer").run(t)
	cliTest(false, false, "bootenvs", "destroy", "local3").run(t)
	cliTest(false, false, "templates", "destroy", "local3-pxelinux.tmpl").run(t)
	cliTest(false, false, "templates", "destroy", "local3-elilo.tmpl").run(t)
	cliTest(false, false, "templates", "destroy", "local3-ipxe.tmpl").run(t)
	cliTest(false, false, "isos", "destroy", "sledgehammer-708de8b878e3818b1c1bb598a56de968939f9d4b.tar").run(t)

	// Make sure that ic exists and iso exists
	// if _, err := os.Stat("ic"); os.IsNotExist(err) {
	//	t.Errorf("FAIL: Failed to create ic directory\n")
	// }
	if _, err := os.Stat("isos"); os.IsNotExist(err) {
		t.Errorf("FAIL: Failed to create isos directory\n")
	}

	os.RemoveAll("bootenvs")
	os.RemoveAll("templates")
	os.RemoveAll("isos")
	os.RemoveAll("ic")
	os.RemoveAll(path.Join(tmpDir, "tftpboot", "sledgehammer"))
	verifyClean(t)
}

func TestBootEnvLookaside(t *testing.T) {
	testFile := "sledgehammer/708de8b878e3818b1c1bb598a56de968939f9d4b/vmlinuz0"
	cliTest(false, false, "profiles", "add", "global", "param", "package-repositories", "to", "-").Stdin(`
- tag: "sledgehammer-708de8b878e3818b1c1bb598a56de968939f9d4b"
  os:
    - "sledgehammer/708de8b878e3818b1c1bb598a56de968939f9d4b"
  installSource: true
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
				t.Logf("tftp: Lookaside from \nsledgehammer/708de8b878e3818b1c1bb598a56de968939f9d4b/vmlinuz0\nto\nhttp://127.0.0.1:10003/hammertime/vmlinuz0\nworked")
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
