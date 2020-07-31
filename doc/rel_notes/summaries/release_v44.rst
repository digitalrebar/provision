.. Copyright (c) 2020 RackN Inc.
.. Licensed under the Apache License, Version 2.0 (the "License");
.. Digital Rebar Provision documentation under Digital Rebar master license
.. index::
  pair: Digital Rebar Provision; Release v4.4
  pair: Digital Rebar Provision; Release Notes


.. _rs_release_v44:

Digital Rebar version 4.4
-------------------------

Release Date: July 30, 2020

Release Themes: Pooling, Secure Boot for VMware

Note about release cadence: RackN is taking steps to maintain a fast (2-3 month) release cadence with Digital Rebar.  The goal of smaller releases is to ensure that important features and fixes are not gated by multi-release feature deliverables.  In many cases, that allows advanced users to have technical preview access to new features inside of stable releases.  As always, we work to keep the primary development branch reliable throughout the development cycle.


In addition to bug fixes and performance improvements, the release includes several customer-driven features.

.. _rs_release_v44_pooling:

Infrastructure Resource Pooling
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

Infrastructure Resource Pooling provides APIs that create a “bare metal cloud” consumable by tools like Terraform. Pooling creates capabilities that abstract Machines into resource groups that can be allocated generically from a single API making the system more elastic.  This functionality enables cloud-like behavior because operators can request a Machine based on an attribute map rather before assigning a specific Machine(s).  The Pooling system also provides operator Workflows for allocation and reallocation processes. Benefits for Operators include controls and transparency to monitor while optimizing and managing end-user consumption in a scalable way. This balance reflects the RackN focus on self-management by helping infrastructure operators create SaaS like offerings. 

.. _rs_release_v44_secure_boot:

Trusted Platform Module (TPM) Secure Boot
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

Trusted Platform Module (TPM) Secure Boot ensures customers no longer have to choose between security or automation. Previously, most companies disregarded the risk of disabling the TPM because it is difficult to automate. This created security risks compromising operating systems especially in edge locations. Compliance issues emerge. RackN solves this flaw by enabling an intrinsically trusted, data center out of the box. Version 4.4 specifically, allows for fully automated processes from first boot to full cluster in under an hour with advanced features making infrastructure “secure by default”. A 100% compliance pass will now be standard. Additionally, out of band management passwords, tls certs, VLAN will all be secure by default with this new upgrade.

.. _rs_release_v44_vmware:

Standardized VMware ESXi Install
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

Building on :ref:`rs_release_v43_vmware`, improvements were made based on field testing and standardization efforts.  This dramatically simplified the overall ESXi installation process and takes advantage of the specialized DRPCLI agent (drpy) for ESXi.

.. _rs_release_v44_multisite:

Preview: Integrated Multi-Site Manager
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

Integrated Multi-Site Manager is in tech preview enabled for v4.4, released planned in v4.5.  This work was one of the major work items for the v4.4 release.  It is incorporated into the code base but disabled by default.  We are actively undergoing user and performance testing on this major feature.

.. _rs_release_v44_otheritems:

Other Items of Note
~~~~~~~~~~~~~~~~~~~

* Updates to Image Deploy - improved flow and configuraiton options for image based deployment and creation including updating Curtin.
* Bootenv options - Improved controls for BootEnv overrides
* Cloud Wrap - standardizes patterns to create cloud instances using Contexts with Terraform and Ansible.
* Prototype Terraform Provider - enables use of the Pooling API via Terraform
* UX Enhancements
  * Significant UX performance rendering improvements - helpful for large scale and systems that are under stress due to error storms
  * Select entire row instead of just checkbox
  * Improve claimed editing
  * Improve editing panels and pages
  * Content pack differencing
* Backup Enhancements - improved tooling around server backups and diagnostics
* Workflow timeout - allows operators to determine if a workflow has stalled
* Agent auto-update - Agent can upgrade itself.  This is very helpful for agents embedded in image deploys and containers.
* Additional Operating Systems
  * Debian 10
  * Centos 8 / RHEL 8
  * Ubuntu 20.04
* Moved into open ecosystem: Cohesity and Solidfire Support