package cli

import (
	"testing"
)

var subnetAddrErrorString string = "Error: GET: subnets: Invalid Address: fred\n\n"
var subnetExpireTimeErrorString string = "Error: GET: subnets: Invalid subnet CIDR: false\n\n"
var subnetShowMissingArgErrorString string = "Error: GET: subnets/ignore: Not Found\n\n"
var subnetExistsMissingIgnoreString string = "Error: GET: subnets/ignore: Not Found\n\n"
var subnetCreateBadJSONErrorString = "Error: Unable to create new subnet: Invalid type passed to subnet create\n\n"
var subnetCreateDuplicateErrorString = "Error: CREATE: subnets/john: already exists\n\n"
var subnetUpdateJohnMissingErrorString string = "Error: GET: subnets/john2: Not Found\n\n"
var subnetPatchJohnMissingErrorString string = "Error: PATCH: subnets/john2: Not Found\n\n"
var subnetDestroyMissingJohnString string = "Error: DELETE: subnets/john: Not Found\n\n"

var subnetInvalidEnabledBooleanListString = "Error: GET: subnets: Enabled must be true or false\n\n"
var subnetInvalidProxyBooleanListString = "Error: GET: subnets: Proxy must be true or false\n\n"
var subnetRangeIPFailureString string = "Error: PATCH: subnets/john: invalid IP address: cq.98.42.1234\n\n"
var subnetRangeIPBadIpString string = "Error: PATCH: subnets/john: invalid IP address: 192.168.100.500\n\n"
var subnetSubnetCIDRFailureString = "Error: 1111.11.2223.544/66666 is not a valid subnet CIDR\n\n"
var subnetStrategyMacFailureErrorString string = "Error: t5:44:llll:b is not a valid strategy\n\n"
var subnetLeasetimesIntFailureString string = "Error: 4x5 could not be read as a number\n\n"
var subnetSetIntFailureErrorString string = "Error: 6tl could not be read as a number\n\n"
var subnetGetToNull string = "Error: option 6 does not exist\n\n"

var subnetDefaultListString string = "[]\n"
var subnetEmptyListString string = "[]\n"

var subnetShowNoArgErrorString string = "Error: drpcli subnets show [id] [flags] requires 1 argument\n"
var subnetShowTooManyArgErrorString string = "Error: drpcli subnets show [id] [flags] requires 1 argument\n"

var subnetShowJohnString string = `{
  "ActiveEnd": "192.168.100.100",
  "ActiveLeaseTime": 60,
  "ActiveStart": "192.168.100.20",
  "Available": true,
  "Enabled": false,
  "Errors": [],
  "Name": "john",
  "NextServer": "3.3.3.3",
  "OnlyReservations": false,
  "Options": [
    {
      "Code": 1,
      "Value": "255.255.255.0"
    },
    {
      "Code": 28,
      "Value": "192.168.100.255"
    }
  ],
  "Pickers": [
    "hint",
    "nextFree",
    "mostExpired"
  ],
  "Proxy": false,
  "ReadOnly": false,
  "ReservedLeaseTime": 7200,
  "Strategy": "MAC",
  "Subnet": "192.168.100.0/24",
  "Validated": true
}
`

var subnetExistsNoArgErrorString string = "Error: drpcli subnets exists [id] [flags] requires 1 argument"
var subnetExistsTooManyArgErrorString string = "Error: drpcli subnets exists [id] [flags] requires 1 argument"
var subnetExistsIgnoreString string = ""

var subnetCreateNoArgErrorString string = "Error: drpcli subnets create [json] [flags] requires 1 argument\n"
var subnetCreateTooManyArgErrorString string = "Error: drpcli subnets create [json] [flags] requires 1 argument\n"
var subnetCreateBadJSONString = "asdgasdg"

