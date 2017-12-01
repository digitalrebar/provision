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
    "Description": "Global profile attached automatically to all machines.",
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
    "Description": "Global profile attached automatically to all machines.",
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
    "Description": "Global profile attached automatically to all machines.",
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
    "Description": "Global profile attached automatically to all machines.",
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
    "Errors": [],
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
    "Errors": [],
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
var contentPack1BadSyntaxUpdateErrorString = `Error: PUT: contents/Pack1
  Unable to load profiles: error unmarshaling JSON: json: cannot unmarshal number into Go struct field Profile.Params of type map[string]interface {}
  Profile p1-prof (at 0) does not exist

`

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

	cliTest(false, false, "contents", "list").run(t)
	cliTest(false, false, "bootenvs", "create", contentMyLocalBootEnvString).run(t)

	cliTest(false, true, "contents", "create", contentPackBadString).run(t)
	cliTest(false, false, "contents", "create", contentPack1String).run(t)
	cliTest(false, false, "profiles", "list").run(t)
	cliTest(false, true, "contents", "create", contentPack2String).run(t)
	cliTest(false, false, "profiles", "list").run(t)

	cliTest(false, false, "machines", "create", contentMachineCreateString).run(t)
	cliTest(false, false, "machines", "addprofile", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "p1-prof").run(t)

	cliTest(false, true, "contents", "update", "Pack1", contentPack1BadSyntaxUpdateString).run(t)
	cliTest(false, false, "contents", "update", "Pack1", contentPack1BadUpdateString).run(t)
	cliTest(false, false, "profiles", "list").run(t)
	cliTest(false, false, "contents", "update", "Pack1", contentPack1UpdateString).run(t)
	cliTest(false, false, "profiles", "list").run(t)

	cliTest(false, false, "machines", "destroy", "3e7031fe-3062-45f1-835c-92541bc9cbd3").run(t)

	cliTest(false, false, "contents", "destroy", "Pack1").run(t)
	cliTest(false, false, "profiles", "list").run(t)

	cliTest(false, false, "bootenvs", "destroy", "mylocal").run(t)
	cliTest(false, false, "contents", "list").run(t)

}
