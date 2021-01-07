.. Copyright (c) 2020 RackN Inc.
.. Licensed under the Apache License, Version 2.0 (the "License");
.. Digital Rebar Platform documentation under Digital Rebar master license
.. index::
  pair: Digital Rebar Platform; Manager Operations

    
.. warning:: The Endpoint license installed on the designated DRP Endpoint Manager must have entitlement
             to the Multi Site Manager and High Availability features.  These are Enterprise license
             features.  For each Endpoint that is managed under Multi-Site Manager, the Endpoint must
             also be listed as an authorized Endpoint in the License.
             
             Please contact RackN support with any licensing questions.
          
.. _rs_manager_ops:

Manager Operations
==================

This section will address usage of the manager system.  The architecture and implementation of the Mulit-Site Manager system is
described at :ref:`rs_manager_arch`.

Manager Enablement - Portal
---------------------------

The *Multi-Site Manager* Preferences must be set to *enabled*.  To do so, please do the following in the Portal:

 * go to the *Info & Preferences* page
 * under the *System Preferences* settting panel to the right
 * set the value named *Multi-Site Manager* to *enabled*
 * ensure you hit **Save** in the upper right

Manager Enablement - CLI
------------------------

To enable the manager functions in an endpoint, you will need to set the **manager** preference to **true**.

  ::

    drpcli prefs set manager true

To disable the manager functions in an endpoint, you will need to set the **manager** prefence to **false**

  ::

    drpcli prefs set manager false

Manager Enablement - API
------------------------

The manager setting can be applied with the ``/api/v3/prefs`` POST endpoint API.

The POST requires an object with the ``manager`` string type value set to either ``"true"`` or ``"false"``.  For example:

  ::

    {
      "manager": "true"
    }


Manager Actions
---------------

The manager system provides the same actions that the manager plugin had in ``v4.4.x`` and earlier versions.
These have also been included in the *endpoints* plugin for backwards compatibility reasons.

.. note:: It is recommended to set up a user account that is dedicated to Manager endpoint management
          functions, and to not use the generic "rocketskates" superuser role.  Create a new User and
          assign it the *superuser* role for this purpose.

.. note:: A valid license with Multi-Site Manager enablement must be installed on each DRP Endpoint.
          The Portal addEndpoint method will automatically install the license from the Manager, on
          the Managed endpoints.
          
          HOWEVER - the CLI and API methods do not do this.  You must FIRST install a valid license
          on the Endpoint before adding them under management.


addEndpoint - Portal
____________________

Use the *Endpoints* menu entry and the blue *Add* button at the top center to add a new Endpoint.

You must provide the fully qualified Endpoint URL (eg ``https://192.168.1.10:8092``) and a
username/password pair that has the *superuser* role.

addEndpoint - CLI
_________________

Utilizing the plugin actions mechanism, it is possible to add a new Endpoint under management in a single CLI step with the following command line:

  ::
  
    # install a valid license on the to be managed Endpoint first
    # assumes rackn-license.json carries a valid license with entitlements
    drpcli contents create rackn-license.json
    # bring the endpoint under management
    drpcli plugins runaction manager addEndpoint manager/url https://192.168.1.10:8092 manager/username manager manager/password manager-password

The ``manager/url`` is the remote Endpoint URL that is being added in to the Manager system for management.

addEndpoint - API
_________________

The `/api/v3/plugins/manager/addEndpoint` POST endpoint allows for the auto-creation of the Endpoint object and
validates connectivity.

You must first install a valid License on the remote Endpoint that will be managed.  To do this; use the API ``/api/v3/contents`` endpoint with a POST operation containing the content pack ``rackn-license.json`` as a payload.

The ``addEndpoint`` POST requires an object with the ``manager/url``, ``manager/username``, and ``manager/password`` values.  For example:

  ::

    {
      "manager/url": "https://192.168.124.10:8092",
      "manager/username": "constantBackUpUsername",
      "manager/password": "constantBackUpPassword"
    }

This will validate the credentials and add the Endpoint object for that system named correctly.  This is the preferred
method to add an endpoint to the manager instead of directly creating the Endpoint object and then populating it with values.


buildCatalog
____________

The ``/api/v3/plugins/manager/buildCatalog`` POST endpoint allows for the building or rebuilding of the local
catalog of content packages and plugin providers.

The Manager will use its local catalog when applying content.  The catalog is actually a content pack that is
loaded into the manager.  The default RackN catalog can be used and it will reference the internet, but often
times content would like to be cached locally or expanded with additional components.

The **buildCatalog** action process the cache directory, ``files/rebar-catalog``, and builds a catalog content
package and stores as ``rackn-catalog.json`` in that directory.  It will then load that content package into the
manager.

The ``files/rebar-catalog`` directory can be populated by the ``drpcli catalog updateLocal``.  This will by default
use the RackN catalog to cache all content locally.  You can also provide options to the command to handle additional
catalogs or firewall and proxies.  This will provide the layout for the catalog directories.

Custom content can be added to the catalog directories.  You will need to follow the format for plugin providers
or content packages.  The **files** api can be used to update the catalog.

The UX has a helper button for this action (on the *Endpoints* menu, as "*Rebuild Catalog*" button).


Proxy Creating an Object on Managed Endpoint
____________________________________________

It is possible to create objects on managed endpoints by using proxy pass-through from the manager.  Details are available in :ref:`rs_api_proxy`.


Manager Common Methods
----------------------

Here are some common manager actions.

Create and Populate the Initial Catalog Cache
_____________________________________________

The Manager requires a local catalog cache to install and manage items on downstream endpoints.  This catalog can be built and initialized
with the following commands.

To reduce the amount of content downloaded, the ``--version`` flag can be used to specify the minimum version to download
from the catalog.  For example, ``--version=v4.5.4`` would download things newer and including ``v4.5.4``.  If left off,
the command will download the whole catalog.

  ::

    # Put the current license into the catalog.
    drpcli contents show rackn-license > /tmp/v0.0.1.json
    drpcli files upload /tmp/v0.0.1.json as rebar-catalog/rackn-license/v0.0.1.json
    rm -f /tmp/v0.0.1.json

    # Create the initial empty catalog
    drpcli plugins runaction manager buildCatalog

    # Populate the local cache from all items found in the system Catalog
    drpcli catalog updateLocal
    # OR limit local cache to v4.5.4 and newer only items
    drpcli catalog updateLocal --version=v4.5.4

    # Rebuild the initial catalog
    drpcli plugins runaction manager buildCatalog


Update Catalog Cache
____________________

Once the catalog is initialized, you can incrementally update the catalog with the following commands.  This can
be put into a cron job to keep the catalog up to date.

To reduce the amount of content downloaded, the ``--version`` flag can be used to specify the minimum version to download
from the catalog.  For example, ``--version=v4.5.4`` would download things newer and including ``v4.5.4``.  If left off,
the command will download the whole catalog.

  ::

    # Populate the local cache from all items found in the system Catalog
    drpcli catalog updateLocal
    # OR limit local cache to v4.5.4 and newer only items
    drpcli catalog updateLocal --version=v4.5.4

    # Rebuild the initial catalog
    drpcli plugins runaction manager buildCatalog

