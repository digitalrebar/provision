.. Copyright (c) 2017 RackN Inc.
.. Licensed under the Apache License, Version 2.0 (the "License");
.. Digital Rebar Provision documentation under Digital Rebar master license
.. index::
  pair: Digital Rebar Provision; Architecture Design

Architecture Overview
---------------------

Architecturally, *dr-provision* is split into several different packages:

.. glossary::

  :ref:`rs_arch_backend`
    is responsible for making sure that all the data is valid
    and gets written to the persistent store whenever things get
    updated, along with storing any non-persistent runtime data we need
    to keep track of.

  :ref:`rs_arch_midlayer`
    is where the TFTP, static HTTP, and DHCP
    services live, along with the plugin management code.

  :ref:`rs_arch_frontend`
    is responsible for providing the REST API.

  :ref:`rs_arch_models`
    define the data models that the other packages use, along
    with some common functionality that can be shared between the client
    and server side.

  :ref:`rs_arch_api`
    defines a client-side Go API for interacting with *dr-provision*.

  :ref:`rs_arch_cli`
    provides our default CLI for interacting with *dr-provision*.

  :ref:`rs_arch_plugin`
    implements the core client code that all plugins should use
    to act as a *dr-provision* plugin.


.. _rs_arch_backend:

backend
~~~~~~~

In Memory Database
^^^^^^^^^^^^^^^^^^

All operating data lives in memory all the time.  The only time
*dr-provision* reads information from persistent storage is when it is
starting up, otherwise we treat the persistent store as write-only.
The only exception to this design principle is streaming log data from
job execution.  We may revisit this design principle if memory
pressure becomes a real constraint.

The backend provides a DataTracker that is responsible for holding all
of the persistible data.  DataTracker also implements indexing
mechanism that the other layers use to ensure that the other layers
can quickly find what they need, along with a locking scheme to ensure
that the data stays consistent.  The backend also provides a
RequestTracker to ensure that logging and locking on a per-request
basis is consistent.

Pluggable Storage
^^^^^^^^^^^^^^^^^

