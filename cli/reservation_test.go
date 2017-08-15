package cli

import (
	"testing"
)

var reservationAddrErrorString string = "Error: Invalid Address: fred\n\n"
var reservationExpireTimeErrorString string = "Error: Invalid Address: false\n\n"

var reservationDefaultListString string = "[]\n"
var reservationEmptyListString string = "[]\n"

var reservationShowNoArgErrorString string = "Error: drpcli reservations show [id] [flags] requires 1 argument\n"
var reservationShowTooManyArgErrorString string = "Error: drpcli reservations show [id] [flags] requires 1 argument\n"
var reservationShowMissingArgErrorString string = "Error: reservations GET: C0A86467: Not Found\n\n"
var reservationShowJohnString string = `{
  "Addr": "192.168.100.100",
  "NextServer": "2.2.2.2",
  "Options": null,
  "Strategy": "MAC",
  "Token": "john"
}
`

var reservationExistsNoArgErrorString string = "Error: drpcli reservations exists [id] [flags] requires 1 argument"
var reservationExistsTooManyArgErrorString string = "Error: drpcli reservations exists [id] [flags] requires 1 argument"
var reservationExistsIgnoreString string = ""
var reservationExistsMissingIgnoreString string = "Error: reservation get: address not valid: ignore\n\n"

var reservationCreateNoArgErrorString string = "Error: drpcli reservations create [json] [flags] requires 1 argument\n"
var reservationCreateTooManyArgErrorString string = "Error: drpcli reservations create [json] [flags] requires 1 argument\n"
var reservationCreateBadJSONString = "asdgasdg"
var reservationCreateBadJSONErrorString = "Error: Unable to create new reservation: Invalid type passed to reservation create\n\n"
var reservationCreateInputString string = `{
  "Addr": "192.168.100.100",
  "NextServer": "2.2.2.2",
  "Strategy": "MAC",
  "Token": "john"
}
`
var reservationCreateJohnString string = `{
  "Addr": "192.168.100.100",
  "NextServer": "2.2.2.2",
  "Options": null,
  "Strategy": "MAC",
  "Token": "john"
}
`
var reservationCreateDuplicateErrorString = "Error: dataTracker create reservations: C0A86464 already exists\n\n"

var reservationListReservationsString = `[
  {
    "Addr": "192.168.100.100",
    "NextServer": "2.2.2.2",
    "Options": null,
    "Strategy": "MAC",
    "Token": "john"
  }
]
`
var reservationListBothEnvsString = `[
  {
    "Addr": "192.168.100.100",
    "NextServer": "2.2.2.2",
    "Options": null,
    "Strategy": "MAC",
    "Token": "john"
  }
]
`

var reservationUpdateNoArgErrorString string = "Error: drpcli reservations update [id] [json] [flags] requires 2 arguments"
var reservationUpdateTooManyArgErrorString string = "Error: drpcli reservations update [id] [json] [flags] requires 2 arguments"
var reservationUpdateBadJSONString = "asdgasdg"
var reservationUpdateBadJSONErrorString = "Error: Unable to merge objects: json: cannot unmarshal string into Go value of type map[string]interface {}\n\n\n"
var reservationUpdateInputString string = `{
  "Options": [ { "Code": 3, "Value": "1.1.1.1" } ]
}
`
var reservationUpdateJohnString string = `{
  "Addr": "192.168.100.100",
  "NextServer": "2.2.2.2",
  "Options": [
    {
      "Code": 3,
      "Value": "1.1.1.1"
    }
  ],
  "Strategy": "MAC",
  "Token": "john"
}
`
var reservationUpdateJohnMissingErrorString string = "Error: reservations GET: C0A86467: Not Found\n\n"

