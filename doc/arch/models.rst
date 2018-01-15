.. Copyright (c) 2017 RackN Inc.
.. Licensed under the Apache License, Version 2.0 (the "License");
.. Digital Rebar Provision documentation under Digital Rebar master license
.. index::
  pair: Digital Rebar Provision; Data Models

.. _rs_data_models:

Data Models
===========

Digital Rebar Provision uses several different data models to manage
the task of discovering and provisioning machines in a data center.

Common Model Fields
<<<<<<<<<<<<<<<<<<<

.. _rs_data_metadata:

Model Metadata
--------------

Virtually every model contains some embedded metadata, available under
the **Meta** field.  This data, which consists of a map of string ->
string pairs, is ignored by the dr-provision server.  The most common
metadata fields you will see are:

icon
  The icon that the UX will use to display instances of this model.
  Users can choose icons from http://fontawesome.io/icons/.

color
  The color the icon will be displayed as

title
  The full name that the UX will use.

feature-flags
  A comma-seperated list of strings that indicate which
  features are available for a model. Provision uses feature
  flags to help the various components and content layers to
  converge on a supported set of avaiable features.


.. _rs_data_validation:

Model Validation
----------------

Models also contain common fields that track the validity and
availability of individual objects.  These fields are:

- Validated: a boolean value that indicates whether a given object is
  semantically valid or not.  Semantically invalid objects will never
  be saved, and if one is returned via the API the Errors field will
  be populated with a list of messages indicating what is invalid.
- Available: a boolean value that indicates whether the object is
  available to be used, not whether it is semantically valid -- an
  object that is invaild can never be available, while an object that
  is not available can be semantically valid.
- Errors: a list of strings that contain any error messages that
  occurred in the process of checking whether a given object is valid
  and available.  Error messages are designed to be human readable.

Objects are checked for validity and availability on initial startup
of dr-provision (when they are all loaded into memory), and thereafter
every time they are updated.  You must check each object returned from
an API interaction to ensure that it is valid and available before
using it.

Other Common Model Fields
-------------------------

Models can contain other common fields that may be present for user
edification and API tracking purposes, but that do not affect how
dr-provision will use or interpret changes to the objects.  These
extra fields are:

- ReadOnly: a boolean value that indicates whether the object can be
  modified via the API. This field is set to True if the object was
  loaded from a read-only content layer.

- Description: A brief description of what the object is for, how it
  should be used, etc.  Descriptions should be one line long.

- Documentation: A longer description of what the object is for and
  how it should be used, generally a few lines to a few paragraphs
  long.  For now, only Params and Tasks have a Documentation field,
  but other models may add them as situations demand.
