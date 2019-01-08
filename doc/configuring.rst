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

The DHCP server has three primary models

#. DHCP Listeners can be set on an IP for each server interface.  These listeners will respond to DHCP broadcasts on the matching network(s).  Operators should ensure that no other DHCP servers are set up on the configured subnets.

#. DHCP Relay allows other DHCP listeners to forward requests to the Digital Rebar Provision server.  In this mode, the server is passive and can easily co-exist with other DHCP servers.  This mode works with the Provisioner by setting the many optional parameters (like next boot) that are needed for PXE boot processes.

#. Proxy DHCP Mode - this allows for other DHCP servers to forward requests to the DRP Endpoint, and fillin necessary PXE related options that may be missing from the original DHCP request options.  In some environments, very limited DHCP servers may be in use (which do not support the appropriate Options), or a DRP Endpoint user may not have administrative authority over the initial DHCP server in use.   In this case, the :ref:`rs_model_subnet` should include the following configuration to enable Proxy DHCP mode:

  ::

    drpcli subnets update mysubnet '{ "Proxy": true }'


Provisioner (bootenvs)
----------------------

The Provisioner is a combination of several services and a template expansion engine.  The primary model is a boot environment (BootEnv) that contains critical metadata to describe an installation process.  This metadata includes templates that are dynamically expanded when machines boot.

Digital Rebar Provision CLI has a process that combines multiple calls to install BootEnvs.  The following steps will configure a system capable to :ref:`rs_provision_discovered`.

  ::

    drpcli bootenvs uploadiso sledgehammer
    drpcli prefs set defaultStage "discover" unknownBootEnv "discovery" defaultBootEnv "sledgehammer"

In addition to the basic *discovery* capability provided by the *sledgehammer* BootEnv, you will most likely want to install an Operating System related BootEnv.  Basic example for this is as follows:

  ::

    drpcli bootenvs uploadiso centos-7-install
    drpcli bootenvs uploadiso ubuntu-16.04-install

Each content pack has various supported operating system BootEnv definitions.  The default *drp-community-content* pack contains the following BootEnvs:

  #. centos-7-install: CentOS 7 (most recent released version)
  #. centos-7.4.1708-install: Centos 7.4.1708 (this may change as new versions are released)
  #. ubuntu-16.04-install: Ubuntu 16.04
  #. debian-8: Debian 8 (Jessie) version
  #. debian-9: Debian 9 (Stretch) version

Additional Operating System versions are available via registered RackN content pack add-ons.  

.. _rs_configuring_default:

Default Template Identity
-------------------------

These settings apply to TEMPLATES only not the API.

The default password for the default o/s templates is **RocketSkates**

The default user for the default ubuntu/debian templates is **rocketskates**
