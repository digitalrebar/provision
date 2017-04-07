.. Copyright (c) 2017 RackN Inc.
.. Licensed under the Apache License, Version 2.0 (the "License");
.. DigitalRebar Provision documentation under Digital Rebar master license
.. index::
  pair: DigitalRebar Provision; Install

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
    sudo ../bin/$(uname -s | tr '[:upper:]' '[:lower:]')/amd64/dr-provision --file-root=`pwd`/tftpboot --data-root=./digitalrebar

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

If your SSL certificate is not valid, then follow the :ref:`rs_gen_cert` steps.

.. note:: On OSX, you may have to add a route for broadcast addresses to work.  This can be done with the following comand.  The 192.168.100.1 is the IP address of the interface that you want to send messages through.

  ::

    sudo route add 255.255.255.255 192.168.100.1


Configuring the Server
~~~~~~~~~~~~~~~~~~~~~~

DigitalRebar Provision provides both DHCP and Provisioning services but can be run with either disabled.  This allows users to work in environments with existing DHCP infrastructure or to use DigitalRebar Provision as an API driven DHCP server.

DHCP Server (subnets)
---------------------

The DHCP server is configured be enabling Subnets that serve IPs and/or additional configuration information.  It is possible to run the DHCP server using only pre-defined IP Reservations or allow the DHCP server to create IP Leases dynamically.  

The DHCP server had two primary models

#. DHCP Listeners can be set on an IP for each server interface.  These listeners will respond to DHCP broadcasts on the matching network(s).  Operators should ensure that no other DHCP servers are set up on the configured subnets.

#. DHCP Relay allows other DHCP listeners to forward requests to the DigitalRebar Provision server.  In this mode, the server is passive and can easily co-exist with other DHCP servers.  This mode works with the Provisioner by setting the many optional parameters (like next boot) that are needed for PXE boot processes.

Provisioner (bootenvs)
----------------------

The Provisioner is a combination of several services and a template expansion engine.  The primary model is a boot environment (BootEnv) that contains crtical metadata to describe an installation process.  This metadata includes templates that are dynamically expanded when machines boot.

DigitalRebar Provision CLI has a process that combines multiple calls to install BootEnvs.  The following steps will configure a system capable to :ref:`rs_provion_discovered`.

  ::

    ../rscli bootenvs install bootenvs/sledgehammer.yml 
    ../rscli bootenvs install bootenvs/discovery.yml 
    ../rscli bootenvs install bootenvs/local.yml 
    ../rscli templates upload templates/local-elilo.tmpl as local-elilo.tmpl
    ../rscli templates upload templates/local-pxelinux.tmpl as local-pxelinux.tmpl
    ../rscli templates upload templates/local-ipxe.tmpl as local-ipxe.tmpl
    ../rscli prefs set unknownBootEnv to "discovery"
 
