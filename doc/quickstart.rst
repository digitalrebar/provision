.. Copyright (c) 2017 RackN Inc.
.. Licensed under the Apache License, Version 2.0 (the "License");
.. Digital Rebar Provision documentation under Digital Rebar master license
.. index::
  pair: Digital Rebar Provision; Quickstart

.. _rs_quickstart:

Quick Start
~~~~~~~~~~~

.. note::  We HIGHLY recommend using the ``latest`` version of the documentation, as it contains the most up to date information.  Use the version selector in the lower     right corner of your browser.

This quick start guide provides a basic installation and start point for further exploration.  The guide has been designed for UNIX systems: Mac OS, Linux OS, Linux VMs and Linux Packet Servers.  While it is possible to install and run Digital Rebar Provision on a Windows instance, we do not cover that here.  The guide employs Curl and Bash commands which are not typically considered safe, but they do provide a simple and quick process for start up.

It is possible to install on hypervisors and in virtualized environments (eg. VirtualBox, VMware Workstation/Fusion, KVM, etc.).  Each of these environments requires careful setup up of your network environment and consideration with regard to competing DHCP services.  The setup of these environments is outside the scope of this document.

For a full install, please see :ref:`rs_install`

Overview
--------

  * Read :ref:`rs_qs_preparation` steps below
  * :ref:`rs_qs_install` DRP Endpoint (in "isolated" mode)
  * :ref:`rs_qs_start` daemon
  * :ref:`rs_qs_bootenvs` (OS media for installation)
  * :ref:`rs_qs_subnet` to answer DHCP requests
  * :ref:`rs_qs_workflow`
  * :ref:`rs_qs_defaults` for defaultBootEnv, defaultStage, unknownBootEnv, and defaultWorkflow
  * boot your first Machine and :ref:`rs_qs_first_machine` an OS on it

This document refers to the ``drpcli`` command line tool for manipulating the ``dr-provision`` service.  We do not specify any paths in the documentation.  However, in our default quickstart for *isolated* mode, the ``drpcli`` tool will NOT be installed in any system PATH locations.  You must do this, or you may use the local setup symbolic link.  For example - simply change ``drpcli`` to ``./drpcli`` in the documentation below.  Or ... copy the binary to a PATH location.

You can perform all of the actions outlined in this document via the hosted web UX Portal.  If you choose to use the web Portal, please ensure the setup is completed by reviewing the output of the ``Info & Preferences`` panel named *System Wizard* are completed successfully.  Basic `documentation on the Web Portal <http://provision.readthedocs.io/en/tip/doc/ux/portal/systemux.html#machines>`_ UX is available. 

.. _rs_qs_preparation:

Preparation
-----------

Please make sure your environment doesn't have any conflicts or issues that might cause PXE booting to fail.  Some things to note:

  * only one DHCP server on a local subnet
  * your Machines should be set to PXE boot the correct NIC (on the correct provisioning network interface)
  * if you customize Reservations - you must also add all of the correct PXE boot options (see :ref:`rs_create_reservation` )
  * you need the network information for the subnet that your target Machines will be on
  * Mac OSX may require additional setup (see notes below)
  * we rely heavily on the ``jq`` tool for use with the Command Line tool (``drpcli``) - install it if you don't have it already

.. _rs_qs_install:

Install
-------

To begin, execute the following commands in a shell or terminal:
  ::

    mkdir drp ; cd drp
    curl -fsSL get.rebar.digital/stable | bash -s -- --isolated install

.. note:: If you want to try the latest code, you can checkout the development tip using ``curl -fsSL get.rebar.digital/tip | bash -s -- --isolated --drp-version=tip install``

The command will pull the *stable* ``dr-provision`` bundle and checksum from github, extract the files, verify prerequisites are installed, and create some initial directories and links.

.. note:: By default the installer will pull in the default Community Content packages.  If you are going to add your own or different (eg RackN registered content), append the ``--nocontent`` flag to the end of the install command.

