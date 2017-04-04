.. Copyright (c) 2017 RackN Inc.
.. Licensed under the Apache License, Version 2.0 (the "License");
.. Rocket Skates documentation under Digital Rebar master license
.. index::
  pair: Rocket Skates; Architecture

.. _rs_architecture:


Architecture
~~~~~~~~~~~~

Rocket Skates is intended to be a very simple service that can run with minimal overhead in nearly any environment.  For this reason, all the needed components are combined into the Golang binary server including the UI and Swagger UI assets.  The binary is can be run as a user process or easily configured as a operating system service.

The service is designed to work with multiple backend data stores.  For stand alone operation, data is stored on the file system.  For Digital Rebar integration, data can be maintained in Consul.

The CLI is provided as a second executable so that it can be used remotely.

By design, there are minimal integrations between core services.  This allows the service to reduce complexity.  Beyond serving IPs and files, the primary action of the service is template expansion for boot environments (bootenv).  The template expansion system allows subsitition properties to be set on a global or per machine basis.

.. _rs_design_restriction:

Design Restrictions
-------------------

Since Rocket Skates is part of the larger Digital Rebar system, it's scope is limited to handling DHCP and Provisioning actions.  Out of band management to control server flow or configure firmware plus other management features will be handled by other Digital Rebar services.

Services
--------

Provisioning requires handoffs between multiple services as described in the :ref:`rs_workflows` section.  Since several of services are standard protocols (DHCP, TFTP, HTTP), it may be difficult to change ports without breaking workflow.

The figure below illustrates the three core Rocket Skates services including protocols and default ports.  The services are:

#. Web - These services provide control for the other services

   #. API: REST endpoints with Swagger definition
   #. UI: User interface and Swagger helpers

#. DHCP: Address management includes numerous additional option fields used to tell systems how to interact with other data center services such as provisioning, DNS, NTP and routing.

#. Provision: sends files on request during provisioning process based on a template system:

   #. TFTP: very simple (but slow) protocol that's used by firmware boot processes because it is very low overhead.
   #. HTTP: faster file transfer protocol used by more advanced boot processes


.. figure::  images/core_services.png
   :alt: Core Rocket Skates Services
   :target: https://docs.google.com/drawings/d/1SVGGwQZxopiVEYjIM3FXC92yG4DKCCejRBDNMsHmxKE/edit?usp=sharing

