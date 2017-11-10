package cli

import (
	"io/ioutil"
	"net/http"
	"testing"
)

var machineEmptyListString = "[]\n"
var machineDefaultListString = "[]\n"

var machineActionMissingActionErrorString = "Error: GET: machines/3e7031fe-3062-45f1-835c-92541bc9cbd3: Action command: Not Found\n\n"
var machineActionMissingMachineErrorString = "Error: machines Action Get: john: Not Found\n\n"
var machineActionMissingErrorString = "Error: GET: machines/john: Action Get: 'command': Not Found\n\n"
var machineActionMissingParameterString = "Error: GET: machines/3e7031fe-3062-45f1-835c-92541bc9cbd3: Action reset_count Missing Parameter incrementer/touched\n\n"
var machineActionNoArgErrorString = "Error: drpcli machines action [id] [action] [flags] requires 2 argument"
var machineActionsMissingMachineErrorString = "Error: GET: machines/john: Not Found\n\n"
var machineActionsNoArgErrorString = "Error: drpcli machines actions [id] [flags] requires 1 argument"
var machineAddProfileJillJeanJillErrorString = "Error: ValidationError: machines/3e7031fe-3062-45f1-835c-92541bc9cbd3: Duplicate profile jill: at 0 and 2\n\n"
var machineAddProfileNoArgErrorString = "Error: drpcli machines addprofile [id] [profile] [flags] requires 2 arguments\n"
var machineAddrErrorString = "Error: GET: machines: Invalid address: fred\n\n"
var machineBadBoolString = "Error: GET: machines: Runnable must be true or false\n\n"
var machineBootEnvMissingMachineErrorString = "Error: GET: machines/john: Not Found\n\n"
var machineBootEnvNoArgErrorString = "Error: drpcli machines bootenv [id] [bootenv] [flags] requires 2 arguments"
var machineCreateBadJSON2ErrorString = "Error: Unable to create new machine: Invalid type passed to machine create\n\n"
var machineCreateBadJSONErrorString = "Error: Invalid machine object: error converting YAML to JSON: yaml: line 1: did not find expected ',' or '}' and error converting YAML to JSON: yaml: line 1: did not find expected ',' or '}'\n\n"
var machineCreateDuplicateErrorString = "Error: CREATE: machines/3e7031fe-3062-45f1-835c-92541bc9cbd3: already exists\n\n"
var machineCreateNoArgErrorString = "Error: drpcli machines create [json] [flags] requires 1 argument\n"
var machineCreateTooManyArgErrorString = "Error: drpcli machines create [json] [flags] requires 1 argument\n"
var machineDestroyMissingJohnString = "Error: DELETE: machines/3e7031fe-3062-45f1-835c-92541bc9cbd3: Not Found\n\n"
var machineDestroyNoArgErrorString = "Error: drpcli machines destroy [id] [flags] requires 1 argument"
var machineDestroyTooManyArgErrorString = "Error: drpcli machines destroy [id] [flags] requires 1 argument"
var machineExistsMachineString = ""
var machineExistsMissingJohnString = "Error: GET: machines/john: Not Found\n\n"
var machineExistsNoArgErrorString = "Error: drpcli machines exists [id] [flags] requires 1 argument"
var machineExistsTooManyArgErrorString = "Error: drpcli machines exists [id] [flags] requires 1 argument"
var machineExpireTimeErrorString = "Error: GET: machines: Invalid UUID: false\n\n"
var machineGetMissingMachineErrorString = "Error: GET: machines/john: Not Found\n\n"
var machineGetNoArgErrorString = "Error: drpcli machines get [id] param [key] [flags] requires 3 arguments"
var machineParamsMissingMachineErrorString = "Error: GET: machines/john2: Not Found\n\n"
var machineParamsNoArgErrorString = "Error: drpcli machines params [id] [json] [flags] requires 1 or 2 arguments\n"
var machinePatchBadBaseJSONErrorString = "Error: Unable to parse drpcli machines patch [objectJson] [changesJson] [flags] JSON asdgasdg\nError: error unmarshaling JSON: json: cannot unmarshal string into Go value of type genmodels.Machine\n\n"
var machinePatchBadPatchJSONErrorString = "Error: Unable to parse drpcli machines patch [objectJson] [changesJson] [flags] JSON asdgasdg\nError: error unmarshaling JSON: json: cannot unmarshal string into Go value of type genmodels.Machine\n\n"
var machinePatchJohnMissingErrorString = "Error: PATCH: machines/3e7031fe-5555-45f1-835c-92541bc9cbd3: Not Found\n\n"
var machinePatchNoArgErrorString = "Error: drpcli machines patch [objectJson] [changesJson] [flags] requires 2 arguments"
var machinePatchTooManyArgErrorString = "Error: drpcli machines patch [objectJson] [changesJson] [flags] requires 2 arguments"
var machineRemoveProfileNoArgErrorString = "Error: drpcli machines removeprofile [id] [profile] [flags] requires 2 arguments\n"
var machineRunActionBadCommandErrorString = "Error: INVOKE: machines/3e7031fe-3062-45f1-835c-92541bc9cbd3: Action command: Not Found\n\n"
var machineRunActionBadJSONThridArgErrorString = "Error: Invalid parameters: error unmarshaling JSON: json: cannot unmarshal string into Go value of type map[string]interface {}\n\n\n"
var machineRunActionBadStepErrorString = "Error: INVOKE: machines/3e7031fe-3062-45f1-835c-92541bc9cbd3: Action increment: Invalid Parameter: incrementer/step: : (root): Invalid type. Expected: integer, given: string\n\n"
var machineRunActionMissingCommandParametersErrorString = "Error: INVOKE: machines/3e7031fe-3062-45f1-835c-92541bc9cbd3: Action reset_count Missing Parameter incrementer/touched\n\n"
var machineRunActionMissingFredErrorString = "Error: INVOKE: machines/fred: Not Found\n\n"
var machineRunActionNoArgsErrorString = "Error: runaction either takes three arguments or a multiple of two, not 0"
var machineRunActionOneArgErrorString = "Error: runaction either takes three arguments or a multiple of two, not 1"
var machineSetMissingMachineErrorString = "Error: GET: machines/john: Not Found\n\n"
var machineSetNoArgErrorString = "Error: drpcli machines set [id] param [key] to [json blob] [flags] requires 5 arguments"
var machineShowMissingArgErrorString = "Error: GET: machines/john: Not Found\n\n"
var machineShowNoArgErrorString = "Error: drpcli machines show [id] [flags] requires 1 argument\n"
var machineShowTooManyArgErrorString = "Error: drpcli machines show [id] [flags] requires 1 argument\n"
var machineStageErrorStageString = "Error: ValidationError: machines/3e7031fe-3062-45f1-835c-92541bc9cbd3: Stage john2 does not exist\n\n"
var machineStageMissingMachineErrorString = "Error: GET: machines/john: Not Found\n\n"
var machineStageNoArgErrorString = "Error: drpcli machines stage [id] [stage] [flags] requires 2 arguments"
var machineUpdateBadJSONErrorString = "Error: Unable to merge objects: json: cannot unmarshal string into Go value of type map[string]interface {}\n\n\n"
var machineUpdateJohnMissingErrorString = "Error: GET: machines/john2: Not Found\n\n"
var machineUpdateNoArgErrorString = "Error: drpcli machines update [id] [json] [flags] requires 2 arguments"
var machineUpdateStagePendingErrorString = "Error: ValidationError: machines/3e7031fe-3062-45f1-835c-92541bc9cbd3: Can not change stages with pending tasks unless forced\n\n"
var machineUpdateTooManyArgErrorString = "Error: drpcli machines update [id] [json] [flags] requires 2 arguments"
var machineWaitBadBoolErrorString = "Error: strconv.ParseBool: parsing \"fred\": invalid syntax\n\n"
var machineWaitBadTimeoutErrorString = "Error: strconv.ParseInt: parsing \"jk\": invalid syntax\n\n"
var machineWaitMissingMachineErrorString = "Error: GET: machines/jk: Not Found\n\n"
var machineWaitNoArgErrorString = "Error: drpcli machines wait [id] [field] [value] [timeout] [flags] requires at least 3 arguments\n"
var machineWaitTooManyArgErrorString = "Error: drpcli machines wait [id] [field] [value] [timeout] [flags] requires at most 4 arguments\n"
var machinesParamsSetMissingMachineString = "Error: POST: machines/john2: Not Found\n\n"

