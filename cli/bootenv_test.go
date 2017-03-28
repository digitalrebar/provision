package cli

import (
	"testing"
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
	}

	for _, test := range tests {
		testCli(t, test)
	}

}

// TODO: Test Install bootenv.
