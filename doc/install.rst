.. Copyright (c) 2017 RackN Inc.
.. Licensed under the Apache License, Version 2.0 (the "License");
.. Digital Rebar Provision documentation under Digital Rebar master license
.. index::
  pair: Digital Rebar Provision; Install

.. _rs_install:

Install
~~~~~~~

The install script does the following steps (in a slightly different order).  See :ref:`rs_quickstart` for details about the script.

Get Code
--------

The code is delivered by zip file with a sha256sum to validate contents.  These are in github under the
`releases <https://github.com/digitalrebar/provision/releases>`_ tab for the Digital Rebar Provision project.

There are at least 3 releases to choose from:

  * **tip** - This is the most recent code.  This is the latest build of master.  It is bleeding edge and while the project attempts to be very stable with master, it can have issues.
  * **stable** - This is the most recent **stable** code.  This is a tag that tracks the version-based tag.
  * **v3.0.0** - There will be a set of Semantic Versioning (aka semver) named releases.

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

Configuration Options
---------------------

Using ``dr-provision --help`` will provide the most complete list of configuration options.  The following common items are provided for reference.  Please note these may change from version to version, check the current scripts options with the ``--help`` flag to verify current options.

  ::

      --version                Print Version and exit
      --disable-provisioner    Disable provisioner
      --disable-dhcp           Disable DHCP
      --static-port=           Port the static HTTP file server should listen on (default: 8091)
      --tftp-port=             Port for the TFTP server to listen on (default: 69)
      --api-port=              Port for the API server to listen on (default: 8092)
      --dhcp-port=             Port for the DHCP server to listen on (default: 67)
      --backend=               Storage backend to use. Can be either 'consul' or 'directory' (default: directory)
      --data-root=             Location we should store runtime information in (default: /var/lib/dr-provision)
      --static-ip=             IP address to advertise for the static HTTP file server (default: 192.168.124.11)
      --file-root=             Root of filesystem we should manage (default: /var/lib/tftpboot)
      --dhcp-ifs=              Comma-seperated list of interfaces to listen for DHCP packets
      --debug-bootenv=         Debug level for the BootEnv System - 0 = off, 1 = info, 2 = debug (default: 0)
      --debug-dhcp=            Debug level for the DHCP Server - 0 = off, 1 = info, 2 = debug (default: 0)
      --debug-renderer=        Debug level for the Template Renderer - 0 = off, 1 = info, 2 = debug (default: 0)
      --tls-key=               The TLS Key File (default: server.key)
      --tls-cert=              The TLS Cert File (default: server.crt)

Prerequisites
-------------

**dr-provision** requires two applications to operate correctly, **bsdtar** and **7z**.  These are used to extract the contents
of iso and tar images to be served by the file server component of **dr-provision**.  The ``install.sh`` script will attempt to ensure these packages are installed by default.  However, if you are installing via manual process or baking your own installer, you must ensure these prerequisistes are met.

For Linux, the **bsdtar** and **p7zip** packages are required.

.. admonition:: ubuntu

  sudo apt-get install -y bsdtar p7zip-full

.. admonition:: centos/redhat

  sudo yum install -y bsdtar p7zip

.. admonition:: Darwin

  The new package, **p7zip** is required, and **tar** must also be updated.  The **tar** program on Darwin is already **bsdtar**

  * 7z - install from homebrew: ``brew install p7zip``
  * libarchive - update from homebrew to get a functional tar: ``brew install libarchive --force ; brew link libarchive --force``

At this point, the server can be started.

.. note:: In a future release, the required packages may be removed, which will help ensure cross-platform compatibility without relying on these external dependencies.

Running The Server
------------------

Additional support materials in :ref:`rs_faq`.

The **install.sh** script provides two options for running **dr-provision**.

The default values install the server and cli in /usr/local/bin.  It will also put a service control file in place.  Once that finishes, the appropriate service start method will run the daemon.  The **install.sh** script prints out the command to run
and enable the service.  The method described in the :ref:`rs_quickstart` can be used to deploy this way if the
*--isolated* flag is removed from the command line.  Look at the internals of the **install.sh** script to see what
is going on.

.. note:: The default location for storing runtime information is ``/var/lib/dr-provision`` unless overridden by ``--data-root``

Alternatively, the **install.sh** script can be passed the *--isolated* flag and it will setup the current directory
as an isolated "test drive" environment.  This will create a symbolic link from the bin directory to the local top-level
directory for the appropriate OS/platform, create a set of directories for data storage and file storage, and
display a command to run.  This is what the :ref:`rs_quickstart` method describes.

The default username & password used for administering the *dr-provision* service is:
  ::

    username: rocketskates
    password: r0cketsk8ts

Please review `--help` for options like disabling services, logging or paths.

.. note:: sudo may be required to handle binding to the TFTP and DHCP ports.

Once running, the following endpoints are available:

* https://127.0.0.1:8092/swagger-ui - swagger-ui to explore the API
* https://127.0.0.1:8092/swagger.json - API Swagger JSON file
* https://127.0.0.1:8092/api/v3 - Raw api endpoint
* https://127.0.0.1:8092/ - Redirects to RackN Portal (available for community use)
* http://127.0.0.1:8091 - Static files served by http from the *test-data/tftpboot* directory
* udp 69 - Static files served from the test-data/tftpboot directory through the tftp protocol
* udp 67 - DHCP Server listening socket - will only serve addresses when once configured.  By default, silent.
* udp 4011 - BINL Server listening socket - will only serve bootfiles when once configured.  By default, silent.

