.. Copyright (c) 2021 RackN Inc.
.. Licensed under the Apache License, Version 2.0 (the "License");
.. Digital Rebar Provision documentation under Digital Rebar master license
.. index::
  pair: Digital Rebar Provision; RackN Licensing

.. _rackn_licensing:

RackN Licensing Overview (v4.6+)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

This document outlines the RackN limited use and commercial licensing information and initial setup steps necessary to access licensed entitlements.  If you have any questions or concerns, please feel free to contact us on Slack, or email us at support@rackn.com.

If you have a license issued before May 1, 2021, please see :ref:`rackn_licensing_pre46`.  You can identify these licenses because they will end in a number sequence greater than 4600000.  Licenses issued prior to v4.6 will continue to function and can be used to create v4.6+ licenses; however, new entitlements cannot be added to earlier licenses.

License Types
=============

Start-up / Embedded licensing
-----------------------------

*Embedded licensing* of RackN Digital Rebar Platform is provided for solely new endpoints to run until you install get an official use license.

The embedded license is built into the platform and has restricted entitlements:

* 10 machines (10)
* no contexts or pools
* restricted use of plugins
* hard coded expiration date based on binary build date

Limited Use licensing
---------------------

*Limited Use licensing* of RackN Digital Rebar Platform is provided for individual users, trial and non-commercial teams.

The trial/communty license is generated via the processes documented below. Self-service
licenses generally start at:

* 20 machines
* 5 contexts
* 5 pools
* 30-days with 1 self-service renewal
* allow access to all publically available plugins in the RackN catalog.

Enterprise features such as Airgap, HA and Secure boot are not included in the trial license.  Contact the RackN solution team if you would like to expand your entitlements.

Commercial Use licensing
------------------------

*Commercial Use licensing* of RackN Digital Rebar Platform is 
provided to named Organizations.  License entitlements are enabled on several different dimensions
including endpoint id, unit counts, contexts, pools, plugins and advanced features like HA and 
secure boot.  The RackN solution team will need to setup an Organization with the correct license entitlements for this type of license.

.. _rackn_licensing_file:

License Enforcement Mechanism
=============================

The RackN License file restricts operatations for a Digital Rebar endpoint based on different entitlements.  The file contains all the information needed to validate the service so *no external access is required*

RackN License File
------------------

The RackN license 

Entitlement files simply Digital Rebar Content Packs named `rackn-license`.  Installing the Content
creates a `rackn license` Profile with two variables:

#. `rackn/license` is the signed and encrypted license that is used by the server
#. `rackn/license-object` is the human readable version of that license file

The Digital Rebar server looks for the `rackn/license` parameter to verify entitements.

Entitlements
------------

The entitlements are verified on several dimensions including:

* Endpoint ID
* Date
* Plugins in Use
* Object Counts.
   * Machines
   * Contexts
   * Pools
   * Subnets
* High Availability (HA) enabled
* Secure Boot enabled
* Air-Gap enabled

Installling RackN License Files
===============================

The RackN License file is a Digital Rebar content pack that contains the license key.  It can be managed in the system using the Contents API and CLI like any other content pack.


Verify and Update
-----------------

The RackN UX integrates with our license generating API and can be used to authenticate and download updated entitlement files.


.. _rackn_licensing_check:

Check a License from the UX
---------------------------

.. note:: You must log in to the DRP Endpoint first.

From the "License Manager" page, check the "License Management" panel.  The panel performs
three levels of checks on your license:

#. Verified: a valid entitlement file is installed.
#. Registered: the current endpoint is registered in that license.
#. Up-to-date: the version of the file matches the most current version known to RackN.

If the license is not up-to-date then click the "Check and Update License" button to
retrieve an updated license from the RackN entitlement service.

.. _rackn_licensing_generate_license:

Generate a New License
----------------------

.. note:: You must log in to the DRP Endpoint first.

