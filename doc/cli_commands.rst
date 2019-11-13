.. Copyright (c) 2017 RackN Inc.
.. Licensed under the Apache License, Version 2.0 (the "License");
.. Digital Rebar Provision documentation under Digital Rebar master license
.. index::
  pair: Digital Rebar Provision; Command Line Interface (CLI)
  pair: Digital Rebar Provision; drpclicommand

.. _rs_cli_command:

CLI Commands Reference
----------------------

.. toctree::
   :glob:
   :numbered:
   :maxdepth: 1

   cli/*

.. _rs_cli_faq:

CLI Frequently Asked Questions (FAQ)
------------------------------------

This is a specialized FAQ section for FAQ concerns.
Please see the general :ref: `rs_faq` page for additional items.

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
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

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
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

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
