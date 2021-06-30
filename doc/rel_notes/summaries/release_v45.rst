.. Copyright (c) 2020 RackN Inc.
.. Licensed under the Apache License, Version 2.0 (the "License");
.. Digital Rebar Provision documentation under Digital Rebar master license
.. index::
  pair: Digital Rebar Provision; Release v4.5
  pair: Digital Rebar Provision; Release Notes


.. _rs_release_v45:

Digital Rebar version 4.5 [Fall 2020]
-------------------------------------

Release Date: September 25, 2020

Release Themes: Multi-Site Manager

In addition to bug fixes and performance improvements, the release includes several customer-driven features.

See :ref:`rs_release_summaries` for a complete list of all releases.

.. _rs_release_v45_deprecations:

Important Notices
~~~~~~~~~~~~~~~~~

Vulnerabilities
+++++++++++++++

The following vulnerabilities were reported and fixed this release.  See the CVE for which releases contain the fix.

* :ref:`rs_cve_20200924a`

Critical Known Issues
+++++++++++++++++++++

* The DRP binary includes a time-limited embedded license that will disable operations if not replaced by a RackN issued license.  In a new installation after the time-limit, DRP API v4.5 and earlier will not allow users to upload their license during installation before stopping.  To workaround this problem, create the full install path (e.g. ``/var/lib/dr-provision/saas-content/``), copy the license file there _before_ installing.,


Deprecations
++++++++++++

The following items are flagged as deprecated in v4.5 and will be removed in v4.6.

* old pattern cluster synchronization with cluster-add, cluster-step and cluster-sync.  Operators should migrate to the new cluster-gate-* patterns.
* terraform-provider-drp based on the DRP v3 API will not be supported


.. _rs_release_v45_removals:

Removals
++++++++

None at this time.


.. _rs_release_v45_multisite:

Integrated Multi-Site Manager
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

**This is a licensed feature**

This release integrated the Multi-Site Management features directly into DRP core instead of providing it as a plug-in.  While this change only has minor impacts on the API, it significantly improved the performance and resilience of the manager.

Multi-Site Manager significantly enhances Digital Rebar from providing a single site integrated infrastructure as code (IaC) provision platform into a distributed infrastructure control plane.  The design maintains site autonomy AND builds a management federation.  By design, each site may have multiple managers to allow for regional and global layering.

Major Features:

* Ability to consolidate data from multiple endpoints into subscribed endpoints
* Ability to change managed endpoint data via calls to the mirrored data in the manager API (uses DRP proxy forward)
* Ability to centrally manage and enforce configuration of remote endpoints via Version Sets
* Integrates with DRP HA features so that HA groups are managed as a single site.

Important note: all managed endpoints must be licensed BEFORE they can be attached to the manager.  When using the UX to add endpoints, it will automatically check and update the relevant license.

.. _rs_release_v45_sl8:

Sledgehammer on Centos 8 (SL8)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

To take advantage of the latest kernel and distro features, customers want Sledgehammer migrated to Centos 8 (aka SL8) instead of Centos 7 (aka SL7).  The change also includes python3.

As the heart of our discovery, imaging and hardware update features, Sledgehammer is a critical component for Digital Rebar operations.  Migrating to an updated Centos version allows operators to take advantage of features in the latest distro; however, the potential for Sledgehammer changes to impact existing workflows means that RackN takes a cautious approach to revisions.

To leverage the new SL8, add the enabling SL8 profile in your workflow.

Note: SL8 is expected to become default in DRP v4.6.  A profile will be provided the enable backwards SL7 support.


For v4.6 updates, :ref:`rs_release_v46_centos8`

.. _rs_release_v45_backup:

WAL Backup
~~~~~~~~~~

DRP v4 introduced the WAL transcation log data store for Digital Rebar data.  This critical core change has enabled endpoint replication features such as muiti-site and high availably.  In this release, RackN added additional tooling to help operators make backups of the WAL.

Operators can now use the `dr-backup` utility to capture DRP snapshots of live systems with confidence that they will detect and work with system transaction boundaries.  It can also be used remotely to ensure off-system and test convenience backups.

Operators are advised to migrate all backup operations to this new utility since it is more reliable than capturing the state of the file system.


.. _rs_release_v45_log_capture:

Log Capture 
~~~~~~~~~~~~

Utility allows operators to collect log information from DRP as a portal unit for analysis by RackN.

When working with customers and community environments, RackN often needs to review comprehensive system logs.  This utility makes it easier for operators to capture and package the correct logs.  This makes a greatly reduces the risk of incomplete captures and ultimately reduces the time to resolution for customers.

This features has been integrated into the v4.5 DRPCLI and documented there.

.. _rs_release_v45_performance:

Startup and API Performance Tuning
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

RackN customers are running systems with thousands of machines and high transaction loads.  With Multi-Site creating aggregate views of these systems, performance at scale is a critical aspect of the v4.5 release with Multi-Site manager.

Completed enhancements include:
  * Significant refactoring was performed to improve DRP start times and loading of content packs to running systems.
  * Stress testing of 1,000+ parallelized active agents was performed.
  * Optimizations and testing of the RackN UX for high object counts and activity levels.
  * Improved plugin initialization and safeties.

.. _rs_release_v45_terraform:

v4.5 Terraform Provider
~~~~~~~~~~~~~~~~~~~~~~~

The Terraform Provider (https://github.com/rackn/terraform-provider-drp) has been completely rewritten to work with Terraform v0.13+.  This new provider requires the v4.4 :ref:`rs_release_v44_pooling` feature.

Terraform is one of several systems that need to request and release Digital Rebar machines in a more abstracted way.  While the Terraform provider is valuable as a stand alone benefit for Terraform users, RackN also uses it to validate the pooling API process and interfaction.

Due to the new 3rd party registration feature, operators will be able to automatically download the updated provider from a RackN maintained registery.  This eliminates the requirement to track builds, download or create a local version of the provider.

Note: While the provider is APLv2 open source, this feature leverages the licensed feature of pre-defined pools.

.. _rs_release_v45_otheritems:

Other Items of Note
~~~~~~~~~~~~~~~~~~~

* `drpcli machines count` optimization bypassing sending data to get counts of machines
* Fixes to `docker-context` plugin to improve start-up and reset operations
* Tuning of the DHCP performance system
* Improved integration with VMware ESXi provisioning
* Significant updates and improvements to this documentation
* Expand ansible-local-playbooks task to use templates
* Updates to filebeat plugin
* Improved stability for self-runner bootstrapping agent
* Improved data collection and communication within HA clusters
* Web UX
   * Improved Params update from Machines List view including setting secure values
   * Numerous rendering and edit page fixes
