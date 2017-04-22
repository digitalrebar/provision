package cli

import (
	"os"
	"testing"
)

var machineDefaultListString string = "[]\n"

var machineShowNoArgErrorString string = "Error: drpcli machines show [id] requires 1 argument\n"
var machineShowTooManyArgErrorString string = "Error: drpcli machines show [id] requires 1 argument\n"
var machineShowMissingArgErrorString string = "Error: machines GET: john: Not Found\n\n"
var machineShowMachineString string = `{
  "BootEnv": "local",
  "Errors": null,
  "Name": "john",
  "Profiles": null,
  "Uuid": "3e7031fe-3062-45f1-835c-92541bc9cbd3"
}
`

var machineExistsNoArgErrorString string = "Error: drpcli machines exists [id] requires 1 argument"
var machineExistsTooManyArgErrorString string = "Error: drpcli machines exists [id] requires 1 argument"
var machineExistsMachineString string = ""
var machineExistsMissingJohnString string = "Error: machines GET: john: Not Found\n\n"

var machineCreateNoArgErrorString string = "Error: drpcli machines create [json] requires 1 argument\n"
var machineCreateTooManyArgErrorString string = "Error: drpcli machines create [json] requires 1 argument\n"
var machineCreateBadJSONString = "asdgasdg"
var machineCreateBadJSONErrorString = "Error: Invalid machine object: error unmarshaling JSON: json: cannot unmarshal string into Go value of type models.Machine\n\n"
var machineCreateInputString string = `{
  "name": "john",
  "Uuid": "3e7031fe-3062-45f1-835c-92541bc9cbd3",
  "bootenv": "local"
}
`
var machineCreateJohnString string = `{
  "BootEnv": "local",
  "Errors": null,
  "Name": "john",
  "Profiles": null,
  "Uuid": "3e7031fe-3062-45f1-835c-92541bc9cbd3"
}
`
var machineCreateDuplicateErrorString = "Error: dataTracker create machines: 3e7031fe-3062-45f1-835c-92541bc9cbd3 already exists\n\n"

var machineListMachinesString = `[
  {
    "BootEnv": "local",
    "Errors": null,
    "Name": "john",
    "Profiles": null,
    "Uuid": "3e7031fe-3062-45f1-835c-92541bc9cbd3"
  }
]
`

var machineUpdateNoArgErrorString string = "Error: drpcli machines update [id] [json] requires 2 arguments"
var machineUpdateTooManyArgErrorString string = "Error: drpcli machines update [id] [json] requires 2 arguments"
var machineUpdateBadJSONString = "asdgasdg"
var machineUpdateBadJSONErrorString = "Error: Unable to merge objects: json: cannot unmarshal string into Go value of type map[string]interface {}\n\n\n"
var machineUpdateInputString string = `{
  "Description": "lpxelinux.0"
}
`
var machineUpdateJohnString string = `{
  "BootEnv": "local",
  "Description": "lpxelinux.0",
  "Errors": null,
  "Name": "john",
  "Profiles": null,
  "Uuid": "3e7031fe-3062-45f1-835c-92541bc9cbd3"
}
`
var machineUpdateJohnMissingErrorString string = "Error: machines GET: john2: Not Found\n\n"

