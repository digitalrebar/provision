.. Copyright (c) 2020 RackN Inc.
.. Licensed under the Apache License, Version 2.0 (the "License");
.. Digital Rebar Provision documentation under Digital Rebar master license
.. index::
  pair: Digital Rebar Provision; Release v4.6
  pair: Digital Rebar Provision; Release Notes


.. _rs_release_v46:

Digital Rebar version 4.6 [planned]
-----------------------------------

Release Date: expected Winter 2020

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


.. _rs_release_v46_removals:

Removals
++++++++

None at this time.


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

Universal Workflow
~~~~~~~~~~~~~~~~~~

Planned, not committed


Restricted Access ISOs
~~~~~~~~~~~~~~~~~~~~~~

Planned, not committed


Integrated Simple DNS
~~~~~~~~~~~~~~~~~~~~~

Planned, not committed


.. _rs_release_v46_otheritems:

Other Items of Note
~~~~~~~~~~~~~~~~~~~

* TBD
