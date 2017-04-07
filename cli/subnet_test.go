package cli

import (
	"testing"
)

var subnetDefaultListString string = "[]\n"

var subnetShowNoArgErrorString string = "Error: drpcli subnets show [id] requires 1 argument\n"
var subnetShowTooManyArgErrorString string = "Error: drpcli subnets show [id] requires 1 argument\n"
var subnetShowMissingArgErrorString string = "Error: subnets GET: ignore: Not Found\n\n"
var subnetShowJohnString string = `{
  "ActiveEnd": "192.168.100.100",
  "ActiveLeaseTime": 60,
  "ActiveStart": "192.168.100.20",
  "Name": "john",
  "NextServer": "",
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
  "Strategy": "MAC",
  "Subnet": "192.168.100.0/24"
}
`

var subnetExistsNoArgErrorString string = "Error: drpcli subnets exists [id] requires 1 argument"
var subnetExistsTooManyArgErrorString string = "Error: drpcli subnets exists [id] requires 1 argument"
var subnetExistsIgnoreString string = ""
var subnetExistsMissingIgnoreString string = "Error: subnets GET: ignore: Not Found\n\n"

var subnetCreateNoArgErrorString string = "Error: drpcli subnets create [json] requires 1 argument\n"
var subnetCreateTooManyArgErrorString string = "Error: drpcli subnets create [json] requires 1 argument\n"
var subnetCreateBadJSONString = "asdgasdg"
var subnetCreateBadJSONErrorString = "Error: Invalid subnet object: error unmarshaling JSON: json: cannot unmarshal string into Go value of type models.Subnet\n\n"
var subnetCreateInputString string = `{
  "Name": "john",
  "ActiveEnd": "192.168.100.100",
  "ActiveStart": "192.168.100.20",
  "ActiveLeaseTime": 60,
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
  "Name": "john",
  "NextServer": "",
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
  "Strategy": "MAC",
  "Subnet": "192.168.100.0/24"
}
`
var subnetCreateDuplicateErrorString = "Error: dataTracker create subnets: john already exists\n\n"

var subnetListBothEnvsString = `[
  {
    "ActiveEnd": "192.168.100.100",
    "ActiveLeaseTime": 60,
    "ActiveStart": "192.168.100.20",
    "Name": "john",
    "NextServer": "",
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
    "Strategy": "MAC",
    "Subnet": "192.168.100.0/24"
  }
]
`

