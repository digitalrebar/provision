.. Copyright (c) 2020 RackN Inc.
.. Licensed under the Apache License, Version 2.0 (the "License");
.. Digital Rebar Provision documentation under Digital Rebar master license
.. index::
  pair: Digital Rebar Provision; Setup

.. _rs_setup_preparing:

Preparing to Run Digital Rebar
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

Overview
--------

This document describes several considerations and preparation for a successful
setup and test (Proof-of-Concept, PoC, Trial, Test, etc.) of the Digital Rebar
Platform (DRP).  This document applies to both test environment setups, and preparing
for a production installation.


Sizing Your Endpoint
--------------------

Generally speaking, DRP by itself does not require many resources from either
your physical or virtual machines.  However, there are several considerations
primarily around the following, roughly in order of importance:

  * CPU
  * Memory
  * I/O (disk and network)
  * Network Throughput
  * Disk Space

For basic setup and testing of small scale environments (below 100 phyiscal or
virtual machines), the following *minimum* resources should be sufficient (more is
always better):

  * 2x (v)CPU at 2.0 GHz or better
  * 4 GByte (v)Memory
  * 100 GByte - this may vary dramatically if provisiong Windows Images
  * Prefer SSD class storage
  * 1 gbps (v)NIC

For a more complete discussion on Endpoint sizing within production environments,
please refer to the :ref:`rs_scaling`.

Please realize that provisiong can have extremely complex performance differences
depending on the type of provisioning (kickstart, preseed, image deployments), size
of artifacts, concurrency or parallel builds, and other external infrastructure
that your Workflows integrate with.  If you have any questions or would like
guidance, please `contact RackN <https://rackn.com/company/contact-us/>`_ for
help.


Supported Operating Systems
---------------------------

DRP is a single binary Golang application, which can be compiled and run on a
number of different operating systems and hardware platform architectures.  However,
there are several environmental considerations and constraints that the service
relies on.  As a consequcne, the following is the Operating System support
information:


Linux Distributions
===================

The following notes apply to Linux.

  * Typically any modern Linux distro should work
  * Only 64 bit OS builds are currently compiled, tested, and supported
  * Most testing is performed on CentOS 7, CentOS 8, RHEL 8, Ubuntu 18.04
  * Intel/AMD (x86_64, amd64) and ARM architectures are supported
  * RaspberryPI (rPI) is supported with custom named ``rpi`` builds

Most of RackN's commercial customers use either CentOS/RHEL, Ubuntu, or
a Docker container.


MacOS X (Darwin)
================

DRP runs and is regularly tested on MacOS X.  However, RackN does not
recommend nor support DRP running in production environments on MacOS X.


Container Environments
======================

DRP is released on Docker Hub as new releases are made.  This allows DRP to run
on any supported environment that runs Docker Containers.  Other container
environments may also work, but are not actively tested by RackN.

.. note:: For container environments, you will need to ensure you map low numbered ports in
          to the container correctly, and provide elevated privileges for the container
          that DRP is running in.

The RackN release docker container supports storing all writable data in a separate
storage volume that is attached to the container.  The ``install.sh`` script provides
installation, upgrade, and removeal support for containers; incluiding the storage
volume.

Please see the :ref:`rs_install` documentation for more complete details.


Miscellaneous Notes
===================

DRP may run on other alternative platforms that are loosely similar to Linux.  For
example, the BSD variants.  However, RackN does not test nor support these distributions
today.


Windows
=======

RackN does not support, nor does the ``dr-provision`` (DRP Endpoint) service run on Windows.
The Command Line tool (``drpcli``) does run, and is supported fully on Windows 64 bit operating
system versions.


Network Connectivity
--------------------

DHCP Broadcast Traffic
======================

DRP provides PXE boot provisioning and DHCP IP addressing for systems.  As such, your DRP Endpoint
will need network connectivity via the locally connected network interfaces, and routing to any
of the network subnets that you wish to provide Provisioning services for.  As long as standard
IP connectivity and reachability can be accomplished, DRP should be able to successfully provision
systems.

If DRP is the IP Address Management (IPAM) service for your systems via the built-in DHCP server, you
must also ensure *IP Helpers* or *DHCP Relay* options on your network switches/routers are correctly
set up to forward Broadcast DHCP traffic to the DRP Endpoint IP address.  This is a standard requirement
for any DHCP and PXE provisioning system, and is not unique to DRP.

Some additional resources related to forwarding DHCP traffic appropriately:

  * `Advanced IP Address Management <https://www.ciscopress.com/articles/article.asp?p=330807&seqNum=9>`_
  * `DHCP Relay Agent in Computer Network <https://www.geeksforgeeks.org/dhcp-relay-agent-in-computer-network/>`_

IP Helper / DHCP Relay issues must be addressed both in Physical environments and Virtual Environments.
Ultimately, your switch (be it virtual or physical) must forward Broadcast DHCP traffic to the DRP Endpoint.


Interface Speed and Duplex
==========================

DRP's communication to support Workflow does not require much bandwith.  Generally, these control messages
are small packets, and do not consume much network resources.  However, the act of provisioning systems
will require bandwidth dependent on the number of systems being provisioned in parallel, and the size
of the provisioning artifacts (packages, images, etc).  If these resources are hosted on the DRP Endpoints
web service (the standard configuration), then you will need to consider these activities in your sizing.

