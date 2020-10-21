.. Copyright (c) 2017 RackN Inc.
.. Licensed under the Apache License, Version 2.0 (the "License");
.. Digital Rebar Provision documentation under Digital Rebar master license
.. index::
  pair: Digital Rebar Provision; Architecture Design

Architecture Overview
---------------------

Architecturally, *dr-provision* is split into several different packages:

.. glossary::

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
https://godoc.org/github.com/digitalrebar/provision for in-depth
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

