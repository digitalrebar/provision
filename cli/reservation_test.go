package cli

import (
	"testing"
)

var reservationAddrErrorString string = "Error: GET: reservations: Invalid Address: fred\n\n"
var reservationExpireTimeErrorString string = "Error: GET: reservations: Invalid Address: false\n\n"
var reservationShowNoArgErrorString string = "Error: drpcli reservations show [id] [flags] requires 1 argument\n"
var reservationShowTooManyArgErrorString string = "Error: drpcli reservations show [id] [flags] requires 1 argument\n"
var reservationShowMissingArgErrorString string = "Error: GET: reservations/C0A86467: Not Found\n\n"
var reservationExistsNoArgErrorString string = "Error: drpcli reservations exists [id] [flags] requires 1 argument"
var reservationExistsTooManyArgErrorString string = "Error: drpcli reservations exists [id] [flags] requires 1 argument"
var reservationExistsMissingIgnoreString string = "Error: GET: reservation get: address not valid: ignore\n\n"
var reservationCreateNoArgErrorString string = "Error: drpcli reservations create [json] [flags] requires 1 argument\n"
var reservationCreateTooManyArgErrorString string = "Error: drpcli reservations create [json] [flags] requires 1 argument\n"
var reservationCreateBadJSONErrorString = "Error: Unable to create new reservation: Invalid type passed to reservation create\n\n"
var reservationCreateDuplicateErrorString = "Error: CREATE: reservations/C0A86464: already exists\n\n"
var reservationUpdateNoArgErrorString string = "Error: drpcli reservations update [id] [json] [flags] requires 2 arguments"
var reservationUpdateTooManyArgErrorString string = "Error: drpcli reservations update [id] [json] [flags] requires 2 arguments"
var reservationUpdateBadJSONErrorString = "Error: Unable to merge objects: json: cannot unmarshal string into Go value of type map[string]interface {}\n\n\n"
var reservationUpdateJohnMissingErrorString string = "Error: GET: reservations/C0A86467: Not Found\n\n"
var reservationPatchNoArgErrorString string = "Error: drpcli reservations patch [objectJson] [changesJson] [flags] requires 2 arguments"
var reservationPatchTooManyArgErrorString string = "Error: drpcli reservations patch [objectJson] [changesJson] [flags] requires 2 arguments"
var reservationPatchJohnMissingErrorString string = "Error: PATCH: reservations/C1A86464: Not Found\n\n"
var reservationDestroyNoArgErrorString string = "Error: drpcli reservations destroy [id] [flags] requires 1 argument"
var reservationDestroyTooManyArgErrorString string = "Error: drpcli reservations destroy [id] [flags] requires 1 argument"
var reservationDestroyMissingJohnString string = "Error: DELETE: reservations/C0A86464: Not Found\n\n"

var reservationDefaultListString string = "[]\n"
var reservationEmptyListString string = "[]\n"

var reservationShowJohnString string = `{
  "Addr": "192.168.100.100",
  "Available": true,
  "Errors": [],
  "NextServer": "2.2.2.2",
  "Options": [],
  "ReadOnly": false,
  "Strategy": "MAC",
  "Token": "john",
  "Validated": true
}
`

var reservationExistsIgnoreString string = ""

var reservationCreateBadJSONString = "asdgasdg"

var reservationCreateInputString string = `{
  "Addr": "192.168.100.100",
  "NextServer": "2.2.2.2",
  "Strategy": "MAC",
  "Token": "john"
}
`
var reservationCreateJohnString string = `{
  "Addr": "192.168.100.100",
  "Available": true,
  "Errors": [],
  "NextServer": "2.2.2.2",
  "Options": [],
  "ReadOnly": false,
  "Strategy": "MAC",
  "Token": "john",
  "Validated": true
}
`

var reservationListReservationsString = `[
  {
    "Addr": "192.168.100.100",
    "Available": true,
    "Errors": [],
    "NextServer": "2.2.2.2",
    "Options": [],
    "ReadOnly": false,
    "Strategy": "MAC",
    "Token": "john",
    "Validated": true
  }
]
`
var reservationListBothEnvsString = `[
  {
    "Addr": "192.168.100.100",
    "Available": true,
    "Errors": [],
    "NextServer": "2.2.2.2",
    "Options": [],
    "ReadOnly": false,
    "Strategy": "MAC",
    "Token": "john",
    "Validated": true
  }
]
`

var reservationUpdateBadJSONString = "asdgasdg"

var reservationUpdateInputString string = `{
  "Options": [ { "Code": 3, "Value": "1.1.1.1" } ]
}
`
var reservationUpdateJohnString string = `{
  "Addr": "192.168.100.100",
  "Available": true,
  "Errors": [],
  "NextServer": "2.2.2.2",
  "Options": [
    {
      "Code": 3,
      "Value": "1.1.1.1"
    }
  ],
  "ReadOnly": false,
  "Strategy": "MAC",
  "Token": "john",
  "Validated": true
}
`