var reservationPatchNoArgErrorString string = "Error: drpcli reservations patch [objectJson] [changesJson] [flags] requires 2 arguments"
var reservationPatchTooManyArgErrorString string = "Error: drpcli reservations patch [objectJson] [changesJson] [flags] requires 2 arguments"
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
  "NextServer": "2.2.2.2",
  "Options": [
    {
      "Code": 3,
      "Value": "1.1.3.1"
    }
  ],
  "Strategy": "MAC",
  "Token": "john"
}
`
var reservationPatchMissingBaseString string = `{
  "Addr": "193.168.100.100",
  "NextServer": "2.2.2.2",
  "Strategy": "NewStrat",
  "Token": "john"
}
`
var reservationPatchJohnMissingErrorString string = "Error: reservations: PATCH C1A86464: Not Found\n\n"

var reservationDestroyNoArgErrorString string = "Error: drpcli reservations destroy [id] [flags] requires 1 argument"
var reservationDestroyTooManyArgErrorString string = "Error: drpcli reservations destroy [id] [flags] requires 1 argument"
var reservationDestroyJohnString string = "Deleted reservation 192.168.100.100\n"
var reservationDestroyMissingJohnString string = "Error: reservations: DELETE C0A86464: Not Found\n\n"

func TestReservationCli(t *testing.T) {
	tests := []CliTest{
		CliTest{true, false, []string{"reservations"}, noStdinString, "Access CLI commands relating to reservations\n", ""},
		CliTest{false, false, []string{"reservations", "list"}, noStdinString, reservationDefaultListString, noErrorString},

		CliTest{true, true, []string{"reservations", "create"}, noStdinString, noContentString, reservationCreateNoArgErrorString},
		CliTest{true, true, []string{"reservations", "create", "john", "john2"}, noStdinString, noContentString, reservationCreateTooManyArgErrorString},
		CliTest{false, true, []string{"reservations", "create", reservationCreateBadJSONString}, noStdinString, noContentString, reservationCreateBadJSONErrorString},
		CliTest{false, false, []string{"reservations", "create", reservationCreateInputString}, noStdinString, reservationCreateJohnString, noErrorString},
		CliTest{false, true, []string{"reservations", "create", reservationCreateInputString}, noStdinString, noContentString, reservationCreateDuplicateErrorString},
		CliTest{false, false, []string{"reservations", "list"}, noStdinString, reservationListBothEnvsString, noErrorString},

		CliTest{false, false, []string{"reservations", "list", "--limit=0"}, noStdinString, reservationEmptyListString, noErrorString},
		CliTest{false, false, []string{"reservations", "list", "--limit=10", "--offset=0"}, noStdinString, reservationListReservationsString, noErrorString},
		CliTest{false, false, []string{"reservations", "list", "--limit=10", "--offset=10"}, noStdinString, reservationEmptyListString, noErrorString},
		CliTest{false, true, []string{"reservations", "list", "--limit=-10", "--offset=0"}, noStdinString, noContentString, limitNegativeError},
		CliTest{false, true, []string{"reservations", "list", "--limit=10", "--offset=-10"}, noStdinString, noContentString, offsetNegativeError},
		CliTest{false, false, []string{"reservations", "list", "--limit=-1", "--offset=-1"}, noStdinString, reservationListReservationsString, noErrorString},
		CliTest{false, false, []string{"reservations", "list", "Strategy=fred"}, noStdinString, reservationEmptyListString, noErrorString},
		CliTest{false, false, []string{"reservations", "list", "Strategy=MAC"}, noStdinString, reservationListReservationsString, noErrorString},
		CliTest{false, false, []string{"reservations", "list", "Token=john"}, noStdinString, reservationListReservationsString, noErrorString},
		CliTest{false, false, []string{"reservations", "list", "Token=false"}, noStdinString, reservationEmptyListString, noErrorString},
		CliTest{false, false, []string{"reservations", "list", "Addr=192.168.100.100"}, noStdinString, reservationListReservationsString, noErrorString},
		CliTest{false, false, []string{"reservations", "list", "Addr=1.1.1.1"}, noStdinString, reservationEmptyListString, noErrorString},
		CliTest{false, true, []string{"reservations", "list", "Addr=fred"}, noStdinString, noContentString, reservationAddrErrorString},
		CliTest{false, false, []string{"reservations", "list", "NextServer=3.3.3.3"}, noStdinString, reservationEmptyListString, noErrorString},
		CliTest{false, false, []string{"reservations", "list", "NextServer=2.2.2.2"}, noStdinString, reservationListReservationsString, noErrorString},
		CliTest{false, true, []string{"reservations", "list", "NextServer=false"}, noStdinString, noContentString, reservationExpireTimeErrorString},

		CliTest{true, true, []string{"reservations", "show"}, noStdinString, noContentString, reservationShowNoArgErrorString},
		CliTest{true, true, []string{"reservations", "show", "john", "john2"}, noStdinString, noContentString, reservationShowTooManyArgErrorString},
		CliTest{false, true, []string{"reservations", "show", "192.168.100.103"}, noStdinString, noContentString, reservationShowMissingArgErrorString},
		CliTest{false, false, []string{"reservations", "show", "192.168.100.100"}, noStdinString, reservationShowJohnString, noErrorString},

		CliTest{true, true, []string{"reservations", "exists"}, noStdinString, noContentString, reservationExistsNoArgErrorString},
		CliTest{true, true, []string{"reservations", "exists", "john", "john2"}, noStdinString, noContentString, reservationExistsTooManyArgErrorString},
		CliTest{false, false, []string{"reservations", "exists", "192.168.100.100"}, noStdinString, reservationExistsIgnoreString, noErrorString},
		CliTest{false, true, []string{"reservations", "exists", "ignore"}, noStdinString, noContentString, reservationExistsMissingIgnoreString},
		CliTest{true, true, []string{"reservations", "exists", "john", "john2"}, noStdinString, noContentString, reservationExistsTooManyArgErrorString},

		CliTest{true, true, []string{"reservations", "update"}, noStdinString, noContentString, reservationUpdateNoArgErrorString},
		CliTest{true, true, []string{"reservations", "update", "john", "john2", "john3"}, noStdinString, noContentString, reservationUpdateTooManyArgErrorString},
		CliTest{false, true, []string{"reservations", "update", "192.168.100.100", reservationUpdateBadJSONString}, noStdinString, noContentString, reservationUpdateBadJSONErrorString},
		CliTest{false, false, []string{"reservations", "update", "192.168.100.100", reservationUpdateInputString}, noStdinString, reservationUpdateJohnString, noErrorString},
		CliTest{false, true, []string{"reservations", "update", "192.168.100.103", reservationUpdateInputString}, noStdinString, noContentString, reservationUpdateJohnMissingErrorString},
		CliTest{false, false, []string{"reservations", "show", "192.168.100.100"}, noStdinString, reservationUpdateJohnString, noErrorString},

		CliTest{true, true, []string{"reservations", "patch"}, noStdinString, noContentString, reservationPatchNoArgErrorString},
		CliTest{true, true, []string{"reservations", "patch", "john", "john2", "john3"}, noStdinString, noContentString, reservationPatchTooManyArgErrorString},
		CliTest{false, true, []string{"reservations", "patch", reservationPatchBaseString, reservationPatchBadPatchJSONString}, noStdinString, noContentString, reservationPatchBadPatchJSONErrorString},
		CliTest{false, true, []string{"reservations", "patch", reservationPatchBadBaseJSONString, reservationPatchInputString}, noStdinString, noContentString, reservationPatchBadBaseJSONErrorString},
		CliTest{false, false, []string{"reservations", "patch", reservationPatchBaseString, reservationPatchInputString}, noStdinString, reservationPatchJohnString, noErrorString},
		CliTest{false, true, []string{"reservations", "patch", reservationPatchMissingBaseString, reservationPatchInputString}, noStdinString, noContentString, reservationPatchJohnMissingErrorString},
		CliTest{false, false, []string{"reservations", "show", "192.168.100.100"}, noStdinString, reservationPatchJohnString, noErrorString},

		CliTest{true, true, []string{"reservations", "destroy"}, noStdinString, noContentString, reservationDestroyNoArgErrorString},
		CliTest{true, true, []string{"reservations", "destroy", "john", "june"}, noStdinString, noContentString, reservationDestroyTooManyArgErrorString},
		CliTest{false, false, []string{"reservations", "destroy", "192.168.100.100"}, noStdinString, reservationDestroyJohnString, noErrorString},
		CliTest{false, true, []string{"reservations", "destroy", "192.168.100.100"}, noStdinString, noContentString, reservationDestroyMissingJohnString},
		CliTest{false, false, []string{"reservations", "list"}, noStdinString, reservationDefaultListString, noErrorString},

		CliTest{false, false, []string{"reservations", "create", "-"}, reservationCreateInputString + "\n", reservationCreateJohnString, noErrorString},
		CliTest{false, false, []string{"reservations", "list"}, noStdinString, reservationListBothEnvsString, noErrorString},
		CliTest{false, false, []string{"reservations", "update", "192.168.100.100", "-"}, reservationUpdateInputString + "\n", reservationUpdateJohnString, noErrorString},
		CliTest{false, false, []string{"reservations", "show", "192.168.100.100"}, noStdinString, reservationUpdateJohnString, noErrorString},

		CliTest{false, false, []string{"reservations", "destroy", "192.168.100.100"}, noStdinString, reservationDestroyJohnString, noErrorString},
		CliTest{false, false, []string{"reservations", "list"}, noStdinString, reservationDefaultListString, noErrorString},
	}

	for _, test := range tests {
		testCli(t, test)
	}

}
