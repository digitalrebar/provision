package cli

import (
	"testing"
)

var leaseDefaultListString string = "[]\n"
var leaseEmptyListString string = "[]\n"

var leaseAddrErrorString string = "Error: Addr must be an IP address\n\n"
var leaseExpireTimeErrorString string = "Error: ExpireTime is not valid: parsing time \"false\" as \"2006-01-02T15:04:05Z07:00\": cannot parse \"false\" as \"2006\"\n\n"

var leaseShowNoArgErrorString string = "Error: drpcli leases show [id] [flags] requires 1 argument\n"
var leaseShowTooManyArgErrorString string = "Error: drpcli leases show [id] [flags] requires 1 argument\n"
var leaseShowMissingArgErrorString string = "Error: leases GET: C0A8646F: Not Found\n\n"
var leaseShowLeaseString string = `{
  "Addr": "192.168.100.110",
  "Available": true,
  "Errors": [],
  "ExpireTime": "2017-03-31T00:11:21.028-05:00",
  "ReadOnly": false,
  "State": "",
  "Strategy": "MAC",
  "Token": "08:00:27:33:77:de",
  "Validated": true
}
`

var leaseExistsNoArgErrorString string = "Error: drpcli leases exists [id] [flags] requires 1 argument"
var leaseExistsTooManyArgErrorString string = "Error: drpcli leases exists [id] [flags] requires 1 argument"
var leaseExistsLeaseString string = ""
var leaseExistsMissingJohnString string = "Error: leases GET: C0A8646F: Not Found\n\n"

var leaseCreateNoArgErrorString string = "Error: drpcli leases create [json] [flags] requires 1 argument\n"
var leaseCreateTooManyArgErrorString string = "Error: drpcli leases create [json] [flags] requires 1 argument\n"
var leaseCreateBadJSONString = "asdgasdg"
var leaseCreateBadJSONErrorString = "Error: Unable to create new lease: Invalid type passed to lease create\n\n"
var leaseCreateInputString string = `{
  "Addr": "192.168.100.110",
  "ExpireTime": "2017-03-31T00:11:21.028-05:00",
  "Strategy": "MAC",
  "Token": "08:00:27:33:77:de"
}
`
var leaseCreateJohnString string = `{
  "Addr": "192.168.100.110",
  "Available": true,
  "Errors": [],
  "ExpireTime": "2017-03-31T00:11:21.028-05:00",
  "ReadOnly": false,
  "State": "",
  "Strategy": "MAC",
  "Token": "08:00:27:33:77:de",
  "Validated": true
}
`
var leaseCreateDuplicateErrorString = "Error: dataTracker create leases: C0A8646E already exists\n\n"

var leaseListLeasesString = `[
  {
    "Addr": "192.168.100.110",
    "Available": true,
    "Errors": [],
    "ExpireTime": "2017-03-31T00:11:21.028-05:00",
    "ReadOnly": false,
    "State": "",
    "Strategy": "MAC",
    "Token": "08:00:27:33:77:de",
    "Validated": true
  }
]
`

var leaseUpdateNoArgErrorString string = "Error: drpcli leases update [id] [json] [flags] requires 2 arguments"
var leaseUpdateTooManyArgErrorString string = "Error: drpcli leases update [id] [json] [flags] requires 2 arguments"
var leaseUpdateBadJSONString = "asdgasdg"
var leaseUpdateBadJSONErrorString = "Error: Unable to merge objects: json: cannot unmarshal string into Go value of type map[string]interface {}\n\n\n"
var leaseUpdateInputString string = `{
  "ExpireTime": "2019-03-31T00:11:21.028-05:00"
}
`
var leaseUpdateJohnString string = `{
  "Addr": "192.168.100.110",
  "Available": true,
  "Errors": [],
  "ExpireTime": "2019-03-31T00:11:21.028-05:00",
  "ReadOnly": false,
  "State": "",
  "Strategy": "MAC",
  "Token": "08:00:27:33:77:de",
  "Validated": true
}
`
var leaseUpdateJohnMissingErrorString string = "Error: leases GET: C0A8646F: Not Found\n\n"

