package cli

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/digitalrebar/provision/backend"
	"github.com/digitalrebar/provision/midlayer"
)

var (
	limitNegativeError                           = "Error: GET: bootenvs: Limit cannot be negative\n\n"
	offsetNegativeError                          = "Error: GET: bootenvs: Offset cannot be negative\n\n"
	bootEnvShowNoArgErrorString                  = "Error: drpcli bootenvs show [id] [flags] requires 1 argument\n"
	bootEnvShowTooManyArgErrorString             = "Error: drpcli bootenvs show [id] [flags] requires 1 argument\n"
	bootEnvShowMissingArgErrorString             = "Error: GET: bootenvs/john: Not Found\n\n"
	bootEnvCreateBadJSONErrorString              = "Error: CREATE: bootenvs: Empty key not allowed\n\n"
	bootEnvExistsNoArgErrorString                = "Error: drpcli bootenvs exists [id] [flags] requires 1 argument"
	bootEnvExistsTooManyArgErrorString           = "Error: drpcli bootenvs exists [id] [flags] requires 1 argument"
	bootEnvExistsMissingJohnString               = "Error: GET: bootenvs/john: Not Found\n\n"
	bootEnvCreateNoArgErrorString                = "Error: drpcli bootenvs create [json] [flags] requires 1 argument\n"
	bootEnvCreateTooManyArgErrorString           = "Error: drpcli bootenvs create [json] [flags] requires 1 argument\n"
	bootEnvCreateDuplicateErrorString            = "Error: CREATE: bootenvs/john: already exists\n\n"
	bootEnvPatchJohnMissingErrorString           = "Error: PATCH: bootenvs/john2: Not Found\n\n"
	bootEnvDestroyNoArgErrorString               = "Error: drpcli bootenvs destroy [id] [flags] requires 1 argument"
	bootEnvDestroyTooManyArgErrorString          = "Error: drpcli bootenvs destroy [id] [flags] requires 1 argument"
	bootEnvPatchBadBaseJSONErrorString           = "Error: Unable to parse drpcli bootenvs patch [objectJson] [changesJson] [flags] JSON asdgasdg\nError: error unmarshaling JSON: json: cannot unmarshal string into Go value of type genmodels.BootEnv\n\n"
	bootEnvPatchBadPatchJSONErrorString          = "Error: Unable to parse drpcli bootenvs patch [objectJson] [changesJson] [flags] JSON asdgasdg\nError: error unmarshaling JSON: json: cannot unmarshal string into Go value of type genmodels.BootEnv\n\n"
	bootEnvUpdateBadJSONErrorString              = "Error: Unable to merge objects: json: cannot unmarshal string into Go value of type map[string]interface {}\n\n\n"
	bootEnvUpdateJohnMissingErrorString          = "Error: GET: bootenvs/john2: Not Found\n\n"
	bootEnvPatchNoArgErrorString                 = "Error: drpcli bootenvs patch [objectJson] [changesJson] [flags] requires 2 arguments"
	bootEnvPatchTooManyArgErrorString            = "Error: drpcli bootenvs patch [objectJson] [changesJson] [flags] requires 2 arguments"
	bootEnvUpdateNoArgErrorString                = "Error: drpcli bootenvs update [id] [json] [flags] requires 2 arguments"
	bootEnvUpdateTooManyArgErrorString           = "Error: drpcli bootenvs update [id] [json] [flags] requires 2 arguments"
	bootEnvDestroyMissingJohnString              = "Error: DELETE: bootenvs/john: Not Found\n\n"
	bootEnvInstallNoArgUsageString               = "Error: drpcli bootenvs install [bootenvFile] [isoPath] [flags] needs at least 1 arg\n"
	bootEnvInstallTooManyArgUsageString          = "Error: drpcli bootenvs install [bootenvFile] [isoPath] [flags] has Too many args\n"
	bootEnvInstallBadBootEnvDirErrorString       = "Error: Error determining whether bootenvs dir exists: stat bootenvs: no such file or directory\n\n"
	bootEnvInstallBootEnvDirIsFileErrorString    = "Error: bootenvs is not a directory\n\n"
	bootEnvInstallNoSledgehammerErrorString      = "Error: No bootenv bootenvs/fredhammer.yml\n\n"
	bootEnvInstallSledgehammerBadJsonErrorString = "Error: Invalid bootenv object: error unmarshaling JSON: json: cannot unmarshal string into Go value of type genmodels.BootEnv\n\n\n"
	bootEnvBadReadOnlyString                     = "Error: GET: bootenvs: ReadOnly must be true or false\n\n"
	bootEnvBadAvailableString                    = "Error: GET: bootenvs: Available must be true or false\n\n"
	bootEnvBadValidString                        = "Error: GET: bootenvs: Valid must be true or false\n\n"
)