The mechanism for storing persistent data should be pluggable, and
only rely on basic key/value store semantics in the absence of
transactions.  The backend relies on an external Go package
(https://github.com/digitalrebar/provision/store) to abstract some basic
behaviour on top of various key-value type stores.  Adding support for
new storage types will be a matter of adding them to that package, not
to dr-provison itself.

Single Point of Validation
^^^^^^^^^^^^^^^^^^^^^^^^^^

To the extent that it is feasible, all object validation happenS in
the backend.  Since the backend is responsible for writing data to the
persistent store, it is also the best place to implement all data
validation.

Static FS with Dynamic Overlay
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

A large part of what *dr-provision* does boils down to rendering
templates and making them available in the right place at the right
time.  Whenever a machine, bootenv, or stage changes, there are
generally templates that have to be rerendered and made available via
static HTTP and TFTP to ensure that machines boot to the corrent
environment over the network, have the right OS installation
templates, load the proper credentials, etc.  We already serve static
HTTP and TFTP content from a user-configurable location to provide
basic files needed to PXE boot a system and provide all the packages
and files needed to install an OS.

For dynamic content, however, we don't always want to write files to a
filesystem where anyone with a web browser can discover them, and
where we have to worry about cleaning up dynamic content whenever
something changes.  To that end, we provide a static FS implementation
that can register for callbacks to be involed whenever someone
accesses a file via TFTP or HTTP.  That allows us to defer template
rendering until the files are actually requested, and it allows us to
transparently proxy requests to remote repositories when we don't
actually have a file tree present locally.

Dynamic Remote IP to Local IP Caching
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

*dr-provision* is designed to work seamlessly on a multi-homed system
and to deal with complex local network configurations.  To that end,
every other subsystem that listens for packets is instrumented to
capture the IP where the request originated and the local IP address
that the request came in on.  Template rendering and DHCP request
handling use this information to make sure that we supply the best IP
address that a client should use to contact *dr-provision* at for any
future communication.

.. _rs_arch_midlayer:

midlayer
~~~~~~~~

The *midlayer* package handles some basic services that DRP provides as well as
the content package management system.

TFTP Service
^^^^^^^^^^^^

The TFTP service ties together the `pin
<https://github.com/pin/tftp>`_ TFTP handling package and the static
FS that the backend provides to handle TFTP requests.  We only allow
clients to get files, uploading them is not allowed.  Remote and local
IP addresses for each connection are cached

Static HTTP Service
^^^^^^^^^^^^^^^^^^^

The Static HTTP service implements a simple high-performance HTTP
server that serves files using the static FS that the backend
provides.  Remove and local IP addresses for each connection are
cached.


DHCP Service
^^^^^^^^^^^^

The DHCP service built in to *dr-provision* is designed to be fully API
driven and to provide all the features needed to manage system IP
address assignments through the complete provisioning lifecycle. As
such, it has a few interesting features that other DHCP servers may
not have:

- The ability to have different ways of determining what unique
  attribute in a DHCP packet to use to allocate an IP address.  When
  you see references to Strategy and Token in the DHCP models,
  Strategy refers to the unique attribute the DHCP server should use,
  and Token refers to the value that the Stategy picked.

  For now, the only implemented Strategy is MAC, which has the DHCP
  server use the MAC address of the network adaptor of the network
  interface as the unique value of the Token.

- The DHCP server is fully API driven.  You can add, remove, and
  modify Reservations and Subnets on the fly, and changes take effect
  immediately.

- Built-in ProxyDHCP support, on a subnet by subnet basis.
  *dr-provision* can coexist with other DHCP servers to only provide PXE
  support for specific address ranges, leaving address management to
  your preexisting DHCP infrastructure.

Plugin Management
^^^^^^^^^^^^^^^^^

*dr-provision* can add extended functionality via external plugins.  The
midlayer implements all of the functionality needed to accept plugin
uploads, interrogate them to discover what functionality they
implement, import any content built in to the plugin, and hand off
requests and events to the plugin for further processing.

.. _rs_plugin_license_events:

Plugin License Events
=====================

When plugins are loaded, they will validate their licenses and fail to load or generate events.  You may see
these events as part of that validation process.

A plugin that determines a license is hard expired will generate an event:

* Type - plugins
* Action - failure
* Key - Name of Plugin
* Object - A data structure.

The object data structure has four fields:

* Type - "license expired (hard)"
* CurrentDate - The current date
* SoftExpireDate - The soft expire date
* HardExpireDate - The hard expire date

The plugin will then exit.

A plugin that detects an exceeded value in the license, e.g. too machines, will generate an event:

* Type - plugins
* Action - exceeded
* Key - Name of Plugin
* Object - A data structure

The object data structure has three fields:

* Type - what was exceeded, e.g. machines
* Current - integer count of current objects
* Expected - integer count of expected objects

The plugin will then exit.

A plugin that determines a license is soft expired will generate an event:

* Type - plugins
* Action - failure
* Key - Name of Plugin
* Object - A data structure.

The object data structure has four fields:

* Type - "license expired (soft)"
* CurrentDate - The current date
* SoftExpireDate - The soft expire date
* HardExpireDate - The hard expire date

The plugin will continue to operate.

.. _rs_arch_content:

Content Package Management
^^^^^^^^^^^^^^^^^^^^^^^^^^

The *Content Package Management* system builds a stack of content layers
that are provided to the :ref:`rs_arch_backend` to provide objects to the rest
of the system.  The data stack has the following layers used in this order:

.. csv-table:: Definitions
   :header: "Heading", "Definition"
   :widths: 20, 80

   "Layer Type", "Type of layer in the data stack as reported in the content layer meta data"
   "Overwritable", "Can layers above overwrite content packages at this layer."
   "Can Override", "Can a content package at this layer override lower layers."
   "Writable", "Can the system receive written objects"
   "Many", "Can multiple content packages be added to this layer"
   "Use", "Who provideds and its use"

.. csv-table:: Content Package Management
   :header: "Layer", "Overwritable", "Can Override", "Writable", "Many", "Use"
   :widths: 20, 10, 10, 10, 10, 50

   "writable", "yes", "no",  "yes", "no",  "Persistent layer"
   "local",    "yes", "yes", "no",  "no",  "Layer providing content from local filesystem, /etc/dr-provision directory"
   "dynamic",  "no",  "yes", "no",  "yes", "Layer providing dynamic content packages provided by the API"
   "default",  "yes", "yes", "no",  "no",  "Layer providing default content that is always present, but replaceable."
   "plugin",   "no",  "yes", "no",  "yes", "Layer providing plugin provided content packages."
   "basic",    "yes", "yes", "no",  "no",  "Layer providing mandatory DRP model objects."

When an object is looked up, the look up code will start walking down the stack until the object is
found and it will be returned.  When an object is to be updated or created, the *Writable* aspect of
the layer will be checked to see if the object can be updated or created.  If the object can be
stored in a layer, it will be used.  The content layer stack places the wriable store at the top of
the stack.

The simplified view of the stack from the API can be boiled down to:

* Create - Created object's key must not exist in the stack.
* Read - Object will be searched from the top down until it is found.
* Update - Updated object must exist only in the writable layer.
* Delete - Deleted Object must exist only in the writable layer.

.. _rs_arch_frontend:

frontend
~~~~~~~~

The DRP frontend implements a REST + JSON API for others to interact
with and manage *dr-provision*.  The *dr-provision* API is available via
HTTPS, and we will upgrade to HTTP v2 opportunistically.

Threaded Logging
^^^^^^^^^^^^^^^^

Each individual request to the API is logged using a unique ID, and
that ID is threaded through to all the code paths that the request
affects.  Detailed logging along with an arbitrary token can also be
enabled on a per-request basis to aid in debugging and audit purposes.

Basic and JWT Token Authentication
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

You can authenticate to the *dr-provision* API via basic auth and via
time-limited JWT tokens.  We also provide means to invalidate tokens
globally and on a per-user basis.


Websocket-based Event Delivery
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

Authenticated users can open a websocket and arrange for a variety of
different events to be watched for.  This eliminates the need to poll
in a loop for a wide variety of different situations.

.. _rs_arch_models:

models
~~~~~~

Every valid *dr-provision* object has a Model that is implemented in
this package.  These models are authoritative, and their JSON
serialization in Go is the canonical wire format.

.. _rs_arch_api:

api
~~~

The API package implements the reference Go client API for
*dr-provision*. You should consult the go docs for the API at
https://godoc.org/github.com/digitalrebar/provision/api for in-depth
discussion on how to use the client API.

.. _rs_arch_cli:

cli
~~~

The CLI package implements the reference Go client CLI for
*dr-provision*.  The main program for *drpcli* includes this
set of functions.

.. _rs_arch_plugin:

plugin
~~~~~~

The plugin package implements the Go core functions needed to create
a *dr-provision* plugin.

