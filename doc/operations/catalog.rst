.. Copyright (c) 2021 RackN Inc.
.. Licensed under the Apache License, Version 2.0 (the "License");
.. Digital Rebar Platform documentation under Digital Rebar master license
.. index::
  pair: Digital Rebar Platform; Catalog Operations

.. _rs_catalog_ops:

Catalog Operations
==================

This section will describe the catalog system and the tools to use them.


Command Line
------------

The DRP CLI has a section dedicated to ``catalog`` operations.

  ::

    drpcli catalog

This will show the usage for the sections.  The subsection item allows
you to download, install, and get information about a specific catalog item.

Examples:

  ::

    drpcli catalog items                  # Lists the items
    drpcli catalog item show drp          # Shows the possible options for drp
    drpcli catalog item show callback     # Shows the possible options for callback
    drpcli catalog item download drp      # Downloads the stable DRP install image
    drpcli catalog item install callback  # Install the stable callback plugin
    drpcli catalog item install callback --version=v4.6.2 # Install the specific callback version


All the commands can take a web proxy parameter on how to access the outside world, ``-D``

The system will default to the RackN default catalog.  Other catalogs can be referenced by
using the ``-c`` option.

There are three other commands availalbe in the catalog CLI.

The first is ``show``.  This dumps the current catalog.

The next is ``create``.  This takes a version and returns a catalog that only goes back to that version number.
This is useful for building a slimmer catalog than all of RackN's history.

The final command is ``updateLocal``.  This takes a catalog and will make a local repo of that catalog in the
configured DRP instance.  It will download all the contents of the catalog and store it in the files section
of the DRP endpoint specified by RS_ENDPOINT and normal creds.  An optional version can be specified to restrict
the historical items included.  ``updateLocal`` is also incremental.  It will only download new things from the last
time it was run.  The catalog repository is build in the `files/rebar-catalog` directory.

This is really useful for use with the manager functionality.  A cronjob on the manager can be run nightly against
the RackN catalog to update the local repository.  Once that is complete, running the command:

  ::

    drpcli plugins runaction manager buildCatalog

This will build a new local catalog inside the manager.  The new catalog file is located at `files/rebar-catalog/rackn-catalog.json`.
This new catalog is loaded into the DRP endpoint as a content pack.  This is **REQUIRED** for the manager to update downstream endpoints.

Example crontab for maintaining a local catalog repo.

  ::

    drpcli catalog updateLocal --version=v4.6.0
    drpcli plugins runaction manager buildCatalog

The first time this runs it could take a while.


Catalog Details
---------------

A catalog is a collection of ``CatalogItem`` options stored in a content pack format.  This allows it to be loaded into DRP as data.

A DRP endpoint running as a manager will use the loaded catalog as the data elements to update downstream DRP endpoints.

The UX will use the catalog as helper for defining available content to update the DRP endpoint.  The UX uses the `catalog_url` parameter in the global
profile to define which catalog it should use to update the currently attached DRP endpoint.

A catalog item defines a piece of DRP that can be updated:

  * DRP itself
  * DRP UX
  * Plugins, e.g. raid, bios, ipmi, ...
  * Content Packs, e.g. drp-community-content, task-library, ...

The Catalog Item defines the version and type of item.  Additionally, it defines URLs to get the package from.  Checksums and other meta data is also stored
in the catalog item.  For binary items, this will include architecture specific checksums.

