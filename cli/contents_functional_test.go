package cli

import (
	"testing"
)

func TestContentsFunctionalCli(t *testing.T) {
	var (
		contentMyLocalBootEnvString = `{
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
		contentPackBadString = `{
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
		contentPack1String = `{
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
		contentPack2String = `{
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
		contentMachineCreateString = `{
  "Name": "greg",
  "BootEnv": "mylocal",
  "Secret": "secret1",
  "Uuid": "3e7031fe-3062-45f1-835c-92541bc9cbd3"
}
`
		contentPack1BadUpdateString = `{
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
		contentPack1BadSyntaxUpdateString = `{
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
		contentPack1UpdateString = `{
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
	)
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
