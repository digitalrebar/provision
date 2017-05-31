.. Copyright (c) 2017 RackN Inc.
.. Licensed under the Apache License, Version 2.0 (the "License");
.. Digital Rebar Provision documentation under Digital Rebar master license
.. index::
  pair: Digital Rebar Provision; configuring

.. _rs_configuring:

Configuring the Server
~~~~~~~~~~~~~~~~~~~~~~

Digital Rebar Provision provides both DHCP and Provisioning services but can be run with either disabled.  This allows users to work in environments with existing DHCP infrastructure or to use Digital Rebar Provision as an API driven DHCP server.

DHCP Server (subnets)
---------------------

The DHCP server is configured to enable Subnets that serve IPs and/or additional configuration information.  It is possible to run the DHCP server using only predefined IP Reservations or allow the DHCP server to create IP Leases dynamically.

The DHCP server had two primary models

#. DHCP Listeners can be set on an IP for each server interface.  These listeners will respond to DHCP broadcasts on the matching network(s).  Operators should ensure that no other DHCP servers are set up on the configured subnets.

#. DHCP Relay allows other DHCP listeners to forward requests to the Digital Rebar Provision server.  In this mode, the server is passive and can easily co-exist with other DHCP servers.  This mode works with the Provisioner by setting the many optional parameters (like next boot) that are needed for PXE boot processes.

Provisioner (bootenvs)
----------------------

The Provisioner is a combination of several services and a template expansion engine.  The primary model is a boot environment (BootEnv) that contains critical metadata to describe an installation process.  This metadata includes templates that are dynamically expanded when machines boot.

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

These settings apply to TEMPLATES only not the API.

The default password for the default o/s templates is **RocketSkates**

The default user for the default ubuntu/debian templates is **rocketskates**