var subnetCreateInputString string = `{
  "Name": "john",
  "ActiveEnd": "192.168.100.100",
  "ActiveStart": "192.168.100.20",
  "ActiveLeaseTime": 60,
  "NextServer": "3.3.3.3",
  "OnlyReservations": false,
  "ReservedLeaseTime": 7200,
  "Subnet": "192.168.100.0/24",
  "Strategy": "MAC"
}
`
var subnetCreateJohnString string = `{
  "ActiveEnd": "192.168.100.100",
  "ActiveLeaseTime": 60,
  "ActiveStart": "192.168.100.20",
  "Available": true,
  "Enabled": false,
  "Errors": [],
  "Name": "john",
  "NextServer": "3.3.3.3",
  "OnlyReservations": false,
  "Options": [
    {
      "Code": 1,
      "Value": "255.255.255.0"
    },
    {
      "Code": 28,
      "Value": "192.168.100.255"
    }
  ],
  "Pickers": [
    "hint",
    "nextFree",
    "mostExpired"
  ],
  "Proxy": false,
  "ReadOnly": false,
  "ReservedLeaseTime": 7200,
  "Strategy": "MAC",
  "Subnet": "192.168.100.0/24",
  "Validated": true
}
`

var subnetListBothEnvsString = `[
  {
    "ActiveEnd": "192.168.100.100",
    "ActiveLeaseTime": 60,
    "ActiveStart": "192.168.100.20",
    "Available": true,
    "Enabled": false,
    "Errors": [],
    "Name": "john",
    "NextServer": "3.3.3.3",
    "OnlyReservations": false,
    "Options": [
      {
        "Code": 1,
        "Value": "255.255.255.0"
      },
      {
        "Code": 28,
        "Value": "192.168.100.255"
      }
    ],
    "Pickers": [
      "hint",
      "nextFree",
      "mostExpired"
    ],
    "Proxy": false,
    "ReadOnly": false,
    "ReservedLeaseTime": 7200,
    "Strategy": "MAC",
    "Subnet": "192.168.100.0/24",
    "Validated": true
  }
]
`

var subnetUpdateNoArgErrorString string = "Error: drpcli subnets update [id] [json] [flags] requires 2 arguments"
var subnetUpdateTooManyArgErrorString string = "Error: drpcli subnets update [id] [json] [flags] requires 2 arguments"
var subnetUpdateBadJSONString = "asdgasdg"
var subnetUpdateBadJSONErrorString = "Error: Unable to merge objects: json: cannot unmarshal string into Go value of type map[string]interface {}\n\n\n"
var subnetUpdateInputString string = `{
  "Strategy": "NewStrat"
}
`
var subnetUpdateJohnString string = `{
  "ActiveEnd": "192.168.100.100",
  "ActiveLeaseTime": 60,
  "ActiveStart": "192.168.100.20",
  "Available": true,
  "Enabled": false,
  "Errors": [],
  "Name": "john",
  "NextServer": "3.3.3.3",
  "OnlyReservations": false,
  "Options": [
    {
      "Code": 1,
      "Value": "255.255.255.0"
    },
    {
      "Code": 28,
      "Value": "192.168.100.255"
    }
  ],
  "Pickers": [
    "hint",
    "nextFree",
    "mostExpired"
  ],
  "Proxy": false,
  "ReadOnly": false,
  "ReservedLeaseTime": 7200,
  "Strategy": "NewStrat",
  "Subnet": "192.168.100.0/24",
  "Validated": true
}
`