.. note:: The "install.sh" script that is executed (either via 'stable' or 'tip' in the initial 'curl' command), has it's own version number independent of the Digital Rebar Provision endpoint version that is installed (also typically called 'tip' or 'stable').  It is NOT recommend to "mix-n-match" the installer and endpoint version that's being installed.

For reference, you can download the installer (``install.sh``), and observe what the shell script is going to do (highly recommended as a prudent security caution), to do so simply:
  ::

    curl -fsSL get.rebar.digital/stable -o install.sh

Once the installer is downloaded, you can execute it with the appropriate ``install`` options (try ``bash ./install.sh --help`` for details).

It is recommended that directory is used for this process.  The ``mkdir drp ; cd drp`` command does this as the ``drp`` directory.  The directory will contain all installed and operating files. The ``drp`` directory can be anything.

Even for *production* installs (without ``--isolated``), it is recommended to run the ``install.sh`` script in a directory to contain all the install files for easy clean-up and removal if Digital Rebar Provision needs to be removed from the system.

.. _rs_qs_start:

Start Digital Rebar Provision service
-------------------------------------

Our quickstart uses *isolated* mode install, and the ``dr-provision`` service is not installed in the system path.  You need to manually start ``dr-provision`` each time the system is booted up.  The *production* mode installation (do not specify the ``--isolated`` install flag) will install in to system directories, and provide helpers to setup ``init``, ``systemd``, etc. start up scripts for the service.

Once the install has completed, your terminal should then display something like this (please use the output from YOUR install version, the below is just an example that may be out of date with the current versions output):

  ::

    # Run the following commands to start up dr-provision in a local isolated way.
    # The server will store information and serve files from the ./drp-data directory.

    sudo ./dr-provision --base-root=`pwd`/drp-data --local-content="" --default-content="" > drp.log 2>&1 &


.. note:: On MAC DARWIN there is one additional step. You may have to add a route for broadcast addresses to work.  This can be done with following command ``sudo route -n add -net 255.255.255.255 192.168.100.1`` In this example, the ``192.168.100.1`` is the IP address of the interface that you want to send messages through. The install script should make suggestions for you.

The next step is to execute the *sudo* command which will start an instance of Digital Rebar Provision service that uses the ``drp-data`` directory for object and file storage.

.. note:: Before trying to install a BootEnv, please verify that the installed BootEnvs matches the above BootEnv Names that can be installed: ``drpcli bootenvs list | jq '.[].Name'``


You may also use the RackN Portal UX by pointing your web browser to:
  ::

    https://<ip_address_of_your_endpoint>:8092/

Please note that your browser will be redirected to the RackN Portal, pointing at your newly installed Endpoint.  Use the below username/password pair to authenticate to the DRP Endpoint.  Additional capabilities and features can be unlocked by also using the RackN Portal Login (upper right "Login" blue button).

The default username & password used for administering the *dr-provision* service is:
  ::

    username: rocketskates
    password: r0cketsk8ts

.. _rs_qs_bootenvs:

Install Boot Environments (bootenvs)
------------------------------------

With Digital Rebar Provision running; it is now time to install the specialized Digital Rebar Provision content, and the required boot environments (BootEnvs).  We generally refer to this as "content".

.. note:: This documentation assumes you are using the default ``drp-community-content`` pack.

During the install step above, the installer output a message on how to install install BootEnvs.  You must install the ``sledgehammer`` BootEnv for Discovery and Workflow.  You may selectively choose to install one of the Community Content BootEnvs that you wish to install to your Machines.  To obtain a full list of Community Content supported BootEnvs, do:
  ::

    drpcli bootenvs list | jq '.[].Name'

  #. install the *sledgehammer* Boot Environment, used for discovery and provisioning workflow
  #. set default preferences for basic discovery booting
  #. download common utilties for Digital Rebar from the RackN catalog
  #. cache the CentOS Boot ISO locally <optional>
  #. cache the Ubuntu Boot ISO locally <optional>

These steps should be performed from the newly installed *dr-provision* endpoint (or via remote *drpcli* binary with the use of the ``--endpoint`` flag):

  ::

    drpcli bootenvs uploadiso sledgehammer
    drpcli prefs set defaultWorkflow discover-base unknownBootEnv discovery
    drpcli contents upload https://api.rackn.io/catalog/content/task-library
    drpcli bootenvs uploadiso ubuntu-18.04-install
    drpcli bootenvs uploadiso centos-7-install

