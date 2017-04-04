.. Copyright (c) 2017 RackN Inc.
.. Licensed under the Apache License, Version 2.0 (the "License");
.. Rocket Skates documentation under Digital Rebar master license
.. index::
  pair: Rocket Skates; Install

.. _rs_install:

Install
~~~~~~~

There are prerequisites for the system to function.

Linux:
* bsdtar - from your local package manager
  * on ubuntu: apt-get install bsdtar
  * on centos/redhat: yum install bsdtar
* 7z - from your local package manager
  * on ubuntu: apt-get install p7zip
  * on centos/redhat: yum install p7zip

Darwin:
* bash4 - install from homebrw: brew install bash
* 7z - install from homebrew: brew install p7zip
* libarchive - update from homebrew to get a functional tar: brew install libarchive


Running The Server
------------------

Additional support materials in :ref:`rs_faq`.

To run a local copy that will use the local filesystem as a storage area, do the following:

  ::

    cd test-data
    sudo ../rocket-skates

Please review `--help` for options like disabling services, logging or paths.

.. note:: I need the sudo to bind the tftp port.  This is configurable, i.e.  *--tftp-port=30000*  

The following pieces endpoints are available:

* https://127.0.0.1:8092/swagger-ui - swagger-ui to explore the API
* https://127.0.0.1:8092/swagger.json - API Swagger JSON file
* https://127.0.0.1:8092/api/v3 - Raw api endpoint
* https://127.0.0.1:8092/ui - User Configuration Pages
* https://127.0.0.1:8091 - Static files served from the test-data/tftpboot directory
* udp 69 or *--tftp-port* - Static files served from the test-data/tftpboot directory through the tftp protocol
* udp 67 - DHCP Server listening socket - will only server addresses when once configured.  By default, silent.

.. note:: If your SSL certificate is not valid, then follow the :ref:`rs_gen_cert` steps.

.. note:: On OSX, you may have to add a route for broadcast addresses to work.  This can be done with the following
comand.  The 192.168.100.1 is the IP address of the interface that you want to send messages through.

  ::

    sudo route add 255.255.255.255 192.168.100.1


.. _rs_gen_cert:

Generate Certificate
--------------------

Sometimes the cert/key pair in the github tree is corrupt or not sufficient for the environment.  You can run
the following command to rebuild a local cert/key pair.

  ::

    sudo openssl req -new -x509 -keyout server.key -out server.crt -days 365 -nodes


You may need to install the openssl tools.


Configuring the Server
~~~~~~~~~~~~~~~~~~~~~~

Rocket Skates provides both DHCP and Provisioning services but can be run with either disabled.  This allows users to work in environments with existing DHCP infrastructure or to use Rocket Skates as an API driven DHCP server.

DHCP Server
-----------

Provisioner
-----------



Download the service & cli
  ::
    curl -o
    curl -o


Install Boot Environments
Upload Templates
Set Proferences
  ::
    cd assets
    ../rscli bootenvs install bootenvs/sledgehammer.yml 
    ../rscli bootenvs install bootenvs/discovery.yml 
    ../rscli bootenvs install bootenvs/local.yml 
    ../rscli templates upload templates/local-elilo.tmpl as 5  ../rscli templates upload templates/local-pxelinux.tmpl as local-pxelinux.tmpl
    ../rscli templates upload templates/local-ipxe.tmpl as 
 local.elilo.tmpl
    ../rscli prefs set unknownBootEnv to "discovery"
local-ipxe.tmpl
 
You can also review the UX via https://127.0.0.1:8092.
