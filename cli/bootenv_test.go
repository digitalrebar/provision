package cli

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/digitalrebar/provision/backend"
	"github.com/digitalrebar/provision/midlayer"
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
  "Kernel": "lpxelinux.0"
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
	cliTest(false, true, "bootenvs", "exists", "john", "john2").run(t)

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

	midlayer.ServeStatic("127.0.0.1:10003", backend.NewFS("test-data", nil), nil, backend.NewPublishers(nil))

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

	// Clean up
	cliTest(false, false, "bootenvs", "destroy", "fredhammer").run(t)
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
	verifyClean(t)
}