var machinePatchNoArgErrorString string = "Error: drpcli machines patch [objectJson] [changesJson] requires 2 arguments"
var machinePatchTooManyArgErrorString string = "Error: drpcli machines patch [objectJson] [changesJson] requires 2 arguments"
var machinePatchBadPatchJSONString = "asdgasdg"
var machinePatchBadPatchJSONErrorString = "Error: Unable to parse drpcli machines patch [objectJson] [changesJson] JSON asdgasdg\nError: error unmarshaling JSON: json: cannot unmarshal string into Go value of type models.Machine\n\n"
var machinePatchBadBaseJSONString = "asdgasdg"
var machinePatchBadBaseJSONErrorString = "Error: Unable to parse drpcli machines patch [objectJson] [changesJson] JSON asdgasdg\nError: error unmarshaling JSON: json: cannot unmarshal string into Go value of type models.Machine\n\n"
var machinePatchBaseString string = `{
  "BootEnv": "local",
  "Description": "lpxelinux.0",
  "Errors": null,
  "Name": "john",
  "Profiles": null,
  "Uuid": "3e7031fe-3062-45f1-835c-92541bc9cbd3"
}
`
var machinePatchInputString string = `{
  "Description": "bootx64.efi"
}
`
var machinePatchJohnString string = `{
  "BootEnv": "local",
  "Description": "bootx64.efi",
  "Errors": null,
  "Name": "john",
  "Profiles": null,
  "Uuid": "3e7031fe-3062-45f1-835c-92541bc9cbd3"
}
`
var machinePatchMissingBaseString string = `{
  "BootEnv": "local",
  "Description": "bootx64.efi",
  "Errors": null,
  "Name": "john",
  "Profiles": null,
  "Uuid": "3e7031fe-5555-45f1-835c-92541bc9cbd3"
}
`
var machinePatchJohnMissingErrorString string = "Error: machines: PATCH 3e7031fe-5555-45f1-835c-92541bc9cbd3: Not Found\n\n"

var machineDestroyNoArgErrorString string = "Error: drpcli machines destroy [id] requires 1 argument"
var machineDestroyTooManyArgErrorString string = "Error: drpcli machines destroy [id] requires 1 argument"
var machineDestroyJohnString string = "Deleted machine 3e7031fe-3062-45f1-835c-92541bc9cbd3\n"
var machineDestroyMissingJohnString string = "Error: machines: DELETE 3e7031fe-3062-45f1-835c-92541bc9cbd3: Not Found\n\n"

var machineBootEnvNoArgErrorString string = "Error: drpcli machines bootenv [id] [bootenv] requires 2 arguments"
var machineBootEnvMissingMachineErrorString string = "Error: machines GET: john: Not Found\n\n"
var machineBootEnvBadBootEnvErrorString string = "Error: Machine 3e7031fe-3062-45f1-835c-92541bc9cbd3 has BootEnv john2, which is not present in the DataTracker\n\n"

var machineGetNoArgErrorString string = "Error: drpcli machines get [id] param [key] requires 3 arguments"
var machineGetMissingMachineErrorString string = "Error: machines GET Params: john: Not Found\n\n"

var machineSetNoArgErrorString string = "Error: drpcli machines set [id] param [key] to [json blob] requires 5 arguments"
var machineSetMissingMachineErrorString string = "Error: machines GET Params: john: Not Found\n\n"

var machineParamsNoArgErrorString string = "Error: drpcli machines params [id] [json] requires 1 or 2 arguments\n"
var machineParamsMissingMachineErrorString string = "Error: machines GET Params: john2: Not Found\n\n"
var machinesParamsSetMissingMachineString string = "Error: machines SET Params: john2: Not Found\n\n"

var machineParamsStartingString string = `{
  "john3": 4
}
`
var machinesParamsNextString string = `{
  "jj": 3
}
`

