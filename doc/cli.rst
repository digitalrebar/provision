.. Copyright (c) 2017 RackN Inc.
.. Licensed under the Apache License, Version 2.0 (the "License");
.. Rocket Skates documentation under Digital Rebar master license
.. index::
  pair: Rocket Skates; Command Line Interface

.. _rs_cli:

Command Line Interface (CLI)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~

> DON'T DO THIS FOR NOW.

To generate the code:
* swagger generate client -P models.Principal -f swagger.json 

This generates the following directories:

* client
* cmd
* models
* restapi

The file that we edit for the client is:

* cmd/rocket-skates-client/main.go 

This is completely our own file.

.. _rs_client:

Building Client
---------------

* go build -o rscli cmd/rocket-skates-client/main.go


Running Client
--------------

* ./rscli

This is a really simple client but shows how to make the calls and get back structures.


