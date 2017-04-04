package cli

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/rackn/rocket-skates/backend"
	"github.com/rackn/rocket-skates/midlayer"
)

var bootEnvDefaultListString string = `[
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
    "OptionalParams": null,
    "RequiredParams": null,
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
    ]
  }
]
`

var bootEnvShowNoArgErrorString string = "Error: rscli bootenvs show [id] requires 1 argument\n"
var bootEnvShowTooManyArgErrorString string = "Error: rscli bootenvs show [id] requires 1 argument\n"
var bootEnvShowMissingArgErrorString string = "Error: bootenvs GET: john: Not Found\n\n"
var bootEnvShowIgnoreString string = `{
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
  "OptionalParams": null,
  "RequiredParams": null,
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
  ]
}
`

var bootEnvExistsNoArgErrorString string = "Error: rscli bootenvs exists [id] requires 1 argument"
var bootEnvExistsTooManyArgErrorString string = "Error: rscli bootenvs exists [id] requires 1 argument"
var bootEnvExistsIgnoreString string = ""
var bootEnvExistsMissingJohnString string = "Error: bootenvs GET: john: Not Found\n\n"

var bootEnvCreateNoArgErrorString string = "Error: rscli bootenvs create [json] requires 1 argument\n"
var bootEnvCreateTooManyArgErrorString string = "Error: rscli bootenvs create [json] requires 1 argument\n"
var bootEnvCreateBadJSONString = "asdgasdg"
var bootEnvCreateBadJSONErrorString = "Error: Invalid bootenv object: error unmarshaling JSON: json: cannot unmarshal string into Go value of type models.BootEnv\n\n"
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
  "OptionalParams": null,
  "RequiredParams": null,
  "Templates": null
}
`
var bootEnvCreateDuplicateErrorString = "Error: dataTracker create bootenvs: john already exists\n\n"

var bootEnvListBothEnvsString = `[
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
    "OptionalParams": null,
    "RequiredParams": null,
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
    ]
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
    "OptionalParams": null,
    "RequiredParams": null,
    "Templates": null
  }
]
`

var bootEnvUpdateNoArgErrorString string = "Error: rscli bootenvs update [id] [json] requires 2 arguments"
var bootEnvUpdateTooManyArgErrorString string = "Error: rscli bootenvs update [id] [json] requires 2 arguments"
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
  "OptionalParams": null,
  "RequiredParams": null,
  "Templates": null
}
`
var bootEnvUpdateJohnMissingErrorString string = "Error: bootenvs GET: john2: Not Found\n\n"