func TestMachineCli(t *testing.T) {
	if err := os.MkdirAll("bootenvs", 0755); err != nil {
		t.Errorf("Failed to create bootenvs dir: %v\n", err)
	}
	if err := os.Symlink("../../assets/bootenvs/local.yml", "bootenvs/local.yml"); err != nil {
		t.Errorf("Failed to create link to local.yml: %v\n", err)
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

	tests := []CliTest{
		CliTest{false, false, []string{"bootenvs", "install", "bootenvs/local.yml"}, noStdinString, bootEnvInstallLocalSuccessString, noErrorString},

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

		CliTest{true, true, []string{"machines", "bootenv"}, noStdinString, noContentString, machineBootEnvNoArgErrorString},
		CliTest{false, true, []string{"machines", "bootenv", "john", "john2"}, noStdinString, noContentString, machineBootEnvMissingMachineErrorString},
		CliTest{false, true, []string{"machines", "bootenv", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "john2"}, noStdinString, noContentString, machineBootEnvBadBootEnvErrorString},
		CliTest{false, false, []string{"machines", "bootenv", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "local"}, noStdinString, machineUpdateJohnString, noErrorString},

		CliTest{true, true, []string{"machines", "get"}, noStdinString, noContentString, machineGetNoArgErrorString},
		CliTest{false, true, []string{"machines", "get", "john", "param", "john2"}, noStdinString, noContentString, machineGetMissingMachineErrorString},
		CliTest{false, false, []string{"machines", "get", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "param", "john2"}, noStdinString, "null\n", noErrorString},

		CliTest{true, true, []string{"machines", "set"}, noStdinString, noContentString, machineSetNoArgErrorString},
		CliTest{false, true, []string{"machines", "set", "john", "param", "john2", "to", "cow"}, noStdinString, noContentString, machineSetMissingMachineErrorString},
		CliTest{false, false, []string{"machines", "set", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "param", "john2", "to", "cow"}, noStdinString, "\"cow\"\n", noErrorString},
		CliTest{false, false, []string{"machines", "get", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "param", "john2"}, noStdinString, "\"cow\"\n", noErrorString},
		CliTest{false, false, []string{"machines", "set", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "param", "john2", "to", "3"}, noStdinString, "3\n", noErrorString},
		CliTest{false, false, []string{"machines", "set", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "param", "john3", "to", "4"}, noStdinString, "4\n", noErrorString},
		CliTest{false, false, []string{"machines", "get", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "param", "john2"}, noStdinString, "3\n", noErrorString},
		CliTest{false, false, []string{"machines", "get", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "param", "john3"}, noStdinString, "4\n", noErrorString},
		CliTest{false, false, []string{"machines", "set", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "param", "john2", "to", "null"}, noStdinString, "null\n", noErrorString},
		CliTest{false, false, []string{"machines", "get", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "param", "john2"}, noStdinString, "null\n", noErrorString},
		CliTest{false, false, []string{"machines", "get", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "param", "john3"}, noStdinString, "4\n", noErrorString},

		CliTest{true, true, []string{"machines", "params"}, noStdinString, noContentString, machineParamsNoArgErrorString},
		CliTest{false, true, []string{"machines", "params", "john2"}, noStdinString, noContentString, machineParamsMissingMachineErrorString},
		CliTest{false, false, []string{"machines", "params", "3e7031fe-3062-45f1-835c-92541bc9cbd3"}, noStdinString, machineParamsStartingString, noErrorString},
		CliTest{false, true, []string{"machines", "params", "john2", machinesParamsNextString}, noStdinString, noContentString, machinesParamsSetMissingMachineString},
		CliTest{false, false, []string{"machines", "params", "3e7031fe-3062-45f1-835c-92541bc9cbd3", machinesParamsNextString}, noStdinString, machinesParamsNextString, noErrorString},
		CliTest{false, false, []string{"machines", "params", "3e7031fe-3062-45f1-835c-92541bc9cbd3"}, noStdinString, machinesParamsNextString, noErrorString},

		CliTest{false, false, []string{"machines", "destroy", "3e7031fe-3062-45f1-835c-92541bc9cbd3"}, noStdinString, machineDestroyJohnString, noErrorString},
		CliTest{false, false, []string{"machines", "list"}, noStdinString, machineDefaultListString, noErrorString},

		CliTest{false, false, []string{"bootenvs", "destroy", "local"}, noStdinString, "Deleted bootenv local\n", noErrorString},
		CliTest{false, false, []string{"templates", "destroy", "local-pxelinux.tmpl"}, noStdinString, "Deleted template local-pxelinux.tmpl\n", noErrorString},
		CliTest{false, false, []string{"templates", "destroy", "local-elilo.tmpl"}, noStdinString, "Deleted template local-elilo.tmpl\n", noErrorString},
		CliTest{false, false, []string{"templates", "destroy", "local-ipxe.tmpl"}, noStdinString, "Deleted template local-ipxe.tmpl\n", noErrorString},
	}

	for _, test := range tests {
		testCli(t, test)
	}

	os.RemoveAll("bootenvs")
	os.RemoveAll("templates")
	os.RemoveAll("isos")
	os.RemoveAll("ic")
}
