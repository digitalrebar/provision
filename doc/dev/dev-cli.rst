.. Copyright (c) 2017 RackN Inc.
.. Licensed under the Apache License, Version 2.0 (the "License");
.. Digital Rebar Provision documentation under Digital Rebar master license
.. index::
  pair: Digital Rebar Provision; Developing the Command Line Interface

.. _rs_dev_cli:

Developing the Command Line Interface (CLI)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

Using the **build.sh** process for the server will generate
the files needed to build the cli.

It generates the following directories:

* client
* models

The cli uses those client files to access the server.  The editable
cli code lives in:

* cli

The hope is that the CLI will use a generated client library based upon
the generated swagger.json file.  This will help ensure that we are building
a valid and viable swagger.json file.  The build.sh tool generates all the
components need for the cli and also builds multiple instances of it.

.. _rs_client:

Building Client
---------------

While a single *go build* command will generate the cli, it is safer to
use the *build.sh* script to ensure that all the parts are accurately generated.

* tools/build.sh

The results are stored in the bin directory based upon OS and platform.  We
currently build windows, linux, and darwin for amd64.

Running Client
--------------

After building the code, use the tools/install.sh script to get a cli
and dr-provision binary in the top-level directory for the platform.

* tools/install.sh --isolated install

Once that has been done a single time, symbolic links are created so that running
commands from the top-level directory should work.

* ./drpcli

For more information, see :ref:`rs_install`.


