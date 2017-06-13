.. Copyright (c) 2017 RackN Inc.
.. Licensed under the Apache License, Version 2.0 (the "License");
.. Digital Rebar Provision documentation under Digital Rebar master license
.. index::
  pair: Digital Rebar Provision; Quickstart

.. _rs_quickstart:

Quick Start
~~~~~~~~~~~

_For full install, please see :ref:`rs_install`!_

Typically, Curl and Bash are unsafe, but they are simple and quick.

  ::

    curl -fsSL https://raw.githubusercontent.com/digitalrebar/provision/master/tools/install.sh | bash -s -- --isolated install

This will pull the latest code bundle and checksum from github, extract the code files,
make sure prerequisites are installed, and create some initial directories and links.

It should display something like this:

  ::

    # Run the following commands to start up dr-provision in a local isolated way.
    # The server will store information and serve files from the drp-data directory.

    sudo ./dr-provision --static-ip=<IP of an Interface> --file-root=`pwd`/drp-data/tftpboot --data-root=drp-data/digitalrebar &

    # Once dr-provision is started, this command will gather and upload the tools required to
    # do discovery-based machine management

    tools/discovery-load.sh

The sudo command will run an instance of Digital Rebar Provision that uses the drp-data
directory for object and file storage.  Additionally, *dr-provision* will attempt
to use the IP address best suited for client interaction, however if that detection fails, the IP
address specified in by *--static-ip* will be used.

.. note:: On MAC DARWIN there are two additional steps. First, use the ``--static-ip=`` flag to help the service understand traffic targets.  Second, you may have to add a route for broadcast addresses to work.  This can be done with the following comand.  The 192.168.100.1 is the IP address of the interface that you want to send messages through. The install script will make suggestions for you.

The default username & password is ``rocketskates & r0cketsk8ts``.

The *tools/discovery-load.sh* command will use the default credentials to install
the discovery, sledgehammer, and local boot environments.  This will download the
sledgehammer tarball and upload into Digital Rebar Provision.  It will also change the
default and unknown boot environments to do dynamic discovery.  This script needs to be
run after beginning *dr-provision*.

.. note:: Quickstart does NOT create a production deployment and the deployment will NOT restart on failure or reboot.

* Remember to check :ref:`rs_install` for general install information.
* Remember to check :ref:`rs_arch_ports` if there are port access issues.


Videos
------

We constantly update and add videos to the
`DR Provision Playlist <https://www.youtube.com/playlist?list=PLXPBeIrpXjfilUi7Qj1Sl0UhjxNRSC7nx>`_
so please check.

Here are quick start specific videos:

  * Mac OSX https://youtu.be/uUWU-4ObGIY
  * Linux https://youtu.be/MPkGCiakXPs