var subnetUpdateNoArgErrorString string = "Error: drpcli subnets update [id] [json] requires 2 arguments"
var subnetUpdateTooManyArgErrorString string = "Error: drpcli subnets update [id] [json] requires 2 arguments"
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
  "Name": "john",
  "NextServer": "",
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
  "Strategy": "NewStrat",
  "Subnet": "192.168.100.0/24"
}
`
var subnetUpdateJohnMissingErrorString string = "Error: subnets GET: john2: Not Found\n\n"

var subnetPatchNoArgErrorString string = "Error: drpcli subnets patch [objectJson] [changesJson] requires 2 arguments"
var subnetPatchTooManyArgErrorString string = "Error: drpcli subnets patch [objectJson] [changesJson] requires 2 arguments"
var subnetPatchBadPatchJSONString = "asdgasdg"
var subnetPatchBadPatchJSONErrorString = "Error: Unable to parse drpcli subnets patch [objectJson] [changesJson] JSON asdgasdg\nError: error unmarshaling JSON: json: cannot unmarshal string into Go value of type models.Subnet\n\n"
var subnetPatchBadBaseJSONString = "asdgasdg"
var subnetPatchBadBaseJSONErrorString = "Error: Unable to parse drpcli subnets patch [objectJson] [changesJson] JSON asdgasdg\nError: error unmarshaling JSON: json: cannot unmarshal string into Go value of type models.Subnet\n\n"
var subnetPatchBaseString string = `{
  "ActiveEnd": "192.168.100.100",
  "ActiveLeaseTime": 60,
  "ActiveStart": "192.168.100.20",
  "Name": "john",
  "NextServer": "",
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
  "Strategy": "NewStrat",
  "Subnet": "192.168.100.0/24"
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
  "Name": "john",
  "NextServer": "",
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
var subnetPatchMissingBaseString string = `{
  "ActiveEnd": "192.168.100.100",
  "ActiveLeaseTime": 60,
  "ActiveStart": "192.168.100.20",
  "Name": "john2",
  "NextServer": "",
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
var subnetPatchJohnMissingErrorString string = "Error: subnets: PATCH john2: Not Found\n\n"

var subnetDestroyNoArgErrorString string = "Error: drpcli subnets destroy [id] requires 1 argument"
var subnetDestroyTooManyArgErrorString string = "Error: drpcli subnets destroy [id] requires 1 argument"
var subnetDestroyJohnString string = "Deleted subnet john\n"
var subnetDestroyMissingJohnString string = "Error: subnets: DELETE john: Not Found\n\n"

func TestSubnetCli(t *testing.T) {
	tests := []CliTest{
		CliTest{true, false, []string{"subnets"}, noStdinString, "Access CLI commands relating to subnets\n", ""},
		CliTest{false, false, []string{"subnets", "list"}, noStdinString, subnetDefaultListString, noErrorString},

		CliTest{true, true, []string{"subnets", "create"}, noStdinString, noContentString, subnetCreateNoArgErrorString},
		CliTest{true, true, []string{"subnets", "create", "john", "john2"}, noStdinString, noContentString, subnetCreateTooManyArgErrorString},
		CliTest{false, true, []string{"subnets", "create", subnetCreateBadJSONString}, noStdinString, noContentString, subnetCreateBadJSONErrorString},
		CliTest{false, false, []string{"subnets", "create", subnetCreateInputString}, noStdinString, subnetCreateJohnString, noErrorString},
		CliTest{false, true, []string{"subnets", "create", subnetCreateInputString}, noStdinString, noContentString, subnetCreateDuplicateErrorString},
		CliTest{false, false, []string{"subnets", "list"}, noStdinString, subnetListBothEnvsString, noErrorString},

		CliTest{true, true, []string{"subnets", "show"}, noStdinString, noContentString, subnetShowNoArgErrorString},
		CliTest{true, true, []string{"subnets", "show", "john", "john2"}, noStdinString, noContentString, subnetShowTooManyArgErrorString},
		CliTest{false, true, []string{"subnets", "show", "ignore"}, noStdinString, noContentString, subnetShowMissingArgErrorString},
		CliTest{false, false, []string{"subnets", "show", "john"}, noStdinString, subnetShowJohnString, noErrorString},

		CliTest{true, true, []string{"subnets", "exists"}, noStdinString, noContentString, subnetExistsNoArgErrorString},
		CliTest{true, true, []string{"subnets", "exists", "john", "john2"}, noStdinString, noContentString, subnetExistsTooManyArgErrorString},
		CliTest{false, false, []string{"subnets", "exists", "john"}, noStdinString, subnetExistsIgnoreString, noErrorString},
		CliTest{false, true, []string{"subnets", "exists", "ignore"}, noStdinString, noContentString, subnetExistsMissingIgnoreString},
		CliTest{true, true, []string{"subnets", "exists", "john", "john2"}, noStdinString, noContentString, subnetExistsTooManyArgErrorString},

		CliTest{true, true, []string{"subnets", "update"}, noStdinString, noContentString, subnetUpdateNoArgErrorString},
		CliTest{true, true, []string{"subnets", "update", "john", "john2", "john3"}, noStdinString, noContentString, subnetUpdateTooManyArgErrorString},
		CliTest{false, true, []string{"subnets", "update", "john", subnetUpdateBadJSONString}, noStdinString, noContentString, subnetUpdateBadJSONErrorString},
		CliTest{false, false, []string{"subnets", "update", "john", subnetUpdateInputString}, noStdinString, subnetUpdateJohnString, noErrorString},
		CliTest{false, true, []string{"subnets", "update", "john2", subnetUpdateInputString}, noStdinString, noContentString, subnetUpdateJohnMissingErrorString},
		CliTest{false, false, []string{"subnets", "show", "john"}, noStdinString, subnetUpdateJohnString, noErrorString},

		CliTest{true, true, []string{"subnets", "patch"}, noStdinString, noContentString, subnetPatchNoArgErrorString},
		CliTest{true, true, []string{"subnets", "patch", "john", "john2", "john3"}, noStdinString, noContentString, subnetPatchTooManyArgErrorString},
		CliTest{false, true, []string{"subnets", "patch", subnetPatchBaseString, subnetPatchBadPatchJSONString}, noStdinString, noContentString, subnetPatchBadPatchJSONErrorString},
		CliTest{false, true, []string{"subnets", "patch", subnetPatchBadBaseJSONString, subnetPatchInputString}, noStdinString, noContentString, subnetPatchBadBaseJSONErrorString},
		CliTest{false, false, []string{"subnets", "patch", subnetPatchBaseString, subnetPatchInputString}, noStdinString, subnetPatchJohnString, noErrorString},
		CliTest{false, true, []string{"subnets", "patch", subnetPatchMissingBaseString, subnetPatchInputString}, noStdinString, noContentString, subnetPatchJohnMissingErrorString},
		CliTest{false, false, []string{"subnets", "show", "john"}, noStdinString, subnetPatchJohnString, noErrorString},

		CliTest{true, true, []string{"subnets", "destroy"}, noStdinString, noContentString, subnetDestroyNoArgErrorString},
		CliTest{true, true, []string{"subnets", "destroy", "john", "june"}, noStdinString, noContentString, subnetDestroyTooManyArgErrorString},
		CliTest{false, false, []string{"subnets", "destroy", "john"}, noStdinString, subnetDestroyJohnString, noErrorString},
		CliTest{false, true, []string{"subnets", "destroy", "john"}, noStdinString, noContentString, subnetDestroyMissingJohnString},
		CliTest{false, false, []string{"subnets", "list"}, noStdinString, subnetDefaultListString, noErrorString},

		CliTest{false, false, []string{"subnets", "create", "-"}, subnetCreateInputString + "\n", subnetCreateJohnString, noErrorString},
		CliTest{false, false, []string{"subnets", "list"}, noStdinString, subnetListBothEnvsString, noErrorString},
		CliTest{false, false, []string{"subnets", "update", "john", "-"}, subnetUpdateInputString + "\n", subnetUpdateJohnString, noErrorString},
		CliTest{false, false, []string{"subnets", "show", "john"}, noStdinString, subnetUpdateJohnString, noErrorString},

		CliTest{false, false, []string{"subnets", "destroy", "john"}, noStdinString, subnetDestroyJohnString, noErrorString},
		CliTest{false, false, []string{"subnets", "list"}, noStdinString, subnetDefaultListString, noErrorString},
	}

	for _, test := range tests {
		testCli(t, test)
	}

}