This is an example JSON entry for ad-auth version stable.

  ::

    "ad-auth-stable": {
      "ActualVersion": "v4.6.2",
      "Available": false,
      "ContentType": "PluginProvider",
      "Endpoint": "",
      "Errors": [],
      "HotFix": true,
      "Id": "ad-auth-stable",
      "Meta": {
        "Author": "RackN",
        "CodeSource": "https://github.com/rackn/provision-plugins",
        "Color": "green",
        "Copyright": "RackN",
        "Description": "Extends Digital Rebar Provision Authentication to Active Directory",
        "DisplayName": "Active Directory SSO",
        "DocUrl": "https://provision.readthedocs.io/en/latest/doc/content-packages/ad-auth.html",
        "Icon": "key",
        "License": "RackN",
        "Name": "ad-auth",
        "Order": "1000",
        "Prerequisites": "",
        "RequiredFeatures": "",
        "Tags": "exterprise,experimental,auth,active,directory,security"
      },
      "NOJQSource": "ad-auth-stable:::https://s3-us-west-2.amazonaws.com/rebar-catalog/ad-auth/v4.6.2",
      "Name": "ad-auth",
      "ReadOnly": false,
      "Shasum256": {
        "amd64/darwin": "d4c64398bad374560e66bd67fca716b24e64c56ae742dc0824239422306a14e6",
        "amd64/linux": "202f8a67d6aefee9cbef9c2de34dfccd4455e378b4622ed5118b40843c953ef8",
        "arm64/linux": "NOTFOUND",
        "ppc64le/linux": "NOTFOUND"
      },
      "Source": "https://s3-us-west-2.amazonaws.com/rebar-catalog/ad-auth/v4.6.2",
      "Tip": false,
      "Type": "catalog_items",
      "Validated": false,
      "Version": "stable"
    },

This object contains the meta of the object for easy of process.  This is an example of a special case object.
The `Id` field defines the name/version pair of this object.  In this case, the version is `stable`.  This references
the stable version of this object when the catalog was built.  In this case, the `ActualVersion` field defines the
actual version, `v4.6.2`.  All RackN built or buildCatalog action built catalogs have `stable` entries for the most recently
released objects.

There is an additional version, `tip`.  This is present in RackN catalogs that represent the experimental/developmental branch
of the product.  This can be installed like other items.

The Catalog Item assumes that the actual version is a SEMVER compatiable version, e.g. v4.3.2.

The `Shasum256` field defines the shasum256 sum for each arch/os for binary objects.  If the object can be referenced for all
architectures, the value, `any/any` is used.


The `Source` field defines the base path for the object.  Depending upon the type of the object, additional data needs to be added:

  * PluginProvider - `Source`/<arch>/<os>/<name> - e.g. `Source`/amd64/linux/ad-auth
  * ContentPackage - `Source` - No changes needed
  * DRPUX - `Source` - No Changes needed
  * DRP - Special see below.

The DRP download `Source` is special.  If the release v4.6.0 or before, the `Source` field is the single DRP install zip file.  This contains
all arch and OS combinations.  For v4.7.0+, the `Source` field looks like the same thing, but download the image, you need to remove the
`.zip` and replace it with `.<arch>.<os>.zip`.  This will get the specific arch/os version.  Additionally, the catalog item has a complete
Shasum256 field in v4.7.0+.


Custom Catalogs
---------------

RackN provides a catalog built content with a history of 20 or so previous releases.  If you need to add additional catalogs, you will need to do one of these
methods.

Customized Local Catalog
++++++++++++++++++++++++

This is useful for Manager environments.  Using the `updateLocal` and `buildCatalog` commands, you can make your own local catalog repository.  With that in place,
additional content packs can be staged in the files/rebar-catalog directory following the common patterns.  Once staged, rerunning the `buildCatalog` command will
add those injected items into the catalog.  Point the UX at that catalog and it will work.

The layout of the files/rebar-catalog directory works like this.  Create a directory for your new content pack with its `id`.  In this example, `my-workflows` is
the name of the content pack.  A directory would be created, `my-workflows`.  In that directory, a file for each version would be added in the *json* format.
For example, Version v1.1.1 of the my-workflows would be published by creating the file `files/rebar-catalog/my-workflows/v1.1.1.json`.

Plugins and DRP itself have different formats, but they should be obvious from the other items.

Once items are loaded in place, rereun the `buildCatalog` action.  The new catalog is built and added to the DRP endpoint.


Customized Catalog for UX-Only
++++++++++++++++++++++++++++++

Alternatively, the RackN catalog could be downloaded and unbundled, like a content pack.  Once unbundled, additional catalog items can be added following
the sytnax of the existing items.  The `Source` URL can reference a file server hosting your content.  It must follow the `Source` expansion rules above.
Once the items are added, rebundle the content pack use the DRP CLI tools.  You can place that new catalog at a place that can be accessed by the UX through
the `catalog_url` parameter.


