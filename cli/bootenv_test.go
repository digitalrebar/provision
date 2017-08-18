package cli

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/digitalrebar/provision/backend"
	"github.com/digitalrebar/provision/midlayer"
)

var limitNegativeError string = "Error: Limit cannot be negative\n\n"
var offsetNegativeError string = "Error: Offset cannot be negative\n\n"

var bootEnvEmptyListString string = "[]\n"
var bootEnvIgnoreOnlyListString string = `[
  {
    "Available": true,
    "BootParams": "",
    "Description": "The boot environment you should use to have unknown machines boot off their local hard drive",
    "Errors": null,
    "Initrds": null,
    "Kernel": "",
    "Name": "ignore",
    "OS": {
      "Name": "ignore"
    },
    "OnlyUnknown": true,
    "OptionalParams": null,
    "RequiredParams": null,
    "Tasks": null,
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
        "Contents": "#!ipxe\nexit\n",
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
    "Errors": null,
    "Initrds": null,
    "Kernel": "",
    "Name": "local",
    "OS": {
      "Name": "local"
    },
    "OnlyUnknown": false,
    "OptionalParams": null,
    "RequiredParams": null,
    "Tasks": null,
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
    "Initrds": null,
    "Kernel": "",
    "Name": "ignore",
    "OS": {
      "Name": "ignore"
    },
    "OnlyUnknown": true,
    "OptionalParams": null,
    "RequiredParams": null,
    "Tasks": null,
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
        "Contents": "#!ipxe\nexit\n",
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
    "Errors": null,
    "Initrds": null,
    "Kernel": "",
    "Name": "local",
    "OS": {
      "Name": "local"
    },
    "OnlyUnknown": false,
    "OptionalParams": null,
    "RequiredParams": null,
    "Tasks": null,
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

var bootEnvShowNoArgErrorString string = "Error: drpcli bootenvs show [id] [flags] requires 1 argument\n"
var bootEnvShowTooManyArgErrorString string = "Error: drpcli bootenvs show [id] [flags] requires 1 argument\n"
var bootEnvShowMissingArgErrorString string = "Error: bootenvs GET: john: Not Found\n\n"
var bootEnvShowIgnoreString string = `{
  "Available": true,
  "BootParams": "",
  "Description": "The boot environment you should use to have unknown machines boot off their local hard drive",
  "Errors": [],
  "Initrds": null,
  "Kernel": "",
  "Name": "ignore",
  "OS": {
    "Name": "ignore"
  },
  "OnlyUnknown": true,
  "OptionalParams": null,
  "RequiredParams": null,
  "Tasks": null,
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
      "Contents": "#!ipxe\nexit\n",
      "Name": "ipxe",
      "Path": "default.ipxe"
    }
  ],
  "Validated": true
}
`

var bootEnvExistsNoArgErrorString string = "Error: drpcli bootenvs exists [id] [flags] requires 1 argument"
var bootEnvExistsTooManyArgErrorString string = "Error: drpcli bootenvs exists [id] [flags] requires 1 argument"
var bootEnvExistsIgnoreString string = ""
var bootEnvExistsMissingJohnString string = "Error: bootenvs GET: john: Not Found\n\n"

var bootEnvCreateNoArgErrorString string = "Error: drpcli bootenvs create [json] [flags] requires 1 argument\n"
var bootEnvCreateTooManyArgErrorString string = "Error: drpcli bootenvs create [json] [flags] requires 1 argument\n"
var bootEnvCreateBadJSONString = "{asdgasdg}"
var bootEnvCreateBadJSONErrorString = "Error: dataTracker create bootenvs: Empty key not allowed\n\n"
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
  "Initrds": null,
  "Kernel": "",
  "Name": "john",
  "OS": {
    "Name": ""
  },
  "OnlyUnknown": false,
  "OptionalParams": null,
  "RequiredParams": null,
  "Tasks": null,
  "Templates": null,
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
  "Initrds": null,
  "Kernel": "",
  "Name": "fred",
  "OS": {
    "Name": ""
  },
  "OnlyUnknown": false,
  "OptionalParams": null,
  "RequiredParams": null,
  "Tasks": null,
  "Templates": null,
  "Validated": true
}
`
var bootEnvDeleteFredString string = "Deleted bootenv fred\n"
var bootEnvCreateDuplicateErrorString = "Error: dataTracker create bootenvs: john already exists\n\n"