var bootEnvEmptyListString string = "[]\n"
var bootEnvIgnoreOnlyListString string = `[
  {
    "Available": true,
    "BootParams": "",
    "Description": "The boot environment you should use to have unknown machines boot off their local hard drive",
    "Errors": [],
    "Initrds": [],
    "Kernel": "",
    "Name": "ignore",
    "OS": {
      "Name": "ignore"
    },
    "OnlyUnknown": true,
    "OptionalParams": [],
    "ReadOnly": true,
    "RequiredParams": [],
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
    ],
    "Validated": true
  }
]
`
var bootEnvLocalOnlyListString string = `[
  {
    "Available": true,
    "BootParams": "",
    "Description": "The boot environment you should use to have known machines boot off their local hard drive",
    "Errors": [],
    "Initrds": [],
    "Kernel": "",
    "Name": "local",
    "OS": {
      "Name": "local"
    },
    "OnlyUnknown": false,
    "OptionalParams": [],
    "ReadOnly": true,
    "RequiredParams": [],
    "Templates": [
      {
        "Contents": "DEFAULT local\nPROMPT 0\nTIMEOUT 10\nLABEL local\nlocalboot 0\n",
        "Name": "pxelinux",
        "Path": "pxelinux.cfg/{{.Machine.HexAddress}}"
      },
      {
        "Contents": "exit",
        "Name": "elilo",
        "Path": "{{.Machine.HexAddress}}.conf"
      },
      {
        "Contents": "#!ipxe\nexit\n",
        "Name": "ipxe",
        "Path": "{{.Machine.Address}}.ipxe"
      }
    ],
    "Validated": true
  }
]
`
var bootEnvDefaultListString string = `[
  {
    "Available": true,
    "BootParams": "",
    "Description": "The boot environment you should use to have unknown machines boot off their local hard drive",
    "Errors": [],
    "Initrds": [],
    "Kernel": "",
    "Name": "ignore",
    "OS": {
      "Name": "ignore"
    },
    "OnlyUnknown": true,
    "OptionalParams": [],
    "ReadOnly": true,
    "RequiredParams": [],
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
    ],
    "Validated": true
  },
  {
    "Available": true,
    "BootParams": "",
    "Description": "The boot environment you should use to have known machines boot off their local hard drive",
    "Errors": [],
    "Initrds": [],
    "Kernel": "",
    "Name": "local",
    "OS": {
      "Name": "local"
    },
    "OnlyUnknown": false,
    "OptionalParams": [],
    "ReadOnly": true,
    "RequiredParams": [],
    "Templates": [
      {
        "Contents": "DEFAULT local\nPROMPT 0\nTIMEOUT 10\nLABEL local\nlocalboot 0\n",
        "Name": "pxelinux",
        "Path": "pxelinux.cfg/{{.Machine.HexAddress}}"
      },
      {
        "Contents": "exit",
        "Name": "elilo",
        "Path": "{{.Machine.HexAddress}}.conf"
      },
      {
        "Contents": "#!ipxe\nexit\n",
        "Name": "ipxe",
        "Path": "{{.Machine.Address}}.ipxe"
      }
    ],
    "Validated": true
  }
]
`
var bootEnvShowIgnoreString string = `{
  "Available": true,
  "BootParams": "",
  "Description": "The boot environment you should use to have unknown machines boot off their local hard drive",
  "Errors": [],
  "Initrds": [],
  "Kernel": "",
  "Name": "ignore",
  "OS": {
    "Name": "ignore"
  },
  "OnlyUnknown": true,
  "OptionalParams": [],
  "ReadOnly": true,
  "RequiredParams": [],
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
  ],
  "Validated": true
}
`

var bootEnvExistsIgnoreString string = ""

var bootEnvCreateBadJSONString = "{asdgasdg}"

