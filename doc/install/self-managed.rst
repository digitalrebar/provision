.. Copyright (c) 2021 RackN Inc.
.. Licensed under the Apache License, Version 2.0 (the "License");
.. Digital Rebar  documentation under Digital Rebar master license
.. index::
  pair: Digital Rebar Provision; Self-Managed

.. _rs_self_managed:

Digital Rebar is Self-Managed Software
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

Digital Rebar is software designed to be self-managed by the infrastructure operator, not remotely by RackN as a SaaS. For this reason, users must install a Digital Rebar service endpoint before using the software.  Further, RackN never has access to the endpoint APIs or credentials.

RackN Digital Rebar offers unique approach to delivering infratructure software.  We're proud to be innovating in ways that allow operators to control their own destiny and keep their secret secrets while also providing the very latest feature set and capabilities.

This page addresses common questions about Digital Rebar as self-managed software.

.. _rs_self_managed_why:

Why do I have to install Digital Rebar locally?
-----------------------------------------------

Since RackN is not software as a service (SaaS), you must install and then manage Digital Rebar yourself.  We work hard to make this as easy as possible.  More importantly, it means that you are in control of your own experience.

.. _rs_self_managed_self_managed:

Why did RackN choose this self-management approach?
---------------------------------------------------

We believe the better question is why should you trust your infrastructure operations, cloud credentials and secrets to a database that you can't touch, you can't see or inspect how it's guarded?

Digital Rebar is self-managed software because we do not believe customers should share confidential information just to get better management software.

To do that, we allow our customers run and install the software themselves on their premises and behind their firewalls. RackN engineers never have access to their passwords, databases, or credentials to core infrastructure control systems.  And that's the way we think it should be.

.. _rs_self_managed_work:

Does self-management mean more work for my operations team?
-----------------------------------------------------------

The purpose of Digital Rebar is removing the common, repetitive, non value-added toil that most IT operations teams struggle with every day.  RackN works day and night to make Digital Rebar operations easy.  It does not require complex components, infrastructure, and knowledge because we think that's the way it should be.

We help operations teams two ways:

#. by providing an automation catalog that covers the most common hard IT challenges like installing operating systems, updating firmware and building inventories.
#. by including a strong Infrastructure as Code process that encourages teams to reuse and standardize processes beyond our catalog

We know that a strong infrastructure foundation leaves companies to focus on the value added work.

.. _rs_self_managed_ux:

Why doees the RackN UX launch from RackN.io if my install is local?
-------------------------------------------------------------------

The RackN UX is a "single page application" based on React.  It does not matter where the inital application is loaded from because that is just downloading the application to your web browser.  Once it's loaded, it connects directly between your client browser and your Digital Rebar endpoint.  Your Digital Rebar API traffic is local to your networks.

Our enterprise licensed customers have the option of running the UX, and every other component, completely behind their firewall.  The Digital Rebar has a static webserver integrated specifically for this purpose.

For new users, we found that keeping the UX hosted one page app react delivered application significantly reduces their startup and maintenance effort.  In fact, most enterprise customers choose to use the RackN.io hosted UX for daily operations instead of maintaining the UX locally.

.. _rs_self_managed_license:

Do I have to have a license to use Digital Rebar?
-------------------------------------------------

Yes.  Digital Rebar requires a license for use and RackN makes trial licenses available at no charge and without pre-authorization.

Licenses do NOT require online connectivity until they need to be updated or changed.  Once created the the license is installed on the local Digital Rebar.  If you cannot access online license portal, then RackN can generate and send a license to you.

.. _rs_self_managed_faq:

Common Questions about Self-Managed Software
--------------------------------------------

* Can Digital Rebar work without internet access? Yes. We call that "air gap" mode and there are degrees of air gapping
* I can host the UX myself?  Yes.  Depending on your license tier, RackN makes the UX available for local install.  Digital Rebar is able to host the UX from the Endpoint.
* Can RackN access my system?  No.  We do not store or maintain any credentials or access to your system.
* What data does RackN collect?
  * From the Endpoint, nothing.  The exception is when the the Billing plugin is installed for activity based users.
  * From the UX, we collect the following types of data:
    * the data available from ``drpcli info get``
    * the license in use on the endpoint
    * app page usage via Google Analytics
    * synchronization with RackN message queue
    * For "air-gap" enabled licenses, all of the following are disabled.
* How can I secure Digital Rebar?  See :ref:`rs_security_faq`

.. _rs_self_managed_install:

How do I install Digital Rebar?
-------------------------------

Common installation paths:

* :ref:`rs_quickstart` is a basic SystemD install for new users
* :ref:`rs_install_dev` for developers running DRP interactively
* :ref:`rs_install_docker` for trial users minimizing their install requirements
* :ref:`rs_install_cloud` is non-PXE / Cloud-Only installation process
* `Edge Lab with RPi <http://edgelab.digital>`_ is self-contained Digital Rebar inexpensive lab using Raspberry Pi computers.

