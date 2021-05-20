:orphan:

.. Copyright (c) 2018 RackN Inc.
.. Licensed under the Apache License, Version 2.0 (the "License");
.. Digital Rebar Provision documentation under Digital Rebar master license
.. index::
  pair: Digital Rebar Provision; RackN Licensing

.. _rackn_licensing_pre46:

RackN Licensing Overview (prior to v4.6)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

.. note:: If you have a license issued after May 1, 2021, please see :ref:`rackn_licensing`.  You can identify these licenses because they will end in a number sequence greater than 4600000.

Licenses issued prior to v4.6 will continue to function and can be used to create v4.6+ licenses; however, new entitlements cannot be added to earlier licenses.

This document outlines pre-v4.6 RackN limited use and commercial licensing information and initial setup steps necessary to access licensed entitlements.  If you have any questions or concerns, please feel free to contact us on Slack, or email us at support@rackn.com.

Watch the `license training video <https://youtu.be/wIGaSQevjfM!>`_

License Types
-------------

Embedded licensing
==================

*Embedded licensing* of RackN Digital Rebar Platform is provided for bootstrapping endpoints

The embedded license is built into the platform and has restricted entitlements:

* 10 machines (10)
* no contexts or pools
* restricted use of plugins

Limited Use licensing
=====================

*Limited Use licensing* of RackN Digital Rebar Platform is provided for individual users, trial and non-commercial teams.

The trial/communty license is generated via the processes documented below. Self-service
licenses start at:

* 20 machines
* 3 contexts
* 1 pools
* 90-day self-service renewal.  

They allow access to all publically available plugins in the RackN catalog.  HA and Secure boot are not included in the trial license.  Contact the RackN solution team if you would like to expand your entitlements.

Commercial Use licensing
========================

*Commercial Use licensing* of RackN Digital Rebar Platform is
provided to named Organizations.  License entitlements are enabled on several different dimensions
including endpoint id, unit counts, contexts, pools, plugins and advanced features like HA and 
secure boot.  The RackN solution team will need to setup an Organization with the correct license entitlements for this type of license.

.. _rackn_licensing_prereqs:

Prerequisites
-------------

Here is a list of the necessary prerequisites that will need to be in place prior to you successfully using any licensed component(s):

#. You must have a Web Portal user account that is registered and functioning (sign up if you do not already have one here: https://portal.rackn.io/#/user/signup)
#. A functioning DRP Endpoint that is managable via the Web Portal

If you cannot meet these prerequisites, please contact RackN for alternative ways to create an
entitlement file.

Entitlement Enforcement Mechanism
---------------------------------

Entitlement files simply Digital Rebar Content Packs named `rackn-license`.  Installing the Content
creates a `rackn license` Profile with two variables:

#. `rackn/license` is the signed and encrypted license that is used by the server
#. `rackn/license-object` is the human readable version of that license file

The Digital Rebar server looks for the `rackn/license` parameter to verify entitements.

The entitlements are verified on several dimensions including:

* Endpoint ID
* Date
* Plugins in Use
* Counts
   * Machines
   * Contexts
   * Pools
* High Availability (HA) enabled
* Secure Boot enabled


.. _rackn_licensing_check_pre46:

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

.. _rackn_licensing_generate_license_pre46:

Generate a License from RackN Account
-------------------------------------

.. note:: You must log in to the DRP Endpoint first.

Log in to your Rackn Account from the "License Manager" page and "Online Activation and
Support" panel.  If you do not have an account, then you will need to create and verify it
before you can continue.

.. note:: Your RackN Account is different from a Digital Rebar login.  They are only used to create or update entitlement files.  RackN can set up multiple accounts for the same organization or have a single account that supports multiple organizations.

The first time that you activate a license entitlement, you will need to "Authorize" new license file from the "License" tab.  This creates the `rackn-license` content and then uploads it to your endpoint.  You will need to perform this step only once for each Organization that you manage that has a license entitlement.

Once completed, you should see the entitlements in the "License Management" panel.

.. _rackn_licensing_update_license_pre46:

