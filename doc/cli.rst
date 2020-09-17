.. Copyright (c) 2017 RackN Inc.
.. Licensed under the Apache License, Version 2.0 (the "License");
.. Digital Rebar Provision documentation under Digital Rebar master license
.. index::
  pair: Digital Rebar Provision; Command Line Interface (CLI)
  pair: Digital Rebar Provision; drpcli

.. _rs_cli:

Digital Rebar Provision Command Line Interface (CLI)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

The Digital Rebar Provision Command Line Interface (drpcli) provides a simplified way to interact with the :ref:`rs_api`.  The command line tool (``drpcli``) is auto-generated from source code via reference of the API.  This means the CLI should implement 100% coverage of the API.

.. _rs_cli_download:

Download DRPCLI
===============

If you've already installed the Digital Rebar server, then the CLI is available automatically from the server's files list.  You should download the CLI directly from the server using `https:\\[drpserveraddress]:8092\files\drpcli.amd64.darwin` or similar depending on your OS and architecture.

.. note:: There is a ``</>`` button on the UX top right corner that will download the right binary from your endpoint.

To install without the Digital Rebar server, you should review the catalog at `https://repo.rackn.io` for the desired version (stable is safest) to use.  The following script can be used to quickly download DRPCLI and then use the catalog function to ugprade to the latest version.

  ::

    #!/usr/bin/env bash
    # Get 'drpcli' binary.
    # Copyright RackN, Inc - 2020

    # simple script to get the 'drpcli' client binary - to change OS Architecture
    # or Platform, set environment variables ARCH and PLATFORM.  Defaults
    # to "amd64" and "linux".
    #
    # Change the base repo location if you have a different Catalog source by
    # setting the REPO environment variable.  Defaults to https://repo.rackn.io
    #
    # Change drpcli version by using one of the version names tagged in the
    # catalog by setting VERSION environment variable.  Defaults to "stable"
    #
    # Some examples
    #
    # ARCH="arm64" PLATFORM="linux" VERSION="tip" ./get-drpcli.sh
    # ARCH="amd64" PLATFORM="darwin" VERSION="v4.4.0" ./get-drpcli.sh
    # ARCH="amd64" PLATFORM="windows" ./get-drpcli.sh


    REPO=${REPO:-"https://repo.rackn.io"}
    ARCH=${ARCH:-"amd64"}
    PLAT=${PLATFORM:-"linux"}
    VER=${VERSION:-"stable"}

    BASE=$(curl -s --compressed $REPO | jq -r '.sections.catalog_items[]| select(.Name == "drpcli") | select(.Version == "'$VER'") | .Source')

    [[ "$PLAT" == "windows" ]] && EXT=".exe" || EXT=""
    [[ -f "drpcli${EXT}" ]] && { echo "drpcli${EXT} exists, not overwriting, exiting"; exit 1; } || true
    curl --progress-bar -kfSL $BASE/$ARCH/$PLAT/drpcli${EXT} -o drpcli${EXT}
    (( $? )) && echo "download failure" || echo "Saved as:  drpcli${EXT}"
    [[ -f "drpcli${EXT}" ]] && chmod 755 drpcli${EXT} || false


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

.. note:: VERY IMPORTANT - the **update** commands use the **PATCH** operation for the objects in the :ref:`rs_api`.  This has the implication that for map like components (Params sections of :ref:`rs_model_machine` and :ref:`rs_model_profile`) the contents are merged with the existing object.  For the Params sections specifically, use the subaction *params* to replace contents.

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

.. note:: It is necessary to specify either a username and password or a token.

Another useful flag is *--format*.  this will change the tool output to YAML instead of JSON.  This can
be helpful when editing files by hand.  e.g. *--format yaml*

For Bash users, the drpcli can generate its own bash completion file.  Once generated, it is necessary to restart
the terminal/shell or reload the completions.

.. admonition:: linux

  ::

    sudo drpcli autocomplete /etc/bash_completion.d/drpcli
    . /etc/bash_completion

.. admonition:: Darwin

  Assuming that Brew is in use to update and manage bash and bash autocompletion.

  ::

    sudo drpcli autocomplete /usr/local/etc/bash_completion.d/drpcli
    . /usr/local/etc/bash_completion