var leasePatchNoArgErrorString string = "Error: drpcli leases patch [objectJson] [changesJson] [flags] requires 2 arguments"
var leasePatchTooManyArgErrorString string = "Error: drpcli leases patch [objectJson] [changesJson] [flags] requires 2 arguments"
var leasePatchBadPatchJSONString = "asdgasdg"
var leasePatchBadPatchJSONErrorString = "Error: Unable to parse drpcli leases patch [objectJson] [changesJson] [flags] JSON asdgasdg\nError: error unmarshaling JSON: json: cannot unmarshal string into Go value of type genmodels.Lease\n\n"
var leasePatchBadBaseJSONString = "asdgasdg"
var leasePatchBadBaseJSONErrorString = "Error: Unable to parse drpcli leases patch [objectJson] [changesJson] [flags] JSON asdgasdg\nError: error unmarshaling JSON: json: cannot unmarshal string into Go value of type genmodels.Lease\n\n"
var leasePatchBaseString string = `{
  "Addr": "192.168.100.110",
  "Available": true,
  "Errors": [],
  "ExpireTime": "2019-03-31T00:11:21.028-05:00",
  "State": "",
  "Strategy": "MAC",
  "Token": "08:00:27:33:77:de",
  "Validated": true
}
`
var leasePatchInputString string = `{
  "ExpireTime": "2018-03-31T00:11:21.028-05:00",
}
`
var leasePatchJohnString string = `{
  "Addr": "192.168.100.110",
  "Available": true,
  "Errors": [],
  "ExpireTime": "2018-03-31T00:11:21.028-05:00",
  "ReadOnly": false,
  "State": "",
  "Strategy": "MAC",
  "Token": "08:00:27:33:77:de",
  "Validated": true
}
`
var leasePatchMissingBaseString string = `{
  "Addr": "192.168.100.111",
  "Available": true,
  "Errors": [],
  "ExpireTime": "2018-03-31T00:11:21.028-05:00",
  "Strategy": "MAC",
  "Token": "08:00:27:33:77:de",
  "Validated": true
}
`
var leasePatchJohnMissingErrorString string = "Error: leases: PATCH C0A8646F: Not Found\n\n"

var leaseDestroyNoArgErrorString string = "Error: drpcli leases destroy [id] [flags] requires 1 argument"
var leaseDestroyTooManyArgErrorString string = "Error: drpcli leases destroy [id] [flags] requires 1 argument"
var leaseDestroyJohnString string = "Deleted lease 192.168.100.110\n"
var leaseDestroyMissingJohnString string = "Error: leases: DELETE C0A8646E: Not Found\n\n"

var leaseShowInvalidAddressErrorString string = "Error: lease get: address not valid: k192.168.100.110\n\n"
var leaseUpdateInvalidAddressErrorString string = "Error: lease get: address not valid: k192.168.100.111\n\n"
var leaseDestroyInvalidAddressErrorString string = "Error: lease delete: address not valid: k192.168.100.110\n\n"