var bootEnvCreateInputString string = `{
  "name": "john"
}
`
var bootEnvCreateJohnString string = `{
  "Available": false,
  "BootParams": "",
  "Errors": [
    "bootenv: Missing elilo or pxelinux template"
  ],
  "Initrds": [],
  "Kernel": "",
  "Name": "john",
  "OS": {
    "Name": ""
  },
  "OnlyUnknown": false,
  "OptionalParams": [],
  "ReadOnly": false,
  "RequiredParams": [],
  "Templates": [],
  "Validated": true
}
`
var bootEnvCreateFredInputString string = `fred`
var bootEnvCreateFredString string = `{
  "Available": false,
  "BootParams": "",
  "Errors": [
    "bootenv: Missing elilo or pxelinux template"
  ],
  "Initrds": [],
  "Kernel": "",
  "Name": "fred",
  "OS": {
    "Name": ""
  },
  "OnlyUnknown": false,
  "OptionalParams": [],
  "ReadOnly": false,
  "RequiredParams": [],
  "Templates": [],
  "Validated": true
}
`
var bootEnvDeleteFredString string = "Deleted bootenv fred\n"

var bootEnvListBothEnvsString = `[
  {
    "Available": true,
    "BootParams": "",
    "Description": "The boot environment you should use to have unknown machines boot off their local hard drive",
    "Errors": [],
    "Initrds": [],
    "Kernel": "",
    "Name": "ignore",
    "OS": {
      "Name": "ignore"
    },
    "OnlyUnknown": true,
    "OptionalParams": [],
    "ReadOnly": true,
    "RequiredParams": [],
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
    ],
    "Validated": true
  },
  {
    "Available": false,
    "BootParams": "",
    "Errors": [
      "bootenv: Missing elilo or pxelinux template"
    ],
    "Initrds": [],
    "Kernel": "",
    "Name": "john",
    "OS": {
      "Name": ""
    },
    "OnlyUnknown": false,
    "OptionalParams": [],
    "ReadOnly": false,
    "RequiredParams": [],
    "Templates": [],
    "Validated": true
  },
  {
    "Available": true,
    "BootParams": "",
    "Description": "The boot environment you should use to have known machines boot off their local hard drive",
    "Errors": [],
    "Initrds": [],
    "Kernel": "",
    "Name": "local",
    "OS": {
      "Name": "local"
    },
    "OnlyUnknown": false,
    "OptionalParams": [],
    "ReadOnly": true,
    "RequiredParams": [],
    "Templates": [
      {
        "Contents": "DEFAULT local\nPROMPT 0\nTIMEOUT 10\nLABEL local\nlocalboot 0\n",
        "Name": "pxelinux",
        "Path": "pxelinux.cfg/{{.Machine.HexAddress}}"
      },
      {
        "Contents": "exit",
        "Name": "elilo",
        "Path": "{{.Machine.HexAddress}}.conf"
      },
      {
        "Contents": "#!ipxe\nexit\n",
        "Name": "ipxe",
        "Path": "{{.Machine.Address}}.ipxe"
      }
    ],
    "Validated": true
  }
]
`

var bootEnvUpdateBadJSONString = "asdgasdg"

var bootEnvUpdateInputString string = `{
  "Kernel": "lpxelinux.0"
}
`
var bootEnvUpdateJohnString string = `{
  "Available": false,
  "BootParams": "",
  "Errors": [
    "bootenv: Missing elilo or pxelinux template"
  ],
  "Initrds": [],
  "Kernel": "lpxelinux.0",
  "Name": "john",
  "OS": {
    "Name": ""
  },
  "OnlyUnknown": false,
  "OptionalParams": [],
  "ReadOnly": false,
  "RequiredParams": [],
  "Templates": [],
  "Validated": true
}
`

var bootEnvPatchBadPatchJSONString = "asdgasdg"

var bootEnvPatchBadBaseJSONString = "asdgasdg"

var bootEnvPatchBaseString string = `{
  "Available": false,
  "BootParams": "",
  "Errors": [
    "bootenv: Missing elilo or pxelinux template"
  ],
  "Initrds": [],
  "Kernel": "lpxelinux.0",
  "Name": "john",
  "OS": {
    "Name": ""
  },
  "OnlyUnknown": false,
  "OptionalParams": [],
  "ReadOnly": false,
  "RequiredParams": [],
  "Templates": [],
  "Validated": true
}
`
var bootEnvPatchInputString string = `{
  "Kernel": "bootx64.efi"
}
`
var bootEnvPatchJohnString string = `{
  "Available": false,
  "BootParams": "",
  "Errors": [
    "bootenv: Missing elilo or pxelinux template"
  ],
  "Initrds": [],
  "Kernel": "bootx64.efi",
  "Name": "john",
  "OS": {
    "Name": ""
  },
  "OnlyUnknown": false,
  "OptionalParams": [],
  "ReadOnly": false,
  "RequiredParams": [],
  "Templates": [],
  "Validated": true
}
`
var bootEnvPatchMissingBaseString string = `{
  "Available": false,
  "BootParams": "",
  "Errors": [
    "bootenv: Missing elilo or pxelinux template"
  ],
  "Initrds": [],
  "Kernel": "bootx64.efi",
  "Name": "john2",
  "OS": {
    "Name": ""
  },
  "OnlyUnknown": false,
  "OptionalParams": [],
  "RequiredParams": [],
  "Templates": [],
  "Validated": true
}
`

