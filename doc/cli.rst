.. Copyright (c) 2017 RackN Inc.
.. Licensed under the Apache License, Version 2.0 (the "License");
.. DigitalRebar Provision documentation under Digital Rebar master license
.. index::
  pair: DigitalRebar Provision; Command Line Interface

.. _rs_cli:

Command Line Interface (CLI)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~

Using the devtool build process for the server will generate
the files needed to build the cli.

It generates the following directories:

* client
* models

The cli uses those client files to access the server.  The editable 
cli code lives in:

* cli

.. _rs_client:

Building Client
---------------

* go build -o drpcli cmds/drpcli.go


Running Client
--------------

* ./drpcli

