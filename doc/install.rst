.. Copyright (c) 2017 RackN Inc.
.. Licensed under the Apache License, Version 2.0 (the "License");
.. Digital Rebar Provision documentation under Digital Rebar master license
.. index::
  pair: Digital Rebar Provision; Install

.. _rs_install:

Install
~~~~~~~

The install script does the following steps (in a slightly different order).

Get Code
--------

The code is delivered by zip file with a sha256sum to validate contents.  These are in github under the
`releases <https://github.com/digitalrebar/provision/releases>`_ tab for the Digital Rebar Provision project.

There are at least 3 releases to choose from:

  * **tip** - This is the most recent code.  This is the latest build of master.  It is bleeding edge and while the project attempts to be very stable with master, it can have issues.
  * **stable** - This is the most recent **stable** code.  This is a tag that tracks the version-based tag.
  * **v3.0.0** - There will be a set of semver named releases.

Previous releases will continue to be available in tag/release history.  For additional information, see
:ref:`rs_release_process`.

When using the **install.sh** script, the version can be specified by the **--drp-version** flag,
e.g. *--drp-version=v3.0.0*.

An example command sequence for Linux would be:

  ::

    mkdir dr-provision-install
    cd dr-provision-install
    curl -fsSL https://github.com/digitalrebar/provision/releases/download/tip/dr-provision.zip -o dr-provision.zip
    curl -fsSL https://github.com/digitalrebar/provision/releases/download/tip/dr-provision.sha256 -o dr-provision.sha256
    sha256sum -c dr-provision.sha256
    unzip dr-provision.zip

At this point, the **install.sh** script is available in the **tools** directory.  It can be used to continue the process or
continue following the steps in the next sections.  *tools/install.sh --help* will provide help and context information.

Prerequisites
-------------

**dr-provision** requires two applications to operate correctly, **bsdtar** and **7z**.  These are used to extract the contents
of iso and tar images to be served by the file server component of **dr-provision**.

For Linux, you will need the **bsdtar** and **p7zip** packages.

.. admonition:: ubuntu

  sudo apt-get install -y bsdtar p7zip-full

.. admonition:: centos/redhat

  sudo yum install -y bsdtar p7zip

.. admonition:: Darwin

  You will need the following new package, **p7zip**, and you will need to update **tar**.  The **tar** program on Darwin
  is already **bsdtar**.

  * 7z - install from homebrew: brew install p7zip
  * libarchive - update from homebrew to get a functional tar: brew install libarchive

At this point, you may start the server.

Running The Server
------------------

Additional support materials in :ref:`rs_faq`.

The **install.sh** script provides two options for running **dr-provision**.  The default values install the
server and cli in /usr/local/bin.  It will also put a service control file in place.  Once that completes,
the appropriate service start method will run the daemon.  The **install.sh** script prints out the command to run
and enable the service.  The method described in the :ref:`rs_quickstart` can be used to deploy this way if the
*--isolated* flag is removed from the command line.  Look at the internals of the **install.sh** script to see what
is going on.

Alternatively, the **install.sh** script can be provided the *--isolated* flag and it will setup the current directory
as an isolated "test drive" environment.  This will create a symbolic link from the bin directory to the local top-level
directory for the appropriate OS/platform, create a set of directories for data storage and file storage, and
display a command to run.  This is what the :ref:`rs_quickstart` method describes.

The default username:password is ``rocketskates:r0cketsk8ts``.

Please review `--help` for options like disabling services, logging or paths.

.. note:: sudo may be required to handle binding to the TFTP and DHCP ports.

Once running, the following endpoints are available:

* https://127.0.0.1:8092/swagger-ui - swagger-ui to explore the API
* https://127.0.0.1:8092/swagger.json - API Swagger JSON file
* https://127.0.0.1:8092/api/v3 - Raw api endpoint
* https://127.0.0.1:8092/ui - User Configuration Pages
* https://127.0.0.1:8091 - Static files served by http from the *test-data/tftpboot* directory
* udp 69 - Static files served from the test-data/tftpboot directory through the tftp protocol
* udp 67 - DHCP Server listening socket - will only serve addresses when once configured.  By default, silent.

The API, File Server, DHCP, and TFTP ports can be configured, but DHCP and TFTP may not function properly on non-standard ports.

If your SSL certificate is not valid, then follow the :ref:`rs_gen_cert` steps.

.. note:: On Darwin, you may have to add a route for broadcast addresses to work.  This can be done with the following comand.  The 192.168.100.1 is the IP address of the interface that you want to send messages through. The install script will make suggestions for you.

  ::

    sudo route add 255.255.255.255 192.168.100.1


Configuring the Server
~~~~~~~~~~~~~~~~~~~~~~

Digital Rebar Provision provides both DHCP and Provisioning services but can be run with either disabled.  This allows users to work in environments with existing DHCP infrastructure or to use Digital Rebar Provision as an API driven DHCP server.

DHCP Server (subnets)
---------------------

The DHCP server is configured be enabling Subnets that serve IPs and/or additional configuration information.  It is possible to run the DHCP server using only pre-defined IP Reservations or allow the DHCP server to create IP Leases dynamically.

The DHCP server had two primary models

#. DHCP Listeners can be set on an IP for each server interface.  These listeners will respond to DHCP broadcasts on the matching network(s).  Operators should ensure that no other DHCP servers are set up on the configured subnets.

#. DHCP Relay allows other DHCP listeners to forward requests to the Digital Rebar Provision server.  In this mode, the server is passive and can easily co-exist with other DHCP servers.  This mode works with the Provisioner by setting the many optional parameters (like next boot) that are needed for PXE boot processes.

Provisioner (bootenvs)
----------------------

The Provisioner is a combination of several services and a template expansion engine.  The primary model is a boot environment (BootEnv) that contains crtical metadata to describe an installation process.  This metadata includes templates that are dynamically expanded when machines boot.

Digital Rebar Provision CLI has a process that combines multiple calls to install BootEnvs.  The following steps will configure a system capable to :ref:`rs_provision_discovered`.

  ::

    cd assets
    drpcli bootenvs install bootenvs/sledgehammer.yml
    drpcli bootenvs install bootenvs/discovery.yml
    drpcli bootenvs install bootenvs/local.yml
    drpcli prefs set unknownBootEnv "discovery" defaultBootEnv "sledgehammer"

.. note:: The tools/discovery_load.sh script does this with the default credentials.


Default Template Identity 
-------------------------

_These settings apply to TEMPLATES not the API._

The default password for the default o/s templates is **RocketSkates**

The default user for the default ubuntu/debian templates is **rocketskates**
