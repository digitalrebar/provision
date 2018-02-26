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

Configuring Subnets is a critical first step in Digital Rebar Provision operation.  The basic UI will show all configured subnets and provide an easy way to add broadcast subnets based on the known interfaces.

To edit or delete a subnet, click on the name of the subnet to populate the editing area below the list.  To create a relay subnet, click on the add subnet link.  To create a broadcast subnet, click on the link provided after the name of the unassigned interfaces.

There are two primary types of subnets: broadcast and relay:

  * Broadcast subnets are associated with the addresses provided by the Digital Rebar Provision host system.  Digital Rebar Provision can have multiple DHCP ranges; however, only one subnet can be assigned per interface _and_ the subnet CIDR must include the IP of the interface (the range should NOT).  By convention, subnets associated with an interface will be named as the interface.
  * Relay subnets answer requests forwarded from DHCP relays such as those provided by switches.  These can be any suitable IP range.  Since the Relay subnets are not broadcast, they do not conflict with existing DHCP servers in the environment.

Digital Rebar Provision can operate in a permissive reservation mode or require users to whitelist systems before they are serviced.  The `OnlyReservations` flag will operate as a reservations required (whitelist) mode when true; otherwise, Digital Rebar Provision permissive reservation mode will give out addresses to any valid DHCP request.

In additional to serving IPs, DHCP servers provide critical configuration (aka `DHCP Options <https://en.wikipedia.org/wiki/Dynamic_Host_Configuration_Protocol#DHCP_options>`_) information to the clients.  Setting Option 67, Next Boot, is essential for Digital Rebar Provision to operate as a Provisioner.  This information includes next boot (67), gateway (3), domain name (15), DNS (6) and other important information.  It is encoded in the responses according to `IETF RFC 2132 <https://tools.ietf.org/html/rfc2132>`_

Consult the `Godocs <https://godoc.org/github.com/digitalrebar/provision/backend#Subnet>`_ for more details about the specific fields.

.. _rs_ui_bootenvs:

Boot Environments (bootenvs)
----------------------------

Configuring at least one Boot Environment is a critical first step in Digital Rebar Provision operation.  The Digital Rebar CentOS based in-memory discovery image, Sledgehammer, will be installed on first use by default.

The UI will show a complete list of potential Boot Environments;

.. _rs_ui_machines:

Machines
--------

Machines are central to the provisioning process, as they connect Boot Environments to incoming IP addresses. Desired BootEnvs can be assigned to machines from this page.

If a machine is not given a BootEnv, it will use the BootEnv listed as *defaultBootEnv* on the BootEnvs Preferences page.

Profiles can also be bound to machines from this view. Machines will relay the parameters of a profile to the templates provided by the selected BootEnv.

.. _rs_ui_profiles:

Profiles
--------

Profiles provide a convenient way to apply sets of parameters to a machine. Multiple profiles can be assigned to one machine, and will be referenced in the order they are listed.

Parameters can be linked to specific profiles through the profiles page, which can then be attached to machines through the machines UI.

.. _rs_ui_templates:

Templates
---------

Templates contain important instructions for the provisioning process, and are comprised of `golang text/template <https://golang.org/pkg/text/template/#hdr-Actions>`_ strings. Once templates are rendered along with any assigned parameters, they are used by the BootEnv to boot the target machine.

Templates may contain other templates, known as sub-templates.

.. _rs_ui_params:

Parameters (params)
-------------------

Parameters are passed to a template from a machine, and help to drive the template's functions. They consist of key/value pairs that provide configuration to the renderer.

Profiles allow params to be applied in bulk, or they can be attached to templates individually.

.. _rs_ui_reservations:

Reservations
------------

Reservations link tokens to specific IP addresses. This view shows a list of existing reservations along with tokens and strategies associated with each. Currently, MAC is the only available strategy.

Reservations may contain options to be applied to connected servers, which are also visible through the UI.

.. _rs_ui_leases:

Leases
------

Leases show individual links between tokens and addresses, created through reservation or subnet strategies. Leases remain valid for short periods of time, and cannot be edited.

The expiration time of each lease is visible through the UI. Once a lease has expired, it may be removed.

.. _rs_ui_tasks:

Tasks
-----

During the boot process, tasks provide additional configuration to machines in the form of templates. BootEnvs will use these sets of templates to construct specific jobs for a machine.

Within a task, templates are processed in the order they are assigned, so it's important to check that templates are attached correctly to a task.

.. _rs_ui_jobs:

Jobs
-----

A job defines a machine's current step in its boot process. After completing a job, the machine creates a new job from the next instruction in the machine's task list.

Machines will only process one job at a time, and jobs aren't created until the instant they are required.
