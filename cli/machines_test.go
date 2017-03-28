package cli

import (
	"testing"
)

var machineDefaultListString string = "[]\n"

var machineShowNoArgErrorString string = "Error: rscli machines show [id] requires 1 argument\n"
var machineShowTooManyArgErrorString string = "Error: rscli machines show [id] requires 1 argument\n"
var machineShowMissingArgErrorString string = "Error: machines GET: john: Not Found\n\n"
var machineShowMachineString string = `{
  "BootEnv": "ignore",
  "Errors": null,
  "Name": "john",
  "Uuid": "3e7031fe-3062-45f1-835c-92541bc9cbd3"
}
`

var machineExistsNoArgErrorString string = "Error: rscli machines exists [id] requires 1 argument"
var machineExistsTooManyArgErrorString string = "Error: rscli machines exists [id] requires 1 argument"
var machineExistsMachineString string = ""
var machineExistsMissingJohnString string = "Error: machines GET: john: Not Found\n\n"

var machineCreateNoArgErrorString string = "Error: rscli machines create [json] requires 1 argument\n"
var machineCreateTooManyArgErrorString string = "Error: rscli machines create [json] requires 1 argument\n"
var machineCreateBadJSONString = "asdgasdg"
var machineCreateBadJSONErrorString = "Error: Invalid machine object: error unmarshaling JSON: json: cannot unmarshal string into Go value of type models.Machine\n\n"
var machineCreateInputString string = `{
  "name": "john",
  "Uuid": "3e7031fe-3062-45f1-835c-92541bc9cbd3",
  "bootenv": "ignore"
}
`
var machineCreateJohnString string = `{
  "BootEnv": "ignore",
  "Errors": null,
  "Name": "john",
  "Uuid": "3e7031fe-3062-45f1-835c-92541bc9cbd3"
}
`
var machineCreateDuplicateErrorString = "Error: dataTracker create machines: 3e7031fe-3062-45f1-835c-92541bc9cbd3 already exists\n\n"

var machineListMachinesString = `[
  {
    "BootEnv": "ignore",
    "Errors": null,
    "Name": "john",
    "Uuid": "3e7031fe-3062-45f1-835c-92541bc9cbd3"
  }
]
`

var machineUpdateNoArgErrorString string = "Error: rscli machines update [id] [json] requires 2 arguments"
var machineUpdateTooManyArgErrorString string = "Error: rscli machines update [id] [json] requires 2 arguments"
var machineUpdateBadJSONString = "asdgasdg"
var machineUpdateBadJSONErrorString = "Error: Unable to merge objects: json: cannot unmarshal string into Go value of type map[string]interface {}\n\n\n"
var machineUpdateInputString string = `{
  "Description": "lpxelinux.0"
}
`
var machineUpdateJohnString string = `{
  "BootEnv": "ignore",
  "Description": "lpxelinux.0",
  "Errors": null,
  "Name": "john",
  "Uuid": "3e7031fe-3062-45f1-835c-92541bc9cbd3"
}
`
var machineUpdateJohnMissingErrorString string = "Error: machines GET: john2: Not Found\n\n"

