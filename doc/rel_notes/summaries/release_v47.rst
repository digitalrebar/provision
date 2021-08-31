.. Copyright (c) 2020 RackN Inc.
.. Licensed under the Apache License, Version 2.0 (the "License");
.. Digital Rebar Provision documentation under Digital Rebar master license
.. index::
  pair: Digital Rebar Provision; Release v4.7
  pair: Digital Rebar Provision; Release Notes


.. _rs_release_v47:

Digital Rebar version 4.7
-------------------------

Release Date: Sept 1, 2021

Release Themes: 

In addition to bug fixes and performance improvements, the release includes several customer-driven features.

See :ref:`rs_release_summaries` for a complete list of all releases.

.. _rs_release_v47_notices:

Important Notices
~~~~~~~~~~~~~~~~~

* Digital Rebar v4.7 adds port 8090 to the list of ports _required_ for provisioning operations. Please verify that port 8090 (default, this can be changed) is accessible for Digital Rebar endpoints.
* Due to changes in the install zip format, the API-based upgrade of DRP to v4.7+ requires usage of most recent https://portal.RackN.io (v4.7 for self-hosted UX users) or the use of DRPCLI v4.6.7+. The v4.7 ``install.sh upgrade`` process also includes theses changes.

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

Table Refactor for RackN Portal UX
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

The v4.7 release includes significant improvements to the table and panel displays for the UX.  significant effort was made to ensure minimal relearning effort for existing operators.

This update provides long requested features including:

* resizable columns for all views
* stable column sizing (will not change size as data is changed)
* panels can be resized
* panels remain open and can be changed by selecting different rows
* improved filter usability (moves to top of page)
* improved performance and scalability for >10,000 machine customers
* dramatically reduced Digital Rebar API load based on better use of local cache

Advanced operators are asked to use this version for testing and feedback.

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

.. _rs_release_v47_splitapi:

Split Static Files & Template Renders for Public Endpoints (Port 8090)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

Digital Rebar v4.7 adds enhanced port based access security by splitting the secure dynamic
file server (port 8090 default) from API server (port 8092 by default).

This change facilitates using public internet facing Digital Rebar endpoints by limiting the
potential exposure of sensitive data in unauthenticated machine provisioning templates.
For public facing endpoints in which users need API access, operators are encouraged to
block access to port 8090 for untrusted access.

Since the Digital Rebar new/discovered machine join process requires access to dynamically
generated templates; provisioning operations _require_ access to port 8090.

To provide backwards compatibility, v4.7 automatically forwards requests for generated files
from port 8092 to 8090.  If this port is closed for security, those requests will be blocked.

See :ref:`rs_arch_ports` for more networking details.

.. _rs_release_v47_bootp:

BOOTP Support
~~~~~~~~~~~~~

Before there was PXE, there was BOOTP for provisioning!  Digital Rebar now supports BOOTP;
however this feature requires use of Reservations.


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

* Added Cloud-Init to fingerprint information
* Improves to NetWrangler to help configure Linux networking
* Ability to filter deeply into Params and Meta data
* Improve consistency of CLI flag formats
* Allow new Multi-Site managers to connect to older Digital Rebar endpoints
* Fix PATCH against Params


