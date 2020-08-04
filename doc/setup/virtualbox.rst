.. Copyright (c) 2020 RackN Inc.
.. Licensed under the Apache License, Version 2.0 (the "License");
.. Digital Rebar Provision documentation under Digital Rebar master license
.. index::
  pair: Digital Rebar Provision; Setup

.. _rs_setup_virtualbox:

VirtualBox Setup Instructions
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

Overview
--------

This document will help define one possible method to setup a VirtualBox
environment to test Digital Rebar Provision with.   There are several valid
ways to set up VirtualBox to accomplish testing Digital Rebar Provision
and we are only suggesting one method that works for us. 

Please feel free to provide any feedback or alternative setup methods.

Architecture
------------

For this document, we will be running DRP in isolated mode (started manually) inside the base OS.  VirtualBox will be
used to provide Virtual Machines that will be able to PXE boot and provide a platform for testing DRP.  The
Virtual Machines will be generally considered isolated by default (no internet access).   The virtual machine can access
the internet by configuring a secondary interface to NAT to the internet. Additionally, this setup can be used to manage
externally attached systems as well.

The system will be running DRP within a shell on the HOST OS.  VirtualBox will provide the virtual machines running
on a HOST-ONLY network.  The HOST-ONLY network will the default operational network for the virtual machine.  An additional
NAT network will used to provide conditional internet access for the virtual machine.

The examples for this use case will also be done from Mac OS running Catalina 10.15.

Setup
-----

This will take you through the components to bring up this environment.

Install VirtualBox
==================

You will need to install VirtualBox.  Start `here <https://www.virtualbox.org/wiki/Downloads>`_ to get VirtualBox for your system.
Follow the installation procedure for your system.  For this, VirtualBox is at version 6.1.12.

Once VirtualBox is installed, you will need to create a host-only network.  To do this, start virtual box and look at the
GUI.  From the TOOLS tab, select network.  Select the *Create* network button.  This will add row into the network table.
This row will be named something like, vboxnet0.  This will have an IPv4 address range and DHCP Server enabled.  You
will need to turn *OFF* the DHCP server by deselecting the "Enable Server" field on the *DHCP Server*.  Additionally, you
will need select that the network be manually configured and choose the first address in the range.  For example, the
network might be 192.168.100.1/24 from the network row.  The IPv4 Address for the HOST should be 192.168.100.1.  This will
be what DRP uses to communicate with the Virtual Machines.

This will NOT create a network until you create your first virtual machine.  More on that later.

Install Digital Rebar Provision
===============================

Using the :ref:`rs_qs_install`, we will use the isolated install method for getting DRP.  You will want to do this from
a shell.  Please follow through getting bootenvs and preparing to add the subnet.  At this point, you should hold off
adding the Subnet until a virtual machine is started.

Create an Initial Virtual Machine
=================================

To complete network setup, the process can be done easily if we create a virtual machine first.

A virtual machine needs a minimum of 20G of DISK and 3GB of memory.  If you are planning on testing Kubernetes or other
bigger platforms, you will want more memory and disk space.  3GB is a minimum for speedy sledgehammer operations.

To create a virtual machine:

* Open the GUI
* Select the New button from the Welcome Screen.
* A new virtual machine will pop up
* Name the virtual machine, test1, or whatever you want.
* This will select a machine folder (you can take the default).
* Select the OS type to *Linux* and Version to the general 64-bit value, *Linux 2.6/3.x/4.x (64-bit)*
* Move the memory slider to 3GB or type 3072.
* Take the default disk selector
* Click Create
* This will pop up a disk size and filename selector.  You can take the defaults, but you may want to increase size to 20GB for small linux installs.  60GB is usually needed for ESXi and Windows installs.
* Click Create

At this point, you have a virtual machine.  You will need to change two more things, boot order and nic attachments.
Open the settings for this virtual machine and select the *System* tab.  Change the BootOrder section to have Network checked
and first.  Hard Disk should be second and checked.  The rest should be unchecked.  This will make the machine default to
PXE boots and the Hard Disk boots.

The second change is to update the NIC settings.  Select the Network tab.  Adapter 1 will be selected and enabled, but on the NAT
network.  Change the NAT value to Host-only Adapter.  Make sure your Host-Only network is selected, e.g. vboxnet0. Select
the Adapter 2 tab and enable the adapter and select NAT.

Close the settings view and start the machine.  This will create the vboxnet0 network if not already present.  You can
then stop the virtual machine for the moment.

NOTE: By default, it will use 1 Virtual CPU.  This is fine for basic testing, but you may want to change that to 4.

Configure Subnet
================

At this point, open the UX and select the Subnet section of the Nav tree.  Click the Add button in the top of the main
panel.  This will pop up a selector interface names.  Select the vboxnet0 entry to create the subnet.  You can take the
defaults.

At this point, you will need to an route for Mac OS systems to function properly.

  ::

    sudo route add 255.255.255.255 192.168.100.1

In this example, the IPv4 Address, 192.168.100.1, is the address supplied in the network configuration section for the
HOST address.

Boot Virtual Machine
====================

If you have followed the quickstart and downloaded the sledgehammer ISO and set the default bootenvs, you should be
able to start the VM and have it boot into sledgehammer.  The Machine should register into the UX.

Common problems:

* VM fails to get DHCP - this usually happens because the route needed for broadcast is missing.  Use netstat -rn to check for the route and re-add it.
* Sledgehammer crashes - Make sure you have at least 3GB of memory for the virtual machine.

At this point, you can create additional machines as needed.  You can also install other OS to the system following the
other documentation sections.

Additional Setups
-----------------

Some additional things that can be done with this setup.

Building Sledgehammer
=====================

To build sledgehammer in this environment, you will need to follow the sledgehammer-builder content pack documentation.
To access the internet, you will need to add the following parameter to enable the NAT network.

  ::

    sledgehammer/extra-ifs = [ "enp0s8" ]

This will enable the NAT interface during the sledgehammer builder.

Managing External Systems
=========================

Using this same DRP, external systems can be managed by attaching the system to an external network.  For the MAC, using
an external ethernet port, create a new network and configure the network parameters on the system.  You can then create
a subnet for that network in DRP as well.  A route is not needed for that interface.  PXE boot the external systems and they
should join as well.  This will use the DRP system as the DHCP server for that network.