The ``uploadiso`` command will fetch the ISO image as specified in the BootEnv JSON spec, download it, and then "explode" it in to the ``drp-data/tftpboot/`` directory for installation use.  You may optionally choose one or both of the CentOS and Ubuntu BootEnvs (or any other Community Content supported BootEnv) to install; depending on which Operating System and Version you wish to test or use.


Try the RackN UX
----------------

At this point, you can switch to the RackN Web UX by pointing your web browser to ``https://<ip_address_of_your_endpoint>:8092/``.

After you accept the endpoint's self-signed certificate, you will be automatically redirected to
the RackN portal at ``https://portal.rackn.io/`` with your endpoint IP included.
This will guide you to ``Digital Rebar Endpoint Login`` for your new endpoint's API. Clicking ``defaults`` will populate the login with default credentials.

.. note:: The RackN portal runs as a single-page app *locally* in your browser so all DRP API calls
remain behind your firewall. RackN *never* has direct access to your DRP endpoint.

.. _rs_qs_subnet:

Configure a Subnet
------------------

A Subnet defines a network boundary that the DRP Endpoint will answer
DHCP queries for.  In this quickstart, we assume you will use the
local network interface as a subnet definition, and that your Machines
are all booted from the local subnet (layer 2 boundary).  A Subnet
specification must include all of the necessary DHCP boot options to
correctly PXE boot a Machine.

.. note:: DRP supports the use of external DHCP servers, DHCP Proxy, etc.  However, this is considered an advanced topic, and not discussed in the QuickStart.

Starting with Stable release v3.7.0 and newer, Digital Rebar Provision
supports "magic" DHCP Boot Options for `next-server` and `bootfile`
(option code 67).  This means that these options should work "magically"
for you without needing to be set.

_HOWEVER_ - VirtualBox has a broken iPXE implementation, but DRP now ships with an iPXE client that will
work with ViritualBox.  The latest version of DRP will boot VirtualBox systems correctly.

If you are creating a subnet for an older version (before v3.7.0) of Digital Rebar
Provision, you must set the `next-server` to your DRP Endpoint IP Address,
and set the Option 67 value to ``lpxelinux.0`` for Legacy BIOS mode
Machines.

If you are *using VirtualBox* and DRP version before v3.10.0, you must set the `next-server`
value to the DRP Endpoint IP address _and_ the DHCP Option 67 value to ``lpxelinux.0``


.. note:: The UX will create a Subnet based on an interface of the DRP Endpoint with sane defaults - it is easier to create a subnet via the UX.

  If you are using a VirtualBox environment, and if you set the Name of the `Subnet` to ``vboxnet0``, the UX will automatically correct the Option 67 bootfile value to support the broken iPXE environment for VirtualBox networks.

  You must still set all of the remaining network values correctly in your Subnet specification, even in the UX.

To create a basic Subnet from command line we must create a JSON blob that
contains the Subnet and DHCP definitions.  Below is a _sample_ you can
use.  *PLEASE ENSURE* you modify the network parameters accordingly.
Ensure you change the network parameters according to your
environment.

  ::

    ###
    #  EXAMPLE - please modify the below values according to your environment  !!!
    ###

    echo '{
      "Name": "local_subnet",
      "Subnet": "10.10.16.10/24",
      "ActiveStart": "10.10.16.100",
      "ActiveEnd": "10.10.16.254",
      "ActiveLeaseTime": 60,
      "Enabled": true,
      "ReservedLeaseTime": 7200,
      "Strategy": "MAC",
      "Options": [
        { "Code": 3, "Value": "10.10.16.1", "Description": "Default Gateway" },
        { "Code": 6, "Value": "8.8.8.8", "Description": "DNS Servers" },
        { "Code": 15, "Value": "example.com", "Description": "Domain Name" }
      ]
    }' > /tmp/local_subnet.json

    # edit the above JSON spec to suit your environment
    #
    # for v3.6.0 and older:
    #  add a next-server after "Name" with the IP address of your DRP Endpoint, like:
    #    NextServer": "10.10.16.10",
    #
    # for v3.6.0 and older:
    #  add DHCP Option 67 to the Options map, like:
    #    { "Code": 67, "Value": "lpxelinux.0", "Description": "Bootfile" },
    #
    vim /tmp/local_subnet.json

    drpcli subnets create - < /tmp/local_subnet.json

