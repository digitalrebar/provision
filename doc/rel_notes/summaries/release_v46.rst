.. Copyright (c) 2020 RackN Inc.
.. Licensed under the Apache License, Version 2.0 (the "License");
.. Digital Rebar Provision documentation under Digital Rebar master license
.. index::
  pair: Digital Rebar Provision; Release v4.6
  pair: Digital Rebar Provision; Release Notes


.. _rs_release_v46:

Digital Rebar version 4.6
-------------------------

Release Date: expected early 2021

Release Themes: Usability, Universal Workflow

In addition to bug fixes and performance improvements, the release includes several customer-driven features.

See :ref:`rs_release_summaries` for a complete list of all releases.

.. _rs_release_v46_notices:

Important Notices
~~~~~~~~~~~~~~~~~

.. _rs_release_v46_vulns:

Vulnerabilities
+++++++++++++++

The following vulnerabilities were reported

* :ref:`rs_cve_20200924a`

.. _rs_release_v46_deprecations:

Deprecations
++++++++++++

The following items are flagged as deprecated in v4.5 and will be removed in v4.6.

* old pattern cluster synchronization with cluster-add, cluster-step and cluster-sync.  Operators should migrate to the new cluster-gate-* patterns.
* terraform-provider-drp based on the DRP v3 API will not be supported
* Centos7 Sledgehammer will not continue to get updates after this release.  Customers should plan to migrate to the Centos8 version.


.. _rs_release_v46_removals:

Removals
++++++++

None at this time.


Sledgehammer Cento8
~~~~~~~~~~~~~~~~~~~

We have made many updates to the Cento8 Sledgehammer this release.  You will need to get an updated version of Sledgehammer as part of this release.  Centos7 Sledgehammer
is still available as a profile that can be added to the global profile or a specific machine to allow for its usage.

Limit client TLS ciphers
~~~~~~~~~~~~~~~~~~~~~~~~

DRP server allows clients to connect with a range of TLS ciphers by default.  Some security teams choose restrict the allowed ciphers.

Operators who wish to restrict use of client ciphers are advised to start with the `--tls-min-version` flags.  Operators can use the `--tls-cipher-list` and `--tls-ciphers-available` command line flags to determine the current and available ciphers.

Note: this feature was also backported to v4.5.2+


Stand Alone High Availability (HA with RAFT)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

The v4.3 :ref:`rs_release_v43_ha` feature was designed to rely on pacemaker or corosync to trigger failover.  This requirement has been elimimated.

In this verion, DRP HA includes integrated leader (re)election capability.  Operators will be able to influence or force changes in leadership.

This enables DRP to be used in site bootstrapping environments or locations where the additional requirement for failover detection was an operational burden.

Support for Client Long Polling
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

Long Polling allows clients that cannot use websockets to monitor DRP server for data changes to select objects.  This option provides a much lower overhead and faster way for clients to monitor DRP for updates than time-based polling.

RackN recommends using websockets when available; however long polling is strongly encouraged to reduce load on the DRP server when websockets cannot be used.

Image Deploy for ESXi
~~~~~~~~~~~~~~~~~~~~~

To improve the speed and consistence of VMware ESXi installation, The Image Deploy workflow has been expanded to include support for the ESXi operating system.  This allows operators to install ESXi directly to disk from a proven image and bypass the time consuming netboot (WEASEL) and post-configuration processes.

The process has specific requirements including the Digital Rebar VMware agent (aka DRPY) and having the correct partition maps.  Please contact RackN for assistance.

UX Improved Performance
~~~~~~~~~~~~~~~~~~~~~~~

The object storage, retrieval and event processing of the UX was significantly refactored to improve performance for larger environments.  In the new model, static objects are cached by the user's browser during initial login and do not have to be (re)retrieved on each page update.  In addition, event subscriptions are limited to the displayed objects only.

Previous versions of the UX subscribed to all system events.  This created a significant load on both browser and DRP server in large scale environments.


UX Improved Task Debugger
~~~~~~~~~~~~~~~~~~~~~~~~~

To better assist Task developers, the UX added a new tab to the machine view that shows live updates of the Jobs running against a machine during a workflow execution.  This helps developeres monitor a whole workflow lifecycle and provides additional tools for debug, stepping and restarting.

Note: This change relies on features that are only available in v4.6.

We expect this view will continue to improve as the communtiy provides feedback.


Universal Workflow
~~~~~~~~~~~~~~~~~~

The components of Universal Workflow are all included the the v4.6 release.  Universal Workflow provides a standardized workflow that can be applied to all provisioning operations in a consistent way.  Within that workflow, the system is able to dynamically adapt to the detected environment and take additional actions.  Unlike previous cloned Workflows, Operators may add their own custom stages and tasks to the Universal Workflow without interfering with standard operating processes.

Note: There are no helpers or added design tools for Universal Workflow in the v4.6 UX.  These will roll out incrementally based on customer design interactions.


.. _rs_release_v46_otheritems:

Other Items of Note
~~~~~~~~~~~~~~~~~~~

* UX
  * Improved alerting if DRP server loses connection
  * Token cached to avoid login if browser is refreshed
  * Machine Debug View (requires v4.6 DRP)
  * Catalog Intelligent Upgrade
  * Catalog Limits Versions
  * Ability to set UX Banner color
