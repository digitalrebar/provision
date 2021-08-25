.. Copyright (c) 2020 RackN Inc.
.. Licensed under the Apache License, Version 2.0 (the "License");
.. Digital Rebar Provision documentation under Digital Rebar master license
.. index::
  pair: Digital Rebar Provision; Release v4.8
  pair: Digital Rebar Provision; Release Notes


.. _rs_release_v48:

Digital Rebar version 4.8 [in process]
--------------------------------------

Release Date: October 2021

Release Themes: Usability, Universal Workflow UX, Cloud Usecase

In addition to bug fixes and performance improvements, the release includes several customer-driven features.

See :ref:`rs_release_summaries` for a complete list of all releases.

.. _rs_release_v48_notices:

Important Notices
~~~~~~~~~~~~~~~~~

* Digital Rebar v4.7 added port 8090 to the list of ports _required_ for provisioning operations. Please verify that port 8090 (default, this can be changed) is accessible for Digital Rebar endpoints.
* Due to changes in the install zip format, the API-based upgrade of DRP to v4.7+ requires usage of most recent https://portal.RackN.io (v4.7 for self-hosted UX users) or the use of DRPCLI v4.6.7+. The v4.7 ``install.sh upgrade`` process also includes theses changes.

.. _rs_release_v48_vulns:

Vulnerabilities
+++++++++++++++

None known

.. _rs_release_v48_deprecations:

Deprecations
++++++++++++

None known

.. _rs_release_v48_removals:

Removals
++++++++

None known


FEATURE TBD
~~~~~~~~~~~

.. _rs_release_v48_terraform:

Terraform and Cloud-Wrapper Updates
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

Cloud Wrapper templates were updated to enable better integration into universal workflow
and simplify importing existing Terraform plans.

1. to use cloud-init instead of relying on Ansible for join-up.
2. to allow creating many machines from a single plan (uses cluster/profile)
3. improve controls after Terraform created instances
4. improve synchronization after Terraform destroys instances

.. _rs_release_v48_secrets:

External Secrets
~~~~~~~~~~~~~~~~

Extend the Digital Rebar parameters to dynamically retrieve information from external sources.

.. _rs_release_v48_pipelines:

UX for Infrastructure Pipelines / Universal Workflow
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

Enhance RackN Portal UX to provide additional user guidance for Universal Workflow components.

.. _rs_release_v48_workunits:

WorkUnits / AsynchActions (Preview)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

Create a new operational mode for the Digital Rebar runner that allows tasks to be performed
asynchronously outside of Workflows.  This allows machines to be addressed as services.

.. _rs_release_v48_vmware:

VMware Cluster Building
~~~~~~~~~~~~~~~~~~~~~~~

Beyond ESXi installation, Digital Rebar will coordinate cluster building activities for completed
vCenter build out including VSAN, NSX-T using vmware lib and other tools.

.. _rs_release_v48_otheritems:

Other Items of Note
~~~~~~~~~~~~~~~~~~~

* TBD
