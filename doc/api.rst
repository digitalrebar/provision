.. Copyright (c) 2017 RackN Inc.
.. Licensed under the Apache License, Version 2.0 (the "License");
.. Digital Rebar Provision documentation under Digital Rebar master license
.. index::
  pair: Digital Rebar Provision; API
  pair: Digital Rebar Provision; REST

.. _rs_api:

Digital Rebar Provision API
~~~~~~~~~~~~~~~~~~~~~~~~~~~

In general, the Digital Rebar Provision API is documented via the Swagger spec and introspectable for machines via `/swagger.json` and for humans via `/swagger-ui`.

All API calls are available under `/api/v3` based on the Digital Rebar API convention.

.. _rs_api_filters:

API Filters
-----------

The API includes index driven filters for large deployments that can be used to pre-filter requests from the API.

The list of available indexes is provided by the ``/api/v3/indexes`` and ``/api/v3/indexes/[model]`` calls.  These hashs provide a list of all keys that can be used for filters with some additional metadata.

To use the index, simply include one or more indexs and values on the request URI.  For example:

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

Filtering by Param Value
------------------------

The API includes specialized filter behavior for Params that allows deep searching models for Param values.

To filter Machines or Profiles by Param values, pass the Param name and value using the normal Field filter specification.  When the Field is not found, the backend will search model's Params keys and evalute the filter against the Param value.

Payload Reduction (slim)
------------------------

Since Params and Meta may contain a lot of data, the API supports the ``?slim=[Params,Meta]`` option to allow requests to leave these fields unpopulated.  If this data is needed, operators will have to request the full object or object's Params or Meta in secondary calls.

  ::

    /api/v3/machines?slim=Params,Meta

Only endpoints that offer the ``slim-objects`` feature flag (v3.9+) will accept this flag.

.. _rs_api_notes:

API Exception & Deprecation Notes
---------------------------------

There are times when the API and models have exceptions or changes that do not follow the normal pattern.  This section is designed to provide a reference for those exceptions.

This section is intended to provide general information about and functional of the API.  We maintain this section for legacy operators, when possible avoid using deprecated features!

*Machines.Profile (deprecated by flag: profileless-machine)*
  What would otherwise be Machine.Params is actually embedded under Machines.Profile.Params.
  This deprecated composition simplifies that precedence calculation for Params by making Machines the
  top of the Profiles stack.  All the other fields in the Machines.Profile are ignored.

.. swaggerv2doc:: https://github.com/digitalrebar/provision/releases/download/tip/swagger.json