For small lab and test environments, a single 1 gbps full duplex network link is likely sufficient.  For
larger production environments, we recommend 2x 10 gbps bonded links for both bandwidth, and reliability.

Please refer to the :ref:`rs_scaling` for additional information.


Provisioning -vs- Baseboard Management Network
==============================================

The Baseboard Management Network (BMC), often times referred to inaccurately as the "*IPMI*" network, provides
*out-of-band* control path and management of the physical machines in your environment.  Typically these
networks are isolated from production network traffic.  DRP can be configured to interact with, and control
physical machine hardware via the BMC.  The only requirement that DRP has, it network reachability to the IP
addresses of the BMC systems themselves.

To accompllish this, the DRP Endpoint can be "multi-homed", or connected via NICs/Network Segments to both
the in-band provisioning network, and the out-of-band provisioning network.  Alternatively, a single network
interface will suffice, if that network routes to the BMC interfaces.  In this case, it would be prudent to
insure that Firewall or ACL rules block access to unknown systems.  If firewall/ACLs are in place, ensure
that your DRP Endpoint(s) IP addresses are whitelisted/allowewd access.


Multiple Network Connections
============================

It is not necessary to add multiple network interfaces on the DRP endpoint to each network, assuming that
the network switches and/or routers are appropriately forwarding DHCP Broadcast traffic to the DRP Endpoint.

If you do have multiple network connections on the DRP endpoint, it is absolutely critical that you evaluate
your networks Layer 3 (routed) topology, and ensure there are not asymmetric routing issues.  If a packet
ingresses one interface, but the DRP Endpoints operating system routing rules forward replies out a different
interface, this will almost always break provisioning services.

You may need to install/use IP Policy Based Routing (PBR) rules on the DRP Endpoint host operating system
to insure inbound/outbound routing of traffic conforms correctly to the required network paths.  For additional
information, and resources, see:

  * `Policy Based Routing: Concepts and Linux Implementation <https://silo.tips/download/advanced-routing-scenarios-policy-based-routing-concepts-and-linux-implementatio>`_


DHCP Services / IP Addressing
-----------------------------

The DRP Endpoint services is also a full fledges DHCP server.  DRP works extremely hard to provide clean,
fast, and accurate DHCP services that are tightly integrated with the Provisioniong process.  The service
is also designed to reduce as much complexity as possible in the setup, operation, and runtime of the
DHCP services.  Generally speaking, we encourage and recommend that customers use the DRP based DCHP
services whenever possible.

However, in environments with existing (legacy) based DHCP services, or with very complex network topologies
or hardware, it may not be possible to use DRP's DHCP services.  DRP provisioning does support use of
external DHCP services.  The basic mechanisms of "``nextserver``" and "``bootfile``" configurations must
be setup correctly in the external DHCP server.  That configuration is generally extremely specific to
the hardware and the DHCP server implementation.  Please consult your documentation on how to forward the
DHCP PXE traffic appropriately.

DRP controls it's internal DHCP services via the definition of *Subnets* which define the start and ending
ranges for IP Address handout during the DHCP negotiations (often times referred to as *DORA*).  The DHCP
server will NOT interfer with other DHCP traffic, as long as a Subnet for that Layer 3 network is not
configured.

Ultimately - there can only be ONE authoritative DHCP server of record for a Layer 3 Subnet/Network.  You
must ensure that there are no other competing DHCP servers or services on the network, otherwise provisioning
activities will likely fail.

DRP does support providing *Proxy* DHCP responses for limited DHCP servers that do not understand how to
provide the PXE required ``nextserver`` and ``bootfile``.


Provisioning Targets (Machines) Requirements
--------------------------------------------

Physical Machine Requirements
=============================

The vast majority of physical hardware that is provisioned and managed with DRP are server class
systems, with a Baseboard Management Controller (BMC, iDRAC, iLO, XCC, etc.).  Typically, you
will want to set these systems in BIOS to boot PXE first on the primary NIC that you designate
as your in-band provisioning network.  Systems with a BMC are not required, as long as they can
be set to PXE boot, in those cases, DRP can take control via the PXE boot path, and in-band
reboots of the system.  External power control will have to be used/implemented outside of DRP
in these cases (contact RackN, there are tools to control other non-standard Power Management
systems).

DRP is capable of managing switches (via ONIE/ZTP boot/install control), and storage devices.  For
these devices, please `contact RackN <https://rackn.com/company/contact-us/>`_ for further details.


Using Virtual Machines
======================

DRP will provision Virtual Machines equally well as physical hardware.  Similar to physical hardware,
the Virtual Machine vBIOS boot order needs to be configured with PXE boot first.

If you are utilizing virtual machines, you are generally free to size your VMs to whatever virtual
hardware sizes you need.  However, note that some Operating Systems that you might provision will have
requirements that may dictate the lower bounds of your VMs sizing configuration.

For example, CentOS/Redhat VMs should be configured with at least 2 GB of vMemory.  This is a requirement
of the CentOS installation and *dracut* tooling.  In general, we would recommend not configuring vMemory
on your Virtual Machines with less than 2 GB as a safety precaution.