var bootEnvListBothEnvsString = `[
  {
    "Available": true,
    "BootParams": "",
    "Description": "The boot environment you should use to have unknown machines boot off their local hard drive",
    "Errors": [],
    "Initrds": null,
    "Kernel": "",
    "Name": "ignore",
    "OS": {
      "Name": "ignore"
    },
    "OnlyUnknown": true,
    "OptionalParams": null,
    "RequiredParams": null,
    "Tasks": null,
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
        "Contents": "#!ipxe\nexit\n",
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
    "Initrds": null,
    "Kernel": "",
    "Name": "john",
    "OS": {
      "Name": ""
    },
    "OnlyUnknown": false,
    "OptionalParams": null,
    "RequiredParams": null,
    "Tasks": null,
    "Templates": null,
    "Validated": true
  },
  {
    "Available": true,
    "BootParams": "",
    "Description": "The boot environment you should use to have known machines boot off their local hard drive",
    "Errors": null,
    "Initrds": null,
    "Kernel": "",
    "Name": "local",
    "OS": {
      "Name": "local"
    },
    "OnlyUnknown": false,
    "OptionalParams": null,
    "RequiredParams": null,
    "Tasks": null,
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

var bootEnvUpdateNoArgErrorString string = "Error: drpcli bootenvs update [id] [json] [flags] requires 2 arguments"
var bootEnvUpdateTooManyArgErrorString string = "Error: drpcli bootenvs update [id] [json] [flags] requires 2 arguments"
var bootEnvUpdateBadJSONString = "asdgasdg"
var bootEnvUpdateBadJSONErrorString = "Error: Unable to merge objects: json: cannot unmarshal string into Go value of type map[string]interface {}\n\n\n"
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
  "Initrds": null,
  "Kernel": "lpxelinux.0",
  "Name": "john",
  "OS": {
    "Name": ""
  },
  "OnlyUnknown": false,
  "OptionalParams": null,
  "RequiredParams": null,
  "Tasks": null,
  "Templates": null,
  "Validated": true
}
`
var bootEnvUpdateJohnMissingErrorString string = "Error: bootenvs GET: john2: Not Found\n\n"

var bootEnvPatchNoArgErrorString string = "Error: drpcli bootenvs patch [objectJson] [changesJson] [flags] requires 2 arguments"
var bootEnvPatchTooManyArgErrorString string = "Error: drpcli bootenvs patch [objectJson] [changesJson] [flags] requires 2 arguments"
var bootEnvPatchBadPatchJSONString = "asdgasdg"
var bootEnvPatchBadPatchJSONErrorString = "Error: Unable to parse drpcli bootenvs patch [objectJson] [changesJson] [flags] JSON asdgasdg\nError: error unmarshaling JSON: json: cannot unmarshal string into Go value of type genmodels.BootEnv\n\n"
var bootEnvPatchBadBaseJSONString = "asdgasdg"
var bootEnvPatchBadBaseJSONErrorString = "Error: Unable to parse drpcli bootenvs patch [objectJson] [changesJson] [flags] JSON asdgasdg\nError: error unmarshaling JSON: json: cannot unmarshal string into Go value of type genmodels.BootEnv\n\n"
var bootEnvPatchBaseString string = `{
  "Available": false,
  "BootParams": "",
  "Errors": [
    "bootenv: Missing elilo or pxelinux template"
  ],
  "Initrds": null,
  "Kernel": "lpxelinux.0",
  "Name": "john",
  "OS": {
    "Name": ""
  },
  "OnlyUnknown": false,
  "OptionalParams": null,
  "RequiredParams": null,
  "Tasks": null,
  "Templates": null,
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
  "Initrds": null,
  "Kernel": "bootx64.efi",
  "Name": "john",
  "OS": {
    "Name": ""
  },
  "OnlyUnknown": false,
  "OptionalParams": null,
  "RequiredParams": null,
  "Tasks": null,
  "Templates": null,
  "Validated": true
}
`
var bootEnvPatchMissingBaseString string = `{
  "Available": false,
  "BootParams": "",
  "Errors": [
    "bootenv: Missing elilo or pxelinux template"
  ],
  "Initrds": null,
  "Kernel": "bootx64.efi",
  "Name": "john2",
  "OS": {
    "Name": ""
  },
  "OnlyUnknown": false,
  "OptionalParams": null,
  "RequiredParams": null,
  "Tasks": null,
  "Templates": null,
  "Validated": true
}
`
var bootEnvPatchJohnMissingErrorString string = "Error: bootenvs: PATCH john2: Not Found\n\n"