Update a License from Existing License
--------------------------------------

.. note:: You must log in to the DRP Endpoint first.

Once a valid license file is installed on a DRP Endpoint, the "Check and Update License" button
on the "License Manager" page can be used to update the license entitlements.

This is typically needed when a new endpoint is being added to an entitlement or RackN has
expanded the entitlements and triggered a new version.

.. _rackn_licensing_save_license_pre46:

Save an Entitlements File for Backup
------------------------------------

.. note:: You must log in to the DRP Endpoint first.

You can download the current entitlement file using the RackN UX by pressing the "Download" button
on the "License Manager" page.  You can also use `drpcli contents show rackn-license > rackn-license.json` from the command line.

.. _rackn_licensing_install_license_pre46:

Install a License from a File
-----------------------------

.. note:: You must log in to the DRP Endpoint first.

The Digital Rebar entitlements file, typically `rackn-license.json`, is just a Content pack.
You can upload it from Catalog Import in RackN UX or using `DRPCLI contents upload rackn-license.json`

Once a valid license is installed, the key in the license can be used to retrieve an updated license.


.. _rackn_licensing_verify_pre46:

Verify Your License Entitlements
--------------------------------

The "License Manager" page will show an overview of the licensed Contents, Features, and Plugin Providers of the installed entitlements file.  Please verify you are using the correct Organization

* "Soft" expire is when initial warning messages about subsequent de-licensing of a given feature will occur.  At this date, the system is considered out of compliance but will continue to operate.
* "Hard" expire is the date at Digital Rebar will disable the relevant features or stop accepting add/update requests.

Many licenses, including trial/community licenses, use the "upto-nodes" module which allows operators to use *any* licensed content up to the stated number of machines.

.. _rackn_licensing_api_upgrade_pre46:

Check or Update an Existing License
------------------------------------

These steps require that you already have a valid RackN license.
The information contained in the license is used to verify your
entitlements and to authorize an updated license.  It relies on
online RackN License Management APIs.

To update manually, visit the UX *License Management* page.
Click the "Check and Update License" button in the top right
corner of the "License Management" panel.  This uses the API
described below to update your license including adding new
endpoints.

To update automatically using the APIs, you must make the
a GET call with the required rackn headers.  If successful,
the call will return the latest valid license.  If a new
license is required, it will be automatically generated.

The most required fields are all avilable in the `sections.profiles.Params`
section of the License JSON file.

* `rackn-ownerid` = `[base].rackn/license-object.OwnerId`
* `rackn-contactid` = `[base].rackn/license-object.ContactId`
* `rackn-key` = `[base].rackn/license`
* `rackn-version` = `[base].rackn/license-object.Version`

The URL for the GET call is subject to change!  The current
(Nov 2019) URL is `https://1p0q9a8qob.execute-api.us-west-2.amazonaws.com/v40/license`

For faster performance, you can also use `https://1p0q9a8qob.execute-api.us-west-2.amazonaws.com/v40/check`
with the same headers to validate the license before asking for
updates.

Required Header Fields:

* `rackn-ownerid`: license ownerid / org [or 'unknown']
* `rackn-contactid`: license contactid / cognitor userid [or 'unknown']
* `rackn-endpointid`: digital rebar endpoint id [or 'unknown']
* `rackn-key`: license key [or 'unknown']
* `rackn-version`: license version [or 'unknown']

.. note:: The `rackn-endpointid` is the endpoint id (aka `drpid`) of the Digital Rebar Provision endpoint to be licensed.  Licenses are issued per endpoint.  You can add endpoints to a license by sending a new endpoint with license information validated for a different endpoint.  This will create a new license that can be applied too all endpoints.

With header values exported, an example CURL call would resemble:

  ::

    curl GET -H "rackn-contactid: $CONTACTID" \
      -H "rackn-ownerid: $OWNERID" \
      -H "rackn-endpointid: $ENDPOINTID" \
      -H "rackn-key: $KEY" \
      -H "rackn-version: $VERSION" \
      https://1p0q9a8qob.execute-api.us-west-2.amazonaws.com/v40/license
