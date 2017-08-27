.. Copyright (c) 2017 RackN Inc.
.. Licensed under the Apache License, Version 2.0 (the "License");
.. Digital Rebar Provision documentation under Digital Rebar master license
.. index::
  pair: Digital Rebar Provision; API
  pair: Digital Rebar Provision; REST

.. _rs_api:

Digital Rebar Provision API
~~~~~~~~~~~~~~~~~~~~~~~~~~~

In general, the Digital Rebar Provision API is documented via the Swagger spec and introspectable for machines via `/swagger.json` and for humans via `/swagger-ui`.

All API calls are available under `/api/v3` based on the Digital Rebar API convention.

.. _rs_api_notes:

API Exception Notes
-------------------

There are times when the API and models have exceptions or changes that do not follow the normal pattern.  This section is designed to provide a reference for those exceptions

This section is intended to provide general information about and functional of the API

*Machines.Profile*
  What would otherwise be Machine.Params is actually embedded under Machines.Profile.Params.
  This composition simplifies that precedence calculation for Params by making Machines the
  top of the Profiles stack.  All the other fields in the Machines.Profile are ignored.


.. swaggerv2doc:: https://github.com/digitalrebar/provision/releases/download/tip/swagger.json

