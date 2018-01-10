.. Copyright (c) 2017 RackN Inc.
.. Licensed under the Apache License, Version 2.0 (the "License");
.. Digital Rebar Provision documentation under Digital Rebar master license
.. index::
  pair: Digital Rebar Provision; DHCP Models

DHCP Management
<<<<<<<<<<<<<<<

The DHCP server built in to dr-provision is designed to be fully API
driven and to provide all the features needed to manage system IP
address assignments through the complete provisioning lifecycle. As
such, it has a few interesting features that other DHCP servers may
not have:

- The ability to have different ways of determining what unique
  attribute in a DHCP packet to use to allocate an IP address.  When
  you see references to Strategy and Token in the DHCP models,
  Strategy refers to the unique attribute the DHCP server should use,
  and Token refers to the value that the Stategy picked.

  For now, the only implemented Strategy is MAC, which has the DHCP
  server use the MAC address of the network adaptor of the network
  interface as the unique value of the Token.

- The DHCP server is fully API driven.  You can add, remove, and
  modify Reservations and Subnets on the fly, and changes take effect
  immediately.

- Built-in ProxyDHCP support, on a subnet by subnet basis.  The
  dr-provision can coexist with other DHCP servers to only provide PXE
  support for specific address ranges, leaving address management to
  your preexisting DHCP infrastructure.

Models
^^^^^^

These models manage how the DHCP server built into dr-provision.  They
determines what IP addresses it can hand out to which systems, what
values to set for DHCP options, and how to handle DHCP requests at
various points in the lifecycle of any given DHCP lease.

DHCP Option
-----------

The DHCP Option object holds templated values for DHCP options that
should be returned to clients in response to requests.  It has the following fields:

- Code: A byte that holds the numeric DHCP option code. See `RFC 2132
  <https://tools.ietf.org/html/rfc2132>`_ and friends for what these
  codes can be.

- Value: A string that will be template-expanded to form a valid value
  to return as the DHCP option.  Template expansion happens in the
  context of the source options.

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

- Strategy: A string that determines how the subnet will uniquely
  identify part of the DHCP request for address assignment.

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

- OnlyReservations: If set to `true`, then Leases in this subnet range
  can only be created if there is a reservation created for the
  requested address.

- Options: A list of DhcpOption objects that should be returned in any
  replies to dhcp requests.

Reservation
-----------

Reservations are what the dr-provision DHCP service uses to ensure
that an IP address is always issued to the same device.  Reservations
have the following fields:

- Strategy: The strategy that the DHCP service should use to determine
  whether this reservations should be used.

- Token: The unique string that the Strategy uniquely identifies a
  network interface with.

- Address: The IP address that is being reserved.

- Options: The DHCP options that should be returned when creating or
  renewing a Lease based on this Reservation.
