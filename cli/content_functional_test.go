package cli

import (
	"testing"
)

var contentMyLocalBootEnvString = `{
  "BootParams": "",
  "Kernel": "",
  "Name": "local",
  "OS": {
    "Name": "local"
  },
  "OnlyUnknown": false,
  "Templates": [
    {
      "Contents": "local-pxelinux.tmpl",
      "Name": "pxelinux",
      "Path": "pxelinux.cfg/{{.Machine.HexAddress}}"
    },
    {
      "Contents": "local-elilo.tmpl",
      "Name": "elilo",
      "Path": "{{.Machine.HexAddress}}.conf"
    },
    {
      "Contents": "local-ipxe.tmpl",
      "Name": "ipxe",
      "Path": "{{.Machine.Address}}.ipxe"
    }
  ]
}
`

var contentPack1String = `{
  "Name": "Pack1",
  "Version": "0.1",
  "Sections": {
    "profiles": {
      "p1-prof": {
	"Description": "pack1",
        "Name": "p1-prof"
      }
    }
  }
}
`

var contentPack1CreateSuccessString = `{
  "Counts": {
    "profiles": 1
  },
  "Name": "Pack1",
  "Version": "0.1"
}
`

var contentPack2String = `{
  "Name": "Pack2",
  "Version": "0.2",
  "Sections": {
    "profiles": {
      "p1-prof": {
	"Description": "pack2",
        "Name": "p1-prof"
      }
    }
  }
}
`

var contentPack2CreateErrorString = "Error: New layer violates key restrictions: keysCannotBeOverridden: p1-prof is already in layer 1\n\n"

var contentPack1ProfileListString = `[
  {
    "Name": "global",
    "Tasks": null
  },
  {
    "Description": "pack1",
    "Name": "p1-prof",
    "Tasks": null
  }
]
`

var contentNoPackProfileListString = `[
  {
    "Name": "global",
    "Tasks": null
  }
]
`

var contentMachineCreateString = `{
  "Name": "greg",
  "BootEnv": "local",
  "Uuid": "3e7031fe-3062-45f1-835c-92541bc9cbd3"
}
`
var contentMachineCreateSuccessString = `{
  "BootEnv": "local",
  "CurrentTask": 0,
  "Errors": null,
  "Name": "greg",
  "Profile": {
    "Name": "",
    "Tasks": null
  },
  "Profiles": null,
  "Runnable": true,
  "Tasks": [],
  "Uuid": "3e7031fe-3062-45f1-835c-92541bc9cbd3"
}
`

var contentMachineAddProfileString = `{
  "BootEnv": "local",
  "CurrentTask": 0,
  "Errors": null,
  "Name": "greg",
  "Profile": {
    "Name": "",
    "Tasks": null
  },
  "Profiles": [
    "p1-prof"
  ],
  "Runnable": true,
  "Tasks": [],
  "Uuid": "3e7031fe-3062-45f1-835c-92541bc9cbd3"
}
`

var contentBootenvGregCreateSuccessString = `{
  "Available": true,
  "BootParams": "",
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
      "Contents": "local-pxelinux.tmpl",
      "Name": "pxelinux",
      "Path": "pxelinux.cfg/{{.Machine.HexAddress}}"
    },
    {
      "Contents": "local-elilo.tmpl",
      "Name": "elilo",
      "Path": "{{.Machine.HexAddress}}.conf"
    },
    {
      "Contents": "local-ipxe.tmpl",
      "Name": "ipxe",
      "Path": "{{.Machine.Address}}.ipxe"
    }
  ],
  "Validated": true
}
`

var contentPack1DestroyErrorString = "Error: Profile p1-prof (at 0) does not exist\n\n"

func TestContentFunctionalCli(t *testing.T) {

	tests := []CliTest{
		CliTest{false, false, []string{"contents", "list"}, noStdinString, contentDefaultListString, noErrorString},
		CliTest{false, false, []string{"bootenvs", "create", contentMyLocalBootEnvString}, noStdinString, contentBootenvGregCreateSuccessString, noErrorString},

		CliTest{false, false, []string{"contents", "create", contentPack1String}, noStdinString, contentPack1CreateSuccessString, noErrorString},
		CliTest{false, false, []string{"profiles", "list"}, noStdinString, contentPack1ProfileListString, noErrorString},
		CliTest{false, true, []string{"contents", "create", contentPack2String}, noStdinString, noContentString, contentPack2CreateErrorString},
		CliTest{false, false, []string{"profiles", "list"}, noStdinString, contentPack1ProfileListString, noErrorString},

		CliTest{false, false, []string{"machines", "create", contentMachineCreateString}, noStdinString, contentMachineCreateSuccessString, noErrorString},
		CliTest{false, false, []string{"machines", "addprofile", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "p1-prof"}, noStdinString, contentMachineAddProfileString, noErrorString},

		CliTest{false, true, []string{"contents", "destroy", "Pack1"}, noStdinString, noContentString, contentPack1DestroyErrorString},
		CliTest{false, false, []string{"profiles", "list"}, noStdinString, contentPack1ProfileListString, noErrorString},

		CliTest{false, false, []string{"machines", "destroy", "3e7031fe-3062-45f1-835c-92541bc9cbd3"}, noStdinString, "Deleted machine 3e7031fe-3062-45f1-835c-92541bc9cbd3\n", noErrorString},

		CliTest{false, false, []string{"contents", "destroy", "Pack1"}, noStdinString, "Deleted content Pack1\n", noErrorString},
		CliTest{false, false, []string{"profiles", "list"}, noStdinString, contentNoPackProfileListString, noErrorString},

		CliTest{false, false, []string{"bootenvs", "destroy", "local"}, noStdinString, "Deleted bootenv local\n", noErrorString},
		CliTest{false, false, []string{"contents", "list"}, noStdinString, contentDefaultListString, noErrorString},
	}

	for _, test := range tests {
		testCli(t, test)
	}

}
