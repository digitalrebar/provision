.. Copyright (c) 2017 RackN Inc.
.. Licensed under the Apache License, Version 2.0 (the "License");
.. Digital Rebar Provision documentation under Digital Rebar master license
.. index::
  pair: Digital Rebar Provision; Developer Install

.. _rs_install_dev:

Developer Install (Console Isolated)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

This console install (aka ``--isolated``) guide provides a developer focused installation process that runs DRP from the Linux or MacOS command line instead of systemd.

This sub-install guide is a :ref:`rs_quickstart` supplement designed for code and operations development with Digital Rebar.  The command line install is especially useful for debugging and rapid testing where a minimial instalation footprint is required.

Please consult other installation paths for more details:

* :ref:`rs_quickstart` is a basic `systemd <https://en.wikipedia.org/wiki/Systemd>`_ install for new users
* :ref:`rs_install` details more complex installs including offline/airgap.
* :ref:`rs_install` details more complex installs including offline/airgap.
* :ref:`rs_install_docker` for trial users minimizing their install requirements
* :ref:`rs_install_cloud` is non-PXE / Cloud-Only installation process (no DHCP required)
* `Edge Lab with RPi <http://edgelab.digital>`_ is self-contained Digital Rebar inexpensive lab using Raspberry Pi computers.

Each of these environments requires careful setup up of your network environment and consideration with regard to competing DHCP services.  The setup of these environments is outside the scope of this document.

You must install Digital Rebar to use it, there is no SaaS version.  :ref:`rs_self_managed_why`

.. _rs_dev_preparation:

Preparation
-----------

Unlike the Quick Start :ref:`rs_qs_preparation`, this guide expects that users to run Digital Rebar in their primary work environment.  This makes it easy to monitor and reset the Digital Rebar endpoint.

This document refers to the :ref:`rs_cli` for manipulating the ``dr-provision`` service; however, the ``--isolated`` installation does *not* install ``drpcli`` automatically.

Please make sure your environment doesn't have any conflicts or issues that might cause port conflicts (see :ref:`rs_arch_ports`) or cause PXE booting to fail.

.. _rs_dev_install:

Install
-------

To begin, execute the following commands in a shell or terminal:

.. code-block:: bash

    mkdir drp ; cd drp
    curl -fsSL get.rebar.digital/stable | bash -s -- install --isolated --version=tip

The command will pull the *latest* ``dr-provision`` bundle and checksum from github, extract the files, verify prerequisites are installed, and create some initial directories and links.

The `install <http://get.rebar.digital/stable/>`_ script used by our installs has many additional options that are documented in its help and explored in other install guides.

Once you have Digital Rebar working in this fashion, you should be able to upgrade the endpoint by simply replacing the ``dr-provision`` binary and restarting the process.  The ability to quickly reset the environment is a primary benefit of this approach.

.. _rs_dev_start:

Manually Start Digital Rebar
----------------------------

Once the install has completed, your terminal should provide next steps similar to those below.

.. code-block:: bash

    # Run the following commands to start up dr-provision in a local isolated way.
    # The server will store information and serve files from the ./drp-data directory.

    sudo ./dr-provision --base-root=`pwd`/drp-data |& tee log.out


.. _rs_dev_license:

Install RackN License
---------------------

If you have obtained a RackN license file using the process from :ref:`rs_qs_license`, then you can bypass this step in subsequent resets by uploading the ``rackn-license.json`` file via the CLI.

.. code-block:: bash

    drpcli contents upload rackn-license.json


.. _rs_dev_next_steps:

Back to Regular Install
-----------------------

Once Digital Rebar is running in isolated mode, it is exactly the same as any other installation

* :ref:`rs_qs_license`
* :ref:`rs_qs_ux_bootstrap`
* :ref:`rs_qs_cli_bootstrap`
* :ref:`rs_qs_first_machine`
* :ref:`rs_qs_next_steps`

.. _rs_dev_cleanup:

Clean Up
--------

Once you are finished exploring Digital Rebar Provision in isolated mode, the system can cleaned or reset by removing the directory containing the isolated install.  In the previous sections, we used ''drp'' as the directory containing the isolated install.  Removing this directory will clean up the installed files.