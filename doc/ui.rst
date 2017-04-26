.. Copyright (c) 2017 RackN Inc.
.. Licensed under the Apache License, Version 2.0 (the "License");
.. Digital Rebar Provision documentation under Digital Rebar master license
.. index::
  pair: Digital Rebar Provision; Developer Environment

.. _rs_ui:

User Interface (UI)
~~~~~~~~~~~~~~~~~~~

The Digital Rebar Provision UI is intended to help users with basic operational needs.  Advanced users should review the API documentation or consider using Digital Rebar for integration and workflow capabilities.

.. _rs_ui_subnets:

Subnets
-------

Configuring Subnets is a critical first step in Digital Rebar Provision operation.  The basic UI will show all configurated subnets and provide an easy way to add broadcast subnets based on the known interfaces.

To edit or delete a subnet, click on the name of the subnet to populate the editing area below the list.  To create a relay subnet, click on the add subnet link.  To create a broadcast subnet, click on the link provided after the name of the unassigned interfaces.

There are two primary types of subnets: broadcast and relay:

  * Broadcast subnets are associated with the addresses provided by the Digital Rebar Provision host system.  Digital Rebar Provision can have multiple DHCP ranges; however, you can only assign one subnet per interface _and_ the subnet CIDR must include the IP of the interface (the range should NOT).  By convention, subnets associated with an interface will be named as the interface.
  * Relay subnets answer requests forwarded from DHCP relays such as those provided by switches.  These can be any suitable IP range.  Since the Relay subnets are not broadcast, they do not conflict with existing DHCP servers in the environment.

Digital Rebar Provision can operate in a permissive reservation mode or require users to whitelist systems before they are serviced.  The `OnlyReservations` flag will operate as a reservations required (whitelist) mode when true; otherwise, Digital Rebar Provision permissive reservation mode will give out addresses to any valid DHCP request.

In additionl to serving IPs, DHCP servers provide critical confinguration (aka `DHCP Options <https://en.wikipedia.org/wiki/Dynamic_Host_Configuration_Protocol#DHCP_options>`_) information to the clients.  Setting Option 67, Next Boot, is essential for Digital Rebar Provision to operate as a Provisioner.  This information includes next boot (67), gateway (3), domain name (15), DNS (6) and other important information.  It is encoded in the responses according to `IETF RFC 2132 <https://tools.ietf.org/html/rfc2132>`_

Consult the `Godocs <https://godoc.org/github.com/digitalrebar/provision/backend#Subnet>`_ for more details about the specific fields.

.. _rs_ui_bootenvs:

Boot Environments (bootenvs)
----------------------------

Configuring at least one Boot Environment is a critical first step in Digital Rebar Provision operation.  The Digital Rebar CentOS based in-memory discovery image, Sledgehammer, will be installed on first use by default.

The UI will show a complete list of potential Boot Environments;

.. _rs_swagger:

Swagger UI
~~~~~~~~~~

The Digital Rebar Provision UI includes Swagger to allow you to explore and test the API.
