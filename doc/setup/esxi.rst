.. Copyright (c) 2020 RackN Inc.
.. Licensed under the Apache License, Version 2.0 (the "License");
.. Digital Rebar Provision documentation under Digital Rebar master license
.. index::
  pair: Digital Rebar Provision; Setup

.. _rs_setup_esxi:

ESXi Setup Instructions
=======================

Overview
--------
In this doc we cover how to use a couple of virtual machines hosted by ESXi for testing out drp.
This document will not cover how to install ESXi. For help with that please refer to it's official
documentation. While we specifically mention using ESXi in this article you could do the same from vCenter.
This document will only highlight one of the possible configurations that can be used for testing
deployments of VM's or physical machines. This document will mention using vlans, but will not be
a requirement unless you need to keep traffic separated.


Getting Started
---------------
I begin by logging into my switch to create a new VLAN. The VLAN I created for this has been given the ID
of 20. Next log into the ESXi web interface. Using the navigator proceed to the Networking. Select the
Port Group tab if its not already selected. Click on the Add port group button to begin the wizard. Name the
portgroup VLAN-20 (or what ever you like), add the port to your existing vSwitch (or make a new one for this
and add it to that vSwitch), and save it.


Set Up Explained
----------------
Following the getting started above should have you at a point of needing DRP and something to provision.
At this point you need a linux box to install DRP, and it can be a VM. Create a virtual machine in ESXi with 1+
CPU, and 1G+ mem 20G + storage (for sizing info: :ref:`rs_scaling`). Once you have your linux dostro installed
you are ready to follow our :ref:`rs_install` Now you need a "machine" to provision. This is where the VLAN part
is handy from above. You can now provision physical machines on this VALN, or you can add virtual machines to it
and go into their BIOS settings and adjust them to network boot using either EFI or BIOS, and both work with DRP.

.. note:: If you find that the boot order is being changed after you provision a VirtualMachine and it no longer trys to netboot by default you may have to manually edit its vmx file and add bios.bootOrder = "ethernet0,hd"