var bootEnvDestroyJohnString string = "Deleted bootenv john\n"

var bootEnvInstallSledgehammerSuccessWithErrorsString string = `RE:
{
  "Available": false,
  "BootParams": "rootflags=loop root=live:/sledgehammer.iso rootfstype=auto ro liveimg rd_NO_LUKS rd_NO_MD rd_NO_DM provisioner.web={{.ProvisionerURL}} rebar.web={{.CommandURL}} rs.uuid={{.Machine.UUID}} rs.api={{.ApiURL}}",
  "Errors": \[
[\s\S]*
  \],
  "Initrds": \[
    "stage1.img"
  \],
[\s\S]*
}
`

var bootEnvInstallSledgehammerSuccessString string = `RE:
{
  "Available": [\s\S]*,
  "BootParams": "rootflags=loop root=live:/sledgehammer.iso rootfstype=auto ro liveimg rd_NO_LUKS rd_NO_MD rd_NO_DM provisioner.web={{.ProvisionerURL}} rebar.web={{.CommandURL}} rs.uuid={{.Machine.UUID}} rs.api={{.ApiURL}}",
  "Errors": [\s\S]*,
  "Initrds": \[
    "stage1.img"
  \],
  "Kernel": "vmlinuz0",
  "Name": "fredhammer",
  "OS": {
    "IsoFile": "sledgehammer-708de8b878e3818b1c1bb598a56de968939f9d4b.tar",
    [\s\S]*,
    "Name": "sledgehammer/708de8b878e3818b1c1bb598a56de968939f9d4b"
  },
  "OnlyUnknown": false,
  "OptionalParams": \[
    "ntp_servers",
    "access_keys"
  \],
  "ReadOnly": false,
  "RequiredParams": [],
  "Templates": \[
[\s\S]*
  \],
  "Validated": [\s\S]*
}
`

var bootEnvInstallLocalMissingTemplatesErrorString string = "Installing template local3-pxelinux.tmpl\nError: Unable to find template: local3-pxelinux.tmpl: open templates/local3-pxelinux.tmpl: no such file or directory\n\n"

var bootEnvInstallLocalSuccessString string = `{
  "Available": true,
  "BootParams": "",
  "Errors": [],
  "Initrds": [],
  "Kernel": "",
  "Name": "local3",
  "OS": {
    "Name": "local3"
  },
  "OnlyUnknown": false,
  "OptionalParams": [],
  "ReadOnly": false,
  "RequiredParams": [],
  "Templates": [
    {
      "ID": "local3-pxelinux.tmpl",
      "Name": "pxelinux",
      "Path": "pxelinux.cfg/{{.Machine.HexAddress}}"
    },
    {
      "ID": "local3-elilo.tmpl",
      "Name": "elilo",
      "Path": "{{.Machine.HexAddress}}.conf"
    },
    {
      "ID": "local3-ipxe.tmpl",
      "Name": "ipxe",
      "Path": "{{.Machine.Address}}.ipxe"
    }
  ],
  "Validated": true
}
`

var bootEnvSkipDownloadErrorString = "Installing bootenv fredhammer\nSkipping ISO download as requested\nUpload with `drpcli isos upload sledgehammer-708de8b878e3818b1c1bb598a56de968939f9d4b.tar as sledgehammer-708de8b878e3818b1c1bb598a56de968939f9d4b.tar` when you have it\n"

var bootEnvDownloadErrorString = "Installing bootenv fredhammer\nDownloading http://127.0.0.1:10003/sledgehammer-708de8b878e3818b1c1bb598a56de968939f9d4b.tar to isos/sledgehammer-708de8b878e3818b1c1bb598a56de968939f9d4b.tar\nDownloaded 5120 bytes\nUploading isos/sledgehammer-708de8b878e3818b1c1bb598a56de968939f9d4b.tar to DigitalRebar Provision\n"

var bootEnvInstallLocal3ErrorString = "Installing template local3-pxelinux.tmpl\nInstalling template local3-elilo.tmpl\nInstalling template local3-ipxe.tmpl\nInstalling bootenv local3\n"

func TestBootEnvCli(t *testing.T) {
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
}
