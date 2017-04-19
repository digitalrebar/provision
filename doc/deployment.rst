.. Copyright (c) 2017 RackN Inc.
.. Licensed under the Apache License, Version 2.0 (the "License");
.. DigitalRebar Provision documentation under Digital Rebar master license
.. index::
  pair: DigitalRebar Provision; Deployments

.. _rs_deployment:


Deployments
~~~~~~~~~~~

DigitalRebar Provision is intended to be deployed as both a DHCP server and a Provisioner.  There are cases where
you may want one or the other only.  Each feature can be disabled by command line flags.

* *--disable-dhcp* - Turns off the DHCP server
* *--disable-provisioner* - Turns off the Provisioner servers (TFTP and HTTP)

The :ref:`rs_api` doesn't change based upon these flags, only the services being provided.


DHCP Disabled
-------------

If you have your own DHCP environment or wish to run in a more declarative mode, you need to be aware of a couple of
things for each case.  For either case, the underlying assumption is that something will direct the node to use
the provisioner as its NextBoot server.

Declarative Mode
================

This is possible and works.  You will need to declare each machine through the :ref:`rs_cli` or the :ref:`rs_api`.
The IP address in the :ref:`rs_model_machine` will be used to populate the :ref:`rs_model_bootenv` information.  The
provisioner will provide these files through TFTP and HTTP.  It is left to the admin to figure out how to get the
node to reference them.


External DHCP
=============

With DHCP disabled, the admin can provide a DHCP server for handing out addresses.  The DHCP will need to do
the following:

* Set NextServer to an IP that routes to DigitalRebar Provision
* Set Option 67 (bootfile) to *lpxelinux.0* for legacy bios boots.  Other options are available for other.  See: :ref:`rs_model_subnet`
* Set Option 15 (dns domain) - This is needed for discovery boots to construct a meaningful FQDN for the node.
* Set Option 6 (dns server) - This is optional, but often useful in conjunction with Option 15.
* Set Option 3 (gateway) - This is optional, but may be required depending on the routing in your network.


Provisioner Disabled
--------------------

In this mode, DigitalRebar Provision acts as a DHCP server only.  The :ref:`rs_dhcp_models` describe how to use the server. 
Set the DHCP options that will direct to your next boot servers and other needs.

