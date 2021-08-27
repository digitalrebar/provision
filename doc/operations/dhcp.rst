.. Copyright (c) 2021 RackN Inc.
.. Licensed under the Apache License, Version 2.0 (the "License");
.. Digital Rebar Platform documentation under Digital Rebar master license
.. index::
  pair: Digital Rebar Platform; DHCP Operations

.. _rs_dhcp_ops:

DHCP Operations
===============

The DHCP server provided by dr-provision provides the basis for many of
our more advanced bare-metal automation capabilities, including dynamic
network boot control and UEFI secure boot support.

Startup Configuration
---------------------

Whether the DHCP server is enabled and what ports it will listen on are controlled
by the following startup options:

1. `--disable-dhcp` to disable basic DHCP functionality.
2. `--disable-pxe` to disable extended PXE/BINL support (mainly used by Windows)
3. `--dhcp-port` to change the DHCP port from the default of 67 to something else.
4. `--binl-port` to change the BINL port from the default of 4011 to something else.

If none of these options are present, then we default to DHCP and BINL being enabled
on their standard ports.  The DHCP server listens on all interfaces, and that
behaviour cannot be modified at the present time.

Runtime Configuration
---------------------

By default (without any Subnets or Reservations defined) the DHCP server will
ignore any traffic it receives.  To start responding to requests, Subnets and/or
Reservations must be defined.

Subnets
~~~~~~~

A Subnet ( :ref:`rs_dhcp_subnet` ) defines an address range that dr-provision attempt to handle DHCP traffic
from and the IP addressing information that should be applied to any interfaces requesting
configuration from within that range.

Whether any given packet maps into a Subnet depends on one of the following conditions being true:

- The DHCP packet was not relayed to us, and there is an address on the interface it was
  received on that is within range of that Subnet.
- The DHCP packet was relayed to us, and the GIADDR field in the packet is within the range of
  that Subnet.
- The DHCP packet was relayed to us, and Option 82 Suboption 5 is set and within the range of the Subnet.

Additionally, if Option 82 Suboption 11 is set, we will format ServerID in the response packet as appropriate.

Reservations
~~~~~~~~~~~~

A Reservation ( :ref:`rs_dhcp_reservation` ) defines a one-to-one mapping between a specific network interface
to be configured and the IP addressing information that should be applied to that interface.

A Reservation can define an IP address that is outside the bounds of any Subnet.

Runtime Behaviour
-----------------

- If there is no Subnet or Lease that covers an incoming DHCP packet, it will be ignored.
- DHCP Options from Subnets and Reservations stack.  If an incoming request could be handled by both
  a Reservation and a Subnet, then the DHCP options on the Reservation will take precedence over the
  ones from the Subnet.
- If you set DHCP logging to Debug or Trace, the dr-provision log will contain all inbound and
  outbound DHCP traffic, along with messages about how any given packet was handled.
- We rely on CHADDR for Lease and Subnet lookups.  DUIDs are ignored, as they are often different between
  the firmware, IPXE/PXELINUX, Sledgehammer, and the final installed OS.
- If the CHADDR in a DHCP packet corresponds to a recorded HardwareAddr for a Machine, the DHCP server will
  react intelligently to network boot requests.

