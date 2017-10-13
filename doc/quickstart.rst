.. Copyright (c) 2017 RackN Inc.
.. Licensed under the Apache License, Version 2.0 (the "License");
.. Digital Rebar Provision documentation under Digital Rebar master license
.. index::
  pair: Digital Rebar Provision; Quickstart

.. _rs_quickstart:

Quick Start
~~~~~~~~~~~

This quick start guide provides a basic installation and start point for further exploration.  The guide has been designed for UNIX systems: Mac OS, Linux OS, Linux VMs and Linux Packet Servers.  The guide employs Curl and Bash commands which are not typically considered safe, but they do provide a simple and quick process for start up.

For a full install, please see :ref:`rs_install`

To begin, execute the following command in a shell or terminal: 
  ::

    curl -fsSL https://raw.githubusercontent.com/digitalrebar/provision/stable/tools/install.sh | bash -s -- --isolated install
    
.. note:: If you want to try the latest code, you can checkout the development tip using ``curl -fsSL https://raw.githubusercontent.com/digitalrebar/provision/tip/tools/install.sh | bash -s -- --isolated install --drp-version=tip``

The command will pull the latest code bundle and checksum from github, extract the code files,
verify prerequisites are installed, and create some initial directories and links.

The terminal should then display something like this:

  ::

    # Run the following commands to start up dr-provision in a local isolated way.
    # The server will store information and serve files from the drp-data directory.

    sudo ./dr-provision --static-ip=<IP of an Interface> --file-root=`pwd`/drp-data/tftpboot --data-root=drp-data/digitalrebar &

    # Once dr-provision is started, this command will gather and upload the tools required to
    # do discovery-based machine management

    tools/discovery-load.sh

The next step is to execute the sudo command which will run an instance of Digital Rebar Provision that uses the drp-data
directory for object and file storage.  Additionally, *dr-provision* will attempt
to use the IP address best suited for client interaction, however if that detection fails, the IP
address specified in by *--static-ip* will be used.  After Digital Rebar Provision has started a prompt for a username and password will appear.  

.. note:: On MAC DARWIN there are two additional steps. First, use the ``--static-ip=`` flag to help the service understand traffic targets.  Second, you may have to add a route for broadcast addresses to work.  This can be done with the following comand.  The 192.168.100.1 is the IP address of the interface that you want to send messages through. The install script will make suggestions for you.

The default username & password is ``rocketskates & r0cketsk8ts``.

With Digital Rebar Provision running it is now time to install the specialized Digital Rebar Provision images, 
and the required boot environments.

The *tools/discovery-load.sh* command will use the default credentials to install
the discovery, sledgehammer, and local boot environments.  This will download the
sledgehammer tarball and upload into Digital Rebar Provision.  It will also change the
default and unknown boot environments to do dynamic discovery.  This script needs to be
run after beginning *dr-provision*.

When the *tools/discovery-load.sh* script finishes Digital Rebar Provision will be installed and operational.  


.. note:: This quick start guide does NOT create a production deployment and the deployment will NOT restart on failure or reboot.

* Remember to check :ref:`rs_install` for general install information.
* Remember to check :ref:`rs_arch_ports` if there are port access issues.


Videos
------
We constantly update and add videos to the
`DR Provision 3.1 Playlist <https://www.youtube.com/playlist?list=PLXPBeIrpXjfj5_8Joyehwq1nnaYSPCnmw>`_
so please check to make sure you have the right version!