var bootEnvDestroyNoArgErrorString string = "Error: drpcli bootenvs destroy [id] [flags] requires 1 argument"
var bootEnvDestroyTooManyArgErrorString string = "Error: drpcli bootenvs destroy [id] [flags] requires 1 argument"
var bootEnvDestroyJohnString string = "Deleted bootenv john\n"
var bootEnvDestroyMissingJohnString string = "Error: bootenvs: DELETE john: Not Found\n\n"

var bootEnvInstallNoArgUsageString string = "Error: drpcli bootenvs install [bootenvFile] [isoPath] [flags] needs at least 1 arg\n"
var bootEnvInstallTooManyArgUsageString string = "Error: drpcli bootenvs install [bootenvFile] [isoPath] [flags] has Too many args\n"
var bootEnvInstallBadBootEnvDirErrorString string = "Error: Error determining whether bootenvs dir exists: stat bootenvs: no such file or directory\n\n"
var bootEnvInstallBootEnvDirIsFileErrorString string = "Error: bootenvs is not a directory\n\n"
var bootEnvInstallNoSledgehammerErrorString string = "Error: No bootenv bootenvs/fredhammer.yml\n\n"
var bootEnvInstallSledgehammerBadJsonErrorString string = "Error: Invalid bootenv object: error unmarshaling JSON: json: cannot unmarshal string into Go value of type genmodels.BootEnv\n\n\n"

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
  "RequiredParams": null,
  "Tasks": null,
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
  "Initrds": null,
  "Kernel": "",
  "Name": "local3",
  "OS": {
    "Name": "local3"
  },
  "OnlyUnknown": false,
  "OptionalParams": null,
  "RequiredParams": null,
  "Tasks": null,
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
	tests := []CliTest{
		CliTest{true, false, []string{"bootenvs"}, noStdinString, "Access CLI commands relating to bootenvs\n", ""},
		CliTest{false, false, []string{"bootenvs", "list"}, noStdinString, bootEnvDefaultListString, noErrorString},
		CliTest{false, false, []string{"bootenvs", "list", "--limit=0"}, noStdinString, bootEnvEmptyListString, noErrorString},
		CliTest{false, false, []string{"bootenvs", "list", "--limit=10", "--offset=0"}, noStdinString, bootEnvDefaultListString, noErrorString},
		CliTest{false, false, []string{"bootenvs", "list", "--limit=10", "--offset=10"}, noStdinString, bootEnvEmptyListString, noErrorString},
		CliTest{false, true, []string{"bootenvs", "list", "--limit=-10", "--offset=0"}, noStdinString, noContentString, limitNegativeError},
		CliTest{false, true, []string{"bootenvs", "list", "--limit=10", "--offset=-10"}, noStdinString, noContentString, offsetNegativeError},
		CliTest{false, false, []string{"bootenvs", "list", "--limit=-1", "--offset=-1"}, noStdinString, bootEnvDefaultListString, noErrorString},
		CliTest{false, false, []string{"bootenvs", "list", "Name=fred"}, noStdinString, bootEnvEmptyListString, noErrorString},
		CliTest{false, false, []string{"bootenvs", "list", "Name=ignore"}, noStdinString, bootEnvIgnoreOnlyListString, noErrorString},
		CliTest{false, false, []string{"bootenvs", "list", "Available=true"}, noStdinString, bootEnvDefaultListString, noErrorString},
		CliTest{false, false, []string{"bootenvs", "list", "Available=false"}, noStdinString, bootEnvEmptyListString, noErrorString},
		CliTest{false, false, []string{"bootenvs", "list", "OnlyUnknown=true"}, noStdinString, bootEnvIgnoreOnlyListString, noErrorString},
		CliTest{false, false, []string{"bootenvs", "list", "OnlyUnknown=false"}, noStdinString, bootEnvLocalOnlyListString, noErrorString},

		CliTest{true, true, []string{"bootenvs", "show"}, noStdinString, noContentString, bootEnvShowNoArgErrorString},
		CliTest{true, true, []string{"bootenvs", "show", "john", "john2"}, noStdinString, noContentString, bootEnvShowTooManyArgErrorString},
		CliTest{false, true, []string{"bootenvs", "show", "john"}, noStdinString, noContentString, bootEnvShowMissingArgErrorString},
		CliTest{false, false, []string{"bootenvs", "show", "ignore"}, noStdinString, bootEnvShowIgnoreString, noErrorString},

		CliTest{true, true, []string{"bootenvs", "exists"}, noStdinString, noContentString, bootEnvExistsNoArgErrorString},
		CliTest{true, true, []string{"bootenvs", "exists", "john", "john2"}, noStdinString, noContentString, bootEnvExistsTooManyArgErrorString},
		CliTest{false, false, []string{"bootenvs", "exists", "ignore"}, noStdinString, bootEnvExistsIgnoreString, noErrorString},
		CliTest{false, true, []string{"bootenvs", "exists", "john"}, noStdinString, noContentString, bootEnvExistsMissingJohnString},
		CliTest{true, true, []string{"bootenvs", "exists", "john", "john2"}, noStdinString, noContentString, bootEnvExistsTooManyArgErrorString},

		CliTest{true, true, []string{"bootenvs", "create"}, noStdinString, noContentString, bootEnvCreateNoArgErrorString},
		CliTest{true, true, []string{"bootenvs", "create", "john", "john2"}, noStdinString, noContentString, bootEnvCreateTooManyArgErrorString},
		CliTest{false, true, []string{"bootenvs", "create", bootEnvCreateBadJSONString}, noStdinString, noContentString, bootEnvCreateBadJSONErrorString},
		CliTest{false, false, []string{"bootenvs", "create", bootEnvCreateInputString}, noStdinString, bootEnvCreateJohnString, noErrorString},
		CliTest{false, false, []string{"bootenvs", "create", bootEnvCreateFredInputString}, noStdinString, bootEnvCreateFredString, noErrorString},
		CliTest{false, false, []string{"bootenvs", "destroy", bootEnvCreateFredInputString}, noStdinString, bootEnvDeleteFredString, noErrorString},
		CliTest{false, true, []string{"bootenvs", "create", bootEnvCreateInputString}, noStdinString, noContentString, bootEnvCreateDuplicateErrorString},
		CliTest{false, false, []string{"bootenvs", "list"}, noStdinString, bootEnvListBothEnvsString, noErrorString},

		CliTest{true, true, []string{"bootenvs", "update"}, noStdinString, noContentString, bootEnvUpdateNoArgErrorString},
		CliTest{true, true, []string{"bootenvs", "update", "john", "john2", "john3"}, noStdinString, noContentString, bootEnvUpdateTooManyArgErrorString},
		CliTest{false, true, []string{"bootenvs", "update", "john", bootEnvUpdateBadJSONString}, noStdinString, noContentString, bootEnvUpdateBadJSONErrorString},
		CliTest{false, false, []string{"bootenvs", "update", "john", bootEnvUpdateInputString}, noStdinString, bootEnvUpdateJohnString, noErrorString},
		CliTest{false, true, []string{"bootenvs", "update", "john2", bootEnvUpdateInputString}, noStdinString, noContentString, bootEnvUpdateJohnMissingErrorString},
		CliTest{false, false, []string{"bootenvs", "show", "john"}, noStdinString, bootEnvUpdateJohnString, noErrorString},

		CliTest{true, true, []string{"bootenvs", "patch"}, noStdinString, noContentString, bootEnvPatchNoArgErrorString},
		CliTest{true, true, []string{"bootenvs", "patch", "john", "john2", "john3"}, noStdinString, noContentString, bootEnvPatchTooManyArgErrorString},
		CliTest{false, true, []string{"bootenvs", "patch", bootEnvPatchBaseString, bootEnvPatchBadPatchJSONString}, noStdinString, noContentString, bootEnvPatchBadPatchJSONErrorString},
		CliTest{false, true, []string{"bootenvs", "patch", bootEnvPatchBadBaseJSONString, bootEnvPatchInputString}, noStdinString, noContentString, bootEnvPatchBadBaseJSONErrorString},
		CliTest{false, false, []string{"bootenvs", "patch", bootEnvPatchBaseString, bootEnvPatchInputString}, noStdinString, bootEnvPatchJohnString, noErrorString},
		CliTest{false, true, []string{"bootenvs", "patch", bootEnvPatchMissingBaseString, bootEnvPatchInputString}, noStdinString, noContentString, bootEnvPatchJohnMissingErrorString},
		CliTest{false, false, []string{"bootenvs", "show", "john"}, noStdinString, bootEnvPatchJohnString, noErrorString},

		CliTest{true, true, []string{"bootenvs", "destroy"}, noStdinString, noContentString, bootEnvDestroyNoArgErrorString},
		CliTest{true, true, []string{"bootenvs", "destroy", "john", "june"}, noStdinString, noContentString, bootEnvDestroyTooManyArgErrorString},
		CliTest{false, false, []string{"bootenvs", "destroy", "john"}, noStdinString, bootEnvDestroyJohnString, noErrorString},
		CliTest{false, true, []string{"bootenvs", "destroy", "john"}, noStdinString, noContentString, bootEnvDestroyMissingJohnString},
		CliTest{false, false, []string{"bootenvs", "list"}, noStdinString, bootEnvDefaultListString, noErrorString},

		CliTest{false, false, []string{"bootenvs", "create", "-"}, bootEnvCreateInputString + "\n", bootEnvCreateJohnString, noErrorString},
		CliTest{false, false, []string{"bootenvs", "list"}, noStdinString, bootEnvListBothEnvsString, noErrorString},
		CliTest{false, false, []string{"bootenvs", "update", "john", "-"}, bootEnvUpdateInputString + "\n", bootEnvUpdateJohnString, noErrorString},
		CliTest{false, false, []string{"bootenvs", "show", "john"}, noStdinString, bootEnvUpdateJohnString, noErrorString},
		CliTest{false, false, []string{"bootenvs", "destroy", "john"}, noStdinString, bootEnvDestroyJohnString, noErrorString},
		CliTest{false, false, []string{"bootenvs", "list"}, noStdinString, bootEnvDefaultListString, noErrorString},

		CliTest{true, true, []string{"bootenvs", "install"}, noStdinString, noContentString, bootEnvInstallNoArgUsageString},
		CliTest{true, true, []string{"bootenvs", "install", "john", "john", "john2"}, noStdinString, noContentString, bootEnvInstallTooManyArgUsageString},
		CliTest{false, true, []string{"bootenvs", "install", "fredhammer"}, noStdinString, noContentString, bootEnvInstallBadBootEnvDirErrorString},
	}

	for _, test := range tests {
		testCli(t, test)
	}

	if f, err := os.Create("bootenvs"); err != nil {
		t.Errorf("Failed to create bootenvs file: %v\n", err)
	} else {
		f.Close()
	}

	tests = []CliTest{
		CliTest{false, true, []string{"bootenvs", "install", "bootenvs/fredhammer.yml"}, noStdinString, noContentString, bootEnvInstallBootEnvDirIsFileErrorString},
	}
	for _, test := range tests {
		testCli(t, test)
	}

	os.RemoveAll("bootenvs")
	if err := os.MkdirAll("bootenvs", 0755); err != nil {
		t.Errorf("Failed to create bootenvs dir: %v\n", err)
	}

	tests = []CliTest{
		CliTest{false, true, []string{"bootenvs", "install", "bootenvs/fredhammer.yml"}, noStdinString, noContentString, bootEnvInstallNoSledgehammerErrorString},
	}
	for _, test := range tests {
		testCli(t, test)
	}

	if err := ioutil.WriteFile("bootenvs/fredhammer.yml", []byte("TEST"), 0644); err != nil {
		t.Errorf("Failed to create bootenvs file: %v\n", err)
	}

	tests = []CliTest{
		CliTest{false, true, []string{"bootenvs", "install", "bootenvs/fredhammer.yml"}, noStdinString, noContentString, bootEnvInstallSledgehammerBadJsonErrorString},
	}
	for _, test := range tests {
		testCli(t, test)
	}

	midlayer.ServeStatic("127.0.0.1:10003", backend.NewFS("test-data", nil), nil, backend.NewPublishers(nil))

	os.RemoveAll("bootenvs/fredhammer.yml")
	if err := os.MkdirAll("bootenvs", 0755); err != nil {
		t.Errorf("Failed to create bootenvs dir: %v\n", err)
	}
	if err := os.Symlink("../test-data/fredhammer.yml", "bootenvs/fredhammer.yml"); err != nil {
		t.Errorf("Failed to create link to fredhammer.yml: %v\n", err)
	}
	if err := os.Symlink("../test-data/local3.yml", "bootenvs/local3.yml"); err != nil {
		t.Errorf("Failed to create link to local3.yml: %v\n", err)
	}
	tests = []CliTest{
		CliTest{false, false, []string{"bootenvs", "install", "--skip-download", "bootenvs/fredhammer.yml"}, noStdinString, bootEnvInstallSledgehammerSuccessWithErrorsString, bootEnvSkipDownloadErrorString},
		CliTest{false, false, []string{"bootenvs", "destroy", "fredhammer"}, noStdinString, "Deleted bootenv fredhammer\n", noErrorString},
	}
	for _, test := range tests {
		testCli(t, test)
	}

	installSkipDownloadIsos = false
	tests = []CliTest{
		CliTest{false, false, []string{"bootenvs", "install", "bootenvs/fredhammer.yml"}, noStdinString, bootEnvInstallSledgehammerSuccessString, bootEnvDownloadErrorString},
		CliTest{false, true, []string{"bootenvs", "install", "bootenvs/local3.yml"}, noStdinString, noContentString, bootEnvInstallLocalMissingTemplatesErrorString},
	}
	for _, test := range tests {
		testCli(t, test)
	}

	if err := os.MkdirAll("templates", 0755); err != nil {
		t.Errorf("Failed to create templates dir: %v\n", err)
	}
	tmpls := []string{"local3-pxelinux.tmpl", "local3-elilo.tmpl", "local3-ipxe.tmpl"}
	for _, tmpl := range tmpls {
		if err := os.Symlink("../test-data/"+tmpl, "templates/"+tmpl); err != nil {
			t.Errorf("Failed to create link to %s: %v\n", tmpl, err)
		}
	}
	tests = []CliTest{
		CliTest{false, false, []string{"bootenvs", "install", "bootenvs/local3.yml", "ic"}, noStdinString, bootEnvInstallLocalSuccessString, bootEnvInstallLocal3ErrorString},
		CliTest{false, false, []string{"bootenvs", "destroy", "fredhammer"}, noStdinString, "Deleted bootenv fredhammer\n", noErrorString},
		CliTest{false, false, []string{"bootenvs", "install", "bootenvs/fredhammer.yml"}, noStdinString, bootEnvInstallSledgehammerSuccessString, "Installing bootenv fredhammer\n"},

		// Clean up
		CliTest{false, false, []string{"bootenvs", "destroy", "fredhammer"}, noStdinString, "Deleted bootenv fredhammer\n", noErrorString},
		CliTest{false, false, []string{"bootenvs", "destroy", "local3"}, noStdinString, "Deleted bootenv local3\n", noErrorString},
		CliTest{false, false, []string{"templates", "destroy", "local3-pxelinux.tmpl"}, noStdinString, "Deleted template local3-pxelinux.tmpl\n", noErrorString},
		CliTest{false, false, []string{"templates", "destroy", "local3-elilo.tmpl"}, noStdinString, "Deleted template local3-elilo.tmpl\n", noErrorString},
		CliTest{false, false, []string{"templates", "destroy", "local3-ipxe.tmpl"}, noStdinString, "Deleted template local3-ipxe.tmpl\n", noErrorString},
		CliTest{false, false, []string{"isos", "destroy", "sledgehammer-708de8b878e3818b1c1bb598a56de968939f9d4b.tar"}, noStdinString, "Deleted iso sledgehammer-708de8b878e3818b1c1bb598a56de968939f9d4b.tar\n", noErrorString},
	}
	for _, test := range tests {
		testCli(t, test)
	}

	// Make sure that ic exists and iso exists
	if _, err := os.Stat("ic"); os.IsNotExist(err) {
		t.Errorf("Failed to create ic directory\n")
	}
	if _, err := os.Stat("isos"); os.IsNotExist(err) {
		t.Errorf("Failed to create isos directory\n")
	}

	os.RemoveAll("bootenvs")
	os.RemoveAll("templates")
	os.RemoveAll("isos")
	os.RemoveAll("ic")
}
