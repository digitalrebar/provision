.. Copyright (c) 2017 RackN Inc.
.. Licensed under the Apache License, Version 2.0 (the "License");
.. DigitalRebar Provision documentation under Digital Rebar master license
.. index::
  pair: DigitalRebar Provision; Server Architecture

.. _rs_server_architecture:

Server Architecture
===================

DigitalRebar Provision is provided by a single binary that contains tools and images needed to operate.
These are expanded on startup and made available by the file server services.

.. _rs_design_restriction:

Design Restrictions
-------------------

Since DigitalRebar Provision is part of the larger Digital Rebar system, it's scope is limited to handling DHCP and Provisioning actions.  Out of band management to control server flow or configure firmware plus other management features will be handled by other Digital Rebar services.

.. _rs_arch_services:

Services
--------

Provisioning requires handoffs between multiple services as described in the :ref:`rs_workflows` section.  Since several of services are standard protocols (DHCP, TFTP, HTTP), it may be difficult to change ports without breaking workflow.

The figure below illustrates the three core DigitalRebar Provision services including protocols and default ports.  The services are:

#. Web - These services provide control for the other services

   #. API: REST endpoints with Swagger definition
   #. UI: User interface and Swagger helpers

#. DHCP: Address management includes numerous additional option fields used to tell systems how to interact with other data center services such as provisioning, DNS, NTP and routing.

#. Provision: sends files on request during provisioning process based on a template system:

   #. TFTP: very simple (but slow) protocol that's used by firmware boot processes because it is very low overhead.
   #. HTTP: faster file transfer protocol used by more advanced boot processes


.. figure::  ../images/core_services.png
   :alt: Core DigitalRebar Provision Services
   :target: https://docs.google.com/drawings/d/1SVGGwQZxopiVEYjIM3FXC92yG4DKCCejRBDNMsHmxKE/edit?usp=sharing


.. _rs_arch_ports:

Ports
-----

The table describes the ports that need to be available to run DigitalRebar Provision.  Firewall rules may need to be altered to enable these services.  The feature column indicates when the port is required.  For example, the DHCP server can be turned off and that port is no longer required.

========  =======   =====================
Ports     Feature   Usage
========  =======   =====================
67/udp    DHCP      DHCP Port
69/udp    PROV      TFTP Port
8091/tcp  PROV      HTTP-base File Server
8092/tcp  Always    DR Provision Mgmt
========  =======   =====================


