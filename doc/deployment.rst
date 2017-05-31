.. Copyright (c) 2017 RackN Inc.
.. Licensed under the Apache License, Version 2.0 (the "License");
.. Digital Rebar Provision documentation under Digital Rebar master license
.. index::
  pair: Digital Rebar Provision; Deployments

.. _rs_deployment:


Deployment Options
~~~~~~~~~~~~~~~~~~

Digital Rebar Provision is intended to be deployed as both a DHCP server and a Provisioner.  There are cases where
one or the other are desired.  Each feature can be disabled by command line flags.

* *--disable-dhcp* - Turns off the DHCP server
* *--disable-provisioner* - Turns off the Provisioner servers (TFTP and HTTP)

The :ref:`rs_api` doesn't change based upon these flags, only the services being provided.


DHCP Disabled
-------------

If a DHCP environment already exists or a more declarative mode is more desirable, there are a couple of things in each case to be aware of.  
For either case, the underlying assumption is that something will direct the node to use
the provisioner as its *NextBoot* server.

Declarative Mode
================

Each machine must be declared through the :ref:`rs_cli` or the :ref:`rs_api`.
The IP address in the :ref:`rs_model_machine` will be used to populate the :ref:`rs_model_bootenv` information.  The
provisioner will provide these files through TFTP and HTTP.  It is left to the admin to figure out how to get the
node to reference them.


External DHCP
=============

With DHCP disabled, the admin can provide a DHCP server for distributing addresses.  The DHCP will need to do
the following:

* Set NextServer to an IP that routes to Digital Rebar Provision
* Set Option 67 (bootfile) to *lpxelinux.0* for legacy bios boots.  Other options are available for other.  See: :ref:`rs_model_subnet`
* Set Option 15 (dns domain) - This is needed for discovery boots to construct a meaningful FQDN for the node.
* Set Option 6 (dns server) - This is optional, but often useful in conjunction with Option 15.
* Set Option 3 (gateway) - This is optional, but may be required depending on the network routing.


Provisioner Disabled
--------------------

In this mode, Digital Rebar Provision acts as a DHCP server only.  The :ref:`rs_dhcp_models` describe how to use the server.
Set the DHCP options that will direct to the next boot servers and other needs.

