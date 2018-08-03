package cli

import (
	"testing"
)

func TestSubnetCli(t *testing.T) {

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

	var subnetUpdateBadJSONString = "asdgasdg"

	var subnetUpdateInputString string = `{
  "Strategy": "NewStrat"
}
`
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
	cliTest(false, true, "subnets", "update", "john2", subnetUpdateInputString).run(t)
	cliTest(false, false, "subnets", "show", "john").run(t)
	cliTest(true, true, "subnets", "destroy").run(t)
	cliTest(true, true, "subnets", "destroy", "john", "june").run(t)
	cliTest(false, false, "subnets", "destroy", "john").run(t)
	cliTest(false, true, "subnets", "destroy", "john").run(t)
	cliTest(false, false, "subnets", "list").run(t)
	cliTest(false, false, "subnets", "create", "-").Stdin(subnetCreateInputString + "\n").run(t)
	cliTest(false, false, "subnets", "list").run(t)
	cliTest(false, true, "subnets", "update", "john", "-").Stdin(subnetUpdateInputString + "\n").run(t)
	cliTest(false, false, "subnets", "show", "john").run(t)
	cliTest(true, true, "subnets", "range").run(t)
	cliTest(true, true, "subnets", "range", "john", "1.24.36.7", "1.24.36.16", "1.24.36.16").run(t)
	cliTest(false, true, "subnets", "range", "john", "192.168.100.10", "192.168.100.500").run(t)
	cliTest(false, true, "subnets", "range", "john", "cq.98.42.1234", "1.24.36.16").run(t)
	cliTest(false, false, "subnets", "range", "john", "192.168.100.10", "192.168.100.200").run(t)
	cliTest(true, true, "subnets", "subnet").run(t)
	cliTest(true, true, "subnets", "subnet", "john", "june", "1.24.36.16").run(t)
	cliTest(false, true, "subnets", "subnet", "john", "192.168.100.0/10").run(t)
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
	cliTest(false, false, "subnets", "set", "john", "option", "6", "to", "67").run(t)
	cliTest(false, false, "subnets", "get", "john", "option", "6").run(t)
	cliTest(false, false, "subnets", "set", "john", "option", "6", "to", "null").run(t)
	cliTest(false, true, "subnets", "get", "john", "option", "6").run(t)
	//End of Helpers
	cliTest(false, false, "reservations", "create", "-").Stdin(`---
Addr: "192.168.100.100"
Strategy: MAC
Token: foo
Scoped: true`).run(t)
	cliTest(false, true, "subnets", "destroy", "john").run(t)
	cliTest(false, false, "reservations", "destroy", "192.168.100.100").run(t)
	cliTest(false, false, "subnets", "destroy", "john").run(t)
	cliTest(false, false, "subnets", "list").run(t)
	verifyClean(t)
}
