.. Copyright (c) 2020 RackN Inc.
.. Licensed under the Apache License, Version 2.0 (the "License");
.. Digital Rebar Platform documentation under Digital Rebar master license
.. index::
  pair: Digital Rebar Platform; Universal Workflow Operations

.. _rs_universal_ops:

Universal Workflow Operations
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

This section will address usage of the Universal Workflow system.  The architecture and implementation of the Universal
Workflow system is described at :ref:`rs_universal_arch`.

The following workflows are the default set of universal workflows.  Using these are starting points that do most default
operations needed to operate a data center.

Each workflow defines what it requires on start and set on completion.  These actions allow for the workflows to work
together.  For example, Discover should set the needed parameters for driving Hardware and the networking in the other workflows.
These are the elements defined by the universal parts of these components.

Universal depends upon:

  * Digital Rebar Community Content Pack
  * Task Library Content Pack
  * FlexiFlow Content Pack
  * Validation Content Pack
  * Callback Plugin
  * Classification Content Pack

Some of the Universal Workflows are:

  * Discover - Handles discovery and inventory
  * Hardware - Handles hardware configuration
  * Build Baseline - Handles generating profiles to be consumed by Hardware
  * ESXI Install - Handles installing ESXI.  A map defines what to install
  * Solidfire RTFI - Handles updating and setting up SolidFire systems through RTFI process
  * Linux Install - Handles installing a linux operating system
  * Image Deploy - Handles deploying an image-based system
  * Maintenance - Sets the system to go into non-destructive Hardware and returnes toe LocalIdle
  * Rebuild - Redrives the system through the process
  * Decommission - Handles decommissioning hardware
  * SledgehammerWait - End state for waiting or debugging

Discover
========

This workflow handles discovery a machine and inventories the machine.  It also will validate the basic initial pieces.
The final piece is classification to drive additional workflows if needed.

The default stages for `discover` are:

  * discover
  * centos-setup-repos
  * ipmi-inventory
  * bios-inventory
  * raid-inventory
  * network-lldp
  * inventory

A common set of post `discover` flexiflow tasks would be

  * rack-discover
  * universal-classify
  * universal-host-crt-callback

Hardware
========

This workflow is used to configure hardware and burnin that state.  It is also the place where BMC certificates are loaded.

The workflow assumes that the parameters of configuration have been set.  See the power of `universal-classify` and `Build Baseline`.

The default stages for `hardware` are:

  * hardware-tools-install
  * ipmi-configure
  * flash
  * raid-enable-encryption
  * raid-configure
  * bios-configure
  * universal-ilo-config
  * burnin
  * burnin-reboot
  * universal-load-certs


Build Baseline
==============

This workflow is used outside of a universal workflow to snapshot the hardware state of a machine and builds profiles
named such they work with the universal-classify stage/task.

ESXI Install
============

Solidfire RTFI
==============

Linux Install
=============

Image Deploy
============

Image Capture
=============

This is not a universal workflow in itself.  By altering the linux-install workflow with the following tasks....

this goes in the linux post install pre reboot.

  - image-reset-package-repos
  - image-update-packages
  - image-install-cloud-init
  - image-builder-cleanup
  - image-capture

CloudInit Post Install
======================

Maintenance
===========

Rebuild
=======

Decommission
============

AutoMaintenance
===============

Local
=====

SledgehammerWait
================