.. note:: Option 67 (bootfile name) specifies the PXE boot file.  The `lpxelinux.0` boot file is for Legacy BIOS machines.  If you are booting a UEFI system, you will need to make more advanced changes to support UEFI boot mode. Please see the FAQ on :ref:`rs_uefi_boot_option`.  DRP v3.7.0 and newer has magic helpers to try and set the Legacy/UEFI bootfile for you, but custom usage or custom/unique PXE implementations may require changes.

.. _rs_qs_workflow:

Create a Workflow
-----------------

No Action Required!

*Workflows* define a series of *Stages* that a Machine transitions through, driven
by Digital Rebar Provision.  Not only do the drive basic Operating System 
installation, but they also allow for advanced application installation and
configuration if desired.  *Workflows* also allow for some basic power management
functions for hardware that does not support IPMI-like functions.  

For our QuickStart use case, you've automatically created three simple *Workflows*:

  #. discover-base (from community content)
  #. centos-base (from RackN task-library)
  #. ubuntu-base (from RackN task-library)

These pre-built basic workflows are imported as read only to get you running quickly.
You'll need to clone them to customize them as you learn Digital Rebar.

.. _rs_qs_defaults:

Set The Defaults 
----------------

No Action Required!

For our QuickStart use case, you've already configured preferences to use base discovery.

One of the basic safety mechanisms for newly installed DRP Endpoints, is to
prevent accidental Installation of a Machine, if it should PXE boot against 
a DRP Endpoint ... **before** you are ready for that to happen!!  So we must
first set the default actions for a few system wide preferences.  One of those
defaults will point to our Discovery Workflow (see :ref:`rs_qs_workflow`).

You can review the preferences using ``drpcli prefs get``.

Any Machine that boots will by default be placed in to the Discovery Workflow,
which will NOT install an Operating System, but will enroll the machine for
management by Digital Rebar Provision

Once set, any "unknown" or new Machines that boot against the DRP Endpoint will
be *discovered* and sit idly by waiting (`sledgehammer-wait`) for your 
next commands (eg. install an operating system).

.. _rs_qs_first_machine:

Install your first Machine
--------------------------

Content configuration is the most complex topic with Digital Rebar Provision.
The basic provisioning setup with the above "ISO" uploads, default preferences,
and simple workflows will allow you to install an operating system on the Machine
with *manual power management* (on/off/reboot etc) transitions.  