var machinePatchNoArgErrorString string = "Error: rscli machines patch [objectJson] [changesJson] requires 2 arguments"
var machinePatchTooManyArgErrorString string = "Error: rscli machines patch [objectJson] [changesJson] requires 2 arguments"
var machinePatchBadPatchJSONString = "asdgasdg"
var machinePatchBadPatchJSONErrorString = "Error: Unable to parse rscli machines patch [objectJson] [changesJson] JSON asdgasdg\nError: error unmarshaling JSON: json: cannot unmarshal string into Go value of type models.Machine\n\n"
var machinePatchBadBaseJSONString = "asdgasdg"
var machinePatchBadBaseJSONErrorString = "Error: Unable to parse rscli machines patch [objectJson] [changesJson] JSON asdgasdg\nError: error unmarshaling JSON: json: cannot unmarshal string into Go value of type models.Machine\n\n"
var machinePatchBaseString string = `{
  "BootEnv": "ignore",
  "Description": "lpxelinux.0",
  "Errors": null,
  "Name": "john",
  "Uuid": "3e7031fe-3062-45f1-835c-92541bc9cbd3"
}
`
var machinePatchInputString string = `{
  "Description": "bootx64.efi"
}
`
var machinePatchJohnString string = `{
  "BootEnv": "ignore",
  "Description": "bootx64.efi",
  "Errors": null,
  "Name": "john",
  "Uuid": "3e7031fe-3062-45f1-835c-92541bc9cbd3"
}
`
var machinePatchMissingBaseString string = `{
  "BootEnv": "ignore",
  "Description": "bootx64.efi",
  "Errors": null,
  "Name": "john",
  "Uuid": "3e7031fe-5555-45f1-835c-92541bc9cbd3"
}
`
var machinePatchJohnMissingErrorString string = "Error: machines: PATCH 3e7031fe-5555-45f1-835c-92541bc9cbd3: Not Found\n\n"

var machineDestroyNoArgErrorString string = "Error: rscli machines destroy [id] requires 1 argument"
var machineDestroyTooManyArgErrorString string = "Error: rscli machines destroy [id] requires 1 argument"
var machineDestroyJohnString string = "Deleted machine 3e7031fe-3062-45f1-835c-92541bc9cbd3\n"
var machineDestroyMissingJohnString string = "Error: machines: DELETE 3e7031fe-3062-45f1-835c-92541bc9cbd3: Not Found\n\n"

