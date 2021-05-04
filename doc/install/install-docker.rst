.. Copyright (c) 2017 RackN Inc.
.. Licensed under the Apache License, Version 2.0 (the "License");
.. Digital Rebar Provision documentation under Digital Rebar master license
.. index::
  pair: Digital Rebar Provision; Install Docker
  pair: Digital Rebar Provision; Install Container
  pair: Digital Rebar Provision; Install Podman

.. _rs_install_docker:

Container Install (Docker/Podman)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

This install guide provides a streamlined installation via a container system such as Docker or Podman.  For trial users, this install method minimizes the installation footprint of Digital Rebar.

Note: while RackN does support this installation approach, we *strongly recommend* using systemd installation for production and multi-site installations.

Other installation paths:

* :ref:`rs_quickstart` is a basic `systemd <https://en.wikipedia.org/wiki/Systemd>`_ install for new users
* :ref:`rs_install` details more complex installs including offline/airgap.
* :ref:`rs_install_dev` for developers running DRP interactively
* :ref:`rs_install_cloud` is non-PXE / Cloud-Only installation process (no DHCP required)
* `Edge Lab with RPi <http://edgelab.digital>`_ is self-contained Digital Rebar inexpensive lab using Raspberry Pi computers.

Each of these environments requires careful setup up of your network environment and consideration with regard to competing DHCP services.  The setup of these environments is outside the scope of this document.

You must install Digital Rebar to use it, there is no SaaS version.  :ref:`rs_self_managed_why`

.. _rs_docker_preparation:

Prepare Your Environment
------------------------

This guide will install the Digital Rebar onto a Linux or MacOS system using a downloaded container pre-populated with Digital Rebar.  It will attach to port 8092 (and these :ref:`rs_arch_ports`).  You must have installed Docker or Podman before starting this installation process.

RackN maintains containers on Dockerhub for stable versions of Digital Rebar.  Builds of the latest (aka tip) version may run behind latest.  If you are planning frequent updates, please consult :ref:`rs_install_dev`.

Note: some container configurations do not forward UDP packets correctly.  This will interfer with DHCP and TFTP services required for PXE provisioning.

.. _rs_docker_install:

Install
-------

To begin, use the RackN install.sh script to start your container.  We recommend using process instead simply pulling the container to ensure correct configuration and port mapping.

  ::

    curl -fsSL get.rebar.digital/stable | bash -s -- --container --version=stable install

The command will pull the *stable* ``dr-provision`` in a container from Dockerhub then verify prerequisites are installed, map the correct ports and create some initial directories and links.


The `install <http://get.rebar.digital/stable/>`_ script used by our installs has many additional container configuration options including setting the type, name, restart, volume, netns, and more.  Users planning to long term installations with containers should carefully review these installation options.

Once the installation script completes, a Digital Rebar endpoint will be running your local system!

.. _rs_docker_next_steps:

Back to Regular Install
-----------------------

Once Digital Rebar is running in a container mode, it is exactly the same as any other installation

* :ref:`rs_qs_license`
* :ref:`rs_qs_ux_bootstrap`
* :ref:`rs_qs_cli_bootstrap`
* :ref:`rs_qs_first_machine`
* :ref:`rs_qs_next_steps`


.. _rs_docker_advanced:

Advanced Container Deployments
------------------~~~~~~~~~---

Installation is perforemed with the ``install.sh`` script with the ``--container`` flag and associated options.  Here are some of the options (please check the latest installer script for updates/details):

  ::

    --container             # Force to install as a container, not zipfile
    --container-type=<string>
                            # Container install type, defaults to "docker"
    --container-name=<string>
                            # Set the "docker run" container name, defaults to "drp"
    --container-restart=<string>
                            # Set the Docker restart option, defaults to "always"
                            # options are:  no, on-failure, always, unless-stopped
                            * see: https://docs.docker.com/config/containers/start-containers-automatically/
    --container-volume=<string>
                            # Volume name to use for backing persistent storage, default "drp-data"
    --container-registry="drp.example.com:5000"
                            # Alternate registry to get container images from, default "index.docker.io"
    --container-env="<string> <string> <string>"
                            # Define a space separated list of environment variables to pass to the
                            # container on start (eg "RS_METRICS_PORT=8888 RS_DRP_ID=fred")
                            # see 'dr-provision --help' for complete list of startup variables
    --container-netns="<string>"
                            # Define Network Namespace to start container in. Defaults to "host"
                            # If set to empty string (""), then disable setting any network namespace

.. note:: WARNING: If you intend to Upgrade DRP in a container based scenarios, it iS IMPORTANT that you retain a copy of the installation command line flags you use for install time.  These flags will have to be specified for the upgrade command to work correctly.

Container based installations will by default name the container ``drp``, and the data backing volume ``drp-data``.  You can change these with appropriate flags.  The writable data store is located in the backing volume, which helps isolate the binary/service environment from the writable content.  See the :ref:`rs_upgrade_container` for more details.

The ``dr-provision`` service binary utilizes environment variables as a mechanism to support customization of the runtime of the service.  This also allows the operator to start the container and modify the runtime via the use of passing Environment variables in to the container.  Here is an example:

  ::

    ./install.sh install --container --container-restart=always --container-netns=host --container-env="RS_METRICS_PORT=8888"

This example modifies the Metrics port to be changed from the default of ``8080`` to relocate to port ``8888``.  See ``dr-provision --help`` for a list of all environment variable options that can be set.


.. _rs_docker_cleanup:

Clean Up
--------

Once you are finished exploring Digital Rebar Provision in container mode, the system can cleaned or reset by removing container.