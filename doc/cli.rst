.. Copyright (c) 2017 RackN Inc.
.. Licensed under the Apache License, Version 2.0 (the "License");
.. Digital Rebar Provision documentation under Digital Rebar master license
.. index::
  pair: Digital Rebar Provision; Command Line Interface (CLI)
  pair: Digital Rebar Provision; drpcli

.. _rs_cli:

Digital Rebar Provision Command Line Interface (CLI)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

The Digital Rebar Provision Command Line Interface (drpcli) prevents a simplified way to interact with the
:ref:`rs_api`.

Overview
========

The CLI provides help for commands and follows a pattern of chained parameters with a few flags for additional 
modifications.

Some examples are:

  ::

    drpcli bootenvs list
    drpcli subnets get mysubnet
    drpcli preferences set defaultBootEnv discovery


The *drpcli* has help at each layer of command and is the easiest way to figure out what can and can not be done.

  ::

    drpcli help
    drpcli bootenvs help


Each object in the :ref:`rs_data_architecture` has a CLI subcommand.

By default, the CLI will attempt to access the *dr-provision* API endpoint on the localhost at port 8092 with
the username and password of *rocketskates* and *r0cketsk8ts*, respectively.
All three of these values can be provided by environment variable or command line flag.

======== ==================== ================ ==============================================================
Option   Environment Variable Flag             Format
======== ==================== ================ ==============================================================
Username RS_KEY               -P or --password String, but when part of RS_KEY it is: username:password
Password RS_KEY               -U or --username String, but when part of RS_KEY it is: username:password
Token    RS_TOKEN             N/A              Base64 encoded string from a generate token API call.
Endpoint RS_ENDPOINT          -E or --endpoint URL for access, https://IP:PORT. e.g. https://127.0.0.1:8092
======== ==================== ================ ==============================================================

.. note:: You must specify either a username and password or a token.

Another useful flag is *--format*.  You can use this to have the tool output YAML instead of JSON.  This can
be helpful when editting files by hand.  e.g. *--format yaml*

For Bash users, the drpcli can generate its own bash completion file.  Once generated, you will need to restart 
your terminal/shell or reload the completions.

.. admonition:: linux

  ::

    sudo drpcli autocomplete /etc/bash_completion.d/drpcli
    . /etc/bash_completion

.. admonition:: Darwin

  Assuming you are using Brew to update and manage bash and bash autocompletion.

  ::

    sudo drpcli autocomplete /usr/local/etc/bash_completion.d/drpcli
    . /usr/local/etc/bash_completion


Commands
========

.. toctree::
   :glob:
   :numbered:
   
   cli/*