More advanced workflows and plugin_providers will allow for complete automation
workflows with complex stages and state transitions.  To keep things "quick", the
below are just bare basics, for more details and information, please see the
Content documentation section.

  1. PXE Boot your Machine

    * ensure your test Machine is on the same Layer 2 subnet as your DRP endpoint, or that you've configured your networks *IP Helper* to forward your DHCP requests to your DRP Endpoint
    * the Machine should be in the same subnet as defined in the Subnets section above (not strictly required, but this is a simplified quickstart environment!)
    * set your test machine or VM instance to PXE boot
    * power the Machine on (or reboot it) and verify from the console that the NIC begins the PXE boot process
    * verify that the DRP Endpoint responds with a DHCP lease to the Machine

    The Machine should boot in to the Sledgehammer BootEnv - which will bring
    the console to a prompt that looks like (the version signature may differ):

      ::

        Digital Rebar: Sledgehammer 6122f34b46b5b74b668d6779e33f5fcd0f44a8cc
        Kernel 3.10.0-693.21.1.el7.x86_64 on an x86_64

        d0c-c4-7a-e5-48-b6 login:

  2. Get your Machines UUID so you can set the Workflow for it

    * once your machine has booted, and received DHCP from the DRP Endpoint, it will now be "registered" with the Endpoint for installation
    * by default, DRP will NOT attempt an OS install unless you explicitly direct it to (for safety's sake!)
    * obtain your Machine's ID, you'll use it to define your BootEnv (see :ref:`rs_filter_gohai` for more detailed/cleaner syntax)

    ::

      drpcli machines list | jq '.[].Uuid'

  3. Set the Workflow to to your Operating System Workflow you defined above;
     replace *<UUID>* with your machines ID from the above command:

    ::

      # example for CentOS 7 workflow
      drpcli machines update <UUID> '{ "Workflow": "centos7" }'

      # example for Debian 9 workflow
      drpcli machines update <UUID> '{ "Workflow": "debian9" }'

      # example for Ubuntu 18.04 workflow
      drpcli machines update <UUID> '{ "Workflow": "ubuntu18" }'

  4. Reboot your Machine - it should now kick off a BootEnv install as you specified above.

    * watch the console, and you should see the appropriate installer running
    * the machine should reboot in to the Operating System you specified once install is completed

.. note:: Digital Rebar Provision is capable of automated workflow management of the boot process, power control, and much more.  This quickstart walks through the simplest process to get you up and running with a single test install.  Please review the rest of the documentation for further configuration details and information on automation of your provisioning environment.

More Advanced Workflow
----------------------

The above procedure uses manual reboot of Machines, and manual application of the BootEnv definition to the Machine for final installation.  A simple workflow can be used to achieve the same effect, but it is a little more complex to setup.  See the :ref:`rs_operation` documentation for further details.

Machine Power Management
------------------------

Fully automated provisioning control requires use of advanced RackN features (plugins) for Power Management actions.  These are done through the IPMI subsystem, with a specific IPMI plugin for a specific environments.  Some existing plugins exist for environments like:

  * bare metal - hardware based BMC (baseboard management controller) functions that implement the IPMI protocol
  * Virtual Box
  * Packet bare metal hosting provider (https://www.packet.net/)
  * Advanced BMC functions are supported for some hardware vendors (eg Dell, HP, IBM, etc)

`Contact RackN <https://www.rackn.com/company/contact-us/>`_ for additional details and information.

Isolated vs Production Install Mode
-----------------------------------

The quickstart guide does NOT create a production deployment and the DRP Endpoint service will NOT restart on failure or reboot.  You will have to start the *dr-provision* service on each system reboot (or add appropiate startup scripts).

A production mode install will install to ``/var/lib/dr-provision`` directory (by default), while an isolated install mode will install to ``$PWD/drp-data``.

For more detailed installation information, see: :ref:`rs_install`

Clean Up
--------

Once you are finished exploring Digital Rebar Provision in isolated mode, the system can cleaned by removing the directory containing the isolated install.  In the previous sections, we used ''drp'' as the directory containing the isolated install.  Removing this directory will clean up the installed files.

For production deployments, the ``install.sh`` script can be run with the ``remove`` argument instead of the ``install`` argument to clean up the system.  This will not remove the data files stored in ``/var/lib/dr-provision``, ``/etc/dr-provision``, or ``/usr/share/dr-provision``.  The ``tools/install.sh`` script is in the directory where you ran the ``install.sh`` script the first time.  The script can be also redownloaded and run through curl | bash.

  ::

    tools/install.sh remove

To additionally remove the data files, run instead:

  ::

    tools/install.sh --remove-data remove

Ports
-----

The Digital Rebar Provision endpoint service requires specific TCP Ports be accessible on the endpoint.  Please see :ref:`rs_arch_ports` for more detailed information.

If you are running in a Containerized environment, please ensure you are forwarding all of the ports appropriately in to the container.  If you have a Firewall or packet filtering service on the node running the DRP Endpoint - ensure the appropriate ports are open.


Videos
------

We constantly update and add videos to the
`DR Provision Playlist <https://www.youtube.com/playlist?list=PLXPBeIrpXjfilUi7Qj1Sl0UhjxNRSC7nx>`_
so please check to make sure you have the right version!