func TestLeaseCli(t *testing.T) {
	tests := []CliTest{
		// Create subnet
		CliTest{false, false, []string{"subnets", "create", subnetCreateInputString}, noStdinString, subnetCreateJohnString, noErrorString},

		CliTest{true, false, []string{"leases"}, noStdinString, "Access CLI commands relating to leases\n", ""},
		CliTest{false, false, []string{"leases", "list"}, noStdinString, leaseDefaultListString, noErrorString},

		CliTest{true, true, []string{"leases", "create"}, noStdinString, noContentString, leaseCreateNoArgErrorString},
		CliTest{true, true, []string{"leases", "create", "john", "john2"}, noStdinString, noContentString, leaseCreateTooManyArgErrorString},
		CliTest{false, true, []string{"leases", "create", leaseCreateBadJSONString}, noStdinString, noContentString, leaseCreateBadJSONErrorString},
		CliTest{false, false, []string{"leases", "create", leaseCreateInputString}, noStdinString, leaseCreateJohnString, noErrorString},
		CliTest{false, true, []string{"leases", "create", leaseCreateInputString}, noStdinString, noContentString, leaseCreateDuplicateErrorString},
		CliTest{false, false, []string{"leases", "list"}, noStdinString, leaseListLeasesString, noErrorString},
		CliTest{false, false, []string{"leases", "list", "--limit=0"}, noStdinString, leaseEmptyListString, noErrorString},
		CliTest{false, false, []string{"leases", "list", "--limit=10", "--offset=0"}, noStdinString, leaseListLeasesString, noErrorString},
		CliTest{false, false, []string{"leases", "list", "--limit=10", "--offset=10"}, noStdinString, leaseEmptyListString, noErrorString},
		CliTest{false, true, []string{"leases", "list", "--limit=-10", "--offset=0"}, noStdinString, noContentString, limitNegativeError},
		CliTest{false, true, []string{"leases", "list", "--limit=10", "--offset=-10"}, noStdinString, noContentString, offsetNegativeError},
		CliTest{false, false, []string{"leases", "list", "--limit=-1", "--offset=-1"}, noStdinString, leaseListLeasesString, noErrorString},
		CliTest{false, false, []string{"leases", "list", "Strategy=fred"}, noStdinString, leaseEmptyListString, noErrorString},
		CliTest{false, false, []string{"leases", "list", "Strategy=MAC"}, noStdinString, leaseListLeasesString, noErrorString},
		CliTest{false, false, []string{"leases", "list", "Token=08:00:27:33:77:de"}, noStdinString, leaseListLeasesString, noErrorString},
		CliTest{false, false, []string{"leases", "list", "Token=false"}, noStdinString, leaseEmptyListString, noErrorString},
		CliTest{false, false, []string{"leases", "list", "Addr=192.168.100.110"}, noStdinString, leaseListLeasesString, noErrorString},
		CliTest{false, false, []string{"leases", "list", "Addr=1.1.1.1"}, noStdinString, leaseEmptyListString, noErrorString},
		CliTest{false, true, []string{"leases", "list", "Addr=fred"}, noStdinString, noContentString, leaseAddrErrorString},
		CliTest{false, false, []string{"leases", "list", "ExpireTime=2016-03-31T00:11:21.028-05:00"}, noStdinString, leaseEmptyListString, noErrorString},
		CliTest{false, false, []string{"leases", "list", "ExpireTime=2017-03-31T00:11:21.028-05:00"}, noStdinString, leaseListLeasesString, noErrorString},
		CliTest{false, true, []string{"leases", "list", "ExpireTime=false"}, noStdinString, noContentString, leaseExpireTimeErrorString},
		CliTest{false, false, []string{"leases", "list", "Available=true"}, noStdinString, leaseListLeasesString, noErrorString},
		CliTest{false, false, []string{"leases", "list", "Available=false"}, noStdinString, leaseEmptyListString, noErrorString},
		CliTest{false, true, []string{"leases", "list", "Available=fred"}, noStdinString, noContentString, bootEnvBadAvailableString},
		CliTest{false, false, []string{"leases", "list", "Valid=true"}, noStdinString, leaseListLeasesString, noErrorString},
		CliTest{false, false, []string{"leases", "list", "Valid=false"}, noStdinString, leaseEmptyListString, noErrorString},
		CliTest{false, true, []string{"leases", "list", "Valid=fred"}, noStdinString, noContentString, bootEnvBadValidString},
		CliTest{false, false, []string{"leases", "list", "ReadOnly=true"}, noStdinString, leaseEmptyListString, noErrorString},
		CliTest{false, false, []string{"leases", "list", "ReadOnly=false"}, noStdinString, leaseListLeasesString, noErrorString},
		CliTest{false, true, []string{"leases", "list", "ReadOnly=fred"}, noStdinString, noContentString, bootEnvBadReadOnlyString},

		CliTest{true, true, []string{"leases", "show"}, noStdinString, noContentString, leaseShowNoArgErrorString},
		CliTest{true, true, []string{"leases", "show", "john", "john2"}, noStdinString, noContentString, leaseShowTooManyArgErrorString},
		CliTest{false, true, []string{"leases", "show", "192.168.100.111"}, noStdinString, noContentString, leaseShowMissingArgErrorString},
		CliTest{false, false, []string{"leases", "show", "192.168.100.110"}, noStdinString, leaseShowLeaseString, noErrorString},
		CliTest{false, true, []string{"leases", "show", "k192.168.100.110"}, noStdinString, noContentString, leaseShowInvalidAddressErrorString},

		CliTest{true, true, []string{"leases", "exists"}, noStdinString, noContentString, leaseExistsNoArgErrorString},
		CliTest{true, true, []string{"leases", "exists", "john", "john2"}, noStdinString, noContentString, leaseExistsTooManyArgErrorString},
		CliTest{false, false, []string{"leases", "exists", "192.168.100.110"}, noStdinString, leaseExistsLeaseString, noErrorString},
		CliTest{false, true, []string{"leases", "exists", "192.168.100.111"}, noStdinString, noContentString, leaseExistsMissingJohnString},
		CliTest{true, true, []string{"leases", "exists", "john", "john2"}, noStdinString, noContentString, leaseExistsTooManyArgErrorString},

		CliTest{true, true, []string{"leases", "update"}, noStdinString, noContentString, leaseUpdateNoArgErrorString},
		CliTest{true, true, []string{"leases", "update", "john", "john2", "john3"}, noStdinString, noContentString, leaseUpdateTooManyArgErrorString},
		CliTest{false, true, []string{"leases", "update", "192.168.100.110", leaseUpdateBadJSONString}, noStdinString, noContentString, leaseUpdateBadJSONErrorString},
		CliTest{false, false, []string{"leases", "update", "192.168.100.110", leaseUpdateInputString}, noStdinString, leaseUpdateJohnString, noErrorString},
		CliTest{false, true, []string{"leases", "update", "192.168.100.111", leaseUpdateInputString}, noStdinString, noContentString, leaseUpdateJohnMissingErrorString},
		CliTest{false, false, []string{"leases", "show", "192.168.100.110"}, noStdinString, leaseUpdateJohnString, noErrorString},
		CliTest{false, true, []string{"leases", "update", "k192.168.100.111", leaseUpdateInputString}, noStdinString, noContentString, leaseUpdateInvalidAddressErrorString},

		CliTest{true, true, []string{"leases", "patch"}, noStdinString, noContentString, leasePatchNoArgErrorString},
		CliTest{true, true, []string{"leases", "patch", "john", "john2", "john3"}, noStdinString, noContentString, leasePatchTooManyArgErrorString},
		CliTest{false, true, []string{"leases", "patch", leasePatchBaseString, leasePatchBadPatchJSONString}, noStdinString, noContentString, leasePatchBadPatchJSONErrorString},
		CliTest{false, true, []string{"leases", "patch", leasePatchBadBaseJSONString, leasePatchInputString}, noStdinString, noContentString, leasePatchBadBaseJSONErrorString},
		CliTest{false, false, []string{"leases", "patch", leasePatchBaseString, leasePatchInputString}, noStdinString, leasePatchJohnString, noErrorString},
		CliTest{false, true, []string{"leases", "patch", leasePatchMissingBaseString, leasePatchInputString}, noStdinString, noContentString, leasePatchJohnMissingErrorString},
		CliTest{false, false, []string{"leases", "show", "192.168.100.110"}, noStdinString, leasePatchJohnString, noErrorString},

		CliTest{true, true, []string{"leases", "destroy"}, noStdinString, noContentString, leaseDestroyNoArgErrorString},
		CliTest{true, true, []string{"leases", "destroy", "john", "june"}, noStdinString, noContentString, leaseDestroyTooManyArgErrorString},
		CliTest{false, false, []string{"leases", "destroy", "192.168.100.110"}, noStdinString, leaseDestroyJohnString, noErrorString},
		CliTest{false, true, []string{"leases", "destroy", "192.168.100.110"}, noStdinString, noContentString, leaseDestroyMissingJohnString},
		CliTest{false, false, []string{"leases", "list"}, noStdinString, leaseDefaultListString, noErrorString},
		CliTest{false, true, []string{"leases", "destroy", "k192.168.100.110"}, noStdinString, noContentString, leaseDestroyInvalidAddressErrorString},

		CliTest{false, false, []string{"leases", "create", "-"}, leaseCreateInputString + "\n", leaseCreateJohnString, noErrorString},
		CliTest{false, false, []string{"leases", "list"}, noStdinString, leaseListLeasesString, noErrorString},
		CliTest{false, false, []string{"leases", "update", "192.168.100.110", "-"}, leaseUpdateInputString + "\n", leaseUpdateJohnString, noErrorString},
		CliTest{false, false, []string{"leases", "show", "192.168.100.110"}, noStdinString, leaseUpdateJohnString, noErrorString},

		CliTest{false, false, []string{"leases", "destroy", "192.168.100.110"}, noStdinString, leaseDestroyJohnString, noErrorString},
		CliTest{false, false, []string{"leases", "list"}, noStdinString, leaseDefaultListString, noErrorString},
		// Teardown subnet
		CliTest{false, false, []string{"subnets", "destroy", "john"}, noStdinString, subnetDestroyJohnString, noErrorString},
	}

	for _, test := range tests {
		testCli(t, test)
	}

}
