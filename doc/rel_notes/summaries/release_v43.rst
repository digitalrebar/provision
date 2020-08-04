.. Copyright (c) 2020 RackN Inc.
.. Licensed under the Apache License, Version 2.0 (the "License");
.. Digital Rebar Provision documentation under Digital Rebar master license
.. index::
  pair: Digital Rebar Provision; Release v4.3
  pair: Digital Rebar Provision; Release Notes


.. _rs_release_v43:

Digital Rebar version 4.3
-------------------------

Release Date: June 3, 2020

Release Themes: Secure Boot (preview), High Availability and Back End Performance

Digital Rebar v4.3 reflects a significant enhancement of backend storage work delivered in v4.0 for a single endpoint.  The long pole for this release has been validation, stress testing and usability enhancements needed to allow other endpoints to subscribe directly to the backend transaction streams.  These changes enabled high availability and enabled bringing the Multi-Site Manager components into the platform core.

Along with improvements to the core platform storage, RackN included numerous usability, bug fixes, and content extensions.  Some of the notable ones include enhanced VMware support and secure boot.  We’ve adapted Digital Rebar to an RPis in a subproject called EdgeLab.digital and expanded the UX to include messaging to RackN outside of Slack.

In addition to bug fixes and performance improvements, the release includes several customer-driven features.

See :ref:`rs_release_summaries` for a complete list of all releases.

.. _rs_release_v43_ha:

High Availability
~~~~~~~~~~~~~~~~~

Digital Rebar v4.3 adds basic support for high-availability configurations.  In the initial release, that support is limited to an active-passive configuration (with zero or more passive endpoints) sharing a single virtual IP address with manual failover between the endpoints.  The intent is to provide the basic features needed to allow automated failover using standard Linux cluster tools such as pacemaker and corosync.  The active and passive endpoints remain in sync using a combination of write-ahead log shipping and replication of important data via synchronous replication from the active to the passive endpoints.


.. _rs_release_v43_ux_inbox:

UX Messaging Feature
~~~~~~~~~~~~~~~~~~~~

The RackN UX Inbox allows operators to communicate directly with RackN via the UX.  This facilitates communication outside of the RackN Slack or email.  It also allows RackN to push notifications about releases, patches and other operational issues directly to users.

Note: UX Messaging is not an official support channel.  Commercial customers should raise Zendesk tickets for urgent production support or tracked issues.


.. _rs_release_v43_secure_boot:

Secure Boot
~~~~~~~~~~~

Digital Rebar adds support for discovering and booting systems that have UEFI Secure Boot enabled.  This support relies on a combination of DHCP server enhancements and boot environment updates.  The Sledgehammer discovery image and the CentOS install boot environments will work with Secure Boot out of the box, and support is planned for VMWare, Windows, and the Debian/Ubuntu boot environments.  A license with proper entitlements is needed to take advantage of Secure Boot, and Digital Rebar must be the DHCP server for Secure Boot support to function.

.. _rs_release_v43_vmware:

VMware V7 Integration
~~~~~~~~~~~~~~~~~~~~~

Support for ESXi 7.0 has been added. Changes have been made that move the agent code from firstboot into a native esxi agent. The agent rework will allow us to support secure boot once the agent code has been reviewed and signed by VMware.  

Note: the underlying requirements for secure boot are available in v4.3; consequently, the VMware v7 secure boot integration will be available as a content pack update.

.. _rs_release_v43_performance:

Performance Enhancements
~~~~~~~~~~~~~~~~~~~~~~~~

Digital Rebar went through performance analysis and benchmarking to increase backend performance and DHCP handling performance under adverse conditions.  Key improvements include:

* Rewriting the backend storage layer to collect long-lived data into a single file instead of having thousands of smaller files.  This was also required to enable streaming replication for HA support.
* Updating the WAL handling code to use more efficient JSON encoding and decoding methods.
* DHCP lease handling code has been made multi-threaded.

.. _rs_release_v43_fingerprinting:

Machine Fingerprinting
~~~~~~~~~~~~~~~~~~~~~~~

Digital Rebar no longer relies exclusively on MAC address or system IP address to identify a machine when it PXE boots.  Instead, we use a combination of system and chassis serial numbers, the MAC addresses of the NICs present in the system, and the serial numbers of any memory DIMMs present in the system.  This allows Digital Rebar to function more reliably in the face of incremental hardware changes, and removes several ways the system could become confused in the face of IP address exhaustion when using an external DHCP server.


.. _rs_release_v43_multisite:

Delayed: Integrated Multi-Site Manager
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

After reviewing customer feedback and testing with the Multi-Site Manager Plugin, RackN accelerated plans to bring that feature set into the core platform.  This work is one of the themes for the v4.4 release and has been gated by the transaction streaming work completed for HA in this release.  RackN has deprecated the MSM plugin and will collaborate with early adopters of the feature during v4.4 development.

.. _rs_release_v43_otheritems:

Other Items of Note
~~~~~~~~~~~~~~~~~~~

* Integrated High Availability
* UX Messaging Feature
* Secure Boot
* VMware v7 Integration
* Other Notable Enhancements 
  * Support for additional DHCP Options
  * Significant performance enhancements on backend storage
  * Machine fingerprint supports constrained external DHCP
  * Re-integration of automated test coverage reports (sustained >70%)
* Integrations and Operational Helpers
  * Integrated log rotation settings with safe defaults
  * Improved Ansible Integrations via API and Contexts
  * Endpoint bootstrapping workflows (was beta in v4.2)
* Hardware Expansions
  * Raspberry Pi Support (exposed via EdgeLab.digital)
  * Netapp Solidfire and Cohesity Support (not in public catalog)