var reservationPatchBadPatchJSONString = "asdgasdg"
var reservationPatchBadPatchJSONErrorString = "Error: Unable to parse drpcli reservations patch [objectJson] [changesJson] [flags] JSON asdgasdg\nError: error unmarshaling JSON: json: cannot unmarshal string into Go value of type genmodels.Reservation\n\n"
var reservationPatchBadBaseJSONString = "asdgasdg"
var reservationPatchBadBaseJSONErrorString = "Error: Unable to parse drpcli reservations patch [objectJson] [changesJson] [flags] JSON asdgasdg\nError: error unmarshaling JSON: json: cannot unmarshal string into Go value of type genmodels.Reservation\n\n"
var reservationPatchBaseString string = `{
  "Addr": "192.168.100.100",
  "Strategy": "MAC",
  "Token": "john"
}
`
var reservationPatchInputString string = `{
  "Options": [ { "Code": 3, "Value": "1.1.3.1" } ]
}
`
var reservationPatchJohnString string = `{
  "Addr": "192.168.100.100",
  "Available": true,
  "Errors": [],
  "NextServer": "2.2.2.2",
  "Options": [
    {
      "Code": 3,
      "Value": "1.1.3.1"
    }
  ],
  "ReadOnly": false,
  "Strategy": "MAC",
  "Token": "john",
  "Validated": true
}
`
var reservationPatchMissingBaseString string = `{
  "Addr": "193.168.100.100",
  "NextServer": "2.2.2.2",
  "Strategy": "NewStrat",
  "Token": "john"
}
`

var reservationDestroyJohnString string = "Deleted reservation 192.168.100.100\n"

func TestReservationCli(t *testing.T) {
	cliTest(true, false, "reservations").run(t)
	cliTest(false, false, "reservations", "list").run(t)
	cliTest(true, true, "reservations", "create").run(t)
	cliTest(true, true, "reservations", "create", "john", "john2").run(t)
	cliTest(false, true, "reservations", "create", reservationCreateBadJSONString).run(t)
	cliTest(false, false, "reservations", "create", reservationCreateInputString).run(t)
	cliTest(false, true, "reservations", "create", reservationCreateInputString).run(t)
	cliTest(false, false, "reservations", "list").run(t)
	cliTest(false, false, "reservations", "list", "Strategy=fred").run(t)
	cliTest(false, false, "reservations", "list", "Strategy=MAC").run(t)
	cliTest(false, false, "reservations", "list", "Token=john").run(t)
	cliTest(false, false, "reservations", "list", "Token=false").run(t)
	cliTest(false, false, "reservations", "list", "Addr=192.168.100.100").run(t)
	cliTest(false, false, "reservations", "list", "Addr=1.1.1.1").run(t)
	cliTest(false, true, "reservations", "list", "Addr=fred").run(t)
	cliTest(false, false, "reservations", "list", "NextServer=3.3.3.3").run(t)
	cliTest(false, false, "reservations", "list", "NextServer=2.2.2.2").run(t)
	cliTest(false, true, "reservations", "list", "NextServer=false").run(t)
	cliTest(true, true, "reservations", "show").run(t)
	cliTest(true, true, "reservations", "show", "john", "john2").run(t)
	cliTest(false, true, "reservations", "show", "192.168.100.103").run(t)
	cliTest(false, false, "reservations", "show", "192.168.100.100").run(t)
	cliTest(true, true, "reservations", "exists").run(t)
	cliTest(true, true, "reservations", "exists", "john", "john2").run(t)
	cliTest(false, false, "reservations", "exists", "192.168.100.100").run(t)
	cliTest(false, true, "reservations", "exists", "ignore").run(t)
	cliTest(true, true, "reservations", "exists", "john", "john2").run(t)
	cliTest(true, true, "reservations", "update").run(t)
	cliTest(true, true, "reservations", "update", "john", "john2", "john3").run(t)
	cliTest(false, true, "reservations", "update", "192.168.100.100", reservationUpdateBadJSONString).run(t)
	cliTest(false, false, "reservations", "update", "192.168.100.100", reservationUpdateInputString).run(t)
	cliTest(false, true, "reservations", "update", "192.168.100.103", reservationUpdateInputString).run(t)
	cliTest(false, false, "reservations", "show", "192.168.100.100").run(t)
	cliTest(false, false, "reservations", "show", "192.168.100.100").run(t)
	cliTest(true, true, "reservations", "destroy").run(t)
	cliTest(true, true, "reservations", "destroy", "john", "june").run(t)
	cliTest(false, false, "reservations", "destroy", "192.168.100.100").run(t)
	cliTest(false, true, "reservations", "destroy", "192.168.100.100").run(t)
	cliTest(false, false, "reservations", "list").run(t)
	cliTest(false, false, "reservations", "create", "-").Stdin(reservationCreateInputString + "\n").run(t)
	cliTest(false, false, "reservations", "list").run(t)
	cliTest(false, false, "reservations", "update", "192.168.100.100", "-").Stdin(reservationUpdateInputString + "\n").run(t)
	cliTest(false, false, "reservations", "show", "192.168.100.100").run(t)
	cliTest(false, false, "reservations", "destroy", "192.168.100.100").run(t)
	cliTest(false, false, "reservations", "list").run(t)

}
