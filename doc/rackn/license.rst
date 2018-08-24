.. Copyright (c) 2018 RackN Inc.
.. Licensed under the Apache License, Version 2.0 (the "License");
.. Digital Rebar Provision documentation under Digital Rebar master license
.. index::
  pair: Digital Rebar Provision; RackN Licensing

.. _rackn_licensing:

RackN Licensing Overview
~~~~~~~~~~~~~~~~~~~~~~~~

This document outlines the RackN commercial licensing information and initial setup steps necessary to access license entitlements.  If you have any questions or concerns, please feel free to contact us on Slack, or email us at support@rackn.com. 

Commercial licensing of Digital Rebar Provision is provided to Organizations.  License entitlements are enabled on the Organization and not on an individual users Portal account.  Your RackN sales and solution team will need to setup an Organization with the correct license entitlements for you.


.. _rackn_licensing_prereqs:

Prerequisites
-------------

Here is a list of the necessary prerequisites that will need to be in place prior to you successfully using any licensed component(s):

1. An `Organization` setup in the Portal with correct license entitlements enabled
2. You must have a Web Portal user account that is registered and functioning (sign up if you do not already have one here: https://portal.rackn.io/#/user/signup)
3. A functioning DRP Endpoint that is managable via the Web Portal

Insure you are logged in to the Rackn Web Portal (using the upper right "login" button.

Log in to the DRP Endpoint - which will be the username/password authentication dialog in the center of the Web Portal if you are not logged in. If you have not changed the default username and password, click the "Defaults" button, then "Login".


.. _rackn_licensing_overview:

Overview of Steps
-----------------

The following are the basic steps you need to perform to generate, enable, and use licensed plugins and contents.

1. Generate a License
2. Enable DRP Endpoints to use Licensed Content
3. Install License Plugin Provider
4. Install Licensed Content and Plugins

.. _rackn_licensing_generate_license:

Generate a License
------------------

The first time that you use a license entitlement, you will need to generate a license.  This creates the and starts the license entitlements based on the terms and condidions of your license (content, plugins, duration of license contract, etc.).  You will need to perform this step only once for each Organization that you manage that has a license entitlement. 

1. Select the Organization in the upper left blue drop down.  For example: "Foo Industries"

.. figure::  images/licensing/01-select-org.png
   :align: left
   :width: 300 px
   :alt: Select Organization

2. Shift-Reload your browser to insure the Org change was successful
3. Go to the "Info & Preferences" menu screen

.. figure::  images/licensing/02-info-prefs.png
   :align: left
   :width: 200 px
   :alt: Info & Preferences Menu Item

4. Verify in the center bottom panel that you see a green check mark and the text *Foo Industries is a Licensed Organization*

.. figure::  images/licensing/03-licensed-org.png
   :align: left
   :width: 300 px
   :alt: Verify Organization is Licensed

5. Click on the blue *Update License* button - there will be a spinning feedback dialog for 10 to 30 seconds

.. figure::  images/licensing/04-spinning.png
   :align: left
   :width: 300 px
   :alt: License Generation Spinner

6. Once completed, you should see the licensed features, contents, and plugin providers

.. figure::  images/licensing/05-generated-license.png
   :align: left
   :width: 300 px
   :alt: Generated Licensed Overview


.. _rackn_licensing_enable_endpoint:

Enable a DRP Endpoint to use Licensed Content
---------------------------------------------

Once you have generated a license, you now need to enable each endpoint that will consume licensed content.  This will allow for Content and Plugins that are licensed to be imported in to the DRP Endpoint and used for provisioning activities. 

.. note:: The DRP Endpoint you initially generated the license on will also be enabled to utilize licensed content and plugins.  You will only need to do this step subsequently for any additional DRP Endpoints that will be using licensed content or plugins.

1. Go to the "Info & Preferences" menu item and click on the "Update License" button for any DRP Endpoint that requires licensed content

.. figure::  images/licensing/02-info-prefs.png
   :align: left
   :width: 200 px
   :alt: Info & Preferences Menu Item


.. _rackn_licensing_license_plugin:

Install License Plugin Provider
-------------------------------

It is necessary to install the *License* Plugin Provider, which works in conjunction with the signed license, plugins, contents, and DRP Endpoint to enable the entitlements specified in the license.  To install the Plugin Provider, do:

1. Go to the *Plugin Providers* menu item

.. figure::  images/licensing/06-plugin-providers.png
   :align: left
   :width: 200 px
   :alt: Plugin Providers Menu Item

2. Locate the *License* plugin in the right side panel (labeled "Organization Plugin Providers")

.. figure::  images/licensing/07-org-license-provider.png
   :align: left
   :width: 350 px
   :alt: Organization License Plugin Provider

3. Click on the "Transfer" link to install the Plugin Provider on the DRP Endpoint (See the above image for details)
4. The *License* plugin provider should now be listed in the "Endpoint Plugin Providers" panel in the center of the page

.. figure::  images/licensing/08-enabled-license-provider.png
   :align: left
   :width: 350 px
   :alt: Endpoint License Plugin Provider


.. _rackn_licensing_use:

Install Licensed Content and Plugins
------------------------------------

Once the above steps have been completed, you may now install licensed Contents and Plugin Providers that you are entitled to use.  This process is very simple, and completed as follows:

1. Go to the *Plugin Providers* menu item

.. figure::  images/licensing/06-plugin-providers.png
   :align: left
   :width: 200 px
   :alt: Plugin Providers Menu Item

2. Select the appropriate plugin from the "Organization Plugin Providers" panel on the right

.. figure::  images/licensing/09-image-deploy-example.png
   :align: left
   :width: 350 px
   :alt: Image Deploy Plugin Provider Example

3. Click "Transfer" to install the Plugin Provider on the DRP Endpoint (see above image for details)
4. Verify the *Plugin Provider* was installed successfully be examining the center "Endpoint Plugin Providers" panel

.. figure::  images/licensing/10-installed-plugin-providers.png
   :align: left
   :width: 350 px
   :alt: Installed Endpoint Plugin Providers

.. _rackn_licensing_verify:

Verify Your License Entitlements
--------------------------------

The "Info & Preferences" page will show an overview of the licensed Contents, Features, and Plugin Providers that the current organization is entitled to.  Please verify you are using the correct Organization to view the licensing rights for that Organization (upper left blue pull down menu item).  If you are currently in the context of your personal Portal account (eg. it shows your email address or account), you will NOT be able to view or manage license entitlements.

Additionally, you can view each individual components entitlements from the overview license page.

1. Click on the "Hamburger" menu in the upper left (three horizontal gray bars)

.. figure::  images/licensing/11-hamburger-menu.png
   :align: left
   :width: 300 px
   :alt: Hamburger Menu

2. Select "Licenses"

.. figure::  images/licensing/12-select-licenses.png
   :align: left
   :width: 200 px
   :alt: Select Licenses

3. Click in the body to the right
4. General license terms will be shown first
5. Each licensed component (feature, content, or plugin provider) will have individual licensing terms and details following the "General" terms

.. figure::  images/licensing/13-license-details.png
   :align: left
   :width: 450 px
   :alt: License Details


The General terms (soft and hard expire dates) will override each individual license expiration terms.  

"Soft" expire is when initial warning messages about subsequent de-licensing of a given feature will occur.

"Hard" expire is the date at which a given featre or term expires and will no longer be active.