The API, File Server, DHCP, BINL,  and TFTP ports can be configured, but DHCP, BINL, and TFTP may not function properly on non-standard ports.

If the SSL certificate is not valid, then follow the :ref:`rs_gen_cert` steps.

.. note:: On MAC DARWIN there is one additional step. You may have to add a route for broadcast addresses to work.  This can be done with the following comand.  The 192.168.100.1 is the IP address of the interface that you want to send messages through. The install script will make suggestions for you.

  ::

    sudo route add 255.255.255.255 192.168.100.1

Production Deployments
----------------------

The following items should be considered for production deployments.  Recommendations may be missing so operators should use their best judgement.

Start DRP Without Root (or sudo)
================================

If you are using DHCPD and TFTPD services of DRP, you will need to be able to bind to port 67 and 69 (respectively).  Typically Unix/Linux systems require root privileges to do this.  DRP doesn't start as root, and then drop privileges with a ``fork()`` to another less privileged user by default.

To enable DRP endpoint to run as a non-privileged user and ensure a higher level of security, it's possible to use the Linux "*setcap*" (Capabilities) system to assign rights for the *dr-provision* binary to open low numbered (privileged) ports.  The process is relatively simple, but does (clearly/obviously) require root permissions initially to enable the capabilities for the binary.  Once the capabilities have been set, the *dr-provision* binary can be run as a standard user.

To enable any non-privileged user to start up the dr-provision binary and bind to privileged ports 67 and 69, do the following:

# in "isolated" mode, as the user you installed DRP as:
  ::

    sudo setcap "cap_net_raw,cap_net_bind_service=+ep" $HOME/bin/linux/amd64/dr-provision

or, in "production" mode:
  ::

    sudo setcap "cap_net_raw,cap_net_bind_service=+ep" /usr/local/bin/dr-provision

Start the "dr-provision" binary as an ordinary user, and now it will have permission to bind to privileged ports 67 and 69.

.. note:: The *setcap* command must reference the actual binary itself, and can not be pointed at a symbolic link.  Additional refinement of the capabilities may be possible.  For extremely security conscious setups, you may want to refer to the StackOverflow discussion (eg setting capabilities on a per-user basis, etc.):
  https://stackoverflow.com/questions/1956732/is-it-possible-to-configure-linux-capabilities-per-user

.. note:: You must run the *setcap* command after very upgrade of DRP, the *setcap* tracks the binary and if it changes, you must rerun for the new binary.

System Logs
===========

The Digital Rebar Provision service logs by sending output to standard error.  To capture system logs, SystemD (or Docker) should be configured to direct this output to the desired log management infrastructrure.

Job Log Rotation
================

If you are using the jobs system, Digital Rebar Provision stores job logs based on the directory configuration of the system.  This data is considered compliance related information; consequently, the system does not automatically remove these records.

Operators should set up a job log rotation mechanism to ensure that these logs to not exhaust available disk space.

Removal of Digital Rebar Provision
==================================

To remove Digital Rebar Provision, you can use the *tools/install.sh* script to remove programs for a ``production`` installs.  The *tools/install.sh* script should be run as root or under sudo unless the ``setcap`` process was used.

  ::

    tools/install.sh remove

To remove programs and data use.

  ::

    tools/install.sh --remove-data remove

For *iolated* installs, remove the directory used to contain the isolated install.  In the example above, the directory *dr-provision-install* was used to isolate the install process.  A command like this would clean up the system.

  ::

    sudo rm -rf dr-provision-install


Running the RackN UX Locally
============================

Setting up DRP to host the RackN UX locally is trivial.  The DRP server includes an embedded web server that can host the UX files from a local directory.  The RackN UX can also be set up using any other HTTP server, however this document only addresses the setup related to using DRP as the HTTP server.

The RackN UX uses the rackn-license content pack for entitlements so no external login to the RacKN SaaS is required.

The RackN UX will still attempt to connect the RackN SaaS for updates and the catalog; however, the system will operate even if these calls fail.  This can be turned off by setting a parameter in the global profile, `ux-air-gap`, to `true`.

Setup
+++++

Before starting, you'll need a copy of the RackN UX and to have installed a `rackn-license.json` content package in the DRP server.  These items require a current RackN license - using them without a valid enterprise or trial license is a copyright violation.

Extract the RackN UX files into a directory named `ux` at the same level as the `drp-data` directory.  The account running your `dr-server` must have read permission for this directory.

It is OK to use a different directory - the different directory can be specified with the `--local-ui` command line option for dr-provision.  The option specifies the directory containing the UX files.  If the path is relative, it will be assumed to be relative to the `data-root` option.


Running the UX from DRP
+++++++++++++++++++++++

By unpacking the files in the `ux` directory within the `data-root` directory or specifying the `--local-ui` option, the DRP endpoint will serve that directory as `/local-ui`.

If you are using the default port, you can access the local UX from `https://127.0.0.1:8092/local-ui`.

The endpoint will detect file changes so no restart is required if you update or change the RackN UX files.


Redirecting URL
+++++++++++++++

If you are hosting a local UX, you should change the DRP endpoint UX redirect.  This is the site that is presented if you visit the DRP endpoints root URL.  To use the local-ui, add `--ui-url=/local-ui` to the `dr-provision` command line arguments.


* Air Gap mode - the RackN UX disables all external calls and only operates against the local DRP endpoint.

