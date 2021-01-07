.. Copyright (c) 2020 RackN Inc.
.. Licensed under the Apache License, Version 2.0 (the "License");
.. Digital Rebar Platform documentation under Digital Rebar master license
.. index::
  pair: Digital Rebar Platform; Manager Operations

.. _rs_manager_ops:

Manager Operations
==================

This section will address usage of the manager system.  The architecture and implementation of the pooling system is
described at :ref:`rs_manager_arch`.

All of these operations are limited by the licenses.  You will need to have a license with the Manager and High
Availability features.

Manager Enablement
------------------

To enable the manager functions in an endpoint, you will need to set the **manager** preference to **true**.

  ::

    drpcli prefs set manager true

To disable the manager functions in an endpoint, you will need to set the **manager** prefence to **false**

  ::

    drpcli prefs set manager false

Manager Actions
---------------

The manager system provides the same actions that the manager plugin had.  These are on the plugin endpoints for
backwards compatibility reasons.

addEndpoint
___________

The `/api/v3/plugins/manager/addEndpoint` POST endpoint allows for the auto-creation of the Endpoint object and
validates connectivity.

The POST requires an object with the `manager/url`, `manager/username`, and `manager/password` values.  For example:

  ::

    {
      "manager/url": "https://192.168.124.2:8092",
      "manager/username": "constantBackUpUsername",
      "manager/password": "constantBackUpPassword"
    }

This will validate the credentials and add the Endpoint object for that system named correctly.  This is the preferred
method to add an endpoint to the manager instead of directly creating the Endpoint object.

The UX has a helper button for doing this action.

buildCatalog
____________

The `/api/v3/plugins/manager/buildCatalog` POST endpoint allows for the building or rebuilding of the local
catalog of content packages and plugin providers.

The Manager will use its local catalog when applying content.  The catalog is actually a content pack that is
loaded into the manager.  The default RackN catalog can be used and it will reference the internet, but often
times content would like to be cached locally or expanded with additional components.

The **buildCatalog** action process the cache directory, `files/rebar-catalog`, and builds a catalog content
package and stores as rackn-catalog.json in that directory.  It will then load that content package into the
manager.

The `files/rebar-catalog` directory can be populated by the `drpcli catalog updateLocal`.  This will by default
use the RackN catalog to cache all content locally.  You can also provide options to the command to handle additional
catalogs or firewall and proxies.  This will provide the layout for the catalog directories.

Custom content can be added to the catalog directories.  You will need to follow the format for plugin providers
or content packages.  The **files** api can be used to update the catalog.

The UX has a helper button for this action.

Proxy Creating an Object on Managed Endpoint
____________________________________________

It is possible to create objects on managed endpoints by using proxy pass-through from the manager.  Details are available in :ref:`rs_api_proxy`.


Manager Common Methods
----------------------

Here are some common manager actions.

Load Catalog Cache
__________________

The Manager requires a local catalog cache to install downstream endpoints.  This catalog can be built and initialized
with the following commands.

To reduce the amount of content downloaded, the `--version` flag can be used to specify the minimum version to download
from the catalog.  For example, `--version=v4.5.3` would download things newer and including `v4.5.3`.  If left off,
the command will download the whole catalog.

  ::

    # Put the current license into the catalog.
    drpcli contents show rackn-license > /tmp/v0.0.1.json
    drpcil files upload /tmp/v0.0.1.json as rebar-catalog/rackn-license/v0.0.1.json
    rm -f /tmp/v0.0.1.json

    # Create the initial catalog
    drpcli plugins runaction manager buildCatalog

    # Load the Cache
    drpcli catalog updateLocal

    # Rebuild the initial catalog
    drpcli plugins runaction manager buildCatalog


Update Catalog Cache
____________________

Once the catalog is initialized, you can incrementally update the catalog with the following commands.  This can
be put into a cron job to keep the catalog up to date.

To reduce the amount of content downloaded, the `--version` flag can be used to specify the minimum version to download
from the catalog.  For example, `--version=v4.5.3` would download things newer and including `v4.5.3`.  If left off,
the command will download the whole catalog.

  ::

    # Load the Cache
    drpcli catalog updateLocal

    # Rebuild the initial catalog
    drpcli plugins runaction manager buildCatalog