var subnetPatchNoArgErrorString string = "Error: drpcli subnets patch [objectJson] [changesJson] [flags] requires 2 arguments"
var subnetPatchTooManyArgErrorString string = "Error: drpcli subnets patch [objectJson] [changesJson] [flags] requires 2 arguments"
var subnetPatchBadPatchJSONString = "asdgasdg"
var subnetPatchBadPatchJSONErrorString = "Error: Unable to parse drpcli subnets patch [objectJson] [changesJson] [flags] JSON asdgasdg\nError: error unmarshaling JSON: json: cannot unmarshal string into Go value of type genmodels.Subnet\n\n"
var subnetPatchBadBaseJSONString = "asdgasdg"
var subnetPatchBadBaseJSONErrorString = "Error: Unable to parse drpcli subnets patch [objectJson] [changesJson] [flags] JSON asdgasdg\nError: error unmarshaling JSON: json: cannot unmarshal string into Go value of type genmodels.Subnet\n\n"
var subnetPatchBaseString string = `{
  "ActiveEnd": "192.168.100.100",
  "ActiveLeaseTime": 60,
  "ActiveStart": "192.168.100.20",
  "Available": true,
  "Enabled": false,
  "Errors": [],
  "Name": "john",
  "NextServer": "3.3.3.3",
  "OnlyReservations": false,
  "Options": [
    {
      "Code": 1,
      "Value": "255.255.255.0"
    },
    {
      "Code": 28,
      "Value": "192.168.100.255"
    }
  ],
  "Pickers": [
    "hint",
    "nextFree",
    "mostExpired"
  ],
  "Proxy": false,
  "ReadOnly": false,
  "ReservedLeaseTime": 7200,
  "Strategy": "NewStrat",
  "Subnet": "192.168.100.0/24",
  "Validated": true
}
`
var subnetPatchInputString string = `{
  "Strategy": "bootx64.efi"
}
`
var subnetPatchJohnString string = `{
  "ActiveEnd": "192.168.100.100",
  "ActiveLeaseTime": 60,
  "ActiveStart": "192.168.100.20",
  "Available": true,
  "Enabled": false,
  "Errors": [],
  "Name": "john",
  "NextServer": "3.3.3.3",
  "OnlyReservations": false,
  "Options": [
    {
      "Code": 1,
      "Value": "255.255.255.0"
    },
    {
      "Code": 28,
      "Value": "192.168.100.255"
    }
  ],
  "Pickers": [
    "hint",
    "nextFree",
    "mostExpired"
  ],
  "Proxy": false,
  "ReadOnly": false,
  "ReservedLeaseTime": 7200,
  "Strategy": "bootx64.efi",
  "Subnet": "192.168.100.0/24",
  "Validated": true
}
`
var subnetPatchMissingBaseString string = `{
  "ActiveEnd": "192.168.100.100",
  "ActiveLeaseTime": 60,
  "ActiveStart": "192.168.100.20",
  "Enabled": false,
  "Name": "john2",
  "NextServer": "3.3.3.3",
  "OnlyReservations": false,
  "Options": [
    {
      "Code": 1,
      "Value": "255.255.255.0"
    },
    {
      "Code": 28,
      "Value": "192.168.100.255"
    }
  ],
  "Pickers": [
    "hint",
    "nextFree",
    "mostExpired"
  ],
  "ReservedLeaseTime": 7200,
  "Strategy": "bootx64.efi",
  "Subnet": "192.168.100.0/24"
}
`

var subnetDestroyNoArgErrorString string = "Error: drpcli subnets destroy [id] [flags] requires 1 argument"
var subnetDestroyTooManyArgErrorString string = "Error: drpcli subnets destroy [id] [flags] requires 1 argument"
var subnetDestroyJohnString string = "Deleted subnet john\n"

var subnetRangeNoArgErrorString string = "Error: drpcli subnets range [subnetName] [startIP] [endIP] [flags] requires 3 arguments\n"
var subnetRangeTooManyArgErrorString string = "Error: drpcli subnets range [subnetName] [startIP] [endIP] [flags] requires 3 arguments\n"
var subnetRangeIPSuccessString string = `{
  "ActiveEnd": "192.168.100.200",
  "ActiveLeaseTime": 60,
  "ActiveStart": "192.168.100.10",
  "Available": true,
  "Enabled": false,
  "Errors": [],
  "Name": "john",
  "NextServer": "3.3.3.3",
  "OnlyReservations": false,
  "Options": [
    {
      "Code": 1,
      "Value": "255.255.255.0"
    },
    {
      "Code": 28,
      "Value": "192.168.100.255"
    }
  ],
  "Pickers": [
    "hint",
    "nextFree",
    "mostExpired"
  ],
  "Proxy": false,
  "ReadOnly": false,
  "ReservedLeaseTime": 7200,
  "Strategy": "NewStrat",
  "Subnet": "192.168.100.0/24",
  "Validated": true
}
`

