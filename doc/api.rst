.. Copyright (c) 2017 RackN Inc.
.. Licensed under the Apache License, Version 2.0 (the "License");
.. Digital Rebar Provision documentation under Digital Rebar master license
.. index::
  pair: Digital Rebar Provision; API
  pair: Digital Rebar Provision; REST

.. _rs_api:

Digital Rebar Provision API
~~~~~~~~~~~~~~~~~~~~~~~~~~~

In general, the Digital Rebar Provision API is documented via the Swagger spec and introspectable for machines via `/swagger.json` and for humans via `/swagger-ui`.  See :ref:`rs_swagger`.

All API calls are available under `/api/v3` based on the Digital Rebar API convention.

.. _rs_api_filters:

API Filters
-----------

The API includes index driven filters for large deployments that can be used to pre-filter requests from the API.

The list of available indexes is provided by the ``/api/v3/indexes`` and ``/api/v3/indexes/[model]`` calls.  These objects provide a list of all keys that can be used for filters with some additional metadata.

To use the index, simply include one or more indexes and values on the request URI.  For example:

  ::

    /api/v3/machines?Runnable=true&Available=true

The filter specification allows more complex filters using functions:

  * Key=Eq(Value) (that is the same as Key=Value)
  * Lt(Value)
  * Lte(Value)
  * Gt(Value)
  * Gte(Value)
  * Ne(Value)
  * Ranges:
    * Between(Value1,Value2) (edited)
    * Except(Value1,Value2)

The query string applies ALL parameters are to be applied (as implied by the & separator).  All must match to be returned.

.. _rs_api_param_filter:

Filtering by Param Value
------------------------

The API includes specialized filter behavior for Params that allows deep searching models for Param values.

To filter Machines or Profiles by Param values, pass the Param name and value using the normal Field filter specification.  When the Field is not found, the backend will search model's Params keys and evalute the filter against the Param value.

.. _rs_api_proxy:

Leveraging Multi-Site Proxy Forwarding
--------------------------------------

In the :ref:`rs_manager_arch`, API calls to a DRP manager for remote objects are automatically forwarded to the correct attached endpoint.  This allows operators to make API calls to remote endpoints from a centralized manager without knowing the actual owner of the object.  The owning endpoint can be determined for any object by inspecting its `Endpoint` property.

For most cases, no additional information or API notation is required to leverage this functionality. The exception is creating objects on attached endpoints.  The create (aka POST) case requires API callers to specify the target remote endpoint (the endpoint must be registered or the request will be rejected) to the manager.

To make an explicit API proxy call, prefix the path of the request with the endpoint identity.  For example:

  ::

    /[target endpoint id]/api/v3/machines


This pattern is also invoked by with the `-u` flag in the DRPCLI.

NOTE: Multi-site is a licensed feature available in DRP v4.5+ and must be enabled for the endpoint.

.. _rs_api_slim:

Payload Reduction (slim)
------------------------

Since Params and Meta may contain a lot of data, the API supports the ``?slim=[Params,Meta]`` option to allow requests to leave these fields unpopulated.  If this data is needed, operators will have to request the full object or object's Params or Meta in secondary calls.

  ::

    /api/v3/machines?slim=Params,Meta

Only endpoints that offer the ``slim-objects`` feature flag (v3.9+) will accept this flag.


Exploration with Curl
---------------------

You can also interact with the API using ``curl``.  The general pattern is:

  ::

    curl -X <method> -k -u <username>:<password> -H `Content-Type: application/json' -H 'Accept: application/json' https://<endpoint addr>:<port>/api/v3/<opject type>/<object ID>

In the remainder of this section, <object type> refers to the lower case, pluralized version of the type of object.  This is `bootenvs` for boot environments, `workflows` for workflows, `machines` for machines, and so on.
<object id> refers to the unique identifier for this object, which is generally the `Name`, `ID` or `Uuid` field of an object.  You can also use any unique index in this field, in the form of `<index name>:<value>`.
A common one to use is `Name:machine.name` for Machine objects instead of the Uuid.


The API follows the usual REST guidelines:

* HEAD /api/v3/<object type> gets you headers containing basic information about how many of <object type> are present in the system.
  You can use filters on this request.
* GET /api/v3/<object type> lists all of the objects of the requested type.  The result is a JSON array.  You can use filters on this request.
* POST /api/v3/<object type> is a request to create a new object.  The body of the payload should be valid JSON for the object type.
* GET /api/v3/<object type>/<object id> fetches the request object.
* HEAD /spi/v3/<object type>/<object id> tests to see if the requested object exists.

.. _rs_api_notes:

API Exception & Deprecation Notes
---------------------------------

There are times when the API and models have exceptions or changes that do not follow the normal pattern.  This section is designed to provide a reference for those exceptions.

This section is intended to provide general information about and functional of the API.  We maintain this section for legacy operators, when possible avoid using deprecated features!

*Machines.Profile (deprecated by flag: profileless-machine)*
  What would otherwise be Machine.Params is actually embedded under Machines.Profile.Params.
  This deprecated composition simplifies that precedence calculation for Params by making Machines the
  top of the Profiles stack.  All the other fields in the Machines.Profile are ignored.