var bootEnvPatchNoArgErrorString string = "Error: rscli bootenvs patch [objectJson] [changesJson] requires 2 arguments"
var bootEnvPatchTooManyArgErrorString string = "Error: rscli bootenvs patch [objectJson] [changesJson] requires 2 arguments"
var bootEnvPatchBadPatchJSONString = "asdgasdg"
var bootEnvPatchBadPatchJSONErrorString = "Error: Unable to parse rscli bootenvs patch [objectJson] [changesJson] JSON asdgasdg\nError: error unmarshaling JSON: json: cannot unmarshal string into Go value of type models.BootEnv\n\n"
var bootEnvPatchBadBaseJSONString = "asdgasdg"
var bootEnvPatchBadBaseJSONErrorString = "Error: Unable to parse rscli bootenvs patch [objectJson] [changesJson] JSON asdgasdg\nError: error unmarshaling JSON: json: cannot unmarshal string into Go value of type models.BootEnv\n\n"
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
  "OptionalParams": null,
  "RequiredParams": null,
  "Templates": null
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
  "OptionalParams": null,
  "RequiredParams": null,
  "Templates": null
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
  "OptionalParams": null,
  "RequiredParams": null,
  "Templates": null
}
`
var bootEnvPatchJohnMissingErrorString string = "Error: bootenvs: PATCH john2: Not Found\n\n"

var bootEnvDestroyNoArgErrorString string = "Error: rscli bootenvs destroy [id] requires 1 argument"
var bootEnvDestroyTooManyArgErrorString string = "Error: rscli bootenvs destroy [id] requires 1 argument"
var bootEnvDestroyJohnString string = "Deleted bootenv john\n"
var bootEnvDestroyMissingJohnString string = "Error: bootenvs: DELETE john: Not Found\n\n"

var bootEnvInstallNoArgUsageString string = "Error: rscli bootenvs install [bootenvFile] [isoPath] needs at least 1 arg\n"
var bootEnvInstallTooManyArgUsageString string = "Error: rscli bootenvs install [bootenvFile] [isoPath] has Too many args\n"
var bootEnvInstallBadBootEnvDirErrorString string = "Error: Error determining whether bootenvs dir exists: stat bootenvs: no such file or directory\n\n"
var bootEnvInstallBootEnvDirIsFileErrorString string = "Error: bootenvs is not a directory\n\n"
var bootEnvInstallNoSledgehammerErrorString string = "Error: No bootenv bootenvs/sledgehammer.yml\n\n"
var bootEnvInstallSledgehammerBadJsonErrorString string = "Error: Invalid bootenv object: error unmarshaling JSON: json: cannot unmarshal string into Go value of type models.BootEnv\n\n\n"

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

var bootEnvInstallSledgehammerSuccessString string = `{
  "Available": true,
  "BootParams": "rootflags=loop root=live:/sledgehammer.iso rootfstype=auto ro liveimg rd_NO_LUKS rd_NO_MD rd_NO_DM provisioner.web={{.ProvisionerURL}} rebar.web={{.CommandURL}} rs.uuid={{.Machine.UUID}} rs.api={{.ApiURL}}",
  "Errors": null,
  "Initrds": [
    "stage1.img"
  ],
  "Kernel": "vmlinuz0",
  "Name": "sledgehammer",
  "OS": {
    "IsoFile": "sledgehammer-708de8b878e3818b1c1bb598a56de968939f9d4b.tar",
    "IsoSha256": "e094e066b24671c461c17482c6f071d78723275c48df70eb9d24125c89e99760",
    "IsoUrl": "http://127.0.0.1:10003/sledgehammer-708de8b878e3818b1c1bb598a56de968939f9d4b.tar",
    "Name": "sledgehammer/708de8b878e3818b1c1bb598a56de968939f9d4b"
  },
  "OptionalParams": [
    "ntp_servers",
    "access_keys"
  ],
  "RequiredParams": null,
  "Templates": [
    {
      "Contents": "DEFAULT discovery\nPROMPT 0\nTIMEOUT 10\nLABEL discovery\n  KERNEL {{.Env.PathFor \"tftp\" .Env.Kernel}}\n  INITRD {{.Env.JoinInitrds \"tftp\"}}\n  APPEND {{.BootParams}}\n  IPAPPEND 2\n",
      "Name": "pxelinux",
      "Path": "pxelinux.cfg/{{.Machine.HexAddress}}"
    },
    {
      "Contents": "delay=2\ntimeout=20\nverbose=5\nimage={{.Env.PathFor \"tftp\" .Env.Kernel}}\ninitrd={{.Env.JoinInitrds \"tftp\"}}\nappend={{.BootParams}}\n",
      "Name": "elilo",
      "Path": "{{.Machine.HexAddress}}.conf"
    },
    {
      "Contents": "#!ipxe\nkernel {{.Env.PathFor \"http\" .Env.Kernel}} {{.BootParams}} BOOTIF=01-${netX/mac:hexhyp}\n{{ range $initrd := .Env.Initrds }}\ninitrd {{$.Env.PathFor \"http\" $initrd}}\n{{ end }}\nboot\n",
      "Name": "ipxe",
      "Path": "{{.Machine.Address}}.ipxe"
    },
    {
      "Contents": "#!/bin/bash\n# Copyright 2017, RackN\n#\n# Licensed under the Apache License, Version 2.0 (the \"License\");\n# you may not use this file except in compliance with the License.\n# You may obtain a copy of the License at\n#\n#  http://www.apache.org/licenses/LICENSE-2.0\n#\n# Unless required by applicable law or agreed to in writing, software\n# distributed under the License is distributed on an \"AS IS\" BASIS,\n# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.\n# See the License for the specific language governing permissions and\n# limitations under the License.\n#\n\n# We get the following variables from start-up.sh\n# MAC BOOTDEV ADMIN_IP DOMAIN HOSTNAME HOSTNAME_MAC MYIP\n\nset -x\nshopt -s extglob\nexport PS4=\"${BASH_SOURCE}@${LINENO}(${FUNCNAME[0]}): \"\ncp /usr/share/zoneinfo/GMT /etc/localtime\n\n# Set up just enough infrastructure to let the jigs work.\n# Allow client to pass http proxy environment variables\necho \"AcceptEnv http_proxy https_proxy no_proxy\" \u003e\u003e /etc/ssh/sshd_config\nservice sshd restart\n\n# Synchronize our date\n{{ if (.ParamExists \"ntp_servers\") }}\nntpdate \"{{index (.Param \"ntp_servers\") 0}}\"\n{{ end }}\n\n{{ if (.ParamExists \"access_keys\") }}\nmkdir -p /root/.ssh\ncat \u003e/root/.ssh/authorized_keys \u003c\u003cEOF\n### BEGIN GENERATED CONTENT\n{{ range $key := .Param \"access_keys\" }}{{$key}}{{ end }}\n#### END GENERATED CONTENT\nEOF\n{{ end }}\n\n# The last line in this script must always be exit 0!!\nexit 0\n",
      "Name": "control.sh",
      "Path": "{{.Machine.Path}}/control.sh"
    }
  ]
}
`

var bootEnvInstallLocalMissingTemplatesErrorString string = "Error: local requires template local-pxelinux.tmpl, but we cannot find it in templates/local-pxelinux.tmpl\n\n"

var bootEnvInstallLocalSuccessString string = `{
  "Available": true,
  "BootParams": "",
  "Errors": null,
  "Initrds": null,
  "Kernel": "",
  "Name": "local",
  "OS": {
    "Name": "local"
  },
  "OptionalParams": null,
  "RequiredParams": null,
  "Templates": [
    {
      "ID": "local-pxelinux.tmpl",
      "Name": "pxelinux",
      "Path": "pxelinux.cfg/{{.Machine.HexAddress}}"
    },
    {
      "ID": "local-elilo.tmpl",
      "Name": "elilo",
      "Path": "{{.Machine.HexAddress}}.conf"
    },
    {
      "ID": "local-ipxe.tmpl",
      "Name": "ipxe",
      "Path": "{{.Machine.Address}}.ipxe"
    }
  ]
}
`

func TestBootEnvCli(t *testing.T) {
	tests := []CliTest{
		CliTest{true, false, []string{"bootenvs"}, noStdinString, "Access CLI commands relating to bootenvs\n", ""},
		CliTest{false, false, []string{"bootenvs", "list"}, noStdinString, bootEnvDefaultListString, noErrorString},

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
		CliTest{false, true, []string{"bootenvs", "install", "sledgehammer"}, noStdinString, noContentString, bootEnvInstallBadBootEnvDirErrorString},
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
		CliTest{false, true, []string{"bootenvs", "install", "bootenvs/sledgehammer.yml"}, noStdinString, noContentString, bootEnvInstallBootEnvDirIsFileErrorString},
	}
	for _, test := range tests {
		testCli(t, test)
	}

	os.RemoveAll("bootenvs")
	if err := os.MkdirAll("bootenvs", 0755); err != nil {
		t.Errorf("Failed to create bootenvs dir: %v\n", err)
	}

	tests = []CliTest{
		CliTest{false, true, []string{"bootenvs", "install", "bootenvs/sledgehammer.yml"}, noStdinString, noContentString, bootEnvInstallNoSledgehammerErrorString},
	}
	for _, test := range tests {
		testCli(t, test)
	}

	if err := ioutil.WriteFile("bootenvs/sledgehammer.yml", []byte("TEST"), 0644); err != nil {
		t.Errorf("Failed to create bootenvs file: %v\n", err)
	}

	tests = []CliTest{
		CliTest{false, true, []string{"bootenvs", "install", "bootenvs/sledgehammer.yml"}, noStdinString, noContentString, bootEnvInstallSledgehammerBadJsonErrorString},
	}
	for _, test := range tests {
		testCli(t, test)
	}

	midlayer.ServeStatic("127.0.0.1:10003", backend.NewFS("test-data", nil), nil)

	os.RemoveAll("bootenvs/sledgehammer.yml")
	if err := os.MkdirAll("bootenvs", 0755); err != nil {
		t.Errorf("Failed to create bootenvs dir: %v\n", err)
	}
	if err := os.Symlink("../test-data/sledgehammer.yml", "bootenvs/sledgehammer.yml"); err != nil {
		t.Errorf("Failed to create link to sledgehammer.yml: %v\n", err)
	}
	if err := os.Symlink("../../assets/bootenvs/local.yml", "bootenvs/local.yml"); err != nil {
		t.Errorf("Failed to create link to local.yml: %v\n", err)
	}
	tests = []CliTest{
		CliTest{false, false, []string{"bootenvs", "install", "--skip-download", "bootenvs/sledgehammer.yml"}, noStdinString, bootEnvInstallSledgehammerSuccessWithErrorsString, noErrorString},
		CliTest{false, false, []string{"bootenvs", "destroy", "sledgehammer"}, noStdinString, "Deleted bootenv sledgehammer\n", noErrorString},
	}
	for _, test := range tests {
		testCli(t, test)
	}

	installSkipDownloadIsos = false
	tests = []CliTest{
		CliTest{false, false, []string{"bootenvs", "install", "bootenvs/sledgehammer.yml"}, noStdinString, bootEnvInstallSledgehammerSuccessString, noErrorString},
		CliTest{false, true, []string{"bootenvs", "install", "bootenvs/local.yml"}, noStdinString, noContentString, bootEnvInstallLocalMissingTemplatesErrorString},
	}
	for _, test := range tests {
		testCli(t, test)
	}

	if err := os.MkdirAll("templates", 0755); err != nil {
		t.Errorf("Failed to create templates dir: %v\n", err)
	}
	tmpls := []string{"local-pxelinux.tmpl", "local-elilo.tmpl", "local-ipxe.tmpl"}
	for _, tmpl := range tmpls {
		if err := os.Symlink("../../assets/templates/"+tmpl, "templates/"+tmpl); err != nil {
			t.Errorf("Failed to create link to %s: %v\n", tmpl, err)
		}
	}
	tests = []CliTest{
		CliTest{false, false, []string{"bootenvs", "install", "bootenvs/local.yml", "ic"}, noStdinString, bootEnvInstallLocalSuccessString, noErrorString},
		CliTest{false, false, []string{"bootenvs", "destroy", "sledgehammer"}, noStdinString, "Deleted bootenv sledgehammer\n", noErrorString},
		CliTest{false, false, []string{"bootenvs", "install", "bootenvs/sledgehammer.yml"}, noStdinString, bootEnvInstallSledgehammerSuccessString, noErrorString},

		// Clean up
		CliTest{false, false, []string{"bootenvs", "destroy", "sledgehammer"}, noStdinString, "Deleted bootenv sledgehammer\n", noErrorString},
		CliTest{false, false, []string{"bootenvs", "destroy", "local"}, noStdinString, "Deleted bootenv local\n", noErrorString},
		CliTest{false, false, []string{"templates", "destroy", "local-pxelinux.tmpl"}, noStdinString, "Deleted template local-pxelinux.tmpl\n", noErrorString},
		CliTest{false, false, []string{"templates", "destroy", "local-elilo.tmpl"}, noStdinString, "Deleted template local-elilo.tmpl\n", noErrorString},
		CliTest{false, false, []string{"templates", "destroy", "local-ipxe.tmpl"}, noStdinString, "Deleted template local-ipxe.tmpl\n", noErrorString},
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