var subnetSubnetNoArgErrorString string = "Error: drpcli subnets subnet [subnetName] [subnet CIDR] [flags] requires 2 arguments\n"
var subnetSubnetTooManyArgErrorString string = "Error: drpcli subnets subnet [subnetName] [subnet CIDR] [flags] requires 2 arguments\n"
var subnetSubnetCIDRSuccessString = `{
  "ActiveEnd": "192.168.100.200",
  "ActiveLeaseTime": 60,
  "ActiveStart": "192.168.100.10",
  "Available": true,
  "Enabled": false,
  "Errors": [],
  "Name": "john",
  "NextServer": "3.3.3.3",
  "OnlyReservations": false,
  "Options": [
    {
      "Code": 1,
      "Value": "255.192.0.0"
    },
    {
      "Code": 28,
      "Value": "192.191.255.255"
    }
  ],
  "Pickers": [
    "hint",
    "nextFree",
    "mostExpired"
  ],
  "Proxy": false,
  "ReadOnly": false,
  "ReservedLeaseTime": 7200,
  "Strategy": "NewStrat",
  "Subnet": "192.168.100.0/10",
  "Validated": true
}
`

var subnetStrategyNoArgErrorString string = "Error: drpcli subnets strategy [subnetName] [MAC] [flags] requires 2 arguments\n"
var subnetStrategyTooManyArgErrorString string = "Error: drpcli subnets strategy [subnetName] [MAC] [flags] requires 2 arguments\n"
var subnetStrategyMacSuccessString string = `{
  "ActiveEnd": "192.168.100.200",
  "ActiveLeaseTime": 60,
  "ActiveStart": "192.168.100.10",
  "Available": true,
  "Enabled": false,
  "Errors": [],
  "Name": "john",
  "NextServer": "3.3.3.3",
  "OnlyReservations": false,
  "Options": [
    {
      "Code": 1,
      "Value": "255.192.0.0"
    },
    {
      "Code": 28,
      "Value": "192.191.255.255"
    }
  ],
  "Pickers": [
    "hint",
    "nextFree",
    "mostExpired"
  ],
  "Proxy": false,
  "ReadOnly": false,
  "ReservedLeaseTime": 7200,
  "Strategy": "MAC",
  "Subnet": "192.168.100.0/10",
  "Validated": true
}
`

var subnetPickersNoArgErrorString string = "Error: drpcli subnets pickers [subnetName] [list] [flags] requires 2 arguments\n"
var subnetPickersTooManyArgErrorString string = "Error: drpcli subnets pickers [subnetName] [list] [flags] requires 2 arguments\n"
var subnetPickersSuccessString string = `{
  "ActiveEnd": "192.168.100.200",
  "ActiveLeaseTime": 60,
  "ActiveStart": "192.168.100.10",
  "Available": true,
  "Enabled": false,
  "Errors": [],
  "Name": "john",
  "NextServer": "3.3.3.3",
  "OnlyReservations": false,
  "Options": [
    {
      "Code": 1,
      "Value": "255.192.0.0"
    },
    {
      "Code": 28,
      "Value": "192.191.255.255"
    }
  ],
  "Pickers": [
    "none",
    "nextFree",
    "mostExpired"
  ],
  "Proxy": false,
  "ReadOnly": false,
  "ReservedLeaseTime": 7200,
  "Strategy": "a3:b3:51:66:7e:11",
  "Subnet": "192.168.100.0/10",
  "Validated": true
}
`

