.. Copyright (c) 2017 RackN Inc.
.. Licensed under the Apache License, Version 2.0 (the "License");
.. DigitalRebar Provision documentation under Digital Rebar master license
.. index::
  pair: DigitalRebar Provision; Quickstart

.. _rs_quickstart:

Quick Start
~~~~~~~~~~~

TL;DR - Curl/Bash is NOT safe, but it is fun and easy.

  ::

    curl -fsSL https://raw.githubusercontent.com/digitalrebar/provision/master/tools/install.sh | bash -s -- --isolated install

This will pull the latest code bundle and checksum from github, extract the code files, make sure prerequistes are installed,
and create some initial directories and links.  It should display something like this:

  ::

    Run the following commands to start up dr-provision in a local isolated way.
    The server will store information and server files in the drp-data directory.

    sudo ./dr-provision --static-ip=<IP of an Interface> --file-root=`pwd`/drp-data/tftpboot --data-root=drp-data/digitalrebar &
    tools/discovery-load.sh

The sudo command will run an instance of Digital Rebar Provision that uses the drp-data directory for object storage and file storage.
Additionally, *dr-provision* will attempt to use the best IP address for client's to talk to it, but if that detection fails, the
IP address specified in by *--static-ip* will be used.

The default user:password is 'rocketskates:r0cket8ts'

The *tools/discovery-load.sh* command will use the default credentials to install the discovery, sledgehammer, and local boot
environments.  This will download the sledgehammer tarball and upload into Digital Rebar Provision.  It will also change the
default and unknown boot environments to do dynamic discovery.  This script needs to be run after starting *dr-provision*.

The default username:password is ``rocketskates:r0cketsk8ts``.

.. note:: This is NOT a production deployment and will NOT restart on failure or reboot.


Videos
------

We constantly update and add videos to the `DR Provision Playlist <https://www.youtube.com/playlist?list=PLXPBeIrpXjfilUi7Qj1Sl0UhjxNRSC7nx>`_ so please check.

Here are quick start specific videos:

  * Mac OSX https://youtu.be/uUWU-4ObGIY
  * Linux https://youtu.be/MPkGCiakXPs