.. Copyright (c) 2017 RackN Inc.
.. Licensed under the Apache License, Version 2.0 (the "License");
.. DigitalRebar Provision documentation under Digital Rebar master license
.. index::
  pair: DigitalRebar Provision; Architecture

.. _rs_architecture:


Architecture
~~~~~~~~~~~~

DigitalRebar Provision is intended to be a very simple service that can run with minimal overhead in nearly any environment.  For this reason, all the needed components are combined into the Golang binary server including the UI and Swagger UI assets.  The binary is can be run as a user process or easily configured as a operating system service.

The service is designed to work with multiple backend data stores.  For stand alone operation, data is stored on the file system.  For Digital Rebar integration, data can be maintained in Consul.

The CLI is provided as a second executable so that it can be used remotely.

By design, there are minimal integrations between core services.  This allows the service to reduce complexity.  Beyond serving IPs and files, the primary action of the service is template expansion for boot environments (bootenv).  The template expansion system allows subsitition properties to be set on a global or per machine basis.

The architecture can be described in terms of the server and its data model.

.. toctree::
  :glob:
  :numbered:

  arch/server
  arch/data

