.. Copyright (c) 2017 RackN Inc.
.. Licensed under the Apache License, Version 2.0 (the "License");
.. Digital Rebar Provision documentation under Digital Rebar master license
.. index::
  pair: Digital Rebar Provision; DHCP Models

DHCP Management Models
<<<<<<<<<<<<<<<<<<<<<<

These models manage how the DHCP server built into dr-provision.  They
determines what IP addresses it can hand out to which systems, what
values to set for DHCP options, and how to handle DHCP requests at
various points in the lifecycle of any given DHCP lease.

Subnet
------

The Subnet Object defines the configuration of a single subnet for the
DHCP server to process.  Multiple subnets are allowed.  The Subnet can
represent a local subnet attached to a local interface (Broadcast
Subnet) to the Digital Rebar Provision server or a subnet that is
being forwarded or relayed (Relayed Subnet) to the Digital Rebar
Provision server.  Subnet objects have the following fields:

- Name: The unique name of this Subnet.

- Enabled: A boolean value that indicates whether this subnet is
  available to hand out new Leases, and whether to allow renewals of
  any Leases in its range.  Setting this to `true` allows the subnet
  to operate normally, and setting it to `false` will cause
  dr-provision to refuse to hand out new Leases for addresses in its
  range and will cause lease renewals for already existing Leases in
  its address range to fail.

- Proxy: A boolean value that indicates that dr-provision should
  respond to requests for addresses in this address range as if it was
  a proxy DHCP server (as defined in section 2 of `the PXE
  specification
  <http://www.pix.net/software/pxeboot/archive/pxespec.pdf>`_).

- Subnet: The network address in CIDR form of this Subnet.  Subnets
  may not have overlapping address ranges.

- ActiveStart: This is the start of the IP address range that this
  subnet will hand out.  It must be within the address range the
  Subnet is responsible for, and it must be less than ActiveEnd.

- ActiveEnd: This is the end of the IP address range that this subnet
  will hand out.  It must be within the address range the Subnet is
  responsible for, and it must be greater than ActiveStart.

- ActiveLeaseTime: This is the time (in seconds) that a lease created
  in this subnet will be valid for.

- ReservedLeaseTime: This is the time (in seconds) that a lease
  created by a Reservation in this subnet will be valid for.  It
  overrides ActiveLeaseTime.

- OnlyReservations: If set to `true`, then Leases in this subnet range can only be created if there is a reservation 