var machineShowMachineString = `{
  "Address": "192.168.100.110",
  "Available": true,
  "BootEnv": "local",
  "CurrentTask": 0,
  "Errors": [],
  "Name": "john",
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

var machineCreateBadJSONString = "{asdgasdg"

var machineCreateBadJSON2String = "[asdgasdg]"

var machineCreateInputString = `{
  "Address": "192.168.100.110",
  "name": "john",
  "Secret": "secret1",
  "Uuid": "3e7031fe-3062-45f1-835c-92541bc9cbd3",
  "bootenv": "local"
}
`
var machineCreateJohnString = `{
  "Address": "192.168.100.110",
  "Available": true,
  "BootEnv": "local",
  "CurrentTask": 0,
  "Errors": [],
  "Name": "john",
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

var machineCreateJohnString2 = `{
  "Address": "192.168.100.110",
  "Available": true,
  "BootEnv": "local",
  "CurrentTask": -1,
  "Errors": [],
  "Name": "john",
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
  "Stage": "stage1",
  "Tasks": [
    "jamie",
    "justine"
  ],
  "Uuid": "3e7031fe-3062-45f1-835c-92541bc9cbd3",
  "Validated": true
}
`

var machineListMachinesString = `[
  {
    "Address": "192.168.100.110",
    "Available": true,
    "BootEnv": "local",
    "CurrentTask": 0,
    "Errors": [],
    "Name": "john",
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
]
`

var machineUpdateBadJSONString = "asdgasdg"

var machineUpdateInputString = `{
  "Description": "lpxelinux.0"
}
`
var machineUpdateJohnString = `{
  "Address": "192.168.100.110",
  "Available": true,
  "BootEnv": "local",
  "CurrentTask": 0,
  "Description": "lpxelinux.0",
  "Errors": [],
  "Name": "john",
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

var machinePatchBadPatchJSONString = "asdgasdg"

var machinePatchBadBaseJSONString = "asdgasdg"

var machinePatchBaseString = `{
  "Address": "192.168.100.110",
  "BootEnv": "local",
  "CurrentTask": 0,
  "Description": "lpxelinux.0",
  "Errors": [],
  "Name": "john",
  "Profile": {
    "Name": ""
  },
  "Profiles": [],
  "Runnable": true,
  "Secret": "secret1",
  "Tasks": [],
  "Uuid": "3e7031fe-3062-45f1-835c-92541bc9cbd3"
}
`
var machinePatchInputString = `{
  "Description": "bootx64.efi"
}
`
var machinePatchJohnString = `{
  "Address": "192.168.100.110",
  "Available": true,
  "BootEnv": "local",
  "CurrentTask": 0,
  "Description": "bootx64.efi",
  "Errors": [],
  "Name": "john",
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
var machinePatchMissingBaseString = `{
  "Address": "192.168.100.110",
  "BootEnv": "local",
  "Description": "lpxelinux.0",
  "Errors": [],
  "Name": "john",
  "Profile": {
    "Name": ""
  },
  "Profiles": [],
  "Runnable": true,
  "Tasks": [],
  "Uuid": "3e7031fe-5555-45f1-835c-92541bc9cbd3"
}
`

var machineAddProfileJill2String = `{
  "Address": "192.168.100.110",
  "Available": true,
  "BootEnv": "local2",
  "CurrentTask": 0,
  "Description": "lpxelinux.0",
  "Errors": [],
  "Name": "john",
  "Profile": {
    "Available": false,
    "Errors": [],
    "Name": "",
    "ReadOnly": false,
    "Validated": false
  },
  "Profiles": [
    "jill"
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
var machineAddProfileJillString = `{
  "Address": "192.168.100.110",
  "Available": true,
  "BootEnv": "local",
  "CurrentTask": 0,
  "Description": "lpxelinux.0",
  "Errors": [],
  "Name": "john",
  "Profile": {
    "Available": false,
    "Errors": [],
    "Name": "",
    "ReadOnly": false,
    "Validated": false
  },
  "Profiles": [
    "jill"
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
var machineAddProfileJillJeanString = `{
  "Address": "192.168.100.110",
  "Available": true,
  "BootEnv": "local",
  "CurrentTask": 0,
  "Description": "lpxelinux.0",
  "Errors": [],
  "Name": "john",
  "Profile": {
    "Available": false,
    "Errors": [],
    "Name": "",
    "ReadOnly": false,
    "Validated": false
  },
  "Profiles": [
    "jill",
    "jean"
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

var machineRemoveProfileJeanString = `{
  "Address": "192.168.100.110",
  "Available": true,
  "BootEnv": "local",
  "CurrentTask": 0,
  "Description": "lpxelinux.0",
  "Errors": [],
  "Name": "john",
  "Profile": {
    "Available": false,
    "Errors": [],
    "Name": "",
    "ReadOnly": false,
    "Validated": false
  },
  "Profiles": [
    "jean"
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
var machineRemoveProfileAllGoneString = `{
  "Address": "192.168.100.110",
  "Available": true,
  "BootEnv": "local",
  "CurrentTask": 0,
  "Description": "lpxelinux.0",
  "Errors": [],
  "Name": "john",
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

var machineRemoveProfileAllGone2String = `{
  "Address": "192.168.100.110",
  "Available": true,
  "BootEnv": "local2",
  "CurrentTask": -1,
  "Description": "lpxelinux.0",
  "Errors": [],
  "Name": "john",
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
  "Tasks": [
    "justine"
  ],
  "Uuid": "3e7031fe-3062-45f1-835c-92541bc9cbd3",
  "Validated": true
}
`

var machineDestroyJohnString = "Deleted machine 3e7031fe-3062-45f1-835c-92541bc9cbd3\n"

var machineBootEnvErrorBootEnvString = `{
  "Address": "192.168.100.110",
  "Available": false,
  "BootEnv": "john2",
  "CurrentTask": 0,
  "Description": "lpxelinux.0",
  "Errors": [
    "Bootenv john2 does not exist"
  ],
  "Name": "john",
  "Profile": {
    "Available": false,
    "Errors": [],
    "Name": "",
    "ReadOnly": false,
    "Validated": false
  },
  "Profiles": [],
  "ReadOnly": false,
  "Runnable": false,
  "Secret": "secret1",
  "Stage": "none",
  "Tasks": [],
  "Uuid": "3e7031fe-3062-45f1-835c-92541bc9cbd3",
  "Validated": true
}
`

var machineParamsStartingString = `{
  "asgdasdg": 1,
  "incrementer/default": 2,
  "incrementer/touched": 3,
  "john3": 4,
  "parm1": 1,
  "parm2": 10,
  "parm5": 20
}
`
var machinesParamsNextString = `{
  "jj": 3
}
`
var machineUpdateJohnWithParamsString = `{
  "Address": "192.168.100.110",
  "Available": true,
  "BootEnv": "local",
  "CurrentTask": 0,
  "Description": "lpxelinux.0",
  "Errors": [],
  "Name": "john",
  "Profile": {
    "Available": false,
    "Errors": [],
    "Name": "",
    "Params": {
      "jj": 3
    },
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

var machineStageProfCreate = `{
  "Available": true,
  "Errors": [],
  "Name": "stage-prof",
  "ReadOnly": false,
  "Validated": true
}
`

var machineJillCreate = `{
  "Available": true,
  "Errors": [],
  "Name": "jill",
  "ReadOnly": false,
  "Validated": true
}
`
var machineJeanCreate = `{
  "Available": true,
  "Errors": [],
  "Name": "jean",
  "ReadOnly": false,
  "Validated": true
}
`
var machineProfileJamieUpdate = `{
  "Available": true,
  "Errors": [],
  "Name": "jill",
  "ReadOnly": false,
  "Tasks": [
    "justine"
  ],
  "Validated": true
}
`

var machineActionsListString = `[
  {
    "Command": "increment",
    "OptionalParams": [
      "incrementer/step",
      "incrementer/parameter"
    ],
    "Provider": "incrementer",
    "RequiredParams": []
  }
]
`
var machineActionShowString = `{
  "Command": "increment",
  "OptionalParams": [
    "incrementer/step",
    "incrementer/parameter"
  ],
  "Provider": "incrementer",
  "RequiredParams": []
}
`

var machineActionsListWithResetString = `[
  {
    "Command": "increment",
    "OptionalParams": [
      "incrementer/step",
      "incrementer/parameter"
    ],
    "Provider": "incrementer",
    "RequiredParams": []
  },
  {
    "Command": "reset_count",
    "OptionalParams": [],
    "Provider": "incrementer",
    "RequiredParams": [
      "incrementer/touched"
    ]
  }
]
`
var machineActionShowResetString = `{
  "Command": "reset_count",
  "OptionalParams": [],
  "Provider": "incrementer",
  "RequiredParams": [
    "incrementer/touched"
  ]
}
`

var machinePluginCreateString = `{
  "Available": true,
  "Errors": [],
  "Name": "incr",
  "PluginErrors": [],
  "Provider": "incrementer",
  "ReadOnly": false,
  "Validated": true
}
`

var machineRunActionMissingParameterStdinString = "{}"
var machineRunActionGoodStdinString = `{
	"incrementer/parameter": "parm5",
	"incrementer/step": 10
}
`

var machineJamieCreate = `{
  "Available": true,
  "Errors": [],
  "Meta": {
    "feature-flags": "original-exit-codes"
  },
  "Name": "jamie",
  "OptionalParams": [],
  "ReadOnly": false,
  "RequiredParams": [],
  "Templates": [],
  "Validated": true
}
`
var machineJustineCreate = `{
  "Available": true,
  "Errors": [],
  "Meta": {
    "feature-flags": "original-exit-codes"
  },
  "Name": "justine",
  "OptionalParams": [],
  "ReadOnly": false,
  "RequiredParams": [],
  "Templates": [],
  "Validated": true
}
`

var machineUpdateLocalWithoutRunnableString = `{
  "Address": "192.168.100.110",
  "Available": true,
  "BootEnv": "local",
  "CurrentTask": 0,
  "Description": "lpxelinux.0",
  "Errors": [],
  "Name": "john",
  "Profile": {
    "Available": false,
    "Errors": [],
    "Name": "",
    "ReadOnly": false,
    "Validated": false
  },
  "Profiles": [],
  "ReadOnly": false,
  "Runnable": false,
  "Secret": "secret1",
  "Stage": "none",
  "Tasks": [],
  "Uuid": "3e7031fe-3062-45f1-835c-92541bc9cbd3",
  "Validated": true
}
`

var machineUpdateLocalString = `{
  "Address": "192.168.100.110",
  "Available": true,
  "BootEnv": "local",
  "CurrentTask": 0,
  "Description": "lpxelinux.0",
  "Errors": [],
  "Name": "john",
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
var machineUpdateLocal3String = `{
  "Address": "192.168.100.110",
  "Available": true,
  "BootEnv": "local",
  "CurrentTask": 0,
  "Description": "lpxelinux.0",
  "Errors": [],
  "Name": "john",
  "Profile": {
    "Available": false,
    "Errors": [],
    "Name": "",
    "ReadOnly": false,
    "Validated": false
  },
  "Profiles": [
    "jill",
    "jean"
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

var machineUpdateStage1WithoutRunnableString = `{
  "Address": "192.168.100.110",
  "Available": true,
  "BootEnv": "local",
  "CurrentTask": -1,
  "Description": "lpxelinux.0",
  "Errors": [],
  "Name": "john",
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
  "Stage": "stage1",
  "Tasks": [
    "jamie",
    "justine"
  ],
  "Uuid": "3e7031fe-3062-45f1-835c-92541bc9cbd3",
  "Validated": true
}
`

var machineUpdateStage1LocalString = `{
  "Address": "192.168.100.110",
  "Available": true,
  "BootEnv": "local",
  "CurrentTask": -1,
  "Description": "lpxelinux.0",
  "Errors": [],
  "Name": "john",
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
  "Stage": "stage1",
  "Tasks": [
    "jamie",
    "justine"
  ],
  "Uuid": "3e7031fe-3062-45f1-835c-92541bc9cbd3",
  "Validated": true
}
`

var machineUpdateLocal2String = `{
  "Address": "192.168.100.110",
  "Available": true,
  "BootEnv": "local",
  "CurrentTask": 0,
  "Description": "lpxelinux.0",
  "Errors": [],
  "Name": "john",
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
  "Stage": "stage2",
  "Tasks": [],
  "Uuid": "3e7031fe-3062-45f1-835c-92541bc9cbd3",
  "Validated": true
}
`

var machineStage1CreateString = `{
	"Name": "stage1",
	"BootEnv": "local",
	"Tasks": [ "jamie", "justine" ]
}
`
var machineStage1CreateSuccessString = `{
  "Available": true,
  "BootEnv": "local",
  "Errors": [],
  "Name": "stage1",
  "OptionalParams": [],
  "Profiles": [],
  "ReadOnly": false,
  "RequiredParams": [],
  "Tasks": [
    "jamie",
    "justine"
  ],
  "Templates": [],
  "Validated": true
}
`

var machineStage2CreateString = `{
  "Name": "stage2",
  "BootEnv": "local",
  "Templates": [
    {
      "Contents": "{{.Param \"sp-param\"}}",
      "Name": "test",
      "Path": "{{.Machine.Path}}/file"
    }
  ]
}
`
var machineStage2CreateSuccessString = `{
  "Available": true,
  "BootEnv": "local",
  "Errors": [],
  "Name": "stage2",
  "OptionalParams": [],
  "Profiles": [],
  "ReadOnly": false,
  "RequiredParams": [],
  "Tasks": [],
  "Templates": [
    {
      "Contents": "{{.Param \"sp-param\"}}",
      "Name": "test",
      "Path": "{{.Machine.Path}}/file"
    }
  ],
  "Validated": true
}
`

var machineStage2AgainSuccessString = `{
  "Address": "192.168.100.110",
  "Available": true,
  "BootEnv": "local",
  "CurrentTask": 0,
  "Description": "lpxelinux.0",
  "Errors": [],
  "Name": "john",
  "Profile": {
    "Available": false,
    "Errors": [],
    "Name": "",
    "ReadOnly": false,
    "Validated": false
  },
  "Profiles": [
    "jill",
    "jean"
  ],
  "ReadOnly": false,
  "Runnable": true,
  "Secret": "secret1",
  "Stage": "stage2",
  "Tasks": [],
  "Uuid": "3e7031fe-3062-45f1-835c-92541bc9cbd3",
  "Validated": true
}
`

var machinesUpdateStageSuccessString = `{
  "Available": true,
  "BootEnv": "local",
  "Errors": [],
  "Name": "stage2",
  "OptionalParams": [],
  "Profiles": [
    "stage-prof"
  ],
  "ReadOnly": false,
  "RequiredParams": [],
  "Tasks": [],
  "Templates": [
    {
      "Contents": "{{.Param \"sp-param\"}}",
      "Name": "test",
      "Path": "{{.Machine.Path}}/file"
    }
  ],
  "Validated": true
}
`

var machineAggregateParamString = `{
  "jill-param": "janga",
  "sp-param": "val"
}
`

var machinesSetDefaultStageString = `{
  "baseTokenSecret": "token-secret-token-secret-token1",
  "debugBootEnv": "0",
  "debugDhcp": "0",
  "debugFrontend": "0",
  "debugPlugins": "0",
  "debugRenderer": "0",
  "defaultBootEnv": "local",
  "defaultStage": "stage1",
  "knownTokenTimeout": "3600",
  "systemGrantorSecret": "system-grantor-secret",
  "unknownBootEnv": "ignore",
  "unknownTokenTimeout": "600"
}
`

var machinesSetDefaultStageBackString = `{
  "baseTokenSecret": "token-secret-token-secret-token1",
  "debugBootEnv": "0",
  "debugDhcp": "0",
  "debugFrontend": "0",
  "debugPlugins": "0",
  "debugRenderer": "0",
  "defaultBootEnv": "local",
  "defaultStage": "none",
  "knownTokenTimeout": "3600",
  "systemGrantorSecret": "system-grantor-secret",
  "unknownBootEnv": "ignore",
  "unknownTokenTimeout": "600"
}
`

func TestMachineCli(t *testing.T) {

	tests := []CliTest{
		CliTest{false, false, []string{"profiles", "create", "jill"}, noStdinString, machineJillCreate, noErrorString},
		CliTest{false, false, []string{"profiles", "create", "jean"}, noStdinString, machineJeanCreate, noErrorString},
		CliTest{false, false, []string{"profiles", "create", "stage-prof"}, noStdinString, machineStageProfCreate, noErrorString},
		CliTest{false, false, []string{"tasks", "create", "jamie"}, noStdinString, machineJamieCreate, noErrorString},
		CliTest{false, false, []string{"tasks", "create", "justine"}, noStdinString, machineJustineCreate, noErrorString},
		CliTest{false, false, []string{"stages", "create", machineStage1CreateString}, noStdinString, machineStage1CreateSuccessString, noErrorString},
		CliTest{false, false, []string{"stages", "create", machineStage2CreateString}, noStdinString, machineStage2CreateSuccessString, noErrorString},

		CliTest{false, false, []string{"plugins", "create", machinePluginCreateString}, noStdinString, machinePluginCreateString, noErrorString},

		CliTest{true, false, []string{"machines"}, noStdinString, "Access CLI commands relating to machines\n", ""},
		CliTest{false, false, []string{"machines", "list"}, noStdinString, machineDefaultListString, noErrorString},

		CliTest{true, true, []string{"machines", "create"}, noStdinString, noContentString, machineCreateNoArgErrorString},
		CliTest{true, true, []string{"machines", "create", "john", "john2"}, noStdinString, noContentString, machineCreateTooManyArgErrorString},
		CliTest{false, true, []string{"machines", "create", machineCreateBadJSONString}, noStdinString, noContentString, machineCreateBadJSONErrorString},
		CliTest{false, true, []string{"machines", "create", machineCreateBadJSON2String}, noStdinString, noContentString, machineCreateBadJSON2ErrorString},
		CliTest{false, false, []string{"machines", "create", machineCreateInputString}, noStdinString, machineCreateJohnString, noErrorString},
		CliTest{false, true, []string{"machines", "create", machineCreateInputString}, noStdinString, noContentString, machineCreateDuplicateErrorString},
		CliTest{false, false, []string{"machines", "list"}, noStdinString, machineListMachinesString, noErrorString},
		CliTest{false, false, []string{"machines", "list", "Name=fred"}, noStdinString, machineEmptyListString, noErrorString},
		CliTest{false, false, []string{"machines", "list", "Name=john"}, noStdinString, machineListMachinesString, noErrorString},
		CliTest{false, false, []string{"machines", "list", "BootEnv=local"}, noStdinString, machineListMachinesString, noErrorString},
		CliTest{false, false, []string{"machines", "list", "BootEnv=false"}, noStdinString, machineEmptyListString, noErrorString},
		CliTest{false, false, []string{"machines", "list", "Address=192.168.100.110"}, noStdinString, machineListMachinesString, noErrorString},
		CliTest{false, false, []string{"machines", "list", "Address=1.1.1.1"}, noStdinString, machineEmptyListString, noErrorString},
		CliTest{false, true, []string{"machines", "list", "Address=fred"}, noStdinString, noContentString, machineAddrErrorString},
		CliTest{false, false, []string{"machines", "list", "UUID=4e7031fe-3062-45f1-835c-92541bc9cbd3"}, noStdinString, machineEmptyListString, noErrorString},
		CliTest{false, false, []string{"machines", "list", "UUID=3e7031fe-3062-45f1-835c-92541bc9cbd3"}, noStdinString, machineListMachinesString, noErrorString},
		CliTest{false, true, []string{"machines", "list", "UUID=false"}, noStdinString, noContentString, machineExpireTimeErrorString},
		CliTest{false, false, []string{"machines", "list", "Runnable=true"}, noStdinString, machineListMachinesString, noErrorString},
		CliTest{false, false, []string{"machines", "list", "Runnable=false"}, noStdinString, machineEmptyListString, noErrorString},
		CliTest{false, true, []string{"machines", "list", "Runnable=fred"}, noStdinString, noContentString, machineBadBoolString},

		CliTest{true, true, []string{"machines", "show"}, noStdinString, noContentString, machineShowNoArgErrorString},
		CliTest{true, true, []string{"machines", "show", "john", "john2"}, noStdinString, noContentString, machineShowTooManyArgErrorString},
		CliTest{false, true, []string{"machines", "show", "john"}, noStdinString, noContentString, machineShowMissingArgErrorString},
		CliTest{false, false, []string{"machines", "show", "3e7031fe-3062-45f1-835c-92541bc9cbd3"}, noStdinString, machineShowMachineString, noErrorString},
		CliTest{false, false, []string{"machines", "show", "Key:3e7031fe-3062-45f1-835c-92541bc9cbd3"}, noStdinString, machineShowMachineString, noErrorString},
		CliTest{false, false, []string{"machines", "show", "Uuid:3e7031fe-3062-45f1-835c-92541bc9cbd3"}, noStdinString, machineShowMachineString, noErrorString},
		CliTest{false, false, []string{"machines", "show", "Name:john"}, noStdinString, machineShowMachineString, noErrorString},

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

		// bootenv tests
		CliTest{true, true, []string{"machines", "bootenv"}, noStdinString, noContentString, machineBootEnvNoArgErrorString},
		CliTest{false, true, []string{"machines", "bootenv", "john", "john2"}, noStdinString, noContentString, machineBootEnvMissingMachineErrorString},
		CliTest{false, false, []string{"machines", "bootenv", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "john2"}, noStdinString, machineBootEnvErrorBootEnvString, noErrorString},
		CliTest{false, false, []string{"machines", "bootenv", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "local"}, noStdinString, machineUpdateLocalWithoutRunnableString, noErrorString},
		CliTest{false, false, []string{"machines", "update", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "{ \"Runnable\": true }"}, noStdinString, machineUpdateLocalString, noErrorString},

		// stage tests
		CliTest{true, true, []string{"machines", "stage"}, noStdinString, noContentString, machineStageNoArgErrorString},
		CliTest{false, true, []string{"machines", "stage", "john", "john2"}, noStdinString, noContentString, machineStageMissingMachineErrorString},
		CliTest{false, true, []string{"machines", "stage", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "john2"}, noStdinString, noContentString, machineStageErrorStageString},
		CliTest{false, false, []string{"machines", "stage", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "stage1"}, noStdinString, machineUpdateStage1WithoutRunnableString, noErrorString},
		CliTest{false, false, []string{"machines", "update", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "{ \"Runnable\": true }"}, noStdinString, machineUpdateStage1LocalString, noErrorString},
		CliTest{false, true, []string{"machines", "stage", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "stage2"}, noStdinString, noContentString, machineUpdateStagePendingErrorString},
		CliTest{false, false, []string{"machines", "stage", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "stage2", "--force"}, noStdinString, machineUpdateLocal2String, noErrorString},
		CliTest{false, false, []string{"machines", "stage", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "", "--force"}, noStdinString, machineUpdateLocalString, noErrorString},

		// Add/Remove Profile tests
		CliTest{true, true, []string{"machines", "addprofile"}, noStdinString, noContentString, machineAddProfileNoArgErrorString},
		CliTest{false, false, []string{"machines", "addprofile", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "jill"}, noStdinString, machineAddProfileJillString, noErrorString},
		CliTest{false, false, []string{"machines", "addprofile", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "jean"}, noStdinString, machineAddProfileJillJeanString, noErrorString},
		CliTest{false, true, []string{"machines", "addprofile", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "jill"}, noStdinString, noContentString, machineAddProfileJillJeanJillErrorString},

		CliTest{false, false, []string{"profiles", "set", "jill", "param", "jill-param", "to", "janga"}, noStdinString, "\"janga\"\n", noErrorString},
		CliTest{false, false, []string{"profiles", "set", "stage-prof", "param", "sp-param", "to", "val"}, noStdinString, "\"val\"\n", noErrorString},
		CliTest{false, false, []string{"machines", "stage", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "stage2", "--force"}, noStdinString, machineStage2AgainSuccessString, noErrorString},
		CliTest{false, false, []string{"stages", "addprofile", "stage2", "stage-prof"}, noStdinString, machinesUpdateStageSuccessString, noErrorString},

		CliTest{false, false, []string{"machines", "params", "3e7031fe-3062-45f1-835c-92541bc9cbd3"}, noStdinString, "{}\n", noErrorString},
		CliTest{false, false, []string{"machines", "params", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "--aggregate"}, noStdinString, machineAggregateParamString, noErrorString},
	}
	for _, test := range tests {
		testCli(t, test)
	}

	tr := &http.Transport{}
	client := &http.Client{Transport: tr}
	req, _ := http.NewRequest("GET", "http://127.0.0.1:10002/machines/3e7031fe-3062-45f1-835c-92541bc9cbd3/file", nil)
	req.SetBasicAuth("rocketskates", "r0cketsk8ts")
	rsp, apierr := client.Do(req)
	if apierr != nil {
		t.Errorf("Failed to query machine file: %s", apierr)
	}
	defer rsp.Body.Close()
	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		t.Errorf("Failed to read all: %s", err)
	}
	if string(body) != "val" {
		t.Errorf("Body was: AA%sAA expected %s", string(body), "val")
	}

	tests2 := []CliTest{
		CliTest{false, false, []string{"machines", "stage", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "", "--force"}, noStdinString, machineUpdateLocal3String, noErrorString},
		CliTest{true, true, []string{"machines", "removeprofile"}, noStdinString, noContentString, machineRemoveProfileNoArgErrorString},
		CliTest{false, false, []string{"machines", "removeprofile", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "justine"}, noStdinString, machineAddProfileJillJeanString, noErrorString},
		CliTest{false, false, []string{"machines", "removeprofile", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "jill"}, noStdinString, machineRemoveProfileJeanString, noErrorString},
		CliTest{false, false, []string{"machines", "removeprofile", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "jean"}, noStdinString, machineRemoveProfileAllGoneString, noErrorString},

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

		CliTest{true, true, []string{"machines", "actions"}, noStdinString, noContentString, machineActionsNoArgErrorString},
		CliTest{false, true, []string{"machines", "actions", "john"}, noStdinString, noContentString, machineActionsMissingMachineErrorString},
		CliTest{false, false, []string{"machines", "actions", "3e7031fe-3062-45f1-835c-92541bc9cbd3"}, noStdinString, machineActionsListString, noErrorString},
		CliTest{true, true, []string{"machines", "action"}, noStdinString, noContentString, machineActionNoArgErrorString},
		CliTest{true, true, []string{"machines", "action", "john"}, noStdinString, noContentString, machineActionNoArgErrorString},
		CliTest{false, true, []string{"machines", "action", "john", "command"}, noStdinString, noContentString, machineActionMissingErrorString},
		CliTest{false, true, []string{"machines", "action", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "command"}, noStdinString, noContentString, machineActionMissingActionErrorString},
		CliTest{false, false, []string{"machines", "action", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "increment"}, noStdinString, machineActionShowString, noErrorString},

		CliTest{true, true, []string{"machines", "runaction"}, noStdinString, noContentString, machineRunActionNoArgsErrorString},
		CliTest{true, true, []string{"machines", "runaction", "fred"}, noStdinString, noContentString, machineRunActionOneArgErrorString},
		CliTest{false, true, []string{"machines", "runaction", "fred", "command"}, noStdinString, noContentString, machineRunActionMissingFredErrorString},
		CliTest{false, true, []string{"machines", "runaction", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "command"}, noStdinString, noContentString, machineRunActionBadCommandErrorString},
		CliTest{false, false, []string{"machines", "runaction", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "increment"}, noStdinString, "{}\n", noErrorString},
		CliTest{false, true, []string{"machines", "runaction", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "increment", "fred"}, noStdinString, noContentString, machineRunActionBadJSONThridArgErrorString},

		CliTest{false, false, []string{"machines", "actions", "3e7031fe-3062-45f1-835c-92541bc9cbd3"}, noStdinString, machineActionsListWithResetString, noErrorString},
		CliTest{false, false, []string{"machines", "action", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "reset_count"}, noStdinString, machineActionShowResetString, noErrorString},
		CliTest{false, false, []string{"machines", "runaction", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "reset_count"}, noStdinString, "{}\n", noErrorString},
		CliTest{false, false, []string{"machines", "actions", "3e7031fe-3062-45f1-835c-92541bc9cbd3"}, noStdinString, machineActionsListString, noErrorString},
		CliTest{false, true, []string{"machines", "action", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "reset_count"}, noStdinString, noContentString, machineActionMissingParameterString},
		CliTest{false, true, []string{"machines", "runaction", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "reset_count"}, noStdinString, noContentString, machineRunActionMissingCommandParametersErrorString},
		CliTest{false, false, []string{"machines", "runaction", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "increment", "incrementer/parameter", "asgdasdg"}, noStdinString, "{}\n", noErrorString},

		CliTest{false, false, []string{"machines", "runaction", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "increment", "incrementer/parameter", "parm1", "extra", "10"}, noStdinString, "{}\n", noErrorString},
		CliTest{false, false, []string{"machines", "get", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "param", "parm1"}, noStdinString, "1\n", noErrorString},
		CliTest{false, true, []string{"machines", "runaction", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "increment", "incrementer/parameter", "parm2", "incrementer/step", "asgdasdg"}, noStdinString, noContentString, machineRunActionBadStepErrorString},
		CliTest{false, false, []string{"machines", "get", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "param", "parm2"}, noStdinString, "null\n", noErrorString},
		CliTest{false, false, []string{"machines", "runaction", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "increment", "incrementer/parameter", "parm2", "incrementer/step", "10"}, noStdinString, "{}\n", noErrorString},
		CliTest{false, false, []string{"machines", "get", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "param", "parm2"}, noStdinString, "10\n", noErrorString},

		CliTest{false, true, []string{"machines", "runaction", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "increment", "-"}, "fred", noContentString, machineRunActionBadJSONThridArgErrorString},
		CliTest{false, false, []string{"machines", "runaction", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "reset_count", "-"}, machineRunActionMissingParameterStdinString, "{}\n", noErrorString},
		CliTest{false, true, []string{"machines", "runaction", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "reset_count", "-"}, machineRunActionMissingParameterStdinString, noContentString, machineRunActionMissingCommandParametersErrorString},
		CliTest{false, false, []string{"machines", "runaction", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "increment", "-"}, machineRunActionMissingParameterStdinString, "{}\n", noErrorString},
		CliTest{false, false, []string{"machines", "runaction", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "increment", "-"}, machineRunActionGoodStdinString, "{}\n", noErrorString},
		CliTest{false, false, []string{"machines", "runaction", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "increment", "-"}, machineRunActionGoodStdinString, "{}\n", noErrorString},
		CliTest{false, false, []string{"machines", "get", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "param", "parm5"}, noStdinString, "20\n", noErrorString},

		CliTest{true, true, []string{"machines", "wait"}, noStdinString, noContentString, machineWaitNoArgErrorString},
		CliTest{true, true, []string{"machines", "wait", "jk"}, noStdinString, noContentString, machineWaitNoArgErrorString},
		CliTest{true, true, []string{"machines", "wait", "jk", "jk"}, noStdinString, noContentString, machineWaitNoArgErrorString},
		CliTest{true, true, []string{"machines", "wait", "jk", "jk", "jk", "jk", "jk"}, noStdinString, noContentString, machineWaitTooManyArgErrorString},
		CliTest{false, true, []string{"machines", "wait", "jk", "jk", "jk", "jk"}, noStdinString, noContentString, machineWaitBadTimeoutErrorString},
		CliTest{false, true, []string{"machines", "wait", "jk", "jk", "jk"}, noStdinString, noContentString, machineWaitMissingMachineErrorString},
		CliTest{false, false, []string{"machines", "wait", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "jk", "jk", "1"}, noStdinString, "timeout\n", noErrorString},
		CliTest{false, false, []string{"machines", "wait", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "BootEnv", "local", "1"}, noStdinString, "complete\n", noErrorString},
		CliTest{false, true, []string{"machines", "wait", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "Runnable", "fred", "1"}, noStdinString, noContentString, machineWaitBadBoolErrorString},

		CliTest{true, true, []string{"machines", "params"}, noStdinString, noContentString, machineParamsNoArgErrorString},
		CliTest{false, true, []string{"machines", "params", "john2"}, noStdinString, noContentString, machineParamsMissingMachineErrorString},
		CliTest{false, false, []string{"machines", "params", "3e7031fe-3062-45f1-835c-92541bc9cbd3"}, noStdinString, machineParamsStartingString, noErrorString},
		CliTest{false, true, []string{"machines", "params", "john2", machinesParamsNextString}, noStdinString, noContentString, machinesParamsSetMissingMachineString},
		CliTest{false, false, []string{"machines", "params", "3e7031fe-3062-45f1-835c-92541bc9cbd3", machinesParamsNextString}, noStdinString, machinesParamsNextString, noErrorString},
		CliTest{false, false, []string{"machines", "params", "3e7031fe-3062-45f1-835c-92541bc9cbd3"}, noStdinString, machinesParamsNextString, noErrorString},

		CliTest{false, false, []string{"machines", "show", "3e7031fe-3062-45f1-835c-92541bc9cbd3"}, noStdinString, machineUpdateJohnWithParamsString, noErrorString},

		CliTest{false, false, []string{"machines", "destroy", "3e7031fe-3062-45f1-835c-92541bc9cbd3"}, noStdinString, machineDestroyJohnString, noErrorString},
		CliTest{false, false, []string{"machines", "list"}, noStdinString, machineDefaultListString, noErrorString},

		CliTest{false, false, []string{"prefs", "set", "defaultStage", "stage1"}, noStdinString, machinesSetDefaultStageString, noErrorString},
		CliTest{false, false, []string{"machines", "create", machineCreateInputString}, noStdinString, machineCreateJohnString2, noErrorString},
		CliTest{false, false, []string{"machines", "destroy", "3e7031fe-3062-45f1-835c-92541bc9cbd3"}, noStdinString, machineDestroyJohnString, noErrorString},
		CliTest{false, false, []string{"machines", "list"}, noStdinString, machineDefaultListString, noErrorString},
		CliTest{false, false, []string{"prefs", "set", "defaultStage", "none"}, noStdinString, machinesSetDefaultStageBackString, noErrorString},

		CliTest{false, false, []string{"plugins", "destroy", "incr"}, noStdinString, "Deleted plugin incr\n", noErrorString},
		CliTest{false, false, []string{"stages", "destroy", "stage1"}, noStdinString, "Deleted stage stage1\n", noErrorString},
		CliTest{false, false, []string{"stages", "destroy", "stage2"}, noStdinString, "Deleted stage stage2\n", noErrorString},
		CliTest{false, false, []string{"profiles", "destroy", "jill"}, noStdinString, "Deleted profile jill\n", noErrorString},
		CliTest{false, false, []string{"profiles", "destroy", "jean"}, noStdinString, "Deleted profile jean\n", noErrorString},
		CliTest{false, false, []string{"profiles", "destroy", "stage-prof"}, noStdinString, "Deleted profile stage-prof\n", noErrorString},
		CliTest{false, false, []string{"tasks", "destroy", "jamie"}, noStdinString, "Deleted task jamie\n", noErrorString},
		CliTest{false, false, []string{"tasks", "destroy", "justine"}, noStdinString, "Deleted task justine\n", noErrorString},
	}
	for _, test := range tests2 {
		testCli(t, test)
	}
}
