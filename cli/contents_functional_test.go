package cli

import (
	"testing"
)

var contentMyLocalBootEnvString = `{
  "BootParams": "",
  "Kernel": "",
  "Name": "mylocal",
  "OS": {
    "Name": "mylocal"
  },
  "OnlyUnknown": false,
  "ReadOnly": false,
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

var contentPackBadString = `{
  "meta": {
    "Name": "PackBad",
    "Version": "0.1",
  },
  "sections": {
    "profiles": {
      "p1-bad": {
	"Description": "packbad",
	"Params": 12
      }
    }
  }
}
`

var contentPackBadCreateErrorString = "Error: Failed to load backing objects from cache: Unable to load profiles: error unmarshaling JSON: json: cannot unmarshal number into Go struct field Profile.Params of type map[string]interface {}\n\n"

var contentPack1String = `{
  "meta": {
    "Name": "Pack1",
    "Version": "0.1",
  },
  "sections": {
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
  "Warnings": [],
  "meta": {
    "Name": "Pack1",
    "Type": "dynamic",
    "Version": "0.1"
  }
}
`

var contentPack2String = `{
  "meta": {
    "Name": "Pack2",
    "Version": "0.2",
  },
  "sections": {
    "profiles": {
      "p1-prof": {
	"Description": "pack2",
        "Name": "p1-prof"
      }
    }
  }
}
`

var contentPack2CreateErrorString = "Error: ValidationError: New layer violates key restrictions: keysCannotBeOverridden: p1-prof is already in layer 1\n\n"

var contentPack1ProfileListString = `[
  {
    "Available": true,
    "Errors": [],
    "Meta": {
      "color": "blue",
      "icon": "world",
      "title": "Digital Rebar Provision"
    },
    "Name": "global",
    "ReadOnly": false,
    "Validated": true
  },
  {
    "Available": true,
    "Description": "pack1",
    "Errors": [],
    "Name": "p1-prof",
    "ReadOnly": true,
    "Validated": true
  }
]
`
var contentPack1ProfileList2String = `[
  {
    "Available": true,
    "Errors": [],
    "Meta": {
      "color": "blue",
      "icon": "world",
      "title": "Digital Rebar Provision"
    },
    "Name": "global",
    "ReadOnly": false,
    "Validated": true
  },
  {
    "Available": true,
    "Description": "pack1-2",
    "Errors": [],
    "Name": "p2-prof",
    "ReadOnly": true,
    "Validated": true
  }
]
`
var contentPack1UpdateProfileListString = `[
  {
    "Available": true,
    "Errors": [],
    "Meta": {
      "color": "blue",
      "icon": "world",
      "title": "Digital Rebar Provision"
    },
    "Name": "global",
    "ReadOnly": false,
    "Validated": true
  },
  {
    "Available": true,
    "Description": "pack1-2",
    "Errors": [],
    "Name": "p1-prof",
    "ReadOnly": true,
    "Validated": true
  }
]
`

var contentNoPackProfileListString = `[
  {
    "Available": true,
    "Errors": [],
    "Meta": {
      "color": "blue",
      "icon": "world",
      "title": "Digital Rebar Provision"
    },
    "Name": "global",
    "ReadOnly": false,
    "Validated": true
  }
]
`

var contentMachineCreateString = `{
  "Name": "greg",
  "BootEnv": "mylocal",
  "Secret": "secret1",
  "Uuid": "3e7031fe-3062-45f1-835c-92541bc9cbd3"
}
`
var contentMachineCreateSuccessString = `{
  "Available": true,
  "BootEnv": "mylocal",
  "CurrentTask": 0,
  "Errors": [],
  "Name": "greg",
  "Profile": {
    "Available": false,
    "Errors": null,
    "Name": "",
    "ReadOnly": false,
    "Validated": false
  },
  "Profiles": [],
  "ReadOnly": false,
  "Runnable": true,
  "Secret": "secret1",
  "Stage": "none",
  "Tasks": [],
  "Uuid": "3e7031fe-3062-45f1-835c-92541bc9cbd3",
  "Validated": true
}
`

var contentMachineAddProfileString = `{
  "Available": true,
  "BootEnv": "mylocal",
  "CurrentTask": 0,
  "Errors": [],
  "Name": "greg",
  "Profile": {
    "Available": false,
    "Errors": null,
    "Name": "",
    "ReadOnly": false,
    "Validated": false
  },
  "Profiles": [
    "p1-prof"
  ],
  "ReadOnly": false,
  "Runnable": true,
  "Secret": "secret1",
  "Stage": "none",
  "Tasks": [],
  "Uuid": "3e7031fe-3062-45f1-835c-92541bc9cbd3",
  "Validated": true
}
`

var contentBootenvGregCreateSuccessString = `{
  "Available": true,
  "BootParams": "",
  "Errors": [],
  "Initrds": [],
  "Kernel": "",
  "Name": "mylocal",
  "OS": {
    "Name": "mylocal"
  },
  "OnlyUnknown": false,
  "OptionalParams": [],
  "ReadOnly": false,
  "RequiredParams": [],
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

var contentPack1BadUpdateString = `{
  "meta": {
    "Name": "Pack1",
    "Version": "0.2",
  },
  "sections": {
    "profiles": {
      "p2-prof": {
	"Description": "pack1-2",
        "Name": "p2-prof"
      }
    }
  }
}
`
var contentPack1BadUpdateSuccessString = `{
  "Counts": {
    "profiles": 1
  },
  "Warnings": [
    "Profile p1-prof (at 0) does not exist"
  ],
  "meta": {
    "Name": "Pack1",
    "Type": "dynamic",
    "Version": "0.2"
  }
}
`

var contentPack1BadSyntaxUpdateString = `{
  "meta": {
    "Name": "Pack1",
    "Version": "0.2",
  },
  "sections": {
    "profiles": {
      "p2-prof": {
	"Description": "pack1-2",
        "Name": "p2-prof",
	"Params": 12
      }
    }
  }
}
`
var contentPack1BadSyntaxUpdateErrorString = "Error: Unable to load profiles: error unmarshaling JSON: json: cannot unmarshal number into Go struct field Profile.Params of type map[string]interface {}\n\n"

var contentPack1UpdateString = `{
  "meta": {
    "Name": "Pack1",
    "Version": "0.3",
  },
  "sections": {
    "profiles": {
      "p1-prof": {
	"Description": "pack1-2",
        "Name": "p1-prof"
      }
    }
  }
}
`
var contentPack1UpdateSuccessString = `{
  "Counts": {
    "profiles": 1
  },
  "Warnings": [],
  "meta": {
    "Name": "Pack1",
    "Type": "dynamic",
    "Version": "0.3"
  }
}
`

func TestContentsFunctionalCli(t *testing.T) {

	tests := []CliTest{
		CliTest{false, false, []string{"contents", "list"}, noStdinString, contentDefaultListString, noErrorString},
		CliTest{false, false, []string{"bootenvs", "create", contentMyLocalBootEnvString}, noStdinString, contentBootenvGregCreateSuccessString, noErrorString},

		CliTest{false, true, []string{"contents", "create", contentPackBadString}, noStdinString, noContentString, contentPackBadCreateErrorString},
		CliTest{false, false, []string{"contents", "create", contentPack1String}, noStdinString, contentPack1CreateSuccessString, noErrorString},
		CliTest{false, false, []string{"profiles", "list"}, noStdinString, contentPack1ProfileListString, noErrorString},
		CliTest{false, true, []string{"contents", "create", contentPack2String}, noStdinString, noContentString, contentPack2CreateErrorString},
		CliTest{false, false, []string{"profiles", "list"}, noStdinString, contentPack1ProfileListString, noErrorString},

		CliTest{false, false, []string{"machines", "create", contentMachineCreateString}, noStdinString, contentMachineCreateSuccessString, noErrorString},
		CliTest{false, false, []string{"machines", "addprofile", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "p1-prof"}, noStdinString, contentMachineAddProfileString, noErrorString},

		CliTest{false, true, []string{"contents", "update", "Pack1", contentPack1BadSyntaxUpdateString}, noStdinString, noContentString, contentPack1BadSyntaxUpdateErrorString},
		CliTest{false, false, []string{"contents", "update", "Pack1", contentPack1BadUpdateString}, noStdinString, contentPack1BadUpdateSuccessString, noErrorString},
		CliTest{false, false, []string{"profiles", "list"}, noStdinString, contentPack1ProfileList2String, noErrorString},
		CliTest{false, false, []string{"contents", "update", "Pack1", contentPack1UpdateString}, noStdinString, contentPack1UpdateSuccessString, noErrorString},
		CliTest{false, false, []string{"profiles", "list"}, noStdinString, contentPack1UpdateProfileListString, noErrorString},

		CliTest{false, false, []string{"machines", "destroy", "3e7031fe-3062-45f1-835c-92541bc9cbd3"}, noStdinString, "Deleted machine 3e7031fe-3062-45f1-835c-92541bc9cbd3\n", noErrorString},

		CliTest{false, false, []string{"contents", "destroy", "Pack1"}, noStdinString, "Deleted content Pack1\n", noErrorString},
		CliTest{false, false, []string{"profiles", "list"}, noStdinString, contentNoPackProfileListString, noErrorString},

		CliTest{false, false, []string{"bootenvs", "destroy", "mylocal"}, noStdinString, "Deleted bootenv mylocal\n", noErrorString},
		CliTest{false, false, []string{"contents", "list"}, noStdinString, contentDefaultListString, noErrorString},
	}

	for _, test := range tests {
		testCli(t, test)
	}

}