var subnetNextserverNoArgErrorString string = "Error: drpcli subnets nextserver [subnetName] [IP] [flags] requires 2 arguments\n"
var subnetNextserverTooManyArgErrorString string = "Error: drpcli subnets nextserver [subnetName] [IP] [flags] requires 2 arguments\n"
var subnetNextserverIPSuccess string = `{
  "ActiveEnd": "192.168.100.200",
  "ActiveLeaseTime": 60,
  "ActiveStart": "192.168.100.10",
  "Available": true,
  "Enabled": false,
  "Errors": [],
  "Name": "john",
  "NextServer": "1.24.36.16",
  "OnlyReservations": false,
  "Options": [
    {
      "Code": 1,
      "Value": "255.192.0.0"
    },
    {
      "Code": 28,
      "Value": "192.191.255.255"
    }
  ],
  "Pickers": [
    "none",
    "nextFree",
    "mostExpired"
  ],
  "Proxy": false,
  "ReadOnly": false,
  "ReservedLeaseTime": 7200,
  "Strategy": "a3:b3:51:66:7e:11",
  "Subnet": "192.168.100.0/10",
  "Validated": true
}
`
var subnetLeasetimesNoArgErrorString string = "Error: drpcli subnets leasetimes [subnetName] [active] [reserved] [flags] requires 3 arguments\n"
var subnetLeasetimesTooManyArgErrorString string = "Error: drpcli subnets leasetimes [subnetName] [active] [reserved] [flags] requires 3 arguments\n"
var subnetLeasetimesSuccessString string = `{
  "ActiveEnd": "192.168.100.200",
  "ActiveLeaseTime": 65,
  "ActiveStart": "192.168.100.10",
  "Available": true,
  "Enabled": false,
  "Errors": [],
  "Name": "john",
  "NextServer": "1.24.36.16",
  "OnlyReservations": false,
  "Options": [
    {
      "Code": 1,
      "Value": "255.192.0.0"
    },
    {
      "Code": 28,
      "Value": "192.191.255.255"
    }
  ],
  "Pickers": [
    "none",
    "nextFree",
    "mostExpired"
  ],
  "Proxy": false,
  "ReadOnly": false,
  "ReservedLeaseTime": 7300,
  "Strategy": "a3:b3:51:66:7e:11",
  "Subnet": "192.168.100.0/10",
  "Validated": true
}
`

var subnetSetNoArgErrorString string = "Error: drpcli subnets set [subnetName] option [number] to [value] [flags] requires 5 arguments\n"
var subnetSetTooManyArgErrorString string = "Error: drpcli subnets set [subnetName] option [number] to [value] [flags] requires 5 arguments\n"

var subnetSetTo66 string = `{
  "ActiveEnd": "192.168.100.200",
  "ActiveLeaseTime": 65,
  "ActiveStart": "192.168.100.10",
  "Available": true,
  "Enabled": false,
  "Errors": [],
  "Name": "john",
  "NextServer": "1.24.36.16",
  "OnlyReservations": false,
  "Options": [
    {
      "Code": 1,
      "Value": "255.192.0.0"
    },
    {
      "Code": 28,
      "Value": "192.191.255.255"
    },
    {
      "Code": 6,
      "Value": "66"
    }
  ],
  "Pickers": [
    "none",
    "nextFree",
    "mostExpired"
  ],
  "Proxy": false,
  "ReadOnly": false,
  "ReservedLeaseTime": 7300,
  "Strategy": "a3:b3:51:66:7e:11",
  "Subnet": "192.168.100.0/10",
  "Validated": true
}
`
var subnetSetToNull string = `{
  "ActiveEnd": "192.168.100.200",
  "ActiveLeaseTime": 65,
  "ActiveStart": "192.168.100.10",
  "Available": true,
  "Enabled": false,
  "Errors": [],
  "Name": "john",
  "NextServer": "1.24.36.16",
  "OnlyReservations": false,
  "Options": [
    {
      "Code": 1,
      "Value": "255.192.0.0"
    },
    {
      "Code": 28,
      "Value": "192.191.255.255"
    }
  ],
  "Pickers": [
    "none",
    "nextFree",
    "mostExpired"
  ],
  "Proxy": false,
  "ReadOnly": false,
  "ReservedLeaseTime": 7300,
  "Strategy": "a3:b3:51:66:7e:11",
  "Subnet": "192.168.100.0/10",
  "Validated": true
}
`
var subnetGetNoArgErrorString string = "Error: drpcli subnets get [subnetName] option [number] [flags] requires 3 arguments\n"
var subnetGetTooManyArgErrorString string = "Error: drpcli subnets get [subnetName] option [number] [flags] requires 3 arguments\n"
var subnetGetTo66 string = "Option 6: 66\n"

