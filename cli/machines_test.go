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
	cliTest(false, false, "profiles", "create", "jill").run(t)
	cliTest(false, false, "profiles", "create", "jean").run(t)
	cliTest(false, false, "profiles", "create", "stage-prof").run(t)
	cliTest(false, false, "tasks", "create", "jamie").run(t)
	cliTest(false, false, "tasks", "create", "justine").run(t)
	cliTest(false, false, "stages", "create", machineStage1CreateString).run(t)
	cliTest(false, false, "stages", "create", machineStage2CreateString).run(t)
	cliTest(false, false, "plugins", "create", machinePluginCreateString).run(t)
	cliTest(true, false, "machines").run(t)
	cliTest(false, false, "machines", "list").run(t)
	cliTest(true, true, "machines", "create").run(t)
	cliTest(true, true, "machines", "create", "john", "john2").run(t)
	cliTest(false, true, "machines", "create", machineCreateBadJSONString).run(t)
	cliTest(false, true, "machines", "create", machineCreateBadJSON2String).run(t)
	cliTest(false, false, "machines", "create", machineCreateInputString).run(t)
	cliTest(false, true, "machines", "create", machineCreateInputString).run(t)
	cliTest(false, false, "machines", "list").run(t)
	cliTest(false, false, "machines", "list", "Name=fred").run(t)
	cliTest(false, false, "machines", "list", "Name=john").run(t)
	cliTest(false, false, "machines", "list", "BootEnv=local").run(t)
	cliTest(false, false, "machines", "list", "BootEnv=false").run(t)
	cliTest(false, false, "machines", "list", "Address=192.168.100.110").run(t)
	cliTest(false, false, "machines", "list", "Address=1.1.1.1").run(t)
	cliTest(false, true, "machines", "list", "Address=fred").run(t)
	cliTest(false, false, "machines", "list", "Uuid=4e7031fe-3062-45f1-835c-92541bc9cbd3").run(t)
	cliTest(false, false, "machines", "list", "Uuid=3e7031fe-3062-45f1-835c-92541bc9cbd3").run(t)
	cliTest(false, true, "machines", "list", "Uuid=false").run(t)
	cliTest(false, false, "machines", "list", "Runnable=true").run(t)
	cliTest(false, false, "machines", "list", "Runnable=false").run(t)
	cliTest(false, true, "machines", "list", "Runnable=fred").run(t)
	cliTest(true, true, "machines", "show").run(t)
	cliTest(true, true, "machines", "show", "john", "john2").run(t)
	cliTest(false, true, "machines", "show", "john").run(t)
	cliTest(false, false, "machines", "show", "3e7031fe-3062-45f1-835c-92541bc9cbd3").run(t)
	cliTest(false, false, "machines", "show", "Key:3e7031fe-3062-45f1-835c-92541bc9cbd3").run(t)
	cliTest(false, false, "machines", "show", "Uuid:3e7031fe-3062-45f1-835c-92541bc9cbd3").run(t)
	cliTest(false, false, "machines", "show", "Name:john").run(t)
	cliTest(true, true, "machines", "exists").run(t)
	cliTest(true, true, "machines", "exists", "john", "john2").run(t)
	cliTest(false, false, "machines", "exists", "3e7031fe-3062-45f1-835c-92541bc9cbd3").run(t)
	cliTest(false, true, "machines", "exists", "john").run(t)
	cliTest(true, true, "machines", "exists", "john", "john2").run(t)
	cliTest(true, true, "machines", "update").run(t)
	cliTest(true, true, "machines", "update", "john", "john2", "john3").run(t)
	cliTest(false, true, "machines", "update", "3e7031fe-3062-45f1-835c-92541bc9cbd3", machineUpdateBadJSONString).run(t)
	cliTest(false, false, "machines", "update", "3e7031fe-3062-45f1-835c-92541bc9cbd3", machineUpdateInputString).run(t)
	cliTest(false, true, "machines", "update", "john2", machineUpdateInputString).run(t)
	cliTest(false, false, "machines", "show", "3e7031fe-3062-45f1-835c-92541bc9cbd3").run(t)
	cliTest(false, false, "machines", "show", "3e7031fe-3062-45f1-835c-92541bc9cbd3").run(t)
	cliTest(true, true, "machines", "destroy").run(t)
	cliTest(true, true, "machines", "destroy", "john", "june").run(t)
	cliTest(false, false, "machines", "destroy", "3e7031fe-3062-45f1-835c-92541bc9cbd3").run(t)
	cliTest(false, true, "machines", "destroy", "3e7031fe-3062-45f1-835c-92541bc9cbd3").run(t)
	cliTest(false, false, "machines", "list").run(t)
	cliTest(false, false, "machines", "create", "-").Stdin(machineCreateInputString + "\n").run(t)
	cliTest(false, false, "machines", "list").run(t)
	cliTest(false, false, "machines", "update", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "-").Stdin(machineUpdateInputString + "\n").run(t)
	cliTest(false, false, "machines", "show", "3e7031fe-3062-45f1-835c-92541bc9cbd3").run(t)
	// bootenv tests
	cliTest(true, true, "machines", "bootenv").run(t)
	cliTest(false, true, "machines", "bootenv", "john", "john2").run(t)
	cliTest(false, false, "machines", "bootenv", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "john2").run(t)
	cliTest(false, false, "machines", "bootenv", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "local").run(t)
	cliTest(false, false, "machines", "update", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "{ \"Runnable\": true }").run(t)
	// stage tests
	cliTest(true, true, "machines", "stage").run(t)
	cliTest(false, true, "machines", "stage", "john", "john2").run(t)
	cliTest(false, true, "machines", "stage", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "john2").run(t)
	cliTest(false, false, "machines", "stage", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "stage1").run(t)
	cliTest(false, false, "machines", "update", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "{ \"Runnable\": true }").run(t)
	cliTest(false, true, "machines", "stage", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "stage2").run(t)
	cliTest(false, false, "machines", "stage", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "stage2", "--force").run(t)
	cliTest(false, false, "machines", "stage", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "", "--force").run(t)
	// Add/Remove Profile tests
	cliTest(true, true, "machines", "addprofile").run(t)
	cliTest(false, false, "machines", "addprofile", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "jill").run(t)
	cliTest(false, false, "machines", "addprofile", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "jean").run(t)
	cliTest(false, true, "machines", "addprofile", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "jill").run(t)
	cliTest(false, false, "profiles", "set", "jill", "param", "jill-param", "to", "janga").run(t)
	cliTest(false, false, "profiles", "set", "stage-prof", "param", "sp-param", "to", "val").run(t)
	cliTest(false, false, "machines", "stage", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "stage2", "--force").run(t)
	cliTest(false, false, "stages", "addprofile", "stage2", "stage-prof").run(t)

	cliTest(false, false, "machines", "params", "3e7031fe-3062-45f1-835c-92541bc9cbd3").run(t)
	cliTest(false, false, "machines", "params", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "--aggregate").run(t)
	tr := &http.Transport{}
	client := &http.Client{Transport: tr}
	req, _ := http.NewRequest("GET", "http://127.0.0.1:10002/machines/3e7031fe-3062-45f1-835c-92541bc9cbd3/file", nil)
	req.SetBasicAuth("rocketskates", "r0cketsk8ts")
	rsp, apierr := client.Do(req)
	if apierr != nil {
		t.Errorf("FAIL: Failed to query machine file: %s", apierr)
	} else {
		defer rsp.Body.Close()
		body, err := ioutil.ReadAll(rsp.Body)
		if err != nil {
			t.Errorf("FAIL: Failed to read all: %s", err)
		}
		if string(body) != "val" {
			t.Errorf("FAIL: Body was: AA%sAA expected %s", string(body), "val")
		}
	}

	cliTest(false, false, "machines", "stage", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "", "--force").run(t)
	cliTest(true, true, "machines", "removeprofile").run(t)
	cliTest(false, false, "machines", "removeprofile", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "justine").run(t)
	cliTest(false, false, "machines", "removeprofile", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "jill").run(t)
	cliTest(false, false, "machines", "removeprofile", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "jean").run(t)
	cliTest(true, true, "machines", "get").run(t)
	cliTest(false, true, "machines", "get", "john", "param", "john2").run(t)
	cliTest(false, false, "machines", "get", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "param", "john2").run(t)
	cliTest(true, true, "machines", "set").run(t)
	cliTest(false, true, "machines", "set", "john", "param", "john2", "to", "cow").run(t)
	cliTest(false, false, "machines", "set", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "param", "john2", "to", "cow").run(t)
	cliTest(false, false, "machines", "get", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "param", "john2").run(t)
	cliTest(false, false, "machines", "set", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "param", "john2", "to", "3").run(t)
	cliTest(false, false, "machines", "set", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "param", "john3", "to", "4").run(t)
	cliTest(false, false, "machines", "get", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "param", "john2").run(t)
	cliTest(false, false, "machines", "get", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "param", "john3").run(t)
	cliTest(false, false, "machines", "set", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "param", "john2", "to", "null").run(t)
	cliTest(false, false, "machines", "get", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "param", "john2").run(t)
	cliTest(false, false, "machines", "get", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "param", "john3").run(t)
	cliTest(true, true, "machines", "actions").run(t)
	cliTest(false, true, "machines", "actions", "john").run(t)
	cliTest(false, false, "machines", "actions", "3e7031fe-3062-45f1-835c-92541bc9cbd3").run(t)
	cliTest(true, true, "machines", "action").run(t)
	cliTest(true, true, "machines", "action", "john").run(t)
	cliTest(false, true, "machines", "action", "john", "command").run(t)
	cliTest(false, true, "machines", "action", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "command").run(t)
	cliTest(false, false, "machines", "action", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "increment").run(t)
	cliTest(true, true, "machines", "runaction").run(t)
	cliTest(true, true, "machines", "runaction", "fred").run(t)
	cliTest(false, true, "machines", "runaction", "fred", "command").run(t)
	cliTest(false, true, "machines", "runaction", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "command").run(t)
	cliTest(false, false, "machines", "runaction", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "increment").run(t)
	cliTest(false, true, "machines", "runaction", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "increment", "fred").run(t)

	cliTest(false, false, "machines", "actions", "3e7031fe-3062-45f1-835c-92541bc9cbd3").run(t)
	cliTest(false, false, "machines", "action", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "reset_count").run(t)
	cliTest(false, false, "machines", "runaction", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "reset_count").run(t)
	cliTest(false, false, "machines", "actions", "3e7031fe-3062-45f1-835c-92541bc9cbd3").run(t)
	cliTest(false, true, "machines", "action", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "reset_count").run(t)
	cliTest(false, true, "machines", "runaction", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "reset_count").run(t)
	cliTest(false, false, "machines", "runaction", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "increment", "incrementer/parameter", "asgdasdg").run(t)
	cliTest(false, false, "machines", "runaction", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "increment", "incrementer/parameter", "parm1", "extra", "10").run(t)
	cliTest(false, false, "machines", "get", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "param", "parm1").run(t)
	cliTest(false, true, "machines", "runaction", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "increment", "incrementer/parameter", "parm2", "incrementer/step", "asgdasdg").run(t)
	cliTest(false, false, "machines", "get", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "param", "parm2").run(t)
	cliTest(false, false, "machines", "runaction", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "increment", "incrementer/parameter", "parm2", "incrementer/step", "10").run(t)
	cliTest(false, false, "machines", "get", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "param", "parm2").run(t)
	cliTest(false, true, "machines", "runaction", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "increment", "-").Stdin("fred").run(t)
	cliTest(false, false, "machines", "runaction", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "reset_count", "-").Stdin(machineRunActionMissingParameterStdinString).run(t)
	cliTest(false, true, "machines", "runaction", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "reset_count", "-").Stdin(machineRunActionMissingParameterStdinString).run(t)
	cliTest(false, false, "machines", "runaction", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "increment", "-").Stdin(machineRunActionMissingParameterStdinString).run(t)
	cliTest(false, false, "machines", "runaction", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "increment", "-").Stdin(machineRunActionGoodStdinString).run(t)
	cliTest(false, false, "machines", "runaction", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "increment", "-").Stdin(machineRunActionGoodStdinString).run(t)
	cliTest(false, false, "machines", "get", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "param", "parm5").run(t)
	cliTest(true, true, "machines", "wait").run(t)
	cliTest(true, true, "machines", "wait", "jk").run(t)
	cliTest(true, true, "machines", "wait", "jk", "jk").run(t)
	cliTest(true, true, "machines", "wait", "jk", "jk", "jk", "jk", "jk").run(t)
	cliTest(false, true, "machines", "wait", "jk", "jk", "jk", "jk").run(t)
	cliTest(false, true, "machines", "wait", "jk", "jk", "jk").run(t)
	cliTest(false, false, "machines", "wait", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "jk", "jk", "1").run(t)
	cliTest(false, false, "machines", "wait", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "BootEnv", "local", "1").run(t)
	cliTest(false, false, "machines", "wait", "3e7031fe-3062-45f1-835c-92541bc9cbd3", "Runnable", "fred", "1").run(t)
	cliTest(true, true, "machines", "params").run(t)
	cliTest(false, true, "machines", "params", "john2").run(t)
	cliTest(false, false, "machines", "params", "3e7031fe-3062-45f1-835c-92541bc9cbd3").run(t)
	cliTest(false, true, "machines", "params", "john2", machinesParamsNextString).run(t)
	cliTest(false, false, "machines", "params", "3e7031fe-3062-45f1-835c-92541bc9cbd3", machinesParamsNextString).run(t)
	cliTest(false, false, "machines", "params", "3e7031fe-3062-45f1-835c-92541bc9cbd3").run(t)
	cliTest(false, false, "machines", "show", "3e7031fe-3062-45f1-835c-92541bc9cbd3").run(t)
	cliTest(false, false, "machines", "destroy", "3e7031fe-3062-45f1-835c-92541bc9cbd3").run(t)
	cliTest(false, false, "machines", "list").run(t)
	cliTest(false, false, "prefs", "set", "defaultStage", "stage1").run(t)
	cliTest(false, false, "machines", "create", machineCreateInputString).run(t)
	cliTest(false, false, "machines", "destroy", "3e7031fe-3062-45f1-835c-92541bc9cbd3").run(t)
	cliTest(false, false, "machines", "list").run(t)
	cliTest(false, false, "prefs", "set", "defaultStage", "none").run(t)
	cliTest(false, false, "plugins", "destroy", "incr").run(t)
	cliTest(false, false, "stages", "destroy", "stage1").run(t)
	cliTest(false, false, "stages", "destroy", "stage2").run(t)
	cliTest(false, false, "profiles", "destroy", "jill").run(t)
	cliTest(false, false, "profiles", "destroy", "jean").run(t)
	cliTest(false, false, "profiles", "destroy", "stage-prof").run(t)
	cliTest(false, false, "tasks", "destroy", "jamie").run(t)
	cliTest(false, false, "tasks", "destroy", "justine").run(t)
}
