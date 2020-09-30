.. Copyright (c) 2020 RackN Inc.
.. Licensed under the Apache License, Version 2.0 (the "License");
.. Digital Rebar Provision documentation under Digital Rebar master license
.. index::
  pair: Digital Rebar Provision; Manager Architecture

.. _rs_manager_arch:

Manager Architecture
<<<<<<<<<<<<<<<<<<<<

Digital Rebar Provision provides a mechanism to manage

This section will focus on the architecture of the manager system.  The operational methods can be found at
:ref:`rs_manager_ops`.

Manager provide two basic functions.

* Single pane of glass-like aggregated API operations
* Aggregated management of additional Digital Rebar Provision endpoints.

By registering the endpoints with the manager, the manager provides a set of API endpoints that aggregate the
machines from all those endpoints.  This allows filtering and operating on machines (and other objects) across
multiple endpoints from the one API target.

Additionally, the manager can act as an update and configuration controller for the attached endpoints.  Using
Infrastructure as Code (IaC) methods, the configuration of downstream endpoints can be defined and automated.  The
manager can also maintain and provide a catalog of resources to install and use.

Manager Hierarchy
-----------------

The **Manager** system allows for chained endpoints with multiple levels of management.  The new **Endpoint** object
defines what endpoints the **Manager** manages.

Once an **Endpoint** object is created in the **Manager**, the **Manager** will use the high-availability data
replication system to maintain a local instance of the endpoint's data.  This will be used to answer "read"
operations.  All write and action operations will be forwarded to the owning endpoint.

Because of the replication of objects, a **Manager** registered to a **Manager** will allow the top-level manager
to control the endpoints attached to the downstream manager.  All communication will be funnelled through the
manager to allow for networking isolation or separation.

For example, if we have 7 Endpoints, A, B, C, D, R1, R2, and M1.  A, B, C, D are leaf endpoints managing machines.
We have two regional endpoint / managers, R1 and R2.  A and B are managed by R1 and C and D are managed by R2.
Manager M1 is managing R1 and R2.

  ::

    A <----|                | ---> C
           R1 <--- M1 ---> R2
    B <----|                |----> D

In this example, R1 has **Endpoint** objects for A and B.  R2 has **Endpoint** objects for C and D.
M1 has **Endpoint** objects for R1 and R2.

Because of the object replication, M1 will see the **Endpoint** objects for A, B, C, and D from R1 and R2.  This will
enable those **Endpoints** to be operated on by M1 (both single pane of glass and configuration).

M1 will present all **Machine** objects from all endpoints.  R1 will present **Machine** objects from A, B, and R1.
R2 will presents **Machine** objects from C, D, and R2.

Manager Functions
-----------------

The two main functions of manager are discussed here.

Single Pane of Glass
====================

When making API requests, the manager will provide results from all the attached endpoints.  Additionally, all objects have a field, **Endpoint**, that is populated with the Endpoint Id (the High Availability Id) of the owning endpoint.  API requests made to objects from attached endpoints are
automatically forward to the correct endpoint.  See :ref:`rs_api_proxy` for more information about
using this automatic API behavior explicity.

Only the local objects are replicated up to the manager, objects provided by content packages and plugins are not
replicated to the manager.  It is assumed that the manager will have all the content packages and plugins loaded to
resolve dependencies and parameters.

Configuration Management
========================

The manager system adds a new object, **VersionSet**, that allows for the definition of configuration, content
packages, plugins, files, and isos in a defined versioned set.

Applying a list of VersionSets to an Endpoint will cause the manager to ensure that the state defined by the
version sets is applied to the system.

Manager works with the High Availability system.  A set of DRP endpoints in an HA cluster are operated on as a single
endpoint.  The configuration is replicated across all cluster members.  The only exception is the version of DRP itself.
The UX provides mechanisms to update the cluster members version of DRP.  Once DRP is updated, the High Availability
system will replicate remaining content.

The Manager uses a local catalog to push configuration elements to the downstream endpoints.  The catalog can be built
with CLI commands.  The local catalog can also be extended to contain custom user content.  These can be reference
by version sets to push out to endpoints.

Manager Implementation
----------------------

The following sections define the implementation details of the Manager system.