The first time that you login to a Digital Rebar endpoint from the RackN UX, you will be prompted to either request new license file or upload an existing license.  If you complete the request information then the service creates a short term `rackn-license` content file and then uploads it to your endpoint.


You will only need to perform this step once because the license file is used to validate your access rather than a RackN specific login or password.  For this reason, it is important to download and store the license file for future use.  Even if a newer license if issued, previous licenses can still be used to validate your identity to the RackN license service.

IMPORTANT: Licenses created via the self-enrollment process have limited time spans.  You will need to contact RackN via the UX or email to verify your self-enrollment license and expand the entitlement settings.

Once completed, you should see the entitlements in the "License Management" panel.

.. _rackn_licensing_update_license:

Update a License from Existing License
--------------------------------------

.. note:: You must log in to the DRP Endpoint first.

Once a valid license file is installed on a DRP Endpoint, the "Check and Update License" button
on the "License Manager" page can be used to update the license entitlements.

This is typically needed when a new endpoint is being added to an entitlement or RackN has
expanded the entitlements and triggered a new version.

.. _rackn_licensing_save_license:

Save an Entitlements File for Backup
------------------------------------

.. note:: You must log in to the DRP Endpoint first.

You can download the current entitlement file using the RackN UX by pressing the "Download" button
on the "License Manager" page.  You can also use `drpcli contents show rackn-license > rackn-license.json` from the command line.

.. _rackn_licensing_install_license:

Install a License from a File
-----------------------------

.. note:: You must log in to the DRP Endpoint first.

The Digital Rebar entitlements file, typically `rackn-license-org.json`, is just a Content pack.
You can upload it from Catalog Import in RackN UX or using `DRPCLI contents upload rackn-license-org.json`
` 
Once a valid license is installed, the key in the license can be used to retrieve an updated license.


.. _rackn_licensing_verify:

Verify Your License Entitlements
--------------------------------

The "License Manager" page will show an overview of the licensed Contents, Features, and Plugin Providers of the installed entitlements file.  Please verify you are using the correct Organization

* "Soft" expire is when initial warning messages about subsequent de-licensing of a given feature will occur.  At this date, the system is considered out of compliance but will continue to operate.
* "Hard" expire is the date at Digital Rebar will disable the relevant features or stop accepting add/update requests.

Many licenses, including trial/community licenses, use the "upto-nodes" module which allows operators to use *any* licensed content up to the stated number of machines.


Check or Update an Existing License
------------------------------------

To update a license, visit the UX *License Management* page.
Click the "Check and Update License" button in the top right
corner of the "License Management" panel.  This uses the API
described below to update your license including adding new
endpoints.

.. _rackn_licensing_api_upgrade:

Non-UX Update an Existing License
=================================

These steps require that you already have a valid RackN license.
The information contained in the license is used to verify your
entitlements and to authorize an updated license.  It relies on
online RackN License Management APIs.


Required Header Fields:

* `rackn-endpointid`: digital rebar endpoint id
* `Authorization`: license key

.. note:: The `rackn-endpointid` is the endpoint id (aka `drpid`) of the Digital Rebar Provision endpoint to be licensed.  Licenses are issued per endpoint.  

  ::

    export ENDPOINTID=$(drpcli info get | jq -r .ha_id)
    export KEY=$(drpcli profiles get rackn-license param rackn/license | jq -r)

    curl -X POST -H "rackn-endpointid: $ENDPOINTID" \
      -H "Authorization: $KEY" \
      -d '$(drpcli info get)' \
      https://cloudia.rackn.io/api/v1/license/update

Adding Endpoints to a License
-----------------------------

Generally, the UX will add endpoints automatically during the Check and Update process on a new endpoint.

If you cannot use the UX to add an endpoint then use the API above.  Add endpoints to a license upto your entitlements by sending a new endpoint with license information validated for a different endpoint.  This will create a new license that can be applied too all endpoints.  


Removing Endpoints from a License
----------------------------------

There is no automated process to REMOVE endpoints from a license.  Contact RackN if you need to do this.