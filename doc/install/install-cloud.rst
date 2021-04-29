.. Copyright (c) 2017 RackN Inc.
.. Licensed under the Apache License, Version 2.0 (the "License");
.. Digital Rebar Provision documentation under Digital Rebar master license
.. index::
  pair: Digital Rebar Provision; Install Cloud

.. _rs_install_cloud:

Non-PXE Install (Cloud & VM)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~

This Non-PXE install (aka Cloud) guide provides a cloud focused installation process that runs DRP inside a Cloud Provider where PXE installation is not possible.

Other installation paths:

* :ref:`rs_quickstart` is a basic `systemd <https://en.wikipedia.org/wiki/Systemd>`_ install for new users
* :ref:`rs_install` details more complex installs including offline/airgap.
* :ref:`rs_install_dev` for developers running DRP interactively
* :ref:`rs_install_docker` for trial users minimizing their install requirements
* `Edge Lab with RPi <http://edgelab.digital>`_ is self-contained Digital Rebar inexpensive lab using Raspberry Pi computers.

Unlike other environments which requires careful setup up of your network environment and consideration with regard to competing DHCP services, this setup does not use DHCP or PXE provisioning.

.. _rs_cloud_preparation:

Preparation
-----------

Acquire a Linux Cloud or Virtual instance.  It should have at least 4 Gb or RAM and 20 Gb of storage available.

You *must* provide access to TCP/8092 (or override the default port) to access Digital Rebar.  Openning ports varies depending on the provider.

For the Cloud Wrapper provisioning, Docker or Podman must be installed.  If you use the self-runner flags and ``bootstrap-advanced`` workflow then Digital Rebar will install the available container system during bootstrapping.

.. _rs_cloud_install:

Install
-------

To begin, execute the following commands in an SSH session or during instance cloud-init process:

.. code-block:: bash

    curl -fsSL get.rebar.digital/stable | bash -s -- install --systemd --version=stable

The command will download the stable Digital Rebar (the ``systemctl`` service name is ``dr-provision``) bundle and checksum from github, extract the files, verify prerequisites are installed, and create needed directories and links under ``/var/lib/dr-provision``.  The ``--systemd`` and ``--version`` flags included for clarity, they are not required for this install.

The `install <http://get.rebar.digital/stable/>`_ script used by our installs has many additional options including ``remove`` that are documented in its help and explored in other install guides.

Once the installation script completes, a Digital Rebar endpoint will be running your instance!

Follow the steps in :ref:`rs_qs_ux_bootstrap` to registered your Digital Rebar endpoint.

.. _rs_cloud_provisioning:

Cloud Wrappers
--------------

Once Digital Rebar is running, it is the same as any other installation; however, different provisioning utilities are required because PXE and Netboot are not available.  In these cases, Digital Rebar will use the cloud providers' APIs to create and destroy machines.

Instead of upload ISOs for provisioning, Cloud instances should review :ref:`rs_cp_cloud_wrappers` for details about building cloud provisioning using Digital Rebar contexts.

.. _rs_cloud_cleanup:

Clean Up
--------

Once you are finished exploring Digital Rebar Provision in the cloud, remove the instance.