Any DRP Endpoint with a license with the manager feature can be come a manager.  Whether an endpoint operates
as a manager is controlled by the preference, **manager**.  Setting this preference to **true** will allow the
endpoint to operate as a manager.  Setting it to **false** (the default) will cause the Manager functions to
stop.  VersionSet and Endpoint Objects will be maintained while not acting as a manager.

Endpoint Object
===============

The **Endpoint** object defines the API connection path (URL, username, and password) to the managed endpoint.  As
the manager operates, the **Endpoint** object is updated with the state and configuration of the endpoint.  The
**Endpoint** object also controls if the manager should update configuration by VersionSets.

The **Endpoint** object has the standard object fields of `Description`, `Documentation`, `Errors`, `Available`,
`Valid`, `Meta`, `Endpoint`, and `Bundle`.  These are like all other objects.

The **Endpoint** object identity fields are:

* **Id** - The name of the endpoint.  This should be the HaId from the endpoint in question.  It will be automatically set by join.
* **HaId** - The HaId of the endpoint.  It is set automatically.
* **Arch** - The architecture of the endpoint.  It is set automatically.
* **Os** - The Os of the endpoint.  It is set automatically.

The **Endpoint** object connection fields are:

* **Params** - Parameters defined for this **Endpoint** object.
   * `manager/url` - defines the URL to use to access this **Endpoint**
   * `manager/username` - defines the username to use to access this **Endpoint**
   * `manager/password` - defines the password to use to access this **Endpoint**
* **ConnectionStatus** - a string of the current state of the endpoint (connecting or updating).

The **Endpoint** object state fields are:

* **DRPVersion** - defines the currently installed DRP Version
* **DRPUXVersion** - defines the currently installed DRP UX Version
* **Components** - a list of Elements installed (these are content packages and plugin providers)
   * **Name** - Name of Element - e.g. `burnin`, `bios`, ...
   * **Type** - Type of Element - e.g. `CP`, `PP`
   * **Version** - Version of the element (short form) - e.g. `tip`, `stable`, `v4.3.2`
   * **ActualVersion** - Actual Version of the element (long form)
* **Plugins** - a list of current plugin configurations.  The list is *Plugin* objects from the system.
* **Prefs** - A map of string pref values that contain all the preferences from the endpoint.
* **Global** - The contents of the global profile on the system.

The **Endpoint** object configuration fields are:

* **Apply** - A boolean field that if true causes the version sets to be applied to the system.
* **VersionSet** - a single *VersionSet* object name.  THIS IS DEPRECATED.
* **VersionSets** - a list of *VersionSet* object names to apply to the system.
* **Actions** - This contains a list of actions that need to be applied to the system.  This can be used to preview changes.


VersionSet Object
=================

The **VersionSet** Object defines a set of configuration state that could be applied to an endpoint.

The **VersionSet** object has the standard object fields of `Description`, `Documentation`, `Errors`, `Available`,
`Valid`, `Meta`, `Endpoint`, and `Bundle`.  These are like all other objects.

The **VersionSet** object identity fields are:

* **Id** - The name of the version set.

The **VersionSet** configuration fields are:

* **DRPVersion** - defines desired installed DRP Version
* **DRPUXVersion** - defines the desired installed DRP UX Version
* **Components** - a list of Elements to install (these are content packages and plugin providers)
   * **Name** - Name of Element - e.g. `burnin`, `bios`, ...
   * **Type** - Type of Element - e.g. `CP`, `PP` - This is required.
   * **Version** - Version of the element (short form) - e.g. `tip`, `stable`, `v4.3.2`.  This is required and will fill in ActualVersion from a catalog.
   * **ActualVersion** - Actual Version of the element (long form)
* **Plugins** - a list of plugin configurations to apply.  The list is *Plugin* objects from the system.
* **Prefs** - A partial map of string pref values to apply.
* **Global** - A partial map of the global values to apply.
* **Files** - a list of File Data objects that should be installed on the system.
   * **Path** - the `isos` or `files` path for this file to live at
   * **Sha256Sum** - The SHA-256 sum of the file
   * **Source** - A URL of the file.  This can be a full URL or self://path to indicate that a local file should be used.
   * **Explode** - A boolean flag that if true will cause `bsdtar` to explode the file at the PATH location.
* **Apply** - A boolean field that if true allows the version set to be applied.  This should be true and the **Endpoint** apply field used instead.