func TestMachineCli(t *testing.T) {
	tests := []CliTest{
		CliTest{true, false, []string{"machines"}, noStdinString, "Access CLI commands relating to machines\n", ""},
		CliTest{false, false, []string{"machines", "list"}, noStdinString, machineDefaultListString, noErrorString},

		CliTest{true, true, []string{"machines", "create"}, noStdinString, noContentString, machineCreateNoArgErrorString},
		CliTest{true, true, []string{"machines", "create", "john", "john2"}, noStdinString, noContentString, machineCreateTooManyArgErrorString},
		CliTest{false, true, []string{"machines", "create", machineCreateBadJSONString}, noStdinString, noContentString, machineCreateBadJSONErrorString},
		CliTest{false, false, []string{"machines", "create", machineCreateInputString}, noStdinString, machineCreateJohnString, noErrorString},
		CliTest{false, true, []string{"machines", "create", machineCreateInputString}, noStdinString, noContentString, machineCreateDuplicateErrorString},
		CliTest{false, false, []string{"machines", "list"}, noStdinString, machineListMachinesString, noErrorString},

		CliTest{true, true, []string{"machines", "show"}, noStdinString, noContentString, machineShowNoArgErrorString},
		CliTest{true, true, []string{"machines", "show", "john", "john2"}, noStdinString, noContentString, machineShowTooManyArgErrorString},
		CliTest{false, true, []string{"machines", "show", "john"}, noStdinString, noContentString, machineShowMissingArgErrorString},
		CliTest{false, false, []string{"machines", "show", "3e7031fe-3062-45f1-835c-92541bc9cbd3"}, noStdinString, machineShowMachineString, noErrorString},

		CliTest{true, true, []string{"machines", "exists"}, noStdinString, noContentString, machineExistsNoArgErrorString},
		CliTest{true, true, []string{"machines", "exists", "john", "john2"}, noStdinString, noContentString, machineExistsTooManyArgErrorString},
		CliTest{false, false, []string{"machines", "exists", "3e7031fe-3062-45f1-835c-92541bc9cbd3"}, noStdinString, machineExistsMachineString, noErrorString},
		CliTest{false, true, []string{"machines", "exists", "john"}, noStdinString, noContentString, machineExistsMissingJohnString},
		CliTest{true, true, []string{"machines", "exists", "john", "john2"}, noStdinString, noContentString, machineExistsTooManyArgErrorString},

		CliTest{true, true, []string{"machines", "update"}, noStdinString, noContentString, machineUpdateNoArgErrorString},
		CliTest{true, true, []string{"machines", "update", "john", "john2", "john3"}, noStdinString, noContentString, machineUpdateTooManyArgErrorString},
		CliTest{false, true, []string{"machines", "update", "3e7031fe-3062-45f1-835c-92541bc9cbd3", machineUpdateBadJSONString}, noStdinString, noContentString, machineUpdateBadJSONErrorString},
		CliTest{false, false, []string{"machines", "update", "3e7031fe-3062-45f1-835c-92541bc9cbd3", machineUpdateInputString}, noStdinString, machineUpdateJohnString, noErrorString},
		CliTest{false, true, []string{"machines", "update", "john2", machineUpdateInputString}, noStdinString, noContentString, machineUpdateJohnMissingErrorString},
		CliTest{false, false, []string{"machines", "show", "3e7031fe-3062-45f1-835c-92541bc9cbd3"}, noStdinString, machineUpdateJohnString, noErrorString},

		CliTest{true, true, []string{"machines", "patch"}, noStdinString, noContentString, machinePatchNoArgErrorString},
		CliTest{true, true, []string{"machines", "patch", "john", "john2", "john3"}, noStdinString, noContentString, machinePatchTooManyArgErrorString},
		CliTest{false, true, []string{"machines", "patch", machinePatchBaseString, machinePatchBadPatchJSONString}, noStdinString, noContentString, machinePatchBadPatchJSONErrorString},
		CliTest{false, true, []string{"machines", "patch", machinePatchBadBaseJSONString, machinePatchInputString}, noStdinString, noContentString, machinePatchBadBaseJSONErrorString},
		CliTest{false, false, []string{"machines", "patch", machinePatchBaseString, machinePatchInputString}, noStdinString, machinePatchJohnString, noErrorString},
		CliTest{false, true, []string{"machines", "patch", machinePatchMissingBaseString, machinePatchInputString}, noStdinString, noContentString, machinePatchJohnMissingErrorString},
		CliTest{false, false, []string{"machines", "show", "3e7031fe-3062-45f1-835c-92541bc9cbd3"}, noStdinString, machinePatchJohnString, noErrorString},

		CliTest{true, true, []string{"machines", "destroy"}, noStdinString, noContentString, machineDestroyNoArgErrorString},
		CliTest{true, true, []string{"machines", "destroy", "john", "june"}, noStdinString, noContentString, machineDestroyTooManyArgErrorString},
		CliTest{false, false, []string{"machines", "destroy", "3e7031fe-3062-45f1-835c-92541bc9cbd3"}, noStdinString, machineDestroyJohnString, noErrorString},
		CliTest{false, true, []string{"machines", "destroy", "3e7031fe-3062-45f1-835c-92541bc9cbd3"}, noStdinString, noContentString, machineDestroyMissingJohnString},
		CliTest{false, false, []string{"machines", "list"}, noStdinString, machineDefaultListString, noErrorString},

		CliTest{false, false, []string{"machines", "create", "-"}, machineCreateInputString + "\n", machineCreateJohnString, noErrorString},
		CliTest{false, false, []string{"machines", "list"}, noStdinString, machineListMachinesString, noErrorString},
		CliTest{false, false, []string{"machines", "update", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "-"}, machineUpdateInputString + "\n", machineUpdateJohnString, noErrorString},
		CliTest{false, false, []string{"machines", "show", "3e7031fe-3062-45f1-835c-92541bc9cbd3"}, noStdinString, machineUpdateJohnString, noErrorString},
		CliTest{false, false, []string{"machines", "destroy", "3e7031fe-3062-45f1-835c-92541bc9cbd3"}, noStdinString, machineDestroyJohnString, noErrorString},
		CliTest{false, false, []string{"machines", "list"}, noStdinString, machineDefaultListString, noErrorString},
	}

	for _, test := range tests {
		testCli(t, test)
	}

}