func TestSubnetCli(t *testing.T) {
	cliTest(true, false, "subnets").run(t)
	cliTest(false, false, "subnets", "list").run(t)
	cliTest(true, true, "subnets", "create").run(t)
	cliTest(true, true, "subnets", "create", "john", "john2").run(t)
	cliTest(false, true, "subnets", "create", subnetCreateBadJSONString).run(t)
	cliTest(false, false, "subnets", "create", subnetCreateInputString).run(t)
	cliTest(false, true, "subnets", "create", subnetCreateInputString).run(t)
	cliTest(false, false, "subnets", "list").run(t)
	cliTest(false, false, "subnets", "list", "Name=fred").run(t)
	cliTest(false, false, "subnets", "list", "Name=john").run(t)
	cliTest(false, false, "subnets", "list", "Strategy=MAC").run(t)
	cliTest(false, false, "subnets", "list", "Strategy=false").run(t)
	cliTest(false, false, "subnets", "list", "NextServer=3.3.3.3").run(t)
	cliTest(false, false, "subnets", "list", "NextServer=1.1.1.1").run(t)
	cliTest(false, true, "subnets", "list", "NextServer=fred").run(t)
	cliTest(false, false, "subnets", "list", "Enabled=false").run(t)
	cliTest(false, false, "subnets", "list", "Enabled=true").run(t)
	cliTest(false, true, "subnets", "list", "Enabled=george").run(t)
	cliTest(false, false, "subnets", "list", "Proxy=false").run(t)
	cliTest(false, false, "subnets", "list", "Proxy=true").run(t)
	cliTest(false, true, "subnets", "list", "Proxy=george").run(t)
	cliTest(false, false, "subnets", "list", "Subnet=192.168.103.0/24").run(t)
	cliTest(false, false, "subnets", "list", "Subnet=192.168.100.0/24").run(t)
	cliTest(false, true, "subnets", "list", "Subnet=false").run(t)
	cliTest(true, true, "subnets", "show").run(t)
	cliTest(true, true, "subnets", "show", "john", "john2").run(t)
	cliTest(false, true, "subnets", "show", "ignore").run(t)
	cliTest(false, false, "subnets", "show", "john").run(t)
	cliTest(true, true, "subnets", "exists").run(t)
	cliTest(true, true, "subnets", "exists", "john", "john2").run(t)
	cliTest(false, false, "subnets", "exists", "john").run(t)
	cliTest(false, true, "subnets", "exists", "ignore").run(t)
	cliTest(true, true, "subnets", "exists", "john", "john2").run(t)
	cliTest(true, true, "subnets", "update").run(t)
	cliTest(true, true, "subnets", "update", "john", "john2", "john3").run(t)
	cliTest(false, true, "subnets", "update", "john", subnetUpdateBadJSONString).run(t)
	cliTest(false, false, "subnets", "update", "john", subnetUpdateInputString).run(t)
	cliTest(false, true, "subnets", "update", "john2", subnetUpdateInputString).run(t)
	cliTest(false, false, "subnets", "show", "john").run(t)
	cliTest(false, false, "subnets", "show", "john").run(t)
	cliTest(true, true, "subnets", "destroy").run(t)
	cliTest(true, true, "subnets", "destroy", "john", "june").run(t)
	cliTest(false, false, "subnets", "destroy", "john").run(t)
	cliTest(false, true, "subnets", "destroy", "john").run(t)
	cliTest(false, false, "subnets", "list").run(t)
	cliTest(false, false, "subnets", "create", "-").Stdin(subnetCreateInputString + "\n").run(t)
	cliTest(false, false, "subnets", "list").run(t)
	cliTest(false, false, "subnets", "update", "john", "-").Stdin(subnetUpdateInputString + "\n").run(t)
	cliTest(false, false, "subnets", "show", "john").run(t)
	cliTest(true, true, "subnets", "range").run(t)
	cliTest(true, true, "subnets", "range", "john", "1.24.36.7", "1.24.36.16", "1.24.36.16").run(t)
	cliTest(false, true, "subnets", "range", "john", "192.168.100.10", "192.168.100.500").run(t)
	cliTest(false, true, "subnets", "range", "john", "cq.98.42.1234", "1.24.36.16").run(t)
	cliTest(false, false, "subnets", "range", "john", "192.168.100.10", "192.168.100.200").run(t)
	cliTest(true, true, "subnets", "subnet").run(t)
	cliTest(true, true, "subnets", "subnet", "john", "june", "1.24.36.16").run(t)
	cliTest(false, false, "subnets", "subnet", "john", "192.168.100.0/10").run(t)
	cliTest(false, true, "subnets", "subnet", "john", "1111.11.2223.544/66666").run(t)
	/* Save for when we have extra strategies other than MAC */
	/*
		cliTest(true, true, "subnets", "strategy").run(t)
		cliTest(true, true, "subnets", "strategy", "john", "june", "MAC").run(t)
		cliTest(false, false, "subnets", "strategy", "john", "MAC").run(t)
		cliTest(false, true, "subnets", "strategy", "john", "t5:44:llll:b").run(t)
	*/
	cliTest(true, true, "subnets", "pickers").run(t)
	cliTest(true, true, "subnets", "pickers", "john", "june", "test1,test2,test3").run(t)
	cliTest(false, false, "subnets", "pickers", "john", "none,nextFree,mostExpired").run(t)
	cliTest(true, true, "subnets", "nextserver").run(t)
	cliTest(true, true, "subnets", "nextserver", "john", "june", "1.24.36.16").run(t)
	cliTest(false, false, "subnets", "nextserver", "john", "1.24.36.16").run(t)
	cliTest(true, true, "subnets", "leasetimes").run(t)
	cliTest(true, true, "subnets", "leasetimes", "john", "june", "32", "55").run(t)
	cliTest(false, false, "subnets", "leasetimes", "john", "65", "7300").run(t)
	cliTest(false, true, "subnets", "leasetimes", "john", "4x5", "55").run(t)
	cliTest(true, true, "subnets", "set").run(t)
	cliTest(true, true, "subnets", "set", "john", "option", "45", "to", "34", "77").run(t)
	cliTest(true, true, "subnets", "get").run(t)
	cliTest(true, true, "subnets", "get", "john", "option", "45", "77").run(t)
	cliTest(false, true, "subnets", "set", "john", "option", "6tl", "to", "66").run(t)
	cliTest(false, false, "subnets", "set", "john", "option", "6", "to", "66").run(t)
	cliTest(false, false, "subnets", "get", "john", "option", "6").run(t)
	cliTest(false, false, "subnets", "set", "john", "option", "6", "to", "null").run(t)
	cliTest(false, true, "subnets", "get", "john", "option", "6").run(t)
	//End of Helpers
	cliTest(false, false, "subnets", "destroy", "john").run(t)
	cliTest(false, false, "subnets", "list").run(t)
}
