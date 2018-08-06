package cli

import (
	"testing"
)

func TestReservationCli(t *testing.T) {
	var reservationCreateBadJSONString = "asdgasdg"
	var reservationCreateInputString string = `{
  "Addr": "192.168.100.100",
  "NextServer": "2.2.2.2",
  "Strategy": "MAC",
  "Token": "john"
}
`
	var reservationUpdateBadJSONString = "asdgasdg"
	var reservationUpdateInputString string = `{
  "Options": [ { "Code": 3, "Value": "1.1.1.1" } ]
}
`
	cliTest(true, false, "reservations").run(t)
	cliTest(false, false, "reservations", "list").run(t)
	cliTest(false, true, "reservations", "create", "-").Stdin(`---
Addr: "192.168.100.101"
Strategy: MAC
Token: frank
Scoped: true`).run(t)
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
	verifyClean(t)
}

func TestScopedReservations(t *testing.T) {
	cliTest(false, false, "subnets", "create", "-").Stdin(`---
Name: aa
Strategy: MAC
Subnet: "192.168.124.0/24"
ActiveStart: "192.168.124.10"
ActiveEnd: "192.168.124.20"`).run(t)
	cliTest(false, false, "subnets", "create", "-").Stdin(`---
Name: bb
Strategy: MAC
Subnet: "192.168.125.0/24"
ActiveStart: "192.168.125.10"
ActiveEnd: "192.168.125.20"`).run(t)
	cliTest(false, true, "reservations", "create", "-").Stdin(`---
Addr: "192.168.123.30"
Strategy: MAC
Scoped: true
Token: baz`).run(t)
	cliTest(false, false, "reservations", "create", "-").Stdin(`---
Addr: "192.168.123.30"
Strategy: MAC
Token: baz`).run(t)
	cliTest(false, true, "reservations", "create", "-").Stdin(`---
Addr: "192.168.124.0"
Strategy: MAC
Scoped: true
Token: foo`).run(t)
	cliTest(false, true, "reservations", "create", "-").Stdin(`---
Addr: "192.168.124.255"
Strategy: MAC
Scoped: true
Token: foo`).run(t)
	cliTest(false, false, "reservations", "create", "-").Stdin(`---
Addr: "192.168.124.1"
Strategy: MAC
Scoped: true
Token: foo`).run(t)
	cliTest(false, true, "reservations", "create", "-").Stdin(`---
Addr: "192.168.124.2"
Strategy: MAC
Scoped: true
Token: foo`).run(t)
	cliTest(false, true, "reservations", "create", "-").Stdin(`---
Addr: "192.168.124.2"
Strategy: MAC
Scoped: true
Token: bar`).run(t)
	cliTest(false, true, "reservations", "create", "-").Stdin(`---
Addr: "192.168.124.3"
Strategy: MAC
Scoped: true
Token: baz`).run(t)
	cliTest(false, true, "reservations", "create", "-").Stdin(`---
Addr: "192.168.124.2"
Strategy: MAC
Scoped: true
Token: bar`).run(t)
	cliTest(false, true, "reservations", "create", "-").Stdin(`---
Addr: "192.168.125.0"
Strategy: MAC
Scoped: true
Token: foo`).run(t)
	cliTest(false, true, "reservations", "create", "-").Stdin(`---
Addr: "192.168.125.255"
Strategy: MAC
Scoped: true
Token: foo`).run(t)
	cliTest(false, false, "reservations", "create", "-").Stdin(`---
Addr: "192.168.125.1"
Strategy: MAC
Scoped: true
Token: foo`).run(t)
	cliTest(false, true, "reservations", "create", "-").Stdin(`---
Addr: "192.168.125.2"
Strategy: MAC
Scoped: true
Token: foo`).run(t)
	cliTest(false, false, "reservations", "create", "-").Stdin(`---
Addr: "192.168.125.2"
Strategy: MAC
Scoped: true
Token: bar`).run(t)
	cliTest(false, true, "reservations", "create", "-").Stdin(`---
Addr: "192.168.125.3"
Strategy: MAC
Scoped: true
Token: baz`).run(t)
	cliTest(false, true, "reservations", "create", "-").Stdin(`---
Addr: "192.168.125.2"
Strategy: MAC
Scoped: true
Token: bar`).run(t)
	cliTest(false, false, "reservations", "destroy", "192.168.125.2").run(t)
	cliTest(false, false, "reservations", "destroy", "192.168.125.1").run(t)
	cliTest(false, false, "reservations", "destroy", "192.168.124.2").run(t)
	cliTest(false, false, "reservations", "destroy", "192.168.124.1").run(t)
	cliTest(false, false, "reservations", "destroy", "192.168.123.30").run(t)
	cliTest(false, false, "subnets", "destroy", "bb").run(t)
	cliTest(false, false, "subnets", "destroy", "aa").run(t)
	verifyClean(t)
}