.. _rs_cli_command:

CLI Commands Reference
======================

.. toctree::
   :glob:
   :maxdepth: 1

   cli/*

.. _rs_cli_filters:

Filtering Results for DRPCLI
============================

There are several ways to filter DRPCLI output;

.. _rs_cli_filters_api:

DRPCLI and API Filters
----------------------

The DRPCLI passes through API filters.  See :ref:`rs_api_filters` for more details.

DRPCLI allows multiple filters to be passed in a single command.  They are considered to be AND operations.

When passing filters via DRPCLI, you may need to protect specialized filters using quotes or
single ticks.

For example, to use the not equal (Ne) filter:

  ::

     drpcli machines list Name='Ne(cluster01)'
     drpcli machines list Name="Ne(cluster01)"

.. _rs_cli_filters_jq:

DRPCLI and JSON Query (JQ)
--------------------------

DRPCLI includes a fully jq parser!  Please see :ref:`rs_kb_00042` for details about using JQ with DRPCLI.

.. note:: Since JQ acts _after_ results have been returned, API filters should be applied first when possible.  Returning large data sets will degrade system performance.

.. _rs_cli_faq:

CLI Frequently Asked Questions (FAQ)
====================================

This is a specialized FAQ section for FAQ concerns.
Please see the general :ref:`rs_faq` page for additional items.

.. _rs_download_rackn_content:

Download Content and Plugins via Command Line
---------------------------------------------

RackN maintains a catalog of open and proprietary Digital Rebar extensions at ``https://repo.rackn.io``.

Content downloads directly from the Catalog as JSON and can be imported directly using the DRP CLI.
  ::
    drpcli contents upload catalog:task-library-tip

or

  ::

    drpcli catalog item install task-library --version tip

Plugin downloads directly from the Catalog work as follows:

  ::

    drpcli plugin_providers upload raid from catalog:raid-stable

or

  ::

    drpcli catalog item install raid

.. _rs_autocomplete:

Turn on autocomplete for the CLI
--------------------------------

The DRP CLI has built in support to generate autocomplete (tab completion) capabilities for the BASH shell.  To enable, you must generate the autocomplete script file, and add it to your system.  This can also be added to your global shell ``rc`` files to enable autocompletion every time you log in.  NOTE that most Linux distros do this slightly differently.  Select the method that works for your distro.

You must specify a filename as an argument to the DRP CLI autocomplete command.  The filename will be created with the autocomplete script.  If you are writing to system areas, you need ``root`` access (eg via `sudo`).

For Debian/Ubuntu and RHEL/CentOS distros:
  ::

    sudo drpcli autocomplete /etc/bash_completion.d/drpcli

For Mac OSX (Darwin):
  ::

    sudo drpcli autocomplete /usr/local/etc/bash_completion.d/drpcli

Once the autocomplete file has been created, either log out and log back in, or ``source`` the created file to enable autocomplete in the current shell session (example for Linux distros, adjust accordingly):
  ::

    source /etc/bash_completion.d/drpcli

.. note:: If you receive an error message when using autocomplete similar to:
    ::

      bash: _get_comp_words_by_ref: command not found

  Then you will need to install the ``bash-completion`` package (eg. ``sudo yum -y install bash-completion`` or ``sudo apt -y install bash-completion``).

  You will also need to log out and then back in to your shell account to correct the bash_completion issue.

.. _rs_cli_faq_zip:

How do I upload multiple files using a zip/tar file?
----------------------------------------------------

The DRP files API allows exploding a compressed file, using
bsdtar, after it has been uploaded.  This can be very
helpful when multiple files or a full directory tree
are being uploaded.

This is a two stage process enabled by the `--explode`
flag.  The first stage simply uploads the compressed
file to the target path and location.  The second stage
explodes the file in that path.

For example: `drpcli files upload my.zip as mypath/my.zip --explode`

The above command will upload the `my.zip` file to the
files `/mypath` location.  It will also expand all
the files in `my.zip` into `/mypath` after upload.
All paths in `my.zip` will be preserved and created
relative to `/mypath/`.
