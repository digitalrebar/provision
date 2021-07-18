.. Copyright (c) 2020 RackN Inc.
.. Licensed under the Apache License, Version 2.0 (the "License");
.. Digital Rebar Provision documentation under Digital Rebar master license
.. index::
  pair: Digital Rebar Provision; Release v4.7
  pair: Digital Rebar Provision; Release Notes


.. _rs_release_v47:

Digital Rebar version 4.7 [in process]
--------------------------------------

Release Date: early August 2021

Release Themes: 

In addition to bug fixes and performance improvements, the release includes several customer-driven features.

See :ref:`rs_release_summaries` for a complete list of all releases.

.. _rs_release_v47_notices:

Important Notices
~~~~~~~~~~~~~~~~~

.. _rs_release_v47_vulns:

Vulnerabilities
+++++++++++++++

None known

.. _rs_release_v47_deprecations:

Deprecations
++++++++++++

None known

.. _rs_release_v47_removals:

Removals
++++++++

The following items are flagged as deprecated in v4.5 and are removed in v4.7.

* old pattern cluster synchronization with cluster-add, cluster-step and cluster-sync.  Operators should migrate to the new `cluster-manager` patterns.


.. _rs_release_v47_rocky:

Rocky and Alma Linux in Universal Workflow
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

With the transition to Centos Streams, we've added alternatives for two new Centos adjacent distros: Rocky and Alma linux.
For operators concerned about long term stability for Centos, these should offer a reliable open alternative.

These are available as applications and pre-defined profiles in Universal Workflows.

.. _rs_release_v47_ux_improvements:

Table Refactor for RackN Portal UX (tech preview)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

The v4.8 release will include significant improvements to the table and panel displays for the UX.  significant effort was made to ensure minimal relearning effort for existing operators.

This update provides long requested features including:

* resizable columns for all views
* stable column sizing (will not change size as data is changed)
* panels can be resized
* panels remain open and can be changed by selecting different rows
* improved filter usability (moves to top of page)
* improved performance and scalability for >10,000 machine customers
* dramatically reduced Digital Rebar API load based on better use of local cache

This work is available as a preview from https://tip.rackn.io.  Advanced operators are asked to use this version for testing and feedback.

As usual, the updated UX maintains compatability with all v4.x versions of Digital Rebar.

.. _rs_release_v47_hardware:

Firmware & Hardware Support
~~~~~~~~~~~~~~~~~~~~~~~~~~~

The following hardware OEMs and generalized tooling was added.

.. _rs_release_v47_firmware:

Generalized Firmware Update Process
-----------------------------------

Firmware support for SuperMicro (see below) was provided in a generalized way will be used for other OEMs.

At this time, no existing OEMs have been migrated to use this new process.

.. _rs_release_v47_lpar:

IBM LPAR (Logical Partitions)
-----------------------------

Allows Digital Rebar to manage LPAR VMs.  This brings control of the LPAR VM inline with physical infrastructure processes.

See: https://www.ibm.com/docs/en/zos-basic-skills?topic=design-mainframe-hardware-logical-partitions-lpars

.. _rs_release_v47_supermicro:

SuperMicro
-----------

Redfish BMC, firmware and RAID configuration of SuperMicro hardware.

.. _rs_release_v47_bootp:

BOOTP Support
~~~~~~~~~~~~~

Before there was PXE, there was BOOTP for provisioning!  Digital Rebar now supports BOOTP;
however this feature requires use of Reservations.

.. _rs_release_v47_terraform:

Terraform and Cloud-Wrapper Updates
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

Cloud Wrapper templates were updated to enable better integration into universal workflow
and simplify importing existing Terraform plans.

1. to use cloud-init instead of relying on Ansible for join-up.
2. to allow creating many machines from a single plan (uses cluster/profile)
3. improve controls after Terraform created instances
4. improve synchronization after Terraform destroys instances

.. _rs_release_v47_bootenv:

Implementation of Profile/Param Overrides for BootEnvs
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

Addition of the ``bootenv-customize`` parameter in v4.6 allowed operators to overlay dynamic customizations
on top of BootEnvs.  This feature was intended to reduce the number of BootEnvs maintained in the system.

Digital Rebar v4.7 included a render helper to make this process easier to apply.  The ESXi BootEnv has
been updated to use this feature.

Universal Workflow BootEnvs will leverage this feature.

.. _rs_release_v47_otheritems:

Other Items of Note
~~~~~~~~~~~~~~~~~~~

* TBD
